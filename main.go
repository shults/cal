package main

import (
	"fmt"
	"os"

	"github.com/shults/cal/calendar"
)

func main() {
	if out, err := calendar.Run(os.Args); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %s", err)
		os.Exit(2)
	} else {
		print(string(out))
	}
}
