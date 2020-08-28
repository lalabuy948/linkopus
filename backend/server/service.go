package server

import (
	"fmt"
	"log"

	"github.com/lalabuy948/linkopus/backend/di"

	"github.com/valyala/fasthttp"
)

type server struct {
	container *di.Container
}

// NewServer injects services container into linkopus server and returns it.
func NewServer(c *di.Container) *server {
	return &server{c}
}

// Start stars server.
func (s *server) Start(port *string, compress *bool) {
	requestHandler := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/api/v1/link":
			s.linkHandler(ctx)
		case "/api/v1/link/stats":
			s.statsHandler(ctx)
		case "/api/v1/link/stats/top":
			s.topHandler(ctx)
		default:
			s.redirectHandler(ctx)
		}
	}

	h := requestHandler
	if *compress {
		h = fasthttp.CompressHandler(requestHandler)
	}

	fmt.Println("starting server on ", *port)
	if err := fasthttp.ListenAndServe(":"+*port, h); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}
