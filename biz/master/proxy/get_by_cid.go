package proxy

import (
	"context"
	"fmt"

	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/dao"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/sirupsen/logrus"
)

// GetProxyByCID get proxy info by client id
func GetProxyByCID(c context.Context, req *pb.GetProxyByCIDRequest) (*pb.GetProxyByCIDResponse, error) {
	logrus.Infof("get proxy by client id, req: [%+v]", req)
	var (
		clientID = req.GetClientId()
		userInfo = common.GetUserInfo(c)
	)

	if len(clientID) == 0 {
		return nil, fmt.Errorf("request invalid")
	}

	proxyList, err := dao.GetProxyByClientID(userInfo, clientID)
	if proxyList == nil || err != nil {
		logrus.WithError(err).Errorf("cannot get proxy, client id: [%s]", clientID)
		return nil, err
	}
	return &pb.GetProxyByCIDResponse{
		Status:     &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
		ProxyInfos: convertProxyList(proxyList),
	}, nil
}
