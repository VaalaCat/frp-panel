package client

import (
	"context"

	"github.com/VaalaCat/frp-panel/biz/common"
	"github.com/VaalaCat/frp-panel/conf"
	"github.com/VaalaCat/frp-panel/pb"
)

func StartPTYConnect(c context.Context, req *pb.CommonRequest) (*pb.CommonResponse, error) {
	return common.StartPTYConnect(c, req, &pb.PTYClientMessage{Base: &pb.PTYClientMessage_ClientBase{
		ClientBase: &pb.ClientBase{
			ClientId:     conf.Get().Client.ID,
			ClientSecret: conf.Get().Client.Secret,
		},
	}})
}
