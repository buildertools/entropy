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
	payload := policy{
		Name:        c.Gin.PostForm(`name`),
		Image:       c.Gin.PostForm(`image`),
		Criteria:    c.Gin.PostForm(`criteria`),
		Frequency:   c.Gin.PostForm(`frequency`),
		Probability: c.Gin.PostForm(`probability`),
		Failures:    c.Gin.PostForm(`failures`)}

	client, err := docker.NewDockerClient(c.Target, nil)
	if err != nil {
		panic(err)
	}

	// generate a policy name
	if payload.Name == "" {
		var name string
		for i := 0; i < 100; i++ {
			name = GenerateName()
			if _, err = client.InspectContainer(name); err != nil {
				break
			}
		}
		payload.Name = name
	}

	// determine which containers match the criteria
	// and create an injector for each using a _N suffix
	containers, err := client.ListContainers(true, false, fmt.Sprintf("{\"label\":[\"%s\"]}", payload.Criteria))
	if err != nil {
		httpError(w, err.Error(), http.StatusInternalServerError)
		return
	} else if len(containers) == 0 {
		// TODO: this should be a different reesponse code
		log.Println("Recieved create request that failed to match containers.")
		w.WriteHeader(http.StatusCreated)
		return
	}

	for _, target := range containers {
		err = CreateInjector(payload, target.Id, c.Target)
		if err != nil {
			httpError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	i := injector{Name: payload.Name}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(struct{ Name string }{Name: i.Name}); err != nil {
		panic(err)
	}
}

func remove(c *context, w http.ResponseWriter, r *http.Request) {
	name := c.Gin.Param(`name`)
	log.Printf("Removing policy with name: %v\n", name)
	if name == `` {
		httpError(w, `No policy name specified.`, http.StatusNotFound)
	}
	client, err := docker.NewDockerClient(c.Target, nil)
	if err != nil {
		httpError(w, err.Error(), http.StatusInternalServerError)
	}
	is, err := GetInjectorsForPolicy(client, name)
	for _, i := range is {
		log.Printf("Deleting injector: %s", i.Name)
		err = client.RemoveContainer(i.CID, true, true)
		if err != nil {
			httpError(w, err.Error(), http.StatusInternalServerError)
		}
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(struct{ Ok bool }{Ok: true}); err != nil {
		panic(err)
	}
}

func CreateInjector(p policy, target, endpoint string) error {
	sTarget := target[0:12]
	client, err := docker.NewDockerClient(endpoint, nil)
	if err != nil {
		return err
	}

	hostConfig := docker.HostConfig{
		CapAdd:      []string{"NET_ADMIN"},
		IpcMode:     fmt.Sprintf(`container:%s`, target),
		NetworkMode: fmt.Sprintf(`container:%s`, target),
	}
	containerConfig := &docker.ContainerConfig{
		HostConfig:  hostConfig,
		Image:       p.Image,
		AttachStdin: false,
		Tty:         false,
		Labels: map[string]string{
			AGENT_LABEL:       p.Name,
			FREQUENCY_LABEL:   p.Frequency,
			PROBABILITY_LABEL: p.Probability,
			FAILURES_LABEL:    p.Failures,
			TARGET_LABEL:      target,
			CRITERIA_LABEL:    p.Criteria},
		Env: []string{
			fmt.Sprintf(`ENTROPY_FREQUENCY=%s`, p.Frequency),
			fmt.Sprintf(`ENTROPY_PROBABILITY=%s`, p.Probability),
			fmt.Sprintf(`ENTROPY_FAILURES=%s`, p.Failures),
			fmt.Sprintf(`ENTROPY_TARGET=%s`, target),
		},
	}
	containerId, err := client.CreateContainer(containerConfig, fmt.Sprintf("%s_%v", p.Name, sTarget), nil)
	if err != nil {
		return err
	}

	err = client.StartContainer(containerId, &docker.HostConfig{
		CapAdd:      []string{"NET_ADMIN"},
		IpcMode:     fmt.Sprintf(`container:%s`, target),
		NetworkMode: fmt.Sprintf(`container:%s`, target),
	})
	if err != nil {
		return err
	}
	return nil
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
