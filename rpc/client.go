package rpc

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/logger"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func CallClientWrapper[R common.RespType](c context.Context, clientID string, event pb.Event, req proto.Message, resp *R) error {
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

func CallClient(c context.Context, clientID string, event pb.Event, msg proto.Message) (*pb.ClientMessage, error) {
	sender := GetClientsManager().Get(clientID)
	if sender == nil {
		logger.Logger(c).Errorf("cannot get client, id: [%s]", clientID)
		return nil, fmt.Errorf("cannot get client, id: [%s]", clientID)
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		logger.Logger(context.Background()).WithError(err).Errorf("cannot marshal")
		return nil, err
	}

	req := &pb.ServerMessage{
		Event:     event,
		Data:      data,
		SessionId: uuid.New().String(),
		ClientId:  clientID,
	}

	recvMap.Store(req.SessionId, make(chan *pb.ClientMessage))
	err = sender.Conn.Send(req)
	if err != nil {
		logger.Logger(context.Background()).WithError(err).Errorf("cannot send")
		GetClientsManager().Remove(clientID)
		return nil, err
	}
	respChAny, ok := recvMap.Load(req.SessionId)
	if !ok {
		logrus.Fatalf("cannot load")
	}

	respCh, ok := respChAny.(chan *pb.ClientMessage)
	if !ok {
		logrus.Fatalf("cannot cast")
	}

	resp := <-respCh
	if resp.Event == pb.Event_EVENT_ERROR {
		return nil, fmt.Errorf("client return error: %s", resp.Data)
	}

	close(respCh)
	recvMap.Delete(req.SessionId)
	return resp, nil
}

var (
	recvMap *sync.Map
)

func init() {
	recvMap = &sync.Map{}
}

func Recv(clientID string) chan bool {
	done := make(chan bool)
	go func() {
		c := context.Background()
		for {
			reciver := GetClientsManager().Get(clientID)
			if reciver == nil {
				logger.Logger(c).Errorf("cannot get client")
				continue
			}
			resp, err := reciver.Conn.Recv()
			if err == io.EOF {
				logger.Logger(c).Infof("finish client recv")
				done <- true
				return
			}
			if err != nil {
				logger.Logger(context.Background()).WithError(err).Errorf("cannot recv, usually means client disconnect")
				done <- true
				return
			}

			respChAny, ok := recvMap.Load(resp.SessionId)
			if !ok {
				logger.Logger(c).Errorf("cannot load")
				continue
			}

			respCh, ok := respChAny.(chan *pb.ClientMessage)
			if !ok {
				logger.Logger(c).Errorf("cannot cast")
				continue
			}
			logger.Logger(c).Infof("recv success, resp: %+v", resp)
			respCh <- resp
		}
	}()
	return done
}
