package main

import (
	"fmt"
	"github.com/gofiber/fiber"
	"github.com/urfave/cli"
	"github.com/zikwall/blogchain/actions"
	"github.com/zikwall/blogchain/middlewares"
	"log"
	"os"
)

// @title Blog Chain swagger documentation for Go service
// @version 1.0
// @description This is a sample server celler server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support Blog Chain
// @contact.url http://www.blogchain.io/support
// @contact.email support@blogchain.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host blogchain.io
func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "host",
				Value: "localhost",
				Usage: "Run service in host",
			},
			&cli.IntFlag{
				Name:  "port",
				Value: 3001,
				Usage: "Run service in port",
			},
		},
	}

	app.Action = func(c *cli.Context) error {
		host := c.String("host")
		port := c.Int("port")

		app := fiber.New()

		app.Use(middlewares.JWT)
		app.Get("/", actions.HelloWorldAction)
		app.Static("/docs", "./docs")

		err := app.Listen(fmt.Sprintf("%s:%d", host, port))

		if err != nil {
			log.Fatal(err)
		}

		return nil
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}