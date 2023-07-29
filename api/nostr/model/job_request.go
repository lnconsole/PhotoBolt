package model

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/lnconsole/photobolt/env"
	"github.com/lnconsole/photobolt/http"
	srvc "github.com/lnconsole/photobolt/service"
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

func (j *JobRequest) Receive(inputEventID string, result string) {
	log.Printf("receiving input: %s", inputEventID)
	if sub, pending := j.PendingJobInputs[inputEventID]; pending {
		itag := j.Itags.GetFirst([]string{"i", inputEventID})
		if itag == nil || len(*itag) < 3 {
			log.Printf("invalid i tag: %v", itag)
			return
		}
		base64, err := j.DownloadUrl(result)
		if err != nil {
			log.Printf("err downloading url: %s, %s", result, err.Error())
			return
		}
		(*itag)[1] = base64
		(*itag)[2] = "text"

		log.Printf("successfully received input: %s", inputEventID)

		sub.Unsub()
		delete(j.PendingJobInputs, inputEventID)
	}
}

func (j *JobRequest) Ready() bool {
	log.Printf("is ready: %v", j.PendingJobInputs)
	return len(j.PendingJobInputs) == 0
}

func (j *JobRequest) DownloadUrl(url string) (string, error) {
	downloaded := srvc.FileLocation{
		Path: fmt.Sprintf("%s/api/nostr", env.PhotoBolt.RepoDirectory),
		Name: uuid.New().String() + ".png",
	}
	if err := http.DownloadFile(url, downloaded); err != nil {
		return "", err
	}
	base64, err := downloaded.ToBase64()
	if err != nil {
		return "", err
	}
	downloaded.Remove()

	return base64, nil
}
