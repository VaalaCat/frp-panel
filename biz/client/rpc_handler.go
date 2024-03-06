package client

import (
	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/sirupsen/logrus"
)

func HandleServerMessage(req *pb.ServerMessage) *pb.ClientMessage {
	logrus.Infof("client get a server message, origin is: [%+v]", req)
	switch req.Event {
	case pb.Event_EVENT_UPDATE_FRPC:
		return common.WrapperServerMsg(req, UpdateFrpcHander)
	case pb.Event_EVENT_REMOVE_FRPC:
		return common.WrapperServerMsg(req, RemoveFrpcHandler)
	case pb.Event_EVENT_START_FRPC:
		return common.WrapperServerMsg(req, StartFRPCHandler)
	case pb.Event_EVENT_STOP_FRPC:
		return common.WrapperServerMsg(req, StopFRPCHandler)
	case pb.Event_EVENT_PING:
		return &pb.ClientMessage{
			Event: pb.Event_EVENT_PONG,
			Data:  []byte("pong"),
		}
	default:
	}

	return &pb.ClientMessage{
		Event: pb.Event_EVENT_ERROR,
		Data:  []byte("unknown event"),
	}
}
