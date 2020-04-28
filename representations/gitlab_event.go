package representations

type GitlabPushEvent struct {
	After      string     `json:"after"`
	Repository Repository `json:"repository"`
}

type GitlabPullRequestEvent struct {
	PullRequest GitlabPullRequest `json:"object_attributes"`
	Repository  Repository        `json:"repository"`
}

type GitlabPullRequest struct {
	Id             int64        `json:"id"`
	LastCommit     GitlabCommit `json:"last_commit"`
	MergeCommitSha string
	MergeStatus    string `json:"merge_status"`
	Merged         bool
	Action         string `json:"action"`
}

type GitlabCommit struct {
	Id      string `json:id`
	message string `json:message`
}
