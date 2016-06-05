package main

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var (
	noUnix   = errors.New(`Unix connections not currently supported.`)
	noTarget = errors.New(`Unspecified client target. Please set the host parameter.`)
)

func ClientCreatePolicy(c *cli.Context) error {
	unix, target, err := processClientTarget(c)
	if err != nil {
		return err
	}

	if unix {
		return noUnix
	} else {
		name := c.String(`name`)
		image := c.String(`image`)
		criteria := c.String(`criteria`)
		frequency := c.String(`frequency`)
		probability := c.String(`probability`)
		failure := c.String(`failure`)

		resp, err := http.PostForm(fmt.Sprintf("http://%s/policy/", target), url.Values{
			`name`:        {name},
			`image`:       {image},
			`criteria`:    {criteria},
			`frequency`:   {frequency},
			`probability`: {probability},
			`failures`:    {failure},
		})
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		fmt.Printf("%s", body)
	}
	return nil
}

func RequestPolicyList(c *cli.Context) error {
	unix, target, err := processClientTarget(c)
	if err != nil {
		return err
	}

	if unix {
		return noUnix
	} else {
		resp, err := http.Get(fmt.Sprintf("http://%s/policy/", target))
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		fmt.Printf("%s", body)
	}
	return nil
}

func RequestInjectorList(c *cli.Context) error {
	unix, target, err := processClientTarget(c)
	if err != nil {
		return err
	}

	if unix {
		return noUnix
	} else {
		resp, err := http.Get(fmt.Sprintf("http://%s/injector/", target))
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		fmt.Printf("%s", body)
	}
	return nil
}

func DeletePolicy(c *cli.Context) error {
	unix, target, err := processClientTarget(c)
	if err != nil {
		return err
	}

	if unix {
		return noUnix
	} else {
		var name string
		if c.NArg() > 0 {
			name = c.Args()[0]
		}

		client := &http.Client{}

		req, err := http.NewRequest("DELETE", fmt.Sprintf("http://%s/policy/%s", target, name), nil)
		if err != nil {
			log.Fatal(err)
		}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		fmt.Printf("%s", body)
	}
	return nil
}

func processClientTarget(c *cli.Context) (isUnix bool, target string, err error) {
	for _, v := range c.GlobalStringSlice(`host`) {
		if strings.HasPrefix(v, "tcp://") {
			isUnix = false
			target = strings.TrimPrefix(v, "tcp://")
			err = nil
		}
		if strings.HasPrefix(v, "unix://") {
			isUnix = true
			target = strings.TrimPrefix(v, "unix://")
			err = nil
		}
	}
	return
}
