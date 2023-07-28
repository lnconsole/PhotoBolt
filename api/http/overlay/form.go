package overlay

import "mime/multipart"

type CombineImagesBody struct {
	Front *multipart.FileHeader `form:"front"`
	Back  *multipart.FileHeader `form:"back"`
}

type CombineImagesResponse struct {
	Image string `json:"image"`
}
