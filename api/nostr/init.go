package istr

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	nimage "github.com/lnconsole/photobolt/api/nostr/image"
	"github.com/lnconsole/photobolt/api/nostr/model"
	"github.com/lnconsole/photobolt/env"
	"github.com/lnconsole/photobolt/http"
	srvc "github.com/lnconsole/photobolt/service"
	pstr "github.com/lnconsole/photobolt/service/nostr"
	"github.com/nbd-wtf/go-nostr"
)

/*
pass in (ch<- for nostr subscription, <-ch for event)
for sub in ch
	subscribe
	for evt in sub
		send to event ch
check param
- if background & clear, or
- if param & overlay
  - if input is text, validate payload is valid
  - if input is job, subscribe to [kind 65001, etag "job request"]. subscription ends when the event is received, and payload is valid
perform job


PendingJobInput map[input_eid]*JobInput
Processor map[kind]func(JobInput)err

init()
	main = make(chan, nostr.event)
	go func()
		for e in unique(main)
			( any error, log & continue)
			if target 65002-66000
				ji = new JobInput(target)
				pendingJobInputEventIds = []string
				for all i
					if job,
						pendingJobinputEventIds.append(i.value())
						sub = subscribe(65001(ref target.id))
						ji.Wait(target.id, sub)
					else
						ji.Add(target.itag)
				if ji.Ready()
					if processor := Processor[target.kind]
						processor(ji)
					else
						error
				else
					for id := pendingJobInputEventIds
						PendingJobInput[id] = jobinput
			if target 65001
				if ji := PendingJobInput[id]
					delete(PendingJobInput, id)
					ji.Receive(target)
					if ji.Ready()
						if processor := Processor[target.kind]
							processor(ji)
						else
							error

Subscribe(filter) sub
	s = nostr.sub
	go func() main <- s.Events
	return s

*/

var (
	// relay         *nostr.Relay
	// relayUrl      = "ws://localhost:7447"
	// mainch        = make(chan nostr.Event)
	kindJobResult = 65001
	// pk            string
)

func Init() error {
	var (
		PendingJobRequest = map[string]*model.JobRequest{}    // map[input_eid]. cache of nostr events we are waiting for. Job input is a pointer that gets updated till there are no more pending inputs
		Processor         = map[int]func(*model.JobRequest){} // map[kind]fn(Jobrequest)error
		// err               error
	)

	Processor[nimage.KindGeneration()] = nimage.ProcessGeneration
	Processor[nimage.KindManipulation()] = nimage.ProcessManipulation

	filters := nostr.Filters{
		nimage.FilterManipulation(),
		nimage.FilterGeneration(),
	}

	pstr.Subscribe(filters)

	go func() {
		for evt := range pstr.Mainch {
			log.Printf("rcv: %v", evt)
			if evt.Kind == nimage.KindGeneration() ||
				evt.Kind == nimage.KindManipulation() {
				var (
					jr                  = &model.JobRequest{Event: evt}
					pendingJobInputEids = []string{}
					itags               = evt.Tags.GetAll([]string{"i"})
				)
				for _, itag := range itags {
					if len(itag) < 3 {
						continue
					}
					inputEventID := itag.Value()
					if itag[2] == "job" {
						pendingJobInputEids = append(pendingJobInputEids, itag.Value())
						sub := pstr.Subscribe(nostr.Filters{{
							Kinds: []int{kindJobResult},
							IDs:   []string{inputEventID}, // eid
						}})
						jr.Wait(inputEventID, sub)
						jr.AddInput(itag)
					} else if itag[2] == "url" {
						downloaded := srvc.FileLocation{
							Path: fmt.Sprintf("%s/api/nostr", env.PhotoBolt.RepoDirectory),
							Name: uuid.New().String() + ".png",
						}
						if err := http.DownloadFile(itag[1], downloaded); err != nil {
							log.Printf("download file: %s", err.Error())
							continue
						}
						base64, err := downloaded.ToBase64()
						if err != nil {
							log.Printf("tobase64: %s", err.Error())
							continue
						}
						downloaded.Remove()
						copy := nostr.Tag{}
						for idx := range itag {
							if idx == 1 {
								copy = append(copy, base64)
							} else if idx == 2 {
								copy = append(copy, "text")
							} else {
								copy = append(copy, itag[idx])
							}
						}
						jr.AddInput(copy)
					} else {
						jr.AddInput(itag)
					}
					if jr.Ready() {
						if process, ok := Processor[evt.Kind]; ok {
							process(jr)
						}
					} else {
						for _, id := range pendingJobInputEids {
							PendingJobRequest[id] = jr
						}
					}
				}
			} else if evt.Kind == kindJobResult { // job result
				if jr, ok := PendingJobRequest[evt.ID]; ok {
					delete(PendingJobRequest, evt.ID)
					jr.Receive(evt)
					if jr.Ready() {
						if process, ok := Processor[evt.Kind]; ok {
							process(jr)
						}
					}
				}
			}
		}
	}()

	return nil
}
