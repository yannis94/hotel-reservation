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
    GetUsers(context.Context) ([]*customtypes.User, error)
    InsertUser(context.Context, *customtypes.User) (*customtypes.User, error)
    UpdateUser(context.Context, bson.M, customtypes.UpdateUserParams) error
    DeleteUser(context.Context, string) error
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

func (db *MongoUserStore) InsertUser(ctx context.Context, user *customtypes.User) (*customtypes.User, error) {
    res, err := db.coll.InsertOne(ctx, user)

    if err != nil {
        return nil, err
    }

    user.ID = res.InsertedID.(primitive.ObjectID)
    return user, nil
}

func (db *MongoUserStore) UpdateUser(ctx context.Context, filter bson.M, params customtypes.UpdateUserParams) error {
    values := params.ToBSON()

    update := bson.D{
        {
            "$set", values,
        },
    }

    _, err := db.coll.UpdateOne(ctx, filter, update)

    if err != nil {
        return err
    }

    return nil
}

func (db *MongoUserStore) DeleteUser(ctx context.Context, id string) error {
    objectId, err := primitive.ObjectIDFromHex(id)

    if err != nil {
        return err
    }

    _, err = db.coll.DeleteOne(ctx, bson.M{ "_id": objectId })

    if err != nil {
        return err
    }

    return nil
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

func (db *MongoUserStore) GetUsers(ctx context.Context) ([]*customtypes.User, error) {
    cur, err := db.coll.Find(ctx, bson.M{})

    if err != nil {
        return nil, err
    }
    var users []*customtypes.User

    if err := cur.All(ctx, &users); err != nil {
        return []*customtypes.User{}, err
    }

    return users, nil
}
