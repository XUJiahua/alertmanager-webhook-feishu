package feishu

import (
	"bytes"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xujiahua/alertmanager-webhook-feishu/config"
	"github.com/xujiahua/alertmanager-webhook-feishu/model"
	"net/http"
	"text/template"
)

//go:embed default.gojson
var f embed.FS

type fsWebhookResponse struct {
	StatusCode    int    `json:"StatusCode"`
	StatusMessage string `json:"StatusMessage"`
}

type Bot struct {
	webhook string
	openIDs []string
	client  http.Client
	tpl     *template.Template
}

func New(bot *config.Bot, helper *EmailHelper) (*Bot, error) {
	openIDs, err := getOpenIDs(bot.Mention, helper)
	if err != nil {
		return nil, err
	}

	tpl := template.Must(template.ParseFS(f, "default.gojson"))

	return &Bot{
		webhook: bot.Webhook,
		openIDs: openIDs,
		client:  http.Client{},
		tpl:     tpl,
	}, nil
}

func getOpenIDs(mention *config.Mention, helper *EmailHelper) ([]string, error) {
	if mention == nil {
		return nil, nil
	}
	if mention.All {
		return []string{"all"}, nil
	}

	openIDs := mention.OpenIDs
	emails := mention.Emails
	if len(emails) != 0 && helper == nil {
		return nil, errors.New("@somebody by email need email flag enabled")
	}
	if len(emails) != 0 {
		remaining, err := helper.Lookup(emails)
		if err != nil {
			return nil, err
		}
		openIDs = append(openIDs, remaining...)
	}
	return openIDs, nil
}

func (b Bot) Send(alerts *model.WebhookMessage) error {
	// attach @xxx
	alerts.OpenIDs = b.openIDs

	fsMessage, err := b.toFeishuCard(alerts)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", b.webhook, bytes.NewBufferString(fsMessage))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := b.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var fsResp fsWebhookResponse
	err = json.NewDecoder(resp.Body).Decode(&fsResp)
	if err != nil {
		return err
	}

	if fsResp.StatusCode != 0 {
		return errors.New(fmt.Sprintf("code: %d, err: %s", fsResp.StatusCode, fsResp.StatusMessage))
	}

	return nil
}

// TODO: template factory
func (b Bot) toFeishuCard(alerts *model.WebhookMessage) (string, error) {
	var buf bytes.Buffer
	err := b.tpl.Execute(&buf, alerts)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
