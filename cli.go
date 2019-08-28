package main

import "github.com/urfave/cli"

func info() {
	app.Name = "Vip manager"
	app.Usage = "A Simple CLI for managing Vips"
	app.Author = "Evolução Infra - Team"
	app.Version = "0.1"
}

func flags() {
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "env, e",
			Value: "qa",
			Usage: "NetworkApi env prod/qa",
		},
		cli.StringFlag{
			Name:   "ldap_user",
			Usage:  "User from ladp to auth on the napi (without domain)",
			EnvVar: "ldap_user",
		},
		cli.StringFlag{
			Name:   "ldap_pass",
			Usage:  "Pass from ladp to auth on the napi",
			EnvVar: "ldap_pass",
		},
	}
}

func commands() {
	app.Commands = []cli.Command{
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "Add machines to the Vip pools",
			Action:  func(c *cli.Context) { updateVip(c, ADD) },
			Flags: []cli.Flag{
				cli.StringFlag{Name: "vip, v"},
				cli.StringFlag{Name: "real, r"},
			},
		},
		{
			Name:    "rem",
			Aliases: []string{"r"},
			Usage:   "Removing machines from Vip pools",
			Action:  func(c *cli.Context) { updateVip(c, REM) },
			Flags: []cli.Flag{
				cli.StringFlag{Name: "vip, v"},
				cli.StringFlag{Name: "real, r"},
			},
		},
	}
}
