package graphql

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/zodius/api-war/model"
	"github.com/zodius/api-war/tools/graph"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func RegisterHandler(service model.Service, app *gin.Engine) {
	app.POST("/graphql", graphqlHandler(service))
	app.GET("/graphiql", playgroundHandler())
}

func graphqlHandler(service model.Service) gin.HandlerFunc {
	h := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		Service: service,
	}}))
	return func(c *gin.Context) {
		// extract token from header
		token := c.GetHeader("X-Api-Token")
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, "token", token)
		request := c.Request.WithContext(ctx)
		h.ServeHTTP(c.Writer, request)
	}
}

func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/graphql")
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
