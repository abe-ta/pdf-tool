package main

import (
	"fmt"
	"os"

	"github.com/abe-ta/pdf-tool/command"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "pdf-tool"
	app.Usage = ""
	app.Description = "A PDF tool that allows you to edit PDF files on the command line."
	app.Version = "1.0.0"
	app.Commands = command.NewCommands()
	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		return
	}
}
