package main

import (
	"github.com/google/uuid"
	"math/rand"
	"os"
	"strconv"
)

var chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-."

func randChars() string {
	amount, _ := strconv.ParseInt(os.Getenv("ESPANSO_CHARS"), 10, 64)
	out := ""
	for i := 0; i < int(amount); i++ {
		out += string(chars[rand.Intn(len(chars))])
	}
	return out
}

func randUuid() string {
	id := uuid.New()
	return id.String()
}
