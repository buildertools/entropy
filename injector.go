package main

import (
	"errors"
	docker "github.com/samalba/dockerclient"
	"strings"
)

const (
	FREQUENCY_LABEL   = "buildertools.entropy.frequency"
	PROBABILITY_LABEL = "buildertools.entropy.probability"
	FAILURES_LABEL    = "buildertools.entropy.failures"
	TARGET_LABEL      = "buildertools.entropy.target"
	CRITERIA_LABEL    = "buildertools.entropy.criteria"
	AGENT_LABEL       = "buildertools.entropy.agent"
)

type injector struct {
	Name        string
	Policy      string
	CID         string
	TargetCID   string
	Criteria    string
	Frequency   string
	Probability string
	Image       string
	Failures    []string
	Status      string
	Running     bool
}

var failures = []string{
	"partition",
	"loss",
	"latency",
	"reordering",
	"pause",
}

func InjectorFromContainer(i docker.Container) injector {
	r := injector{Name: strings.Join(i.Names, ", "), CID: i.Id, Image: i.Image}
	r.Policy = i.Labels[AGENT_LABEL]
	r.TargetCID = i.Labels[TARGET_LABEL]
	r.Criteria = i.Labels[CRITERIA_LABEL]
	r.Frequency = i.Labels[FREQUENCY_LABEL]
	r.Probability = i.Labels[PROBABILITY_LABEL]
	r.Failures = strings.Split(i.Labels[FAILURES_LABEL], ",")
	r.Status = i.Status
	return r
}

func InjectorFromContainerInfo(i docker.ContainerInfo) injector {
	r := injector{Name: i.Name, CID: i.Id, Image: i.Image}
	r.Policy = i.Config.Labels[AGENT_LABEL]
	r.TargetCID = i.Config.Labels[TARGET_LABEL]
	r.Criteria = i.Config.Labels[CRITERIA_LABEL]
	r.Frequency = i.Config.Labels[FREQUENCY_LABEL]
	r.Probability = i.Config.Labels[PROBABILITY_LABEL]
	r.Failures = strings.Split(i.Config.Labels[FAILURES_LABEL], ",")
	r.Status = i.State.String()
	r.Running = i.State.Running
	return r
}

func InjectorFromEvent(e docker.Event) (injector, error) {
	if e.Type != "container" {
		return injector{}, errors.New("The provided event is not a container type.")
	}
	// determine if the container is an injector
	if _, present := e.Actor.Attributes[AGENT_LABEL]; !present {
		return injector{}, errors.New("The provided event is not for an injector.")
	}
	return injector{
		Name:        e.Actor.Attributes["name"],
		Policy:      e.Actor.Attributes[AGENT_LABEL],
		Image:       e.Actor.Attributes["image"],
		CID:         e.Actor.ID,
		TargetCID:   e.Actor.Attributes[TARGET_LABEL],
		Criteria:    e.Actor.Attributes[CRITERIA_LABEL],
		Frequency:   e.Actor.Attributes[FREQUENCY_LABEL],
		Probability: e.Actor.Attributes[PROBABILITY_LABEL],
		Failures:    strings.Split(e.Actor.Attributes[FAILURES_LABEL], ","),
		Status:      e.Status}, nil
}
