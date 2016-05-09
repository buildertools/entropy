package main

import (
	"github.com/gin-gonic/gin"
	"os"
	"os/signal"
)

type manager struct {
	Unix string
	Tcp  string
}

func (m *manager) Start() {
	// Create a synchronization channel
	c := make(chan bool)

	registerInterruptHandler(func() {
		c <- true
	})

	// Start the UNIX server if configured
	if m.Unix != "" {
		go func() {
			router := createRoutes()
			router.RunUnix(m.Unix)
			c <- true
		}()
	}

	// Start the TCP server if configured
	if m.Tcp != "" {
		go func() {
			router := createRoutes()
			router.Run(m.Tcp)
			c <- true
		}()
	}
	// Block until one of the servers stop or the manager is interrupted
	<-c
	// TODO: Do cleanup

}

func registerInterruptHandler(f func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			f()
		}
	}()
}

func createRoutes() *gin.Engine {
	router := gin.Default()
	for _, f := range API {
		// Wrap with a gin.HandlerFunc
		wrap := func(c *gin.Context) {
			ctx := &context{}
			f.Handler(ctx, c.Writer, c.Request)
		}

		switch {
		case "GET" == f.Method:
			router.GET(f.Route, wrap)
		case "POST" == f.Method:
			router.POST(f.Route, wrap)
		case "DELETE" == f.Method:
			router.DELETE(f.Route, wrap)
		case "PUT" == f.Method:
			router.PUT(f.Route, wrap)
		}
	}
	return router
}
