package client

import (
	"context"

	"github.com/VaalaCat/frp-panel/defs"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
)

func PullWireGuards(appInstance app.Application, clientID, clientSecret string) error {
	ctx := app.NewContext(context.Background(), appInstance)
	log := ctx.Logger().WithField("op", "PullWireGuards")

	log.Debugf("start to pull wireguards belong to client, clientID: [%s]", clientID)

	cli := ctx.GetApp().GetMasterCli()
	resp, err := cli.Call().ListClientWireGuards(ctx, &pb.ListClientWireGuardsRequest{
		Base: &pb.ClientBase{
			ClientId:     clientID,
			ClientSecret: clientSecret,
		},
	})
	if err != nil {
		log.WithError(err).Errorf("cannot list client wireguards, do not change anything")
		return err
	}

	if len(resp.GetWireguardConfigs()) == 0 {
		log.Debugf("client [%s] has no wireguards", clientID)
		return nil
	}

	log.Debugf("client [%s] has [%d] wireguards, check their status", clientID, len(resp.GetWireguardConfigs()))
	log.Tracef("wireguardConfigs: %v", resp.GetWireguardConfigs())
	wgMgr := ctx.GetApp().GetWireGuardManager()
	successCnt := 0
	for _, wireGuard := range resp.GetWireguardConfigs() {
		wgCfg := &defs.WireGuardConfig{WireGuardConfig: wireGuard}
		wgSvc, ok := wgMgr.GetService(wireGuard.GetInterfaceName())
		if ok {
			log.Debugf("wireguard [%s] already exists, skip create, update peers if need", wireGuard.GetInterfaceName())
			wgSvc.PatchPeers(wgCfg.GetParsedPeers())
			continue
		}

		wgSvc, err := wgMgr.CreateService(&defs.WireGuardConfig{WireGuardConfig: wireGuard})
		if err != nil {
			log.WithError(err).Errorf("create wireguard service failed")
			continue
		}
		err = wgSvc.Start()
		if err != nil {
			log.WithError(err).Errorf("start wireguard service failed")
			continue
		}
		successCnt++
	}

	log.Debugf("pull wireguards belong to client success, clientID: [%s], [%d] wireguards created", clientID, successCnt)

	return nil
}
