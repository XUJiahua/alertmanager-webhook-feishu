package feishu

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xujiahua/alertmanager-webhook-feishu/config"
	"github.com/xujiahua/alertmanager-webhook-feishu/model"
	"github.com/xujiahua/alertmanager-webhook-feishu/tmpl"
	"net/http"
	"text/template"
)

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
	// @xxx
	openIDs, err := getOpenIDs(bot.Mention, helper)
	if err != nil {
		return nil, err
	}

	// template
	tpl, err := getTemplate(bot.Template)

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

func getTemplate(tmplConf *config.Template) (*template.Template, error) {
	if tmplConf != nil && tmplConf.CustomPath != "" {
		t, err := tmpl.GetCustomTemplate(tmplConf.CustomPath)
		if err != nil {
			return nil, err
		}
		return t, nil
	}

	filename := "default.tmpl"
	if tmplConf != nil && tmplConf.EmbedFilename != "" {
		filename = tmplConf.EmbedFilename
	}

	t, err := tmpl.GetEmbedTemplate(filename)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (b Bot) Send(alerts *model.WebhookMessage) error {
	// attach @xxx
	alerts.OpenIDs = b.openIDs

	var buf bytes.Buffer
	err := b.tpl.Execute(&buf, alerts)
	if err != nil {
		return err
	}

	// TODO: move to sdk
	req, err := http.NewRequest("POST", b.webhook, &buf)
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
