package db

import (
	"context"

	customtypes "github.com/yannis94/hotel-reservation/custom_types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookingStore interface {
    InsertBooking(context.Context, *customtypes.Booking) (*customtypes.Booking, error)
    GetBookings(context.Context, bson.M) ([]*customtypes.Booking, error)
    GetBookingByID(context.Context, string) (*customtypes.Booking, error)
}

type MongoBookingStore struct {
    client *mongo.Client
    coll *mongo.Collection
}

func NewMongoBookingStore(client *mongo.Client) *MongoBookingStore {
    return &MongoBookingStore{
        client: client,
        coll: client.Database(DbName).Collection(bookingColl),
    }
}

func (db *MongoBookingStore) InsertBooking(ctx context.Context, booking *customtypes.Booking) (*customtypes.Booking, error) {
    res, err := db.coll.InsertOne(ctx, booking)

    if err != nil {
        return nil, err
    }

    booking.ID = res.InsertedID.(primitive.ObjectID)
    
    return booking, nil
}

func (db *MongoBookingStore) GetBookings(ctx context.Context, filter bson.M) ([]*customtypes.Booking, error) {
    cur, err := db.coll.Find(ctx, filter)

    if err != nil {
        return nil, err
    }

    var bookings []*customtypes.Booking

    if err := cur.All(ctx, &bookings); err != nil {
        return nil, err
    }

    return bookings, nil
}
 func (db *MongoBookingStore) GetBookingByID(ctx context.Context, id string) (*customtypes.Booking, error) {
     objectID, err := primitive.ObjectIDFromHex(id)

     if err != nil {
         return nil, err
     }

     var booking customtypes.Booking

     if err := db.coll.FindOne(ctx, bson.M{"_id": objectID}).Decode(&booking); err != nil {
        return nil, err
     }

     return &booking, nil
 }
