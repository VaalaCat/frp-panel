package dao

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"

	"github.com/VaalaCat/frp-panel/models"
	"github.com/VaalaCat/frp-panel/utils"
	"github.com/sirupsen/logrus"
)

func (q *queryImpl) InitCert(template *x509.Certificate) *tls.Config {
	var (
		certPem []byte
		keyPem  []byte
	)
	cnt, err := q.CountCerts()
	if err != nil {
		logrus.Fatal(err)
	}
	if cnt == 0 {
		certPem, keyPem, err = GenX509Info(template)
		if err != nil {
			logrus.Fatal(err)
		}
		if err = q.ctx.GetApp().GetDBManager().GetDefaultDB().Create(&models.Cert{
			Name:     "default",
			CertFile: certPem,
			CaFile:   certPem,
			KeyFile:  keyPem,
		}).Error; err != nil {
			logrus.Fatal(err)
		}
	} else {
		keyPem, certPem, err = q.GetDefaultKeyPair()
		if err != nil {
			logrus.Fatal(err)
		}
	}

	resp, err := utils.TLSServerCert(certPem, keyPem)
	if err != nil {
		logrus.Fatal(err)
	}
	return resp
}

func GenX509Info(template *x509.Certificate) (certPem []byte, keyPem []byte, err error) {

	// priv, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	// if err != nil {
	// 	return nil, nil, err
	// }

	priv, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, err
	}

	cert, err := x509.CreateCertificate(rand.Reader, template, template,
		priv.Public(), priv)
	if err != nil {
		return nil, nil, err
	}

	var certBuf bytes.Buffer
	pem.Encode(&certBuf, &pem.Block{
		Type: "CERTIFICATE", Bytes: cert,
	})

	var keyBuf bytes.Buffer
	pem.Encode(&keyBuf, utils.PemBlockForPrivKey(priv))
	return certBuf.Bytes(), keyBuf.Bytes(), nil
}

func (q *queryImpl) CountCerts() (int64, error) {
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	var count int64
	err := db.Model(&models.Cert{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (q *queryImpl) GetDefaultKeyPair() (keyPem []byte, certPem []byte, err error) {
	resp := &models.Cert{}
	err = q.ctx.GetApp().GetDBManager().GetDefaultDB().Model(&models.Cert{}).
		Where(&models.Cert{Name: "default"}).First(resp).Error
	if err != nil {
		return nil, nil, err
	}
	return resp.KeyFile, resp.CertFile, nil
}
