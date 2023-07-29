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
	"github.com/lnconsole/photobolt/service/ffmpeg"
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
s <- nostr.sub{ 65007 }

for e in e
	( any error, log & continue )
	if 65007
		check i
			if text, process()
			if job, s <- nostr.sub{ 65001(ref 65007)}
	if 65001
		process()

process(ji)
	go routine
	if ji.Jobrequest.param, (background & clear) || (overlay)
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

func KindManipulation() int {
	return 65007
}

func FilterManipulation() nostr.Filter {
	now := time.Now()
	return nostr.Filter{
		Kinds: []int{KindManipulation()},
		Since: &now,
	}
}

func ProcessManipulation(jr *model.JobRequest) {
	go func() {
		paramTag := jr.Event.Tags.GetFirst([]string{"param"})
		jobType := paramTag.Value()
		if jobType != "background" &&
			jobType != "overlay" {
			log.Printf("unknown param for kind %d: %s", jr.Event.Kind, paramTag.Value())
			return
		}
		if (jobType == "background" && len(*paramTag) < 3) ||
			(jobType == "overlay" && len(*paramTag) < 3) {
			log.Printf("incomplete param: %v", paramTag)
			return
		}
		// Generate Invoice
		chargeMsat := 10000
		lnData, err := ln.CreateInvoice(ln.InvoiceParams{
			AmountMsat: chargeMsat,
			Memo:       fmt.Sprintf("'%s' Service Fee: %d sats", jobType, chargeMsat/1000),
		})
		if err != nil {
			log.Printf("create inv: %s", err)
			return
		}
		time.Sleep(time.Duration(shared.RandInt(3)) * time.Second)
		// send 65000
		if _, err := pstr.Publish(context.Background(), nostr.Event{
			CreatedAt: time.Now(),
			Kind:      pstr.KindJobFeedback,
			Content:   "",
			Tags: nostr.Tags{
				{"status", "payment-required", fmt.Sprintf("I would like to process this job for %d sats", chargeMsat/1000)},
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
				filename := uuid.New().String() + ".png"
				if jobType == "background" {
					// get input base64 data
					inputText := jr.Itags.GetFirst([]string{"i"})
					// save png to disk
					inputLoc := srvc.FileLocation{
						Path: fmt.Sprintf("%s/api/nostr/image", env.PhotoBolt.RepoDirectory),
						Name: filename,
					}
					if err := inputLoc.SavePNG(inputText.Value()); err != nil {
						log.Printf("save png: %s", err)
						return
					}
					// remove bg
					outputLoc, err := rembg.RemoveBackground(inputLoc)
					if err != nil {
						log.Printf("os.Create: %s", err)
						return
					}
					// disk to base64
					outputBase64, err := outputLoc.ToBase64()
					if err != nil {
						log.Printf("to base64: %s", err)
						return
					}
					inputLoc.Remove()
					outputLoc.Remove()

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
				} else if jobType == "overlay" {
					inputTexts := jr.Itags.GetAll([]string{"i"})
					if len(inputTexts) != 2 {
						log.Printf("should only be 2 input texts: %v", inputTexts)
						return
					}
					var (
						frontBase64       string
						frontFileName     = uuid.New().String() + ".png"
						frontFileLocation = srvc.FileLocation{
							Path: fmt.Sprintf("%s/%s", env.PhotoBolt.RepoDirectory, "api/nostr/image"),
							Name: frontFileName,
						}
						backBase64       string
						backFileName     = uuid.New().String() + ".png"
						backFileLocation = srvc.FileLocation{
							Path: fmt.Sprintf("%s/%s", env.PhotoBolt.RepoDirectory, "api/nostr/image"),
							Name: backFileName,
						}
					)
					for _, inputText := range inputTexts {
						if len(inputText) < 4 {
							log.Printf("len should be at least 4: %v", inputText)
							return
						}
						if inputText[3] == "front" {
							frontBase64 = inputText.Value()
						} else if inputText[3] == "back" {
							backBase64 = inputText.Value()
						}
					}
					if err := frontFileLocation.SavePNG(frontBase64); err != nil {
						log.Printf("save png: %s", frontFileLocation.FullPath())
						return
					}
					if err := backFileLocation.SavePNG(backBase64); err != nil {
						log.Printf("save png: %s", frontFileLocation.FullPath())
						return
					}
					var (
						width  int
						height int
					)
					if (*paramTag)[2] == "full" {
						backImg, err := shared.GetImage(backBase64)
						if err != nil {
							log.Printf("error getimage: %v", err)
							return
						}
						width = backImg.Bounds().Dx()
						height = backImg.Bounds().Dy()
					} else { // logo
						width = 256
						height = 256
					}
					// convert backgroundless image to mask
					overlayOutput, err := ffmpeg.OverlayImages(
						frontFileLocation,
						backFileLocation,
						width, height,
					)
					if err != nil {
						log.Printf("overlay images: %s", err)
						return
					}
					frontFileLocation.Remove()
					backFileLocation.Remove()
					// disk to base64
					outputBase64, err := overlayOutput.ToBase64()
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
				}
				break
			} else if update.InvoiceStatus == invoices.ContractCanceled.String() {
				break
			}
			log.Printf("pending: %s", update.InvoiceStatus)
		}
	}()
}
