package main

import (
	"os"

	_ "github.com/kbsink-org/douyin-plugin/internal/douyinregister"
	"github.com/kbsink-org/kbsink/cmd/kb-sink-md/cli"
)

func main() {
	os.Exit(cli.Run(os.Args))
}
