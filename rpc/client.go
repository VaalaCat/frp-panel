package rpc

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/VaalaCat/frp-panel/pb"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

func CallClient(c context.Context, clientID string, event pb.Event, msg proto.Message) (*pb.ClientMessage, error) {
	sender := GetClientsManager().Get(clientID)
	if sender == nil {
		logrus.Errorf("cannot get client, id: [%s]", clientID)
		return nil, fmt.Errorf("cannot get client, id: [%s]", clientID)
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		logrus.WithError(err).Errorf("cannot marshal")
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
		logrus.WithError(err).Errorf("cannot send")
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
		for {
			reciver := GetClientsManager().Get(clientID)
			if reciver == nil {
				logrus.Errorf("cannot get client")
				continue
			}
			resp, err := reciver.Conn.Recv()
			if err == io.EOF {
				logrus.Infof("finish client recv")
				done <- true
				return
			}
			if err != nil {
				logrus.WithError(err).Errorf("cannot recv, usually means client disconnect")
				done <- true
				return
			}

			respChAny, ok := recvMap.Load(resp.SessionId)
			if !ok {
				logrus.Errorf("cannot load")
				continue
			}

			respCh, ok := respChAny.(chan *pb.ClientMessage)
			if !ok {
				logrus.Errorf("cannot cast")
				continue
			}
			logrus.Infof("recv success, resp: %+v", resp)
			respCh <- resp
		}
	}()
	return done
}
