package db

const (
    DbUri = "mongodb://localhost:27017"
    DbName = "hotel-reservation"
    userColl = "users"
    hotelColl = "hotels"
    roomColl = "rooms"
    bookingColl = "bookings"
)

type Store struct {
    User UserStore
    Hotel HotelStore
    Room RoomStore
    Booking BookingStore
}
