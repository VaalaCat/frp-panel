# FRP-Panel

FRP-Panel is a visualization management dashboard for FRP, offering centralized configuration, unified credentials, dynamic scheduling, and edge Worker support—making NAT traversal and service exposure simpler, safer, and more efficient.

[Detailed Documentation (Wiki)](https://vaala.cat/frp-panel/en/) · [Development Blog](https://vaala.cat/posts/frp-panel-doc/) · [Screenshots & Videos](https://vaala.cat/posts/frp-panel-doc/en/screenshots) · QQ Group: 830620423

English | [中文](./README_zh.md)

<div align="center">
  <a href="https://trendshift.io/repositories/7147" target="_blank">
    <img src="https://trendshift.io/api/badge/repositories/7147" alt="VaalaCat/frp-panel | Trendshift" width="250" height="55"/>
  </a>
</div>

## Key Advantages

| Advantage               | Description                                                                 |
|:-----------------------|:----------------------------------------------------------------------------|
| Centralized Configuration | All client/server configs are managed by Master—no manual JSON editing       |
| Multi-node Management     | Monitor and orchestrate any number of frpc (clients) and frps (servers)      |
| Visual Interface          | Create, edit, and monitor tunnels and Workers via Web UI, with real-time logs and stats |
| Simplified Credential Distribution | Auto-generate and distribute startup commands—no manual parameter passing |
| Edge Worker Deployment    | Deploy custom Workers on Clients, expose them via Server, and adjust configs live via Master |

## Architecture Overview

![Architecture Diagram](./docs/public/images/arch.svg)

1. **Master** – Centralized management and authentication; requires access from all Servers and Clients  
2. **Server** – Public-facing entry point that handles traffic for Clients  
3. **Client** – Internal proxy that supports deploying Workers  

## Community & Sponsorship

FRP-Panel is fully open source—welcome Stars, Issues, and PRs.  
If FRP-Panel brings you value, consider sponsoring the author:

- Email: me@vaala.cat

Sponsored by [NodeSupport](https://github.com/NodeSeekDev/NodeSupport)

<div align="left">
  <a href="https://yxvm.com/">
    <img src="https://github.com/user-attachments/assets/0bd7087a-7994-4caf-a465-a428af19c5aa" width="300"/>
  </a>
</div>

## Project Status

[![Star History](https://api.star-history.com/svg?repos=vaalacat/frp-panel&type=Date)](https://www.star-history.com/#vaalacat/frp-panel&Date)

---

For more deployment, usage, and configuration details, see the Wiki → [FRP-Panel Wiki](https://vaala.cat/frp-panel/en/)