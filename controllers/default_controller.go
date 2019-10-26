package controllers

import "github.com/gin-gonic/gin"

var (
	Router *gin.Engine
)

type Route struct {
	Method string
	Path   string
}

func Home(c *gin.Context) {
	var routes []Route
	for _, route := range Router.Routes() {
		routes = append(routes, Route{
			Method: route.Method,
			Path:   route.Path,
		})
	}
	c.JSON(200, gin.H{
		"Routes": routes,
	})
}
