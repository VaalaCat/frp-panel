package common

import (
	"fmt"
	"net/http"

	"github.com/VaalaCat/frp-panel/defs"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
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
		pb.StartFRPCResponse | pb.StopFRPCResponse | pb.StartFRPSResponse | pb.StopFRPSResponse |
		pb.GetProxyStatsByClientIDResponse | pb.GetProxyStatsByServerIDResponse |
		pb.CreateProxyConfigResponse | pb.ListProxyConfigsResponse | pb.UpdateProxyConfigResponse |
		pb.DeleteProxyConfigResponse | pb.GetProxyConfigResponse | pb.SignTokenResponse
}

func OKResp[T RespType](c *gin.Context, origin *T) {
	c.Header(defs.TraceIDKey, c.GetString(defs.TraceIDKey))
	if c.ContentType() == "application/x-protobuf" {
		c.ProtoBuf(http.StatusOK, origin)
	} else {
		c.JSON(http.StatusOK, OK(defs.ReqSuccess).WithBody(origin))
	}
}

func ErrResp[T RespType](c *gin.Context, origin *T, err string) {
	c.Header(defs.TraceIDKey, c.GetString(defs.TraceIDKey))
	if c.ContentType() == "application/x-protobuf" {
		c.ProtoBuf(http.StatusInternalServerError, origin)
	} else {
		c.JSON(http.StatusOK, Err(err).WithBody(origin))
	}
}

func ErrUnAuthorized(c *gin.Context, err string) {
	c.Header(defs.TraceIDKey, c.GetString(defs.TraceIDKey))
	if c.ContentType() == "application/x-protobuf" {
		c.ProtoBuf(http.StatusUnauthorized,
			&pb.CommonResponse{Status: &pb.Status{Code: pb.RespCode_RESP_CODE_UNAUTHORIZED, Message: err}})
	} else {
		c.JSON(http.StatusOK, Err(err))
	}
}

func ProtoResp[T RespType](origin *T) (*pb.ClientMessage, error) {
	event, msg, err := getEvent(origin)
	if err != nil {
		return nil, err
	}

	rawData, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return &pb.ClientMessage{
		Event: event,
		Data:  rawData,
	}, nil
}

func getEvent(origin interface{}) (pb.Event, protoreflect.ProtoMessage, error) {
	switch ptr := any(origin).(type) {
	case *pb.CommonResponse:
		return pb.Event_EVENT_DATA, ptr, nil
	case *pb.UpdateFRPCResponse:
		return pb.Event_EVENT_UPDATE_FRPC, ptr, nil
	case *pb.RemoveFRPCResponse:
		return pb.Event_EVENT_REMOVE_FRPC, ptr, nil
	case *pb.UpdateFRPSResponse:
		return pb.Event_EVENT_UPDATE_FRPC, ptr, nil
	case *pb.RemoveFRPSResponse:
		return pb.Event_EVENT_REMOVE_FRPC, ptr, nil
	case *pb.StartFRPCResponse:
		return pb.Event_EVENT_START_FRPC, ptr, nil
	case *pb.StopFRPCResponse:
		return pb.Event_EVENT_STOP_FRPC, ptr, nil
	case *pb.StartFRPSResponse:
		return pb.Event_EVENT_START_FRPS, ptr, nil
	case *pb.StopFRPSResponse:
		return pb.Event_EVENT_STOP_FRPS, ptr, nil
	case *pb.GetProxyConfigResponse:
		return pb.Event_EVENT_GET_PROXY_INFO, ptr, nil
	default:
		return 0, nil, fmt.Errorf("cannot unmarshal unknown type: %T", origin)
	}
}
