package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                 primitive.ObjectID `bson:"_id" json:"_id"`
	FirstName          string             `bson:"first_name" json:"first_name"`
	LastName           string             `bson:"last_name" json:"last_name"`
	Email              string             `bson:"email" json:"email"`
	Password           string             `bson:"password" json:"-"`
	IsActive           bool               `bson:"is_active" json:"-"`
	RefreshToken       string             `bson:"refresh_token" json:"-"`
	EmailVerified      bool               `bson:"email_verified" json:"email_verified"`
	ResetPasswordToken string             `bson:"reset_password_token" json:"-"`
	JoinedOn           primitive.DateTime `bson:"joined_on" json:"joined_on"`
	UserRole           string             `bson:"user_role" json:"user_role"`
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
