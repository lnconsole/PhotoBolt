package nimage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/lightningnetwork/lnd/invoices"
	"github.com/lnconsole/photobolt/api/nostr/model"
	"github.com/lnconsole/photobolt/env"
	"github.com/lnconsole/photobolt/http"
	srvc "github.com/lnconsole/photobolt/service"
	"github.com/lnconsole/photobolt/service/automatic1111"
	"github.com/lnconsole/photobolt/service/ln"
	pstr "github.com/lnconsole/photobolt/service/nostr"
	"github.com/lnconsole/photobolt/service/rembg"
	"github.com/lnconsole/photobolt/shared"
	"github.com/nbd-wtf/go-nostr"
)

/*
s = make(chan, nostr.sub)
e = make(chan, nostr.event)
model.sub(s, e)
s <- nostr.sub{ 65005 }

for e in e
	( any error, log & continue )
	if 65005
		check i
			if text, process()
			if job, s <- nostr.sub{ 65001(ref 65007)}
	if 65001
		process()

process(input)
	go routine
	if param, (background & clear) || (overlay)
		pch = generate invoice
		send 65000
		for status in pch:
			if success
				send 65000 processing
				if param, (bg & clear)
					loc = save file
					rembg.removebackground(file loc)
					send 65001
				elif param, (overlay)
					parse front & back
					ffmpeg.OverlayImages
					send 65001
				break
			elif failed, break
			log
*/

func KindGeneration() int {
	return 65005
}

func FilterGeneration() nostr.Filter {
	now := time.Now()
	return nostr.Filter{
		Kinds: []int{KindGeneration()},
		Since: &now,
	}
}

