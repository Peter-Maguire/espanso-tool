package main

import (
	"encoding/hex"
	"golang.org/x/sys/windows"
	"hash/fnv"
	"math"
	"regexp"
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

func formatEmailName(text string) string {
	re := regexp.MustCompile("[^a-z]|signup|register")
	text = strings.ToLower(text)
	lastDash := strings.LastIndex(text, "-")
	if lastDash > -1 {
		text = text[:lastDash]
	}
	text = re.ReplaceAllString(text, "")
	return text[:int(math.Min(32, float64(len(text))))]
}

func emailHash(text string) string {
	h := fnv.New32a()
	h.Write([]byte(emailSecret))
	h.Write([]byte(text))
	return hex.EncodeToString([]byte{byte(h.Sum32())})
}

func email() string {
	if hwnd := getWindow("GetForegroundWindow"); hwnd != 0 {
		text := GetWindowText(HWND(hwnd))
		return formatEmailName(text) + "-" + emailHash(text)
	}
	return "Couldn't get foreground window"
}
