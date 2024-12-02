package platform

import (
	"context"
	"time"

	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/logger"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/rpc"
)

func GetClientsStatus(c context.Context, req *pb.GetClientsStatusRequest) (*pb.GetClientsStatusResponse, error) {
	userInfo := common.GetUserInfo(c)
	if !userInfo.Valid() || req == nil || len(req.GetClientIds()) == 0 || req.GetClientType() == pb.ClientType_CLIENT_TYPE_UNSPECIFIED {
		return &pb.GetClientsStatusResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "request invalid"},
		}, nil
	}

	var (
		clientIDs = req.GetClientIds()
		resps     = map[string]*pb.ClientStatus{}
	)

	for _, clientID := range clientIDs {
		conn := rpc.GetClientsManager().Get(clientID)
		if conn == nil {
			resps[clientID] = &pb.ClientStatus{
				ClientType: req.GetClientType(),
				ClientId:   clientID,
				Status:     pb.ClientStatus_STATUS_OFFLINE,
				Ping:       -1,
			}
			continue
		}
		startTime := time.Now()
		tresp, err := rpc.CallClient(c, clientID, pb.Event_EVENT_PING, &pb.CommonRequest{})
		endTime := time.Now()
		pingTime := endTime.Sub(startTime).Milliseconds()
		if err != nil || tresp == nil {
			logger.Logger(context.Background()).WithError(err).Errorf("get client status error, client id: [%s]", clientID)
			resps[clientID] = &pb.ClientStatus{
				ClientType: req.GetClientType(),
				ClientId:   clientID,
				Status:     pb.ClientStatus_STATUS_ERROR,
				Ping:       int32(pingTime),
			}
			continue
		}
		resps[clientID] = &pb.ClientStatus{
			ClientType: req.GetClientType(),
			ClientId:   clientID,
			Status:     pb.ClientStatus_STATUS_ONLINE,
			Ping:       int32(pingTime),
		}
	}

	return &pb.GetClientsStatusResponse{
		Status:  &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
		Clients: resps,
	}, nil
}
