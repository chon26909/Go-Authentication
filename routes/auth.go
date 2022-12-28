package routes

import (
	"go-auth/controllers"

	"github.com/gofiber/fiber/v2"
)

type authRoutes struct {
	authController controllers.AuthController
}

func NewAuthRoutes(authController controllers.AuthController) Routes {
	return &authRoutes{authController: authController}
}

func (r *authRoutes) Install(app *fiber.App) {
	app.Post("/signup", r.authController.SignUp)
	app.Post("/signin", r.authController.SignIn)
	app.Get("/me", Authenticate, r.authController.GetProfile)
}
