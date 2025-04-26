package client

import (
	"context"

	"github.com/VaalaCat/frp-panel/app"
	"github.com/VaalaCat/frp-panel/dao"
	"github.com/VaalaCat/frp-panel/defs"
	"github.com/VaalaCat/frp-panel/logger"
	"github.com/VaalaCat/frp-panel/models"
	"github.com/samber/lo"
)

func SyncTunnel(ctx *app.Context, userInfo models.UserInfo) error {
	clis, err := dao.NewQuery(ctx).GetAllClients(userInfo)
	if err != nil {
		return err
	}
	lo.ForEach(clis, func(cli *models.ClientEntity, _ int) {
		cfg, err := cli.GetConfigContent()
		if err != nil {
			logger.Logger(context.Background()).WithError(err).Errorf("cannot get client config content, id: [%s]", cli.ClientID)
			return
		}

		cfg.User = userInfo.GetUserName()
		cfg.Metadatas = map[string]string{
			defs.FRPAuthTokenKey: userInfo.GetToken(),
		}
		if err := cli.SetConfigContent(*cfg); err != nil {
			logger.Logger(context.Background()).WithError(err).Errorf("cannot set client config content, id: [%s]", cli.ClientID)
			return
		}

		if err := dao.NewQuery(ctx).UpdateClient(userInfo, cli); err != nil {
			logger.Logger(context.Background()).WithError(err).Errorf("cannot update client, id: [%s]", cli.ClientID)
			return
		}
		logger.Logger(ctx).Infof("update client success, id: [%s]", cli.ClientID)
	})
	return nil
}
