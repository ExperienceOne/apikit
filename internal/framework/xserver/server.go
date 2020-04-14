package xserver

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"encoding/json"
	"runtime/debug"

	"github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/fault"
	"github.com/ExperienceOne/apikit/internal/framework/xhttperror"
)

type (
	ServerOpts struct {
		ErrorHandler ErrorHandler
		Middleware   []Middleware
		OnStart      func(router *routing.Router)
	}

	Middleware struct {
		Handler routing.Handler
		After   bool
	}

	RouteDescription struct {
		Path       string
		Handler    routing.Handler
		Middleware []Middleware
		Method     string
	}

	Server struct {
		ErrorLogger func(v ...interface{})
		server      *http.Server
		after       []routing.Handler
		before      []routing.Handler
		SwaggerSpec string
		Router      *routing.Router
		OnStart     func(router *routing.Router)
	}
)

type ErrorHandler func(v ...interface{})

// NewServer initializes a new Server instance with a middleware handler.
// middleware represents middleware handler that will be executed before  or after the actual
// handler of each server endpoint.
func NewServer(opts *ServerOpts) *Server {

	if opts == nil {
		return &Server{
			ErrorLogger: func(v ...interface{}) {},
		}
	}

	if opts.ErrorHandler == nil {
		opts.ErrorHandler = func(v ...interface{}) {}
	}

	server := &Server{
		ErrorLogger: opts.ErrorHandler,
	}

	if opts.OnStart != nil {
		server.OnStart = opts.OnStart
	}

	if len(opts.Middleware) != 0 {
		before := make([]routing.Handler, 0)
		after := make([]routing.Handler, 0)

		for _, m := range opts.Middleware {
			if m.After {
				after = append(after, m.Handler)
			} else {
				before = append(before, m.Handler)
			}
		}

		server.after = after
		server.before = before
	}

	return server
}

func (server *Server) makeRouter(routes []RouteDescription) (*routing.Router, error) {

	router := routing.New()
	router.UseEscapedPath = true

	logError := func(format string, a ...interface{}) {
		msg := fmt.Sprintf(format, a...)
		server.ErrorLogger(msg)
	}

	var beforeStack []routing.Handler
	beforeStack = append(beforeStack, errorHandler(logError))
	if server.before != nil {
		beforeStack = append(beforeStack, server.before...)
	}
	router.Use(beforeStack...)

	var afterStack []routing.Handler
	if server.after != nil {
		afterStack = append(afterStack, server.after...)
	}

	router.Get("/spec", func(c *routing.Context) error {
		return c.Write(server.SwaggerSpec)
	})

	for _, route := range routes {

		var before, after []routing.Handler

		if route.Middleware != nil {
			for _, m := range route.Middleware {
				if m.After {
					after = append(after, m.Handler)
				} else {
					before = append(before, m.Handler)
				}
			}
		}

		var handler []routing.Handler
		handler = append(handler, before...)
		handler = append(handler, route.Handler)
		handler = append(handler, after...)
		handler = append(handler, afterStack...)

		router.To(route.Method, route.Path, handler...)
	}

	return router, nil
}

func (server *Server) Start(port int, routes []RouteDescription) error {

	router, err := server.makeRouter(routes)
	if err != nil {
		return err
	}
	if server.OnStart != nil {
		server.OnStart(router)
	}
	server.Router = router

	httpServer := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: router,
	}
	server.server = httpServer

	return httpServer.ListenAndServe()
}

func (server *Server) Stop() error {

	if server.server != nil {
		deadline, _ := context.WithTimeout(context.TODO(), 30*time.Second)
		return server.server.Shutdown(deadline)
	}

	return nil
}

//errorHandler overwrites global error handling of router
func errorHandler(logf fault.LogFunc) func(c *routing.Context) (err error) {
	return func(c *routing.Context) (err error) {
		defer func() {

			if e := recover(); e != nil {

				if logf != nil {
					logf("recovered from panic: %v", string(debug.Stack()))
				}
				c.Response.WriteHeader(http.StatusInternalServerError)
				err = nil
				c.Abort()

			} else if err != nil {

				switch errType := err.(type) {
				case *xhttperror.HttpJsonError:
					c.Response.Header()["Content-Type"] = []string{"application/json"}
					c.Response.WriteHeader(errType.StatusCode())
					if e := json.NewEncoder(c.Response).Encode(errType.Message); e != nil && logf != nil {
						logf("failed to write error message: %v", errType.Message)
					}
				case *xhttperror.HttpCodeError:
					c.Response.Header()["Content-Type"] = []string{""}
					c.Response.WriteHeader(errType.StatusCode())
				case routing.HTTPError:
					c.Response.Header()["Content-Type"] = []string{"text/plain; charset=utf-8"}
					c.Response.WriteHeader(errType.StatusCode())
					if _, e := c.Response.Write([]byte(errType.Error())); e != nil {
						logf("failed to write error message: %v", errType.Error())
					}
				}

				err = nil
				c.Abort() // call Abort() to stop further middleware processing
			}
		}()

		return c.Next()
	}
}
