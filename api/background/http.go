package background

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lnconsole/photobolt/env"
	srvc "github.com/lnconsole/photobolt/service"
	"github.com/lnconsole/photobolt/service/ffmpeg"
	"github.com/lnconsole/photobolt/service/rembg"
)

func Replace(c *gin.Context) {
	/*
		remove bg
		add white bg
		generate mask and add white bg
		call img2img
		return img
	*/
	var payload ReplaceBackgroundBody

	if err := c.ShouldBind(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": "ERROR",
			"reason": err.Error(),
		})
		return
	}

	var (
		file = payload.File
		// Retrieve file information
		extension = filepath.Ext(file.Filename)
		// Generate random file name for the new uploaded file so it doesn't override the old file with same name
		newFileName = uuid.New().String() + extension
		inputPath   = fmt.Sprintf("%s/%s", env.PhotoBolt.RepoDirectory, "api/background")
		inputName   = newFileName
	)
	if err := c.SaveUploadedFile(file, fmt.Sprintf("%s/%s", inputPath, inputName)); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to save the file",
		})
		return
	}

	// remove image background
	rembgOutput, err := rembg.RemoveBackground(srvc.FileLocation{
		Path: inputPath,
		Name: inputName,
	})
	if err != nil {
		log.Printf("rembg err: %s", err)
		return
	}
	// convert backgroundless image to mask
	maskOutput, err := ffmpeg.ConvertToMask(srvc.FileLocation{
		Path: rembgOutput.Path,
		Name: rembgOutput.Name,
	})
	if err != nil {
		log.Printf("maskoutput err: %s", err)
		return
	}
	// add white background to backgroundless image
	whitebg, err := ffmpeg.InsertWhiteBackground(srvc.FileLocation{
		Path: rembgOutput.Path,
		Name: rembgOutput.Name,
	})
	if err != nil {
		log.Printf("whitebg err: %s", err)
		return
	}
	// add white background to mask image
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
		// overlay.Path, overlay.Name,
	)
	c.Status(http.StatusOK)
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
