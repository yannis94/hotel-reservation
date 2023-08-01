package api

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	customtypes "github.com/yannis94/hotel-reservation/custom_types"
	"github.com/yannis94/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
    userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
    return &UserHandler{
        userStore: userStore,
    }
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
    userID := c.Params("id")

    if err := h.userStore.DeleteUser(c.Context(), userID); err != nil {
        if errors.Is(err, mongo.ErrNoDocuments) {
            return c.JSON(map[string]string{ "message": "Not found" })
        }
        return c.JSON(err)
    }

    return c.JSON(map[string]string{ "deleted": userID })
}

func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
    var (
        params = customtypes.UpdateUserParams{}
        userID = c.Params("id")
    )

    objectID, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        return c.JSON(err)
    }

    if err := c.BodyParser(&params); err != nil {
        return c.JSON(err)
    }

    fmt.Println(params)

    filter := bson.M{ "_id": objectID}

    if err := h.userStore.UpdateUser(c.Context(), filter, params); err != nil {
        return c.JSON(err)
    }

    return c.JSON(map[string]string{ "updated": string(userID) })
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
    var params customtypes.CreateUserParams

    if err := c.BodyParser(&params); err != nil {
        return err
    }

    if errors := params.Validate(); len(errors) > 0 {
        return c.JSON(errors)
    }

    user, err := customtypes.NewUserFromParams(params)

    if err != nil {
        return err
    }

    insertedUser, err := h.userStore.InsertUser(c.Context(), user)

    if err != nil {
        return err
    }

    return c.JSON(insertedUser)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
    users, err := h.userStore.GetUsers(c.Context())

    if err != nil {
        if errors.Is(err, mongo.ErrNoDocuments) {
            return c.JSON(map[string]string{ "message": "Not found" })
        }
        return c.JSON(err)
    }

    return c.JSON(users)
}

func (h *UserHandler) HandleGetUserById(c *fiber.Ctx) error {
    var (
        id = c.Params("id")
    )

    user, err := h.userStore.GetUserByID(c.Context(), id)

    if err != nil {
        if errors.Is(err, mongo.ErrNoDocuments) {
            return c.JSON(map[string]string{ "message": "Not found" })
        }
        return c.JSON(err)
    }

    return c.JSON(user)
}
