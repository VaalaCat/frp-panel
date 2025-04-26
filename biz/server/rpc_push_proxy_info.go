package server

import (
	"context"

	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/utils/logger"
	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/samber/lo"
)

var proxyTypeList = []v1.ProxyType{
	v1.ProxyTypeTCP,
	v1.ProxyTypeUDP,
	v1.ProxyTypeTCPMUX,
	v1.ProxyTypeHTTP,
	v1.ProxyTypeHTTPS,
	v1.ProxyTypeSTCP,
	v1.ProxyTypeXTCP,
	v1.ProxyTypeSUDP,
}

func PushProxyInfo(appInstance app.Application, serverID, serverSecret string) error {
	proxyInfos := []*pb.ProxyInfo{}

	if cli := appInstance.GetServerController().Get(serverID); cli != nil {
		firstSync := cli.IsFirstSync()
		for _, proxyType := range proxyTypeList {
			proxyStatsList := cli.GetProxyStatsByType(proxyType)
			for _, proxyStats := range proxyStatsList {
				if proxyStats == nil {
					continue
				}
				proxyInfos = append(proxyInfos, &pb.ProxyInfo{
					Name:            lo.ToPtr(proxyStats.Name),
					Type:            lo.ToPtr(proxyStats.Type),
					TodayTrafficIn:  lo.ToPtr(proxyStats.TodayTrafficIn),
					TodayTrafficOut: lo.ToPtr(proxyStats.TodayTrafficOut),
					FirstSync:       lo.ToPtr(firstSync),
				})
			}
		}
	}

	if len(proxyInfos) > 0 {
		ctx := context.Background()
		cli := appInstance.GetMasterCli()
		_, err := cli.Call().PushProxyInfo(ctx, &pb.PushProxyInfoReq{
			Base: &pb.ServerBase{
				ServerId:     serverID,
				ServerSecret: serverSecret,
			},
			ProxyInfos: proxyInfos,
		})
		if err != nil {
			logger.Logger(context.Background()).WithError(err).Error("cannot push proxy info")
			return err
		}
	}
	return nil
}
