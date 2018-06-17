package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		GenerateKeyCommand,
		CreateElectionCommand,
		AddPollCommand,
		VoteCommand,
		QueryElectionsCommand,
		QueryLatestElectionCommand,
		QueryPollsCommand,
		QueryLatestPollCommand,
		QueryResultsCommand,
	}
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
