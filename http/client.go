package http

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func Post(url string, body interface{}, response interface{}) error {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewBuffer(bodyBytes),
	)
	if err != nil {
		return err
	}

	client := &http.Client{}
	httpResponse, err := client.Do(request)
	if err != nil {
		return err
	}

	defer func() {
		if err := httpResponse.Body.Close(); err != nil {
			log.Print(err)
		}
	}()

	if err := json.NewDecoder(httpResponse.Body).Decode(response); err != nil {
		return err
	}

	return nil
}
