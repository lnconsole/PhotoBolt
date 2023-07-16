package automatic1111

import (
	"github.com/lnconsole/photobolt/http"
)

func Img2Img(input *Img2ImgInput) (*ImgOutput, error) {
	var (
		output = &ImgOutput{}
	)

	if err := http.Post("http://127.0.0.1:7860/sdapi/v1/img2img", input, output); err != nil {
		return nil, err
	}

	return output, nil
}
