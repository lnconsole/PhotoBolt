package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

	if response != nil {
		defer func() {
			if err := httpResponse.Body.Close(); err != nil {
				log.Print(err)
			}
		}()

		if err := json.NewDecoder(httpResponse.Body).Decode(response); err != nil {
			return err
		}
	}

	return nil
}

func PostBytes(url string, b bytes.Buffer, contentType string, response interface{}) error {
	request, err := http.NewRequest(
		http.MethodPost,
		url,
		&b,
	)
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", contentType)

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

	if httpResponse.StatusCode != 200 {
		bodyBytes, err := io.ReadAll(httpResponse.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("status code: %d; err: %s", httpResponse.StatusCode, string(bodyBytes))
	}

	if err := json.NewDecoder(httpResponse.Body).Decode(response); err != nil {
		return err
	}

	return nil
}

func Get(url string, response interface{}) error {
	request, err := http.NewRequest(
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return err
	}

	client := &http.Client{}
	httpResponse, err := client.Do(request)
	if err != nil {
		return err
	}

	if httpResponse.StatusCode < 200 || httpResponse.StatusCode > 299 {
		return fmt.Errorf("http response code: %d", httpResponse.StatusCode)
	}

	if response != nil {
		if err := json.NewDecoder(httpResponse.Body).Decode(response); err != nil {
			return err
		}

		defer func() {
			if err := httpResponse.Body.Close(); err != nil {
				log.Print(err)
			}
		}()
	}

	return nil
}
