package server

import (
	"fmt"

	"github.com/fagongzi/gateway-plugin-consul/pkg/conf"
	"github.com/fagongzi/gateway-plugin-consul/pkg/consul"
	"github.com/labstack/echo"
	sd "github.com/labstack/echo/engine/standard"
	mw "github.com/labstack/echo/middleware"
)

// Server http interface server
type Server struct {
	cnf    *conf.Conf
	consul *consul.Backend
	e      *echo.Echo
}

// NewServer create a plugin server
func NewServer(cnf *conf.Conf) *Server {
	server := &Server{
		cnf: cnf,
		e:   echo.New(),
	}

	c, err := consul.NewBackend(cnf)
	if err != nil {
		panic(err)
	}

	server.consul = c

	server.initHTTPServer()

	return server
}

func (server *Server) initHTTPServer() {
	server.e.Use(mw.Recover())
	server.e.Use(mw.Logger())
	server.initAPIRoute()
}

// Start start the admin server
func (server *Server) Start() {
	fmt.Printf("start at %s\n", server.cnf.Addr)
	server.e.Run(sd.New(server.cnf.Addr))
}
