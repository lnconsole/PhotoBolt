package icon

import (
	"fmt"
	"image/png"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lnconsole/photobolt/env"
	srvc "github.com/lnconsole/photobolt/service"
	"github.com/lnconsole/photobolt/service/automatic1111"
	"github.com/lnconsole/photobolt/service/rembg"
	"github.com/lnconsole/photobolt/shared"
)

func Generate(automatic1111Url string) gin.HandlerFunc {
	/*
		img2img
		return img
	*/
	return func(c *gin.Context) {
		var payload GenerateIconBody

		if err := c.ShouldBind(&payload); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		if payload.File == nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Empty file received",
			})
		}

		var (
			file = payload.File
			// Retrieve file information
			extension = filepath.Ext(file.Filename)
			// Generate random file name for the new uploaded file so it doesn't override the old file with same name
			newFileName       = uuid.New().String() + extension
			inputFilelocation = srvc.FileLocation{
				Path: fmt.Sprintf("%s/%s", env.PhotoBolt.RepoDirectory, "api/icon"),
				Name: newFileName,
			}
		)
		if err := c.SaveUploadedFile(file, inputFilelocation.FullPath()); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": fmt.Sprintf("Unable to save the file: %v", err),
			})
			return
		}
		defer func() { inputFilelocation.Remove() }()

		inputFileBytes, err := os.ReadFile(inputFilelocation.FullPath())
		if err != nil {
			log.Printf("error opening whitebg file: %v", err)
			return
		}

		inputFileBase64, err := shared.EncodeImageToBase64(inputFileBytes)
		if err != nil {
			log.Printf("error base64 encoding whitebg bytes: %v", err)
			return
		}

		// all img2img input preparation
		img2imgInput := automatic1111.NewImg2ImgInput()
		img2imgInput.SDModelCheckpoint = automatic1111.SDModelDreamShaperV7
		img2imgInput.Prompt = automatic1111.LoraColoredIcons(0.9, payload.Prompt)
		img2imgInput.SamplerName = automatic1111.SamplerDPMPP2MKarras
		img2imgInput.InitImages = []string{inputFileBase64}

		iconOutput := srvc.FileLocation{
			Path: fmt.Sprintf("%s/%s", env.PhotoBolt.RepoDirectory, "api/icon"),
			Name: uuid.New().String() + ".png",
		}
		imageOutput, err := automatic1111.Img2Img(
			automatic1111Url,
			img2imgInput,
		)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("%v", err),
			})
			return
		}
		img, err := shared.GetImage(imageOutput.Images[0])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("%v", err),
			})
			return
		}

		f, err := os.Create(iconOutput.FullPath())
		png.Encode(f, img)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("%v", err),
			})
			return
		}
		defer func() { os.Remove(iconOutput.FullPath()) }()

		// remove image background
		rembgOutput, err := rembg.RemoveBackground(iconOutput)
		if err != nil {
			log.Printf("rembg err: %s", err)
			return
		}
		defer func() { rembgOutput.Remove() }()

		rembgBytes, err := os.ReadFile(rembgOutput.FullPath())
		if err != nil {
			log.Printf("error opening rembg file: %v", err)
			return
		}

		rembgBase64, err := shared.EncodeImageToBase64(rembgBytes)
		if err != nil {
			log.Printf("error base64 encoding rembg bytes: %v", err)
			return
		}

		split := strings.Split(rembgBase64, ",")

		c.JSON(http.StatusOK, &GenerateIconResponse{
			Image: split[1],
		})
	}
}
