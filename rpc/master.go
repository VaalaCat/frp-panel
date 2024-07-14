package rpc

import (
	"context"

	"github.com/VaalaCat/frp-panel/conf"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/imroc/req/v3"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

var (
	masterCli pb.MasterClient
)

func newMasterCli() {
	conn, err := grpc.NewClient(conf.RPCCallAddr(),
		grpc.WithTransportCredentials(conf.ClientCred))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}

	masterCli = pb.NewMasterClient(conn)
}

func MasterCli(c context.Context) (pb.MasterClient, error) {
	if masterCli == nil {
		newMasterCli()
	}
	return masterCli, nil
}

func GetClientCert(clientID, clientSecret string, clientType pb.ClientType) []byte {
	apiEndpoint := conf.GetAPIURL()
	c := req.C()
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
