package client

import (
	"context"

	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/samber/lo"
)

func ReportWireGuardRuntimeInfo(appInstance app.Application, clientID, clientSecret string) error {
	ctx := app.NewContext(context.Background(), appInstance)
	log := ctx.Logger().WithField("op", "ReportWireGuardRuntimeInfo")

	log.Debugf("start to report wireguard runtime info, clientID: [%s]", clientID)

	cli := ctx.GetApp().GetMasterCli()

	wgs := ctx.GetApp().GetWireGuardManager().GetAllServices()
	for _, wg := range wgs {
		runtimeInfo, err := wg.GetWGRuntimeInfo()
		if err != nil {
			log.WithError(err).Errorf("failed to get wireguard runtime info")
			continue
		}

		resp, err := cli.Call().ReportWireGuardRuntimeInfo(ctx, &pb.ReportWireGuardRuntimeInfoReq{
			Base: &pb.ClientBase{
				ClientId:     clientID,
				ClientSecret: clientSecret,
			},
			InterfaceName: lo.ToPtr(wg.GetBaseIfceConfig().InterfaceName),
			RuntimeInfo:   runtimeInfo,
		})
		if err != nil {
			log.WithError(err).Errorf("failed to report wireguard runtime info")
			return err
		}
		log.Debugf("report wireguard runtime info success, clientID: [%s], resp: %s", clientID, resp.String())
	}

	return nil
}
