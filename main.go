package main

import (
	"context"
	"flag"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/yannis94/hotel-reservation/api"
	"github.com/yannis94/hotel-reservation/api/middleware"
	"github.com/yannis94/hotel-reservation/db"
)

var config = fiber.Config{
    // Override default error handler
    ErrorHandler: func(ctx *fiber.Ctx, err error) error {
        return ctx.JSON(map[string]string{"error": err.Error()})
    },
}

func main() {
    listenAddr := flag.String("listenAddr", ":3000", "The server's listen address (default :3000).")
    flag.Parse()
    
    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DbUri))
    if err != nil {
        panic(err)
    }

    var (
        hotelStore = db.NewMongoHotelStore(client)
        roomStore = db.NewMongoRoomStore(client, hotelStore)
        userStore = db.NewMongoUserStore(client)
        store = &db.Store{
            User: userStore,
            Hotel: hotelStore,
            Room: roomStore,
        }
        hotelHandler = api.NewHotelHandler(store)
        userHandler = api.NewUserHandler(userStore)
        authHandler = api.NewAuthHandler(userStore)
        app = fiber.New(config)
        api = app.Group("/api")
        apiv1 = app.Group("/api/v1", middleware.JWTAuthentication)
    )


    api.Post("/auth", authHandler.HandleAuthenticate)

    apiv1.Get("/user", userHandler.HandleGetUsers)
    apiv1.Get("/user/:id", userHandler.HandleGetUserById)
    apiv1.Post("/user", userHandler.HandlePostUser)
    apiv1.Put("/user/:id", userHandler.HandlePutUser)
    apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)

    apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
    apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotelByID)
    apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)
    
    app.Listen(*listenAddr)
}
