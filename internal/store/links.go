package store

type Link struct {
	ID          int    `json:"id"`
	OriginalURL string `json:"original_url"`
	ShortCode   string `json:"short_code"`
	CreatedAt   string `json:"created_at"`
}

func (l *Link) String() string {
	return l.OriginalURL + " -> " + l.ShortCode
}
