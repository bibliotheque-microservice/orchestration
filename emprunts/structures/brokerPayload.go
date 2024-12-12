package structures

type PenaltyMessage struct {
	UserId  int     `json:"userId"`
	LivreId int     `json:"livreId"`
	Amount  float64 `json:"amount"`
}
