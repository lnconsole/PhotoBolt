package automatic1111

import "fmt"

func LoraColoredIcons(strength float64, prompt string) string {
	return fmt.Sprintf("<lora:Colored_Icons:%f> coloredic0n icon %s", strength, prompt)
}
