# 快速开始

## 开始前必读

`frp-panel` 有三个模块：

1. `master`: 中控模块，负责分发配置文件和控制所有其他模块
2. `server`: 对应`frps`，负责提供流量接入点
3. `client`: 对应`frpc`，可以将本地的服务暴露给`server`上某一个接入点

> 部署`master`时，`master`会启动一个默认的`default server` 供`client`连接，因此`master`一般不会独立存在，但你可以选择不使用

部署时，我们一般从 `master` 开始。`master` 负责管理的 `server` 和 `client` 部署时，需要用到成功部署后的 `master` 控制页面中自动生成的内容。

对于 `frp-panel` 我们**推荐所有的组件都使用 `docker` 部署**，并且**使用 `host` 网络**模式，除非你需要远程终端控制远端的机器时才使用服务安装到客户机

## 架构图

![](./public/images/arch.svg)
