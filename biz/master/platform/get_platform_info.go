package platform

import (
	"fmt"

	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/conf"
	"github.com/VaalaCat/frp-panel/dao"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/gin-gonic/gin"
)

func GetPlatformInfo(c *gin.Context) {
	resp, err := getPlatformInfo(c)
	if err != nil {
		common.ErrResp(c, &pb.CommonResponse{Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: err.Error()}}, err.Error())
		return
	}
	common.OKResp(c, resp)
}

func getPlatformInfo(c *gin.Context) (*pb.GetPlatformInfoResponse, error) {
	userInfo := common.GetUserInfo(c)
	if !userInfo.Valid() {
		return &pb.GetPlatformInfoResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "invalid user"},
		}, nil
	}
	totalServers, err := dao.CountServers(userInfo)
	if err != nil {
		return nil, err
	}
	totalClients, err := dao.CountClients(userInfo)
	if err != nil {
		return nil, err
	}

	configuredServers, err := dao.CountConfiguredServers(userInfo)
	if err != nil {
		return nil, err
	}
	configuredClients, err := dao.CountConfiguredClients(userInfo)
	if err != nil {
		return nil, err
	}

	unconfiguredServers := totalServers - configuredServers

	unconfiguredClients := totalClients - configuredClients

	clientRPCUrl := conf.Get().Client.RPCUrl
	clientAPIUrl := conf.Get().Client.APIUrl

	if len(clientRPCUrl) == 0 {
		clientRPCUrl = fmt.Sprintf("grpc://%s:%d", conf.Get().Master.RPCHost, conf.Get().Master.RPCPort)
	}

	if len(clientAPIUrl) == 0 {
		clientAPIUrl = fmt.Sprintf("%s://%s:%d", conf.Get().Master.APIScheme, conf.Get().Master.RPCHost, conf.Get().Master.APIPort)
	}

	return &pb.GetPlatformInfoResponse{
		Status:                  &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
		TotalClientCount:        int32(totalClients),
		TotalServerCount:        int32(totalServers),
		UnconfiguredClientCount: int32(unconfiguredClients),
		UnconfiguredServerCount: int32(unconfiguredServers),
		ConfiguredClientCount:   int32(configuredClients),
		ConfiguredServerCount:   int32(configuredServers),
		GlobalSecret:            conf.MasterDefaultSalt(),
		MasterRpcHost:           conf.Get().Master.RPCHost,
		MasterRpcPort:           int32(conf.Get().Master.RPCPort),
		MasterApiPort:           int32(conf.Get().Master.APIPort),
		MasterApiScheme:         conf.Get().Master.APIScheme,
		ClientRpcUrl:            clientRPCUrl,
		ClientApiUrl:            clientAPIUrl,
	}, nil
}
