package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type userStatusResponse struct {
	Valid   bool   `json:"valid"`
	Message string `json:"message,omitempty"`
}

// Vérifie si l'utilisateur est actif et s'il n'a pas trop de pénalités
func CheckUserStatus(userID int) (bool, string, error) {
	apiBookURL := os.Getenv("API_USER")
	url := fmt.Sprintf("http://%s/valid-user/%d", apiBookURL, userID)
	resp, err := http.Get(url)
	if err != nil {
		return false, "", fmt.Errorf("failed to get user status: %w", err)
	}
	defer resp.Body.Close()

	var userStatus userStatusResponse

	if err := json.NewDecoder(resp.Body).Decode(&userStatus); err != nil {
		return false, "", fmt.Errorf("failed to decode response: %w", err)
	}

	if !userStatus.Valid {
		return false, userStatus.Message, fmt.Errorf(userStatus.Message)
	}

	return true, userStatus.Message, nil
}
