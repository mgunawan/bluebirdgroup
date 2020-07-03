package database

import (
	"log"
	"time"

	"bluebirdgroup/bbone/masterdata/vehicle/proto"
	"bluebirdgroup/bbone/masterdata/vehicle/repository"

	"context"
	"fmt"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mutex = &sync.RWMutex{}
var ctx context.Context

type mgReaderWriter struct {
	db *mongo.Database
	//client *mongo.Client
	mutex sync.RWMutex
}

func init() {
	fmt.Println("[INFO] Connect to mongodb")
}

//MGConnection ...
func MGConnection(conf repository.DBConfiguration) (repository.DBReaderWriter, error) {
	//LOCAL
	//clientOptions := options.Client().ApplyURI("mongodb://mongodb:27017")
	//client, err := mongo.Connect(context.TODO(), clientOptions)

	//MONGODB CLOUD
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://mongouser:mongopassword@cluster0-qw0g6.mongodb.net/MasterData?retryWrites=true&w=majority"))
	if err != nil {
		fmt.Println("[INFO] Error connection to MongoDB")
		return nil, err
	}
	db := client.Database("MasterData")
	//defer client.Disconnect(ctx) // error client.disconnect bila di pasang

	fmt.Println("[INFO] Sucessfull connect to mongoDB")
	return &mgReaderWriter{
		db: db,
		//client: client,
		mutex: sync.RWMutex{},
	}, nil
}

func (p *mgReaderWriter) Close() error {
	return nil
}

func (p *mgReaderWriter) Add(r *proto.VehicleModel) (int64, error) {
	collection := p.db.Collection("vehicle")

	insertResult, err := collection.InsertOne(context.TODO(), &r)
	if err != nil {
		return 0, err
	}
	fmt.Println("Inserted a Single Document: ", insertResult.InsertedID)

	return 999, nil
}

func (p *mgReaderWriter) Edit(request *proto.VehicleModel) (int64, error) {
	collection := p.db.Collection("vehicle")
	filter := bson.M{"nolambung": request.Nolambung}
	update := bson.M{
		"$set": bson.M{"nopolisi": request.Nopolisi, "jenisservice": request.Jenisservice,
			"kodeperusahaan": request.Kodeperusahaan, "kodebbm1": request.Kodebmb1,
			"kodebbm2": request.Kodebmm2, "issticker": request.Issticker},
	}

	upsert := true
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	result := collection.FindOneAndUpdate(ctx, filter, update, &opt)
	if result.Err() != nil {
		return 0, result.Err()
	}
	return 1, nil
}

func (p *mgReaderWriter) Delete(request *proto.ReqByNoLambung) (int64, error) {
	result, err := p.db.Collection("vehicle").DeleteOne(ctx, bson.M{"nolambung": request.Nolambung})
	if err != nil {
		log.Fatal(err)
	}
	return result.DeletedCount, nil
}

func (p *mgReaderWriter) GetByNoLambung(request *proto.ReqByNoLambung) (*proto.VehicleModel, error) {
	filter := bson.M{"nolambung": request.Nolambung}
	cursor, err := p.db.Collection("vehicle").Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	var vehicle proto.VehicleModel
	for cursor.Next(context.TODO()) {
		err := cursor.Decode(&vehicle)
		if err != nil {
			log.Fatal(err)
		}
	}
	return &vehicle, nil
}

func (p *mgReaderWriter) GetByPool(request *proto.ReqByPool) (*proto.ResVehicleList, error) {
	var vehicles = []*proto.VehicleModel{}
	filter := bson.M{"kodepool": request.Kodepool}
	cursor, err := p.db.Collection("vehicle").Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	for cursor.Next(context.TODO()) {
		var veh proto.VehicleModel
		err := cursor.Decode(&veh)
		if err != nil {
			log.Fatal(err)
		}

		vehicles = append(vehicles, &veh)
	}
	return &proto.ResVehicleList{VehicleList: vehicles}, nil
}

func (p *mgReaderWriter) GetAllByPrsh(request *proto.ReqAllByPrsh) (*proto.ResVehicleList, error) {
	var vehicles = []*proto.VehicleModel{}
	filter := bson.M{"kodeperusahaan": request.Kodeperusahaan}
	cursor, err := p.db.Collection("vehicle").Find(context.TODO(), filter)

	// findOptions := options.Find()
	// findOptions.SetLimit(50)
	// cursor, err := p.db.Collection("vehicle").Find(context.TODO(), bson.D{{}}, findOptions)

	if err != nil {
		log.Fatal(err)
	}

	for cursor.Next(context.TODO()) {
		var veh proto.VehicleModel
		err := cursor.Decode(&veh)
		if err != nil {
			log.Fatal(err)
		}

		vehicles = append(vehicles, &veh)
	}
	return &proto.ResVehicleList{VehicleList: vehicles}, nil
}

// func (p *mgReaderWriter) GetById(request *proto.PoolById) (*proto.Pool, error) {
// 	filter := bson.M{"kode": request.Kodepool}
// 	cursor, err := p.db.Collection("pool").Find(context.TODO(), filter)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	var pool proto.Pool
// 	for cursor.Next(context.TODO()) {
// 		err := cursor.Decode(&pool)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}
// 	return &pool, nil
// }

// func (p *mgReaderWriter) GetAllByRegion(request *proto.PoolRegion) (*proto.PoolList, error) {
// 	var pools = []*proto.Pool{}
// 	filter := bson.M{"region": request.Region}
// 	cursor, err := p.db.Collection("pool").Find(context.TODO(), filter)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	for cursor.Next(context.TODO()) {
// 		var pool proto.Pool
// 		err := cursor.Decode(&pool)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		pools = append(pools, &pool)
// 	}
// 	return &proto.PoolList{List: pools}, nil
// }

// func (p *mgReaderWriter) GetPoolInduk(request *proto.PoolSatelit) (*proto.Pool, error) {
// 	// collection := p.client.Database("masterdata").Collection("pool")
// 	// fmt.Println(collection)

// 	// for pool := range collection {
// 	// 	fmt.Println(pool)
// 	// }

// 	return &proto.Pool{}, nil
// }

//function
// func toDoc(v interface{}) (doc *bson.Document, err error) {
// 	data, err := bson.Marshal(v)
// 	if err != nil {
// 		return
// 	}

// 	err = bson.Unmarshal(data, &doc)
// 	return
// }
