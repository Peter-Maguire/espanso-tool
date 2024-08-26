package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/sys/windows"
	"io"
	"math"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

var (
	mod                     = windows.NewLazyDLL("user32.dll")
	procGetWindowText       = mod.NewProc("GetWindowTextW")
	procGetWindowTextLength = mod.NewProc("GetWindowTextLengthW")
)

type (
	HANDLE uintptr
	HWND   HANDLE
)

func GetWindowTextLength(hwnd HWND) int {
	ret, _, _ := procGetWindowTextLength.Call(
		uintptr(hwnd))

	return int(ret)
}

func GetWindowText(hwnd HWND) string {
	textLen := GetWindowTextLength(hwnd) + 1

	buf := make([]uint16, textLen)
	procGetWindowText.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(textLen))

	return syscall.UTF16ToString(buf)
}

func getWindow(funcName string) uintptr {
	proc := mod.NewProc(funcName)
	hwnd, _, _ := proc.Call()
	return hwnd
}

var chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-."

func main() {
	arg := os.Args[1]
	if arg == "uuid" {
		id := uuid.New()
		fmt.Println(id.String())
		return
	}

	if arg == "rand" {
		amount, _ := strconv.ParseInt(os.Getenv("ESPANSO_CHARS"), 10, 64)
		out := ""
		for i := 0; i < int(amount); i++ {
			out += string(chars[rand.Intn(len(chars))])
		}
		fmt.Println(out)
		return
	}

	if arg == "email" {
		re := regexp.MustCompile("[^a-zA-Z]")
		if hwnd := getWindow("GetForegroundWindow"); hwnd != 0 {
			text := GetWindowText(HWND(hwnd))
			lastDash := strings.LastIndex(text, "-")
			if lastDash > -1 {
				text = text[:lastDash]
			}
			text = re.ReplaceAllString(text, "")
			fmt.Println(text[:int(math.Min(32, float64(len(text))))])
		}
		return
	}

	if arg == "video" {
		url := os.Getenv("ESPANSO_CLIPBOARD")
		data, err := DownloadCobaltFile(url)
		if err != nil {
			fmt.Println(err)
			return
		}

		fileName, err := CreateTempFile(data)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = exec.Command("powershell", "-NoProfile", fmt.Sprintf("Set-Clipboard -Path %s", fileName)).Run()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("")
		return
	}

	if arg == "upload" {
		url := os.Getenv("ESPANSO_CLIPBOARD")
		data, err := DownloadCobaltFile(url)
		if err != nil {
			fmt.Println(err)
			return
		}

		out, err := UploadFile(data)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(out)
		return
	}
}

func CreateTempFile(input io.ReadCloser) (string, error) {
	tempFile, err := os.CreateTemp(os.TempDir(), "Cobalt*.mp4")
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
