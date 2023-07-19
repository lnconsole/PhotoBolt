package automatic1111

type Sampler string

const (
	SamplerEuler         Sampler = "Euler"
	SamplerEulerA        Sampler = "Euler a"
	SamplerDPMPP2M       Sampler = "DPM++ 2M"
	SamplerDPMPP2MKarras Sampler = "DPM++ 2M Karras"

	ResizeModeJustResize              = 0
	ResizeModeCropAndResize           = 1
	ResizeModeResizeAndFill           = 2
	ResizeModeJustResizeLatentUpscale = 3

	InpaintingFillFill          = 0
	InpaintingFillOriginal      = 1
	InpaintingFillLatentNoise   = 2
	InpaintingFillLatentNothing = 3

	ControlNetModuleCanny = "canny"

	ControlNetModelCanny = "control_v11p_sd15_canny"

	ControlNetModeBalanced                  = 0
	ControlNetModeMyPromptIsMoreImportant   = 1
	ControlNetModeControlNetIsMoreImportant = 2
)

type Text2ImgInput struct {
	Prompt         string   `json:"prompt"`
	NegativePrompt string   `json:"negative_prompt"`
	Styles         []string `json:"styles"`
	Steps          int      `json:"steps"`
	Seed           int      `json:"seed"`
	CFGScale       float64  `json:"cfg_scale"`
	SamplerName    Sampler  `json:"sampler_name"`
}

type Img2ImgInput struct {
	InitImages            []string               `json:"init_images"`
	Mask                  string                 `json:"mask"`
	MaskBlur              int                    `json:"mask_blur"`
	Prompt                string                 `json:"prompt"`
	NegativePrompt        string                 `json:"negative_prompt"`
	Styles                []string               `json:"styles"`
	BatchSize             int                    `json:"batch_size"`
	Steps                 int                    `json:"steps"`
	Seed                  int                    `json:"seed"`
	CFGScale              float64                `json:"cfg_scale"`
	ImageCFGScale         float64                `json:"image_cfg_scale"`
	SamplerName           Sampler                `json:"sampler_name"`
	ResizeMode            int                    `json:"resize_mode"`
	Width                 int                    `json:"width"`
	Height                int                    `json:"height"`
	DenoisingStrength     float64                `json:"denoising_strength"`
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

type ImgOutput struct {
	Images []string `json:"images"`
}

// NewImg2ImgInput builds a new Img2Img input with default automatic1111 values
func NewImg2ImgInput() *Img2ImgInput {
	return &Img2ImgInput{
		InitImages:            []string{},
		Mask:                  "",
		MaskBlur:              4,
		Prompt:                "",
		NegativePrompt:        "",
		Styles:                []string{},
		BatchSize:             1,
		Steps:                 20,
		Seed:                  -1,
		CFGScale:              7,
		ImageCFGScale:         1.5,
		SamplerName:           SamplerEuler,
		ResizeMode:            ResizeModeJustResize,
		Width:                 512,
		Height:                512,
		DenoisingStrength:     0.75,
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

func (i *Img2ImgInput) AddControlNetUnit(unit *ControlNetUnit) {
	i.AlwaysOnScripts.ControlNet.Args = append(i.AlwaysOnScripts.ControlNet.Args, unit)
}
