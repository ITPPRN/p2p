package servers

import (
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"

	"p2p-back-end/logs"
	_authCon "p2p-back-end/modules/auth/controller"
	_authRe "p2p-back-end/modules/auth/repository"
	_authSer "p2p-back-end/modules/auth/service"
	"p2p-back-end/pkg/middlewares"
)


func (s *server) Handlers() error {

	v1 := s.App.Group("/v1")
	// Register swagger handler
	v1.Get("/swagger/*", fiberSwagger.WrapHandler)
	v1.Use(middlewares.NewCorsOriginMiddleWare())
	v1.Use(logs.LogHttp)

	if s.Cfg.App.Mode == "release" {
		s.App.Use(fiberzap.New(fiberzap.Config{Logger: logs.Logger}))
	} else {
		s.App.Use(middlewares.NewLoggerMiddleWare())
	}

	authRepo := _authRe.NewAuthRepositoryDB(s.Db)
	authSrv := _authSer.NewAuthService(s.Keycloak, s.Cfg, authRepo)
	_authCon.NewUserController(v1.Group("/auth"), authSrv)

	s.App.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     "error, end point not found",
			"result":      nil,
		})
	})

	return nil
}
