package server

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	e := gin.New()
	e.Use(cors.Default())

	store := cookie.NewStore([]byte("top_secret"))
	e.Use(sessions.Sessions("player", store))

	e.Use(gin.Logger())
	e.Use(gin.Recovery())

	e.NoRoute(handleNoRoute)
	return e
}

func handleNoRoute(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", "")
}
