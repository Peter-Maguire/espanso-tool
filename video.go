package main

import (
	"fmt"
	"os"
	"os/exec"
)

func video() string {
	url := os.Getenv("ESPANSO_CLIPBOARD")
	data, err := DownloadCobaltFile(url)
	if err != nil {
		return err.Error()
	}

	fileName, err := CreateTempFile(data, "Cobalt*.mp4")
	if err != nil {
		return err.Error()
	}
	err = exec.Command("powershell", "-NoProfile", fmt.Sprintf("Set-Clipboard -Path %s", fileName)).Run()
	if err != nil {
		return err.Error()
	}
	return ""
}

func upload() string {
	url := os.Getenv("ESPANSO_CLIPBOARD")
	data, err := DownloadCobaltFile(url)
	if err != nil {
		return err.Error()
	}

	out, err := UploadFile(data)
	if err != nil {
		return err.Error()
	}
	return out
}
