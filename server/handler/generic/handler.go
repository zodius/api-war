package generic

import (
	"strconv"

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

	app.GET("/scoreboard", handler.CorsMiddleware(), handler.GetScoreboard)
	app.GET("/me", handler.CorsMiddleware(), handler.GetMe)
	app.GET("/map", handler.CorsMiddleware(), handler.GetMap)
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

func (h *Handler) GetMe(c *gin.Context) {
	token := c.GetHeader("X-Api-Token")
	username, err := h.Service.GetMe(token)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"username": username})
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
	startPos := 0
	endPos := 0

	// get start and end from query
	startParam := c.Query("start")
	endParam := c.Query("end")
	if startParam != "" || endParam != "" {
		// convert start and end to int
		start, err := strconv.Atoi(startParam)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid start parameter"})
			return
		}
		end, err := strconv.Atoi(endParam)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid end parameter"})
			return
		}
		if start > end {
			c.JSON(400, gin.H{"error": "Invalid start and end parameter"})
			return
		}
		if start < 0 || end > model.FieldCount {
			c.JSON(400, gin.H{"error": "Invalid start and end parameter"})
			return
		}
		startPos, endPos = start, end
	}

	mapObject, err := h.Service.GetCurrentMap(startPos, endPos)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, mapObject.Representation())
}
