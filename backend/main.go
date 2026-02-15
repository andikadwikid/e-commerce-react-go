package main

import (
	"github.com/gin-gonic/gin"

	"backend-commerce/config"
	"backend-commerce/database"
	seeders "backend-commerce/database/seeder"

)

func main() {
	// load config .env
	config.LoadEnv()

	// inisialisasi database
	database.InitDB()

	//run seeders
	seeders.Seed()

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	router.Run(":" + config.GetEnv("APP_PORT", "3000"))
}
