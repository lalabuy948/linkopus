package entity

// LinkMap entity
type LinkMap struct {
	Link     string `json:"link"`
	LinkHash string `json:"link_hash"`
}

func NewLinkMap(link string, linkHash string) LinkMap {
	return LinkMap{
		Link:     link,
		LinkHash: linkHash,
	}
}
