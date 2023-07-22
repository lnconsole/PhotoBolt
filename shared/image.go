package shared

import (
	"encoding/base64"
	"image"
	"strings"
)

func GetImage(imageFromBase64 string) (image.Image, error) {
	img, _, err := image.Decode(base64.NewDecoder(base64.StdEncoding, strings.NewReader(imageFromBase64)))
	return img, err
}
