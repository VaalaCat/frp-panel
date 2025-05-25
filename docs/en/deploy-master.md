# Master Deployment

We recommend deploying the Master via Docker! Direct installation on the host is not recommended.

You have three deployment options—choose one.

After deployment, there is no default user. The first registered user becomes the administrator. For security, multi-user registration is disabled by default.

By default, the program stores data in its working directory. To change this, see the configuration reference.

**Important:** If you want to deploy the Master and have it also act as a Server, remember to configure the `default` server in the Web UI after starting the Master.

## Prerequisites

Open the following ports on your server:

- **WEB UI port**: TCP 9000  
- **RPC port**: TCP 9001  
- **frps API port**: any free port (example uses TCP/UDP 7000)  
- **frps service ports**: any port range (example uses TCP/UDP 26999–27050)  

If you use a reverse proxy, you can ignore WEB UI and RPC ports—just open 80/443.  
- The WEB UI port can also accept h2c RPC connections.  
- The RPC port can also accept self-signed HTTPS API connections.  
- Both can be fronted by a TLS-terminating reverse proxy.  

To secure communication, set the environment variables `CLIENT_RPC_URL` and `CLIENT_API_URL`. First deploy normally, then adjust these variables.

![](../public/images/frp-panel-platform-connection-env.svg)

> To test if a port (e.g. 8080) is open, run on the server:  
> ```shell
> python3 -m http.server 8080
> ```  
> Then from another host:  
> ```shell
> curl http://SERVER_IP:8080 -I
> ```  
> A successful response looks like:  
> ```
> HTTP/1.0 200 OK
> Server: SimpleHTTP/0.6 Python/3.11.0
> Date: Sat, 12 Apr 2025 17:12:15 GMT
> Content-type: text/html; charset=utf-8
> Content-Length: 8225
> ```

---

## Deploying on Linux

### Option 1: Docker Compose

Install Docker and Docker Compose, then create `docker-compose.yaml`:

```yaml
version: "3"

services:
  frpp-master:
    image: vaalacat/frp-panel:latest
    network_mode: host
    environment:
      APP_GLOBAL_SECRET: your_secret
      MASTER_RPC_HOST: 1.2.3.4       # external IP or domain
      MASTER_RPC_PORT: 9001
      MASTER_API_HOST: 1.2.3.4       # external IP or domain
      MASTER_API_PORT: 9000
      MASTER_API_SCHEME: http
    volumes:
      - ./data:/data               # data directory
    restart: unless-stopped
    command: master
```

### Option 2: Docker CLI

Install Docker. We recommend `host` network mode:

```bash
docker run -d \
  --network=host \
  --restart=unless-stopped \
  -v /opt/frp-panel:/data \
  -e APP_GLOBAL_SECRET=your_secret \
  -e MASTER_RPC_HOST=0.0.0.0 \
  vaalacat/frp-panel
```

If you cannot use `host` network mode:

```bash
docker run -d \
  -p 9000:9000 \        # API
  -p 9001:9001 \        # RPC
  -p 7000:7000 \        # frps API
  -p 27000-27050:27000-27050 \  # frps service ports
  --restart=unless-stopped \
  -v /opt/frp-panel:/data \
  -e APP_GLOBAL_SECRET=your_secret \
  -e MASTER_RPC_HOST=0.0.0.0 \
  vaalacat/frp-panel
```

### Option 3: Docker + Reverse-Proxy TLS (Traefik Example)

Create a Docker network for Traefik:

```bash
docker network create traefik
```

Create `docker-compose.yaml`:

```yaml
version: '3'

services:
  traefik:
    image: traefik:v3.3
    restart: unless-stopped
    networks:
      - traefik
    command:
      - --entryPoints.web.address=:80
      - --entryPoints.websecure.address=:443
      - --entryPoints.websecure.http2.maxConcurrentStreams=250
      - --providers.docker
      - --providers.docker.network=traefik
      - --certificatesresolvers.le.acme.email=me@example.com
      - --certificatesresolvers.le.acme.storage=/etc/traefik/conf/acme.json
      - --certificatesresolvers.le.acme.httpchallenge=true
    ports:
      - "80:80"
      - "443:443"
      - "8080:8080"        # Traefik dashboard (remove in production)
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - ./conf:/etc/traefik/conf

  frpp-master:
    image: vaalacat/frp-panel:latest
    networks:
      - traefik
    volumes:
      - ./data:/data
    restart: unless-stopped
    command: master
    environment:
      APP_GLOBAL_SECRET: your_secret
      MASTER_RPC_HOST: frpp-rpc.example.com
      MASTER_API_HOST: frpp.example.com
      MASTER_API_PORT: 443
      MASTER_API_SCHEME: https
    ports:
      - 7000:7000
      - 7000:7000/udp
      - 26999-27050:26999-27050
      - 26999-27050:26999-27050/udp
    labels:
      # API router
      - traefik.http.routers.frp-panel-api.rule=Host(`frpp.example.com`)
      - traefik.http.routers.frp-panel-api.tls=true
      - traefik.http.routers.frp-panel-api.tls.certresolver=le
      - traefik.http.routers.frp-panel-api.entrypoints=websecure
      - traefik.http.services.frp-panel-api.loadbalancer.server.port=9000
      - traefik.http.services.frp-panel-api.loadbalancer.server.scheme=http

      # RPC router
      - traefik.http.routers.frp-panel-rpc.rule=Host(`frpp-rpc.example.com`)
      - traefik.http.routers.frp-panel-rpc.tls=true
      - traefik.http.routers.frp-panel-rpc.tls.certresolver=le
      - traefik.http.routers.frp-panel-rpc.entrypoints=websecure
      - traefik.http.services.frp-panel-rpc.loadbalancer.server.port=9000
      - traefik.http.services.frp-panel-rpc.loadbalancer.server.scheme=h2c

      # Tunnel router (optional HTTP proxy for frpc)
      - traefik.http.routers.frp-panel-tunnel.rule=HostRegexp(`.*.frpp.example.com`)
      - traefik.http.routers.frp-panel-tunnel.tls.domains[0].sans=*.frpp.example.com
      - traefik.http.routers.frp-panel-tunnel.tls=true
      - traefik.http.routers.frp-panel-tunnel.tls.certresolver=le
      - traefik.http.routers.frp-panel-tunnel.entrypoints=websecure
      - traefik.http.services.frp-panel-tunnel.loadbalancer.server.port=26999
      - traefik.http.services.frp-panel-tunnel.loadbalancer.server.scheme=http

networks:
  traefik:
    external: true
    name: traefik
```

After starting, visit `SERVER_IP:8080` to view Traefik’s dashboard.

Then configure the `default` server in the Master Web UI:

| Setting               | Value                  |
|-----------------------|------------------------|
| FRPs listen port      | 7000                   |
| FRPs listen address   | 0.0.0.0                |
| Proxy listen address  | 0.0.0.0                |
| HTTP listen port      | 26999                  |
| Domain suffix         | frpp.example.com       |

---

## Deploying on Windows

### Direct Execution

In the folder containing the executable, create a `.env` file (no extension) with:

```
APP_GLOBAL_SECRET=your_secret
MASTER_RPC_HOST=IP
DB_DSN=data.db
```

Then run:

```
frp-panel-amd64.exe master
```