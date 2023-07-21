package automatic1111

import (
	"fmt"

	"github.com/lnconsole/photobolt/http"
)

type GetOptionsOutput struct {
	SDModelCheckpoint string `json:"sd_model_checkpoint"`
}

func GetOptions(automatic1111Url string) (*GetOptionsOutput, error) {
	var (
		output = &GetOptionsOutput{}
		err    error
	)

	err = http.Get(fmt.Sprintf("%s/sdapi/v1/options", automatic1111Url), output)

	return output, err
}
