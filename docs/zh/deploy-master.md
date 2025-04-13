# Master 部署

Master 推荐使用 docker 部署！不推荐直接安装到服务器中

会给出三种部署方式，任选一种即可

部署后没有默认用户，注册的第一个用户即为管理员，为了安全，默认不开启多用户注册

## 前期准备

### 服务器开放公网端口：

- **WEBUI 端口**: 默认 `TCP 9000`
- **RPC 端口**: 默认 `TCP 9001`
- **frps 的API端口**：没有默认，请随意预留，例子使用 `TCP/UDP 7000`
- **frps 对外开放的服务端口**：没有默认，请随意预留，例子使用 `TCP/UDP 26999-27050`

如果使用反向代理，请忽略 WEBUI 和 RPC 端口，放行 80/443 即可

WEBUI 端口也可以处理 h2c 格式的 RPC 连接

RPC 端口也可以处理自签名 HTTPS 的 API 连接

二者都可使用反向代理服务器连接并提供TLS

> 测试端口是否开放的方法（以8080为例），在服务器上运行：
> ```shell
> python3 -m http.server 8080
> ```
> 然后在另一台电脑/服务器上执行：
> ```shell
> curl http://服务器公网IP/域名:8080 -I
> ```
> 成功的话，输出类似
> ```
> HTTP/1.0 200 OK
> Server: SimpleHTTP/0.6 Python/3.11.0
> Date: Sat, 12 Apr 2025 17:12:15 GMT
> Content-type: text/html; charset=utf-8
> Content-Length: 8225
> ```

## 方式一：Docker Compose 部署

服务器需要安装docker和docker compose

首先创建一个`docker-compose.yaml`文件，写入以下内容

```yaml
version: "3"

services:
  frpp-master:
    image: vaalacat/frp-panel:latest
    environment:
      APP_GLOBAL_SECRET: your_secret 
      MASTER_RPC_HOST: 1.2.3.4 #服务器的外部IP或域名
      MASTER_RPC_PORT: 9001
      MASTER_API_HOST: 1.2.3.4 #服务器的外部IP或域名
      MASTER_API_PORT: 9000
      MASTER_API_SCHEME: http
    volumes:
      - ./data:/data
    restart: unless-stopped
    command: master
```

## 方式二：Docker 命令部署

服务器需要安装 docker，我们推荐使用 host 网络模式部署 `Master`

```bash
# 推荐
# MASTER_RPC_HOST要改成你服务器的外部IP
# APP_GLOBAL_SECRET注意不要泄漏，客户端和服务端的是通过Master生成的
docker run -d \
	--network=host \
	--restart=unless-stopped \
	-v /opt/frp-panel:/data \
	-e APP_GLOBAL_SECRET=your_secret \
	-e MASTER_RPC_HOST=0.0.0.0 \
	vaalacat/frp-panel
```

如果你不想使用 host 网络模式，请参考使用下面的命令修改

```bash
# 或者
# 运行时记得删除命令中的中文
docker run -d -p 9000:9000 \ # API控制台端口
	-p 9001:9001 \ # rpc端口
	-p 7000:7000 \ # frps 端口
	-p 27000-27050:27000-27050 \ # 给frps预留的端口
	--restart=unless-stopped \
	-v /opt/frp-panel:/data \ # 数据存储位置
	-e APP_GLOBAL_SECRET=your_secret \ # Master的secret注意不要泄漏，客户端和服务端的是通过Master生成的
	-e MASTER_RPC_HOST=0.0.0.0 \ # 这里要改成你服务器的外部IP
	vaalacat/frp-panel
```

## 方式三：使用 docker 反向代理 TLS 加密部署

