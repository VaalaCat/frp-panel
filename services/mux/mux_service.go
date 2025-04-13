package mux

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type MuxServer interface {
	Run()
	Stop()
}

type muxImpl struct {
	srv *http.Server
	lis net.Listener
	tls bool
}

func NewMux(grpcServer, apiServer http.Handler, lis net.Listener, creds *tls.Config) MuxServer {
	tlsServer := grpcHandlerFunc(grpcServer, apiServer)
	tlsServer.TLSConfig = creds
	return &muxImpl{
		srv: tlsServer,
		lis: lis,
		tls: creds != nil,
	}
}

func (m *muxImpl) Run() {
	if m.tls {
		if err := m.srv.ServeTLS(m.lis, "", ""); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := m.srv.Serve(m.lis); err != nil {
			log.Fatal(err)
		}
	}
}

func (m *muxImpl) Stop() {
}

func grpcHandlerFunc(grpcServer http.Handler, httpHandler http.Handler) *http.Server {
	return &http.Server{Handler: h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// fmt.Printf("proto major: %d,  %s , %s\n", r.ProtoMajor, r.RequestURI, r.Header.Get("Content-Type"))
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			httpHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{}), ReadHeaderTimeout: time.Second * 30}
}
