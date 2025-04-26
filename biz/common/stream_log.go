package common

import (
	"sync"

	"github.com/VaalaCat/frp-panel/app"
	"github.com/VaalaCat/frp-panel/logger"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/sirupsen/logrus"
)

type HookMgr struct {
	*sync.Mutex
	hook *logger.StreamLogHook
}

func (h *HookMgr) Close() {
	if h.Mutex == nil {
		h.Mutex = &sync.Mutex{}
	}
	h.Lock()
	defer h.Unlock()
	if h.hook != nil {
		h.hook.Close()
		h.hook = nil
	}
	logger.Instance().ReplaceHooks(logrus.LevelHooks{})
}

func (h *HookMgr) AddStream(send func(msg string), closeSend func()) {
	if h.Mutex == nil {
		h.Mutex = &sync.Mutex{}
	}
	h.Lock()
	defer h.Unlock()
	h.hook = logger.NewStreamLogHook(send, closeSend)
	logger.Instance().AddHook(h.hook)
	go h.hook.Send()
}

func StartSteamLogHandler(ctx *app.Context, req *pb.CommonRequest, initStreamLogFunc func(*app.Context, app.StreamLogHookMgr)) (*pb.CommonResponse, error) {
	logger.Logger(ctx).Infof("get a start stream log request, origin is: [%+v]", req)

	StopSteamLogHandler(ctx, req)
	hookMgr := ctx.GetApp().GetStreamLogHookMgr()
	initStreamLogFunc(ctx, hookMgr)

	return &pb.CommonResponse{}, nil
}

func StopSteamLogHandler(ctx *app.Context, req *pb.CommonRequest) (*pb.CommonResponse, error) {
	logger.Logger(ctx).Infof("get a stop stream log request, origin is: [%+v]", req)
	h := ctx.GetApp().GetStreamLogHookMgr()
	h.Close()
	return &pb.CommonResponse{}, nil
}
