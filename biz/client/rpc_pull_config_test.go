package client

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/tunnel"
	"github.com/samber/lo"
	"google.golang.org/grpc"
)

type fakeMaster struct {
	pb.MasterClient
	pull func(context.Context, *pb.PullClientConfigReq) (*pb.PullClientConfigResp, error)
}

func (f fakeMaster) PullClientConfig(ctx context.Context, r *pb.PullClientConfigReq, _ ...grpc.CallOption) (*pb.PullClientConfigResp, error) {
	return f.pull(ctx, r)
}

type fakeMasterClient struct{ cli pb.MasterClient }

func (f fakeMasterClient) Call() pb.MasterClient { return f.cli }

type trackedHandler struct {
	app.ClientHandler
	stopped atomic.Int32
}

func (h *trackedHandler) Stop() { h.stopped.Add(1) }
func TestPullConfigRemovesOnlyStaleChildren(t *testing.T) {
	a := app.NewApp()
	ctrl := tunnel.NewClientController()
	a.SetClientController(ctrl)
	activeA, activeB, stale := &trackedHandler{}, &trackedHandler{}, &trackedHandler{}
	ctrl.Add("a", "server", activeA)
	ctrl.Add("b", "server", activeB)
	ctrl.Add("stale", "server", stale)
	a.SetMasterCli(fakeMasterClient{fakeMaster{pull: func(_ context.Context, r *pb.PullClientConfigReq) (*pb.PullClientConfigResp, error) {
		c := &pb.Client{Id: lo.ToPtr(r.GetBase().GetClientId()), OriginClientId: lo.ToPtr("root")}
		if r.GetBase().GetClientId() == "root" {
			c.OriginClientId = nil
			c.ClientIds = []string{"a", "b", "a"}
		}
		return &pb.PullClientConfigResp{Client: c}, nil
	}}})
	if err := PullConfig(a, "root", "secret"); err != nil {
		t.Fatal(err)
	}
	if activeA.stopped.Load() != 0 || activeB.stopped.Load() != 0 {
		t.Fatal("active children stopped")
	}
	if ctrl.GetByClient("stale") != nil || stale.stopped.Load() != 1 {
		t.Fatal("stale child not removed")
	}
}
func TestFailedConfigPullKeepsExistingTunnels(t *testing.T) {
	a := app.NewApp()
	ctrl := tunnel.NewClientController()
	a.SetClientController(ctrl)
	active := &trackedHandler{}
	ctrl.Add("a", "server", active)
	a.SetMasterCli(fakeMasterClient{fakeMaster{pull: func(context.Context, *pb.PullClientConfigReq) (*pb.PullClientConfigResp, error) {
		return nil, errors.New("offline")
	}}})
	if PullConfig(a, "root", "secret") == nil {
		t.Fatal("expected failure")
	}
	if active.stopped.Load() != 0 || ctrl.GetByClient("a") == nil {
		t.Fatal("working config lost")
	}
}
func TestConcurrentConfigPullsAreSerialized(t *testing.T) {
	a := app.NewApp()
	a.SetClientController(tunnel.NewClientController())
	var active atomic.Int32
	a.SetMasterCli(fakeMasterClient{fakeMaster{pull: func(context.Context, *pb.PullClientConfigReq) (*pb.PullClientConfigResp, error) {
		if active.Add(1) != 1 {
			t.Error("overlapping reconciliation")
		}
		time.Sleep(time.Millisecond)
		active.Add(-1)
		return &pb.PullClientConfigResp{Client: &pb.Client{Id: lo.ToPtr("root")}}, nil
	}}})
	var wg sync.WaitGroup
	for range 12 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := PullConfig(a, "root", "secret"); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
}
