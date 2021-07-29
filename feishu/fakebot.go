package feishu

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/xujiahua/alertmanager-webhook-feishu/model"
)

type FakeBot struct {
}

func (f FakeBot) Send(message *model.WebhookMessage) error {
	spew.Dump(message)
	return nil
}
