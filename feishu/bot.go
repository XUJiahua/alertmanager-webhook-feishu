package feishu

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/xujiahua/alertmanager-webhook-feishu/config"
	"github.com/xujiahua/alertmanager-webhook-feishu/model"
	"github.com/xujiahua/alertmanager-webhook-feishu/tmpl"
	"strings"
	"text/template"
)

type Bot struct {
	webhook  string
	openIDs  []string
	sdk      *Sdk
	tpl      *template.Template
	alertTpl *template.Template
}

func New(bot *config.Bot, helper *EmailHelper) (*Bot, error) {
	// @xxx
	openIDs, err := getOpenIDs(bot.Mention, helper)
	if err != nil {
		return nil, err
	}

	// template
	tpl, alertTpl, err := getTemplates(bot.Template)
	if err != nil {
		return nil, err
	}

	return &Bot{
		webhook:  bot.Webhook,
		openIDs:  openIDs,
		sdk:      NewSDK("", ""),
		tpl:      tpl,
		alertTpl: alertTpl,
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

func getTemplates(tmplConf *config.Template) (*template.Template, *template.Template, error) {
	if tmplConf != nil && tmplConf.CustomPath != "" {
		t, err := tmpl.GetCustomTemplate(tmplConf.CustomPath)
		if err != nil {
			return nil, nil, err
		}
		return t, nil, nil
	}

	// by default, use two tmpls, one is for alert
	dt, err := tmpl.GetEmbedTemplate("default.tmpl")
	if err != nil {
		return nil, nil, err
	}

	dat, err := tmpl.GetEmbedTemplate("default_alert.tmpl")
	if err != nil {
		return nil, nil, err
	}

	return dt, dat, nil
}

func (b Bot) Send(alerts *model.WebhookMessage) error {
	// attach @xxx
	alerts.OpenIDs = b.openIDs

	// prepare data
	err := b.preprocessAlerts(alerts)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = b.tpl.Execute(&buf, alerts)
	if err != nil {
		return err
	}
	if logrus.IsLevelEnabled(logrus.DebugLevel) {
		fmt.Println(buf.String())
	}

	return b.sdk.WebhookV2(b.webhook, &buf)
}

func (b Bot) preprocessAlerts(alerts *model.WebhookMessage) error {
	if b.alertTpl == nil {
		return nil
	}

	// preprocess using alert template
	for _, alert := range alerts.Alerts.Firing() {
		var buf bytes.Buffer
		err := b.alertTpl.Execute(&buf, alert)
		if err != nil {
			return err
		}
		res := strings.ReplaceAll(buf.String(), "\n", "\\n")
		alerts.FiringAlerts = append(alerts.FiringAlerts, res)
	}
	for _, alert := range alerts.Alerts.Resolved() {
		var buf bytes.Buffer
		err := b.alertTpl.Execute(&buf, alert)
		if err != nil {
			return err
		}
		res := strings.ReplaceAll(buf.String(), "\n", "\\n")
		alerts.ResolvedAlerts = append(alerts.ResolvedAlerts, res)
	}

	return nil
}
