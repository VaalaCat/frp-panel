package wg

import (
	"errors"
	"net/netip"

	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/models"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/dao"
)

func CreateNetwork(ctx *app.Context, req *pb.CreateNetworkRequest) (*pb.CreateNetworkResponse, error) {
	log := ctx.Logger().WithField("op", "CreateNetwork")

	userInfo := common.GetUserInfo(ctx)
	if !userInfo.Valid() {
		log.Errorf("invalid user")
		return nil, errors.New("invalid user")
	}
	if req.GetNetwork() == nil || len(req.GetNetwork().GetName()) == 0 || len(req.GetNetwork().GetCidr()) == 0 {
		return nil, errors.New("invalid request")
	}

	if _, err := netip.ParsePrefix(req.GetNetwork().GetCidr()); err != nil {
		log.WithError(err).Errorf("invalid cidr")
		return nil, errors.New("invalid cidr")
	}

	entity := &models.NetworkEntity{
		Name:     req.GetNetwork().GetName(),
		CIDR:     req.GetNetwork().GetCidr(),
		UserId:   uint32(userInfo.GetUserID()),
		TenantId: uint32(userInfo.GetTenantID()),
		ACL:      models.JSON[*pb.AclConfig]{Data: req.GetNetwork().GetAcl()},
	}

	if err := dao.NewMutation(ctx).CreateNetwork(userInfo, entity); err != nil {
		log.WithError(err).Errorf("create network error")
		return nil, err
	}

	log.Debugf("create network success")

	return &pb.CreateNetworkResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "success"},
		Network: &pb.Network{
			Id: 0, UserId: uint32(userInfo.GetUserID()),
			TenantId: uint32(userInfo.GetTenantID()),
			Name:     entity.Name, Cidr: entity.CIDR,
		},
	}, nil
}
