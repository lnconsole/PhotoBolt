package background

import (
	"fmt"
	"image/png"
	"strings"

	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/lnconsole/photobolt/env"
	srvc "github.com/lnconsole/photobolt/service"
	"github.com/lnconsole/photobolt/service/automatic1111"
	"github.com/lnconsole/photobolt/service/ffmpeg"
	"github.com/lnconsole/photobolt/service/rembg"
	"github.com/lnconsole/photobolt/shared"
)

func Replace(automatic1111Url string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload ReplaceBackgroundBody

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
				Path: fmt.Sprintf("%s/%s", env.PhotoBolt.RepoDirectory, "api/background"),
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

		// remove image background
		rembgOutput, err := rembg.RemoveBackground(srvc.FileLocation{
			Path: inputFilelocation.Path,
			Name: inputFilelocation.Name,
		})
		if err != nil {
			log.Printf("rembg err: %s", err)
			return
		}
		defer func() { rembgOutput.Remove() }()

		// add white background to backgroundless image
		whitebg, err := ffmpeg.InsertWhiteBackground(srvc.FileLocation{
			Path: rembgOutput.Path,
			Name: rembgOutput.Name,
		})
		if err != nil {
			log.Printf("whitebg err: %s", err)
			return
		}
		defer func() { whitebg.Remove() }()

		whiteBgFileBytes, err := os.ReadFile(whitebg.FullPath())
		if err != nil {
			log.Printf("error opening whitebg file: %v", err)
			return
		}

		whiteBgFileBase64, err := shared.EncodeImageToBase64(whiteBgFileBytes)
		if err != nil {
			log.Printf("error base64 encoding whitebg bytes: %v", err)
			return
		}

		split := strings.Split(whiteBgFileBase64, ",")

		inputImg, err := shared.GetImage(split[1])
		if err != nil {
			log.Printf("error getimage: %v", err)
			return
		}
		// all txt2img input preparation
		txt2img := automatic1111.NewText2ImgControlNetInput()
		txt2img.SDModelCheckpoint = automatic1111.SDModelPhotonV1
		txt2img.Prompt = payload.Prompt
		txt2img.NegativePrompt = automatic1111.SDModelPhotonV1.NegativePrompt()
		txt2img.BatchSize = 1
		txt2img.Steps = 25
		txt2img.Seed = -1
		txt2img.CFGScale = 3
		txt2img.SamplerName = automatic1111.SamplerDPMPP2M
		txt2img.Width = inputImg.Bounds().Dx()
		txt2img.Height = inputImg.Bounds().Dy()

		cannyCNUnit := automatic1111.NewControlNetUnit()
		cannyCNUnit.InputImage = whiteBgFileBase64
		cannyCNUnit.Weight = 2
		cannyCNUnit.ControlMode = automatic1111.ControlNetModeBalanced
		cannyCNUnit.ProcessorRes = inputImg.Bounds().Dx()
		cannyCNUnit.ThresholdA = 100
		cannyCNUnit.ThresholdB = 200
		txt2img.AddControlNetUnit(cannyCNUnit)

		imageOutput, err := automatic1111.Text2ImgControlNet(
			automatic1111Url,
			txt2img,
		)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("%v", err),
			})
			return
		}

		replacedOutput := srvc.FileLocation{
			Path: fmt.Sprintf("%s/%s", env.PhotoBolt.RepoDirectory, "api/background"),
			Name: uuid.New().String() + ".png",
		}

		img, err := shared.GetImage(imageOutput.Images[0])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("%v", err),
			})
			return
		}
		f, err := os.Create(replacedOutput.FullPath())
		png.Encode(f, img)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("%v", err),
			})
			return
		}
		defer func() { os.Remove(replacedOutput.FullPath()) }()

		// overlay target object on generated background
		overlayOutput, err := ffmpeg.OverlayImages(
			rembgOutput,
			replacedOutput,
			txt2img.Width, txt2img.Height,
		)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": fmt.Sprintf("overlayOutput err: %v", err),
			})
			return
		}
		defer func() { overlayOutput.Remove() }()

		overlayFileBytes, err := os.ReadFile(overlayOutput.FullPath())
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": fmt.Sprintf("error opening overlayOutput file: %v", err),
			})
			return
		}

		overlayFileBase64, err := shared.EncodeImageToBase64(overlayFileBytes)
		if err != nil {
			log.Printf("error base64 encoding overlayOutput bytes: %v", err)
			return
		}

		split = strings.Split(overlayFileBase64, ",")

		c.JSON(http.StatusOK, &ReplaceBackgroundResponse{
			Image: []string{split[1]},
		})
	}
}
