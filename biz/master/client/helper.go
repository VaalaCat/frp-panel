package client

import (
	"context"
	"fmt"

	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/dao"
	"github.com/VaalaCat/frp-panel/logger"
	"github.com/VaalaCat/frp-panel/models"
	"github.com/samber/lo"
	"github.com/tiendc/go-deepcopy"
)

type ValidateableClientRequest interface {
	GetClientSecret() string
	GetClientId() string
}

func ValidateClientRequest(req ValidateableClientRequest) (*models.ClientEntity, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request")
	}

	if req.GetClientId() == "" || req.GetClientSecret() == "" {
		return nil, fmt.Errorf("invalid request")
	}

	var (
		cli *models.ClientEntity
		err error
	)

	if cli, err = dao.ValidateClientSecret(req.GetClientId(), req.GetClientSecret()); err != nil {
		return nil, err
	}

	return cli, nil
}

func MakeClientShadowed(c context.Context, serverID string, clientEntity *models.ClientEntity) (*models.ClientEntity, error) {
	userInfo := common.GetUserInfo(c)

	var clientID = clientEntity.ClientID
	var childClient *models.ClientEntity
	var err error
	if len(clientEntity.ConfigContent) != 0 {
		childClient, err = ChildClientForServer(c, serverID, clientEntity)
		if err != nil {
			logger.Logger(c).WithError(err).Errorf("cannot create child client, id: [%s]", clientID)
			return nil, err
		}

		if err := dao.RebuildProxyConfigFromClient(userInfo, &models.Client{ClientEntity: childClient}); err != nil {
			logger.Logger(c).WithError(err).Errorf("cannot rebuild proxy config from client, id: [%s]", childClient.ClientID)
			return nil, err
		}
	}

	clientEntity.IsShadow = true
	clientEntity.ConfigContent = nil
	clientEntity.ServerID = ""
	if err := dao.UpdateClient(userInfo, clientEntity); err != nil {
		logger.Logger(c).WithError(err).Errorf("cannot update client, id: [%s]", clientID)
		return nil, err
	}

	if err := dao.DeleteProxyConfigsByClientID(userInfo, clientID); err != nil {
		logger.Logger(c).WithError(err).Errorf("cannot delete proxy configs, id: [%s]", clientID)
		return nil, err
	}

	return childClient, nil
}

// ChildClientForServer 支持传入serverID和任意类型client，返回serverID对应的client in shadow，如果不存在则新建
func ChildClientForServer(c context.Context, serverID string, clientEntity *models.ClientEntity) (*models.ClientEntity, error) {
	userInfo := common.GetUserInfo(c)

	originClientID := clientEntity.ClientID
	if len(clientEntity.OriginClientID) != 0 {
		originClientID = clientEntity.OriginClientID
	}

	existClient, err := dao.GetClientByFilter(userInfo, &models.ClientEntity{
		ServerID:       serverID,
		OriginClientID: originClientID,
	}, lo.ToPtr(false))
	if err == nil {
		return existClient, nil
	}

	shadowCount, err := dao.CountClientsInShadow(userInfo, originClientID)
	if err != nil {
		logger.Logger(c).WithError(err).Errorf("cannot count shadow clients, id: [%s]", originClientID)
		return nil, err
	}

	copiedClient := &models.ClientEntity{}
	if err := deepcopy.Copy(copiedClient, clientEntity); err != nil {
		logger.Logger(c).WithError(err).Errorf("cannot copy client, id: [%s]", originClientID)
		return nil, err
	}
	copiedClient.ServerID = serverID
	copiedClient.ClientID = common.ShadowedClientID(originClientID, shadowCount+1)
	copiedClient.OriginClientID = originClientID
	copiedClient.IsShadow = false
	copiedClient.Stopped = false
	if err := dao.CreateClient(userInfo, copiedClient); err != nil {
		logger.Logger(c).WithError(err).Errorf("cannot create child client, id: [%s]", copiedClient.ClientID)
		return nil, err
	}

	return copiedClient, nil
}