这里我们以 [Traefik](https://traefik.io/traefik/) 为例

> `Traefik` 可以实时自动识别 Docker 容器的端口并热更新配置，非常适合 Docker 服务的反向代理

首先创建一个名为`traefik`的反向代理专用网络
```bash
docker network create traefik
```
然后启动反向代理和 Master 服务
- `docker-compose.yaml`

```yaml
version: '3'

services:
  traefk-reverse-proxy:
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
      - --api.insecure # 在生产环境请删除这一行
	  # 这下面使用 80 端口做ACME HTTP DNS证书验证
      - --certificatesresolvers.le.acme.email=me@example.com
      - --certificatesresolvers.le.acme.storage=/etc/traefik/conf/acme.json
      - --certificatesresolvers.le.acme.httpchallenge=true
    ports:
      # 反向代理的 HTTP 端口
      - "80:80"
	  # 反向代理的 HTTPS 端口
	  - "443:443"
      # Traefik 的 Web UI (--api.insecure=true 会使用这个端口)
	  # 生产环境请删除这个端口
      - "8080:8080"
    volumes:
      # 挂载 docker.sock，这样 Traefik 可以自动识别主机上所有 docker 容器反向代理配置
      - /var/run/docker.sock:/var/run/docker.sock
	  # 保存 Traefik 申请的证书
	  - ./conf:/etc/traefik/conf

  frpp-master:
    image: vaalacat/frp-panel:latest # 这里换成你想使用的版本
    environment:
      APP_GLOBAL_SECRET: your_secret
	# 因为 api 和 rpc 使用的协议不一样
	# 我们需要对 api 和 rpc 使用两个域名
	# 以便反向代理正确识别需要转发的协议
      MASTER_RPC_HOST: frpp.example.com
      MASTER_API_PORT: 443
      MASTER_API_HOST: frpp-rpc.example.com
      MASTER_API_SCHEME: https
    networks:
      - traefik
    volumes:
      - ./data:/data
    ports:
	  # 无需为 master 预留 api 和 rpc 端口
	  # 预留frps api端口
      - 7000:7000
      - 7000:7000/udp
	  # 预留frps的业务端口
	  # 26999 端口是留给 frps 的http代理端口
      - 26999-27050:26999-27050
      - 26999-27050:26999-27050/udp
    restart: unless-stopped
    command: master
    labels:
	  # API
      - traefik.http.routers.frp-panel-api.rule=Host(`frpp.example.com`)
      - traefik.http.routers.frp-panel-api.tls=true
      - traefik.http.routers.frp-panel-api.tls.certresolver=le
      - traefik.http.routers.frp-panel-api.entrypoints=websecure
      - traefik.http.routers.frp-panel-api.service=frp-panel-api
      - traefik.http.services.frp-panel-api.loadbalancer.server.port=9000
      - traefik.http.services.frp-panel-api.loadbalancer.server.scheme=http
	  # RPC
      - traefik.http.routers.frp-panel-rpc.rule=Host(`frpp-rpc.example.com`)
      - traefik.http.routers.frp-panel-rpc.tls=true
      - traefik.http.routers.frp-panel-rpc.tls.certresolver=le
      - traefik.http.routers.frp-panel-rpc.entrypoints=websecure
      - traefik.http.routers.frp-panel-rpc.service=frp-panel-rpc
      - traefik.http.services.frp-panel-rpc.loadbalancer.server.port=9000
      - traefik.http.services.frp-panel-rpc.loadbalancer.server.scheme=h2c
      # 下方如果你用不到 frps 的http代理，可以不要
	  # 需要配置域名 *.frpp.example.com 泛解析到你服务器的公网IP
	  # 这样可以实现使用 .frpp.example.com 结束的域名，在 443 端口，转发多个服务到多个 frpc
      - traefik.http.routers.frp-panel-tunnel.rule=HostRegexp(`.*.frpp.example.com`)
      - traefik.http.routers.frp-panel-tunnel.tls.domains[0].sans=*.frpp.example.com
      - traefik.http.routers.frp-panel-tunnel.tls=true
      - traefik.http.routers.frp-panel-tunnel.tls.certresolver=le
      - traefik.http.routers.frp-panel-tunnel.entrypoints=websecure
      - traefik.http.routers.frp-panel-tunnel.service=frp-panel-tunnel
      - traefik.http.services.frp-panel-tunnel.loadbalancer.server.port=26999
      - traefik.http.services.frp-panel-tunnel.loadbalancer.server.scheme=http
networks:
  traefik:
    external: true
    name: traefik
```

上方的 `docker-compose.yaml` 部署完成后，可以访问 `服务器公网IP/域名:8080` 查看反向代理状态

随后配置 default server 即可实现 frp 子域名转发：

| 配置项 | 值 |
|----|-----|
|	FRPs 监听端口	|	7000	|
|	FRPs 监听地址	|	0.0.0.0	|
|	代理监听地址	|	0.0.0.0	|
| 	HTTP 监听端口	|	26999	|
|	域名后缀		|	frpp.example.com	|
