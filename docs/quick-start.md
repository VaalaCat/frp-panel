# 快速开始

## 开始前必读

`frp-panel` 有三个模块：

1. `master`: 中控模块，负责分发配置文件和控制所有其他模块
2. `server`: 对应`frps`，负责提供流量接入点
3. `client`: 对应`frpc`，可以将本地的服务暴露给`server`上某一个接入点

> 部署`master`时，`master`会启动一个默认的`default server` 供`client`连接，因此`master`一般不会独立存在，但你可以选择不使用

部署时，我们一般从 `master` 开始。`master` 负责管理的 `server` 和 `client` 部署时，需要用到成功部署后的 `master` 控制页面中自动生成的内容。

对于 `frp-panel` 我们**推荐所有的组件都使用 `docker` 部署**，并且**使用 `host` 网络**模式，除非你需要远程终端控制远端的机器时，才使用服务安装到客户机

## 文件下载说明

frp-panel 可选 docker 和直接运行模式部署，直接部署请到 release 下载文件：[release](https://github.com/VaalaCat/frp-panel/releases)

注意：二进制有两种，一种是仅客户端，一种是全功能可执行文件，推荐使用全功能可执行文件。

客户端版只能执行 client 命令(无需 client 参数)，仅客户端版的名字会带有 client 标识

启动过后默认例子的访问地址为 `http://IP:9000`

默认第一个注册的用户是管理员。且默认不开放注册多用户，如果需要，请在 Master 启动命令或配置文件中添加参数：`APP_ENABLE_REGISTER=true`

> 如果在部署过程中，对配置有疑问，请参考 [配置说明](./all-configs.md)
> 
> 推荐单独打开一个页面随时参考

## 架构图

![](./public/images/arch.svg)
