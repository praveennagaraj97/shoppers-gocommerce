package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Locality struct {
	Name          string `json:"name" bson:"name"`
	PostalCode    string `json:"postal_code" bson:"postal_code"`
	StreetAddress string `json:"street_address" bson:"street_address"`
	City          string `json:"city" bson:"city"`
}

type Country struct {
	Name string `json:"name" bson:"name"`
	Code string `json:"code" bson:"code"`
}

type AddressModel struct {
	ID          primitive.ObjectID `bson:"_id" json:"_id"`
	Uid         primitive.ObjectID `json:"uid" bson:"uid"`
	UserName    string             `json:"user_name" bson:"user_name"`
	Building    string             `json:"building" bson:"building"`
	Country     Country            `json:"country" bson:"country"`
	Locality    Locality           `json:"locality" bson:"locality"`
	State       string             `json:"state" bson:"state"`
	Phone       string             `json:"phone" bson:"phone"`
	AddressType string             `json:"address_type" bson:"address_type"`
	CreatedAt   primitive.DateTime `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt   primitive.DateTime `json:"updated_at" bson:"updated_at"`
}
