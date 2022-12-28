package controllers

import (
	"go-auth/dto"
	"go-auth/models"
	"go-auth/repository"
	"go-auth/security"
	"go-auth/util"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/asaskevich/govalidator.v9"
)

type AuthController interface {
	SignUp(ctx *fiber.Ctx) error
	SignIn(ctx *fiber.Ctx) error
	GetProfile(ctx *fiber.Ctx) error
}

type authController struct {
	usersRepository repository.UsersRepository
}

func NewAuthController(usersRepo repository.UsersRepository) AuthController {
	return &authController{usersRepository: usersRepo}
}

func (c *authController) SignUp(ctx *fiber.Ctx) error {

	var newUser models.User
	err := ctx.BodyParser(&newUser)

	log.Println(newUser)

	if err != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(util.NewError(err))
	}

	newUser.Email = util.NormalizeEmail(newUser.Email)

	if !govalidator.IsEmail(newUser.Email) {
		return ctx.Status(http.StatusBadRequest).JSON(util.NewError(util.ErrInvalidEmail))
	}

	exists, err := c.usersRepository.GetByEmail(newUser.Email)

	log.Println("find email err", err)

	if err == mongo.ErrNoDocuments {
		if strings.TrimSpace(newUser.Password) == "" {
			return ctx.Status(http.StatusBadRequest).JSON(util.NewError(util.ErrEmptyPassword))
		}

		newUser.Password, err = security.EncryptPassword(newUser.Password)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(util.NewError(err))
		}

		newUser.CreateAt = time.Now()
		newUser.UpdateAt = time.Now()
		newUser.Id = primitive.NewObjectID()

		err = c.usersRepository.Save(&newUser)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(util.NewError(err))
		}

		token, err := security.NewToken(newUser.Id.Hex())
		if err != nil {
			log.Printf("%s signin failed: %v\n", newUser.Email, err.Error())
			return ctx.Status(http.StatusUnprocessableEntity).JSON(util.NewError(err))
		}

		userData := &dto.User{
			Id:    newUser.Id,
			Email: newUser.Email,
		}

		return ctx.Status(http.StatusCreated).JSON(fiber.Map{
			"token": token,
			"user":  userData,
		})
	}

	if exists != nil {
		err = util.ErrEmailAlreadyExists
	}

	return ctx.Status(http.StatusBadRequest).JSON(util.NewError(err))
}

func (c *authController) SignIn(ctx *fiber.Ctx) error {
	var input models.User
	err := ctx.BodyParser(&input)

	if err != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(util.NewError(err))
	}

	input.Email = util.NormalizeEmail(input.Email)

	user, err := c.usersRepository.GetByEmail(input.Email)
	if err != nil {
		log.Printf("%s signin failed: %v\n", input.Email, err.Error())
		return ctx.Status(http.StatusUnprocessableEntity).JSON(util.NewError(util.ErrInvalidCredentials))
	}

	err = security.VerifyPassword(user.Password, input.Password)
	if err != nil {
		log.Printf("%s signin failed: %v\n", input.Email, err.Error())
		return ctx.Status(http.StatusUnprocessableEntity).JSON(util.NewError(err))
	}
	token, err := security.NewToken(user.Id.Hex())
	if err != nil {
		log.Printf("%s signin failed: %v\n", input.Email, err.Error())
		return ctx.Status(http.StatusUnprocessableEntity).JSON(util.NewError(err))
	}

	response := &dto.User{
		Id:    user.Id,
		Email: user.Email,
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"user":  response,
		"token": token,
	})
}

func (c *authController) GetProfile(ctx *fiber.Ctx) error {

	uid, err := AuthRequestWithId(ctx)

	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(util.NewError(err))
	}

	user, err := c.usersRepository.GetById(uid)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(util.NewError(err))
	}

	profile := &dto.User{
		Id:    user.Id,
		Email: user.Email,
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "success",
		"profile": profile,
	})
}

func AuthRequestWithId(ctx *fiber.Ctx) (string, error) {

	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	uid := claims["id"].(string)

	return uid, nil
}
