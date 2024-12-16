package client

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/VaalaCat/frp-panel/logger"
	"github.com/VaalaCat/frp-panel/utils"
	"github.com/fatedier/frp/client"
	"github.com/fatedier/frp/client/proxy"
	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/fatedier/frp/pkg/config/v1/validation"
	"github.com/fatedier/frp/pkg/util/log"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/sourcegraph/conc"
)

type ClientHandler interface {
	Run()
	Stop()
	Wait()
	Running() bool
	Update([]v1.ProxyConfigurer, []v1.VisitorConfigurer)
	AddProxy(v1.ProxyConfigurer)
	AddVisitor(v1.VisitorConfigurer)
	RemoveProxy(v1.ProxyConfigurer)
	RemoveVisitor(v1.VisitorConfigurer)
	GetProxyStatus(string) (*proxy.WorkingStatus, bool)
	GetCommonCfg() *v1.ClientCommonConfig
	GetProxyCfgs() map[string]v1.ProxyConfigurer
	GetVisitorCfgs() map[string]v1.VisitorConfigurer
}

type Client struct {
	cli         *client.Service
	Common      *v1.ClientCommonConfig
	ProxyCfgs   map[string]v1.ProxyConfigurer
	VisitorCfgs map[string]v1.VisitorConfigurer
	done        chan bool
	running     bool
}

var (
	cli *Client
)

func InitGlobalClientService(commonCfg *v1.ClientCommonConfig,
	proxyCfgs []v1.ProxyConfigurer,
	visitorCfgs []v1.VisitorConfigurer) {
	ctx := context.Background()
	if cli != nil {
		logger.Logger(ctx).Warn("client has been initialized")
		return
	}
	cli = NewClientHandler(commonCfg, proxyCfgs, visitorCfgs)
}

func GetGlobalClientSerivce() ClientHandler {
	if cli == nil {
		logrus.Panic("client has not been initialized")
	}
	return cli
}

func NewClientHandler(commonCfg *v1.ClientCommonConfig,
	proxyCfgs []v1.ProxyConfigurer,
	visitorCfgs []v1.VisitorConfigurer) *Client {
	ctx := context.Background()

	warning, err := validation.ValidateAllClientConfig(commonCfg, proxyCfgs, visitorCfgs)
	if warning != nil {
		logger.Logger(ctx).WithError(err).Warnf("validate client config warning: %+v", warning)
	}
	if err != nil {
		logger.Logger(ctx).Panic(err)
	}

	log.InitLogger(commonCfg.Log.To, commonCfg.Log.Level, int(commonCfg.Log.MaxDays), commonCfg.Log.DisablePrintColor)
	cli, err := client.NewService(client.ServiceOptions{
		Common:      commonCfg,
		ProxyCfgs:   proxyCfgs,
		VisitorCfgs: visitorCfgs,
	})
	if err != nil {
		logger.Logger(ctx).Panic(err)
	}

	return &Client{
		cli:         cli,
		Common:      commonCfg,
		ProxyCfgs:   lo.SliceToMap(proxyCfgs, utils.TransformProxyConfigurerToMap),
		VisitorCfgs: lo.SliceToMap(visitorCfgs, utils.TransformVisitorConfigurerToMap),
	}
}

func (c *Client) Run() {
	if c.running {
		logger.Logger(context.Background()).Warn("client is running, skip run")
		return
	}

	shouldGracefulClose := c.Common.Transport.Protocol == "kcp" || c.Common.Transport.Protocol == "quic"
	if shouldGracefulClose {
		var wg conc.WaitGroup
		wg.Go(func() { handleTermSignal(c.cli) })
	}
	c.running = true
	c.done = make(chan bool)

	defer func() {
		c.running = false
		close(c.done)
	}()

	var wg conc.WaitGroup
	wg.Go(func() {
		ctx := context.Background()
		logger.Logger(ctx).Infof("start to run client")
		if err := c.cli.Run(ctx); err != nil {
			logger.Logger(ctx).Errorf("run client error: %v", err)
		}
	})
	wg.Wait()
}

func (c *Client) Stop() {
	wg := conc.NewWaitGroup()
	wg.Go(func() { c.cli.Close() })
	wg.Wait()
}

func (c *Client) Update(proxyCfgs []v1.ProxyConfigurer, visitorCfgs []v1.VisitorConfigurer) {
	c.ProxyCfgs = lo.SliceToMap(proxyCfgs, utils.TransformProxyConfigurerToMap)
	c.VisitorCfgs = lo.SliceToMap(visitorCfgs, utils.TransformVisitorConfigurerToMap)
	c.cli.UpdateAllConfigurer(proxyCfgs, visitorCfgs)
}

func (c *Client) AddProxy(proxyCfg v1.ProxyConfigurer) {
	c.ProxyCfgs[proxyCfg.GetBaseConfig().Name] = proxyCfg
	c.cli.UpdateAllConfigurer(lo.Values(c.ProxyCfgs), lo.Values(c.VisitorCfgs))
}

func (c *Client) AddVisitor(visitorCfg v1.VisitorConfigurer) {
	c.VisitorCfgs[visitorCfg.GetBaseConfig().Name] = visitorCfg
	c.cli.UpdateAllConfigurer(lo.Values(c.ProxyCfgs), lo.Values(c.VisitorCfgs))
}

func (c *Client) RemoveProxy(proxyCfg v1.ProxyConfigurer) {
	old := c.ProxyCfgs
	delete(old, proxyCfg.GetBaseConfig().Name)

	c.ProxyCfgs = old
	c.cli.UpdateAllConfigurer(lo.Values(c.ProxyCfgs), lo.Values(c.VisitorCfgs))
}

func (c *Client) RemoveVisitor(visitorCfg v1.VisitorConfigurer) {
	old := c.VisitorCfgs
	delete(old, visitorCfg.GetBaseConfig().Name)

	c.VisitorCfgs = old
	c.cli.UpdateAllConfigurer(lo.Values(c.ProxyCfgs), lo.Values(c.VisitorCfgs))
}

func (c *Client) GetProxyStatus(name string) (*proxy.WorkingStatus, bool) {
	return c.cli.StatusExporter().GetProxyStatus(name)
}

func (c *Client) GetCommonCfg() *v1.ClientCommonConfig {
	return c.Common
}

func (c *Client) GetProxyCfgs() map[string]v1.ProxyConfigurer {
	return c.ProxyCfgs
}

func (c *Client) GetVisitorCfgs() map[string]v1.VisitorConfigurer {
	return c.VisitorCfgs
}

func (c *Client) Running() bool {
	return c.running
}

func (c *Client) Wait() {
	<-c.done
}

func handleTermSignal(svr *client.Service) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	svr.GracefulClose(500 * time.Millisecond)
}
