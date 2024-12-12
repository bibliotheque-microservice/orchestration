package handlers

import (
	"fmt"
	"time"

	"github.com/Bibliotheque-microservice/emprunts/database"
	"github.com/Bibliotheque-microservice/emprunts/models"
	rabbitmq "github.com/Bibliotheque-microservice/emprunts/rabbitmq"
	"github.com/Bibliotheque-microservice/emprunts/structures"
	"gorm.io/gorm"
)

func CheckPenalities() ([]models.Emprunt, error) {

	const penality = 0.2

	// Get current time
	current_time := time.Now()

	fmt.Print(current_time)

	// Get all emprunts that returned date is lower than today
	// and effective return is null
	var late_emprunts []models.Emprunt
	err := database.DB.Db.Where("date_retour_effectif IS NULL AND date_retour_prevu < ?", current_time).Find(&late_emprunts).Error

	// For each row
	if err != nil {
		return nil, err
	}

	for _, emprunt := range late_emprunts {
		// Verify if penality already exists
		var penalty models.Penalite
		err := database.DB.Db.Where("emprunt_id = ?", emprunt.IDEmprunt).First(&penalty).Error

		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// Any penalities has been found, create a new penality
				newPenalty := models.Penalite{
					EmpruntID:     emprunt.IDEmprunt,
					UtilisateurID: emprunt.UtilisateurID,
					Montant:       penality,
				}

				database.DB.Db.Create(&newPenalty)

				var penaltyRecorded models.Penalite
				err := database.DB.Db.Preload("Emprunt").First(&penaltyRecorded, "emprunt_id = ? AND date_paiement IS NULL", newPenalty.EmpruntID).Error
				// For each row
				if err != nil {
					return nil, err
				}

				payload := structures.Penality_payload{
					PenalityID: int(penaltyRecorded.IDPenalite),
					EmpruntID:  int(penaltyRecorded.EmpruntID),
					Amount:     penaltyRecorded.Montant,
					UserID:     int(penaltyRecorded.UtilisateurID),
					CreatedAt:  penaltyRecorded.CreatedAt,
					UpdatedAt:  penaltyRecorded.UpdatedAt,
				}
				rabbitmq.PublishMessage("penality_exchange", "user.v1.penalities.new", payload)
			} else {
				return nil, err
			}
		} else {

			// If penality exists, just update the fee amount
			penalty.Montant += penality
			database.DB.Db.Save(&penalty)

			payload := structures.Penality_payload{
				PenalityID: int(penalty.IDPenalite),
				EmpruntID:  int(penalty.EmpruntID),
				Amount:     penalty.Montant,
				UserID:     int(penalty.UtilisateurID),
				CreatedAt:  penalty.CreatedAt,
				UpdatedAt:  penalty.UpdatedAt,
			}

			rabbitmq.PublishMessage("penality_exchange", "user.v1.penalities.updated", payload)
		}

	}

	return late_emprunts, nil
}

func RemovePenality(penality_id int) {

	err := database.DB.Db.Model(&models.Penalite{}).
		Where("id_penalite = ?", penality_id).
		Update("date_paiement", time.Now()).Error

	if err != nil {
		fmt.Printf("Erreur modification : %v", err)
	} else {

		var late_emprunts []models.Penalite
		err := database.DB.Db.Model(&models.Penalite{}).
			Where("id_penalite = ?", penality_id).Find(&late_emprunts).Error

		fmt.Print(late_emprunts)
		fmt.Printf("Successfully updated %v", err)

	}

}
