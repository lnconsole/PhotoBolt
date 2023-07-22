package automatic1111

type AlwaysOnScripts struct {
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
