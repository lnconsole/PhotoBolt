package automatic1111

import (
	"fmt"

	"github.com/lnconsole/photobolt/http"
)

type Img2ImgInput struct {
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
	Mask                  string                 `json:"mask"`
	MaskBlur              int                    `json:"mask_blur"`
	Styles                []string               `json:"styles"`
	ImageCFGScale         float64                `json:"image_cfg_scale"`
	InpaintingFill        int                    `json:"inpainting_fill"`
	InpaintFullRes        int                    `json:"inpaint_full_res"`
	InpaintFullResPadding int                    `json:"inpaint_full_res_padding"`
	InpaintingMaskInvert  int                    `json:"inpainting_mask_invert"`
	AlwaysOnScripts       Img2ImgAlwaysOnScripts `json:"alwayson_scripts"`
}

type Img2ImgAlwaysOnScripts struct {
	ControlNet ControlNetScript `json:"controlnet"`
}

type ControlNetScript struct {
	Args []*ControlNetUnit `json:"args"`
}

type ControlNetUnit struct {
	Module        string  `json:"module"`
	Model         string  `json:"model"`
	Weight        float64 `json:"weight"`
	ControlMode   int     `json:"control_mode"`
	ProcessorRes  int     `json:"processor_res"`
	ThresholdA    int     `json:"threshold_a"`
	ThresholdB    int     `json:"threshold_b"`
	InputImage    string  `json:"input_image"`
	GuidanceStart float64 `json:"guidance_start"`
	GuidanceEnd   float64 `json:"guidance_end"`
}

// NewImg2ImgInput builds a new Img2Img input with default automatic1111 values
func NewImg2ImgInput() *Img2ImgInput {
	return &Img2ImgInput{
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

func NewControlNetUnit() *ControlNetUnit {
	return &ControlNetUnit{
		Module:        ControlNetModuleCanny,
		Model:         ControlNetModelCanny,
		Weight:        1.0,
		ControlMode:   ControlNetModeBalanced,
		ProcessorRes:  64,
		ThresholdA:    64,
		ThresholdB:    64,
		InputImage:    "",
		GuidanceStart: 0,
		GuidanceEnd:   1.0,
	}
}

func (i *Img2ImgInpaintUploadInput) AddControlNetUnit(unit *ControlNetUnit) {
	i.AlwaysOnScripts.ControlNet.Args = append(i.AlwaysOnScripts.ControlNet.Args, unit)
}

// Img2Img send img2img request to an automatic1111 instance listening at automatic1111Url
func Img2Img(automatic1111Url string, input *Img2ImgInput) (*ImgOutput, error) {
	var (
		output = &ImgOutput{}
	)

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

	if err := http.Post(fmt.Sprintf("%s/sdapi/v1/img2img", automatic1111Url), input, output); err != nil {
		return nil, err
	}

	return output, nil
}
