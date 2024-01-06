package rpc

import (
	"context"

	"github.com/VaalaCat/frp-panel/conf"
	"github.com/VaalaCat/frp-panel/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func MasterCli(c context.Context) (pb.MasterClient, error) {
	conn, err := grpc.Dial(conf.RPCCallAddr(),
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		))
	if err != nil {
		return nil, err
	}

	client := pb.NewMasterClient(conn)
	return client, nil
}
