# Client Deployment

Client is recommended to be deployed using docker, directly on the client machine. You can execute commands with root privileges directly on the client machine through a remote terminal, which makes upgrading and management convenient.

## Preparation

Open the Master's webui and log in. If you don't have an account, register directly - the first user will be the administrator.

Navigate to `Clients` in the sidebar, click `New` at the top and enter a unique identifier for the client, then click save.

![](../public/images/en_client_list.png)

After refreshing, the new client will appear in the list.

Before deploying the client, you need to modify the server configuration, otherwise the client cannot connect properly.

## Direct Deployment

Click on the `ID (click to view installation command)` column of the corresponding client. A popup will show installation commands for different systems. Copy the appropriate command to the corresponding terminal to install. Here's an example for Linux:

```
curl -fSL https://raw.githubusercontent.com/VaalaCat/frp-panel/main/install.sh | bash -s --  client -s abc -i user.s.client1 -a 123123 -r frpp-rpc.example.com -c 9001 -p 9000 -e http
```

If you're in China, you can add a GitHub accelerator to the installation script:

```
curl -fSL https://ghfast.top/https://raw.githubusercontent.com/VaalaCat/frp-panel/main/install.sh | bash -s --  client -s abc -i user.s.client1 -a 123123 -r frpp-rpc.example.com -c 9001 -p 9000 -e http
```

Note, if you use a reverse proxy with TLS, you need to modify this command to something like:

```bash
curl -fSL https://ghfast.top/https://raw.githubusercontent.com/VaalaCat/frp-panel/main/install.sh | bash -s --  frp-panel client -s abc -i user.s.client1 -a 123123 -t frpp.example.com -r frpp-rpc.example.com -c 443 -p 443 -e https
```

## Docker Compose Deployment

Click on the hidden field in the `Secret (click to view startup command)` column of the corresponding client, and copy a startup command similar to the following for later use:

```bash
frp-panel client -s abc -i user.s.client1 -a 123123 -r frpp-rpc.example.com -c 9001 -p 9000 -e http
```

Note, if you use a reverse proxy with TLS, you need to modify this command to something like:

```bash
frp-panel client -s abc -i user.s.client1 -a 123123 -t frpp.example.com -r frpp-rpc.example.com -c 443 -p 443 -e https
```

docker-compose.yaml

```yaml
version: '3'
services:
  frp-panel-client:
    image: vaalacat/frp-panel
    container_name: frp-panel-client
    network_mode: host
    restart: unless-stopped
    command: client -s abc -i user.s.client1 -a 123123 -t frpp.example.com -r frpp-rpc.example.com -c 443 -p 443 -e https
```
