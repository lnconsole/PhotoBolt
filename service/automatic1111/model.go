package automatic1111

type Sampler string

const (
	SamplerEuler         Sampler = "Euler"
	SamplerEulerA        Sampler = "Euler a"
	SamplerDPMPP2M       Sampler = "DPM++ 2M"
	SamplerDPMPP2MKarras Sampler = "DPM++ 2M Karras"
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
	InitImages        []string `json:"init_images"`
	Prompt            string   `json:"prompt"`
	NegativePrompt    string   `json:"negative_prompt"`
	Styles            []string `json:"styles"`
	Steps             int      `json:"steps"`
	Seed              int      `json:"seed"`
	CFGScale          float64  `json:"cfg_scale"`
	SamplerName       Sampler  `json:"sampler_name"`
	ResizeMode        int      `json:"resize_mode"`
	DenoisingStrength float64  `json:"denoising_strength"`
	InpaintingFill    int      `json:"inpainting_fill"`
	Mask              string   `json:"mask"`
}

type Img2ImgInputControlNet struct {
	Img2ImgInput
}

type ImgOutput struct {
	Images []string `json:"images"`
}
