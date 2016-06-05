package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	docker "github.com/samalba/dockerclient"
)

type event struct {
	DateTime  string
	Policy    policy
	Type      string
	TargetCID string
}

func (e event) String() string {
	return fmt.Sprintf("[%v] %v - %v: %v", e.DateTime, e.Policy.Name, e.Type, e.TargetCID)
}

type EntropyEventContext struct {
	observers []chan event
}

type DockerEventContext struct {
	observers []chan docker.Event
	started   bool
	target    string
}

func (ec *DockerEventContext) Register(c chan docker.Event) {
	ec.observers = append(ec.observers, c)
}

func (ec *EntropyEventContext) Register(c chan event) {
	ec.observers = append(ec.observers, c)
}

// send a copy of the event to every registered observer
func (ec *DockerEventContext) notify(e docker.Event) {
	for _, o := range ec.observers {
		go func(o chan docker.Event, e docker.Event) {
			// TODO: add timeout handler to remove observer and terminate goroutine
			o <- e
		}(o, e)
	}
}

// send a copy of the event to every registered observer
func (ec *EntropyEventContext) notify(e event) {
	for _, o := range ec.observers {
		go func(o chan event, e event) {
			// TODO: add timeout handler to remove observer and terminate goroutine
			o <- e
		}(o, e)
	}
}

func startDockerEventLogger(ec *DockerEventContext) {
	o := make(chan docker.Event)
	ec.Register(o)
	go func(c chan docker.Event) {
		for v := range o {
			log.Debug(v)
		}
	}(o)
}

func startEntropyEventLogger(ec *EntropyEventContext) {
	o := make(chan event)
	ec.Register(o)
	go func(c chan event) {
		for v := range o {
			log.Info(v)
		}
	}(o)
}

func (ec *DockerEventContext) start(t string) {
	if ec.started {
		return
	}
	ec.target = t
	client, err := docker.NewDockerClient(t, nil)
	if err != nil {
		panic(err)
	}
	echan := make(chan error)
	go func() {
		<-echan
		client.StopAllMonitorEvents()
	}()
	client.StartMonitorEvents(func(de *docker.Event, echan chan error, args ...interface{}) {
		ec.notify(*de)
	}, echan)
}
