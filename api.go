package main

type feature struct {
	Path         string
	PathTemplate string
	Method       string
	Route        string
	Handler      Handler
}

var API = map[string]feature{
	"ping":    {PathTemplate: "/_ping", Method: "GET", Route: "/_ping", Handler: ping},
	"events":  {PathTemplate: "/events", Method: "GET", Route: "/events", Handler: handlerNotYetImplemented},
	"info":    {PathTemplate: "/info", Method: "GET", Route: "/info", Handler: info},
	"version": {PathTemplate: "/version", Method: "GET", Route: "/version", Handler: version},

	"list":   {PathTemplate: "/policy/", Method: "GET", Route: "/policy/", Handler: list},
	"policy": {PathTemplate: "/policy/{{.Cid}}/json", Method: "GET", Route: "/policy/:name/json", Handler: handlerNotYetImplemented},
	"create": {PathTemplate: "/policy/", Method: "POST", Route: "/policy/", Handler: create},
	"kill":   {PathTemplate: "/policy/{{.Cid}}/kill", Method: "POST", Route: "/policy/:name/kill", Handler: handlerNotYetImplemented},
	"start":  {PathTemplate: "/policy/{{.Cid}}/start", Method: "POST", Route: "/policy/:name/start", Handler: handlerNotYetImplemented},
	"stop":   {PathTemplate: "/policy/{{.Cid}}/stop", Method: "POST", Route: "/policy/:name/stop", Handler: handlerNotYetImplemented},
	"update": {PathTemplate: "/policy/{{.Cid}}/update", Method: "POST", Route: "/policy/:name/update", Handler: handlerNotYetImplemented},
	"delete": {PathTemplate: "/policy/{{.Cid}}", Method: "DELETE", Route: "/policy/:name", Handler: remove},

	"lsi":      {PathTemplate: "/injector/", Method: "GET", Route: "/injector/", Handler: lsi},
	"injector": {PathTemplate: "/injector/{{.Cid}}/json", Method: "GET", Route: "/injector/:name/json", Handler: handlerNotYetImplemented},
	"logs":     {PathTemplate: "/injector/{{.Cid}}/logs", Method: "GET", Route: "/injector/:name/logs", Handler: handlerNotYetImplemented},
	"pause":    {PathTemplate: "/injector/{{.Cid}}/pause", Method: "POST", Route: "/injector/:name/pause", Handler: handlerNotYetImplemented},
	"unpause":  {PathTemplate: "/injector/{{.Cid}}/unpause", Method: "POST", Route: "/injector/:name/unpause", Handler: handlerNotYetImplemented},
	"restart":  {PathTemplate: "/injector/{{.Cid}}/restart", Method: "POST", Route: "/injector/:name/restart", Handler: handlerNotYetImplemented},
}
