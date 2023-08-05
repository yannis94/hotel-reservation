package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yannis94/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct {
    store *db.Store
}

func NewBookingHandler(s *db.Store) *BookingHandler {
    return &BookingHandler{
        store: s,
    }
}

func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
    bookings, err := h.store.Booking.GetBookings(c.Context(), bson.M{})

    if err != nil {
        return err
    }
    return c.JSON(bookings)
}

func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
    id := c.Params("id")
    booking, err := h.store.Booking.GetBookingByID(c.Context(), id)

    if err != nil {
        return err
    }

    return c.JSON(booking)
}
