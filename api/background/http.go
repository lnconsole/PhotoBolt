package background

import (
	"fmt"

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
	/*
		remove bg
		add white bg
		generate mask and add white bg
		call img2img
		return img
	*/
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

		// convert backgroundless image to mask
		maskOutput, err := ffmpeg.ConvertToMask(srvc.FileLocation{
			Path: rembgOutput.Path,
			Name: rembgOutput.Name,
		})
		if err != nil {
			log.Printf("maskoutput err: %s", err)
			return
		}
		defer func() { maskOutput.Remove() }()

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

		// add white background to mask image
		maskwhitebg, err := ffmpeg.InsertWhiteBackground(srvc.FileLocation{
			Path: maskOutput.Path,
			Name: maskOutput.Name,
		})
		if err != nil {
			log.Printf("maskwhitebg err: %s", err)
			return
		}
		defer func() { maskwhitebg.Remove() }()

		log.Printf(
			"Ready for automatic1111\nPic with white background: %s\nMask with white background: %s\n",
			whitebg.FullPath(),
			maskwhitebg.FullPath(),
		)

		whiteBgFileBytes, err := os.ReadFile(whitebg.FullPath())
		if err != nil {
			log.Printf("error opening whitebg file: %v", err)
			return
		}

		maskFileBytes, err := os.ReadFile(maskwhitebg.FullPath())
		if err != nil {
			log.Printf("error opening mask file: %v", err)
			return
		}

		whiteBgFileBase64, err := shared.EncodeImageToBase64(whiteBgFileBytes)
		if err != nil {
			log.Printf("error base64 encoding whitebg bytes: %v", err)
			return
		}

		maskFileBase64, err := shared.EncodeImageToBase64(maskFileBytes)
		if err != nil {
			log.Printf("error base64 encoding mask bytes: %v", err)
			return
		}

		// all img2img input preparation
		img2imgInput := automatic1111.NewImg2ImgInpaintUploadInput()
		img2imgInput.Prompt = payload.Prompt
		img2imgInput.SamplerName = automatic1111.SamplerDPMPP2MKarras
		img2imgInput.InitImages = []string{whiteBgFileBase64}
		img2imgInput.Mask = maskFileBase64
		img2imgInput.InpaintFullResPadding = 40
		img2imgInput.MaskBlur = 0
		img2imgInput.DenoisingStrength = 0.75
		img2imgInput.CFGScale = 6.0
		img2imgInput.Steps = 30

		cannyCNUnit := automatic1111.NewControlNetUnit()
		cannyCNUnit.InputImage = whiteBgFileBase64
		cannyCNUnit.Weight = 2
		cannyCNUnit.ControlMode = 2
		cannyCNUnit.ProcessorRes = 512
		cannyCNUnit.ThresholdA = 100
		cannyCNUnit.ThresholdB = 200
		img2imgInput.AddControlNetUnit(cannyCNUnit)

		imageOutput, err := automatic1111.Img2ImgInpaintUpload(
			automatic1111Url,
			img2imgInput,
		)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("%v", err),
			})
			return
		}

		c.JSON(http.StatusOK, &ReplaceBackgroundResponse{
			Image: imageOutput.Images,
		})
	}
}

// overlay logo on white background image
// overlay, err := ffmpeg.OverlayImages(
// 	srvc.FileLocation{
// 		Path: inputPath,
// 		Name: "beatzcoin.png",
// 	},
// 	srvc.FileLocation{
// 		Path: whitebg.Path,
// 		Name: whitebg.Name,
// 	},
// )
// if err != nil {
// 	log.Printf("overlay err: %s", err)
// 	return
// }
