package model

import "github.com/prometheus/alertmanager/template"

type WebhookMessage struct {
	template.Data
	// @某人
	OpenIDs []string
}
type Alert template.Alert
type Alerts template.Alerts
