package main

import (
	"fmt"
	"os"
)

var commands = map[string]func() string{
	"uuid":     randUuid,
	"rand":     randChars,
	"email":    email,
	"video":    video,
	"upload":   upload,
	"frinkiac": frinkiac,
}

func main() {
	arg := os.Args[1]
	fmt.Println(commands[arg]())
}
