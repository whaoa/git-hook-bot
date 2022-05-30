package git

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"time"
)

type (
	Raw struct {
		Type string
		Body io.ReadCloser
	}

	User struct {
		Id       int    `json:"id"`
		Name     string `json:"name"`
		Avatar   string `json:"avatar_url"`
		Username string `json:"username"`
		FullName string `json:"full_name"`
	}

	Repository struct {
		Name        string `json:"name"`
		FullName    string `json:"full_name"`
		Description string `json:"description"`
		WebUrl      string `json:"html_url"`
	}

	Branch struct {
		Label string `json:"label"`
		Ref   string `json:"ref"`
		Sha1  string `json:"sha"`
	}

	Commit struct {
		Sha1      string    `json:"id"`
		Message   string    `json:"message"`
		Url       string    `json:"url"`
		Timestamp time.Time `json:"timestamp"`
		Removed   []string  `json:"removed"`
		Added     []string  `json:"added"`
		Modified  []string  `json:"modified"`
		Committer User      `json:"committer"`
	}

	PullRequest struct {
		Id        int        `json:"id"`
		Url       string     `json:"url"`
		Number    int        `json:"number"`
		User      User       `json:"user"`
		State     string     `json:"state"`
		Base      Branch     `json:"base"`
		Head      Branch     `json:"head"`
		Title     string     `json:"title"`
		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt time.Time  `json:"updated_at"`
		ClosedAt  *time.Time `json:"closed_at"`
		Merged    bool       `json:"merged"`
		MergedAt  *time.Time `json:"merged_at"`
		MergedBy  *User      `json:"merged_by"`
	}

	WebHookEvent struct {
		raw         []byte
		Type        string      `json:"type"`
		Ref         string      `json:"ref"`
		BeforeSha1  string      `json:"before"`
		AfterSha1   string      `json:"after"`
		Commits     []Commit    `json:"commits"`
		Sender      User        `json:"sender"`
		Repository  Repository  `json:"repository"`
		PullRequest PullRequest `json:"pull_request"`
	}
)

func (u User) NameText() string {
	if u.Username != "" {
		return u.Username
	}
	if u.Name != "" {
		return u.Name
	}
	if u.FullName != "" {
		return u.FullName
	}
	return ""
}

func ParseWebHook(params Raw) *WebHookEvent {
	body, _ := ioutil.ReadAll(params.Body)

	hook := &WebHookEvent{raw: body, Type: params.Type}
	_ = json.Unmarshal(hook.raw, hook)

	return hook
}
