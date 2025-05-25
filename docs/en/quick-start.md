# Quick Start

## Before You Begin

`frp-panel` consists of three modules:

1. `master`: the central control module, responsible for distributing configuration files and managing all other modules  
2. `server`: corresponds to `frps`, responsible for providing traffic entry points  
3. `client`: corresponds to `frpc`, which exposes local services to a specific entry point on the `server`  

> When you deploy the `master`, it will automatically start a default “default server” for clients to connect to. Therefore, the `master` is normally not used on its own, though you can choose to disable this feature.

In a typical deployment, we start with the `master`. When deploying `server` and `client` instances managed by the `master`, you will need the information automatically generated in the `master`’s web console after it has been successfully deployed.

For `frp-panel`, **we recommend deploying all components via Docker** and **using the `host` network mode**, unless you need remote terminal access to the target machine, in which case you may install the services directly on the host.

## Download Instructions

frp-panel supports deployment via Docker or direct execution. To deploy directly, download the release files here: [release](https://github.com/VaalaCat/frp-panel/releases)

Note: There are two binary versions—one is client-only, and the other is a full-featured executable. We recommend using the full-featured executable.

The client-only version can only execute the `client` command (no client parameters required) and its filename includes the “client” identifier.

After starting, the default example can be accessed at `http://IP:9000`

The first registered user is the administrator by default. Registration of additional users is disabled by default. To enable it, add the parameter `APP_ENABLE_REGISTER=true` in the Master’s startup command or configuration file.

> If you have questions about the configuration during deployment, please refer to the [Configuration Guide](./all-configs.md)  
> We recommend keeping this page open for reference.

## Architecture Diagram

![](../public/images/arch.svg)