package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	customtypes "github.com/yannis94/hotel-reservation/custom_types"
	"github.com/yannis94/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
    dbtestUri = "mongodb://localhost:27017"
    dbName = "hotel-reservation-test"
)

type testdb struct {
    db.UserStore
}

func (tdb *testdb) teardown(t *testing.T) {
    if err := tdb.Drop(context.TODO()); err != nil {
        t.Fatal(err)
    }
}

func setup(t *testing.T) *testdb {
    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbtestUri))
    if err != nil {
        panic(err)
    }
    return &testdb{
        UserStore: db.NewMongoUserStore(client),
    }
}

func TestPostUser(t *testing.T) {
    tdb := setup(t)
    userHandler := NewUserHandler(tdb.UserStore)
    app := fiber.New()

    defer tdb.teardown(t)

    app.Post("/", userHandler.HandlePostUser)

    params := customtypes.CreateUserParams{
        FirstName: "John",
        LastName: "Doe",
        Email: "john.doe@mail.com",
        Password: "thesuperpassword",
    }

    data, _ := json.Marshal(params)

    req := httptest.NewRequest("POST", "/", bytes.NewReader(data))
    req.Header.Add("Content-Type", "application/json")
    res, _ := app.Test(req)

    var user customtypes.User

    json.NewDecoder(res.Body).Decode(&user)

    if len(user.ID) == 0 {
        t.Errorf("Expecting a user id to be set")
    }
    if len(user.Password) > 0 {
        t.Errorf("The user's password shouldn't be sended back")
    }
    if user.FirstName != params.FirstName {
        t.Errorf("Expected user name %s but got %s", params.FirstName, user.FirstName)
    }
    if user.LastName != params.LastName {
        t.Errorf("Expected user name %s but got %s", params.LastName, user.LastName)
    }
    if user.Email != params.Email {
        t.Errorf("Expected user name %s but got %s", params.Email, user.Email)
    }
}
