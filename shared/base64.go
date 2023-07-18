package shared

import (
	"encoding/base64"
	"errors"
	"net/http"
)

func EncodeImageToBase64(bytes []byte) (string, error) {
	var (
		mimeType         = http.DetectContentType(bytes)
		base64EncodedImg string
	)

	switch mimeType {
	case "image/jpeg":
		base64EncodedImg += "data:image/jpeg;base64,"
	case "image/png":
		base64EncodedImg += "data:image/png;base64,"
	default:
		return "", errors.New("unsupported image format")
	}

	base64EncodedImg += base64.StdEncoding.EncodeToString(bytes)

	return base64EncodedImg, nil
}
