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

// temporary middleware to disable CORS
func (h *Handler) CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
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
	mapObject, err := h.Service.GetCurrentMap()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, mapObject.Representation())
}
