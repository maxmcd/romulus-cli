package main

import (
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "romulus"
	app.Usage = "A tool to manage romulus sites"
	app.EnableBashCompletion = true
	app.Action = func(c *cli.Context) {
		println("Hello friend!")
	}
	app.Commands = []cli.Command{
		{
			Name:      "login",
			ShortName: "l",
			Usage:     "login to romulus",
			Action: func(c *cli.Context) {
				println("username: ", c.Args().First(), "password:", c.Args().Get(1))
			},
		},
	}

	app.Run(os.Args)
}
