package icon

import "mime/multipart"

type GenerateIconBody struct {
	File   *multipart.FileHeader `form:"file"`
	Prompt string                `form:"prompt"`
}

type GenerateIconResponse struct {
	Image string `json:"image"`
}
