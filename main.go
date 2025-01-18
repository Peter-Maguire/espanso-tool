package main

import (
	"fmt"
	"os"
)

var commands = map[string]func() string{
	"uuid":        randUuid,
	"rand":        randChars,
	"email":       email,
	"video":       video,
	"upload":      upload,
	"frinkiac":    frinkiac,
	"sort":        sort,
	"upper":       capitalise,
	"lower":       uncaptialise,
	"zalgo":       zalgo,
	"fullwidth":   transform("fullwidth"),
	"bold":        transform("bold"),
	"italic":      transform("italic"),
	"pirate":      transform("pirate"),
	"fancy":       transform("fancy"),
	"smallcaps":   transform("smallcaps"),
	"subscript":   transform("subscript"),
	"superscript": transform("superscript"),
}

func main() {
	arg := os.Args[1]
	fun, ok := commands[arg]
	if !ok {
		fmt.Printf("Unknown command: %s\n", arg)
	} else {
		fmt.Println(fun())
	}
}
