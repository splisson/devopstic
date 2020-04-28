package representations

type Repository struct {
	Id   int64  `json:"id    "`
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Head struct {
	Sha string `json:"sha"`
}
