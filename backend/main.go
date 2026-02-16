package main

import (
	"backend-commerce/config"
	"backend-commerce/database"
	seeders "backend-commerce/database/seeder"
	"backend-commerce/routes"

)

func main() {
	// load config .env
	config.LoadEnv()

	// inisialisasi database
	database.InitDB()

	//run seeders
	seeders.Seed()

	router := routes.SetupRouter()

	router.Run(":" + config.GetEnv("APP_PORT", "3000"))
}
