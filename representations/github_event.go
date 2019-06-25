package representations

type GithubPushEvent struct {
	After      string     `json:"after"`
	Repository Repository `json:"repository"`
}

type Repository struct {
	Id   int64  `json:"id    "`
	Name string `json:"name"`
}

type GithubPullRequestEvent struct {
	Action      string      `json:"action"`
	PullRequest PullRequest `json:"pull_request"`
	Repository  Repository  `json:"repository"`
}

type Head struct {
	Sha string `json:"sha"`
}

type PullRequest struct {
	Head           Head   `json:"head"`
	Id             int64  `json:"id"`
	Number         int64  `json:"number"`
	MergeCommitSha string `json:"merge_commit_sha"`
	Merged         bool   `json:"merged"`
}
