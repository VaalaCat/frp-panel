package clientrpc

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/utils/logger"
	"github.com/google/uuid"
)

func registClientToMaster(stream pb.Master_ServerSendClient, event pb.Event, clientID, secret string) error {
	if err := stream.Send(&pb.ClientMessage{Event: event, ClientId: clientID, SessionId: uuid.NewString(), Secret: secret}); err != nil {
		return err
	}
	resp, err := stream.Recv()
	if err != nil {
		return err
	}
	if resp.GetEvent() != event || resp.GetClientId() != clientID || len(resp.GetData()) != 0 {
		return fmt.Errorf("master rejected registration for %s", clientID)
	}
	return nil
}

func runCLientRpcHandler(ctx context.Context, appInstance app.Application, stream pb.Master_ServerSendClient, cancel context.CancelFunc, clientID string,
	handler func(app.Application, *pb.ServerMessage) *pb.ClientMessage) {
	var sendMu sync.Mutex
	workers := make(chan struct{}, 32)
	for {
		req, err := stream.Recv()
		if err != nil {
			return
		}
		select {
		case workers <- struct{}{}:
		case <-ctx.Done():
			return
		}
		go func(req *pb.ServerMessage) {
			defer func() { <-workers }()
			var msg *pb.ClientMessage
			func() {
				defer func() {
					if recovered := recover(); recovered != nil {
						logger.Logger(ctx).Errorf("RPC handler panic for event %s", req.GetEvent())
						msg = &pb.ClientMessage{Event: pb.Event_EVENT_ERROR, Data: []byte("client request failed")}
					}
				}()
				msg = handler(appInstance, req)
			}()
			if msg == nil || ctx.Err() != nil {
				return
			}
			msg.ClientId, msg.SessionId = clientID, req.GetSessionId()
			sendMu.Lock()
			defer sendMu.Unlock()
			if ctx.Err() != nil {
				return
			}
			if err := stream.Send(msg); err != nil {
				cancel()
			}
		}(req)
	}
}

func startClientRpcHandler(appInstance app.Application, client app.MasterClient, ctx context.Context, clientID, secret string, event pb.Event,
	handler func(app.Application, *pb.ServerMessage) *pb.ClientMessage) {
	for ctx.Err() == nil {
		streamCtx, cancel := context.WithCancel(ctx)
		stream, err := client.Call().ServerSend(streamCtx)
		if err == nil {
			// Registration has a deadline, while the successfully registered stream is long-lived.
			timer := time.AfterFunc(15*time.Second, cancel)
			err = registClientToMaster(stream, event, clientID, secret)
			timer.Stop()
			if err == nil {
				logger.Logger(ctx).Infof("registered RPC client %s", clientID)
				runCLientRpcHandler(streamCtx, appInstance, stream, cancel, clientID, handler)
			}
			cancel()
			_ = stream.CloseSend()
		} else {
			cancel()
		}
		if ctx.Err() != nil {
			return
		}
		logger.Logger(ctx).Warnf("RPC stream disconnected for %s; reconnecting", clientID)
		timer := time.NewTimer(3 * time.Second)
		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
		}
	}
}
