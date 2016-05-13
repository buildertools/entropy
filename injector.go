package main

import (
	docker "github.com/samalba/dockerclient"
	"strings"
)

const (
	FREQUENCY_LABEL   = "buildertools.entropy.frequency"
	PROBABILITY_LABEL = "buildertools.entropy.probability"
	FAULTS_LABEL      = "buildertools.entropy.faults"
	TARGET_LABEL      = "buildertools.entropy.target"
	CRITERIA_LABEL    = "buildertools.entropy.criteria"
	AGENT_LABEL       = "buildertools.entropy.agent"
)

type injector struct {
	Name        string
	TargetCID   string
	Criteria    string
	Frequency   string
	Probability string
	Image       string
	Faults      []string
	Status      string
}

var faults = []string{
	"partition",
	"loss",
	"latency",
	"reordering",
	"pause",
}

func InjectorFromContainer(i docker.Container) injector {
	r := injector{Name: i.Id, Image: i.Image}
	r.TargetCID = i.Labels[TARGET_LABEL]
	r.Criteria = i.Labels[CRITERIA_LABEL]
	r.Frequency = i.Labels[FREQUENCY_LABEL]
	r.Probability = i.Labels[PROBABILITY_LABEL]
	r.Faults = strings.Split(i.Labels[FAULTS_LABEL], ",")
	r.Status = i.Status
	return r
}