func ProcessGeneration(jr *model.JobRequest) {
	go func() {
		paramTag := jr.Event.Tags.GetFirst([]string{"param"})
		itag := jr.Itags.GetFirst([]string{"i"})
		// Generate Invoice
		chargeMsat := 10000
		lnData, err := ln.CreateInvoice(ln.InvoiceParams{
			AmountMsat: chargeMsat,
			Memo:       fmt.Sprintf("'Generative AI' Service Fee: %d sats", chargeMsat/1000),
		})
		if err != nil {
			log.Printf("create inv: %s", err)
			return
		}
		// send 65000
		if _, err := pstr.Publish(context.Background(), nostr.Event{
			CreatedAt: time.Now(),
			Kind:      pstr.KindJobFeedback,
			Content:   "",
			Tags: nostr.Tags{
				{"status", "payment-required", "I would like to process this job for you! Please pay up 🙂"},
				{"amount", fmt.Sprintf("%d", chargeMsat), lnData.Bolt11},
				{"e", jr.Event.ID, pstr.RelayUrl},
				{"p", jr.Event.PubKey},
			},
		}); err != nil {
			log.Printf("publish: %s", err)
			return
		}
		for update := range lnData.Update {
			if update.InvoiceStatus == invoices.ContractSettled.String() {
				if _, err := pstr.Publish(context.Background(), nostr.Event{
					CreatedAt: time.Now(),
					Kind:      pstr.KindJobFeedback,
					Content:   "",
					Tags: nostr.Tags{
						{"status", "processing", "Processing! It'll be ready in a moment"},
						{"e", jr.Event.ID, pstr.RelayUrl},
						{"p", jr.Event.PubKey},
					},
				}); err != nil {
					log.Printf("publish: %s", err)
					return
				}
				if paramTag == nil {
					// img2img
					img2imgInput := automatic1111.NewImg2ImgInput()
					img2imgInput.SDModelCheckpoint = automatic1111.SDModelDreamShaperV7
					img2imgInput.Prompt = automatic1111.LoraColoredIcons(0.9, jr.Event.Content)
					img2imgInput.SamplerName = automatic1111.SamplerDPMPP2MKarras
					img2imgInput.InitImages = []string{itag.Value()}

					outputLoc := srvc.FileLocation{
						Path: fmt.Sprintf("%s/%s", env.PhotoBolt.RepoDirectory, "api/nostr/image"),
						Name: uuid.New().String() + ".png",
					}
					imageOutput, err := automatic1111.Img2Img(
						env.PhotoBolt.Automatic1111URL,
						img2imgInput,
					)
					if err != nil {
						log.Printf("automatic img2img: %s", err)
						return
					}
					if err := outputLoc.SavePNG(imageOutput.Images[0]); err != nil {
						log.Printf("save png: %s", outputLoc.FullPath())
						return
					}
					rembgOutput, err := rembg.RemoveBackground(outputLoc)
					if err != nil {
						log.Printf("rembg err: %s", err)
						return
					}
					outputLoc.Remove()
					outputBase64, err := rembgOutput.ToBase64()
					if err != nil {
						log.Printf("to base64: %s", err)
						return
					}

					fileUrl, err := http.UploadFileImgBB(outputBase64)
					if err != nil {
						log.Printf("upload file: %s", err)
						return
					}

					if _, err := pstr.Publish(context.Background(), nostr.Event{
						CreatedAt: time.Now(),
						Kind:      pstr.KindJobResult,
						Content:   fileUrl,
						Tags: nostr.Tags{
							{"request", string(shared.CleanMarshal(jr.Event))},
							{"e", jr.Event.ID, pstr.RelayUrl},
							{"p", jr.Event.PubKey},
						},
					}); err != nil {
						log.Printf("publish: %s", err)
						return
					}
				} else {
					// txt2img
					if len(*paramTag) < 3 {
						log.Printf("incomplete param: %v", paramTag)
						return
					}
					if (*paramTag)[1] != "control-net" ||
						(*paramTag)[2] != "canny" {
						log.Printf("unknown param: %v", paramTag)
						return
					}
					inputImg, err := shared.GetImage(itag.Value())
					if err != nil {
						log.Printf("error getimage: %v", err)
						return
					}
					// all txt2img input preparation
					txt2img := automatic1111.NewText2ImgControlNetInput()
					txt2img.SDModelCheckpoint = automatic1111.SDModelPhotonV1
					txt2img.Prompt = jr.Event.Content
					txt2img.NegativePrompt = automatic1111.SDModelPhotonV1.NegativePrompt()
					txt2img.BatchSize = 1
					txt2img.Steps = 25
					txt2img.Seed = -1
					txt2img.CFGScale = 3
					txt2img.SamplerName = automatic1111.SamplerDPMPP2M
					txt2img.Width = inputImg.Bounds().Dx()
					txt2img.Height = inputImg.Bounds().Dy()

					cannyCNUnit := automatic1111.NewControlNetUnit()
					cannyCNUnit.InputImage = itag.Value()
					cannyCNUnit.Weight = 2
					cannyCNUnit.ControlMode = automatic1111.ControlNetModeBalanced
					cannyCNUnit.ProcessorRes = inputImg.Bounds().Dx()
					cannyCNUnit.ThresholdA = 100
					cannyCNUnit.ThresholdB = 200
					txt2img.AddControlNetUnit(cannyCNUnit)

					imageOutput, err := automatic1111.Text2ImgControlNet(
						env.PhotoBolt.Automatic1111URL,
						txt2img,
					)
					if err != nil {
						log.Printf("automatic.txt2img: %s", err)
						return
					}

					fileUrl, err := http.UploadFileImgBB(imageOutput.Images[0])
					if err != nil {
						log.Printf("upload file: %s", err)
						return
					}

					if _, err := pstr.Publish(context.Background(), nostr.Event{
						CreatedAt: time.Now(),
						Kind:      pstr.KindJobResult,
						Content:   fileUrl,
						Tags: nostr.Tags{
							{"request", string(shared.CleanMarshal(jr.Event))},
							{"e", jr.Event.ID, pstr.RelayUrl},
							{"p", jr.Event.PubKey},
						},
					}); err != nil {
						log.Printf("publish: %s", err)
						return
					}
				}
			} else if update.InvoiceStatus == invoices.ContractCanceled.String() {
				break
			}
		}
	}()
}
