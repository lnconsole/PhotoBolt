package overlay

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lnconsole/photobolt/env"
	srvc "github.com/lnconsole/photobolt/service"
	"github.com/lnconsole/photobolt/service/ffmpeg"
	"github.com/lnconsole/photobolt/shared"
)

func Combine() gin.HandlerFunc {
	/*
		overlay images into a image
		return image
	*/
	return func(c *gin.Context) {
		var payload CombineImagesBody

		if err := c.ShouldBind(&payload); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status": "ERROR",
				"reason": err.Error(),
			})
			return
		}

		var (
			// front
			front             = payload.Front
			frontExtension    = filepath.Ext(front.Filename)
			frontFileName     = uuid.New().String() + frontExtension
			frontFileLocation = srvc.FileLocation{
				Path: fmt.Sprintf("%s/%s", env.PhotoBolt.RepoDirectory, "api/overlay"),
				Name: frontFileName,
			}
			// back
			back             = payload.Back
			backExtension    = filepath.Ext(back.Filename)
			backFileName     = uuid.New().String() + backExtension
			backFileLocation = srvc.FileLocation{
				Path: fmt.Sprintf("%s/%s", env.PhotoBolt.RepoDirectory, "api/overlay"),
				Name: backFileName,
			}
		)
		if err := c.SaveUploadedFile(front, frontFileLocation.FullPath()); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": fmt.Sprintf("Unable to save the file: %v", err),
			})
			return
		}
		defer func() { frontFileLocation.Remove() }()

		if err := c.SaveUploadedFile(back, backFileLocation.FullPath()); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": fmt.Sprintf("Unable to save the file: %v", err),
			})
			return
		}
		defer func() { backFileLocation.Remove() }()

		// convert backgroundless image to mask
		overlayOutput, err := ffmpeg.OverlayImages(
			frontFileLocation,
			backFileLocation,
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

		split := strings.Split(overlayFileBase64, ",")

		c.JSON(http.StatusOK, &CombineImagesResponse{
			Image: split[1],
		})
	}
}
