package httpdelivery

import (
	"fmt"

	"github.com/daronenko/backend-template/internal/app/config"
	"github.com/daronenko/backend-template/internal/models"
	user "github.com/daronenko/backend-template/internal/modules/auth/usecase/v1"
	session "github.com/daronenko/backend-template/internal/modules/session/usecase/v1"
	"github.com/daronenko/backend-template/internal/pkg/errs"
	"github.com/daronenko/backend-template/internal/server/svr"
	"github.com/daronenko/backend-template/pkg/logger"
	"github.com/daronenko/backend-template/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/fx"
)

const (
	UserKey     = "user"
	UserIDParam = "user_id"
)

type Auth struct {
	fx.In

	Conf   *config.Config
	Logger logger.Logger

	UserUsecase    user.Usecase
	SessionUsecase session.Usecase
}

func RegisterAuth(v1 *svr.V1, d Auth) {
	result := v1.Group("/auth")

	result.Post("/register", d.Register)
	result.Post("/login", d.Login)
	result.Post("/logout", d.Logout)
	result.Put(fmt.Sprintf("/:%s", UserIDParam), d.Update)
	result.Delete(fmt.Sprintf("/:%s", UserIDParam), d.Delete)
	result.Get("/search", d.FindByUsername)
	result.Get("/all", d.GetUsers)
	result.Get("/me", d.GetMe)
}

// Register godoc
// @Summary Register new user
// @Description register new user, returns user and token
// @Tags Auth
// @Accept json
// @Produce json
// @Success 201 {object} models.User
// @Router /api/v1/auth/register [post]
func (d *Auth) Register(c *fiber.Ctx) error {
	ctx := utils.GetRequestCtx(c)

	user := &models.User{}
	if err := utils.ReadRequest(c, user); err != nil {
		utils.LogResponseError(c, d.Logger, err)
		return errs.ErrInvalidReq
	}

	createdUser, err := d.UserUsecase.Register(ctx, user)
	if err != nil {
		utils.LogResponseError(c, d.Logger, err)
		return MapUsecaseError(err)
	}

	sess, err := d.SessionUsecase.Create(ctx, &models.Session{
		UserID: createdUser.User.ID,
	})
	if err != nil {
		utils.LogResponseError(c, d.Logger, err)
		return MapUsecaseError(err)
	}

	c.Cookie(utils.CreateSessionCookie(d.Conf, sess))

	return c.Status(fiber.StatusCreated).JSON(createdUser)
}

// Login godoc
// @Summary Login new user
// @Description login user, returns user and set session
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} models.User
// @Router /api/v1/auth/login [post]
func (d *Auth) Login(c *fiber.Ctx) error {
	type Login struct {
		Email    string `json:"email" db:"email" validate:"omitempty,lte=60,email"`
		Password string `json:"password,omitempty" db:"password" validate:"required,gte=6"`
	}

	ctx := utils.GetRequestCtx(c)

	login := &Login{}
	if err := utils.ReadRequest(c, login); err != nil {
		utils.LogResponseError(c, d.Logger, err)
		return errs.ErrInvalidReq
	}

	userWithToken, err := d.UserUsecase.Login(ctx, &models.User{
		Email:    login.Email,
		Password: login.Password,
	})
	if err != nil {
		utils.LogResponseError(c, d.Logger, err)
		return MapUsecaseError(err)
	}

	sess, err := d.SessionUsecase.Create(ctx, &models.Session{
		UserID: userWithToken.User.ID,
	})
	if err != nil {
		utils.LogResponseError(c, d.Logger, err)
		return MapUsecaseError(err)
	}

	c.Cookie(utils.CreateSessionCookie(d.Conf, sess))

	return c.JSON(userWithToken)
}

// Logout godoc
// @Summary Logout user
// @Description logout user removing session
// @Tags Auth
// @Accept  json
// @Produce  json
// @Success 200 {string} string	"ok"
// @Router /api/v1/auth/logout [post]
func (d *Auth) Logout(c *fiber.Ctx) error {
	ctx := utils.GetRequestCtx(c)

	sessionIDStr := c.Cookies(d.Conf.Service.Auth.Session.Cookie.Name)
	if sessionIDStr == "" {
		return errs.ErrUnauthorized
	}
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		utils.LogResponseError(c, d.Logger, err)
		return errs.ErrInvalidReq
	}

	if err := d.SessionUsecase.DeleteByID(ctx, sessionID); err != nil {
		utils.LogResponseError(c, d.Logger, err)
		return MapUsecaseError(err)
	}

	c.ClearCookie(d.Conf.Service.Auth.Session.Cookie.Name)

	return c.SendStatus(fiber.StatusOK)
}

