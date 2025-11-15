# FRP-Panel

FRP-Panel 是一款基于 FRP 的可视化管理面板，提供中心化配置、统一凭证、动态调度和边缘 Worker 支持，让内网穿透和服务暴露更简单、更安全、更高效。

[详细使用文档 (Wiki)](https://vaala.cat/frp-panel) | [Blog 开发记录](https://vaala.cat/posts/frp-panel-doc/) | [截图/视频展示](https://vaala.cat/frp-panel/screenshots) | QQ 群：830620423

中文文档 | [English](./README.md)

<div align="center">
<a href="https://trendshift.io/repositories/7147" target="_blank"><img src="https://trendshift.io/api/badge/repositories/7147" alt="VaalaCat%2Ffrp-panel | Trendshift" style="width: 250px; height: 55px;" width="250" height="55"/></a>
</div>


## 核心优势

| 优势                       | 描述                                                         |
|:--------------------------|:------------------------------------------------------------|
| 中央化配置                 | 所有客户端/服务端配置由 Master 管理，无需手动编辑 JSON 文件          |
| 多节点统一管理             | 支持任意数量的 frpc（客户端）与 frps（服务端）节点集中监控与调度     |
| 可视化界面                 | Web UI 一键创建、编辑、监控隧道和Worker，实时日志与统计一目了然                  |
| 简化凭证分发               | 自动生成并分发启动命令，无须手动传参             |
| 边缘 Worker 自部署         | 在 Client 上部署自定义 Worker，Server 将其暴露到公网，Master 可实时调整配置 |

## 架构概览

![arch](docs/public/images/arch.png)

1. **Master** – 集中管理与鉴权，要求所有 Server 和 Client 可访问；
2. **Server** – 承载业务流量，作为公网入口，为 Client 提供服务；
3. **Client** – 内网代理，支持部署 Worker；

## 社区与赞助

本项目完全开源，欢迎 Star、Issues、PR。
若 FRP-Panel 为您带来价值，欢迎赞助作者：

-  邮箱：me@vaala.cat

[NodeSupport](https://github.com/NodeSeekDev/NodeSupport) / [林枫云](https://www.dkdun.cn) 赞助了该项目

<div align="left">
  <a href="https://yxvm.com/">
    <img src="https://github.com/user-attachments/assets/0bd7087a-7994-4caf-a465-a428af19c5aa" width="300" />
  </a>
</div>
<div align="left">
  <a href="https://www.dkdun.cn">
    <img src="https://www.dkdun.cn/themes/web/www/upload/local68c2dbb2ab148.png" width="300" />
  </a>
</div>

## 项目状态

[![Star History](https://api.star-history.com/svg?repos=vaalacat/frp-panel&type=Date)](https://www.star-history.com/#vaalacat/frp-panel&Date)

---

*更多部署、使用与配置细节，请移步 Wiki → [FRP-Panel WiKi](https://vaala.cat/frp-panel)*
