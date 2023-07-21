package automatic1111

import "fmt"

const (
	LoraColoredIcons = "Colored_Icons_by_vizsumit"
)

type Lora struct {
	Name     string
	Strength int
}

func (l *Lora) String() string {
	return fmt.Sprintf("<lora:%s:%d>", l.Name, l.Strength)
}
