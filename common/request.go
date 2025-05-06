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
		pb.GetProxyStatsByClientIDRequest | pb.GetProxyStatsByServerIDRequest |
		pb.CreateProxyConfigRequest | pb.ListProxyConfigsRequest | pb.UpdateProxyConfigRequest |
		pb.DeleteProxyConfigRequest | pb.GetProxyConfigRequest | pb.SignTokenRequest |
		pb.StartProxyRequest | pb.StopProxyRequest |
		pb.CreateWorkerRequest | pb.RemoveWorkerRequest | pb.RunWorkerRequest | pb.StopWorkerRequest | pb.UpdateWorkerRequest | pb.GetWorkerRequest |
		pb.ListWorkersRequest | pb.CreateWorkerIngressRequest | pb.GetWorkerIngressRequest |
		pb.GetWorkerStatusRequest | pb.InstallWorkerdRequest |
		pb.StartSteamLogRequest
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
	msg, ok := any(r).(protoreflect.ProtoMessage)
	if !ok {
		return fmt.Errorf("type does not implement protoreflect.ProtoMessage")
	}

	err = trans(b, msg)
	if err != nil {
		return err
	}
	return nil
}
