## Project Development Guide

### Platform Architecture Design

After choosing the tech stack, the next step is to design the program architecture. As mentioned in the background, frp itself has frpc and frps (client and server), these two roles are indispensable. Then we need to add something new to manage them, so frp-panel introduces a new master role. The master will be responsible for managing various frpc and frps, as well as centrally storing configuration files and connection information.

Next, we have frpc and frps. The original version requires writing configuration files on both sides. Since the original version already supports this, we don't need to follow the original approach. We will directly not support configuration files, and all configurations must be obtained from the master.

In addition, we also need to consider the compatibility with the original version. The client/server of frp-panel must be able to connect to the official frpc/frps service. In this way, both configuration file and non-configuration file modes can work perfectly.
Overall, the architecture is quite simple.

![arch](../public/images/arch.png)

### Development

The project includes three roles:

1. Master: The control node, accepts requests from the frontend and is responsible for managing Client and Server.
2. Server: The server side, controlled by the control node, responsible for providing services to clients, including frps and rpc (for connecting to the Master) services.
3. Client: The client side, controlled by the control node, including frpc and rpc (for connecting to the Master) services.

Next, we will provide the functionality of each package in the project:

```
.
|-- biz                 # Main business logic
|   |-- client          # Client logic (here referring to the frp-panel client)
|   |-- master          # frp-panel control plane, responsible for handling frontend requests, and using rpc to manage frp-panel's server and client
|   |   |-- auth        # Authentication module, including user authentication and client authentication
|   |   |-- client      # Client module, including various APIs for the frontend to manage clients
|   |   |-- server      # Server module, including various APIs for the frontend to manage servers
|   |   `-- user        # User module, including user management, user information retrieval, etc.
|   `-- server          # Server logic
|-- cache               # Cache, used to store frps authentication tokens
|-- cmd                 # Command line entry, where the main function is located, responsible for starting various modules as needed
|-- common
|-- conf
|-- dao                 # Data access object, any operations related to the database will call this library
|-- doc                 # Documentation
|-- idl                 # IDL definitions
|-- middleware          # API middleware, including JWT and context-related, used to process API requests. After authentication passes, user information will be injected into the context and can be obtained through the common package.
|-- models              # Database models, used to define database tables. Also includes entity definitions.
|-- pb                  # Generated protobuf pb files
|-- rpc                 # Location of various rpcs, including the logic for Client/Server to call Master, as well as the logic for Master to use Stream to call Client and Server
|-- services            # Various modules that need to run persistently in memory, this package can manage the running/stopping of various services
|   |-- api             # API service, requires an external ginRouter to run
|   |-- client          # frp client, i.e., frpc, can control various configurations/start and stop of frpc
|   |-- master          # Master service, including the rpc server definition, after receiving an rpc request, it will call the biz package to handle the logic
|   |-- rpcclient       # Stateful rpc client, because the rpc clients don't have public IP addresses, the rpcclient will call the master's stream long-connection rpc when starting, and after the connection is established, the Master and Client communicate through this package
|   `-- server          # frp server, i.e., frps, can control various configurations/start and stop of frps
|-- tunnel              # Tunnel module, used to manage tunnels, i.e., manage frpc and frps services
|-- utils
|-- watcher             # Scheduled tasks, e.g., updating configuration files every 30 seconds
`-- www
    |-- api
    |-- components      # There is an apitest component here for testing
    |   `-- ui
    |-- lib
    |   `-- pb
    |-- pages
    |-- public
    |-- store
    |-- styles
    `-- types
```

### Debugging and Startup Methods:

- master: `go run cmd/*.go master`
  > For client and server, please copy the content from the master webui
- client: `go run cmd/*.go client -i <clientID> -s <clientSecret>`
- server: `go run cmd/*.go server -i <serverID> -s <serverSecret>`

The project configuration file will read the .env file in the current folder by default. The project includes a sample configuration file, which can be modified according to your needs.

Detailed architecture call diagram:

![structure](../public/images/callvis.svg)

### Core Configuration Explanation

[settings.go](https://github.com/VaalaCat/frp-panel/blob/main/conf/settings.go)
This file contains detailed explanations of the configuration parameters. Please refer to this file if you need to further modify the configuration.
