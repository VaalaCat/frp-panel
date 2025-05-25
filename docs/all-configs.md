# 配置说明

## frp隧道高级模式配置

本面板完全兼容 frp 原本的`json`格式配置，仅需要将配置文件内容粘贴到服务端/客户端高级模式编辑框内，更新即可，详细的使用参考：[frp 文档](https://gofrp.org/zh-cn/docs/features/common/configure/)

## 程序启动配置文件

程序会按顺序读取以下文件内容作为配置文件：`.env`,`/etc/frpp/.env`

## 程序配置说明

> 文档可能有点老。。。
> 
> 完整的最新配置参考这个文件：[settings.go](https://github.com/VaalaCat/frp-panel/blob/main/conf/settings.go)

| 类型   | 环境变量名                             | 默认值               | 描述                                                             |
|--------|-------------------------------------|--------------------|----------------------------------------------------------------|
| string | `APP_SECRET`                       | -                  | 应用密钥，用于客户端和服务器的和Master的通信加密                        |
| string | `APP_GLOBAL_SECRET`                | `frp-panel`        | 全局密钥，用于管理生成密钥，需妥善保管                                 |
| int    | `APP_COOKIE_AGE`                   | `86400`            | Cookie 的有效期（秒），默认值为 1 天                                  |
| string | `APP_COOKIE_NAME`                  | `frp-panel-cookie` | Cookie 名称                                                        |
| string | `APP_COOKIE_PATH`                  | `/`                | Cookie 路径                                                       |
| string | `APP_COOKIE_DOMAIN`                | -                  | Cookie 域                                                         |
| bool   | `APP_COOKIE_SECURE`                | `false`            | Cookie 是否安全                                                   |
| bool   | `APP_COOKIE_HTTP_ONLY`             | `true`             | Cookie 是否仅限 HTTP                                             |
| bool   | `APP_ENABLE_REGISTER`              | `false`            | 是否启用注册，仅允许第一个管理员注册                               |
| int    | `MASTER_API_PORT`                  | `9000`             | 主节点 API 端口                                                  |
| string | `MASTER_API_HOST`                  | -                  | 主节点域名，可以在反向代理和CDN后                                 |
| string | `MASTER_API_SCHEME`                | `http`             | 主节点 API 协议（注意，这里不影响主机行为，设置为https只是为了方便复制客户端启动命令，HTTPS需要自行反向代理）|
| int    | `MASTER_CACHE_SIZE`                | `10`               | 缓存大小（MB）                                                   |
| string | `MASTER_RPC_HOST`                  | `127.0.0.1`        | Master节点公共 IP 或域名                                          |
| int    | `MASTER_RPC_PORT`                  | `9001`             | Master节点 RPC 端口                                            |
| bool   | `MASTER_COMPATIBLE_MODE`           | `false`            | 兼容模式，用于官方 frp 客户端                                     |
| string | `MASTER_INTERNAL_FRP_SERVER_HOST`  | -                  | Master内置 frps 服务器主机，用于客户端连接                                |
| int    | `MASTER_INTERNAL_FRP_SERVER_PORT`  | `9002`             | Master内置 frps 服务器端口，用于客户端连接                                |
| string | `MASTER_INTERNAL_FRP_AUTH_SERVER_HOST` | `127.0.0.1`    | Master内置 frps 认证服务器主机                                          |
| int    | `MASTER_INTERNAL_FRP_AUTH_SERVER_PORT` | `8999`          | Master内置 frps 认证服务器端口                                          |
| string | `MASTER_INTERNAL_FRP_AUTH_SERVER_PATH` | `/auth`         | Master内置 frps 认证服务器路径                                          |
| int    | `SERVER_API_PORT`                  | `8999`             | 服务器 API 端口                                                  |
| string | `DB_TYPE`                          | `sqlite3`         | 数据库类型，如 mysql postgres 或 sqlite3 等                                 |
| string | `DB_DSN`                           | `data.db`         | 数据库 DSN，默认使用sqlite3，数据默认存储在可执行文件同目录下，对于 sqlite 是路径，其他数据库为 DSN，参见 [MySQL DSN](https://github.com/go-sql-driver/mysql#dsn-data-source-name) |
| string | `CLIENT_ID`                        | -                  | 客户端 ID                                                        |
| string | `CLIENT_SECRET`                   | -                  | 客户端密钥                                                       |
