package automatic1111

import (
	"fmt"

	"github.com/lnconsole/photobolt/http"
)

type Text2ImgInput struct {
	SDModelCheckpoint string  `json:"-"`
	Prompt            string  `json:"prompt"`
	NegativePrompt    string  `json:"negative_prompt"`
	BatchSize         int     `json:"batch_size"`
	Steps             int     `json:"steps"`
	Seed              int     `json:"seed"`
	CFGScale          float64 `json:"cfg_scale"`
	SamplerName       Sampler `json:"sampler_name"`
	Width             int     `json:"width"`
	Height            int     `json:"height"`
}

func NewText2ImgInput() *Text2ImgInput {
	return &Text2ImgInput{
		SDModelCheckpoint: SDModelPhotonV1,
		Prompt:            "",
		NegativePrompt:    "",
		BatchSize:         1,
		Steps:             20,
		Seed:              -1,
		CFGScale:          7,
		SamplerName:       SamplerEuler,
		Width:             512,
		Height:            512,
	}
}

func Text2Img(automatic1111Url string, input *Text2ImgInput) (*ImgOutput, error) {
	var (
		output = &ImgOutput{}
	)

	if err := SetOptions(automatic1111Url, &SetOptionsInput{SDModelCheckpoint: input.SDModelCheckpoint}); err != nil {
		return nil, err
	}

	if err := http.Post(fmt.Sprintf("%s/sdapi/v1/txt2img", automatic1111Url), input, output); err != nil {
		return nil, err
	}

	return output, nil
}
