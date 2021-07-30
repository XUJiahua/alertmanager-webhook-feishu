package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/xujiahua/alertmanager-webhook-feishu/feishu"
	"github.com/xujiahua/alertmanager-webhook-feishu/model"
	"net/http"
	"time"
)

type Server struct {
	bots map[string]feishu.IBot
}

func New(bots map[string]feishu.IBot) *Server {
	s := &Server{
		bots: bots,
	}
	return s
}

func (s Server) hook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	group := vars["group"]
	bot, ok := s.bots[group]
	if !ok {
		logrus.Errorf("group not found: %s", group)
		http.Error(w, "group not found", http.StatusBadRequest)
		return
	}

	var alerts model.WebhookMessage
	err := json.NewDecoder(r.Body).Decode(&alerts)
	if err != nil {
		logrus.Errorf("cannot parse content, %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logrus.Debug(alerts)

	err = bot.Send(&alerts)
	if err != nil {
		logrus.Errorf("cannot send alerts, %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = fmt.Fprintf(w, "ok")
}

func (s Server) health(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (s Server) reload(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (s Server) Start(address string) error {
	r := mux.NewRouter()
	r.HandleFunc("/hook/{group}", s.hook).Methods("POST")

	// management etc...
	sr := r.PathPrefix("/-").Subrouter()
	sr.HandleFunc("/healthz", s.health).Methods("GET")
	sr.HandleFunc("/reload", s.health).Methods("GET")

	// prometheus
	r.Handle("/metrics", promhttp.Handler()).Methods("GET")

	srv := &http.Server{
		Handler:      r,
		Addr:         address,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return srv.ListenAndServe()
}
