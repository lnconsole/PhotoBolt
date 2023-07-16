package automatic1111

import (
	"fmt"
	"github.com/lnconsole/photobolt/http"
)

func Text2Img(automatic1111Url string, input *Text2ImgInput) (*ImgOutput, error) {
	var (
		output = &ImgOutput{}
	)

	if err := http.Post(fmt.Sprintf("%s/sdapi/v1/txt2img", automatic1111Url), input, output); err != nil {
		return nil, err
	}

	return output, nil
}
