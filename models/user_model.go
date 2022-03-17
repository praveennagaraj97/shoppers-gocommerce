package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                 primitive.ObjectID `json:"_id"`
	FirstName          string             `json:"first_name"`
	LastName           string             `json:"last_name"`
	Email              string             `json:"email"`
	Password           string             `json:"-"`
	IsActive           bool               `json:"-"`
	RefreshToken       string             `json:"-" `
	EmailVerified      bool               `json:"email_verified"`
	ResetPasswordToken string             `json:"-"`
	JoinedOn           primitive.DateTime `json:"joined_on"`
	UserRole           string             `json:"-"`
}

func (u *User) HashPassword() error {

	byte, err := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
	if err != nil {
		return err
	}

	hashedPass := string(byte)

	u.Password = hashedPass

	return nil
}

func (u *User) CompareHashedPassword(pwd string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pwd))
}
