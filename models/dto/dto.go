package dto

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
	FirstName             string `form:"first_name,omitempty" json:"first_name,omitempty"`
	LastName              string `form:"last_name,omitempty" json:"last_name,omitempty"`
	StreetNumber          string `form:"street_number,omitempty" json:"street_number,omitempty"`
	Building              string `form:"building,omitempty" json:"building,omitempty"`
	Phone                 string `form:"phone,omitempty" json:"phone,omitempty"`
	State                 string `form:"state,omitempty" json:"state,omitempty"`
	AddressType           string `form:"address_type,omitempty" json:"address_type,omitempty"`
	CountryName           string `form:"country_name,omitempty" json:"country_name,omitempty"`
	CountryCode           string `form:"country_code,omitempty" json:"country_code,omitempty"`
	LocalityName          string `form:"locality_name,omitempty" json:"locality_name,omitempty"`
	LocalityPostalCode    string `form:"locality_postal_code,omitempty" json:"locality_postal_code,omitempty"`
	LocalityStreetAddress string `form:"locality_street_address,omitempty" json:"locality_street_address,omitempty"`
	LocalityCity          string `form:"locality_city,omitempty" json:"locality_city,omitempty"`
}

type NewAssetDTO struct {
	Title            string `json:"title" form:"title"`
	Description      string `json:"description" form:"description"`
	ContainerName    string `json:"container_name" form:"container_name" binding:"required"`
	BlurDataRequired bool   `json:"blur_data_required" form:"blur_data_required" strings:"boolean"`
	FileFieldName    string `json:"file_field_name" form:"file_field_name" binding:"required"`
	Published        bool   `json:"published" bson:"published" form:"published"`
}
