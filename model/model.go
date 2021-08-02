package model

import "github.com/prometheus/alertmanager/template"

type WebhookMessage struct {
	// reference: https://prometheus.io/docs/alerting/latest/notifications/
	template.Data
	// @某人
	OpenIDs []string
	// 仅内置模板中使用，自定义模板中访问是空数组
	// 目前没有发现在 {{template defined_name .}} 后对其结果进行进一步处理的方式
	// 首先，通过模板，将每个 Alert 转为字符串，大段文本都在 content 字段，需要注意转义。
	FiringAlerts   []string
	ResolvedAlerts []string
}

type Alert template.Alert
