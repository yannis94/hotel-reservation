package middleware

import (
	"errors"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(c *fiber.Ctx) error {
    fmt.Println("JWT token")

    token, ok := c.GetReqHeaders()["X-Api-Access-Token"]

    if !ok {
        return c.JSON(map[string]string{"message": "unauthorized"})
    }

    claims, err := validateToken(token) 

    if err != nil {
        return err
    }

    fmt.Println("token claims:", claims)

    return c.Next()
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
    token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            fmt.Println("Invalid signion method", token.Header["alg"])
            return nil, fmt.Errorf("unauthorized")
        }

        secret := os.Getenv("JWT_SECRET")
        return []byte(secret), nil
    })

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        return claims, nil
    }
    if err != nil {
        if errors.Is(err, jwt.ErrTokenExpired) {
            return nil, fmt.Errorf("token expired")
        }
        fmt.Println("fail parse token", err)
        return nil, fmt.Errorf("unauthorized")
    }
    return nil, fmt.Errorf("unauthorized")
}
