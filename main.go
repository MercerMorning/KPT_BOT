package main

import (
	"KPT_BOT/clients"
	"KPT_BOT/config"
	"KPT_BOT/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
)

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
	}))
	bot := clients.Init()
	handlers.Init(bot)
	log.Fatal(app.Listen(":" + config.Config("PORT")))
}
