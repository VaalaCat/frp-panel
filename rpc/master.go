package rpc

import (
	"context"

	"github.com/VaalaCat/frp-panel/conf"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/imroc/req/v3"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

func MasterCli(c context.Context) (pb.MasterClient, error) {
	conn, err := grpc.Dial(conf.RPCCallAddr(),
		grpc.WithTransportCredentials(conf.ClientCred))
	if err != nil {
		return nil, err
	}

	client := pb.NewMasterClient(conn)
	return client, nil
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
