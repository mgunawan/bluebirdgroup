package grpc

import (
	"bluebirdgroup/bbone/commons/cert"
	"bluebirdgroup/bbone/commons/config"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	bluebirdIDDomain = "bbone"
)

//NewServerConnection ...
func NewServerConnection(addr, ca, cer, prikey string) (*grpc.ClientConn, error) {
	var creds credentials.TransportCredentials
	var err error
	if strings.Contains(addr, bluebirdIDDomain) {
		creds, err = credentials.NewClientTLSFromFile(config.Get("BB1_ID_CERT", "tls.crt"), "")
		if err != nil {
			return nil, err
		}
	} else {
		_, credsVal, _ := cert.CreateTransportCredentials(ca, cer, prikey)
		creds = *credsVal
	}
	var conn *grpc.ClientConn
	if creds == nil {
		conn, err = grpc.Dial(addr, grpc.WithInsecure())
	} else {
		conn, err = grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	}
	if err != nil {
		return nil, err
	}
	return conn, nil
}
