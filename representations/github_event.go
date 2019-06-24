package representations

type GithubEvent struct {
	Name       string     `json:"name"`
	Action     string     `json:"action"`
	Head       string     `json:"head"`
	Repository Repository `json:"repository"`
}

type Repository struct {
	Id   int64  `json:"id    "`
	Name string `json:"name"`
}
