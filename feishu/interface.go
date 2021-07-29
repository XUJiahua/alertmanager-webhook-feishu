package feishu

import "github.com/xujiahua/alertmanager-webhook-feishu/model"

type IBot interface {
	Send(*model.WebhookMessage) error
}
