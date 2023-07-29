package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lnconsole/photobolt/api/http/background"
	"github.com/lnconsole/photobolt/api/http/icon"
	"github.com/lnconsole/photobolt/api/http/overlay"
	istr "github.com/lnconsole/photobolt/api/nostr"
	"github.com/lnconsole/photobolt/env"
	"github.com/lnconsole/photobolt/service/ln"
	pstr "github.com/lnconsole/photobolt/service/nostr"
)

func main() {
	log.Printf("winning")

	if err := env.Init("env/.env"); err != nil {
		log.Printf("env err: %s", err)
		return
	}
	// lightning init
	if err := ln.Init(
		env.PhotoBolt.LNDMacaroonHex,
		env.PhotoBolt.LNDCertPath,
		env.PhotoBolt.LNDGrpcAddr,
		env.PhotoBolt.LndClientNetwork(),
		env.PhotoBolt.LnNetwork(),
	); err != nil {
		log.Printf("lnd init: %s", err)
		return
	}
	// nostr init
	if err := pstr.Init(
		env.PhotoBolt.NostrRelay,
		env.PhotoBolt.NostrPrivateKey,
	); err != nil {
		log.Printf("nostr init: %s", err)
		return
	}
	// nostr service provider init
	if err := istr.Init(); err != nil {
		log.Printf("nostr service provider init: %s", err)
		return
	}

	engine := gin.Default()
	setupRoutes(engine)
	engine.Run(":" + env.PhotoBolt.ServerPort)
}

func setupRoutes(engine *gin.Engine) {
	api := engine.Group("/api")
	// BACKGROUND
	api.POST("/background", background.Replace(env.PhotoBolt.Automatic1111URL))
	// ICON
	api.POST("/icon", icon.Generate(env.PhotoBolt.Automatic1111URL))
	// OVERLAY
	api.POST("/overlay", overlay.Combine())
}
