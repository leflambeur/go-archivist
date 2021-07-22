package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/leflambeur/go-archivist/archivist"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "go-rkvst",
		Usage: "Utility for Interacting with Jitsuin RKVST",

		Commands: []*cli.Command{
			{
				Name:    "init",
				Aliases: []string{"i"},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "interactive",
						Aliases: []string{"i"},
						Usage:   "Interactive",
					},
					&cli.BoolFlag{
						Name:    "write-file",
						Aliases: []string{"w"},
						Usage:   "Writes the token to .auth_token in the current directory",
					},
					&cli.BoolFlag{
						Name:    "use-client-secret",
						Aliases: []string{"s"},
						Usage:   "Chooses Client Secret Authentication",
					},
					&cli.BoolFlag{
						Name:    "base64",
						Aliases: []string{"b64"},
						Usage:   "Base64's the token",
					},
					&cli.StringFlag{
						Name:    "read-file",
						Aliases: []string{"r"},
						Value:   ".jitsuin",
						Usage:   "Select a `FILE` with Client Secret config",
					},
					&cli.StringFlag{
						Name:    "profile",
						Aliases: []string{"p"},
						Usage:   "Name a `PROFILE` to read in",
						Value:   "default",
					},
					&cli.BoolFlag{
						Name:    "inline",
						Aliases: []string{"l"},
						Usage:   "Requires RKVST_URL, RKVST_TENANT_ID, RKVST_CLIENT_ID and RKVST_CLIENT_SECRET to be set either in env or using flags",
					},
					&cli.StringFlag{
						Name:    "rkvst-url",
						Aliases: []string{"url"},
						Usage:   "RKVST Tenant ID",
						EnvVars: []string{"RKVST_URL"},
					},
					&cli.StringFlag{
						Name:    "tenant-id",
						Aliases: []string{"tid"},
						Usage:   "RKVST Tenant ID",
						EnvVars: []string{"RKVST_TENANT_ID"},
					},
					&cli.StringFlag{
						Name:    "client-id",
						Aliases: []string{"cid"},
						Usage:   "RKVST Client ID",
						EnvVars: []string{"RKVST_CLIENT_ID"},
					},
					&cli.StringFlag{
						Name:    "client-secret",
						Aliases: []string{"csec"},
						Usage:   "RKVST Cient Secret",
						EnvVars: []string{"RKVST_CLIENT_SECRET"},
					},
				},
				Action: func(c *cli.Context) error {
					var archivistInit string
					var err error

					if c.Bool("interactive") {
						if c.Bool("use-client-secret") {
							archivistInit, err = archivist.InitAsk("authTypeSecret")
							if err != nil {
								return err
							}
						} else {
							archivistInit, err = archivist.InitAsk("")
							if err != nil {
								return err
							}
						}
					} else if c.Bool("inline") {
						rkvstURL := c.String("rkvst-url")
						tenantID := c.String("tenant-id")
						clientID := c.String("client-id")
						clientSecret := c.String("client-secret")
						archivistInit, err = archivist.ClientSecretLogin(rkvstURL, tenantID, clientID, clientSecret)
						if err != nil {
							return err
						}
					} else {
						archivistInit, err = archivist.InitRead(c.String("read-file"), c.String("profile"))
						if err != nil {
							return err
						}
					}
					if c.Bool("base64") {
						encoded := base64.StdEncoding.EncodeToString([]byte(archivistInit))
						fmt.Println(encoded)
					} else {
						fmt.Println("\nToken: \n\n" + archivistInit)
					}
					if c.Bool("write-file") {
						f, err := os.Create(".auth_token")
						if err != nil {
							return err
						}
						defer f.Close()
						tokenWrite, err := f.WriteString(archivistInit)
						if err != nil {
							return err
						}
						fmt.Printf("\nWrote %d bytes to .auth_token\n", tokenWrite)
					}
					return nil
				},
			},
			{
				Name:    "asset",
				Aliases: []string{"a"},
				Action: func(c *cli.Context) error {
					err := archivist.CheckTokenExists()
					if err != nil {
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

	app.EnableBashCompletion = true
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
