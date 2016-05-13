package main

import (
	docker "github.com/samalba/dockerclient"
	"strings"
)

type policy struct {
	Name        string
	Criteria    string
	Frequency   string
	Probability string
	Image       string
	Faults      string
}

func PolicyFromInjector(i injector) policy {
	return policy{
		Criteria:    i.Criteria,
		Frequency:   i.Frequency,
		Probability: i.Probability,
		Image:       i.Image,
		Faults:      strings.Join(i.Faults, ",")}
}

type criteriaMap map[string]imageMap
type imageMap map[string]faultMap
type faultMap map[string]frequencyMap
type frequencyMap map[string]probabilityMap
type probabilityMap map[string]policy

func PoliciesFromContainers(cs []docker.Container) []policy {
	policies := criteriaMap{}
	r := []policy{}

	for _, v := range cs {
		i := InjectorFromContainer(v)
		p := PolicyFromInjector(i)

		if _, ok := policies[p.Criteria]; !ok {
			policies[p.Criteria] = imageMap{ p.Image: faultMap{p.Faults: frequencyMap{p.Frequency: probabilityMap{p.Probability: p}}}}
			r = append(r, p)
		}

		if _, ok := policies[p.Criteria][p.Image]; !ok {
			policies[p.Criteria][p.Image] = faultMap{p.Faults: frequencyMap{p.Frequency: probabilityMap{p.Probability: p}}}
			r = append(r, p)
		}

		if _, ok := policies[p.Criteria][p.Image][p.Faults]; !ok {
			policies[p.Criteria][p.Image][p.Faults] = frequencyMap{p.Frequency: probabilityMap{p.Probability: p}}
			r = append(r, p)
		}

		if _, ok := policies[p.Criteria][p.Image][p.Faults][p.Frequency]; !ok {
			policies[p.Criteria][p.Image][p.Faults][p.Frequency] = probabilityMap{p.Probability: p}
			r = append(r, p)
		}

		if _, ok := policies[p.Criteria][p.Image][p.Faults][p.Frequency][p.Probability]; !ok {
			policies[p.Criteria][p.Image][p.Faults][p.Frequency][p.Probability] = p
			r = append(r, p)
		}
	}

	return r
}
