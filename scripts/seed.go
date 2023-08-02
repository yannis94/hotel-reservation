package main

import (
	"context"
	"fmt"

	customtypes "github.com/yannis94/hotel-reservation/custom_types"
	"github.com/yannis94/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
    client *mongo.Client
    hotelStore db.HotelStore
    roomStore db.RoomStore
    userStore db.UserStore
    ctx = context.Background()
)

func seedUser(firstname, lastname, email, password string) {
    userParams := customtypes.CreateUserParams{
        FirstName: firstname,
        LastName: lastname,
        Email: email,
        Password: password,
    }
    user, err := customtypes.NewUserFromParams(userParams)

    if err != nil {
        panic(err)
    }

    _, err = userStore.InsertUser(ctx, user)

    if err != nil {
        panic(err)
    }
}

func seedHotel(name, location string, rate int) {
    hotel := customtypes.Hotel{
        Name: name,
        Location: location,
        Rooms: []primitive.ObjectID{},
        Rating: rate,
    }

    rooms := []customtypes.Room{
        {
            Seaside: true,
            Price: 300.00,
            Size: "kingsize",
        },
        {
            Price: 75.00,
            Size: "small",
        },
        {
            Price: 175.00,
            Size: "normal",
        },
    }
    insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)

    if err != nil {
        panic(err)
    }

    fmt.Println(insertedHotel)

    for _, room := range rooms {
        room.HotelID = insertedHotel.ID
        insertedRoom, err := roomStore.InsertRoom(ctx, &room)

        if err != nil {
            panic(err)
        }
        fmt.Println(insertedRoom)
    }
}

func main() {
    roomStore.Drop(ctx)
    hotelStore.Drop(ctx)
    userStore.Drop(ctx)
    seedHotel("Ibis", "Montpelier", 5)
    seedHotel("Belagio", "Las Vegas", 9)
    seedHotel("Hotel du Palais", "Biarritz", 10)
    seedUser("James", "Bond", "jamesbond@mi6.en", "right2kill")
}

func init() {
    var err error
    client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DbUri))
    if err != nil {
        panic(err)
    }
    userStore = db.NewMongoUserStore(client)
    hotelStore = db.NewMongoHotelStore(client)
    roomStore = db.NewMongoRoomStore(client, hotelStore)
}
