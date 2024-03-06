package client

import (
	"context"
	"fmt"

	"github.com/VaalaCat/frp-panel/models"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/samber/lo"
)

func RPCPullConfig(ctx context.Context, req *pb.PullClientConfigReq) (*pb.PullClientConfigResp, error) {
	var (
		err error
		cli *models.ClientEntity
	)

	if cli, err = ValidateClientRequest(req.GetBase()); err != nil {
		return nil, err
	}

	if cli.Stopped {
		return nil, fmt.Errorf("client is stopped")
	}

	return &pb.PullClientConfigResp{
		Client: &pb.Client{
			Id:     lo.ToPtr(cli.ClientID),
			Config: lo.ToPtr(string(cli.ConfigContent)),
		},
	}, nil
}
