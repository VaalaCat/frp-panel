package client

import (
	"fmt"

	"github.com/VaalaCat/frp-panel/app"
	"github.com/VaalaCat/frp-panel/dao"
	"github.com/VaalaCat/frp-panel/logger"
	"github.com/VaalaCat/frp-panel/models"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/samber/lo"
)

func RPCPullConfig(ctx *app.Context, req *pb.PullClientConfigReq) (*pb.PullClientConfigResp, error) {
	var (
		err       error
		cli       *models.ClientEntity
		clientIDs []string
	)

	if cli, err = ValidateClientRequest(ctx, req.GetBase()); err != nil {
		logger.Logger(ctx).WithError(err).Errorf("cannot validate client request")
		return nil, err
	}

	if cli.IsShadow {
		proxies, err := dao.NewQuery(ctx).AdminListProxyConfigsWithFilters(&models.ProxyConfigEntity{
			OriginClientID: cli.ClientID,
		})
		if err != nil {
			logger.Logger(ctx).Infof("cannot get client ids in shadow, maybe not a shadow client, id: [%s]", cli.ClientID)
		}
		clientIDs = lo.Map(proxies, func(p *models.ProxyConfig, _ int) string { return p.ClientID })
	}

	if cli.Stopped && cli.IsShadow {
		return nil, fmt.Errorf("client is stopped")
	}

	return &pb.PullClientConfigResp{
		Client: &pb.Client{
			Id:             lo.ToPtr(cli.ClientID),
			ServerId:       lo.ToPtr(cli.ServerID),
			Config:         lo.ToPtr(string(cli.ConfigContent)),
			OriginClientId: lo.ToPtr(cli.OriginClientID),
			ClientIds:      clientIDs,
		},
	}, nil
}
