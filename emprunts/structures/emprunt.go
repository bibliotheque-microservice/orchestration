package structures

type EmpruntReturned struct {
	EmpruntID int  `json:"empruntId"`
	Returned  bool `json:"returned"`
}

type EmpruntRequest struct {
	BookID int `json:"bookId"`
	UserID int `json:"userId"`
}
