package main

import (
	"github.com/chen-keinan/beacon/internal/cli"
)

func main() {
	cli.InitCLI(cli.ArgsSanitizer)
}
