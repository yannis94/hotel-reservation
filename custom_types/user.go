package customtypes

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
    bcryptCost = 12
    minFirstNameLen = 2
    minLastNameLen = 2
    minPasswordLength = 7
)

type CreateUserParams struct {
    FirstName string `json:"first_name"`
    LastName string `json:"last_name"`
    Email string `json:"email"`
    Password string `json:"password"`
}

type UpdateUserParams struct {
    FirstName string `json:"first_name"`
    LastName string `json:"last_name"`
}

type User struct {
    ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    FirstName string `bson:"first_name" json:"first_name"`
    LastName string `bson:"last_name" json:"last_name"`
    Email string `bson:"email" json:"email"`
    Password string `bson:"password" json:"-"`
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
    encryptPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
    
    if err != nil {
        return nil, err
    }

    return &User{
        FirstName: params.FirstName,
        LastName: params.LastName,
        Email: params.Email,
        Password: string(encryptPassword),
    }, nil
}

func (params UpdateUserParams) ToBSON() bson.M {
    document := bson.M{}

    if len(params.FirstName) > 0 {
        document["first_name"] = params.FirstName
    }
    if len(params.LastName) > 0 {
        document["last_name"] = params.LastName
    }

    return document
}

func (p *CreateUserParams) Validate() map[string]string {
    var errors = make(map[string]string)

    if len(p.FirstName) < minFirstNameLen {
        err := fmt.Sprintf("First name length should be at least %d characters", minFirstNameLen)
        errors["first_name"] = err
    }
    if len(p.LastName) < minLastNameLen {
        err := fmt.Sprintf("Last name length should be at least %d characters", minLastNameLen)
        errors["last_name"] = err
    }
    if len(p.Password) < minPasswordLength {
        err := fmt.Sprintf("Password's length should be at least %d characters", minPasswordLength)
        errors["password"] =  err
    }
    if !isEmailValid(p.Email) {
        err := fmt.Sprintf("Email address invalid.")
        errors["email"] = err
    }

    return errors
}

func isEmailValid(email string) bool {
    emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}
