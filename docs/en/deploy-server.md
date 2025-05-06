# Server Deployment

Server is recommended to be deployed using docker! Direct installation on the server is not recommended.

> If you only have one public network server to manage, you can use the `default server` that comes with `master`, no need to deploy a separate `server`.

## 1. Preparation

Open the Master's webui and log in. If you don't have an account, register directly - the first user will be the administrator.

Navigate to `Servers` in the sidebar, click `New` at the top and enter a unique identifier for the server and the IP/domain name that can be accessed from the public network, then click save.

![](../public/images/en_server_list.png)

After refreshing, the new server will appear in the list. Click on the hidden field in the `Secret (click to view startup command)` column of the corresponding server, and copy a startup command similar to the following for later use:

```bash
frp-panel server -s abc -i user.s.server1 -a 123123 -r frpp-rpc.example.com -c 9001 -p 9000 -e http
```

Note, if you use a reverse proxy with TLS, you need to modify this command to something like:

```bash
frp-panel server -s abc -i user.s.server1 -a 123123 -t frpp.example.com -r frpp-rpc.example.com -c 443 -p 443 -e https
```

## 2. Program Installation

### Docker Compose Deployment

docker-compose.yaml

```yaml
version: '3'
services:
  frp-panel-server:
    image: vaalacat/frp-panel
    container_name: frp-panel-server
    network_mode: host
    restart: unless-stopped
    command: server -s abc -i user.s.server1 -a 123123 -t frpp.example.com -r frpp-rpc.example.com -c 443 -p 443 -e https
```

### Direct Deployment

If you want to deploy directly, please refer to the client deployment steps.

## 3. Server Configuration

After installation, you need to modify the server configuration according to your network and requirements, otherwise the client cannot connect properly.