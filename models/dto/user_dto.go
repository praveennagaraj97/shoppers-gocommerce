package dto

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserDTO struct {
	FirstName          string             `bson:"first_name" json:"first_name" form:"first_name"`
	LastName           string             `bson:"last_name" json:"last_name" form:"last_name"`
	Email              string             `bson:"email" json:"email" form:"email"`
	Password           string             `bson:"password" json:"password" form:"password"`
	ID                 primitive.ObjectID `bson:"_id"`
	IsActive           bool               `bson:"is_active"`
	RefreshToken       string             `bson:"refresh_token"`
	EmailVerified      bool               `bson:"email_verified"`
	ResetPasswordToken string             `bson:"reset_password_token"`
	JoinedOn           primitive.DateTime `bson:"joined_on"`
	UserRole           string             `bson:"user_role"`
}

func (u *CreateUserDTO) SetData(role string) error {
	byte, err := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
	if err != nil {
		return err
	}

	hashedPass := string(byte)
	u.Password = hashedPass

	u.ID = primitive.NewObjectID()
	u.JoinedOn = primitive.NewDateTimeFromTime(time.Now())
	u.UserRole = role

	return nil
}

type LoginDTO struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type UpdateUserDTO struct {
	FirstName string `form:"first_name,omitempty" json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName  string `form:"last_name,omitempty" json:"last_name,omitempty" bson:"last_name,omitempty"`
	Email     string `form:"email,omitempty" json:"email,omitempty" bson:"email,omitempty"`
}

type UpdatePasswordDTO struct {
	CurrentPassword string `form:"current_password" json:"current_password" bson:"current_password"`
	NewPassword     string `form:"new_password" json:"new_password" bson:"new_password"`
}
