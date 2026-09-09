package client

import (
	"sync"

	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/dao"
	"github.com/VaalaCat/frp-panel/services/rpc"
	"github.com/VaalaCat/frp-panel/utils/logger"
)

func UpgradeFrppHandler(ctx *app.Context, req *pb.UpgradeFrppRequest) (*pb.UpgradeFrppResponse, error) {
	userInfo := common.GetUserInfo(ctx)
	clientIds := req.GetClientIds()
	log := logger.Logger(ctx)
	log.Infof("upgrade frpp called, user=%v clientIds=%v", userInfo, clientIds)

	if len(clientIds) == 0 {
		return &pb.UpgradeFrppResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "client_ids is empty"},
		}, nil
	}

	// 默认值处理：proto3 optional 未设置时 GetXXX 返回零值
	backup := true
	if req.Backup != nil {
		backup = req.GetBackup()
	}
	restartService := true
	if req.RestartService != nil {
		restartService = req.GetRestartService()
	}
	useGithubProxy := true
	if req.UseGithubProxy != nil {
		useGithubProxy = req.GetUseGithubProxy()
	}

	// 并发下发，提高多 client 批量升级速度
	var (
		wg      sync.WaitGroup
		errOnce error
		mu      sync.Mutex
	)

	for _, cid := range clientIds {
		clientId := cid
		wg.Add(1)
		go func() {
			defer wg.Done()

			_, err := dao.NewQuery(ctx).GetClientByClientID(userInfo, clientId)
			if err != nil {
				mu.Lock()
				if errOnce == nil {
					errOnce = err
				}
				mu.Unlock()
				return
			}

			// 每次只给单个 client 下发（减少 payload 混淆）
			reqForClient := &pb.UpgradeFrppRequest{
				ClientIds:      []string{clientId},
				Version:        req.Version,
				DownloadUrl:    req.DownloadUrl,
				GithubProxy:    req.GithubProxy,
				UseGithubProxy: &useGithubProxy,
				HttpProxy:      req.HttpProxy,
				TargetPath:     req.TargetPath,
				Backup:         &backup,
				ServiceName:    req.ServiceName,
				RestartService: &restartService,
				Workdir:        req.Workdir,
				ServiceArgs:    req.GetServiceArgs(),
			}

			resp := &pb.UpgradeFrppResponse{}
			if err := rpc.CallClientWrapper(ctx, clientId, pb.Event_EVENT_UPGRADE_FRPP, reqForClient, resp); err != nil {
				mu.Lock()
				if errOnce == nil {
					errOnce = err
				}
				mu.Unlock()
				return
			}
		}()
	}
	wg.Wait()

	if errOnce != nil {
		log.WithError(errOnce).Error("upgrade frpp dispatch failed")
		return nil, errOnce
	}

	return &pb.UpgradeFrppResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
	}, nil
}
