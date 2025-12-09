package controller

import (
	"github.com/gofiber/fiber/v2"

	"p2p-back-end/logs"
	"p2p-back-end/modules/entities/models"
	"p2p-back-end/pkg/middlewares"
)

type authController struct {
	authSrv models.AuthService
}

func NewUserController(router fiber.Router, authSrv models.AuthService) {
	controller := &authController{authSrv: authSrv}
	router.Post("/register", controller.register)
	router.Post("/login", controller.login)
	router.Get("/test11", middlewares.JwtAuthentication(controller.test11))
	// router.Get("/user/check-by-id", middlewares.JwtAuthentication(controller.checkUserByID))
	// router.Put("/user/reset-password", middlewares.JwtAuthentication(controller.resetPassword))
}

func (h authController) test11(c *fiber.Ctx, userInfo *models.UserInfo) error {
	m := "hello"
	return responseSuccess(c, m)
}

// @Summary User registration
// @Description Register a new user
// @Tags User
// @Accept json
// @Produce json
// @Param register body models.RegisterReq true "Registration request"
// @Success 200 {object} models.ResponseData{data=string}
// @Failure 400 {object} models.ResponseError
// @Router /v1/auth/register [post]
func (h authController) register(c *fiber.Ctx) error {
	var req models.RegisterReq
	if err := c.BodyParser(&req); err != nil {
		logs.Info("Invalid request: " + err.Error())
		return badReqErrResponse(c, "Invalid request: "+err.Error())
	}
	m, err := h.authSrv.Register(&req)
	if err != nil {
		return responseWithError(c, err)
	}

	return responseSuccess(c, m)
}

// @Summary User login
// @Description User login with username and password
// @Tags User
// @Accept json
// @Produce json
// @Param login body models.LoginReq true "Login request"
// @Success 200 {object} models.ResponseData{data=string}
// @Failure 400 {object} models.ResponseError
// @Router /v1/auth/login [post]
func (h authController) login(c *fiber.Ctx) error {
	var req models.LoginReq
	if err := c.BodyParser(&req); err != nil {
		logs.Info("Invalid request: " + err.Error())
		return badReqErrResponse(c, "Invalid request: "+err.Error())
	}
	m, err := h.authSrv.Login(&req)
	if err != nil {
		return responseWithError(c, err)
	}

	return responseSuccess(c, m)
}

// // @Summary Check if user exists
// // @Description Check if a user exists by user ID
// // @Tags User
// // @Accept json
// // @Produce json
// // @Param Authorization header string true "Bearer {token}"
// // @Success 200 {object} models.ResponseData{data=bool}
// // @Failure 400 {object} models.ResponseError
// // @Router /v1/user/check-by-id [get]
// // @Security ApiKeyAuth
// func (h authController) checkUserByID(c *fiber.Ctx, userInfo *models.UserInfo) error {

// 	m, err := h.authSrv.IsUserExistByID(userInfo.UserId)
// 	if err != nil {
// 		return responseWithError(c, err)
// 	}

// 	return responseSuccess(c, m)
// }

// // @Summary Reset Password
// // @Description Reset the password for a user
// // @Tags User
// // @Accept json
// // @Produce json
// // @Param Authorization header string true "Bearer {token}"
// // @Param req body models.ChangePasswordReq true "Change Password Request"
// // @Success 200 {string} string "Password Changed Successfully"
// // @Failure 400 {object} models.ResponseError
// // @Router /v1/user/reset-password [put]
// // @Security ApiKeyAuth
// func (h authController) resetPassword(c *fiber.Ctx, userInfo *models.UserInfo) error {

// 	var req models.ChangePasswordReq
// 	if err := c.BodyParser(&req); err != nil {
// 		logs.Info("Invalid request: " + err.Error())
// 		return badReqErrResponse(c, "Invalid request: "+err.Error())
// 	}

// 	m, err := h.authSrv.ChangePassword(&req, userInfo)
// 	if err != nil {
// 		return responseWithError(c, err)
// 	}
// 	return responseSuccess(c, m)
// }
