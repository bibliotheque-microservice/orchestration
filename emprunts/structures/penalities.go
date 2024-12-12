package structures

import "time"

type Emprunt_en_retard struct {
	userId    int
	empruntId int
}

type Penality_payload struct {
	PenalityID int       `json:"penalityId"`
	EmpruntID  int       `json:"empruntId"`
	Amount     float64   `json:"amount"`
	UserID     int       `json:"userId"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Penality_paye_payload struct {
	PenalityID int `json:"id_penalite"`
}
