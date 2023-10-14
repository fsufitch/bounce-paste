package main

import (
	"fmt"
	"os"

	"github.com/fsufitch/bounce-paste/db"

	"github.com/urfave/cli/v2"
)

var commandLineApp = cli.App{
	Action: mainAction,
}

func mainAction(ctx *cli.Context) error {
	fmt.Println("hello")
	fmt.Println(db.GetSomeText())
	return nil
}

func main() {
	commandLineApp.Run(os.Args)
}
