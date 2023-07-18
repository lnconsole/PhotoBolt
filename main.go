package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lnconsole/photobolt/api/background"
	"github.com/lnconsole/photobolt/api/icon"
	"github.com/lnconsole/photobolt/api/overlay"
	"github.com/lnconsole/photobolt/env"
)

func main() {
	log.Printf("winning")

	if err := env.Init(); err != nil {
		log.Printf("env err: %s", err)
		return
	}

	engine := gin.Default()
	setupRoutes(engine)
	engine.Run()
}

func setupRoutes(engine *gin.Engine) {
	api := engine.Group("/api")
	// BACKGROUND
	api.POST("/background", background.Replace(env.PhotoBolt.Automatic1111URL))
	// ICON
	api.POST("/icon", icon.Generate)
	// OVERLAY
	api.POST("/overlay", overlay.Combine)
}
