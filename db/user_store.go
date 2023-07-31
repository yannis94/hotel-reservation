package db

import (
	"context"

	customtypes "github.com/yannis94/hotel-reservation/custom_types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserStore interface {
    GetUserByID(context.Context, string) (*customtypes.User, error)
}

type MongoUserStore struct {
    client *mongo.Client
    coll *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client) *MongoUserStore {
    return &MongoUserStore{
        client: client,
        coll: client.Database(dbName).Collection(userColl),
    }
}

func (db *MongoUserStore) GetUserByID(ctx context.Context, id string) (*customtypes.User, error) {
    objectId, err := primitive.ObjectIDFromHex(id)

    if err != nil {
        return nil, err
    }

    var user customtypes.User

    if err := db.coll.FindOne(ctx, bson.M{"_id": objectId}).Decode(&user); err != nil {
        return nil, err
    }

    return &user, nil
}
