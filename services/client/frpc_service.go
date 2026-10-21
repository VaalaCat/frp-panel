package client

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/utils"
	"github.com/VaalaCat/frp-panel/utils/logger"
	"github.com/fatedier/frp/client"
	"github.com/fatedier/frp/client/proxy"
	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/fatedier/frp/pkg/config/v1/validation"
	"github.com/fatedier/frp/pkg/featuregate"
	"github.com/samber/lo"
)

type clientImpl struct {
	cli         *client.Service
	Common      *v1.ClientCommonConfig
	ProxyCfgs   map[string]v1.ProxyConfigurer
	VisitorCfgs map[string]v1.VisitorConfigurer
	done        chan struct{}
	running     atomic.Bool
	runOnce     sync.Once
	cfgMu       sync.RWMutex
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewClientHandler(commonCfg *v1.ClientCommonConfig,
	proxyCfgs []v1.ProxyConfigurer,
	visitorCfgs []v1.VisitorConfigurer) app.ClientHandler {
	ctx := context.Background()

	if len(commonCfg.FeatureGates) > 0 {
		if err := featuregate.SetFromMap(commonCfg.FeatureGates); err != nil {
			logger.Logger(ctx).WithError(err).Errorf("there's a feature gate settings, but set failed: %+v, skip", commonCfg.FeatureGates)
		}
	}

	warning, err := validation.ValidateAllClientConfig(commonCfg, proxyCfgs, visitorCfgs)
	if warning != nil {
		logger.Logger(ctx).WithError(err).Warnf("validate client config warning: %+v", warning)
	}
	if err != nil {
		logger.Logger(ctx).Panic(err)
	}

	cli, err := client.NewService(client.ServiceOptions{
		Common:      commonCfg,
		ProxyCfgs:   proxyCfgs,
		VisitorCfgs: visitorCfgs,
	})
	if err != nil {
		logger.Logger(ctx).Panic(err)
	}

	runCtx, cancel := context.WithCancel(context.Background())
	return &clientImpl{
		ctx: runCtx, cancel: cancel, done: make(chan struct{}),
		cli:         cli,
		Common:      commonCfg,
		ProxyCfgs:   lo.SliceToMap(proxyCfgs, utils.TransformProxyConfigurerToMap),
		VisitorCfgs: lo.SliceToMap(visitorCfgs, utils.TransformVisitorConfigurerToMap),
	}
}

func (c *clientImpl) Run() {
	c.runOnce.Do(func() {
		c.running.Store(true)
		defer func() { c.running.Store(false); close(c.done) }()
		if err := c.cli.Run(c.ctx); err != nil {
			logger.Logger(c.ctx).WithError(err).Error("run client failed")
		}
	})
}

func (c *clientImpl) Stop() {
	c.cancel()
	// A stopped handler must not start later when its scheduled Run wins the race.
	c.runOnce.Do(func() { close(c.done) })
	<-c.done
}

func (c *clientImpl) Update(proxyCfgs []v1.ProxyConfigurer, visitorCfgs []v1.VisitorConfigurer) {
	c.cfgMu.Lock()
	defer c.cfgMu.Unlock()
	c.ProxyCfgs = lo.SliceToMap(proxyCfgs, utils.TransformProxyConfigurerToMap)
	c.VisitorCfgs = lo.SliceToMap(visitorCfgs, utils.TransformVisitorConfigurerToMap)
	c.cli.UpdateAllConfigurer(proxyCfgs, visitorCfgs)
}

func (c *clientImpl) AddProxy(proxyCfg v1.ProxyConfigurer) {
	c.cfgMu.Lock()
	defer c.cfgMu.Unlock()
	c.ProxyCfgs[proxyCfg.GetBaseConfig().Name] = proxyCfg
	c.cli.UpdateAllConfigurer(lo.Values(c.ProxyCfgs), lo.Values(c.VisitorCfgs))
}

func (c *clientImpl) AddVisitor(visitorCfg v1.VisitorConfigurer) {
	c.cfgMu.Lock()
	defer c.cfgMu.Unlock()
	c.VisitorCfgs[visitorCfg.GetBaseConfig().Name] = visitorCfg
	c.cli.UpdateAllConfigurer(lo.Values(c.ProxyCfgs), lo.Values(c.VisitorCfgs))
}

func (c *clientImpl) RemoveProxy(proxyCfg v1.ProxyConfigurer) {
	c.cfgMu.Lock()
	defer c.cfgMu.Unlock()
	old := c.ProxyCfgs
	delete(old, proxyCfg.GetBaseConfig().Name)

	c.ProxyCfgs = old
	c.cli.UpdateAllConfigurer(lo.Values(c.ProxyCfgs), lo.Values(c.VisitorCfgs))
}

func (c *clientImpl) RemoveVisitor(visitorCfg v1.VisitorConfigurer) {
	c.cfgMu.Lock()
	defer c.cfgMu.Unlock()
	old := c.VisitorCfgs
	delete(old, visitorCfg.GetBaseConfig().Name)

	c.VisitorCfgs = old
	c.cli.UpdateAllConfigurer(lo.Values(c.ProxyCfgs), lo.Values(c.VisitorCfgs))
}

func (c *clientImpl) GetProxyStatus(name string) (*proxy.WorkingStatus, bool) {
	return c.cli.StatusExporter().GetProxyStatus(name)
}

func (c *clientImpl) GetCommonCfg() *v1.ClientCommonConfig {
	return c.Common
}

func (c *clientImpl) GetProxyCfgs() map[string]v1.ProxyConfigurer {
	c.cfgMu.RLock()
	defer c.cfgMu.RUnlock()
	return lo.Assign(c.ProxyCfgs)
}

func (c *clientImpl) GetVisitorCfgs() map[string]v1.VisitorConfigurer {
	c.cfgMu.RLock()
	defer c.cfgMu.RUnlock()
	return lo.Assign(c.VisitorCfgs)
}

func (c *clientImpl) Running() bool {
	return c.running.Load()
}

func (c *clientImpl) Wait() {
	<-c.done
}
