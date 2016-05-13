package main

import (
	//"errors"
	"fmt"
	"io/ioutil"
	//log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"net/http"
)

func requestPolicyList(c *cli.Context) {
	resp, err := http.Get(fmt.Sprintf("http://%s/policy/", c.GlobalString("host")))
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	fmt.Printf("%v", body)
}

func requestInjectorList(c *cli.Context) {
	resp, err := http.Get(fmt.Sprintf("http://%s/injector/", c.GlobalString("host")))
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	fmt.Printf("%v", body)
}
