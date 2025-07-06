package main

import (
	"KPT_BOT/clients"
	"KPT_BOT/config"
	"KPT_BOT/handlers"
	"KPT_BOT/services"
	"KPT_BOT/session"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
	"strconv"
)

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
	}))
	bot := clients.Init()
	usrSessions := map[int64]*session.Session{}
	go func() {
		handlers.Init(bot, usrSessions)
	}()
	app.Get("/:chatId?", func(c *fiber.Ctx) error {
		if c.Params("chatId") != "" {
			code := c.Query("code")
			chatId, _ := strconv.Atoi(c.Params("chatId"))
			usrSession, _ := usrSessions[int64(chatId)]
			services.InitTable(code, bot, int64(chatId), usrSession)
		}
		return c.SendString("OK")
	})
	log.Fatal(app.Listen(":" + config.Config("PORT")))
}
