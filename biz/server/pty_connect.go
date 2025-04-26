package server

import (
	"github.com/VaalaCat/frp-panel/biz/common"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
)

func StartPTYConnect(c *app.Context, req *pb.CommonRequest) (*pb.CommonResponse, error) {
	return common.StartPTYConnect(c, req, &pb.PTYClientMessage{Base: &pb.PTYClientMessage_ServerBase{
		ServerBase: &pb.ServerBase{
			ServerId:     c.GetApp().GetConfig().Client.ID,
			ServerSecret: c.GetApp().GetConfig().Client.Secret,
		},
	}})
}
