package common

import (
	"fmt"

	"github.com/VaalaCat/frp-panel/pb"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type ReqType interface {
	pb.UpdateFRPCRequest | pb.RemoveFRPCRequest |
		pb.UpdateFRPSRequest | pb.RemoveFRPSRequest |
		pb.CommonRequest | pb.RegisterRequest | pb.LoginRequest |
		pb.InitClientRequest | pb.ListClientsRequest | pb.GetClientRequest |
		pb.DeleteClientRequest |
		pb.InitServerRequest | pb.ListServersRequest | pb.GetServerRequest |
		pb.DeleteServerRequest |
		pb.GetUserInfoRequest | pb.UpdateUserInfoRequest |
		pb.GetPlatformInfoRequest | pb.GetClientsStatusRequest |
		pb.GetClientCertRequest |
		pb.StartFRPCRequest | pb.StopFRPCRequest | pb.StartFRPSRequest | pb.StopFRPSRequest |
		pb.GetProxyByCIDRequest | pb.GetProxyBySIDRequest
}

func GetProtoRequest[T ReqType](c *gin.Context) (r *T, err error) {
	r = new(T)
	if c.ContentType() == "application/x-protobuf" {
		err = c.Copy().ShouldBindWith(r, binding.ProtoBuf)
		if err != nil {
			return nil, err
		}
	} else {
		b, err := c.Copy().GetRawData()
		if err != nil {
			return nil, err
		}

		err = GetServerMessageRequest(b, r, protojson.Unmarshal)
		if err != nil {
			return nil, err
		}
	}
	return r, nil
}

func GetServerMessageRequest[T ReqType](b []byte, r *T, trans func(b []byte, m protoreflect.ProtoMessage) error) (err error) {
	switch ptr := any(r).(type) {
	case *pb.CommonRequest:
		return trans(b, ptr)
	case *pb.UpdateFRPCRequest:
		return trans(b, ptr)
	case *pb.RemoveFRPCRequest:
		return trans(b, ptr)
	case *pb.UpdateFRPSRequest:
		return trans(b, ptr)
	case *pb.RemoveFRPSRequest:
		return trans(b, ptr)
	case *pb.RegisterRequest:
		return trans(b, ptr)
	case *pb.LoginRequest:
		return trans(b, ptr)
	case *pb.InitClientRequest:
		return trans(b, ptr)
	case *pb.ListClientsRequest:
		return trans(b, ptr)
	case *pb.GetClientRequest:
		return trans(b, ptr)
	case *pb.DeleteClientRequest:
		return trans(b, ptr)
	case *pb.InitServerRequest:
		return trans(b, ptr)
	case *pb.ListServersRequest:
		return trans(b, ptr)
	case *pb.GetServerRequest:
		return trans(b, ptr)
	case *pb.DeleteServerRequest:
		return trans(b, ptr)
	case *pb.GetUserInfoRequest:
		return trans(b, ptr)
	case *pb.UpdateUserInfoRequest:
		return trans(b, ptr)
	case *pb.GetPlatformInfoRequest:
		return trans(b, ptr)
	case *pb.GetClientsStatusRequest:
		return trans(b, ptr)
	case *pb.GetClientCertRequest:
		return trans(b, ptr)
	case *pb.StartFRPCRequest:
		return trans(b, ptr)
	case *pb.StopFRPCRequest:
		return trans(b, ptr)
	case *pb.StartFRPSRequest:
		return trans(b, ptr)
	case *pb.StopFRPSRequest:
		return trans(b, ptr)
	case *pb.GetProxyByCIDRequest:
		return trans(b, ptr)
	case *pb.GetProxyBySIDRequest:
		return trans(b, ptr)
	default:
	}
	return fmt.Errorf("cannot unmarshal unknown type: %T", r)
}
