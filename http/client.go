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

	"github.com/lnconsole/photobolt/env"
	srvc "github.com/lnconsole/photobolt/service"
	"github.com/lnconsole/photobolt/shared"
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

// add auth token
func PostForm(url string, authToken string, body map[string]string, response interface{}) (string, error) {
	var (
		b   bytes.Buffer
		err error
		w   = multipart.NewWriter(&b)
	)

	for key, value := range body {
		if key == "file" || key == "front" || key == "back" {
			var (
				fw   io.Writer
				file = openFile(value)
			)
			if file == nil {
				continue
			}

			if fw, err = w.CreateFormFile(key, file.Name()); err != nil {
				return "", fmt.Errorf("error creating writer: %v", err)
			}
			if _, err = io.Copy(fw, file); err != nil {
				return "", fmt.Errorf("error with io.Copy: %v", err)
			}
		} else {
			promptField, err := w.CreateFormField(key)
			if err != nil {
				return "", fmt.Errorf("error creating writer: %v", err)
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
		return "", err
	}
	request.Header.Set("Content-Type", w.FormDataContentType())
	if authToken != "" {
		// if auth token provided, add Authorization: LSAT bla bla
		request.Header.Set("Authorization", authToken)
	}

	client := &http.Client{}
	httpResponse, err := client.Do(request)
	if err != nil {
		return "", err
	}

	defer func() {
		if err := httpResponse.Body.Close(); err != nil {
			log.Print(err)
		}
	}()

	if httpResponse.StatusCode == 402 {
		// return www-authorizaiton value
		l402 := httpResponse.Header.Get("Www-Authenticate")
		return l402, nil
	} else if httpResponse.StatusCode != 200 {
		bodyBytes, err := io.ReadAll(httpResponse.Body)
		if err != nil {
			return "", err
		}
		return "", fmt.Errorf("status code: %d; err: %s", httpResponse.StatusCode, string(bodyBytes))
	}

	if err := json.NewDecoder(httpResponse.Body).Decode(response); err != nil {
		return "", err
	}

	return "", nil
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

func DownloadFile(url string, loc srvc.FileLocation) error {
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
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode < 200 || httpResponse.StatusCode > 299 {
		return fmt.Errorf("http response code: %d", httpResponse.StatusCode)
	}

	// Create a empty file
	file, err := os.Create(loc.FullPath())
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the bytes to the fiel
	_, err = io.Copy(file, httpResponse.Body)
	if err != nil {
		return err
	}

	return nil
}

func UploadFileImgBB(fileBase64 string) (string, error) {
	form := new(bytes.Buffer)
	writer := multipart.NewWriter(form)
	formField, err := writer.CreateFormField("image")
	if err != nil {
		return "", err
	}
	_, err = formField.Write([]byte(fileBase64))
	if err != nil {
		return "", err
	}
	writer.Close()
	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.imgbb.com/1/upload?expiration=600&key=%s", env.PhotoBolt.ImgbbSecret), form)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	res := result{}
	shared.CleanUnmarshal(bodyText, &res)

	log.Printf("file url: %s", res.D.Url)

	return res.D.Url, nil
}

type Data struct {
	Url string `json:"url"`
}

type result struct {
	D Data `json:"data"`
}

func openFile(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		// pwd, _ := os.Getwd()
		// fmt.Println("PWD: ", pwd)
		return nil
	}
	return r
}
