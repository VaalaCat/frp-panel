package proxy

import (
	"context"
	"fmt"

	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/dao"
	"github.com/VaalaCat/frp-panel/logger"
	"github.com/VaalaCat/frp-panel/pb"
)

// GetProxyBySID get proxy info by server id
func GetProxyBySID(c context.Context, req *pb.GetProxyBySIDRequest) (*pb.GetProxyBySIDResponse, error) {
	logger.Logger(c).Infof("get proxy by server id, req: [%+v]", req)
	var (
		serverID = req.GetServerId()
		userInfo = common.GetUserInfo(c)
	)

	if len(serverID) == 0 {
		return nil, fmt.Errorf("request invalid")
	}

	proxyList, err := dao.GetProxyByServerID(userInfo, serverID)
	if proxyList == nil || err != nil {
		logger.Logger(context.Background()).WithError(err).Errorf("cannot get proxy, server id: [%s]", serverID)
		return nil, err
	}
	return &pb.GetProxyBySIDResponse{
		Status:     &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
		ProxyInfos: convertProxyList(proxyList),
	}, nil
}
