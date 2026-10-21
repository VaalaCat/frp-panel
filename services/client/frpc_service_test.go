package client

import (
	"testing"
	"time"

	v1 "github.com/fatedier/frp/pkg/config/v1"
)

func TestStopBeforeScheduledRunDoesNotStartClient(t *testing.T) {
	cfg := &v1.ClientCommonConfig{}
	cfg.Complete()
	handler := NewClientHandler(cfg, nil, nil)
	handler.Stop()
	handler.Stop()
	done := make(chan struct{})
	go func() { handler.Run(); handler.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("stopped client started or Wait blocked")
	}
	if handler.Running() {
		t.Fatal("stopped client reports running")
	}
}
