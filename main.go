package main

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"os"
	"path"
	"strings"
)

func main() {
	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Usage = "An entropy and failure injection management API for Docker platforms."
	app.Version = VERSION
	app.Authors = []cli.Author{{Name: "Jeff Nickoloff", Email: "jeff@allingeek.com"}}
	app.Flags = flags
	app.Commands = commands
	app.Before = func(c *cli.Context) error {
		log.SetOutput(os.Stderr)
		level, err := log.ParseLevel(c.String("log-level"))
		if err != nil {
			log.Fatalf(err.Error())
		}
		log.SetLevel(level)

		// If a log level wasn't specified and we are running in debug mode,
		// enforce log-level=debug.
		if !c.IsSet("log-level") && !c.IsSet("l") && c.Bool("debug") {
			log.SetLevel(log.DebugLevel)
		}

		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func commandNotYetImplemented(c *cli.Context) error {
	fmt.Printf("Invocation: %s [xxx] %s\n", c.Command.Name, c.Args())
	fmt.Println("This command has yet to be implemented.")
	return nil
}

var (
	flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "debug mode",
			EnvVar: "DEBUG",
		},
		cli.StringFlag{
			Name:  "log-level, l",
			Value: "info",
			Usage: fmt.Sprintf("Log level (options: debug, info, panic)"),
		},
		cli.StringSliceFlag{
			Name:   "host, H",
			EnvVar: "ENTROPY_HOST",
			Value:  &cli.StringSlice{"tcp://:2476"},
			Usage:  "Entropy endpoint host.",
		},
	}
	// entropy run -f 30s -p .10 --fault recv_drop --target label=service=myserv
	commands = []cli.Command{
		{
			Name:      "manage",
			ShortName: "m",
			Usage:     "Start the Entropy manager",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "cacert",
					Usage: "",
					Value: "",
				},
				cli.StringFlag{
					Name:  "cert",
					Usage: "",
					Value: "",
				},
				cli.StringFlag{
					Name:  "key",
					Usage: "",
					Value: "",
				},
				cli.BoolTFlag{
					Name:  "tlsverify",
					Usage: "",
				},
				cli.StringFlag{
					Name:  "image",
					Usage: "Gremlin image",
					Value: "qualimente/gremlins",
				},
			},
			Before: func(c *cli.Context) error {
				args := c.Args()
				if len(args) == 0 {
					return errors.New("Required target is missing.")
				}
				return nil
			},
			Action: func(c *cli.Context) error {
				m := &manager{Target: c.Args()[0]}
				log.Printf("Target has been set to: %s", m.Target)
				for _, v := range c.GlobalStringSlice("host") {
					if strings.HasPrefix(v, "tcp://") {
						m.Tcp = strings.TrimPrefix(v, "tcp://")
					}
					if strings.HasPrefix(v, "unix://") {
						m.Unix = strings.TrimPrefix(v, "unix://")
					}
				}
				m.Start()
				log.Println("Done")
				return nil
			},
		},
		{
			Name:   "create",
			Usage:  "Create a failure injector",
			Action: commandNotYetImplemented,
		},
		{
			Name:   "ls",
			Usage:  "List failure injectors",
			Action: commandNotYetImplemented,
		},
		{
			Name:   "run",
			Usage:  "Create and start failure injector",
			Action: commandNotYetImplemented,
		},
		{
			Name:   "start",
			Usage:  "Start a failure injector",
			Action: commandNotYetImplemented,
		},
		{
			Name:   "stop",
			Usage:  "Stop a failure injector",
			Action: commandNotYetImplemented,
		},
		{
			Name:      "version",
			ShortName: "v",
			Usage:     "Show version",
			Action: func(c *cli.Context) error {
				PrintVersion()
				return nil
			},
		},
	}
)
