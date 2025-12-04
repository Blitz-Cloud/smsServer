package main

import (
	"log"

	dbTypes "github.com/Blitz-Cloud/smsServer/db"
	"github.com/Blitz-Cloud/smsServer/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Env file not loaded or missing\n Server startup aborted\nExit 1")
	}	


	app := fiber.New()
	db, err := gorm.Open(sqlite.Open("main.db"), &gorm.Config{})
	db.AutoMigrate(&dbTypes.User{}, &dbTypes.Node{})
	if err != nil {
		log.Fatal(err)
	}
	app.Use(func(c *fiber.Ctx)error{
		c.Locals("db",db)
		return c.Next()
	})

	routes.RegisterApiRoutes(app)


	app.Listen(":3000")
}
