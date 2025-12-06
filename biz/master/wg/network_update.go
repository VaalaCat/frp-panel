package wg

import (
	"errors"

	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/models"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/dao"
)

func UpdateNetwork(ctx *app.Context, req *pb.UpdateNetworkRequest) (*pb.UpdateNetworkResponse, error) {
	userInfo := common.GetUserInfo(ctx)
	if !userInfo.Valid() {
		return nil, errors.New("invalid user")
	}
	n := req.GetNetwork()
	if n == nil || n.GetId() == 0 || len(n.GetName()) == 0 || len(n.GetCidr()) == 0 {
		return nil, errors.New("invalid network")
	}
	entity := &models.NetworkEntity{Name: n.GetName(), CIDR: n.GetCidr(), ACL: models.JSON[*pb.AclConfig]{Data: n.GetAcl()}}
	if err := dao.NewMutation(ctx).UpdateNetwork(userInfo, uint(n.GetId()), entity); err != nil {
		return nil, err
	}

	e := &models.Network{NetworkEntity: entity}
	return &pb.UpdateNetworkResponse{Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "success"},
		Network: e.ToPB(),
	}, nil
}
