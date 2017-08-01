package main

import (
	"fmt"
	"github.com/simonshyu/circle/circle/agent"
	"github.com/simonshyu/circle/circle/server"
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "circle"
	app.Usage = "Simple CI tool !"

	app.Commands = []cli.Command{
		server.Command,
		agent.Command,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
