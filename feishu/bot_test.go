package feishu

import (
	"github.com/prometheus/alertmanager/template"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/xujiahua/alertmanager-webhook-feishu/config"
	"github.com/xujiahua/alertmanager-webhook-feishu/model"
	"testing"
	"time"
)

func getConf() *config.Config {
	conf, err := config.Load("../config.yml")
	if err != nil {
		panic(err)
	}
	return conf
}

func getBotConf() *config.Bot {
	for _, bot := range getConf().Bots {
		if bot.Mention != nil {
			continue
		}
		return bot
	}
	panic("expect at least one")
}

func getAppConf() *config.App {
	return getConf().App
}

func TestBot_Send(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	bot, err := New(getBotConf(), nil)
	require.Nil(t, err)
	alerts := model.WebhookMessage{Data: newAlerts()}
	err = bot.Send(&alerts)
	require.Nil(t, err)
}

// copyright: https://github.com/tomtom-international/alertmanager-webhook-logger/blob/master/main_test.go#L132
func newAlerts() template.Data {
	return template.Data{
		Alerts: template.Alerts{
			template.Alert{
				Status:       "firing",
				Annotations:  map[string]string{"a_key": "a_value"},
				Labels:       map[string]string{"l_key": "l_value"},
				StartsAt:     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				EndsAt:       time.Date(2000, 1, 1, 0, 0, 1, 0, time.UTC),
				GeneratorURL: "file://generatorUrl",
			},
			template.Alert{
				Annotations: map[string]string{"a_key_warn": "a_value_warn"},
				Labels:      map[string]string{"l_key_warn": "l_value_warn"},
				Status:      "resolved",
			},
		},
		CommonAnnotations: map[string]string{"ca_key": "ca_value"},
		CommonLabels:      map[string]string{"cl_key": "cl_value"},
		GroupLabels:       map[string]string{"gl_key": "gl_value"},
		ExternalURL:       "file://externalUrl",
		Receiver:          "test-receiver",
	}
}
