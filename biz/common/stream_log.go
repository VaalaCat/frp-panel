package common

import (
	"sync"

	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/utils/logger"
	"github.com/sirupsen/logrus"
)

type HookMgr struct {
	*sync.Mutex
	hook *logger.StreamLogHook
	pkgs []string
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
	if h.pkgs == nil {
		h.pkgs = make([]string, 0)
	}
	h.hook = logger.NewStreamLogHook(send, closeSend, h.pkgs...)
	logger.Instance().AddHook(h.hook)
	go h.hook.Send()
}

func (h *HookMgr) SetPkgs(pkgs []string) {
	if h.Mutex == nil {
		h.Mutex = &sync.Mutex{}
	}
	h.Lock()
	defer h.Unlock()
	h.pkgs = pkgs
}

func StartSteamLogHandler(ctx *app.Context, req *pb.StartSteamLogRequest, initStreamLogFunc func(*app.Context, app.StreamLogHookMgr)) (*pb.CommonResponse, error) {
	logger.Logger(ctx).Infof("get a start stream log request, origin is: [%s]", req.String())

	StopSteamLogHandler(ctx, &pb.CommonRequest{})
	hookMgr := ctx.GetApp().GetStreamLogHookMgr()
	hookMgr.SetPkgs(req.GetPkgs())

	initStreamLogFunc(ctx, hookMgr)

	return &pb.CommonResponse{}, nil
}

func StopSteamLogHandler(ctx *app.Context, req *pb.CommonRequest) (*pb.CommonResponse, error) {
	logger.Logger(ctx).Infof("get a stop stream log request, origin is: [%s]", req.String())
	h := ctx.GetApp().GetStreamLogHookMgr()
	h.Close()
	return &pb.CommonResponse{}, nil
}
