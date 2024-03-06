package common

import (
	"fmt"
	"net/http"

	"github.com/VaalaCat/frp-panel/pb"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

type RespType interface {
	pb.UpdateFRPCResponse | pb.RemoveFRPCResponse |
		pb.UpdateFRPSResponse | pb.RemoveFRPSResponse |
		pb.CommonResponse | pb.RegisterResponse | pb.LoginResponse |
		pb.InitClientResponse | pb.ListClientsResponse | pb.GetClientResponse |
		pb.DeleteClientResponse |
		pb.InitServerResponse | pb.ListServersResponse | pb.GetServerResponse |
		pb.DeleteServerResponse |
		pb.GetUserInfoResponse | pb.UpdateUserInfoResponse |
		pb.GetPlatformInfoResponse | pb.GetClientsStatusResponse |
		pb.GetClientCertResponse |
		pb.StartFRPCResponse | pb.StopFRPCResponse | pb.StartFRPSResponse | pb.StopFRPSResponse
}

func OKResp[T RespType](c *gin.Context, origin *T) {
	c.Header(TraceIDKey, c.GetString(TraceIDKey))
	if c.ContentType() == "application/x-protobuf" {
		c.ProtoBuf(http.StatusOK, origin)
	} else {
		c.JSON(http.StatusOK, OK(ReqSuccess).WithBody(origin))
	}
}

func ErrResp[T RespType](c *gin.Context, origin *T, err string) {
	c.Header(TraceIDKey, c.GetString(TraceIDKey))
	if c.ContentType() == "application/x-protobuf" {
		c.ProtoBuf(http.StatusInternalServerError, origin)
	} else {
		c.JSON(http.StatusOK, Err(err).WithBody(origin))
	}
}

func ErrUnAuthorized(c *gin.Context, err string) {
	c.Header(TraceIDKey, c.GetString(TraceIDKey))
	if c.ContentType() == "application/x-protobuf" {
		c.ProtoBuf(http.StatusUnauthorized,
			&pb.CommonResponse{Status: &pb.Status{Code: pb.RespCode_RESP_CODE_UNAUTHORIZED, Message: err}})
	} else {
		c.JSON(http.StatusOK, Err(err))
	}
}

func ProtoResp[T RespType](origin *T) (*pb.ClientMessage, error) {
	switch ptr := any(origin).(type) {
	case *pb.CommonResponse:
		rawData, err := proto.Marshal(ptr)
		if err != nil {
			return nil, err
		}
		return &pb.ClientMessage{
			Event: pb.Event_EVENT_DATA,
			Data:  rawData,
		}, nil
	case *pb.UpdateFRPCResponse:
		rawData, err := proto.Marshal(ptr)
		if err != nil {
			return nil, err
		}
		return &pb.ClientMessage{
			Event: pb.Event_EVENT_UPDATE_FRPC,
			Data:  rawData,
		}, nil
	case *pb.RemoveFRPCResponse:
		rawData, err := proto.Marshal(ptr)
		if err != nil {
			return nil, err
		}
		return &pb.ClientMessage{
			Event: pb.Event_EVENT_REMOVE_FRPC,
			Data:  rawData,
		}, nil
	case *pb.UpdateFRPSResponse:
		rawData, err := proto.Marshal(ptr)
		if err != nil {
			return nil, err
		}
		return &pb.ClientMessage{
			Event: pb.Event_EVENT_UPDATE_FRPC,
			Data:  rawData,
		}, nil
	case *pb.RemoveFRPSResponse:
		rawData, err := proto.Marshal(ptr)
		if err != nil {
			return nil, err
		}
		return &pb.ClientMessage{
			Event: pb.Event_EVENT_REMOVE_FRPC,
			Data:  rawData,
		}, nil
	case *pb.StartFRPCResponse:
		rawData, err := proto.Marshal(ptr)
		if err != nil {
			return nil, err
		}
		return &pb.ClientMessage{
			Event: pb.Event_EVENT_START_FRPC,
			Data:  rawData,
		}, nil
	case *pb.StopFRPCResponse:
		rawData, err := proto.Marshal(ptr)
		if err != nil {
			return nil, err
		}
		return &pb.ClientMessage{
			Event: pb.Event_EVENT_STOP_FRPC,
			Data:  rawData,
		}, nil
	case *pb.StartFRPSResponse:
		rawData, err := proto.Marshal(ptr)
		if err != nil {
			return nil, err
		}
		return &pb.ClientMessage{
			Event: pb.Event_EVENT_START_FRPS,
			Data:  rawData,
		}, nil
	case *pb.StopFRPSResponse:
		rawData, err := proto.Marshal(ptr)
		if err != nil {
			return nil, err
		}
		return &pb.ClientMessage{
			Event: pb.Event_EVENT_STOP_FRPS,
			Data:  rawData,
		}, nil
	default:
	}
	return nil, fmt.Errorf("cannot unmarshal unknown type: %T", origin)
}
