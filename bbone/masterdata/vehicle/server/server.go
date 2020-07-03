package server

import (
	"bluebirdgroup/bbone/commons/cert"
	"bluebirdgroup/bbone/commons/config"
	"bluebirdgroup/bbone/commons/constant"
	g "bluebirdgroup/bbone/commons/grpc"
	"context"

	"bluebirdgroup/bbone/masterdata/vehicle/database"
	"bluebirdgroup/bbone/masterdata/vehicle/proto"
	"bluebirdgroup/bbone/masterdata/vehicle/repository"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	repo repository.DBReaderWriter
}

func init() {
	fmt.Println("[INFO] Running service vehicle")
}

//RunGrpcVehicle .....
func RunGrpcVehicle() {
	var grpcServer *grpc.Server
	cacert := config.Get(constant.CACertKey, "")
	serverCert := config.Get(constant.ServerCertKey, "")
	keyCert := config.Get(constant.KeyCertKey, "")
	if cacert != "" && serverCert != "" && keyCert != "" {
		_, creds, err := cert.CreateTransportCredentials(cacert, serverCert, keyCert)
		if err != nil {
			log.Fatal(err)
		}
		grpcServer = grpc.NewServer(append(g.Recovery(), grpc.Creds(*creds))...)
	} else {
		grpcServer = grpc.NewServer(g.Recovery()...)
	}

	repo, err := database.MGConnection(repository.DBConfiguration{
		DBHost:     "",
		DBName:     "",
		DBOptions:  "",
		DBPassword: "password",
		DBPort:     "password",
		DBUser:     "user",
	})

	if err != nil {
		log.Fatal(err)
		return
	}
	defer repo.Close()

	proto.RegisterVehicleServer(grpcServer, &server{repo: repo})

	reflection.Register(grpcServer)

	port := "9089"
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal("Failed created server")
		return
	}
	fmt.Println("Server run on " + port)
	grpcServer.Serve(lis)
	defer grpcServer.GracefulStop()
}

func (s *server) Add(_ context.Context, req *proto.VehicleModel) (*proto.ResEmpty, error) {
	_, err := s.repo.Add(req)
	if err != nil {
		return nil, err
	}
	return &proto.ResEmpty{}, nil
}

func (s *server) Edit(_ context.Context, req *proto.VehicleModel) (*proto.ResEmpty, error) {
	_, err := s.repo.Edit(req)
	if err != nil {
		return nil, err
	}
	return &proto.ResEmpty{}, nil
}

func (s *server) Delete(_ context.Context, req *proto.ReqByNoLambung) (*proto.ResEmpty, error) {
	_, err := s.repo.Delete(req)
	if err != nil {
		return nil, err
	}
	return &proto.ResEmpty{}, nil
}

func (s *server) GetByNoLambung(_ context.Context, req *proto.ReqByNoLambung) (*proto.VehicleModel, error) {
	resp, err := s.repo.GetByNoLambung(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *server) GetByPool(_ context.Context, req *proto.ReqByPool) (*proto.ResVehicleList, error) {
	resp, err := s.repo.GetByPool(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *server) GetAllByPrsh(_ context.Context, req *proto.ReqAllByPrsh) (*proto.ResVehicleList, error) {
	resp, err := s.repo.GetAllByPrsh(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
