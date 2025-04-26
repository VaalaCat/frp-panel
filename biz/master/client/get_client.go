package client

import (
	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/models"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/dao"
	"github.com/VaalaCat/frp-panel/utils/logger"
	"github.com/samber/lo"
)

func GetClientHandler(ctx *app.Context, req *pb.GetClientRequest) (*pb.GetClientResponse, error) {
	logger.Logger(ctx).Infof("get client, req: [%+v]", req)

	var (
		userInfo = common.GetUserInfo(ctx)
		clientID = req.GetClientId()
		serverID = req.GetServerId()
	)

	if !userInfo.Valid() {
		return &pb.GetClientResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "invalid user"},
		}, nil
	}

	if len(clientID) == 0 {
		return &pb.GetClientResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "invalid client id"},
		}, nil
	}

	respCli := &pb.Client{}
	if len(serverID) == 0 {
		client, err := dao.NewQuery(ctx).GetClientByClientID(userInfo, clientID)
		if err != nil {
			return nil, err
		}
		clientIDs, err := dao.NewQuery(ctx).GetClientIDsInShadowByClientID(userInfo, clientID)
		if err != nil {
			logger.Logger(ctx).WithError(err).Errorf("cannot get client ids in shadow, id: [%s]", clientID)
		}

		respCli = &pb.Client{
			Id:        lo.ToPtr(client.ClientID),
			Secret:    lo.ToPtr(client.ConnectSecret),
			Config:    lo.ToPtr(string(client.ConfigContent)),
			ServerId:  lo.ToPtr(client.ServerID),
			Stopped:   lo.ToPtr(client.Stopped),
			Comment:   lo.ToPtr(client.Comment),
			ClientIds: clientIDs,
		}
	} else {
		client, err := dao.NewQuery(ctx).GetClientByFilter(userInfo, &models.ClientEntity{
			OriginClientID: clientID,
			ServerID:       serverID,
		}, lo.ToPtr(false))
		if err != nil {
			client, err = dao.NewQuery(ctx).GetClientByFilter(userInfo, &models.ClientEntity{
				ClientID: clientID,
				ServerID: serverID,
			}, nil)
			if err != nil {
				return nil, err
			}
		}

		respCli = &pb.Client{
			Id:        lo.ToPtr(client.ClientID),
			Secret:    lo.ToPtr(client.ConnectSecret),
			Config:    lo.ToPtr(string(client.ConfigContent)),
			ServerId:  lo.ToPtr(client.ServerID),
			Stopped:   lo.ToPtr(client.Stopped),
			Comment:   lo.ToPtr(client.Comment),
			ClientIds: nil,
		}
	}

	return &pb.GetClientResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
		Client: respCli,
	}, nil
}
