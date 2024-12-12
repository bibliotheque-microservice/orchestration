package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Bibliotheque-microservice/emprunts/cron"
	"github.com/Bibliotheque-microservice/emprunts/database"
	"github.com/Bibliotheque-microservice/emprunts/handlers"
	middleware "github.com/Bibliotheque-microservice/emprunts/middleware"
	rabbitmq "github.com/Bibliotheque-microservice/emprunts/rabbitmq"
	"github.com/Bibliotheque-microservice/emprunts/structures"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func main() {

	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)

	database.ConnectDb()

	rabbitmq.InitRabbitMQ()
	defer rabbitmq.CloseRabbitMQ()

	// go func() {
	// 	penality_msgs := rabbitmq.ConsumeMessages("user_penalties_queue")

	// 	for msg := range penality_msgs {
	// 		switch msg.RoutingKey {
	// 		case "user.v1.penalities.new":
	// 			var jsonData structures.PenaltyMessage
	// 			err := json.Unmarshal(msg.Body, &jsonData)
	// 			if err != nil {
	// 				log.Printf("Erreur de parsing JSON : %v", err)
	// 				continue
	// 			}
	// 			log.Printf("Message JSON reçu : %+v", jsonData)
	// 			log.Print(jsonData.Amount)

	// 		case "user.v1.penalities.paid":
	// 			log.Printf("Message reçu : %s", string(msg.Body))
	// 		default:
	// 			log.Printf("Message non géré avec Routing Key : %s", msg.RoutingKey)
	// 			log.Printf("Contenu brut du message : %s", string(msg.Body))
	// 		}
	// 	}
	// }()

	// go func() {
	// 	emprunts_msg := rabbitmq.ConsumeMessages("emprunts_finished_queue")

	// 	for msg := range emprunts_msg {
	// 		switch msg.RoutingKey {
	// 		case "emprunts.v1.finished":
	// 			var jsonData interface{}
	// 			err := json.Unmarshal(msg.Body, &jsonData)
	// 			if err != nil {
	// 				log.Printf("Erreur de parsing JSON : %v", err)
	// 				continue
	// 			}
	// 			log.Printf("Message JSON reçu : %+v", jsonData)
	// 		default:
	// 			log.Printf("Message reçu : %+v", msg)
	// 			log.Printf("Message non géré avec Routing Key : %s", msg.RoutingKey)
	// 			log.Printf("Contenu brut du message : %s", string(msg.Body))
	// 		}
	// 	}
	// }()

	// go func() {
	// 	emprunts_msg := rabbitmq.ConsumeMessages("emprunts_created_queue")

	// 	for msg := range emprunts_msg {
	// 		// Vérifie la clé de routage (RoutingKey) du message
	// 		switch msg.RoutingKey {
	// 		case "emprunts.v1.finished":
	// 			// Déclarer une variable pour stocker les données désérialisées
	// 			var messageData map[string]interface{}

	// 			// Désérialiser le corps du message JSON
	// 			err := json.Unmarshal(msg.Body, &messageData)
	// 			if err != nil {
	// 				log.Printf("Erreur de parsing JSON : %v", err)
	// 				continue
	// 			}

	// 			// Afficher les données reçues après la désérialisation
	// 			log.Printf("Message JSON reçu : %+v", messageData)

	// 		default:
	// 			var messageData map[string]interface{}

	// 			// Désérialiser le corps du message JSON
	// 			err := json.Unmarshal(msg.Body, &messageData)
	// 			if err != nil {
	// 				log.Printf("Erreur de parsing JSON : %v", err)
	// 				continue
	// 			}

	// 			// Afficher les données reçues après la désérialisation
	// 			log.Printf("Message JSON reçu : %+v", messageData)

	// 		}
	// 	}
	// }()

	go func() {
		penalities_msg := rabbitmq.ConsumeMessages("paiements_queue")

		for msg := range penalities_msg {
			var jsonData structures.Penality_paye_payload
			err := json.Unmarshal(msg.Body, &jsonData)
			if err != nil {
				log.Printf("Erreur de parsing JSON : %v", err)
				continue
			}
			log.Printf("Message JSON reçu paiement : %+v", jsonData)
			handlers.RemovePenality(jsonData.PenalityID)
		}
	}()

	go func() {
		cron.StartCron()
	}()
	app := fiber.New()

	app.Use(middleware.LoggerMiddleware(log))

	setupRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Info("Application démarrée")
	log.Warn("Ceci est un message de niveau Warn")

	app.Listen(fmt.Sprintf(":%s", port))
}
