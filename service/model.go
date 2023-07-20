package srvc

import (
	"fmt"
	"os"
)

type FileLocation struct {
	Path string
	Name string
}

func (f FileLocation) FullPath() string {
	return fmt.Sprintf("%s/%s", f.Path, f.Name)
}

func (f FileLocation) Remove() error {
	return os.Remove(f.FullPath())
}
