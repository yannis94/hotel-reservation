package api

import (
	"github.com/gofiber/fiber/v2"
	customtypes "github.com/yannis94/hotel-reservation/custom_types"
)

func HandleGetUsers(c *fiber.Ctx) error {
    user1 := &customtypes.User{
        FirstName: "Elliot",
        LastName: "Alderson",
    }
    user2 := &customtypes.User{
        FirstName: "Elliot",
        LastName: "Alderson",
    }
    res := []customtypes.User{ *user1, *user2}
    return c.JSON(res)
}

func HandleGetUserById(c *fiber.Ctx) error {
    user := &customtypes.User{
        FirstName: "Elliot",
        LastName: "Alderson",
    }
    return c.JSON(user)
}
