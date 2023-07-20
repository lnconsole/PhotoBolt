package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strings"

	"github.com/lnconsole/photobolt/api/background"
	"github.com/lnconsole/photobolt/http"
)

const (
	serverUrl = "http://localhost:8080"
)

func main() {
	var (
		inputFilePath = "/Users/ChongjinChua/Downloads/a.png"
		inputPrompt   = "A bouquet of flowers at a beach, surrounded by waves, in front of a volcano"
		output        = &background.ReplaceBackgroundResponse{}
	)

	b, w, err := createMultipartFormData(inputFilePath, inputPrompt)
	if err != nil {
		log.Fatal("create multiform: " + err.Error())
	}

	if err := http.PostBytes(
		fmt.Sprintf("%s/api/background", serverUrl),
		b,
		w.FormDataContentType(),
		output,
	); err != nil {
		log.Fatal("postbytes: " + err.Error())
	}

	outputImage := output.Image[0]
	img, err := getImage(outputImage)
	if err != nil {
		log.Fatal("error: " + err.Error())
	}
	f, err := os.Create("output.png")
	png.Encode(f, img)

	if err != nil {
		log.Fatal("error: " + err.Error())
	}

	log.Println("output png generated: test/output.png")
}

func getImage(imageFromBase64 string) (image.Image, error) {
	img, _, err := image.Decode(base64.NewDecoder(base64.StdEncoding, strings.NewReader(imageFromBase64)))
	return img, err
}

func createMultipartFormData(fileName, prompt string) (bytes.Buffer, *multipart.Writer, error) {
	var b bytes.Buffer
	var err error
	w := multipart.NewWriter(&b)
	// Create a form field writer for field prompt
	promptField, err := w.CreateFormField("prompt")
	if err != nil {
		return bytes.Buffer{}, nil, fmt.Errorf("error creating writer: %v", err)
	}
	// Write prompt field
	promptField.Write([]byte(prompt))
	var fw io.Writer
	file := mustOpen(fileName)
	if fw, err = w.CreateFormFile("file", file.Name()); err != nil {
		return bytes.Buffer{}, nil, fmt.Errorf("error creating writer: %v", err)
	}
	if _, err = io.Copy(fw, file); err != nil {
		return bytes.Buffer{}, nil, fmt.Errorf("error with io.Copy: %v", err)
	}
	w.Close()

	return b, w, nil
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		pwd, _ := os.Getwd()
		fmt.Println("PWD: ", pwd)
		panic(err)
	}
	return r
}
