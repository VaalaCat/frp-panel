package proxy

import (
	"context"
	"fmt"

	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/dao"
	"github.com/VaalaCat/frp-panel/utils/logger"
)

// GetProxyStatsByServerID get proxy info by server id
func GetProxyStatsByServerID(c *app.Context, req *pb.GetProxyStatsByServerIDRequest) (*pb.GetProxyStatsByServerIDResponse, error) {
	logger.Logger(c).Infof("get proxy by server id, req: [%+v]", req)
	var (
		serverID = req.GetServerId()
		userInfo = common.GetUserInfo(c)
	)

	if len(serverID) == 0 {
		return nil, fmt.Errorf("request invalid")
	}

	proxyList, err := dao.NewQuery(c).GetProxyStatsByServerID(userInfo, serverID)
	if proxyList == nil || err != nil {
		logger.Logger(context.Background()).WithError(err).Errorf("cannot get proxy, server id: [%s]", serverID)
		return nil, err
	}
	return &pb.GetProxyStatsByServerIDResponse{
		Status:     &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
		ProxyInfos: convertProxyStatsList(proxyList),
	}, nil
}
