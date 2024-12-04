package streamlog

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/VaalaCat/frp-panel/biz/master/client"
	"github.com/VaalaCat/frp-panel/biz/master/server"
	"github.com/VaalaCat/frp-panel/logger"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/utils"
)

const (
	CacheBufSize = 4096
)

type ClientLogManager struct {
	sLogMap        *utils.SyncMap[string, chan string]
	clientLocksMap *utils.SyncMap[string, *sync.Mutex]
}

func (c *ClientLogManager) GetClientLock(clientId string) *sync.Mutex {
	lock, _ := c.clientLocksMap.LoadOrStore(clientId, &sync.Mutex{})
	return lock
}

var (
	mgr *ClientLogManager
)

func Mgr() *ClientLogManager {
	if mgr == nil {
		Init()
	}
	return mgr
}

func Init() {
	mgr = &ClientLogManager{
		sLogMap:        &utils.SyncMap[string, chan string]{},
		clientLocksMap: &utils.SyncMap[string, *sync.Mutex]{},
	}
}

func PushClientStreamLog(sender pb.Master_PushClientStreamLogServer) error {
	for {
		req, err := sender.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Logger(context.Background()).WithError(err).Errorf("cannot recv from client, id: [%+v]", req.GetBase())
			return err
		}

		_, err = client.ValidateClientRequest(req.GetBase())
		if err != nil {
			logger.Logger(context.Background()).WithError(err).Errorf("cannot validate client, id: [%+v]", req.GetBase())
			return err
		}

		ch, ok := Mgr().sLogMap.Load(req.GetBase().GetClientId())
		if !ok {
			return fmt.Errorf("push client stream log cannot find client, id: [%s]", req.GetBase().GetClientId())
		}

		ch <- string(req.GetLog())
	}
	return nil
}

func PushServerStreamLog(sender pb.Master_PushServerStreamLogServer) error {
	for {
		req, err := sender.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Logger(context.Background()).WithError(err).Errorf("cannot recv from server, req: [%+v]", req.GetBase())
			return err
		}

		_, err = server.ValidateServerRequest(req.GetBase())
		if err != nil {
			logger.Logger(context.Background()).WithError(err).Errorf("cannot validate server, req: [%+v]", req.GetBase())
			return err
		}

		ch, ok := Mgr().sLogMap.Load(req.GetBase().GetServerId())
		if !ok {
			return fmt.Errorf("push server stream log cannot find server, id: [%s]", req.GetBase().GetServerId())
		}
		ch <- string(req.GetLog())
	}
	return nil
}
