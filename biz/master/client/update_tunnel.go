package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/VaalaCat/frp-panel/app"
	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/dao"
	"github.com/VaalaCat/frp-panel/defs"
	"github.com/VaalaCat/frp-panel/logger"
	"github.com/VaalaCat/frp-panel/models"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/rpc"
	"github.com/VaalaCat/frp-panel/utils"
	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/samber/lo"
)

func UpdateFrpcHander(c *app.Context, req *pb.UpdateFRPCRequest) (*pb.UpdateFRPCResponse, error) {
	logger.Logger(c).Infof("update frpc, req: [%+v]", req)
	var (
		content     = req.GetConfig()
		serverID    = req.GetServerId()
		reqClientID = req.GetClientId() // may be shadow or child
		userInfo    = common.GetUserInfo(c)
	)

	cliCfg, err := utils.LoadClientConfigNormal(content, true)
	if err != nil {
		logger.Logger(c).WithError(err).Errorf("cannot load config")
		return &pb.UpdateFRPCResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: err.Error()},
		}, err
	}

	cli, err := dao.NewQuery(c).GetClientByClientID(userInfo, reqClientID)
	if err != nil {
		logger.Logger(c).WithError(err).Errorf("cannot get client, id: [%s]", reqClientID)
		return &pb.UpdateFRPCResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "cannot get client"},
		}, fmt.Errorf("cannot get client")
	}

	if cli.IsShadow {
		cli, err = ChildClientForServer(c, serverID, cli)
		if err != nil {
			logger.Logger(c).WithError(err).Errorf("cannot get child client, id: [%s]", reqClientID)
			return &pb.UpdateFRPCResponse{
				Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "cannot get child client"},
			}, fmt.Errorf("cannot get child client")
		}
	}

	if cli.IsShadow && len(cli.ConfigContent) == 0 {
		logger.Logger(c).Warnf("client is shadowed, cannot update, id: [%s]", reqClientID)
		return nil, fmt.Errorf("client is shadowed, cannot update")
	}

	if !cli.IsShadow && len(cli.OriginClientID) == 0 {
		logger.Logger(c).Warnf("client is not shadowed, make it shadow, id: [%s]", reqClientID)
		cli, err = MakeClientShadowed(c, serverID, cli)
		if err != nil {
			logger.Logger(c).WithError(err).Errorf("cannot make client shadowed, id: [%s]", reqClientID)
			return &pb.UpdateFRPCResponse{
				Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "cannot make client shadowed"},
			}, fmt.Errorf("cannot make client shadowed")
		}
	}

	srv, err := dao.NewQuery(c).GetServerByServerID(userInfo, req.GetServerId())
	if err != nil || srv == nil || len(srv.ServerIP) == 0 || len(srv.ConfigContent) == 0 {
		logger.Logger(c).WithError(err).Errorf("cannot get server, server is not prepared, id: [%s]", req.GetServerId())
		return &pb.UpdateFRPCResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "cannot get server"},
		}, fmt.Errorf("cannot get server")
	}

	srvConf, err := srv.GetConfigContent()
	if srvConf == nil || err != nil {
		logger.Logger(c).WithError(err).Errorf("cannot get server, id: [%s]", serverID)
		return nil, err
	}

	cliCfg.ServerAddr = srv.ServerIP
	switch cliCfg.Transport.Protocol {
	case "tcp":
		cliCfg.ServerPort = srvConf.BindPort
	case "kcp":
		cliCfg.ServerPort = srvConf.KCPBindPort
	case "quic":
		cliCfg.ServerPort = srvConf.QUICBindPort
	default:
		cliCfg.ServerPort = srvConf.BindPort
	}

	cliCfg.User = userInfo.GetUserName()

	if cliCfg.Metadatas == nil {
		cliCfg.Metadatas = make(map[string]string)
	}

	cliCfg.Metadatas[defs.FRPAuthTokenKey] = userInfo.GetToken()
	cliCfg.Metadatas[defs.FRPClientIDKey] = reqClientID

	newCfg := struct {
		v1.ClientCommonConfig
		Proxies  []v1.ProxyConfigurer   `json:"proxies,omitempty"`
		Visitors []v1.VisitorConfigurer `json:"visitors,omitempty"`
	}{
		ClientCommonConfig: cliCfg.ClientCommonConfig,
		Proxies:            lo.Map(cliCfg.Proxies, func(item v1.TypedProxyConfig, _ int) v1.ProxyConfigurer { return item.ProxyConfigurer }),
		Visitors:           lo.Map(cliCfg.Visitors, func(item v1.TypedVisitorConfig, _ int) v1.VisitorConfigurer { return item.VisitorConfigurer }),
	}

	rawCliConf, err := json.Marshal(newCfg)
	if err != nil {
		logger.Logger(c).WithError(err).Error("cannot marshal config")
		return nil, err
	}

	cli.ConfigContent = rawCliConf
	cli.ServerID = serverID
	if req.Comment != nil {
		cli.Comment = req.GetComment()
	}

	if err := dao.NewQuery(c).UpdateClient(userInfo, cli); err != nil {
		logger.Logger(c).WithError(err).Errorf("cannot update client, id: [%s]", cli.ClientID)
		return nil, err
	}

	if err := dao.NewQuery(c).RebuildProxyConfigFromClient(userInfo, &models.Client{ClientEntity: cli}); err != nil {
		logger.Logger(c).WithError(err).Errorf("cannot rebuild proxy config from client, id: [%s]", cli.ClientID)
		return nil, err
	}

	cliReq := &pb.UpdateFRPCRequest{
		ClientId: lo.ToPtr(cli.ClientID),
		ServerId: lo.ToPtr(serverID),
		Config:   rawCliConf,
	}

	go func() {
		childCtx := app.NewContext(context.Background(), c.GetApp())
		cliToUpdate, err := dao.NewQuery(childCtx).GetClientByFilter(userInfo, &models.ClientEntity{ClientID: cli.OriginClientID}, nil)
		if err != nil {
			logger.Logger(childCtx).WithError(err).Errorf("cannot get origin client, id: [%s]", cliToUpdate.OriginClientID)
			return
		}

		if cliToUpdate.Stopped {
			logger.Logger(childCtx).Infof("client [%s] is stopped, do not send update event", cliToUpdate.OriginClientID)
			return
		}

		resp, err := rpc.CallClient(childCtx, cliToUpdate.ClientID, pb.Event_EVENT_UPDATE_FRPC, cliReq)
		if err != nil {
			logger.Logger(childCtx).WithError(err).Errorf("update event send to client error, server: [%s], client: [%+v], updated client: [%+v]", serverID, cliToUpdate, cli)
		}

		if resp == nil {
			logger.Logger(childCtx).Errorf("cannot get response, server: [%s], client: [%s]", serverID, cliToUpdate.OriginClientID)
		}
	}()

	logger.Logger(c).Infof("update frpc success, client id: [%s]", reqClientID)
	return &pb.UpdateFRPCResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
	}, nil
}
