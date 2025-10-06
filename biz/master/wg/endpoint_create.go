package wg

import (
	"errors"

	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/models"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/dao"
)

func CreateEndpoint(ctx *app.Context, req *pb.CreateEndpointRequest) (*pb.CreateEndpointResponse, error) {
	userInfo := common.GetUserInfo(ctx)
	if !userInfo.Valid() {
		return nil, errors.New("invalid user")
	}
	e := req.GetEndpoint()
	if e == nil || len(e.GetHost()) == 0 || e.GetPort() == 0 {
		return nil, errors.New("invalid endpoint params")
	}

	entity := &models.EndpointEntity{Host: e.GetHost(), Port: e.GetPort(), ClientID: e.GetClientId()}
	if err := dao.NewQuery(ctx).CreateEndpoint(userInfo, entity); err != nil {
		return nil, err
	}
	return &pb.CreateEndpointResponse{Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "success"}, Endpoint: &pb.Endpoint{Id: 0, Host: entity.Host, Port: entity.Port, ClientId: entity.ClientID}}, nil
}
