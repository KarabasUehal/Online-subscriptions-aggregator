package main

import (
	agghandlers "aggregator/agg-handlers"
	database "aggregator/database"
	"log"
	"os"

	_ "aggregator/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {

	if err := database.InitDatabase(); err != nil {
		log.Printf("Failed to initializing database: %e", err)
	}

	router := gin.Default()

	router.GET("/subscriptions", agghandlers.GetAllSubs(database.DB))
	router.GET("/subscribe/:id", agghandlers.GetSubsById(database.DB))
	router.POST("/subscribe", agghandlers.AddSubscription(database.DB))
	router.PUT("/subscribe/:id", agghandlers.UpdateSubscriptionById(database.DB))
	router.DELETE("/subscribe/:id", agghandlers.DeleteSubById(database.DB))
	router.GET("/total-cost", agghandlers.CalculateSubsPerPeriod(database.DB))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	router.Run(":" + port)

}
