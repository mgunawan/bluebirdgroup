package cert

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"google.golang.org/grpc/credentials"
)

//CreateTransportCredentials create transport credentials
func CreateTransportCredentials(ca, cert, prikey string) (*tls.Certificate, *credentials.TransportCredentials, error) {
	var tlsCert tls.Certificate
	var transportCredentials credentials.TransportCredentials

	if len(ca) > 0 && len(cert) > 0 && len(prikey) > 0 {
		var err error
		tlsCert, err = tls.LoadX509KeyPair(cert, prikey)
		if err != nil {
			return nil, nil, err
		}
		transportCredentials, err = tlsCredentialFromKeyPair(ca, tlsCert, true)
		if err != nil {
			return nil, nil, err
		}
	}
	return &tlsCert, &transportCredentials, nil
}

func tlsCredentialFromKeyPair(cacert string, cert tls.Certificate, mutual bool) (credentials.TransportCredentials, error) {

	rawCaCert, err := ioutil.ReadFile(cacert)
	if err != nil {
		return nil, err
	}

	return tlsCredential(rawCaCert, cert, mutual), nil
}

func tlsCredential(cacert []byte, cert tls.Certificate, mutual bool) credentials.TransportCredentials {
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(cacert)

	tlsCfg := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    caCertPool,
		RootCAs:      caCertPool,
	}

	if mutual {
		tlsCfg.ClientAuth = tls.RequireAndVerifyClientCert
	}

	return credentials.NewTLS(tlsCfg)
}
