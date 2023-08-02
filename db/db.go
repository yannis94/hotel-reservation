package db

const (
    DbUri = "mongodb://localhost:27017"
    DbName = "hotel-reservation"
    userColl = "users"
    hotelColl = "hotels"
    roomColl = "rooms"
)

type Store struct {
    User UserStore
    Hotel HotelStore
    Room RoomStore
}
