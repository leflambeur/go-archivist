package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := &cli.App{

		Name:  "go-archivist",
		Usage: "Utility for Interacting with Jitsuin Archivist",

		Commands: []cli.Command{
			{
				Name:    "init",
				Aliases: []string{"i"},
				Action: func(c *cli.Context) error {
					//archivistInit, err := archivist.Init()
					//if err != nil {
					//	return err
					//}
					//fmt.Println(archivistInit)
					return nil
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
