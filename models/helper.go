package models

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/VaalaCat/frp-panel/utils"
	"github.com/VaalaCat/frp-panel/utils/logger"
	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/tiendc/go-deepcopy"
)

func ParseProxyConfigFromClient(client *Client) ([]*ProxyConfigEntity, error) {
	proxyCfg, err := utils.LoadProxiesFromContent(client.ConfigContent)
	if err != nil {
		return nil, err
	}

	resp := []*ProxyConfigEntity{}

	for _, cfg := range proxyCfg {
		tmpProxyEntity := &ProxyConfigEntity{}

		if err := tmpProxyEntity.FillClientConfig(client.ClientEntity); err != nil {
			return nil, err
		}

		if err := tmpProxyEntity.FillTypedProxyConfig(cfg); err != nil {
			return nil, err
		}

		resp = append(resp, tmpProxyEntity)
	}

	return resp, nil
}

func BuildClientConfigFromProxyConfig(client *Client, proxyCfgs []*ProxyConfig) (*Client, error) {
	if client == nil || len(proxyCfgs) == 0 {
		return nil, errors.New("client or proxy config is nil")
	}

	resp := &Client{}
	if err := deepcopy.Copy(resp, client); err != nil {
		return nil, err
	}

	cliCfg, err := utils.LoadClientConfigNormal(client.ConfigContent, true)
	if err != nil {
		return nil, err
	}

	pxyCfgs := []v1.TypedProxyConfig{}
	for _, proxyCfg := range proxyCfgs {
		pxy, err := utils.LoadProxiesFromContent(proxyCfg.Content)
		if err != nil {
			logger.Logger(context.Background()).WithError(err).Errorf("cannot load proxy config, name: [%s]", proxyCfg.Name)
			continue
		}

		pxyCfgs = append(pxyCfgs, pxy...)
	}

	cliCfg.Proxies = pxyCfgs
	cliCfgBytes, err := json.Marshal(cliCfg)
	if err != nil {
		return nil, err
	}

	client.ConfigContent = cliCfgBytes

	return client, nil
}
