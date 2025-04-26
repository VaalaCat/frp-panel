package rpc

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"

	"github.com/VaalaCat/frp-panel/conf"
	"github.com/VaalaCat/frp-panel/defs"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/utils/logger"
	"github.com/VaalaCat/frp-panel/utils/wsgrpc"
	"github.com/imroc/req/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

type masterClient struct {
	cli         pb.MasterClient
	inited      bool
	appInstance app.Application
}

func (m *masterClient) Call() pb.MasterClient {
	if !m.inited {
		m.cli = newMasterCli(m.appInstance)
		m.inited = true
	}
	return m.cli
}

func NewMasterCli(appInstance app.Application) *masterClient {
	return &masterClient{
		inited:      false,
		appInstance: appInstance,
	}
}

func newMasterCli(appInstance app.Application) pb.MasterClient {
	connInfo := conf.GetRPCConnInfo(appInstance.GetConfig())
	ctx := context.Background()

	opt := []grpc.DialOption{}

	switch connInfo.Scheme {
	case conf.GRPC:
		if appInstance.GetConfig().Client.TLSRpc {
			logger.Logger(ctx).Infof("use tls rpc")
			opt = append(opt, grpc.WithTransportCredentials(appInstance.GetRPCCred()))
		} else {
			logger.Logger(ctx).Infof("use insecure rpc")
			opt = append(opt, grpc.WithTransportCredentials(insecure.NewCredentials()))
		}
	case conf.WS, conf.WSS:
		logger.Logger(ctx).Infof("use ws/wss rpc")

		wsURL := fmt.Sprintf("%s://%s/wsgrpc", connInfo.Scheme, connInfo.Host)
		header := http.Header{}
		wsDialer := wsgrpc.WebsocketDialer(wsURL, header, appInstance.GetConfig().Client.TLSInsecureSkipVerify)
		opt = append(opt, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(wsDialer))
	}

	conn, err := grpc.NewClient(connInfo.Host, opt...)

	if err != nil {
		logger.Logger(ctx).Fatalf("did not connect: %v", err)
	}

	return pb.NewMasterClient(conn)
}

func httpCli() *req.Client {
	c := req.C()
	c.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	return c
}

func GetClientCert(appInstance app.Application, clientID, clientSecret string, clientType pb.ClientType) []byte {
	apiEndpoint := conf.GetAPIURL(appInstance.GetConfig())
	c := httpCli()

	rawReq, err := proto.Marshal(&pb.GetClientCertRequest{
		ClientId:     clientID,
		ClientSecret: clientSecret,
		ClientType:   clientType,
	})
	if err != nil {
		return nil
	}
	r, err := c.R().SetHeader("Content-Type", "application/x-protobuf").
		SetBodyBytes(rawReq).Post(apiEndpoint + "/api/v1/auth/cert")
	if err != nil {
		return nil
	}

	resp := &pb.GetClientCertResponse{}
	err = proto.Unmarshal(r.Bytes(), resp)
	if err != nil {
		return nil
	}
	return resp.Cert
}

func InitClient(appInstance app.Application, clientID, joinToken string) (*pb.InitClientResponse, error) {
	apiEndpoint := conf.GetAPIURL(appInstance.GetConfig())

	c := httpCli()

	rawReq, err := proto.Marshal(&pb.InitClientRequest{
		ClientId: &clientID,
	})
	if err != nil {
		return nil, err
	}

	r, err := c.R().SetHeader("Content-Type", "application/x-protobuf").
		SetHeader(defs.AuthorizationKey, joinToken).
		SetBodyBytes(rawReq).Post(apiEndpoint + "/api/v1/client/init")
	if err != nil {
		return nil, err
	}

	resp := &pb.InitClientResponse{}
	err = proto.Unmarshal(r.Bytes(), resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GetClient(appInstance app.Application, clientID, joinToken string) (*pb.GetClientResponse, error) {
	apiEndpoint := conf.GetAPIURL(appInstance.GetConfig())
	c := httpCli()

	rawReq, err := proto.Marshal(&pb.GetClientRequest{
		ClientId: &clientID,
	})
	if err != nil {
		return nil, err
	}

	r, err := c.R().SetHeader("Content-Type", "application/x-protobuf").
		SetHeader(defs.AuthorizationKey, joinToken).
		SetBodyBytes(rawReq).Post(apiEndpoint + "/api/v1/client/get")
	if err != nil {
		return nil, err
	}

	resp := &pb.GetClientResponse{}
	err = proto.Unmarshal(r.Bytes(), resp)
	if err != nil {
		return nil, err
	}
	if resp.GetStatus().GetCode() != pb.RespCode_RESP_CODE_SUCCESS {
		return nil, errors.New(resp.GetStatus().GetMessage())
	}
	return resp, nil
}
