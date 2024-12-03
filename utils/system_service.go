package utils

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/kardianos/service"
	"github.com/sirupsen/logrus"
)

type SystemService struct{}

func (ss *SystemService) Start(s service.Service) error { return nil }

func (ss *SystemService) Stop(s service.Service) error { return nil }

func createSystemService(args []string) (service.Service, error) {
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

	ss := &SystemService{}
	s, err := service.New(ss, svcConfig)
	if err != nil {
		return nil, fmt.Errorf("service New failed, err: %v", err)
	}
	return s, nil
}

func ControlSystemService(args []string, action string) error {
	logrus.Info("try to ", action, " service, args:", args)
	s, err := createSystemService(args)
	if err != nil {
		logrus.WithError(err).Error("create service controller failed")
		return err
	}

	if err := service.Control(s, action); err != nil {
		logrus.WithError(err).Errorf("controller %v service failed", action)
		return err
	}
	logrus.Infof("controller %v service success", action)
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
