package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
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

func PostForm(url string, body map[string]string, response interface{}) error {
	var (
		b   bytes.Buffer
		err error
		w   = multipart.NewWriter(&b)
	)

	for key, value := range body {
		if key == "file" || key == "front" || key == "back" {
			var (
				fw   io.Writer
				file = mustOpen(value)
			)
			if fw, err = w.CreateFormFile(key, file.Name()); err != nil {
				return fmt.Errorf("error creating writer: %v", err)
			}
			if _, err = io.Copy(fw, file); err != nil {
				return fmt.Errorf("error with io.Copy: %v", err)
			}
		} else {
			promptField, err := w.CreateFormField(key)
			if err != nil {
				return fmt.Errorf("error creating writer: %v", err)
			}
			promptField.Write([]byte(value))
		}
	}
	w.Close()

	request, err := http.NewRequest(
		http.MethodPost,
		url,
		&b,
	)
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", w.FormDataContentType())

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

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		pwd, _ := os.Getwd()
		fmt.Println("PWD: ", pwd)
		panic(err)
	}
	return r
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
