package site

import (
	"github.com/gofiber/fiber/v2"

	"github.com/th0th/poeticmetric/backend/pkg/restapi/middleware/permission"
)

func Add(app *fiber.App) {
	group := app.Group("/sites")

	group.Get("/public/:domain", readPublic)
	group.Get("/google-search-console-sites", listGoogleSearchConsoleSites)

	group.Use(permission.UserAuthenticated, permission.OrganizationSubscription)

	group.Get("/", list)
	group.Post("/", create)

	group.Get("/:id", read)
	group.Delete("/:id", destroy)
	group.Patch("/:id", update)
}
