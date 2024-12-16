package utils

import (
	"github.com/fatedier/frp/pkg/config"
	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/samber/lo"
	"k8s.io/apimachinery/pkg/util/sets"
)

func LoadContentWithTemplate(content []byte, values *config.Values) ([]byte, error) {
	return config.RenderWithTemplate(content, values)
}

func LoadConfigureFromContent(content []byte, c any, strict bool) error {
	ans, err := LoadContentWithTemplate(content, config.GetValues())
	if err != nil {
		return err
	}
	return config.LoadConfigure(ans, c, strict)
}

func LoadProxiesFromContent(content []byte) ([]v1.TypedProxyConfig, error) {
	allCfg := &v1.ClientConfig{}

	if err := LoadConfigureFromContent(content, allCfg, true); err != nil {
		return nil, err
	}

	return allCfg.Proxies, nil
}

func LoadVisitorsFromContent(content []byte) ([]v1.TypedVisitorConfig, error) {
	allCfg := &v1.ClientConfig{}

	if err := LoadConfigureFromContent(content, allCfg, true); err != nil {
		return nil, err
	}

	return allCfg.Visitors, nil
}

func LoadClientConfigNormal(content []byte, strict bool) (*v1.ClientConfig, error) {
	var (
		cliCfg *v1.ClientCommonConfig
	)

	allCfg := v1.ClientConfig{}
	if err := LoadConfigureFromContent(content, &allCfg, strict); err != nil {
		return nil, err
	}
	cliCfg = &allCfg.ClientCommonConfig
	cliCfg.Complete()
	allCfg.ClientCommonConfig = *cliCfg
	return &allCfg, nil
}

func LoadClientConfig(content []byte, strict bool) (
	*v1.ClientCommonConfig,
	[]v1.ProxyConfigurer,
	[]v1.VisitorConfigurer,
	error,
) {
	var (
		cliCfg      *v1.ClientCommonConfig
		proxyCfgs   = make([]v1.ProxyConfigurer, 0)
		visitorCfgs = make([]v1.VisitorConfigurer, 0)
	)

	allCfg := v1.ClientConfig{}
	if err := LoadConfigureFromContent(content, &allCfg, strict); err != nil {
		return nil, nil, nil, err
	}
	cliCfg = &allCfg.ClientCommonConfig
	for _, c := range allCfg.Proxies {
		proxyCfgs = append(proxyCfgs, c.ProxyConfigurer)
	}
	for _, c := range allCfg.Visitors {
		visitorCfgs = append(visitorCfgs, c.VisitorConfigurer)
	}

	// Filter by start
	if len(cliCfg.Start) > 0 {
		startSet := sets.New(cliCfg.Start...)
		proxyCfgs = lo.Filter(proxyCfgs, func(c v1.ProxyConfigurer, _ int) bool {
			return startSet.Has(c.GetBaseConfig().Name)
		})
		visitorCfgs = lo.Filter(visitorCfgs, func(c v1.VisitorConfigurer, _ int) bool {
			return startSet.Has(c.GetBaseConfig().Name)
		})
	}

	cliCfg.Complete()

	for _, c := range proxyCfgs {
		c.Complete(cliCfg.User)
	}
	for _, c := range visitorCfgs {
		c.Complete(cliCfg)
	}
	return cliCfg, proxyCfgs, visitorCfgs, nil
}

func LoadServerConfig(content []byte, strict bool) (*v1.ServerConfig, error) {
	var (
		svrCfg = &v1.ServerConfig{}
	)
	if err := LoadConfigureFromContent(content, svrCfg, strict); err != nil {
		return nil, err
	}

	svrCfg.Complete()

	return svrCfg, nil
}
