package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func AdminAuth(c *fiber.Ctx) error {
    isAdmin, ok := c.Context().UserValue("admin").(bool)

    if !ok {
        return fmt.Errorf("unauthorized")
    }
    if !isAdmin {
        return fmt.Errorf("unauthorized")
    }

    return c.Next()
}
