package service

// LinkMap entity
type LinkMap struct {
	Link     string `json:"link"`
	LinkHash string `json:"link_hash"`
}

// LinkView entity
type LinkView struct {
	Link   string `json:"link"`
	Date   string `json:"date"`
	Amount int    `json:"amount"`
}
