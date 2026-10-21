package clientrpc

import (
	"context"
	"io"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"google.golang.org/grpc"
)

type clientTestStream struct {
	grpc.ClientStream
	ctx      context.Context
	requests chan *pb.ServerMessage
	send     func(*pb.ClientMessage) error
}

func (s *clientTestStream) Send(m *pb.ClientMessage) error { return s.send(m) }
func (s *clientTestStream) Recv() (*pb.ServerMessage, error) {
	select {
	case r, ok := <-s.requests:
		if !ok {
			return nil, io.EOF
		}
		return r, nil
	case <-s.ctx.Done():
		return nil, s.ctx.Err()
	}
}
func (s *clientTestStream) CloseSend() error { return nil }
func TestHandlerReturnsOnEOF(t *testing.T) {
	requests := make(chan *pb.ServerMessage)
	close(requests)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream := &clientTestStream{ctx: ctx, requests: requests}
	done := make(chan struct{})
	go func() { runCLientRpcHandler(ctx, nil, stream, cancel, "client", nil); close(done) }()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("EOF caused busy loop")
	}
}
func TestHandlerSerializesRepliesAndStopsOnCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream := &clientTestStream{ctx: ctx, requests: make(chan *pb.ServerMessage, 32)}
	var sends atomic.Int32
	var received sync.WaitGroup
	received.Add(32)
	stream.send = func(*pb.ClientMessage) error {
		if sends.Add(1) != 1 {
			t.Error("concurrent reply Send")
		}
		time.Sleep(time.Millisecond)
		sends.Add(-1)
		received.Done()
		return nil
	}
	done := make(chan struct{})
	go func() {
		runCLientRpcHandler(ctx, nil, stream, cancel, "client", func(app.Application, *pb.ServerMessage) *pb.ClientMessage {
			return &pb.ClientMessage{Event: pb.Event_EVENT_PONG}
		})
		close(done)
	}()
	for range 32 {
		stream.requests <- &pb.ServerMessage{Event: pb.Event_EVENT_PING}
	}
	received.Wait()
	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("cancellation did not stop idle receiver")
	}
}
func TestRegistrationRejectsErrorAcknowledgement(t *testing.T) {
	stream := &clientTestStream{ctx: context.Background(), requests: make(chan *pb.ServerMessage, 1), send: func(*pb.ClientMessage) error { return nil }}
	stream.requests <- &pb.ServerMessage{Event: pb.Event_EVENT_REGISTER_CLIENT, ClientId: "client", Data: []byte("invalid token")}
	if registClientToMaster(stream, pb.Event_EVENT_REGISTER_CLIENT, "client", "secret") == nil {
		t.Fatal("rejected registration treated as success")
	}
}

type reconnectTestMaster struct {
	pb.MasterClient
	open func(context.Context) (pb.Master_ServerSendClient, error)
}

func (m reconnectTestMaster) ServerSend(ctx context.Context, _ ...grpc.CallOption) (pb.Master_ServerSendClient, error) {
	return m.open(ctx)
}

type reconnectTestClient struct{ master pb.MasterClient }

func (c reconnectTestClient) Call() pb.MasterClient { return c.master }

func TestRegistrationEOFReopensStream(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	var opened atomic.Int32
	reconnected := make(chan struct{})
	master := reconnectTestMaster{open: func(streamCtx context.Context) (pb.Master_ServerSendClient, error) {
		attempt := opened.Add(1)
		requests := make(chan *pb.ServerMessage, 1)
		if attempt == 1 {
			close(requests)
		} else {
			requests <- &pb.ServerMessage{Event: pb.Event_EVENT_REGISTER_CLIENT, ClientId: "client"}
			close(reconnected)
		}
		return &clientTestStream{ctx: streamCtx, requests: requests, send: func(*pb.ClientMessage) error { return nil }}, nil
	}}
	done := make(chan struct{})
	go func() {
		startClientRpcHandler(nil, reconnectTestClient{master: master}, ctx, "client", "secret", pb.Event_EVENT_REGISTER_CLIENT, nil)
		close(done)
	}()
	select {
	case <-reconnected:
	case <-ctx.Done():
		t.Fatal("registration EOF did not reopen the RPC stream")
	}
	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("reconnect loop did not stop after cancellation")
	}
}
