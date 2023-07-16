package automatic1111

import (
	"github.com/lnconsole/photobolt/http"
)

func Text2Img(input *Text2ImgInput) (*ImgOutput, error) {
	var (
		output = &ImgOutput{}
	)

	if err := http.Post("http://127.0.0.1:7860/sdapi/v1/txt2img", input, output); err != nil {
		return nil, err
	}

	return output, nil
}
