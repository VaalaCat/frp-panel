package streamlog

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/logger"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/rpc"
	"github.com/gin-gonic/gin"
	"github.com/sourcegraph/conc"
)

func GetLogHander(c *gin.Context) {
	id := c.Query("id")
	logger.Logger(c).Infof("user try to get stream log, id: [%s]", id)

	if id == "" {
		c.JSON(http.StatusBadRequest, common.Err("id is empty"))
		return
	}

	Mgr().GetClientLock(id).Lock()
	defer Mgr().GetClientLock(id).Unlock()

	ch := make(chan string, CacheBufSize)
	if oldCh, ok := Mgr().sLogMap.LoadAndDelete(id); ok {
		close(oldCh)
	}
	Mgr().sLogMap.Store(id, ch)

	_, err := rpc.CallClient(c, id, pb.Event_EVENT_START_STREAM_LOG, &pb.CommonRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Err(err.Error()))
		return
	}

	defer func() {
		Mgr().sLogMap.Delete(id)
		close(ch)
		rpc.CallClient(context.Background(), id, pb.Event_EVENT_STOP_STREAM_LOG, &pb.CommonRequest{})
	}()

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Encoding", "none")
	c.Writer.Flush()

	var wg conc.WaitGroup

	wg.Go(func() {
		for l := range ch {
			k, _ := json.Marshal(l)
			_, err := c.Writer.WriteString(string(k) + "\r\n")
			if err != nil {
				logger.Logger(c).Errorf("write log error: %v", err)
				break
			}
			c.Writer.Flush()
		}
	})

	select {
	case <-c.Request.Context().Done():
		return
	case <-c.Writer.CloseNotify():
		return
	}
}
