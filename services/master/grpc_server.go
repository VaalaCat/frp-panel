package master

import (
	"context"
	"fmt"
	"io"
	"net"

	"github.com/VaalaCat/frp-panel/biz/master/client"
	masterserver "github.com/VaalaCat/frp-panel/biz/master/server"
	"github.com/VaalaCat/frp-panel/biz/master/shell"
	"github.com/VaalaCat/frp-panel/biz/master/streamlog"
	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/conf"
	"github.com/VaalaCat/frp-panel/dao"
	"github.com/VaalaCat/frp-panel/logger"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/rpc"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type server struct {
	pb.UnimplementedMasterServer
}

func NewRpcServer(creds credentials.TransportCredentials) *grpc.Server {
	s := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterMasterServer(s, &server{})
	return s
}

func RunRpcServer(s *grpc.Server) {
	lis, err := net.Listen("tcp", conf.RPCListenAddr())
	if err != nil {
		logrus.Fatalf("rpc server failed to listen: %v", err)
	}

	logrus.Infof("start server")
	if err := s.Serve(lis); err != nil {
		logrus.Fatalf("failed to serve: %v", err)
	}
}

// PullClientConfig implements pb.MasterServer.
func (s *server) PullClientConfig(ctx context.Context, req *pb.PullClientConfigReq) (*pb.PullClientConfigResp, error) {
	logger.Logger(ctx).Infof("pull client config, clientID: [%+v]", req.GetBase().GetClientId())
	return client.RPCPullConfig(ctx, req)
}

// PullServerConfig implements pb.MasterServer.
func (s *server) PullServerConfig(ctx context.Context, req *pb.PullServerConfigReq) (*pb.PullServerConfigResp, error) {
	logger.Logger(ctx).Infof("pull server config, serverID: [%+v]", req.GetBase().GetServerId())
	return masterserver.RPCPullConfig(ctx, req)
}

// FRPCAuth implements pb.MasterServer.
func (*server) FRPCAuth(ctx context.Context, req *pb.FRPAuthRequest) (*pb.FRPAuthResponse, error) {
	logger.Logger(ctx).Infof("frpc auth, user: [%+v] ,serverID: [%+v]", req.GetUser(), req.GetBase().GetServerId())
	return masterserver.FRPAuth(ctx, req)
}

// ServerSend implements pb.MasterServer.
func (s *server) ServerSend(sender pb.Master_ServerSendServer) error {
	ctx := context.Background()
	logger.Logger(ctx).Infof("server get a client connected")
	var done chan bool
	for {
		req, err := sender.Recv()
		if err == io.EOF {
			logger.Logger(ctx).Infof("finish server send, client id: [%s]", req.GetClientId())
			return nil
		}

		if err != nil {
			logger.Logger(context.Background()).WithError(err).Errorf("cannot recv from client, id: [%s]", req.GetClientId())
			return err
		}

		cliType := ""

		if req.GetEvent() == pb.Event_EVENT_REGISTER_CLIENT || req.GetEvent() == pb.Event_EVENT_REGISTER_SERVER {
			if len(req.GetSecret()) == 0 {
				logger.Logger(ctx).Errorf("rpc auth token is empty")
				sender.Send(&pb.ServerMessage{
					Event: req.GetEvent(),
					Data:  []byte("rpc auth token is invalid"),
				})
				return fmt.Errorf("rpc auth token is invalid")
			}
			var secret string
			switch req.GetEvent() {
			case pb.Event_EVENT_REGISTER_CLIENT:
				cli, err := dao.AdminGetClientByClientID(req.GetClientId())
				if err != nil {
					logger.Logger(context.Background()).WithError(err).Errorf("cannot get client, %s id: [%s]", req.GetEvent().String(), req.GetClientId())
					sender.Send(&pb.ServerMessage{
						Event: req.GetEvent(),
						Data:  []byte("rpc auth token is invalid"),
					})
					return err
				}
				secret = cli.ConnectSecret
				cliType = common.CliTypeClient
			case pb.Event_EVENT_REGISTER_SERVER:
				srv, err := dao.AdminGetServerByServerID(req.GetClientId())
				if err != nil {
					logger.Logger(context.Background()).WithError(err).Errorf("cannot get server, %s id: [%s]", req.GetEvent().String(), req.GetClientId())
					sender.Send(&pb.ServerMessage{
						Event: req.GetEvent(),
						Data:  []byte("rpc auth token is invalid"),
					})
					return err
				}
				secret = srv.ConnectSecret
				cliType = common.CliTypeServer
			}

			if secret != req.GetSecret() {
				logger.Logger(ctx).Errorf("invalid secret, %s id: [%s]", req.GetEvent().String(), req.GetClientId())
				sender.Send(&pb.ServerMessage{
					Event: req.GetEvent(),
					Data:  []byte("rpc auth token is invalid"),
				})
				return fmt.Errorf("invalid secret, %s id: [%s]", req.GetEvent().String(), req.GetClientId())
			}

			rpc.GetClientsManager().Set(req.GetClientId(), cliType, sender)
			done = rpc.Recv(req.GetClientId())
			sender.Send(&pb.ServerMessage{
				Event:     req.GetEvent(),
				ClientId:  req.GetClientId(),
				SessionId: req.GetClientId(),
			})
			logger.Logger(ctx).Infof("register success, req: [%+v]", req)
			break
		}
	}
	<-done
	return nil
}

// PushProxyInfo implements pb.MasterServer.
func (s *server) PushProxyInfo(ctx context.Context, req *pb.PushProxyInfoReq) (*pb.PushProxyInfoResp, error) {
	logger.Logger(ctx).Infof("push proxy info, req server: [%+v]", req.GetProxyInfos())
	return masterserver.PushProxyInfo(ctx, req)
}

func (s *server) PushClientStreamLog(sender pb.Master_PushClientStreamLogServer) error {
	return streamlog.PushClientStreamLog(sender)
}

func (s *server) PushServerStreamLog(sender pb.Master_PushServerStreamLogServer) error {
	return streamlog.PushServerStreamLog(sender)
}

func (s *server) PTYConnect(sender pb.Master_PTYConnectServer) error {
	return shell.PTYConnect(sender)
}
