package main

import (
	"github.com/Bibliotheque-microservice/emprunts/handlers"
	"github.com/gofiber/fiber/v2"
)

func setupRoutes(app *fiber.App) {
	app.Get("/v1", handlers.Home)

	// Route pour mettre à jour l'emprunt
	app.Put("/v1", handlers.UpdateEmprunts)

	// Nouvelle route pour créer un emprunt
	app.Post("/v1/emprunt", handlers.CreateEmprunt)

	// Middleware global pour gérer les 404
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Route not found",
			"method":  c.Method(),
			"path":    c.Path(),
			"message": "Please check the API documentation.",
		})
	})

}
