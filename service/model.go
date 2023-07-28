package srvc

import (
	"fmt"
	"image/png"
	"os"
	"strings"

	"github.com/lnconsole/photobolt/shared"
)

type FileLocation struct {
	Path string
	Name string
}

func (f FileLocation) FullPath() string {
	return fmt.Sprintf("%s/%s", f.Path, f.Name)
}

func (f FileLocation) Remove() error {
	return os.Remove(f.FullPath())
}

func (loc FileLocation) SavePNG(base64data string) error {
	input := base64data
	split := strings.Split(base64data, ",")
	if len(split) > 1 {
		input = split[1]
	}
	img, err := shared.GetImage(input)
	if err != nil {
		return err
	}
	f, err := os.Create(loc.FullPath())
	if err != nil {
		return err
	}
	png.Encode(f, img)

	return nil
}

func (loc FileLocation) ToBase64() (string, error) {
	fileBytes, err := os.ReadFile(loc.FullPath())
	if err != nil {
		return "", err
	}
	fileBase64, err := shared.EncodeImageToBase64(fileBytes)
	if err != nil {
		return "", err
	}

	output := fileBase64
	split := strings.Split(fileBase64, ",")
	if len(split) > 0 {
		output = split[1]
	}

	return output, nil
}
