package main

import (
	"fmt"
	docker "github.com/samalba/dockerclient"
	"strings"
)

type policy struct {
	Name        string
	Criteria    string
	Frequency   string
	Probability string
	Image       string
	Failures    string
	Injectors   int
}

func PolicyFromInjector(i injector) policy {
	return policy{
		Name:        i.Policy,
		Criteria:    i.Criteria,
		Frequency:   i.Frequency,
		Probability: i.Probability,
		Image:       i.Image,
		Failures:    strings.Join(i.Failures, ",")}
}

func GetInjectorsForPolicy(client *docker.DockerClient, name string) ([]injector, error) {
	is := []injector{}
	cs, err := client.ListContainers(true, false, fmt.Sprintf(`{"label":["%s=%s"]}`, AGENT_LABEL, name))
	if err != nil {
		return is, err
	}
	for _, c := range cs {
		is = append(is, InjectorFromContainer(c))
	}
	return is, nil
}

type criteriaMap map[string]imageMap
type imageMap map[string]failureMap
type failureMap map[string]frequencyMap
type frequencyMap map[string]probabilityMap
type probabilityMap map[string]policy

func PoliciesFromContainers(cs []docker.Container) []policy {
	policies := criteriaMap{}
	r := []policy{}

	for _, v := range cs {
		i := InjectorFromContainer(v)
		p := PolicyFromInjector(i)

		if _, ok := policies[p.Criteria]; !ok {
			p.Injectors = 1
			policies[p.Criteria] = imageMap{p.Image: failureMap{p.Failures: frequencyMap{p.Frequency: probabilityMap{p.Probability: p}}}}
		} else if _, ok := policies[p.Criteria][p.Image]; !ok {
			p.Injectors = 1
			policies[p.Criteria][p.Image] = failureMap{p.Failures: frequencyMap{p.Frequency: probabilityMap{p.Probability: p}}}
		} else if _, ok := policies[p.Criteria][p.Image][p.Failures]; !ok {
			p.Injectors = 1
			policies[p.Criteria][p.Image][p.Failures] = frequencyMap{p.Frequency: probabilityMap{p.Probability: p}}
		} else if _, ok := policies[p.Criteria][p.Image][p.Failures][p.Frequency]; !ok {
			p.Injectors = 1
			policies[p.Criteria][p.Image][p.Failures][p.Frequency] = probabilityMap{p.Probability: p}
		} else if _, ok := policies[p.Criteria][p.Image][p.Failures][p.Frequency][p.Probability]; !ok {
			p.Injectors = 1
			policies[p.Criteria][p.Image][p.Failures][p.Frequency][p.Probability] = p
		} else {
			p.Injectors = 1 + policies[p.Criteria][p.Image][p.Failures][p.Frequency][p.Probability].Injectors
			policies[p.Criteria][p.Image][p.Failures][p.Frequency][p.Probability] = p
		}
	}

	for c, _ := range policies {
		for i, _ := range policies[c] {
			for f, _ := range policies[c][i] {
				for fr, _ := range policies[c][i][f] {
					for p, _ := range policies[c][i][f][fr] {
						r = append(r, policies[c][i][f][fr][p])
					}
				}
			}
		}
	}

	return r
}
