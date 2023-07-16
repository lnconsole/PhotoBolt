package overlay

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Combine(c *gin.Context) {
	/*
		overlay images into a image
		return image
	*/
	var payload CombineImages

	if err := c.ShouldBind(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": "ERROR",
			"reason": err.Error(),
		})
		return
	}

}
