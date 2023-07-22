package automatic1111

import (
	"fmt"

	"github.com/lnconsole/photobolt/http"
)

type Img2ImgInput struct {
	SDModelCheckpoint SDModel  `json:"-"`
	InitImages        []string `json:"init_images"`
	Prompt            string   `json:"prompt"`
	NegativePrompt    string   `json:"negative_prompt"`
	BatchSize         int      `json:"batch_size"`
	Steps             int      `json:"steps"`
	Seed              int      `json:"seed"`
	CFGScale          float64  `json:"cfg_scale"`
	SamplerName       Sampler  `json:"sampler_name"`
	ResizeMode        int      `json:"resize_mode"`
	Width             int      `json:"width"`
	Height            int      `json:"height"`
	DenoisingStrength float64  `json:"denoising_strength"`
}

type Img2ImgInpaintUploadInput struct {
	Img2ImgInput
	Mask                  string   `json:"mask"`
	MaskBlur              int      `json:"mask_blur"`
	Styles                []string `json:"styles"`
	ImageCFGScale         float64  `json:"image_cfg_scale"`
	InpaintingFill        int      `json:"inpainting_fill"`
	InpaintFullRes        int      `json:"inpaint_full_res"`
	InpaintFullResPadding int      `json:"inpaint_full_res_padding"`
	InpaintingMaskInvert  int      `json:"inpainting_mask_invert"`
}

type Img2ImgInpaintUploadControlNetInput struct {
	Img2ImgInpaintUploadInput
	AlwaysOnScripts AlwaysOnScripts `json:"alwayson_scripts"`
}

// NewImg2ImgInput builds a new Img2Img input with default automatic1111 values
func NewImg2ImgInput() *Img2ImgInput {
	return &Img2ImgInput{
		SDModelCheckpoint: SDModelPhotonV1,
		InitImages:        []string{},
		Prompt:            "",
		NegativePrompt:    "",
		BatchSize:         1,
		Steps:             20,
		Seed:              -1,
		CFGScale:          7,
		SamplerName:       SamplerEuler,
		ResizeMode:        ResizeModeJustResize,
		Width:             512,
		Height:            512,
		DenoisingStrength: 0.75,
	}
}

// NewImg2ImgInpaintUploadInput builds a new Img2Img inpaint upload input with default automatic1111 values
func NewImg2ImgInpaintUploadInput() *Img2ImgInpaintUploadInput {
	return &Img2ImgInpaintUploadInput{
		Img2ImgInput:          *NewImg2ImgInput(),
		Mask:                  "",
		MaskBlur:              4,
		Styles:                []string{},
		ImageCFGScale:         1.5,
		InpaintingFill:        InpaintingFillFill,
		InpaintFullRes:        1,
		InpaintFullResPadding: 32,
	}
}

// NewImg2ImgInpaintUploadInput builds a new Img2Img inpaint upload input with default automatic1111 values
func NewImg2ImgInpaintUploadControlNetInput() *Img2ImgInpaintUploadControlNetInput {
	return &Img2ImgInpaintUploadControlNetInput{
		Img2ImgInpaintUploadInput: *NewImg2ImgInpaintUploadInput(),
	}
}

func (i *Img2ImgInpaintUploadControlNetInput) AddControlNetUnit(unit *ControlNetUnit) {
	i.AlwaysOnScripts.ControlNet.Args = append(i.AlwaysOnScripts.ControlNet.Args, unit)
}

// Img2Img send img2img request to an automatic1111 instance listening at automatic1111Url
func Img2Img(automatic1111Url string, input *Img2ImgInput) (*ImgOutput, error) {
	var (
		output = &ImgOutput{}
	)

	if err := SetOptions(automatic1111Url, &SetOptionsInput{SDModelCheckpoint: input.SDModelCheckpoint}); err != nil {
		return nil, err
	}

	if err := http.Post(fmt.Sprintf("%s/sdapi/v1/img2img", automatic1111Url), input, output); err != nil {
		return nil, err
	}

	return output, nil
}

// Img2ImgInpaintUpload send img2img request to an automatic1111 instance listening at automatic1111Url
func Img2ImgInpaintUpload(automatic1111Url string, input *Img2ImgInpaintUploadInput) (*ImgOutput, error) {
	var (
		output = &ImgOutput{}
	)

	if err := SetOptions(automatic1111Url, &SetOptionsInput{SDModelCheckpoint: input.SDModelCheckpoint}); err != nil {
		return nil, err
	}

	if err := http.Post(fmt.Sprintf("%s/sdapi/v1/img2img", automatic1111Url), input, output); err != nil {
		return nil, err
	}

	return output, nil
}

// Img2ImgInpaintUpload send img2img request to an automatic1111 instance listening at automatic1111Url
func Img2ImgInpaintUploadControlNet(automatic1111Url string, input *Img2ImgInpaintUploadControlNetInput) (*ImgOutput, error) {
	var (
		output = &ImgOutput{}
	)

	if err := SetOptions(automatic1111Url, &SetOptionsInput{SDModelCheckpoint: input.SDModelCheckpoint}); err != nil {
		return nil, err
	}

	if err := http.Post(fmt.Sprintf("%s/sdapi/v1/img2img", automatic1111Url), input, output); err != nil {
		return nil, err
	}

	return output, nil
}
