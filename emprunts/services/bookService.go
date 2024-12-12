// vérification la disponibilité du livre et le statut de l'utilisateur
package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Vérifie la disponibilité du livre via le service externe
func CheckBookAvailability(bookID int) (bool, error) {
	apiBookURL := os.Getenv("API_BOOK")
	if apiBookURL == "" {
		return false, fmt.Errorf("la variable d'environnement API_BOOK n'est pas définie")
	}

	// Construction de l'URL
	url := fmt.Sprintf("http://%s/books/%d/availability", apiBookURL, bookID)
	fmt.Printf("URL construite pour la disponibilité des livres : %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var availability struct {
		Available bool `json:"availability"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&availability); err != nil {
		return false, err
	}

	return availability.Available, nil
}
