package api

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/betterde/clio/api/routes"
	"github.com/betterde/clio/config"
	"github.com/betterde/clio/internal/journal"
	"github.com/betterde/clio/internal/response"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"go.uber.org/zap"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var ServerInstance *Server

type Server struct {
	Engine *fiber.App
}

func InitServer(name, version string) {
	caser := cases.Title(language.Und, cases.NoLower)
	ServerInstance = &Server{
		Engine: fiber.New(fiber.Config{
			AppName:       name,
			ServerHeader:  fmt.Sprintf("%s %s", caser.String(name), version),
			CaseSensitive: true,
			// Override default error handler
			ErrorHandler: func(ctx *fiber.Ctx, err error) error {
				// Status code defaults to 500
				code := fiber.StatusInternalServerError

				// Retrieve the custom status code if it's a fiber.*Error
				var e *fiber.Error
				if errors.As(err, &e) {
					code = e.Code
				}

				if err != nil {
					if code >= fiber.StatusInternalServerError {
						journal.Logger.Errorw("Analysis server runtime error:", zap.Error(err))
					}

					// In case the SendFile fails
					return ctx.Status(code).JSON(response.Send(code, err.Error(), nil))
				}

				return nil
			},
		}),
	}

	routes.RegisterRoutes(ServerInstance.Engine)
}

func (s *Server) Run(verbose bool) {
	ServerInstance.Engine.Use(cors.New())

	if verbose {
		ServerInstance.Engine.Use(logger.New())
	}
	ServerInstance.Engine.Use(pprof.New())
	ServerInstance.Engine.Use(recover.New())
	ServerInstance.Engine.Use(requestid.New())

	go func() {
		if config.Conf.HTTP.TLSKey != "" && config.Conf.HTTP.TLSCert != "" {
			cert, err := tls.LoadX509KeyPair(config.Conf.HTTP.TLSCert, config.Conf.HTTP.TLSKey)
			if err != nil {
				journal.Logger.Panicw("Failed to start orbit server:", err)
			}

			// Create custom listener
			ln, err := tls.Listen("tcp", config.Conf.HTTP.Listen, &tls.Config{Certificates: []tls.Certificate{cert}})
			if err != nil {
				journal.Logger.Panicw("Failed to start orbit server:", err)
			}

			err = ServerInstance.Engine.Listener(ln)
			if err != nil {
				journal.Logger.Panicw("Failed to start orbit server:", err)
			}
		} else {
			err := ServerInstance.Engine.Listen(config.Conf.HTTP.Listen)
			if err != nil {
				journal.Logger.Panicw("Failed to start orbit server:", err)
			}
		}
	}()
}
