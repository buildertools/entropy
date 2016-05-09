package main

type feature struct {
	Path         string
	PathTemplate string
	Method       string
	Route        string
	Handler      Handler
}

var API = map[string]feature{
	"ping":     {PathTemplate: "/_ping", Method: "GET", Route: "/_ping", Handler: handlerNotYetImplemented},
	"events":   {PathTemplate: "/events", Method: "GET", Route: "/events", Handler: handlerNotYetImplemented},
	"info":     {PathTemplate: "/info", Method: "GET", Route: "/info", Handler: handlerNotYetImplemented},
	"version":  {PathTemplate: "/version", Method: "GET", Route: "/version", Handler: handlerNotYetImplemented},
	"list":     {PathTemplate: "/injectors/", Method: "GET", Route: "/injectors/", Handler: handlerNotYetImplemented},
	"injector": {PathTemplate: "/injectors/{{.Cid}}/json", Method: "GET", Route: "/injectors/:name/json", Handler: handlerNotYetImplemented},
	"logs":     {PathTemplate: "/injectors/{{.Cid}}/logs", Method: "GET", Route: "/injectors/:name/logs", Handler: handlerNotYetImplemented},
	"create":   {PathTemplate: "/injectors/", Method: "POST", Route: "/injectors/", Handler: handlerNotYetImplemented},
	"kill":     {PathTemplate: "/injectors/{{.Cid}}/kill", Method: "POST", Route: "/injectors/:name/kill", Handler: handlerNotYetImplemented},
	"pause":    {PathTemplate: "/injectors/{{.Cid}}/pause", Method: "POST", Route: "/injectors/:name/pause", Handler: handlerNotYetImplemented},
	"unpause":  {PathTemplate: "/injectors/{{.Cid}}/unpause", Method: "POST", Route: "/injectors/:name/unpause", Handler: handlerNotYetImplemented},
	"restart":  {PathTemplate: "/injectors/{{.Cid}}/restart", Method: "POST", Route: "/injectors/:name/restart", Handler: handlerNotYetImplemented},
	"start":    {PathTemplate: "/injectors/{{.Cid}}/start", Method: "POST", Route: "/injectors/:name/start", Handler: handlerNotYetImplemented},
	"stop":     {PathTemplate: "/injectors/{{.Cid}}/stop", Method: "POST", Route: "/injectors/:name/stop", Handler: handlerNotYetImplemented},
	"update":   {PathTemplate: "/injectors/{{.Cid}}/update", Method: "POST", Route: "/injectors/:name/update", Handler: handlerNotYetImplemented},
	"delete":   {PathTemplate: "/injectors/{{.Cid}}", Method: "DELETE", Route: "/injectors/:name", Handler: handlerNotYetImplemented},
}
