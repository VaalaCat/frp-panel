package common

import (
	"context"
	"fmt"

	"github.com/VaalaCat/frp-panel/logger"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

func GlobalClientID(username, clientType, clientID string) string {
	return fmt.Sprintf("%s.%s.%s", username, clientType, clientID)
}

func Wrapper[T ReqType, U RespType](handler func(context.Context, *T) (*U, error)) func(c *gin.Context) {
	return func(c *gin.Context) {
		req, err := GetProtoRequest[T](c)
		if err != nil {
			ErrResp(c, &pb.CommonResponse{
				Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: err.Error()},
			}, err.Error())
			return
		}

		resp, err := handler(c, req)
		if err != nil {
			ErrResp(c, resp, err.Error())
			return
		}

		OKResp(c, resp)
	}
}

func WrapperServerMsg[T ReqType, U RespType](req *pb.ServerMessage, handler func(context.Context, *T) (*U, error)) *pb.ClientMessage {
	r := new(T)
	GetServerMessageRequest(req.GetData(), r, proto.Unmarshal)
	if err := GetServerMessageRequest(req.GetData(), r, proto.Unmarshal); err != nil {
		logger.Logger(context.Background()).WithError(err).Errorf("cannot unmarshal")
		return nil
	}

	ctx := context.Background()
	resp, err := handler(ctx, r)
	if err != nil {
		logger.Logger(context.Background()).WithError(err).Errorf("handler error")
		return &pb.ClientMessage{
			Event: pb.Event_EVENT_ERROR,
			Data:  []byte(err.Error()),
		}
	}

	cliMsg, err := ProtoResp(resp)
	if err != nil {
		logger.Logger(context.Background()).WithError(err).Errorf("cannot marshal")
		return &pb.ClientMessage{
			Event: pb.Event_EVENT_ERROR,
			Data:  []byte(err.Error()),
		}
	}
	return cliMsg
}
