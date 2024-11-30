package dao

import (
	"fmt"
	"strings"
	"time"

	"github.com/VaalaCat/frp-panel/models"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/utils"
	"github.com/samber/lo"
	"github.com/sourcegraph/conc/pool"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetProxyByClientID(userInfo models.UserInfo, clientID string) ([]*models.ProxyEntity, error) {
	if clientID == "" {
		return nil, fmt.Errorf("invalid client id")
	}
	db := models.GetDBManager().GetDefaultDB()
	list := []*models.Proxy{}
	err := db.
		Where(&models.Proxy{ProxyEntity: &models.ProxyEntity{
			UserID:   userInfo.GetUserID(),
			TenantID: userInfo.GetTenantID(),
			ClientID: clientID,
		}}).
		Or(&models.Proxy{ProxyEntity: &models.ProxyEntity{
			UserID:   0,
			TenantID: userInfo.GetTenantID(),
			ClientID: clientID,
		}}).
		Find(&list).Error
	if err != nil {
		return nil, err
	}
	return lo.Map(list, func(item *models.Proxy, _ int) *models.ProxyEntity {
		return item.ProxyEntity
	}), nil
}

func GetProxyByServerID(userInfo models.UserInfo, serverID string) ([]*models.ProxyEntity, error) {
	if serverID == "" {
		return nil, fmt.Errorf("invalid server id")
	}
	db := models.GetDBManager().GetDefaultDB()
	list := []*models.Proxy{}
	err := db.
		Where(&models.Proxy{ProxyEntity: &models.ProxyEntity{
			UserID:   userInfo.GetUserID(),
			TenantID: userInfo.GetTenantID(),
			ServerID: serverID,
		}}).Or(&models.Proxy{ProxyEntity: &models.ProxyEntity{
		UserID:   0,
		TenantID: userInfo.GetTenantID(),
		ServerID: serverID,
	}}).
		Find(&list).Error
	if err != nil {
		return nil, err
	}
	return lo.Map(list, func(item *models.Proxy, _ int) *models.ProxyEntity {
		return item.ProxyEntity
	}), nil
}

func AdminUpdateProxy(srv *models.ServerEntity, inputs []*pb.ProxyInfo) error {
	if srv.ServerID == "" {
		return fmt.Errorf("invalid server id")
	}

	db := models.GetDBManager().GetDefaultDB()
	return db.Transaction(func(tx *gorm.DB) error {

		queryResults := make([]interface{}, 3)
		p := pool.New().WithErrors()
		p.Go(
			func() error {
				user := models.User{}
				if err := tx.Where(&models.User{
					UserEntity: &models.UserEntity{
						UserID: srv.UserID,
					},
				}).First(&user).Error; err != nil {
					return err
				}
				queryResults[0] = user
				return nil
			},
		)
		p.Go(
			func() error {
				clients := []*models.Client{}
				if err := tx.
					Where(&models.Client{ClientEntity: &models.ClientEntity{
						UserID:   srv.UserID,
						ServerID: srv.ServerID,
					}}).Find(&clients).Error; err != nil {
					return err
				}
				queryResults[1] = clients
				return nil
			},
		)
		p.Go(
			func() error {
				oldProxy := []*models.Proxy{}
				if err := tx.
					Where(&models.Proxy{ProxyEntity: &models.ProxyEntity{
						UserID:   srv.UserID,
						ServerID: srv.ServerID,
					}}).Find(&oldProxy).Error; err != nil {
					return err
				}
				oldProxyMap := lo.SliceToMap(oldProxy, func(p *models.Proxy) (string, *models.Proxy) {
					return p.Name, p
				})
				queryResults[2] = oldProxyMap
				return nil
			},
		)
		if err := p.Wait(); err != nil {
			return err
		}

		user := queryResults[0].(models.User)
		clients := queryResults[1].([]*models.Client)
		oldProxyMap := queryResults[2].(map[string]*models.Proxy)

		inputMap := map[string]*pb.ProxyInfo{}
		proxyMap := map[string]*models.ProxyEntity{}
		for _, proxyInfo := range inputs {
			if proxyInfo == nil {
				continue
			}
			proxyName := strings.TrimPrefix(proxyInfo.GetName(), user.UserName+".")
			proxyMap[proxyName] = &models.ProxyEntity{
				ServerID:        srv.ServerID,
				Name:            proxyName,
				Type:            proxyInfo.GetType(),
				UserID:          srv.UserID,
				TenantID:        srv.TenantID,
				TodayTrafficIn:  proxyInfo.GetTodayTrafficIn(),
				TodayTrafficOut: proxyInfo.GetTodayTrafficOut(),
			}
			inputMap[proxyName] = proxyInfo
		}

		proxyEntityMap := map[string]*models.ProxyEntity{}
		for _, client := range clients {
			cliCfg, err := client.GetConfigContent()
			if err != nil || cliCfg == nil {
				continue
			}
			for _, cfg := range cliCfg.Proxies {
				if proxy, ok := proxyMap[cfg.GetBaseConfig().Name]; ok {
					proxy.ClientID = client.ClientID
					proxyEntityMap[proxy.Name] = proxy
				}
			}
		}

		nowTime := time.Now()
		results := lo.Values(lo.MapValues(proxyEntityMap, func(p *models.ProxyEntity, name string) *models.Proxy {
			item := &models.Proxy{
				ProxyEntity: p,
			}
			if oldProxy, ok := oldProxyMap[name]; ok {
				item.ProxyID = oldProxy.ProxyID
				firstSync := inputMap[name].GetFirstSync()
				isSameDay := utils.IsSameDay(nowTime, oldProxy.UpdatedAt)

				item.HistoryTrafficIn = oldProxy.HistoryTrafficIn
				item.HistoryTrafficOut = oldProxy.HistoryTrafficOut
				if !isSameDay || firstSync {
					item.HistoryTrafficIn += oldProxy.TodayTrafficIn
					item.HistoryTrafficOut += oldProxy.TodayTrafficOut
				}
			}
			return item
		}))

		if len(results) > 0 {
			return tx.Save(results).Error
		}
		return nil
	})
}

func AdminGetTenantProxies(tenantID int) ([]*models.ProxyEntity, error) {
	db := models.GetDBManager().GetDefaultDB()
	list := []*models.Proxy{}
	err := db.
		Where(&models.Proxy{ProxyEntity: &models.ProxyEntity{
			TenantID: tenantID,
		}}).
		Find(&list).Error
	if err != nil {
		return nil, err
	}
	return lo.Map(list, func(item *models.Proxy, _ int) *models.ProxyEntity {
		return item.ProxyEntity
	}), nil
}

func AdminGetAllProxies(tx *gorm.DB) ([]*models.ProxyEntity, error) {
	db := tx
	list := []*models.Proxy{}
	err := db.Clauses(clause.Locking{Strength: "UPDATE"}).
		Find(&list).Error
	if err != nil {
		return nil, err
	}
	return lo.Map(list, func(item *models.Proxy, _ int) *models.ProxyEntity {
		return item.ProxyEntity
	}), nil
}
