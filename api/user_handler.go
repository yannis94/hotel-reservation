package api

import (
	"context"

	"github.com/gofiber/fiber/v2"
	customtypes "github.com/yannis94/hotel-reservation/custom_types"
	"github.com/yannis94/hotel-reservation/db"
)

type UserHandler struct {
    userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
    return &UserHandler{
        userStore: userStore,
    }
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
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

func (h *UserHandler) HandleGetUserById(c *fiber.Ctx) error {
    var (
        id = c.Params("id")
        ctx = context.Background()
    )

    user, err := h.userStore.GetUserByID(ctx, id)

    if err != nil {
        return err
    }

    return c.JSON(user)
}
