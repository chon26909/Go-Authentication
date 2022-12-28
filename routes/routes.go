package routes

import (
	"go-auth/security"
	"go-auth/util"
	"net/http"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

type Routes interface {
	Install(app *fiber.App)
}

func Authenticate(ctx *fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		SigningKey:    security.JwtSecretKey,
		SigningMethod: security.JwtSigningMethod.Name,
		TokenLookup:   "header:Authorization",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(http.StatusUnauthorized).JSON(util.NewError(err))
		},
	})(ctx)
}
