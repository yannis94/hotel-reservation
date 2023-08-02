package middleware

import (
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

    fmt.Println(token)
    if err := parseJWTToken(token); err != nil {
        return err
    }

    return nil
}

func parseJWTToken(tokenStr string) error {
    token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            fmt.Println("Invalid signion method", token.Header["alg"])
            return nil, fmt.Errorf("unauthorized")
        }

        secret := os.Getenv("JWT_SECRET")
        return []byte(secret), nil
    })

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        fmt.Println(claims)
    }
    if err != nil {
        fmt.Println("fail parse token", err)
        return fmt.Errorf("unauthorized")
    }
    return fmt.Errorf("unauthorized")
}
