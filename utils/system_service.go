package utils

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/VaalaCat/frp-panel/utils/logger"
	"github.com/kardianos/service"
)

type SystemService struct {
	run func()
	service.Service
}

func (ss *SystemService) Start(s service.Service) error {
	go ss.iRun()
	return nil
}

func (ss *SystemService) Stop(s service.Service) error { return nil }

func (ss *SystemService) iRun() {
	defer func() {
		if service.Interactive() {
			ss.Stop(ss.Service)
		} else {
			ss.Service.Stop()
		}
	}()
	ss.run()
}

func CreateSystemService(args []string, run func()) (service.Service, error) {
	currentPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("get current path failed, err: %v", err)
	}

	svcConfig := &service.Config{
		Name:             "frpp",
		DisplayName:      "frp-panel",
		Description:      "this is frp-panel service, developed by [VaalaCat] - https://github.com/VaalaCat/frp-panel",
		Arguments:        args,
		WorkingDirectory: path.Dir(currentPath),
	}

	ss := &SystemService{
		run: run,
	}

	s, err := service.New(ss, svcConfig)
	if err != nil {
		return nil, fmt.Errorf("service New failed, err: %v", err)
	}
	return s, nil
}

func ControlSystemService(args []string, action string, run func()) error {
	ctx := context.Background()

	logger.Logger(ctx).Info("try to ", action, " service, args:", args)
	s, err := CreateSystemService(args, run)
	if err != nil {
		logger.Logger(ctx).WithError(err).Error("create service controller failed")
		return err
	}

	if err := service.Control(s, action); err != nil {
		logger.Logger(ctx).WithError(err).Errorf("controller %v service failed", action)
		return err
	}
	logger.Logger(ctx).Infof("controller %v service success", action)
	return nil
}

func InstallToSystemPath(installPath string) error {
	currentPath, err := os.Executable()
	if err != nil {
		return err
	}

	targetPath := path.Join(installPath, filepath.Base(currentPath))

	src, err := os.Open(currentPath)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	err = os.Chmod(targetPath, 0755)
	if err != nil {
		return err
	}

	return nil
}
