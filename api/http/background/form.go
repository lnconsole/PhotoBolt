package background

import "mime/multipart"

type ReplaceBackgroundBody struct {
	File   *multipart.FileHeader `form:"file"`
	Prompt string                `form:"prompt"`
}

type ReplaceBackgroundResponse struct {
	Image []string `json:"image"`
}
