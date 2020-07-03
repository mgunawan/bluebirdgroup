package repository

import (
	"io"

	pb "bluebirdgroup/bbone/masterdata/vehicle/proto"
)

//DBReaderWriter contract repository
type DBReaderWriter interface {
	io.Closer
	Add(request *pb.VehicleModel) (int64, error)
	Edit(request *pb.VehicleModel) (int64, error)
	Delete(request *pb.ReqByNoLambung) (int64, error)
	GetByNoLambung(request *pb.ReqByNoLambung) (*pb.VehicleModel, error)
	GetByPool(request *pb.ReqByPool) (*pb.ResVehicleList, error)
	GetAllByPrsh(*pb.ReqAllByPrsh) (*pb.ResVehicleList, error)
}

//DBConfiguration ...
type DBConfiguration struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBOptions  string
}
