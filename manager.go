package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"os"
	"os/signal"
)

type manager struct {
	Unix   string
	Tcp    string
	Target string
	Image  string
	Events *EventContext
}

func (m *manager) Start() {
	// Create a synchronization channel
	c := make(chan bool)

	registerInterruptHandler(func() {
		c <- true
	})

	startEventLogger(m.Events)

	// Start the UNIX server if configured
	if m.Unix != "" {
		go func() {
			router := m.createRoutes()
			router.RunUnix(m.Unix)
			c <- true
		}()
	}

	// Start the TCP server if configured
	if m.Tcp != "" {
		go func() {
			router := m.createRoutes()
			router.Run(m.Tcp)
			c <- true
		}()
	}
	// Block until one of the servers stop or the manager is interrupted
	<-c
	// TODO: Do cleanup

}

func startEventLogger(ec *EventContext) {
	o := make(chan event)
	ec.registerEventObserver(o)
	go func(c chan event) {
		for {
			log.Info(<-o)
		}
	}(o)
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

func (m *manager) createRoutes() *gin.Engine {
	router := gin.Default()
	for _, f := range API {
		lf := f

		// Wrap with a gin.HandlerFunc
		wrap := func(c *gin.Context) {
			ctx := &context{Target: m.Target, Image: m.Image, Gin: c}
			lf.Handler(ctx, c.Writer, c.Request)
		}

		switch {
		case "GET" == lf.Method:
			router.GET(lf.Route, wrap)
		case "POST" == lf.Method:
			router.POST(lf.Route, wrap)
		case "DELETE" == lf.Method:
			router.DELETE(lf.Route, wrap)
		case "PUT" == lf.Method:
			router.PUT(lf.Route, wrap)
		}
	}
	return router
}
