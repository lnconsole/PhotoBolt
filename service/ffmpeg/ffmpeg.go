package ffmpeg

import (
	"fmt"
	"os/exec"

	"github.com/lnconsole/photobolt/env"
	srvc "github.com/lnconsole/photobolt/service"
)

func ConvertToMask(input srvc.FileLocation) (srvc.FileLocation, error) {
	var (
		outputPath = outputPath()
		outputFile = fmt.Sprintf("bare-mask-%s", input.Name)
	)
	_, err := exec.Command(
		env.PhotoBolt.FfmpegFullPath,
		"-i",
		fmt.Sprintf("%s/%s", input.Path, input.Name),
		"-vf",
		"lutrgb=r=0:g=0:b=0",
		"-y",
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

func InsertWhiteBackground(input srvc.FileLocation) (srvc.FileLocation, error) {
	var (
		outputPath = outputPath()
		outputFile = fmt.Sprintf("white-%s", input.Name)
	)
	_, err := exec.Command(
		env.PhotoBolt.FfmpegFullPath,
		"-hide_banner",
		"-i",
		fmt.Sprintf("%s/%s", input.Path, input.Name),
		"-filter_complex",
		"color=white[c];[c][0]scale2ref[cs][0s];[cs][0s]overlay=shortest=1",
		"-y",
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

func OverlayImages(front srvc.FileLocation, back srvc.FileLocation) (srvc.FileLocation, error) {
	var (
		outputPath = outputPath()
		outputFile = fmt.Sprintf("combined-%s", back.Name)
	)
	_, err := exec.Command(
		env.PhotoBolt.FfmpegFullPath,
		"-i",
		fmt.Sprintf("%s/%s", back.Path, back.Name),
		"-i",
		fmt.Sprintf("%s/%s", front.Path, front.Name),
		"-filter_complex",
		"[0:v][1:v] overlay=0:0",
		"-y",
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

func outputPath() string {
	return fmt.Sprintf("%s/service/ffmpeg/output", env.PhotoBolt.RepoDirectory)
}
