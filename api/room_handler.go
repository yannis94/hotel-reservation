package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	customtypes "github.com/yannis94/hotel-reservation/custom_types"
	"github.com/yannis94/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookRoomParams struct {
    FromDate time.Time `json:"from_date"`
    TillDate time.Time `json:"till_date"`
    NumPersons int `json:"num_persons"`
}

type RoomHandler struct {
    store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
    return &RoomHandler{
        store: store,
    }
}

func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
    rooms, err := h.store.Room.GetRooms(c.Context(), bson.M{})

    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(err)
    }

    return c.JSON(rooms)
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
    var params BookRoomParams

    if err := c.BodyParser(&params); err != nil {
        fmt.Println("cannot parse body")
        return c.Status(http.StatusBadRequest).JSON(fmt.Errorf("bad request"))
    }

    if err := params.validate(); err != nil {
        return c.Status(http.StatusBadRequest).JSON(fmt.Errorf("bad request"))
    }

    roomID := c.Params("id")
    roomObjectID, err := primitive.ObjectIDFromHex(roomID)

    if err != nil {
        return c.Status(http.StatusBadRequest).JSON(fmt.Errorf("invalid room id"))
    }

    userID, ok := c.Context().UserValue("user_id").(string)

    if !ok {
        return c.Status(http.StatusInternalServerError).JSON(fmt.Errorf("internal server error"))
    }

    user, err := h.store.User.GetUserByID(c.Context(), userID)

    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fmt.Errorf("internal server error"))
    }

    roomAvailable, err := h.isRoomAvailable(c.Context(), roomObjectID, params)
    fmt.Println(roomAvailable)

    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(err)
    }

    if !roomAvailable {
        return c.Status(http.StatusBadRequest).JSON(fmt.Errorf("room already booked"))
    }

    booking := &customtypes.Booking{
        UserID: user.ID,
        RoomID: roomObjectID,
        FromDate: params.FromDate,
        TillDate: params.TillDate,
        NumPersons: params.NumPersons,
    }

    inserted, err := h.store.Booking.InsertBooking(c.Context(), booking)

    if err != nil {
        return c.Status(http.StatusBadRequest).JSON(fmt.Errorf("internal server error"))
    }
    fmt.Printf("New booking: %+v\n", booking)
    return c.JSON(inserted)
}

func (p BookRoomParams) validate() error {
    now := time.Now()
    if now.After(p.FromDate) || now.After(p.TillDate){
        return fmt.Errorf("cannot book a room in the past")
    }
    return nil
}

func (h *RoomHandler) isRoomAvailable(ctx context.Context, roomID primitive.ObjectID, params BookRoomParams) (bool, error) {
    filter := bson.M{
        "room_id": roomID,
        "from_date": bson.M{
            "$gte": params.FromDate,
        },
        "till_date": bson.M{
            "$lte": params.TillDate,
        },
    }

    bookings, err := h.store.Booking.GetBookings(ctx, filter)

    if err != nil {
        return false, err
    }

    if len(bookings) > 0 {
        return false, nil
    }

    return true, nil
}
