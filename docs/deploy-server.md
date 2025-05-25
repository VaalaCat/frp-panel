# Server 部署

Server 推荐使用 docker 部署！不推荐直接安装到服务器中

注意 ⚠️：client 和 server 的启动指令可能会随着项目更新而改变，虽然在项目迭代时会注意前后兼容，但仍难以完全适配，因此 client 和 server 的启动指令以 master 生成为准

> 如果只有一台公网服务器需要管理，那么使用 `master` 自带的 `default` 服务端即可，无需单独部署 `server`，但要注意在 `master` 启动后要配置 `default` 服务端

## 在 Linux 上部署

### 1. 准备

打开 Master 的 webui 并登录，如果没有账号，请直接注册，第一个用户即为管理员

在侧边栏跳转到 `服务端`，点击上方的 `新建` 并输入 服务端 的唯一识别ID和 服务端 能够被公网访问的 IP/域名，然后点击保存

![](./public/images/cn_server_list.png)

刷新后，新的服务端会出现在列表中。点击对应服务端的`密钥 (点击查看启动命令)`一列中的隐藏字段，复制类似的启动命令如下备用：

```bash
frp-panel server -s abc -i user.s.server1 --api-url http://frpp.example.com:9000 --rpc-url grpc://frpp-rpc.example.com:9001
```

注意，如果你使用 反向代理 TLS，需要以 http 上游的形式，外部 443 端口代理 `master` 的 9000(API) 端口，且修改启动/安装命令类似如下：

```bash
frp-panel server -s abc -i user.s.server1 --api-url https://frpp.example.com:443 --rpc-url wss://frpp.example.com:443
```

### 2. 程序安装

#### Docker Compose 部署

docker-compose.yaml

```yaml
version: '3'
services:
  frp-panel-server:
    image: vaalacat/frp-panel
    container_name: frp-panel-server
    network_mode: host
    restart: unless-stopped
    command: server -s abc -i user.s.server1 --api-url http://frpp.example.com:9000 --rpc-url grpc://frpp-rpc.example.com:9001
```

#### 直接运行

如果你想要直接运行，不使用服务管理工具，请参考 client 直接运行的步骤

#### 安装为 systemd 服务

请参考 client 部署 systemd 的步骤

### 3. 服务端配置

安装完后需要按你的网络和需求，修改服务端的配置，否则客户端无法正常连接

## 在 Windows 上部署

请参考 client 部署的步骤
