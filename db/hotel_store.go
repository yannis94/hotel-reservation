package db

import (
	"context"

	customtypes "github.com/yannis94/hotel-reservation/custom_types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelStore interface {
    Dropper 

    GetHotelByID(context.Context, bson.M) (*customtypes.Hotel, error)
    GetHotels(context.Context, bson.M) ([]*customtypes.Hotel, error)
    InsertHotel(context.Context, *customtypes.Hotel) (*customtypes.Hotel, error)
    UpdateHotel(context.Context, bson.M, bson.M) error
    DeleteHotel(context.Context, string) error
}


type MongoHotelStore struct {
    client *mongo.Client
    coll *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client) *MongoHotelStore {
    return &MongoHotelStore{
        client: client,
        coll: client.Database(DbName).Collection(hotelColl),
    }
}

func (db *MongoHotelStore) InsertHotel(ctx context.Context, hotel *customtypes.Hotel) (*customtypes.Hotel, error) {
    res, err := db.coll.InsertOne(ctx, hotel)

    if err != nil {
        return nil, err
    }

    hotel.ID = res.InsertedID.(primitive.ObjectID)
    
    return hotel, nil
}

func (db *MongoHotelStore) GetHotelByID(ctx context.Context, filter bson.M) (*customtypes.Hotel, error) {
    var hotel customtypes.Hotel

    if err := db.coll.FindOne(ctx, filter).Decode(&hotel); err != nil {
        return &hotel, err
    }

    return &hotel, nil
}

func (db *MongoHotelStore) GetHotels(ctx context.Context, filter bson.M) ([]*customtypes.Hotel, error) {
    var hotels []*customtypes.Hotel
    resp, err := db.coll.Find(ctx, filter)

    if err != nil {
        return hotels, err
    }

    if err := resp.All(ctx, &hotels); err != nil {
        return hotels, err
    }

    return hotels, nil
}

func (db *MongoHotelStore) UpdateHotel(ctx context.Context, filter, update bson.M) error {
    _, err := db.coll.UpdateOne(ctx, filter, update)
    
    if err != nil {
        return err
    }

    return nil
}

func (db *MongoHotelStore) DeleteHotel(ctx context.Context, id string) error {
    return nil 
}

func (db *MongoHotelStore) Drop(ctx context.Context) error {
    return db.coll.Drop(ctx)
}
