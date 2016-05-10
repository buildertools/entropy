package main

import (
	docker "github.com/samalba/dockerclient"
	"strings"
)

const (
	FREQUENCY_LABEL   = "buildertools.entropy.frequency"
	PROBABILITY_LABEL = "buildertools.entropy.probability"
	FAULTS_LABEL      = "buildertools.entropy.faults"
	AGENT_LABEL       = "buildertools.entropy.agent"
)

type injector struct {
	Name        string
	Frequency   string
	Probability string
	Image       string
	Faults      []string
}

var faults = []string{
	"partition",
	"loss",
	"latency",
	"reordering",
	"pause",
}

func FromContainer(i docker.Container) injector {
	r := injector{Name: i.Id, Image: i.Image}
	r.Frequency = i.Labels[FREQUENCY_LABEL]
	r.Probability = i.Labels[PROBABILITY_LABEL]
	r.Faults = strings.Split(i.Labels[FAULTS_LABEL], ",")
	return r
}
