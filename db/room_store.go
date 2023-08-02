package db

import (
	"context"
	"log"

	customtypes "github.com/yannis94/hotel-reservation/custom_types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomStore interface {
    Dropper 

    GetRoomByID(context.Context, string) (*customtypes.Room, error)
    GetRooms(context.Context, bson.M) ([]*customtypes.Room, error)
    InsertRoom(context.Context, *customtypes.Room) (*customtypes.Room, error)
    UpdateRoom(context.Context, bson.M, customtypes.UpdateRoomParams) error
    DeleteRoom(context.Context, string) error
}


type MongoRoomStore struct {
    client *mongo.Client
    coll *mongo.Collection

    HotelStore
}

func NewMongoRoomStore(client *mongo.Client, hotelStore HotelStore) *MongoRoomStore {
    return &MongoRoomStore{
        client: client,
        coll: client.Database(DbName).Collection(roomColl),
        HotelStore: hotelStore,
    }
}

func (db *MongoRoomStore) InsertRoom(ctx context.Context, room *customtypes.Room) (*customtypes.Room, error) {
    res, err := db.coll.InsertOne(ctx, room)

    if err != nil {
        return nil, err
    }

    room.ID = res.InsertedID.(primitive.ObjectID)

    filter := bson.M{ "_id": room.HotelID }
    update := bson.M{ "$push": bson.M{ "rooms": room.ID } }

    if err := db.UpdateHotel(ctx, filter, update); err != nil {
        log.Println(err)
        return nil, err
    }
    
    return room, nil
}

func (db *MongoRoomStore) GetRoomByID(ctx context.Context, id string) (*customtypes.Room, error) {
    return nil, nil
}

func (db *MongoRoomStore) GetRooms(ctx context.Context, filter bson.M) ([]*customtypes.Room, error) {
    var rooms []*customtypes.Room
    res, err := db.coll.Find(ctx, filter)

    if err != nil {
        return rooms, err
    }

    if  err := res.All(ctx, &rooms); err != nil {
        return rooms, err
    }

    return rooms, nil
}

func (db *MongoRoomStore) UpdateRoom(ctx context.Context, filter bson.M, values customtypes.UpdateRoomParams) error {
    return nil
}

func (db *MongoRoomStore) DeleteRoom(ctx context.Context, id string) error {
    return nil 
}

func (db *MongoRoomStore) Drop(ctx context.Context) error {
    return db.coll.Drop(ctx)
}
