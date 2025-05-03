package app

import (
	"context"
	"fmt"

	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/utils/logger"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

func GlobalClientID(username, clientType, clientID string) string {
	return fmt.Sprintf("%s.%s.%s", username, clientType, clientID)
}

func ShadowedClientID(clientID string, shadowCount int64) string {
	return fmt.Sprintf("%s@%d", clientID, shadowCount)
}

func Wrapper[T common.ReqType, U common.RespType](appInstance Application, handler func(*Context, *T) (*U, error)) func(c *gin.Context) {
	return func(c *gin.Context) {
		req, err := common.GetProtoRequest[T](c)
		if err != nil {
			common.ErrResp(c, &pb.CommonResponse{
				Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: err.Error()},
			}, err.Error())
			return
		}

		resp, err := handler(NewContext(c, appInstance), req)
		if err != nil {
			common.ErrResp(c, resp, err.Error())
			return
		}

		common.OKResp(c, resp)
	}
}

func WrapperServerMsg[T common.ReqType, U common.RespType](appInstance Application, req *pb.ServerMessage,
	handler func(*Context, *T) (*U, error)) *pb.ClientMessage {
	r := new(T)
	common.GetServerMessageRequest(req.GetData(), r, proto.Unmarshal)
	if err := common.GetServerMessageRequest(req.GetData(), r, proto.Unmarshal); err != nil {
		logger.Logger(context.Background()).WithError(err).Errorf("cannot unmarshal")
		return nil
	}

	ctx := context.Background()
	appCtx := NewContext(ctx, appInstance)
	resp, err := handler(appCtx, r)
	if err != nil {
		logger.Logger(context.Background()).WithError(err).Errorf("handler error")
		return &pb.ClientMessage{
			Event: pb.Event_EVENT_ERROR,
			Data:  []byte(err.Error()),
		}
	}

	cliMsg, err := common.ProtoResp(resp)
	if err != nil {
		logger.Logger(context.Background()).WithError(err).Errorf("cannot marshal, may need to add this type to [getEvent] function")
		return &pb.ClientMessage{
			Event: pb.Event_EVENT_ERROR,
			Data:  []byte(err.Error()),
		}
	}
	return cliMsg
}
