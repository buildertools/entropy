package main

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	docker "github.com/samalba/dockerclient"
	"net/http"
)

type Handler func(c *context, w http.ResponseWriter, r *http.Request)
type context struct {
	Target string
	Image  string
	Gin    *gin.Context
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

	p := PoliciesFromContainers(containers)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(p); err != nil {
		panic(err)
	}
}

func lsi(c *context, w http.ResponseWriter, r *http.Request) {
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
		l = append(l, InjectorFromContainer(v))
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(l); err != nil {
		panic(err)
	}
}

func create(c *context, w http.ResponseWriter, r *http.Request) {
	payload := injector{
		Criteria:    c.Gin.PostForm(`criteria`),
		Frequency:   c.Gin.PostForm(`frequency`),
		Probability: c.Gin.PostForm(`probability`),
		Faults:      []string{c.Gin.PostForm(`fault`)}}

	client, err := docker.NewDockerClient(c.Target, nil)
	if err != nil {
		panic(err)
	}

	containerConfig := &docker.ContainerConfig{
		Image:       c.Image,
		Cmd:         []string{"sleep", "40"},
		AttachStdin: false,
		Tty:         false,
		Labels: map[string]string{
			AGENT_LABEL:       "",
			FREQUENCY_LABEL:   payload.Frequency,
			PROBABILITY_LABEL: payload.Probability,
			FAULTS_LABEL:      payload.Faults[0],
			TARGET_LABEL:      "",
			CRITERIA_LABEL:    payload.Criteria},
	}
	containerId, err := client.CreateContainer(containerConfig, GenerateName(), nil)
	if err != nil {
		httpError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = client.StartContainer(containerId, &docker.HostConfig{})
	if err != nil {
		httpError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	i := injector{Name: containerId}

	// do something more here

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(struct{ Name string }{Name: i.Name}); err != nil {
		panic(err)
	}
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
