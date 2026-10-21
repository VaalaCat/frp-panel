package rpc

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/VaalaCat/frp-panel/conf"
	"github.com/VaalaCat/frp-panel/defs"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/utils"
	"github.com/VaalaCat/frp-panel/utils/logger"
	"github.com/VaalaCat/frp-panel/utils/wsgrpc"
	"github.com/imroc/req/v3"
	"github.com/samber/lo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/protobuf/proto"
)

type masterClient struct {
	cli         pb.MasterClient
	once        sync.Once
	appInstance app.Application
}

func (m *masterClient) Call() pb.MasterClient {
	m.once.Do(func() { m.cli = newMasterCli(m.appInstance) })
	return m.cli
}

func NewMasterCli(appInstance app.Application) *masterClient {
	logger.Logger(context.Background()).Debugf("creating new master client")
	return &masterClient{
		appInstance: appInstance,
	}
}

func newMasterCli(appInstance app.Application) pb.MasterClient {
	connInfo := conf.GetRPCConnInfo(appInstance.GetConfig())
	ctx := context.Background()

	opt := []grpc.DialOption{grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoke grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if _, bounded := ctx.Deadline(); !bounded {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
			defer cancel()
		}
		return invoke(ctx, method, req, reply, cc, opts...)
	}), grpc.WithKeepaliveParams(keepalive.ClientParameters{Time: 30 * time.Second, Timeout: 10 * time.Second, PermitWithoutStream: true})}

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
		wsDialer := wsgrpc.WebsocketDialer(wsURL,
			header,
			appInstance.GetConfig().Client.TLSInsecureSkipVerify,
			logger.Logger(ctx),
		)
		opt = append(opt, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(wsDialer))
	}

	logger.Logger(ctx).Debugf("creating new grpc client to [%s]", utils.MarshalForJson(connInfo))
	conn, err := grpc.NewClient(connInfo.Host, opt...)

	if err != nil {
		logger.Logger(ctx).Fatalf("did not connect: %v", err)
	}

	logger.Logger(ctx).Debugf("grpc client created")

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

func InitClient(cfg conf.Config, clientID, joinToken string, ephemeral *bool) (*pb.InitClientResponse, error) {
	apiEndpoint := conf.GetAPIURL(cfg)

	c := httpCli()

	if ephemeral == nil {
		ephemeral = lo.ToPtr(true) // default to ephemeral
	}

	rawReq, err := proto.Marshal(&pb.InitClientRequest{
		ClientId:  &clientID,
		Ephemeral: ephemeral,
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

func GetClient(cfg conf.Config, clientID, joinToken string) (*pb.GetClientResponse, error) {
	apiEndpoint := conf.GetAPIURL(cfg)
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
