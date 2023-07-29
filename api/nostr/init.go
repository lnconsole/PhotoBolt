package istr

import (
	"log"

	nimage "github.com/lnconsole/photobolt/api/nostr/image"
	"github.com/lnconsole/photobolt/api/nostr/model"
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
		PendingJobRequestID_Consumer = map[string][]*model.JobRequest{}  // map[input_eid]. cache of nostr events we are waiting for. Job input is a pointer that gets updated till there are no more pending inputs
		Processor                    = map[int]func(*model.JobRequest){} // map[kind]fn(Jobrequest)error
		seen                         = map[string]bool{}
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
			if seen[evt.ID] {
				continue
			}
			seen[evt.ID] = true

			log.Printf("rcv: %v", evt)
			if evt.Kind == nimage.KindGeneration() ||
				evt.Kind == nimage.KindManipulation() {
				var (
					jr = &model.JobRequest{
						PendingJobInputs: map[string]*nostr.Subscription{},
						Event:            evt,
					}
					// pendingJobInputEids = []string{}
					itags = evt.Tags.GetAll([]string{"i"})
				)
				for _, itag := range itags {
					if len(itag) < 3 {
						continue
					}
					inputEventID := itag.Value()
					if itag[2] == "job" {
						newConsumers := []*model.JobRequest{jr}
						if consumers, ok := PendingJobRequestID_Consumer[inputEventID]; ok {
							newConsumers = append(newConsumers, consumers...)
						}
						PendingJobRequestID_Consumer[inputEventID] = newConsumers
						// pendingJobInputEids = append(pendingJobInputEids, inputEventID)
						sub := pstr.Subscribe(nostr.Filters{{
							Kinds: []int{kindJobResult},
							Tags: nostr.TagMap{
								"e": []string{inputEventID},
							},
						}})
						jr.Wait(inputEventID, sub)
						copy := nostr.Tag{}
						for idx := range itag {
							copy = append(copy, itag[idx])
						}
						jr.AddInput(copy)
					} else if itag[2] == "url" {
						base64, err := jr.DownloadUrl(itag[1])
						if err != nil {
							log.Printf("download url: %s", err.Error())
							continue
						}
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
					}
				}
				// for _, id := range pendingJobInputEids {
				// 	log.Printf("adding pending job request id: %s", id)
				// 	PendingJobRequest[id] = jr
				// }
			} else if evt.Kind == kindJobResult { // job result
				etag := evt.Tags.GetFirst([]string{"e"})
				eid := etag.Value()
				log.Printf("got 65001: %s", eid)
				if consumers, ok := PendingJobRequestID_Consumer[eid]; ok {
					log.Printf("%d consumers are waiting", len(consumers))
					for _, consumer := range consumers {
						consumer.Receive(eid, evt.Content)
						if consumer.Ready() {
							if process, ok := Processor[consumer.Event.Kind]; ok {
								process(consumer)
							}
						}
					}
					delete(PendingJobRequestID_Consumer, eid)
				}
			}
		}
	}()

	return nil
}
