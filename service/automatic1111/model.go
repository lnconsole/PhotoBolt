package automatic1111

type Sampler string
type SDModel string

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

	SDModelPhotonV1      SDModel = "photon_v1"
	SDModelDreamShaperV7 SDModel = "dreamshaper_7"
)

type ImgOutput struct {
	Images []string `json:"images"`
}

func (sdm SDModel) NegativePrompt() string {
	if sdm == SDModelPhotonV1 {
		return "cartoon, painting, illustration, (worst quality, low quality, normal quality:2)"
	} else {
		return ""
	}
}
