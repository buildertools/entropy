package main

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	docker "github.com/samalba/dockerclient"
	"net/http"
)

type Handler func(c *context, w http.ResponseWriter, r *http.Request)
type context struct {
	Target string
}

func ping(c *context, w http.ResponseWriter, r *http.Request) {
	w.Write([]byte{'O', 'K'})
}

func version(c *context, w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(VERSION))
}

func info(c *context, w http.ResponseWriter, r *http.Request) {
	i := struct {
		Version string
	}{
		Version: VERSION,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(i); err != nil {
		panic(err)
	}
}

func list(c *context, w http.ResponseWriter, r *http.Request) {
	client, err := docker.NewDockerClient(c.Target, nil)
	if err != nil {
		panic(err)
	}

	containers, err := client.ListContainers(true, false, fmt.Sprintf("{\"label\":[\"%s\"]}", AGENT_LABEL))
	if err != nil {
		panic(err)
	}

	l := []injector{}
	for _, v := range containers {
		l = append(l, FromContainer(v))
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(l); err != nil {
		panic(err)
	}
}

func create(c *context, w http.ResponseWriter, r *http.Request) {

	i := injector{Name: "4gfewg43"}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "{%q:%q}", "Name", i.Name)
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
