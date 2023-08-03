package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	customtypes "github.com/yannis94/hotel-reservation/custom_types"
)

func createTestUser(app *fiber.App) *customtypes.User {
    params := customtypes.CreateUserParams{
        FirstName: "Elliot",
        LastName: "Alderson",
        Email: "elliot.alderson@fsociety.com",
        Password: "thesuperpassword",
    }

    data, _ := json.Marshal(params)

    req := httptest.NewRequest("POST", "/api/v1/user", bytes.NewReader(data))
    req.Header.Add("Content-Type", "application/json")
    res, _ := app.Test(req)

    defer res.Body.Close()

    var user customtypes.User

    json.NewDecoder(res.Body).Decode(&user)

    return &user
}

func TestAuthenticate(t *testing.T) {
    tdb := setup(t)
    defer tdb.teardown(t)

    userHandler := NewUserHandler(tdb.UserStore)
    authHandler := NewAuthHandler(tdb.UserStore)
    app := fiber.New()
    app.Post("/api/v1/user", userHandler.HandlePostUser)
    app.Post("/api/auth", authHandler.HandleAuthenticate)

    userTest := createTestUser(app)
    fmt.Println(userTest)


    params := AuthParams{
        Email: "elliot.alderson@fsociety.com",
        Password: "thesuperpassword",
    }

    data, _ := json.Marshal(params)

    req := httptest.NewRequest("POST", "/api/auth", bytes.NewReader(data))
    req.Header.Add("Content-Type", "application/json")
    res, err := app.Test(req)

    if err != nil {
        t.Fatal(err)
    }
    if res.StatusCode != http.StatusOK {
        t.Errorf("Expect http status 200 but got %d", res.StatusCode)
    }

    defer res.Body.Close()

    var authResp AuthResponse

    if err := json.NewDecoder(res.Body).Decode(&authResp); err != nil {
        t.Error(err)
    }
    if authResp.User.FirstName != userTest.FirstName {
        t.Errorf("The first name return is %s but should be %s", authResp.User.FirstName, userTest.FirstName)
    }
    if authResp.User.LastName != userTest.LastName {
        t.Errorf("The first name return is %s but should be %s", authResp.User.LastName, userTest.LastName)
    }
    if authResp.User.Email != userTest.Email {
        t.Errorf("The first name return is %s but should be %s", authResp.User.Email, userTest.Email)
    }
    if authResp.Token == "" {
        t.Error("The authentication should provide a JWT token")
    }
    
    fmt.Println(authResp)
}
