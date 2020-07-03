package service

import (
	"context"

	pb "bluebirdgroup/bbone/masterdata/vehicle/proto"
)

//VehicleService ............
type VehicleService interface {
	Add(context.Context, *pb.VehicleModel) error
	Edit(context.Context, *pb.VehicleModel) error
	Delete(context.Context, *pb.ReqByNoLambung) error
	GetByNoLambung(context.Context, *pb.ReqByNoLambung) error
	GetByPool(context.Context, *pb.ReqByPool) error
	GetAllByPrsh(context.Context, *pb.ReqAllByPrsh) (*pb.ResVehicleList, error)
}
