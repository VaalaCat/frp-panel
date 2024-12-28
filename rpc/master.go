package rpc

import (
	"context"
	"errors"

	"github.com/VaalaCat/frp-panel/common"
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

func InitClient(clientID, joinToken string) (*pb.InitClientResponse, error) {
	apiEndpoint := conf.GetAPIURL()
	c := req.C()
	rawReq, err := proto.Marshal(&pb.InitClientRequest{
		ClientId: &clientID,
	})
	if err != nil {
		return nil, err
	}

	r, err := c.R().SetHeader("Content-Type", "application/x-protobuf").
		SetHeader(common.AuthorizationKey, joinToken).
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

func GetClient(clientID, joinToken string) (*pb.GetClientResponse, error) {
	apiEndpoint := conf.GetAPIURL()
	c := req.C()
	rawReq, err := proto.Marshal(&pb.GetClientRequest{
		ClientId: &clientID,
	})
	if err != nil {
		return nil, err
	}

	r, err := c.R().SetHeader("Content-Type", "application/x-protobuf").
		SetHeader(common.AuthorizationKey, joinToken).
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
