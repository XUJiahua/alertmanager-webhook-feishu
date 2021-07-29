package feishu

import (
	"bytes"
	"errors"
	"github.com/xujiahua/alertmanager-webhook-feishu/config"
	"github.com/xujiahua/alertmanager-webhook-feishu/model"
	"github.com/xujiahua/alertmanager-webhook-feishu/tmpl"
	"text/template"
)

type Bot struct {
	webhook string
	openIDs []string
	sdk     *Sdk
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
		sdk:     NewSDK("", ""),
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

	// prepare data
	var buf bytes.Buffer
	err := b.tpl.Execute(&buf, alerts)
	if err != nil {
		return err
	}

	return b.sdk.WebhookV2(b.webhook, &buf)
}
