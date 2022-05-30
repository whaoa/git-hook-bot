package endpoint

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"

	"github.com/whaoa/git-hook-bot/pkg/webhook/git"
)

const larkSimpleBotApi = "https://open.feishu.cn/open-apis/bot/v2/hook/%s"

var larkSimpleBotTemplate *template.Template

type larkSimpleBotApiResponse struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	StatusCode int    `json:"StatusCode"`
}

type LarkSimpleBot struct {
	secret string
}

func (b *LarkSimpleBot) name() string {
	return NameLarkSimpleBot
}

func (b *LarkSimpleBot) send(content string) error {
	resp, err := resty.New().R().SetBody(content).Post(fmt.Sprintf(larkSimpleBotApi, b.secret))
	if err != nil {
		return err
	}

	res := &larkSimpleBotApiResponse{}
	if err = json.Unmarshal(resp.Body(), res); err != nil {
		return err
	}

	if res.Code != 0 {
		return errors.New(res.Message)
	}
	return nil
}

func (b *LarkSimpleBot) generateContent(hook *git.WebHookEvent) (string, error) {
	if larkSimpleBotTemplate == nil {
		tmpl, err := template.New(NameLarkSimpleBot).Parse(larkSimpleBotTemplateText)

		if err != nil {
			return "", err
		}

		larkSimpleBotTemplate = tmpl
	}

	dateFormat := "2006-01-02 15:04:05"
	data := map[string]string{
		"endpoint":          NameLarkSimpleBot,
		"event_type":        hook.Type,
		"event_ref":         hook.Ref,
		"repo_full_name":    hook.Repository.FullName,
		"repo_html_url":     hook.Repository.WebUrl,
		"sender_name":       hook.Sender.NameText(),
		"commits_text":      "",
		"pull_request_text": "",
	}
	// Commits Text
	for index, commit := range hook.Commits {
		message := ""
		msg := strings.Split(commit.Message, "\n")
		if len(msg) > 0 {
			message = msg[0]
		}
		data["commits_text"] += fmt.Sprintf(
			"%d. [%s](%s) %s\\n",
			index+1,
			commit.Sha1[:8],
			commit.Url,
			strings.TrimSpace(message),
		)
	}
	// Pull Request Text
	if hook.PullRequest.Number != 0 {
		data["pull_request_text"] = fmt.Sprintf(
			"[%s] [#%d: %s](%s)\\n",
			hook.PullRequest.State,
			hook.PullRequest.Number,
			hook.PullRequest.Title,
			hook.PullRequest.Url,
		)
		data["pull_request_text"] += fmt.Sprintf("- description: *%s <- %s (by %s)*\\n",
			hook.PullRequest.Base.Label,
			hook.PullRequest.Head.Label,
			hook.PullRequest.User.NameText(),
		)
		data["pull_request_text"] += "- created at: *" + hook.PullRequest.CreatedAt.Format(dateFormat) + "*\\n"
		data["pull_request_text"] += "- updated at: *" + hook.PullRequest.UpdatedAt.Format(dateFormat) + "*\\n"

		if hook.PullRequest.Merged {
			merged := ""
			if hook.PullRequest.MergedAt != nil {
				merged += hook.PullRequest.MergedAt.Format(dateFormat)
			}
			if hook.PullRequest.MergedBy != nil {
				name := hook.PullRequest.MergedBy.NameText()
				if name != "" {
					if len(merged) > 0 {
						merged += " "
					}
					merged += "(by " + name + ")"
				}
			}
			if len(merged) > 0 {
				data["pull_request_text"] += "- merged at: *" + merged + "*\\n"
			}
		} else if hook.PullRequest.ClosedAt != nil {
			data["pull_request_text"] += "- closed at: *" + hook.PullRequest.ClosedAt.Format(dateFormat) + "*\\n"
		}
	}

	writer := bytes.NewBufferString("")
	err := larkSimpleBotTemplate.Execute(writer, data)

	if err != nil {
		return "", err
	}

	contentBytes, err := ioutil.ReadAll(writer)

	if err != nil {
		return "", err
	}

	return string(contentBytes), nil
}

const larkSimpleBotTemplateText = `
{
  "msg_type": "interactive",
  "card": {
	"config": { "wide_screen_mode": true },
	  "header": {
		"template": "green",
		"title": { "content": "Git Webhook for {{ .endpoint }}", "tag": "plain_text" }
	  },
	  "elements": [
		{
		  "tag": "div",
		  "fields": [
			{
			  "is_short": true,
			  "text": { "tag": "lark_md", "content": "**Repository：**\n[{{ .repo_full_name }}]({{ .repo_html_url }})" }
			},
			{
			  "is_short": true,
			  "text": { "tag": "lark_md", "content": "**Event Type：**\n{{ .event_type }}" }
			},
			{ "is_short": false, "text": { "content": "", "tag": "lark_md" } },
			{
			  "is_short": true,
			  "text": { "tag": "lark_md", "content": "**Pusher：**\n{{ .sender_name }}" }
			},
			{
			  "is_short": true,
			  "text": { "tag": "lark_md", "content": "**Event Ref：**\n{{ .event_ref }}" }
			}
		  ]
		},

        {{- if (gt (len .commits_text) 0) }}
		{ "tag": "hr" },
		{ "tag": "markdown", "content": "**Commits**\n{{ .commits_text }}" },
		{{ end -}}

        {{- if (gt (len .pull_request_text) 0) }}
		{ "tag": "hr" },
		{ "tag": "markdown", "content": "**Pull Request**\n{{ .pull_request_text }}" },
        {{ end -}}

		{ "tag": "hr" },
		{
		  "tag": "note",
		  "elements": [{ "tag": "plain_text", "content": "Git Webhook for {{ .endpoint }} @git-hook-bot" }]
		}
	  ]
	}
  }
}
`
