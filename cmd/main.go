package main

import (
	"os"

	"github.com/vvatanabe/backlog-bulk-issue-registration-cli/bbir"
)

func main() {
	bbir.NewCLI().Run(os.Args)
}
