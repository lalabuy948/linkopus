package entity

// LinkView entity
type LinkView struct {
	Link     string `json:"link"`
	LinkHash string `json:"linkHash"`
	Date     string `json:"date"`
	Amount   int    `json:"amount"`
}

func NewLinkView(link string, linkHash string, date string, amount int) LinkView {
	return LinkView{
		Link:     link,
		LinkHash: linkHash,
		Date:     date,
		Amount:   amount,
	}
}
