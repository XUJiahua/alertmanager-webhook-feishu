package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/alertmanager/template"
	"github.com/stretchr/testify/require"
	"github.com/xujiahua/alertmanager-webhook-feishu/feishu"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestServer_hook(t *testing.T) {
	bots := make(map[string]feishu.IBot)
	bots["test"] = &feishu.FakeBot{}
	s := New(bots)

	tt := []struct {
		group      string
		shouldPass bool
	}{
		{"test", true},
		{"fail", false},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/hook/%s", tc.group)
		req, err := http.NewRequest("POST", path, newBody())
		require.Nil(t, err)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		// Need to create a router that we can pass the request through so that the vars will be added to the context
		router := mux.NewRouter()
		router.HandleFunc("/hook/{group}", s.hook)
		router.ServeHTTP(rr, req)

		require.Equal(t, tc.shouldPass, rr.Code == http.StatusOK)
	}
}

// test real server
func TestServer(t *testing.T) {
	tt := []struct {
		group      string
		shouldPass bool
	}{
		{"webhook", true},
		{"webhook_mention_all", true},
		{"webhook_mention_openids", true},
		{"webhook_mention_emails", true},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("http://localhost:8000/hook/%s", tc.group)
		req, err := http.NewRequest("POST", path, newBody())
		require.Nil(t, err)
		req.Header.Set("Content-Type", "application/json")

		do, err := http.DefaultClient.Do(req)
		require.Nil(t, err)

		require.Equal(t, tc.shouldPass, do.StatusCode == http.StatusOK)
	}
}
func newBody() io.Reader {
	bs, err := json.Marshal(newAlerts())
	if err != nil {
		panic(err)
	}
	return bytes.NewBuffer(bs)
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
