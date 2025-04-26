package utils

import (
	"context"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"

	"github.com/VaalaCat/frp-panel/utils/logger"
	"google.golang.org/grpc/credentials"
)

func PublicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}

func PemBlockForPrivKey(priv interface{}) *pem.Block {
	ctx := context.Background()
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			logger.Logger(ctx).Fatalf("Unable to marshal ECDSA private key: %v", err)
		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
	default:
		return nil
	}
}

func TLSServerCert(certPem, keyPem []byte) (*tls.Config, error) {
	cert, err := tls.X509KeyPair(certPem, keyPem)
	if err != nil {
		return nil, err
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.NoClientCert,
	}
	return config, nil
}

func TLSClientCert(caPem []byte) (credentials.TransportCredentials, error) {
	certpool := x509.NewCertPool()
	certpool.AppendCertsFromPEM(caPem)
	return credentials.NewClientTLSFromCert(certpool, ""), nil
}

func TLSClientCertNoValidate(caPem []byte) (credentials.TransportCredentials, error) {
	certpool := x509.NewCertPool()
	certpool.AppendCertsFromPEM(caPem)

	config := &tls.Config{
		RootCAs:            certpool,
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS12,
	}

	return credentials.NewTLS(config), nil
}
