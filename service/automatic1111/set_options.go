package automatic1111

import (
	"fmt"

	"github.com/lnconsole/photobolt/http"
)

type SetOptionsInput struct {
	SDModelCheckpoint string `json:"sd_model_checkpoint"`
}

func SetOptions(automatic1111Url string, input *SetOptionsInput) error {
	return http.Post(fmt.Sprintf("%s/sdapi/v1/options", automatic1111Url), input, nil)
}
