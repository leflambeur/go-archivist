package main

import (
	"fmt"
	"github.com/leflambeur/go-archivist/archivist"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	app := &cli.App{

		Name:  "go-archivist",
		Usage: "Utility for Interacting with Jitsuin Archivist",

		Commands: []cli.Command{
			{
				Name:    "init",
				Aliases: []string{"i}"},
				Action: func(c *cli.context) error {
					archivistInit, err := archivist.Init()
					if err != nil {
						return err
					}
					fmt.Println(archivistInit)
					return nil
				},
			},
		},
	}
}
