package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/urfave/cli/v2"
	"github.com/zikwall/blogchain/src/actions"
	"github.com/zikwall/blogchain/src/constants"
	"github.com/zikwall/blogchain/src/lib"
	"github.com/zikwall/blogchain/src/middlewares"
	"github.com/zikwall/blogchain/src/service"
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
	application := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "bind-address",
				Required: true,
				Usage:    "Run service in host",
				EnvVars:  []string{"BIND_ADDRESS"},
			},
			// database
			&cli.StringFlag{
				Name:     "database-host",
				Required: true,
				Usage:    "Database host",
				EnvVars:  []string{"DATABASE_HOST"},
			},
			&cli.StringFlag{
				Name:     "database-user",
				Required: true,
				Usage:    "Database user",
				EnvVars:  []string{"DATABASE_USER"},
			},
			&cli.StringFlag{
				Name:     "database-password",
				Required: true,
				Usage:    "Database password",
				EnvVars:  []string{"DATABASE_PASSWORD"},
			},
			&cli.StringFlag{
				Name:     "database-name",
				Required: true,
				Usage:    "Database name",
				EnvVars:  []string{"DATABASE_NAME"},
			},
			&cli.StringFlag{
				Name:     "database-dialect",
				Required: true,
				Usage:    "Database dialect: mysql, postgres, sqlite3, sqlserver etc",
				EnvVars:  []string{"DATABASE_DIALECT"},
			},
			&cli.StringFlag{
				Name:     "rsa-public-key",
				Required: false,
				Usage:    "Container secret key for JWT, and etc.",
				EnvVars:  []string{"RSA_PUBLIC_KEY"},
			},
			&cli.StringFlag{
				Name:     "rsa-private-key",
				Required: false,
				Usage:    "Container secret key for JWT, and etc.",
				EnvVars:  []string{"RSA_PRIVATE_KEY"},
			},
		},
	}

	application.Action = func(c *cli.Context) error {
		blogchain, err := service.NewBlogchainServiceInstance(
			service.BlogchainServiceConfiguration{
				BloghainDatabaseConfiguration: service.BloghainDatabaseConfiguration{
					Host:     c.String("database-host"),
					User:     c.String("database-user"),
					Password: c.String("database-password"),
					Name:     c.String("database-name"),
					Dialect:  c.String("database-dialect"),
				},
				BlogchainHttpAccessControl: service.BlogchainHttpAccessControl{
					AllowOrigins:     "*",
					AllowMethods:     "*",
					AllowHeaders:     "*",
					AllowCredentials: false,
					ExposeHeaders:    "",
					MaxAge:           0,
				},
				BlogchainContainer: service.BlogchainServiceContainerConfiguration{},
			},
		)

		if err != nil {
			return err
		}

		app := fiber.New()
		app.Static("/docs", "./src/public/docs")
		app.Static("/uploads", "./src/public/uploads")
		app.Get("/metrics", actions.PrometheusWithFastHTTPAdapter())

		app.Use(
			middlewares.WithBlogchainCORSPolicy(blogchain),
			middlewares.WithBlogchainXHeaderPolicy(blogchain),
		)

		rsa := lib.NewBlogchainRSAContainer(
			constants.TestPublicKey, constants.TestPrivateKey,
		)

		actionProvider := actions.NewBlogchainActionProvider(actions.ActionsRequiredInstances{
			RSA: &rsa,
			Db:  blogchain.GetBlogchainDatabaseInstance(),
		})

		api := app.Group("/api",
			middlewares.WithBlogchainJWTAuthorization(&rsa),
			middlewares.WithBlogchainUserIdentity(blogchain),
		)
		{
			api.Get("/healthcheck", actions.HealthCheck)
			api.Get("/runtime", actions.BlogchainRuntimeStatistic(
				blogchain.Container.GetStartedAt(),
			))

			v1 := api.Group("/v1")
			{
				v1.Get("/profile/:username", actionProvider.Profile)
				v1.Get("/content/:id", actionProvider.Content)
				v1.Get("/contents/:page?", actionProvider.Contents)
				v1.Get("/tag/:tag/:page?", actionProvider.Contents)
				v1.Get("/tags", actionProvider.Tags)
				v1.Get("/contents/user/:id/:page?", actionProvider.ContentsUser)
			}

			withAccessControlPolicy := api.Use(
				middlewares.UseBlogchainAccessControlPolicy,
			)

			editor := withAccessControlPolicy.Group("/editor")
			{
				editor.Get("/content/:id", actionProvider.ContentInformation)
				editor.Post("/content/add", actionProvider.ContentCreate)
				editor.Post("/content/update/:id", actionProvider.ContentUpdate)
			}
		}

		auth := app.Group("/auth", middlewares.UseBlogchainSignPolicy)
		{
			auth.Post("/register", actionProvider.Register)
			auth.Post("/login", actionProvider.Login)
			auth.Post("/logout", actionProvider.Logout)
		}

		go func() {
			if err := app.Listen(c.String("bind-address")); err != nil {
				blogchain.GetInternalLogger().Error(err)
			}
		}()

		blogchain.WaitBlogchainSystemNotify()
		blogchain.ShutdownBlogchainServer()

		return nil
	}

	err := application.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}
