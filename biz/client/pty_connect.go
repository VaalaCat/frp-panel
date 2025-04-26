package client

import (
	"github.com/VaalaCat/frp-panel/app"
	"github.com/VaalaCat/frp-panel/biz/common"
	"github.com/VaalaCat/frp-panel/pb"
)

func StartPTYConnect(c *app.Context, req *pb.CommonRequest) (*pb.CommonResponse, error) {
	return common.StartPTYConnect(c, req, &pb.PTYClientMessage{Base: &pb.PTYClientMessage_ClientBase{
		ClientBase: &pb.ClientBase{
			ClientId:     c.GetApp().GetConfig().Client.ID,
			ClientSecret: c.GetApp().GetConfig().Client.Secret,
		},
	}})
}
