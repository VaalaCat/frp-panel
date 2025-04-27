package client

import (
	"fmt"
	"net/url"

	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/models"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/dao"
	"github.com/VaalaCat/frp-panel/utils/logger"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/tiendc/go-deepcopy"
)

type ValidateableClientRequest interface {
	GetClientSecret() string
	GetClientId() string
}

func ValidateClientRequest(ctx *app.Context, req ValidateableClientRequest) (*models.ClientEntity, error) {
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

	if cli, err = dao.NewQuery(ctx).ValidateClientSecret(req.GetClientId(), req.GetClientSecret()); err != nil {
		return nil, err
	}

	return cli, nil
}

func MakeClientShadowed(c *app.Context, serverID string, clientEntity *models.ClientEntity) (*models.ClientEntity, error) {
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

		if err := dao.NewQuery(c).RebuildProxyConfigFromClient(userInfo, &models.Client{ClientEntity: childClient}); err != nil {
			logger.Logger(c).WithError(err).Errorf("cannot rebuild proxy config from client, id: [%s]", childClient.ClientID)
			return nil, err
		}
	}

	clientEntity.IsShadow = true
	clientEntity.ConfigContent = nil
	clientEntity.ServerID = ""
	if err := dao.NewQuery(c).UpdateClient(userInfo, clientEntity); err != nil {
		logger.Logger(c).WithError(err).Errorf("cannot update client, id: [%s]", clientID)
		return nil, err
	}

	if err := dao.NewQuery(c).DeleteProxyConfigsByClientID(userInfo, clientID); err != nil {
		logger.Logger(c).WithError(err).Errorf("cannot delete proxy configs, id: [%s]", clientID)
		return nil, err
	}

	return childClient, nil
}

// ChildClientForServer 支持传入serverID和任意类型client，返回serverID对应的client in shadow，如果不存在则新建
func ChildClientForServer(c *app.Context, serverID string, clientEntity *models.ClientEntity) (*models.ClientEntity, error) {
	userInfo := common.GetUserInfo(c)

	originClientID := clientEntity.ClientID
	if len(clientEntity.OriginClientID) != 0 {
		originClientID = clientEntity.OriginClientID
	}

	existClient, err := dao.NewQuery(c).GetClientByFilter(userInfo, &models.ClientEntity{
		ServerID:       serverID,
		OriginClientID: originClientID,
	}, lo.ToPtr(false))
	if err == nil {
		return existClient, nil
	}

	shadowCount, err := dao.NewQuery(c).CountClientsInShadow(userInfo, originClientID)
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
	copiedClient.ClientID = app.ShadowedClientID(originClientID, shadowCount+1)
	copiedClient.OriginClientID = originClientID
	copiedClient.IsShadow = false
	copiedClient.Stopped = false
	if err := dao.NewQuery(c).CreateClient(userInfo, copiedClient); err != nil {
		logger.Logger(c).WithError(err).Errorf("cannot create child client, id: [%s]", copiedClient.ClientID)
		return nil, err
	}

	return copiedClient, nil
}

func ValidFrpsScheme(scheme string) bool {
	return scheme == "tcp" ||
		scheme == "kcp" || scheme == "quic" ||
		scheme == "websocket" || scheme == "wss"
}

func ValidateFrpsUrl(urlStr string) (*url.URL, error) {
	parsedFrpsUrl, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("parse frps url error")
	}

	if !ValidFrpsScheme(parsedFrpsUrl.Scheme) {
		return nil, fmt.Errorf("invalid frps scheme")
	}

	if len(parsedFrpsUrl.Host) == 0 {
		return nil, fmt.Errorf("invalid frps host")
	}

	if len(parsedFrpsUrl.Hostname()) == 0 {
		return nil, fmt.Errorf("invalid frps hostname")
	}

	if cast.ToInt(parsedFrpsUrl.Port()) == 0 {
		return nil, fmt.Errorf("invalid frps port")
	}

	return parsedFrpsUrl, nil
}
