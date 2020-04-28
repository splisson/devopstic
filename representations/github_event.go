package representations

type GithubPushEvent struct {
	After      string     `json:"after"`
	Repository Repository `json:"repository"`
}

type GithubPullRequestEvent struct {
	Action      string            `json:"action"`
	PullRequest GithubPullRequest `json:"pull_request"`
	Repository  Repository        `json:"repository"`
}

type GithubPullRequest struct {
	Head           Head   `json:"head"`
	Id             int64  `json:"id"`
	Number         int64  `json:"number"`
	MergeCommitSha string `json:"merge_commit_sha"`
	Merged         bool   `json:"merged"`
}
