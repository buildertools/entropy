package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	docker "github.com/samalba/dockerclient"
	"time"
)

func startPolicyEnforcer(dec *DockerEventContext, eec *EntropyEventContext, endpoint string) {
	d := make(chan docker.Event)
	dec.Register(d)
	go func() {
		log.Info("Starting policy control loop...")
		for de := range d {
			if de.Type != "container" {
				continue
			}
			i, err := InjectorFromEvent(de)
			if err != nil {
				handleContainerEvent(de, eec, endpoint)
			} else {
				handleInjectorEvent(i, eec, endpoint)
			}
		}
	}()
}

func handleInjectorEvent(i injector, ec *EntropyEventContext, endpoint string) {
	log.Debugf("Processing injector event: %s, with status %s", i.Name, i.Status)
	switch i.Status {
	case "start":
		ec.notify(NewInjectorEvent(i, "started"))
	case "die":
		ec.notify(NewInjectorEvent(i, "stop"))
		//		if !isTargetRunning(i.TargetCID, endpoint) {
		//			ec.notify(NewInjectorEvent(i, "stop"))
		//		} else {
		//			// injector failed
		//			ec.notify(NewInjectorEvent(i, "recovering"))
		//			restartInjector(i.CID, endpoint)
		//		}
	case "destroy":
		ec.notify(NewInjectorEvent(i, "destroyed"))
	}
}

func handleContainerEvent(de docker.Event, ec *EntropyEventContext, endpoint string) {
	log.Debugf("Processing container event: %s, with status %s", de.Actor.ID, de.Status)

	switch de.Status {
	case "die":
		injectors := fetchInjectorsForContainer(de.Actor.ID, endpoint, false)
		for _, i := range injectors {
			log.Infof("Stopping injector: %s", i.CID)
			stopInjector(i.CID, endpoint)
		}
	case "start":
		// start existing injectors
		injectors := fetchInjectorsForContainer(de.Actor.ID, endpoint, true)
		for _, i := range injectors {
			if !i.Running {
				log.Infof("Restarting injector: %s", i.CID)
				restartInjector(i.CID, endpoint)
			}
		}
		createInjectorsForContainer(de.Actor.ID, injectors, endpoint)
	case "pause":
	case "unpause":
	case "destroy":
		injectors := fetchInjectorsForContainer(de.Actor.ID, endpoint, true)
		for _, i := range injectors {
			destroyInjector(i.CID, endpoint)
		}
	}
}

func NewInjectorEvent(i injector, t string) event {
	return event{
		DateTime:  time.Now().String(),
		Policy:    PolicyFromInjector(i),
		TargetCID: i.TargetCID,
		Type:      t}
}

func createInjectorsForContainer(id string, is []injector, e string) error {
	client, err := docker.NewDockerClient(e, nil)
	if err != nil {
		log.Fatal(err)
	}
	containers, err := client.ListContainers(true, false, fmt.Sprintf("{\"label\":[\"%s\"]}", AGENT_LABEL))
	if err != nil {
		log.Fatal(err)
	}
	ps := PoliciesFromContainers(containers)

policy:
	for _, p := range ps {
		// for each policy, query the system for its criteria and determine if t matches
		matches, err := client.ListContainers(true, false, fmt.Sprintf(`{"label":["%s"]}`, p.Criteria))
		if err != nil {
			log.Error(err)
			continue
		}

		match := false
		for _, o := range matches {
			if o.Id == id {
				match = true
				break
			}
		}

		if !match {
			continue
		}

		// skip policies that have already been applied
		for _, i := range is {
			if i.Policy == p.Name {
				continue policy
			}
		}
		// apply the policy
		CreateInjector(p, id, e)
	}
	return nil
}

func isTargetRunning(t, e string) bool {
	if t == "" {
		return false
	}
	if e == "" {
		log.Fatal("Unspecified Docker endpoint.")
	}

	client, err := docker.NewDockerClient(e, nil)
	if err != nil {
		log.Fatal(err)
	}

	ci, err := client.InspectContainer(t)
	if err != nil {
		if err == docker.ErrNotFound {
			return false
		}
		log.Fatal(err)
	}

	return ci.State.Running
}

func restartInjector(t, e string) {
	if t == "" {
		log.Error("Unspecified CID.")
		return
	}
	if e == "" {
		log.Error("Unspecified Docker endpoint.")
		return
	}

	client, err := docker.NewDockerClient(e, nil)
	if err != nil {
		log.Error(err)
		return
	}

	err = client.RestartContainer(t, 30)
	if err != nil {
		if err == docker.ErrNotFound {
			log.Info("Unable to recover deleted injector: %s", t)
		} else {
			log.Error(err)
		}
	}
}

func stopInjector(t, e string) {
	if t == "" {
		log.Error("Unspecified CID.")
		return
	}
	if e == "" {
		log.Error("Unspecified Docker endpoint.")
		return
	}

	client, err := docker.NewDockerClient(e, nil)
	if err != nil {
		log.Error(err)
		return
	}

	err = client.StopContainer(t, 30)
	if err != nil {
		if err == docker.ErrNotFound {
			log.Info("Unable to recover deleted injector: %s", t)
		} else {
			log.Error(err)
		}
	}
}

func destroyInjector(t, e string) {
	if t == "" {
		log.Error("Unspecified CID.")
		return
	}
	if e == "" {
		log.Error("Unspecified Docker endpoint.")
		return
	}

	client, err := docker.NewDockerClient(e, nil)
	if err != nil {
		log.Error(err)
		return
	}

	err = client.RemoveContainer(t, true, true)
	if err != nil {
		log.Error(err)
	}
}

func fetchInjectorsForContainer(t, e string, all bool) []injector {
	client, err := docker.NewDockerClient(e, nil)
	if err != nil {
		log.Fatal(err)
	}

	containers, err := client.ListContainers(all, false, fmt.Sprintf("{\"label\":[\"%s=%s\"]}", TARGET_LABEL, t))
	injectors := []injector{}
	for _, c := range containers {
		ci, err := client.InspectContainer(c.Id)
		if err != nil {
			continue
		}
		injectors = append(injectors, InjectorFromContainerInfo(*ci))
	}
	return injectors
}
