package dto

import (
	"fmt"
	"time"

	"github.com/praveennagaraj97/shoppers-gocommerce/models"
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

type UserAddressDTO struct {
	ID          primitive.ObjectID  `bson:"_id,omitempty"`
	UId         *primitive.ObjectID `bson:"uid"`
	FirstName   string              `form:"first_name,omitempty" json:"first_name,omitempty" bson:"first_name"`
	LastName    string              `form:"last_name,omitempty" json:"last_name,omitempty" bson:"last_name"`
	Building    string              `form:"building,omitempty" json:"building,omitempty" bson:"building"`
	Phone       string              `form:"phone,omitempty" json:"phone,omitempty" bson:"phone"`
	State       string              `form:"state,omitempty" json:"state,omitempty" bson:"state"`
	AddressType string              `form:"address_type,omitempty" json:"address_type,omitempty" bson:"address_type"`

	CountryName           string `form:"country_name,omitempty" json:"country_name,omitempty" bson:"-"`
	CountryCode           string `form:"country_code,omitempty" json:"country_code,omitempty" bson:"-"`
	LocalityName          string `form:"locality_name,omitempty" json:"locality_name,omitempty" bson:"-"`
	LocalityPostalCode    string `form:"locality_postal_code,omitempty" json:"locality_postal_code,omitempty" bson:"-"`
	LocalityStreetAddress string `form:"locality_street_address,omitempty" json:"locality_street_address,omitempty" bson:"-"`
	LocalityCity          string `form:"locality_city,omitempty" json:"locality_city,omitempty" bson:"-"`

	UserName  string              `bson:"user_name"`
	Country   *models.Country     `bson:"country"`
	Locality  *models.Locality    `bson:"locality"`
	CreatedAt *primitive.DateTime `bson:"created_at,omitempty"`
	UpdatedAt *primitive.DateTime `bson:"updated_at"`
}

func (a *UserAddressDTO) SetData(userID *primitive.ObjectID, isUpdate bool) {

	a.UId = userID
	a.UserName = fmt.Sprintf("%s %s", a.FirstName, a.LastName)
	a.Country = &models.Country{
		Name: a.CountryName,
		Code: a.CountryCode,
	}
	a.Locality = &models.Locality{
		Name:          a.LocalityName,
		PostalCode:    a.LocalityPostalCode,
		StreetAddress: a.LocalityStreetAddress,
		City:          a.LocalityCity,
	}

	currentTime := primitive.NewDateTimeFromTime(time.Now())

	if !isUpdate {

		a.ID = primitive.NewObjectID()
		a.CreatedAt = &currentTime
	}

	a.UpdatedAt = &currentTime

}
