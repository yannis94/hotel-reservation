package customtypes

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
    ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    RoomID primitive.ObjectID `bson:"room_id" json:"room_id"`
    UserID primitive.ObjectID `bson:"user_id" json:"user_id"`
    FromDate time.Time `bson:"from_date" json:"from_date"`
    TillDate time.Time `bson:"till_date" json:"till_date"`
    NumPersons int `bson:"num_persons" json:"num_person"`
}
