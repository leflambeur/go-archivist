package main

import (
	"fmt"
	"log"
	"os"

	"github.com/leflambeur/go-archivist/archivist"
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
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "w",
						Usage: "Writes the token to .jitsuin in the current directory",
					},
				},
				Action: func(c *cli.Context) error {
					archivistInit, err := archivist.Init()
					if err != nil {
						return err
					}
					fmt.Println("\nToken:\n\nBearer " + archivistInit)
					if c.Bool("w") {
						f, err := os.Create(".jitsuin")
						if err != nil {
							return err
						}
						defer f.Close()
						tokenWrite, err := f.WriteString("Authorization: Bearer " + archivistInit)
						if err != nil {
							return err
						}
						fmt.Printf("\nWrote %d bytes to .jitsuin\n", tokenWrite)
					}
					return nil
				},
			},
			{
				Name:    "asset",
				Aliases: []string{"a"},
				Action: func(c *cli.Context) error {
					tokenExists, err := archivist.CheckTokenExists()
					if !tokenExists {
						return err
					}
					archivistAsset, err := archivist.AssetGetAll()
					if err != nil {
						return err
					}
					fmt.Printf("%s", archivistAsset)
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
