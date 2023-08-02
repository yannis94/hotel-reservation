package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yannis94/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelHandler struct {
    store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
    return &HotelHandler{
        store: store,
    }
}

func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
    hotelID := c.Params("id")
    objectID, err := primitive.ObjectIDFromHex(hotelID)
    
    if err != nil {
        return c.JSON(err)
    }

    filter := bson.M{ "hotel_id": objectID }
    rooms, err := h.store.Room.GetRooms(c.Context(), filter)
    
    if err != nil {
        return c.JSON(err)
    }

    return c.JSON(rooms)
}

func (h *HotelHandler) HandleGetHotelByID(c *fiber.Ctx) error {
    hotelID := c.Params("id")
    objectID, err := primitive.ObjectIDFromHex(hotelID)
    
    if err != nil {
        return c.JSON(err)
    }

    filter := bson.M{ "_id": objectID }
    hotel, err := h.store.Hotel.GetHotelByID(c.Context(), filter)
    
    if err != nil {
        return c.JSON(err)
    }

    return c.JSON(hotel)
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
    hotels, err := h.store.Hotel.GetHotels(c.Context(), nil)

    if err != nil {
        return c.JSON(err)
    }

    return c.JSON(hotels)
}
