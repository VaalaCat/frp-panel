package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/dao"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/rpc"
	"github.com/VaalaCat/frp-panel/utils"
	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

func UpdateFrpcHander(c context.Context, req *pb.UpdateFRPCRequest) (*pb.UpdateFRPCResponse, error) {
	logrus.Infof("update frpc, req: [%+v]", req)
	var (
		content  = req.GetConfig()
		serverID = req.GetServerId()
		clientID = req.GetClientId()
		userInfo = common.GetUserInfo(c)
	)

	cliCfg, err := utils.LoadClientConfigNormal(content, true)
	if err != nil {
		logrus.WithError(err).Errorf("cannot load config")
		return &pb.UpdateFRPCResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: err.Error()},
		}, err
	}

	cli, err := dao.GetClientByClientID(userInfo, req.GetClientId())
	if err != nil {
		logrus.WithError(err).Errorf("cannot get client, id: [%s]", req.GetClientId())
		return &pb.UpdateFRPCResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "cannot get client"},
		}, fmt.Errorf("cannot get client")
	}

	srv, err := dao.GetServerByServerID(userInfo, req.GetServerId())
	if err != nil || srv == nil || len(srv.ServerIP) == 0 || len(srv.ConfigContent) == 0 {
		logrus.WithError(err).Errorf("cannot get server, server is not prepared, id: [%s]", req.GetServerId())
		return &pb.UpdateFRPCResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "cannot get server"},
		}, fmt.Errorf("cannot get server")
	}

	srvConf, err := srv.GetConfigContent()
	if srvConf == nil || err != nil {
		logrus.WithError(err).Errorf("cannot get server, id: [%s]", serverID)
		return nil, err
	}

	cliCfg.ServerAddr = srv.ServerIP
	cliCfg.ServerPort = srvConf.BindPort
	cliCfg.User = userInfo.GetUserName()
	cliCfg.Auth = v1.AuthClientConfig{}
	cliCfg.Metadatas = map[string]string{
		common.FRPAuthTokenKey: userInfo.GetToken(),
		common.FRPClientIDKey:  clientID,
	}

	newCfg := struct {
		v1.ClientCommonConfig
		Proxies  []v1.ProxyConfigurer   `json:"proxies,omitempty"`
		Visitors []v1.VisitorBaseConfig `json:"visitors,omitempty"`
	}{
		ClientCommonConfig: cliCfg.ClientCommonConfig,
		Proxies: lo.Map(cliCfg.Proxies, func(item v1.TypedProxyConfig, _ int) v1.ProxyConfigurer {
			return item.ProxyConfigurer
		}),
		Visitors: lo.Map(cliCfg.Visitors, func(item v1.TypedVisitorConfig, _ int) v1.VisitorBaseConfig {
			return *item.GetBaseConfig()
		}),
	}

	rawCliConf, err := json.Marshal(newCfg)
	if err != nil {
		logrus.WithError(err).Error("cannot marshal config")
		return nil, err
	}

	cli.ConfigContent = rawCliConf
	cli.ServerID = serverID
	cli.Comment = req.GetComment()

	if err := dao.UpdateClient(userInfo, cli); err != nil {
		logrus.WithError(err).Errorf("cannot update client, id: [%s]", clientID)
		return nil, err
	}

	cliReq := &pb.UpdateFRPCRequest{
		ClientId: lo.ToPtr(clientID),
		ServerId: lo.ToPtr(serverID),
		Config:   rawCliConf,
	}

	go func() {
		resp, err := rpc.CallClient(context.Background(), req.GetClientId(), pb.Event_EVENT_UPDATE_FRPC, cliReq)
		if err != nil {
			logrus.WithError(err).Errorf("update event send to client error, server: [%s], client: [%s]", serverID, req.GetClientId())
		}

		if resp == nil {
			logrus.Errorf("cannot get response, server: [%s], client: [%s]", serverID, req.GetClientId())
		}
	}()

	logrus.Infof("update frpc success, client id: [%s]", req.GetClientId())
	return &pb.UpdateFRPCResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
	}, nil
}
