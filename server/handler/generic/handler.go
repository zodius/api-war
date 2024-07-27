package generic

import (
	"github.com/gin-gonic/gin"
	"github.com/zodius/api-war/model"
)

type Handler struct {
	Service model.Service
}

func RegisterHandler(service model.Service, app *gin.Engine) {
	handler := Handler{
		Service: service,
	}

	app.GET("/scoreboard", handler.GetScoreboard)
	app.GET("/map", handler.GetMap)
}

func (h *Handler) GetScoreboard(c *gin.Context) {
	scoreList, err := h.Service.GetScoreboard()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"scoreList": scoreList})
}

func (h *Handler) GetMap(c *gin.Context) {
	mapStringRepresentation, err := h.Service.GetCurrentMap()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"map": mapStringRepresentation})
}
