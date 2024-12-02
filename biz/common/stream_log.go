package common

import (
	"context"
	"sync"

	"github.com/VaalaCat/frp-panel/logger"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/sirupsen/logrus"
)

var (
	h = &HookMgr{
		Mutex: &sync.Mutex{},
	}
)

type HookMgr struct {
	*sync.Mutex
	hook *logger.StreamLogHook
}

func (h *HookMgr) Close() {
	h.Lock()
	defer h.Unlock()
	if h.hook != nil {
		h.hook.Close()
		h.hook = nil
	}
	logger.Instance().ReplaceHooks(logrus.LevelHooks{})
}

func (h *HookMgr) AddStream(send func(msg string), closeSend func()) {
	h.Lock()
	defer h.Unlock()
	h.hook = logger.NewStreamLogHook(send, closeSend)
	logger.Instance().AddHook(h.hook)
	go h.hook.Send()
}

func StartSteamLogHandler(ctx context.Context, req *pb.CommonRequest, initStreamLogFunc func(*HookMgr)) (*pb.CommonResponse, error) {
	logger.Logger(ctx).Infof("get a start stream log request, origin is: [%+v]", req)

	StopSteamLogHandler(ctx, req)
	initStreamLogFunc(h)

	return &pb.CommonResponse{}, nil
}

func StopSteamLogHandler(ctx context.Context, req *pb.CommonRequest) (*pb.CommonResponse, error) {
	logger.Logger(ctx).Infof("get a stop stream log request, origin is: [%+v]", req)
	h.Close()
	return &pb.CommonResponse{}, nil
}
