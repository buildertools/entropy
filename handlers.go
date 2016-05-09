package main

import (
	log "github.com/Sirupsen/logrus"
	"net/http"
)

type Handler func(c *context, w http.ResponseWriter, r *http.Request)
type context struct{}

func ping(c *context, w http.ResponseWriter, r *http.Request) {
	w.Write([]byte{'O', 'K'})
}

func handlerNotYetImplemented(c *context, w http.ResponseWriter, r *http.Request) {
	httpError(w, "Not implemented yet.", http.StatusNotImplemented)
}

func notImplementedHandler(c *context, w http.ResponseWriter, r *http.Request) {
	httpError(w, "Not supported in clustering mode.", http.StatusNotImplemented)
}

func httpError(w http.ResponseWriter, err string, status int) {
	log.WithField("status", status).Errorf("HTTP error: %v", err)
	http.Error(w, err, status)
}
