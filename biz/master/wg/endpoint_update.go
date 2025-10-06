package wg

import (
	"errors"

	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/dao"
)

func UpdateEndpoint(ctx *app.Context, req *pb.UpdateEndpointRequest) (*pb.UpdateEndpointResponse, error) {
	userInfo := common.GetUserInfo(ctx)
	if !userInfo.Valid() {
		return nil, errors.New("invalid user")
	}
	e := req.GetEndpoint()
	if e == nil || e.GetId() == 0 || len(e.GetHost()) == 0 || e.GetPort() == 0 {
		return nil, errors.New("invalid endpoint params")
	}

	oldEndpoint, err := dao.NewQuery(ctx).GetEndpointByID(userInfo, uint(e.GetId()))
	if err != nil {
		return nil, err
	}

	// 只更新 host 和 port
	oldEndpoint.Host = e.GetHost()
	oldEndpoint.Port = e.GetPort()
	if err := dao.NewQuery(ctx).UpdateEndpoint(userInfo, uint(e.GetId()), oldEndpoint.EndpointEntity); err != nil {
		return nil, err
	}

	return &pb.UpdateEndpointResponse{Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "success"},
		Endpoint: oldEndpoint.ToPB(),
	}, nil
}
