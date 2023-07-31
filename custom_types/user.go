package customtypes

type User struct {
    ID string `bson:"_id" json:"_"`
    FirstName string `bson:"first_name" json:"first_name"`
    LastName string `bson:"last_name" json:"last_name"`
}
