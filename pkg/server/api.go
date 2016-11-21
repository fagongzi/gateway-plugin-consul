package server

import (
	"net/http"

	"fmt"

	"github.com/fagongzi/gateway/pkg/model"
	"github.com/labstack/echo"
)

func (server *Server) initAPIRoute() {
	server.e.POST("/plugins/service-discovery", server.pluginRegistry())

	server.e.GET("/plugins/service-discovery/clusters", server.getClusters())
	server.e.GET("/plugins/service-discovery/clusters/:name", server.getServers())
}

func (server *Server) pluginRegistry() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	}
}

func (server *Server) getClusters() echo.HandlerFunc {
	return func(c echo.Context) error {
		clusters, err := server.consul.GetClusters()

		if nil != err {
			fmt.Println(err)
			return c.NoContent(http.StatusInternalServerError)
		}

		return c.JSON(http.StatusOK, &model.Clusters{
			Clusters: clusters,
		})
	}
}

func (server *Server) getServers() echo.HandlerFunc {
	return func(c echo.Context) error {
		servers, err := server.consul.GetServers(c.Param("name"))

		if nil != err {
			return c.NoContent(http.StatusInternalServerError)
		}

		return c.JSON(http.StatusOK, &model.Servers{
			Servers: servers,
		})
	}
}
