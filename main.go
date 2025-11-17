// Package main Netlify function
package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"

	"github.com/crispgm/foosbot/internal/app"
	"github.com/crispgm/foosbot/internal/def"
)

var ginLambda *ginadapter.GinLambda

func init() {
	gin.SetMode(gin.ReleaseMode)

	err := def.LoadVariables()
	if err != nil {
		log.Panic(err)
	}

	router := gin.Default()
	app.LoadRoutes(router)
	ginLambda = ginadapter.New(router)
}

func main() {
	lambda.Start(Handler)
}

// Handler aws handler
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}
