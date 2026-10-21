package rpc

import (
	"context"
	"errors"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"google.golang.org/grpc"
)

type testStream struct {
	grpc.ServerStream
	ctx  context.Context
	send func(*pb.ServerMessage) error
	recv chan *pb.ClientMessage
}

func (s *testStream) Context() context.Context         { return s.ctx }
func (s *testStream) Send(msg *pb.ServerMessage) error { return s.send(msg) }
func (s *testStream) Recv() (*pb.ClientMessage, error) {
	select {
	case msg, ok := <-s.recv:
		if !ok {
			return nil, io.EOF
		}
		return msg, nil
	case <-s.ctx.Done():
		return nil, s.ctx.Err()
	}
}
func testApp() app.Application {
	a := app.NewApp()
	a.SetClientsManager(NewClientsManager())
	a.SetClientRecvMap(&sync.Map{})
	return a
}
func assertNoPending(t *testing.T, a app.Application) {
	t.Helper()
	a.GetClientRecvMap().Range(func(_, _ any) bool { t.Error("pending RPC leaked"); return true })
}
func TestCallTimeoutLateResponseDoesNotBlockNextCall(t *testing.T) {
	a := testApp()
	streamCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sent := make(chan *pb.ServerMessage, 4)
	stream := &testStream{ctx: streamCtx, recv: make(chan *pb.ClientMessage, 4), send: func(msg *pb.ServerMessage) error { sent <- msg; return nil }}
	a.GetClientsManager().Set("client", "client", stream)
	done := Recv(a, "client")
	ctx, stop := context.WithTimeout(context.Background(), 30*time.Millisecond)
	_, err := CallClient(app.NewContext(ctx, a), "client", pb.Event_EVENT_PING, &pb.CommonRequest{})
	stop()
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected deadline, got %v", err)
	}
	assertNoPending(t, a)
	old := <-sent
	// The stale and duplicate replies must not consume the receiver for the next call.
	stream.recv <- &pb.ClientMessage{SessionId: old.SessionId, Event: pb.Event_EVENT_PONG}
	stream.recv <- &pb.ClientMessage{SessionId: old.SessionId, Event: pb.Event_EVENT_PONG}
	result := make(chan error, 1)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_, err := CallClient(app.NewContext(ctx, a), "client", pb.Event_EVENT_PING, &pb.CommonRequest{})
		result <- err
	}()
	next := <-sent
	stream.recv <- &pb.ClientMessage{SessionId: next.SessionId, Event: pb.Event_EVENT_PONG}
	if err := <-result; err != nil {
		t.Fatal(err)
	}
	assertNoPending(t, a)
	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("receive loop did not exit")
	}
}
func TestSendFailureAndClientErrorCleanPending(t *testing.T) {
	for _, sendFails := range []bool{true, false} {
		t.Run(map[bool]string{true: "send", false: "reply"}[sendFails], func(t *testing.T) {
			a := testApp()
			stream := &testStream{ctx: context.Background()}
			connector := a.GetClientsManager().Set("client", "client", stream)
			stream.send = func(msg *pb.ServerMessage) error {
				if sendFails {
					return io.ErrClosedPipe
				}
				deliverResponse(a.GetClientRecvMap(), connector, &pb.ClientMessage{SessionId: msg.SessionId, Event: pb.Event_EVENT_ERROR, Data: []byte("unavailable")})
				return nil
			}
			_, err := CallClient(app.NewContext(context.Background(), a), "client", pb.Event_EVENT_PING, &pb.CommonRequest{})
			if err == nil {
				t.Fatal("expected error")
			}
			assertNoPending(t, a)
		})
	}
}
func TestOldDisconnectDoesNotRemoveNewRegistration(t *testing.T) {
	a := testApp()
	old := a.GetClientsManager().Set("client", "client", &testStream{ctx: context.Background()})
	fresh := a.GetClientsManager().Set("client", "client", &testStream{ctx: context.Background()})
	a.GetClientsManager().RemoveIf("client", old)
	if a.GetClientsManager().Get("client") != fresh {
		t.Fatal("old receiver removed replacement")
	}
	select {
	case <-old.Done():
	default:
		t.Fatal("old stream not closed")
	}
}
func TestConnectorSerializesConcurrentSends(t *testing.T) {
	a := testApp()
	var mu sync.Mutex
	active := 0
	stream := &testStream{ctx: context.Background(), send: func(*pb.ServerMessage) error {
		mu.Lock()
		active++
		if active != 1 {
			t.Error("concurrent Send")
		}
		mu.Unlock()
		time.Sleep(time.Millisecond)
		mu.Lock()
		active--
		mu.Unlock()
		return nil
	}}
	conn := a.GetClientsManager().Set("client", "client", stream)
	var wg sync.WaitGroup
	for range 30 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := conn.Send(context.Background(), &pb.ServerMessage{}); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
}
