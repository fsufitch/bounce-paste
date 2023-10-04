package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

var commandLineApp = cli.App{
	Action: mainAction,
}

func mainAction(ctx *cli.Context) error {
	fmt.Println("hello")
	return nil
}

func main() {
	commandLineApp.Run(os.Args)
}
