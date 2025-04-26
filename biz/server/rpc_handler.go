package server

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/VaalaCat/frp-panel/conf"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/utils/logger"
	"google.golang.org/protobuf/proto"
)

func HandleServerMessage(appInstance app.Application, req *pb.ServerMessage) *pb.ClientMessage {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("\n--------------------\ncatch panic !!! \nhandle server message error: %v, stack: %s\n--------------------\n", err, debug.Stack())
		}
	}()

	ctx := context.Background()
	logger.Logger(ctx).Infof("client get a server message, origin is: [%+v]", req)

	switch req.Event {
	case pb.Event_EVENT_UPDATE_FRPS:
		return app.WrapperServerMsg(appInstance, req, UpdateFrpsHander)
	case pb.Event_EVENT_REMOVE_FRPS:
		return app.WrapperServerMsg(appInstance, req, RemoveFrpsHandler)
	case pb.Event_EVENT_START_STREAM_LOG:
		return app.WrapperServerMsg(appInstance, req, StartSteamLogHandler)
	case pb.Event_EVENT_STOP_STREAM_LOG:
		return app.WrapperServerMsg(appInstance, req, StopSteamLogHandler)
	case pb.Event_EVENT_START_PTY_CONNECT:
		return app.WrapperServerMsg(appInstance, req, StartPTYConnect)
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
