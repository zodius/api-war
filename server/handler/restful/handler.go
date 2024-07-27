package restful

import (
	"errors"
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

	api := app.Group("/api/v1")

	api.POST("/register", handler.Register)
	api.POST("/login", handler.Login)
	api.GET("/userlist", handler.GetUserList)
	api.POST("/conquer/:id", handler.Conquer)
	api.GET("/me", handler.Me)

	// leak openapi spec
	app.GET("/openapi.yml", func(c *gin.Context) {
		c.File("handler/restful/openapi.yml")
	})
}

func (h *Handler) Conquer(c *gin.Context) {
	token := c.GetHeader("X-API-TOKEN")
	if token == "" {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	fieldID := c.Param("id")
	if fieldID == "" {
		c.JSON(400, gin.H{"error": "field id is required"})
		return
	}

	fieldIDInt, err := strconv.Atoi(fieldID)
	if err != nil {
		c.JSON(400, gin.H{"error": "field id must be integer"})
		return
	}

	if err := h.Service.ConquerField(token, fieldIDInt, model.TypeRestful); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{})
}

func (h *Handler) Me(c *gin.Context) {
	token := c.GetHeader("X-API-TOKEN")
	if token == "" {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	fields, err := h.Service.GetUserConquerField(token, model.TypeRestful)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"fields": fields})
}

func (h *Handler) Register(c *gin.Context) {
	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req request
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.Register(req.Username, req.Password); err != nil {
		if errors.Is(err, model.ErrUserExist) {
			c.JSON(400, gin.H{"error": err.Error()})
		} else {
			c.JSON(500, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(200, gin.H{})
}

func (h *Handler) Login(c *gin.Context) {
	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req request
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	token, err := h.Service.Login(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, model.ErrInvalidCredentials) {
			c.JSON(401, gin.H{"error": "invalid credentials"})
		} else {
			c.JSON(500, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(200, gin.H{"token": token})
}

func (h *Handler) GetUserList(c *gin.Context) {
	token := c.GetHeader("X-API-TOKEN")
	if token == "" {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	userList, err := h.Service.GetUserList(token)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"userList": userList})
}
