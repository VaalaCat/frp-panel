package rpc

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/defs"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func CallClientWrapper[R common.RespType](c *app.Context, clientID string, event pb.Event, req proto.Message, resp *R) error {
	cresp, err := CallClient(c, clientID, event, req)
	if err != nil {
		return err
	}

	protoMsgRef, ok := any(resp).(protoreflect.ProtoMessage)
	if !ok {
		return fmt.Errorf("type does not implement protoreflect.ProtoMessage")
	}

	return proto.Unmarshal(cresp.GetData(), protoMsgRef)
}

const CallTimeout = 15 * time.Second

var ErrClientOffline = errors.New("client is offline")

type pendingResponse struct {
	connector *defs.Connector
	response  chan *pb.ClientMessage
}

func CallClient(ctx *app.Context, clientID string, event pb.Event, msg proto.Message) (*pb.ClientMessage, error) {
	sender := ctx.GetApp().GetClientsManager().Get(clientID)
	if sender == nil {
		return nil, fmt.Errorf("%w: %s", ErrClientOffline, clientID)
	}
	timeout := CallTimeout
	switch event {
	case pb.Event_EVENT_INSTALL_WORKERD, pb.Event_EVENT_UPGRADE_FRPP:
		timeout = 5 * time.Minute
	case pb.Event_EVENT_UPDATE_FRPC, pb.Event_EVENT_START_FRPC, pb.Event_EVENT_UPDATE_FRPS:
		timeout = time.Minute
	}
	callCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}
	req := &pb.ServerMessage{Event: event, Data: data, SessionId: uuid.NewString(), ClientId: clientID}
	pending := &pendingResponse{connector: sender, response: make(chan *pb.ClientMessage, 1)}
	responses := ctx.GetApp().GetClientRecvMap()
	responses.Store(req.SessionId, pending)
	defer responses.Delete(req.SessionId)
	if err := sender.Send(callCtx, req); err != nil {
		return nil, fmt.Errorf("send to client %s: %w", clientID, err)
	}
	select {
	case resp := <-pending.response:
		if resp.GetEvent() == pb.Event_EVENT_ERROR {
			return nil, fmt.Errorf("client returned error: %s", resp.GetData())
		}
		return resp, nil
	case <-callCtx.Done():
		return nil, fmt.Errorf("client %s did not respond: %w", clientID, callCtx.Err())
	case <-sender.Done():
		return nil, fmt.Errorf("%w: %s", ErrClientOffline, clientID)
	case <-sender.Conn.Context().Done():
		return nil, fmt.Errorf("%w: %s", ErrClientOffline, clientID)
	}
}

func deliverResponse(responses *sync.Map, connector *defs.Connector, resp *pb.ClientMessage) {
	if resp == nil {
		return
	}
	value, ok := responses.Load(resp.GetSessionId())
	if !ok {
		return
	}
	pending, ok := value.(*pendingResponse)
	if !ok || pending.connector != connector || !responses.CompareAndDelete(resp.GetSessionId(), pending) {
		return
	}
	// Late or duplicate replies must never block the stream receive loop.
	select {
	case pending.response <- resp:
	default:
	}
}

func Recv(appInstance app.Application, clientID string) chan bool {
	return RecvConnector(appInstance, appInstance.GetClientsManager().Get(clientID))
}

func RecvConnector(appInstance app.Application, connector *defs.Connector) chan bool {
	done := make(chan bool)
	if connector == nil {
		close(done)
		return done
	}
	go func() {
		defer close(done)
		defer appInstance.GetClientsManager().RemoveIf(connector.CliID, connector)
		for {
			resp, err := connector.Conn.Recv()
			if err != nil {
				return
			}
			deliverResponse(appInstance.GetClientRecvMap(), connector, resp)
		}
	}()
	return done
}
