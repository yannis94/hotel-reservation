package api

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	customtypes "github.com/yannis94/hotel-reservation/custom_types"
	"github.com/yannis94/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
    userStore db.UserStore
}

type AuthParams struct {
    Email string `json:"email"`
    Password string `json:"password"`
}

type AuthResponse struct {
    User *customtypes.User
    Token string
}

func NewAuthHandler(us db.UserStore) *AuthHandler {
    return &AuthHandler{
        userStore: us,
    }
}

func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
    var authParams AuthParams

    if err := c.BodyParser(&authParams); err != nil {
        return c.JSON(err)
    }

    fmt.Println(authParams)
    user, err := h.userStore.GetUserByEmail(c.Context(), authParams.Email)

    if err != nil {
        if errors.Is(err, mongo.ErrNoDocuments) {
            return fmt.Errorf("invalid authentication")
        }
        return c.JSON(err)
    }

    if !customtypes.IsPasswordValid(user.Password, authParams.Password) {
        return fmt.Errorf("invalid authentication")
    }

    token := createTokenFromUser(user)
    
    if token == "" {
        return fmt.Errorf("internal server error")
    }

    resp := AuthResponse{
        User: user,
        Token: token,
    }

    return c.JSON(resp)
}

func createTokenFromUser(user *customtypes.User) string {
    now := time.Now()
    expireAt := now.Add(time.Hour * 4)
    claims := jwt.MapClaims{
        "id": user.ID,
        "email": user.Email,
        "iat": now.Unix(),
        "exp": expireAt.Unix(),
    }
    tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    secret := os.Getenv("JWT_SECRET")
    tknString, err := tkn.SignedString([]byte(secret))

    if err != nil {
        fmt.Println("Failed to sign token")
    }

    return tknString
}
