package url

type Url struct {
	Source      string `json:"source"`
	ShortUrl    string `json:"short_url"`
	DateCreated string `json:"date_created"`
	User        int64  `json:"user"`
}
