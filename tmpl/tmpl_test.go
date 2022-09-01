package tmpl

import (
	"github.com/prometheus/alertmanager/template"
	"github.com/stretchr/testify/require"
	"github.com/xujiahua/alertmanager-webhook-feishu/model"
	"os"
	"testing"
	"time"
	"path/filepath"
	"runtime"
)

func TestFeishuCard(t *testing.T) {
	alerts := model.WebhookMessage{Data: newAlerts()}
	et := embedTemplates["default.tmpl"]
	err := et.Execute(os.Stdout, alerts)
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
				Status:      "warning",
			},
		},
		CommonAnnotations: map[string]string{"ca_key": "ca_value"},
		CommonLabels:      map[string]string{"cl_key": "cl_value"},
		GroupLabels:       map[string]string{"gl_key": "gl_value"},
		ExternalURL:       "file://externalUrl",
		Receiver:          "test-receiver",
	}
}

func unittestRelativePath() string {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		file = ""
	}

	return filepath.Dir(filepath.FromSlash(file))
}

func TestGetCustomTemplate(t *testing.T) {
	require := require.New(t)

	tmpl, err := GetCustomTemplate(filepath.Join(unittestRelativePath(), "templates", "default.tmpl"))
	require.NoError(err)
	require.NotNil(tmpl)

	alerts := model.WebhookMessage{Data: newAlerts()}
	err = tmpl.Execute(os.Stdout, alerts)
	require.NoError(err)
}
