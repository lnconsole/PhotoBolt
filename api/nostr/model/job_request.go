package model

import (
	"log"

	"github.com/nbd-wtf/go-nostr"
)

type JobRequest struct {
	Itags            nostr.Tags
	PendingJobInputs map[string]*nostr.Subscription
	Event            nostr.Event
}

func (j *JobRequest) Wait(eventID string, sub *nostr.Subscription) {
	j.PendingJobInputs[eventID] = sub
}

func (j *JobRequest) AddInput(itag nostr.Tag) {
	j.Itags = j.Itags.AppendUnique(itag)
}

func (j *JobRequest) Receive(input nostr.Event) {
	if sub, pending := j.PendingJobInputs[input.ID]; pending {
		itag := j.Itags.GetFirst([]string{"i", input.ID})
		if itag == nil || len(*itag) < 3 {
			log.Printf("invalid i tag: %v", itag)
			return
		}
		(*itag)[1] = input.Content
		(*itag)[2] = "url"

		sub.Unsub()
		delete(j.PendingJobInputs, input.ID)
	}
}

func (j *JobRequest) Ready() bool {
	return len(j.PendingJobInputs) == 0
}
