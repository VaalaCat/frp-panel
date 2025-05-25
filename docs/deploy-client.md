# Client 部署

Client 推荐使用 docker 部署

但直接部署在客户机中，可以通过远程终端直接在客户机以 root 权限执行命令，方便升级和管理。

注意 ⚠️：client 和 server 的启动指令可能会随着项目更新而改变，虽然在项目迭代时会注意前后兼容，但仍难以完全适配，因此 client 和 server 的启动指令以 master 生成为准

## 准备

打开 Master 的 webui 并登录，如果没有账号，请直接注册，第一个用户即为管理员

在侧边栏跳转到 `服务端`，点击上方的 `新建` 并输入 客户端 的唯一识别ID然后点击保存

![](./public/images/cn_client_list.png)

刷新后，新的客户端会出现在列表中。

部署client之前，需要修改服务端的配置，否则客户端无法正常连接

## 在 Linux 上部署

### 直接运行

首先在系统上创建一个专用目录

点击对应客户端的 `ID (点击查看安装命令)` 一列，弹出不同系统的安装命令，粘贴到对应终端即可安装，这里以 Linux 为例

```
curl -fSL https://raw.githubusercontent.com/VaalaCat/frp-panel/main/install.sh | bash -s --  client -s abc -i user.c.client1 --api-url http://frpp.example.com:9000 --rpc-url grpc://frpp-rpc.example.com:9001
```

如果你在国内，可以在WebUI中配置增加github加速到安装脚本前，以ghfast为例，配置后复制的内容可能类似下方：

```bash
curl -fSL https://ghfast.top/https://raw.githubusercontent.com/VaalaCat/frp-panel/main/install.sh | bash -s -- --github-proxy https://ghfast.top/  client -s abc -i user.c.client1 --api-url http://frpp.example.com:9000 --rpc-url grpc://frpp-rpc.example.com:9001
```

注意，如果你使用 反向代理 TLS，需要以 http 上游的形式，外部 443 端口代理 `master` 的 9000(API) 端口，且修改启动/安装命令类似如下：

```bash
curl -fSL https://ghfast.top/https://raw.githubusercontent.com/VaalaCat/frp-panel/main/install.sh | bash -s -- --github-proxy https://ghfast.top/  client -s abc -i user.c.client1 --api-url https://frpp.example.com:443 --rpc-url wss://frpp.example.com:443
```

### Docker Compose 部署

点击对应客户端的 `密钥 (点击查看启动命令)` 一列中的隐藏字段，复制类似的启动命令如下备用：

```bash
./frp-panel client -s abc -i user.c.client1 --api-url http://frpp.example.com:9000 --rpc-url grpc://frpp-rpc.example.com:9001
```

注意，如果你使用 反向代理 TLS，需要修改这行命令类似如下：

```bash
./frp-panel client -s abc -i user.c.client1 --api-url https://frpp.example.com:443 --rpc-url wss://frpp.example.com:443
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
    command: client -s abc -i user.c.client1 --api-url https://frpp.example.com:443 --rpc-url wss://frpp.example.com:443
```

### 安装为 systemd 服务

frp-panel 拥有管理 systemd 服务的能力，服务名为 `frpp`，内置了很多命令，请使用 `frp-panel --help` 查看支持的命令。这里给出一些例子：

- 安装特定参数的 client 到 systemd （支持任意的参数，包括server）
```bash
sudo ./frp-panel install [client 参数]
# eg. frp-panel install client -s abc -i user.c.client1 --api-url https://frpp.example.com:443 --rpc-url wss://frpp.example.com:443
```

- 卸载 frpp 服务
```bash
sudo ./frp-panel uninstall
```

- 启动 frpp 服务
```bash
sudo ./frp-panel start
```

- 停止 frpp 服务
```bash
sudo ./frp-panel stop
```

- 重启 frpp 服务
```bash
sudo ./frp-panel restart
```

## 在 Windows 上部署

### 直接运行

在 powershell 中，可执行文件的同目录下运行 WebUI 中复制的启动命令

```powershell
.\frp-panel.exe client -s abc -i user.c.client1 --api-url https://frpp.example.com:443 --rpc-url wss://frpp.example.com:443
```

### 安装为服务

与上方 Linux 的命令一致，修改文件名，去掉sudo执行即可

Windows 安装后使用示例：

```
C:/frpp/frpp.exe stop
C:/frpp/frpp.exe start
C:/frpp/frpp.exe uninstall
```
