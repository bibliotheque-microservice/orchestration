package handlers

import (
	"time"

	"github.com/Bibliotheque-microservice/emprunts/database"
	"github.com/Bibliotheque-microservice/emprunts/models"
	"github.com/Bibliotheque-microservice/emprunts/rabbitmq"
	"github.com/Bibliotheque-microservice/emprunts/services" // Import des services pour vérifier livre et utilisateur
	"github.com/Bibliotheque-microservice/emprunts/structures"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	// Configuration du logger
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)
}

// Route pour vérifier et créer un emprunt
func CreateEmprunt(c *fiber.Ctx) error {
	// Parse la requête JSON pour obtenir l'ID de l'utilisateur et du livre
	var request structures.EmpruntRequest
	if err := c.BodyParser(&request); err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Error("Invalid input")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Vérifier la disponibilité du livre
	available, err := services.CheckBookAvailability(request.BookID)
	if err != nil || !available {
		log.WithFields(logrus.Fields{
			"book_id": request.BookID,
		}).Warn("Le livre n'est pas disponible")
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Le livre n'est pas disponible"})
	}

	// Vérifier l'état de l'utilisateur (actif et pas de pénalités)
	userAuthorized, msg, err := services.CheckUserStatus(request.UserID)
	if err != nil || !userAuthorized {

		log.WithFields(logrus.Fields{
			"error": err,
		}).Error(msg)
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": msg})
	}

	// Créer l'emprunt dans la base de données
	emprunt := models.Emprunt{
		UtilisateurID:   uint(request.UserID),
		LivreID:         uint(request.BookID),
		DateEmprunt:     time.Now(),
		DateRetourPrevu: time.Now().Add(14 * 24 * time.Hour), // 2 semaines de durée
	}

	err = database.DB.Db.Create(&emprunt).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Erreur lors de la création de l'emprunt"})
	}

	// Publier un message via RabbitMQ pour notifier l'emprunt
	empruntMessage := map[string]interface{}{
		"livreId":       request.BookID,
		"disponible":    false,
		"idUtilisateur": request.UserID,
	}

	// Publier le message à RabbitMQ
	rabbitmq.PublishMessage("emprunts_exchange", "emprunts.v1.created", empruntMessage)

	// Répondre avec succès et autoriser l'emprunt
	return c.JSON(fiber.Map{"autorisé": true, "message": "Emprunt créé avec succès"})
}

func Home(c *fiber.Ctx) error {
	return c.SendString("Hello ma vie!")
}

func UpdateEmprunts(c *fiber.Ctx) error {
	// Parse la requete
	var empruntRequest structures.EmpruntReturned
	if err := c.BodyParser(&empruntRequest); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Si l'emprunt a été retourné, updater la date de retour
	if empruntRequest.Returned {
		err := database.DB.Db.Model(&models.Emprunt{}).
			Where("id_emprunt = ?", empruntRequest.EmpruntID).
			Update("date_retour_effectif", time.Now()).Error

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Erreur lors de la mise à jour de l'emprunt avec bdd.",
				"error":   err.Error(),
			})
		}

		// retrieve the whole emprunt
		var updatedEmprunt models.Emprunt
		err = database.DB.Db.First(&updatedEmprunt, "id_emprunt = ?", empruntRequest.EmpruntID).Error

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Erreur lors de la mise à jour de l'emprunt avec bdd.",
				"error":   err.Error(),
			})
		}

		rabbitmq.PublishMessage("emprunts_exchange", "emprunts.v1.finished", updatedEmprunt)

		// Envoyer un message aux cosnumers pour indiquer que le livre a été retourné
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Date de retour mise à jour avec succès.",
		})
	} else {
		return c.Status(400).SendString("Livre not returned")
	}

}
