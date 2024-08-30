package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
)

func Frinkiac(query string) string {
	result, err := SearchFrinkiac(query)
	if err != nil {
		return err.Error()
	}

	caption, err := GetCaption(result.Episode, result.Timestamp)
	if err != nil {
		return err.Error()
	}

	fileName, err := CreateTempFile(DownloadFile(GetGIFUrl(*result, caption)), "Frinkiac*.gif")
	if err != nil {
		return err.Error()
	}
	err = exec.Command("powershell", "-NoProfile", fmt.Sprintf("Set-Clipboard -Path %s", fileName)).Run()
	if err != nil {
		return err.Error()
	}
	return ""
}

type FrinkiacSearch struct {
	Id        int
	Episode   string
	Timestamp int
}

func SearchFrinkiac(term string) (*FrinkiacSearch, error) {

	resp, err := http.Get(fmt.Sprintf("https://www.frinkiac.com/api/search?q=%s", url.QueryEscape(term)))

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var searchResults []FrinkiacSearch
	err = json.Unmarshal(body, &searchResults)
	if err != nil {
		return nil, err
	}

	if len(searchResults) < 1 {
		return nil, errors.New("no results")
	}

	firstResult := searchResults[0]

	return &firstResult, nil
}

type FrinkiacCaptions struct {
	Subtitles []FrinkiacSub
}

type FrinkiacSub struct {
	Content string
}

func GetCaption(episode string, timestamp int) (string, error) {
	resp, err := http.Get(fmt.Sprintf("https://www.frinkiac.com/api/caption?e=%s&t=%d", episode, timestamp))

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var captions FrinkiacCaptions
	err = json.Unmarshal(body, &captions)
	if err != nil {
		return "", err
	}

	subsList := make([]string, len(captions.Subtitles))

	for i, sub := range captions.Subtitles {
		subsList[i] = sub.Content
	}
	return strings.Join(subsList, "\n"), nil
}

func GetGIFUrl(result FrinkiacSearch, caption string) string {
	return fmt.Sprintf("https://frinkiac.com/gif/%s/%d/%d.gif?b64lines=%s", result.Episode, result.Timestamp-2000, result.Timestamp+2000, base64.StdEncoding.EncodeToString([]byte(caption)))
}

func DownloadFile(url string) io.ReadCloser {
	resp, err := http.Get(url)
	if err != nil {
		return nil
	}

	return resp.Body
}
