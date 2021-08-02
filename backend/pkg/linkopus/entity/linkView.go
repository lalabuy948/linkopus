package entity

// LinkView entity
type LinkView struct {
	Link     string `json:"link"`
	LinkHash string `json:"linkHash"`
	Date     string `json:"date"`
	Amount   int    `json:"amount"`
}
