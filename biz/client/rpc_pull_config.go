package client

import (
	"context"
	"reflect"

	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/rpc"
	"github.com/VaalaCat/frp-panel/services/client"
	"github.com/VaalaCat/frp-panel/tunnel"
	"github.com/VaalaCat/frp-panel/utils"
	"github.com/sirupsen/logrus"
)

func PullConfig(clientID, clientSecret string) error {
	logrus.Infof("start to pull client config, clientID: [%s]", clientID)
	ctx := context.Background()
	cli, err := rpc.MasterCli(ctx)
	if err != nil {
		logrus.WithError(err).Error("cannot get master client")
		return err
	}
	resp, err := cli.PullClientConfig(ctx, &pb.PullClientConfigReq{
		Base: &pb.ClientBase{
			ClientId:     clientID,
			ClientSecret: clientSecret,
		},
	})
	if err != nil {
		logrus.WithError(err).Error("cannot pull client config")
		return err
	}

	if len(resp.GetClient().GetConfig()) == 0 {
		logrus.Infof("client [%s] config is empty, wait for server init", clientID)
		return nil
	}

	c, p, v, err := utils.LoadClientConfig([]byte(resp.GetClient().GetConfig()), true)
	if err != nil {
		logrus.WithError(err).Error("cannot load client config")
		return err
	}

	ctrl := tunnel.GetClientController()

	if t := ctrl.Get(clientID); t == nil {
		ctrl.Add(clientID, client.NewClientHandler(c, p, v))
		ctrl.Run(clientID)
	} else {
		if !reflect.DeepEqual(t.GetCommonCfg(), c) {
			logrus.Infof("client %s config changed, will recreate it", clientID)
			tcli := ctrl.Get(clientID)
			if tcli != nil {
				tcli.Stop()
				ctrl.Delete(clientID)
			}
			ctrl.Add(clientID, client.NewClientHandler(c, p, v))
			ctrl.Run(clientID)
		} else {
			logrus.Infof("client %s already exists, update if need", clientID)
			tcli := ctrl.Get(clientID)
			if tcli == nil || !tcli.Running() {
				if tcli != nil {
					tcli.Stop()
					ctrl.Delete(clientID)
				}
				ctrl.Add(clientID, client.NewClientHandler(c, p, v))
				ctrl.Run(clientID)
			} else {
				tcli.Update(p, v)
			}
		}
	}

	logrus.Infof("pull client config success, clientID: [%s]", clientID)
	return nil
}
