package dto

type NewAssetDTO struct {
	Title            string `json:"title" form:"title"`
	Description      string `json:"description" form:"description"`
	ContainerName    string `json:"container_name" form:"container_name" binding:"required"`
	BlurDataRequired bool   `json:"blur_data_required" form:"blur_data_required" strings:"boolean"`
	FileFieldName    string `json:"file_field_name" form:"file_field_name" binding:"required"`
	Published        bool   `json:"published" bson:"published" form:"published"`
}
