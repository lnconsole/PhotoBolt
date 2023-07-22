package main

import (
	"fmt"
	"image/png"
	"log"
	"os"

	"github.com/lnconsole/photobolt/api/background"
	"github.com/lnconsole/photobolt/api/icon"
	"github.com/lnconsole/photobolt/api/overlay"
	"github.com/lnconsole/photobolt/env"
	"github.com/lnconsole/photobolt/http"
	"github.com/lnconsole/photobolt/shared"
)

const (
	// CHANGE ME
	replacedBg_inputFile   = "/Users/ChongjinChua/Downloads/whisky.png"
	replacedBg_inputPrompt = "A bouquet of flowers at a beach, surrounded by waves, in front of a volcano"
	icon_inputFile         = "/Users/ChongjinChua/Downloads/ibex.png"
	icon_inputPrompt       = "ibex"

	// CONST
	serverUrl             = "http://localhost:8080"
	replacedBg_outputFile = "replaced-bg.png"
	icon_outputFile       = "icon.png"
	combined_outputFile   = "combined.png"
)

func main() {
	if err := env.Init("../env/.env"); err != nil {
		log.Printf("env err: %s", err)
		return
	}

	if err := replaceBackground(); err != nil {
		log.Fatalf(err.Error())
	}

	if err := generateIcon(); err != nil {
		log.Fatalf(err.Error())
	}

	if err := combineImages(); err != nil {
		log.Fatalf(err.Error())
	}
}

func replaceBackground() error {
	var (
		inputFilePath = replacedBg_inputFile
		inputPrompt   = replacedBg_inputPrompt
		output        = &background.ReplaceBackgroundResponse{}
	)

	if err := http.PostForm(
		fmt.Sprintf("%s/api/background", serverUrl),
		map[string]string{
			"file":   inputFilePath,
			"prompt": inputPrompt,
		},
		output,
	); err != nil {
		return fmt.Errorf("postForm: %s", err.Error())
	}

	if len(output.Image) == 0 {
		return fmt.Errorf("no image generated")
	}

	outputImage := output.Image[0]
	img, err := shared.GetImage(outputImage)
	if err != nil {
		return fmt.Errorf("getimage: %s", err.Error())
	}

	f, err := os.Create(replacedBg_outputFile)
	png.Encode(f, img)
	if err != nil {
		return fmt.Errorf("create: %s", err.Error())
	}

	log.Printf("replaced background png generated: %s/client/%s", env.PhotoBolt.RepoDirectory, replacedBg_outputFile)
	return nil
}

func generateIcon() error {
	var (
		inputFilePath = icon_inputFile
		inputPrompt   = icon_inputPrompt
		output        = &icon.GenerateIconResponse{}
	)

	if err := http.PostForm(
		fmt.Sprintf("%s/api/icon", serverUrl),
		map[string]string{
			"file":   inputFilePath,
			"prompt": inputPrompt,
		},
		output,
	); err != nil {
		return fmt.Errorf("postForm: %s", err.Error())
	}

	if len(output.Image) == 0 {
		return fmt.Errorf("no image generated")
	}

	outputImage := output.Image
	img, err := shared.GetImage(outputImage)
	if err != nil {
		return fmt.Errorf("getimage: %s", err.Error())
	}

	f, err := os.Create(icon_outputFile)
	png.Encode(f, img)
	if err != nil {
		return fmt.Errorf("create: %s", err.Error())
	}

	log.Printf("replaced background png generated: %s/client/%s", env.PhotoBolt.RepoDirectory, icon_outputFile)
	return nil
}

func combineImages() error {
	var (
		frontFilePath = fmt.Sprintf("%s/client/%s", env.PhotoBolt.RepoDirectory, icon_outputFile)
		backFilePath  = fmt.Sprintf("%s/client/%s", env.PhotoBolt.RepoDirectory, replacedBg_outputFile)
		output        = &overlay.CombineImagesResponse{}
	)

	if err := http.PostForm(
		fmt.Sprintf("%s/api/overlay", serverUrl),
		map[string]string{
			"front": frontFilePath,
			"back":  backFilePath,
		},
		output,
	); err != nil {
		return fmt.Errorf("postForm: %s", err.Error())
	}

	outputImage := output.Image
	img, err := shared.GetImage(outputImage)
	if err != nil {
		return fmt.Errorf("getimage: %s", err.Error())
	}

	f, err := os.Create(combined_outputFile)
	png.Encode(f, img)
	if err != nil {
		return fmt.Errorf("create: %s", err.Error())
	}

	log.Printf("combined png generated: %s/client/%s", env.PhotoBolt.RepoDirectory, combined_outputFile)
	return nil
}
