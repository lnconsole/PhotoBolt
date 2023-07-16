package main

import (
	"log"

	"github.com/lnconsole/photobolt/env"
	srvc "github.com/lnconsole/photobolt/service"
	"github.com/lnconsole/photobolt/service/ffmpeg"
	"github.com/lnconsole/photobolt/service/rembg"
)

const (
	inputPath = "/Users/ChongjinChua/Downloads" // /path/to/your/file/directory
	inputName = "bottle-keras.png"              // name.png
)

func main() {
	log.Printf("winning")

	if err := env.Init(); err != nil {
		log.Printf("env err: %s", err)
		return
	}

	if inputPath == "" || inputName == "" {
		log.Printf("Please provide input")
		return
	}

	rembgOutput, err := rembg.RemoveBackground(srvc.FileLocation{
		Path: inputPath,
		Name: inputName,
	})
	if err != nil {
		log.Printf("rembg err: %s", err)
		return
	}

	maskOutput, err := ffmpeg.ConvertToMask(srvc.FileLocation{
		Path: rembgOutput.Path,
		Name: rembgOutput.Name,
	})
	if err != nil {
		log.Printf("maskoutput err: %s", err)
		return
	}

	whitebg, err := ffmpeg.InsertWhiteBackground(srvc.FileLocation{
		Path: rembgOutput.Path,
		Name: rembgOutput.Name,
	})
	if err != nil {
		log.Printf("whitebg err: %s", err)
		return
	}

	maskwhitebg, err := ffmpeg.InsertWhiteBackground(srvc.FileLocation{
		Path: maskOutput.Path,
		Name: maskOutput.Name,
	})
	if err != nil {
		log.Printf("maskwhitebg err: %s", err)
		return
	}

	log.Printf(
		"Ready for automatic1111\nPic with white background: %s/%s\nMask with white background: %s/%s\n",
		whitebg.Path, whitebg.Name,
		maskwhitebg.Path, maskwhitebg.Name,
	)
}
