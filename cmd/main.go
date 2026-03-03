// Package main CLI entry for foosbot
package main

import (
	"github.com/gin-gonic/gin"

	"github.com/crispgm/foosbot/internal/app"
	"github.com/crispgm/foosbot/internal/def"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	err := def.LoadVariables()
	if err != nil {
		panic(err)
	}

	router := gin.Default()
	app.LoadRoutes(router)
	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
