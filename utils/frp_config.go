package utils

import (
	"fmt"

	v1 "github.com/fatedier/frp/pkg/config/v1"
)

func NewBaseFRPServerConfig(port int, token string) *v1.ServerConfig {
	resp := &v1.ServerConfig{
		BindPort: port,
		Auth: v1.AuthServerConfig{
			Method: v1.AuthMethodToken,
			Token:  token,
		},
	}
	resp.Complete()
	return resp
}

func NewBaseFRPServerUserAuthConfig(port int, opts []v1.HTTPPluginOptions) *v1.ServerConfig {
	resp := &v1.ServerConfig{
		BindPort:    port,
		HTTPPlugins: opts,
	}
	resp.Complete()
	return resp
}

func NewBaseFRPClientConfig(serverAddr string, serverPort int, token string) *v1.ClientCommonConfig {
	resp := &v1.ClientCommonConfig{
		Auth: v1.AuthClientConfig{
			Method: v1.AuthMethodToken,
			Token:  token,
		},
		ServerAddr: serverAddr,
		ServerPort: serverPort,
	}
	resp.Complete()
	return resp
}

func NewBaseFRPClientUserAuthConfig(serverAddr string, serverPort int, user, token string) *v1.ClientCommonConfig {
	resp := &v1.ClientCommonConfig{
		User: user,
		Metadatas: map[string]string{
			string(v1.AuthMethodToken): token,
		},
		ServerAddr: serverAddr,
		ServerPort: serverPort,
	}
	resp.Complete()
	return resp
}

func TransformProxyConfigurerToMap(origin v1.ProxyConfigurer) (key string, r v1.ProxyConfigurer) {
	key = origin.GetBaseConfig().Name
	r = origin
	return
}

func TransformVisitorConfigurerToMap(origin v1.VisitorConfigurer) (key string, r v1.VisitorConfigurer) {
	key = origin.GetBaseConfig().Name
	r = origin
	return
}

func NewProxyKey(clientID, serverID, proxyName string) string {
	return fmt.Sprintf("%s/%s/%s", clientID, serverID, proxyName)
}
