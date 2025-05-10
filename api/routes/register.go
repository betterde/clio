package routes

import (
	"github.com/betterde/clio/api/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func RegisterRoutes(app *fiber.App) {
	api := app.Group("api/v1")

	api.Get("/health", handlers.HealthCheck).Name("api.health.check")

	auth := api.Group("auth")
	auth.Post("signup", handlers.SignUp).Name("api.auth.signup")
	auth.Post("signin", handlers.SignIn).Name("api.auth.signin")
	auth.Get(":idp/callback", handlers.Callback).Name("api.auth.callback")

	// Swagger UI router
	app.Get("/docs/*", swagger.New(swagger.Config{
		URL:          "/swagger/user.swagger.json",
		DeepLinking:  false,
		DocExpansion: "none",
	}))

	// Swagger API specification file router
	//app.Get("/swagger/*", filesystem.New(filesystem.Config{
	//	Root:               docs.Serve(),
	//	Index:              "user.swagger.json",
	//	NotFoundFile:       "user.swagger.json",
	//	ContentTypeCharset: "UTF-8",
	//})).Name("Swagger JSON Schema")

	// Embed SPA static resource
	//app.Get("*", filesystem.New(filesystem.Config{
	//	Root:               spa.Serve(),
	//	Index:              "index.html",
	//	NotFoundFile:       "index.html",
	//	ContentTypeCharset: "UTF-8",
	//})).Name("SPA static resource")
}
