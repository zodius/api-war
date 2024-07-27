package main

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/zodius/api-war/handler/generic"
	"github.com/zodius/api-war/handler/restful"
	"github.com/zodius/api-war/repo"
	"github.com/zodius/api-war/service"
)

func main() {
	app := gin.Default()
	redisClient := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	defer redisClient.Close()

	repo := repo.NewRepo(redisClient)

	service := service.NewService(repo)

	generic.RegisterHandler(service, app)
	restful.RegisterHandler(service, app)

	app.Run(":8971")
}
