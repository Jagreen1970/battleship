package endpoints

import (
	"net/http"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	gameAPI GameAPI
}

const (
	sessionKeyPlayerName = "playerName"
)

func NewController(api GameAPI) *Controller {
	return &Controller{
		gameAPI: api,
	}
}

func (c *Controller) Register(engine *gin.Engine) {
	engine.Use(static.Serve("/", static.LocalFile("./frontend/build", true)))

	api := engine.Group("/api")
	{
		api.GET("/", func(context *gin.Context) {
			context.JSON(http.StatusOK, gin.H{"message": "pong"})
		})
		api.POST("/login", c.Login)
		api.GET("/logout", c.Logout)
		api.GET("/scoreboard", c.Scoreboard)
		api.GET("/games", c.Games)
		api.POST("/games", c.CreateGame)
		api.GET("/games/:id", c.GetGame)
		api.PATCH("/games/:id", c.JoinGame)
		api.PUT("/games/:id/pin/:pin", c.PlacePin)
		api.DELETE("/games/:id/pin/:pin", c.RecoverPin)
		api.GET("/games/:id/start", c.StartGame)
		api.POST("/games/:id/target", c.Target)
	}
}
