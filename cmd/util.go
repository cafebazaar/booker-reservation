package cmd

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/cafebazaar/booker-reservation/common"
)

var (
	_keyPair  *tls.Certificate
	_certPool *x509.CertPool
)

func keyPair() (*tls.Certificate, error) {
	if _keyPair == nil {
		pair, err := tls.LoadX509KeyPair(common.ConfigString("CERT_FILE"), common.ConfigString("KEY_FILE"))
		if err != nil {
			return nil, fmt.Errorf("Failed to load tls key pair: %s", err)
		}

		_keyPair = &pair
	}

	return _keyPair, nil
}

func certPool() (*x509.CertPool, error) {
	if _certPool == nil {
		newCertPool := x509.NewCertPool()

		caContent, err := ioutil.ReadFile(common.ConfigString("CA_FILE"))
		if err != nil {
			return nil, fmt.Errorf("Failed to load tls ca file: %s", err)
		}

		ok := newCertPool.AppendCertsFromPEM(caContent)
		if !ok {
			return nil, errors.New("Failed to append tls ca certs")
		}

		_certPool = newCertPool
	}
	return _certPool, nil
}
