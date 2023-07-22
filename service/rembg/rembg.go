package rembg

import (
	"fmt"
	"os/exec"

	"github.com/lnconsole/photobolt/env"
	srvc "github.com/lnconsole/photobolt/service"
)

func RemoveBackground(input srvc.FileLocation) (srvc.FileLocation, error) {
	var (
		outputPath = fmt.Sprintf("%s/service/rembg/output", env.PhotoBolt.RepoDirectory)
		outputFile = input.Name
	)
	_, err := exec.Command(
		env.PhotoBolt.RembgFullPath,
		"i",
		"-m",
		"isnet-general-use",
		fmt.Sprintf("%s/%s", input.Path, input.Name),
		fmt.Sprintf("%s/%s", outputPath, outputFile),
	).Output()

	if err != nil {
		return srvc.FileLocation{}, err
	}

	return srvc.FileLocation{
		Path: outputPath,
		Name: outputFile,
	}, nil
}