// Update godoc
// @Summary Update user
// @Description update existing user
// @Tags Auth
// @Accept json
// @Param id path int true "user_id"
// @Produce json
// @Success 200 {object} models.User
// @Router /api/v1/auth/{id} [put]
func (d *Auth) Update(c *fiber.Ctx) error {
	ctx := utils.GetRequestCtx(c)

	uID, err := uuid.Parse(c.Params(UserIDParam))
	if err != nil {
		utils.LogResponseError(c, d.Logger, err)
		return errs.ErrInvalidReq
	}

	user := &models.User{}
	user.ID = uID

	if err = utils.ReadRequest(c, user); err != nil {
		utils.LogResponseError(c, d.Logger, err)
		return errs.ErrInvalidReq
	}

	updatedUser, err := d.UserUsecase.Update(ctx, user)
	if err != nil {
		utils.LogResponseError(c, d.Logger, err)
		return MapUsecaseError(err)
	}

	return c.Status(fiber.StatusOK).JSON(updatedUser)
}

// Delete
// @Summary Delete user account
// @Description some description
// @Tags Auth
// @Accept json
// @Param id path int true "user_id"
// @Produce json
// @Success 200 {string} string	"ok"
// @Failure 500 {object} httpErrors.RestError
// @Router /api/v1/auth/{id} [delete]
func (d *Auth) Delete(c *fiber.Ctx) error {
	ctx := utils.GetRequestCtx(c)

	uID, err := uuid.Parse(c.Params(UserIDParam))
	if err != nil {
		utils.LogResponseError(c, d.Logger, err)
		return errs.ErrInvalidReq
	}

	if err = d.UserUsecase.Delete(ctx, uID); err != nil {
		utils.LogResponseError(c, d.Logger, err)
		return MapUsecaseError(err)
	}

	return c.SendStatus(fiber.StatusOK)
}

// FindByName godoc
// @Summary Find by name
// @Description Find user by name
// @Tags Auth
// @Accept json
// @Param name query string false "username" Format(username)
// @Produce json
// @Success 200 {object} models.UsersList
// @Failure 500 {object} httpErrors.RestError
// @Router /auth/search [get]
func (d *Auth) FindByUsername(c *fiber.Ctx) error {
	ctx := utils.GetRequestCtx(c)

	if c.Query("username") == "" {
		utils.LogResponseError(c, d.Logger, ErrMissingUsernameQuery)
		return ErrMissingUsernameQuery
	}

	paginationQuery, err := utils.GetPaginationFromCtx(c)
	if err != nil {
		utils.LogResponseError(c, d.Logger, err)
		return errs.ErrInvalidReq
	}

	response, err := d.UserUsecase.FindByUsername(ctx, c.Query("username"), paginationQuery)
	if err != nil {
		utils.LogResponseError(c, d.Logger, err)
		return MapUsecaseError(err)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// GetUsers godoc
// @Summary Get users
// @Description Get the list of all users
// @Tags Auth
// @Accept json
// @Param page query int false "page number" Format(page)
// @Param size query int false "number of elements per page" Format(size)
// @Param orderBy query int false "filter name" Format(orderBy)
// @Produce json
// @Success 200 {object} models.UsersList
// @Failure 500 {object} httpErrors.RestError
// @Router /auth/all [get]
func (d *Auth) GetUsers(c *fiber.Ctx) error {
	ctx := utils.GetRequestCtx(c)

	paginationQuery, err := utils.GetPaginationFromCtx(c)
	if err != nil {
		utils.LogResponseError(c, d.Logger, err)
		return errs.ErrPaginationQueryMissing
	}

	usersList, err := d.UserUsecase.GetUsers(ctx, paginationQuery)
	if err != nil {
		utils.LogResponseError(c, d.Logger, err)
		return MapUsecaseError(err)
	}

	return c.Status(fiber.StatusOK).JSON(usersList)
}

// GetMe godoc
// @Summary Get user by id
// @Description Get current user by id
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} models.User
// @Failure 500 {object} httpErrors.RestError
// @Router /auth/me [get]
func (d *Auth) GetMe(c *fiber.Ctx) error {
	user, ok := c.Locals(UserKey).(*models.User)
	if !ok {
		utils.LogResponseError(c, d.Logger, errs.ErrUnauthorized)
		return errs.ErrUnauthorized
	}

	return c.Status(fiber.StatusOK).JSON(user)
}
