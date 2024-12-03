package server

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/conf"
	"github.com/VaalaCat/frp-panel/logger"
	"github.com/VaalaCat/frp-panel/pb"
	"google.golang.org/protobuf/proto"
)

func HandleServerMessage(req *pb.ServerMessage) *pb.ClientMessage {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("\n--------------------\ncatch panic !!! \nhandle server message error: %v, stack: %s\n--------------------\n", err, debug.Stack())
		}
	}()

	ctx := context.Background()
	logger.Logger(ctx).Infof("client get a server message, origin is: [%+v]", req)

	switch req.Event {
	case pb.Event_EVENT_UPDATE_FRPS:
		return common.WrapperServerMsg(req, UpdateFrpsHander)
	case pb.Event_EVENT_REMOVE_FRPS:
		return common.WrapperServerMsg(req, RemoveFrpsHandler)
	case pb.Event_EVENT_START_STREAM_LOG:
		return common.WrapperServerMsg(req, StartSteamLogHandler)
	case pb.Event_EVENT_STOP_STREAM_LOG:
		return common.WrapperServerMsg(req, StopSteamLogHandler)
	case pb.Event_EVENT_START_PTY_CONNECT:
		return common.WrapperServerMsg(req, StartPTYConnect)
	case pb.Event_EVENT_PING:
		rawData, _ := proto.Marshal(conf.GetVersion().ToProto())
		return &pb.ClientMessage{
			Event: pb.Event_EVENT_PONG,
			Data:  rawData,
		}
	default:
	}

	return &pb.ClientMessage{
		Event: pb.Event_EVENT_ERROR,
		Data:  []byte("unknown event"),
	}
}
