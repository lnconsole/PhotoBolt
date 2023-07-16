package overlay

import "mime/multipart"

type CombineImages struct {
	Front *multipart.FileHeader `form:"front"`
	Back  *multipart.FileHeader `form:"back"`
}
