package main

import (
	"go-auth/controllers"
	"go-auth/db"
	"go-auth/repository"
	"go-auth/routes"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Panicln(err)
	}
}

func main() {

	conn := db.NewConnection()

	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "Hello World"})
	})

	usersRepository := repository.NewUserRepository(conn)
	authController := controllers.NewAuthController(usersRepository)
	authRoutes := routes.NewAuthRoutes(authController)
	authRoutes.Install(app)

	log.Fatal(app.Listen(":4000"))
}
