package api

import (
	"borda/internal/config"
	"borda/internal/services"
	"borda/internal/usecase"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
)

type Handler struct {
	AuthService  *services.AuthService
	UserUsecase  *usecase.UserUsecase
	AdminUsecase *usecase.AdminUsecase
}

func NewHandler(authService *services.AuthService,
	userUsecase *usecase.UserUsecase, adminUsecase *usecase.AdminUsecase) *Handler {
	return &Handler{
		AuthService:  authService,
		UserUsecase:  userUsecase,
		AdminUsecase: adminUsecase,
	}
}

func (h *Handler) Init(app *fiber.App) {
	app.Use(logger.New())

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			c.Path(): "pong",
			"time":   time.Now().Format(time.UnixDate),
		})
	})

	api := app.Group("/api")
	v1 := api.Group("/v1")

	h.initAuthRoutes(v1)

	v1.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(config.JWT().SigningKey),
	}))

	h.initUserRoutes(v1)
	h.initTaskRoutes(v1)
}

func AuthRequired(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["sub"].(string)
	fmt.Println("Welcome " + name)
	return c.Next()
}
