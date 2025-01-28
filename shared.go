package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

func CreateTempFile(input io.ReadCloser, name string) (string, error) {
	tempFile, err := os.CreateTemp(os.TempDir(), name)
	defer tempFile.Close()
	defer input.Close()
	if err != nil {
		return "", err
	}

	_, err = io.Copy(tempFile, input)
	if err != nil {
		return "", err
	}
	return tempFile.Name(), nil
}

func UploadFile(input io.ReadCloser) (string, error) {
	data, err := io.ReadAll(input)
	if err != nil {
		return "", err
	}
	result, err := sendPostRequest(shareUrl, content{
		fname: fmt.Sprintf("%s.mp4", uuid.New().String()),
		ftype: "video/mp4",
		fdata: data,
	})

	if err != nil {
		return "", err
	}

	return strings.ReplaceAll(strings.ReplaceAll(string(result), "share.unacc.eu", "from.bi.gp"), "\n", "/video.mp4"), nil
}

func DownloadCobaltFile(input string) (io.ReadCloser, error) {
	_, _ = http.NewRequest("GET", cobaltUrl, nil)
	body := []byte(fmt.Sprintf(`{"url": "%s"}`, input))
	req, err := http.NewRequest("POST", cobaltUrl, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	output := map[string]interface{}{}
	err = json.Unmarshal(bodyBytes, &output)
	if err != nil {
		return nil, err
	}

	streamUrl, ok := output["url"]
	if !ok {
		text, ok := output["text"]
		if ok {
			return nil, errors.New(text.(string))
		}
		cobaltError, ok := output["error"]
		if ok {
			return nil, errors.New(cobaltError.(string))
		}
		return nil, errors.New(output["error"].(string))
	}

	resp, err := http.Get(streamUrl.(string))
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

type content struct {
	fname string
	ftype string
	fdata []byte
}

func sendPostRequest(url string, file content) ([]byte, error) {
	var (
		buf = new(bytes.Buffer)
		w   = multipart.NewWriter(buf)
	)

	part, err := w.CreateFormFile("file", file.fname)
	if err != nil {
		return []byte{}, err
	}

	_, err = part.Write(file.fdata)
	if err != nil {
		return []byte{}, err
	}

	err = w.Close()
	if err != nil {
		return []byte{}, err
	}

	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return []byte{}, err
	}
	req.Header.Add("Content-Type", w.FormDataContentType())

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()

	cnt, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}
	return cnt, nil
}
