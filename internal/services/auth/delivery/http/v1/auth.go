package httpdelivery

import (
	"fmt"

	"github.com/daronenko/backend-template/internal/app/config"
	"github.com/daronenko/backend-template/internal/model/v1"
	"github.com/daronenko/backend-template/internal/pkg/errs"
	"github.com/daronenko/backend-template/internal/server/svr"
	user "github.com/daronenko/backend-template/internal/services/auth/usecase/v1"
	session "github.com/daronenko/backend-template/internal/services/session/usecase/v1"
	"github.com/daronenko/backend-template/internal/util"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

const (
	UserKey     = "user"
	UserIDParam = "user_id"
)

type Auth struct {
	fx.In

	Conf   *config.Config
	Tracer trace.Tracer

	UserUsecase    user.Usecase
	SessionUsecase session.Usecase
}

func InitAuth(d Auth, v1 *svr.V1) {
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
// @Success 201 {object} model.User
// @Router /api/v1/auth/register [post]
func (d *Auth) Register(c *fiber.Ctx) error {
	ctx := c.UserContext()

	user := &model.User{}
	if err := util.ReadRequest(c, user); err != nil {
		return errs.ErrInvalidReq
	}

	createdUser, err := d.UserUsecase.Register(ctx, user)
	if err != nil {
		return MapUsecaseError(err)
	}

	sess, err := d.SessionUsecase.Create(ctx, &model.Session{
		UserID: createdUser.User.ID,
	})
	if err != nil {
		return MapUsecaseError(err)
	}

	c.Cookie(util.CreateSessionCookie(d.Conf, sess))

	return c.Status(fiber.StatusCreated).JSON(createdUser)
}

// Login godoc
// @Summary Login new user
// @Description login user, returns user and set session
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} model.User
// @Router /api/v1/auth/login [post]
func (d *Auth) Login(c *fiber.Ctx) error {
	type Login struct {
		Email    string `json:"email" db:"email" validate:"omitempty,lte=60,email"`
		Password string `json:"password,omitempty" db:"password" validate:"required,gte=6"`
	}

	ctx, span := d.Tracer.Start(c.UserContext(), "delivery.Auth.Login")
	defer span.End()

	login := &Login{}
	if err := util.ReadRequest(c, login); err != nil {
		return errs.ErrInvalidReq
	}

	userWithToken, err := d.UserUsecase.Login(ctx, &model.User{
		Email:    login.Email,
		Password: login.Password,
	})
	if err != nil {
		return MapUsecaseError(err)
	}

	sess, err := d.SessionUsecase.Create(ctx, &model.Session{
		UserID: userWithToken.User.ID,
	})
	if err != nil {
		return MapUsecaseError(err)
	}

	c.Cookie(util.CreateSessionCookie(d.Conf, sess))

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
	ctx := c.UserContext()

	sessionIDStr := c.Cookies(d.Conf.App.Auth.Session.Cookie.Name)
	if sessionIDStr == "" {
		return errs.ErrUnauthorized
	}
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		return errs.ErrInvalidReq
	}

	if err := d.SessionUsecase.DeleteByID(ctx, sessionID); err != nil {
		return MapUsecaseError(err)
	}

	c.ClearCookie(d.Conf.App.Auth.Session.Cookie.Name)

	return c.SendStatus(fiber.StatusOK)
}

// Update godoc
// @Summary Update user
// @Description update existing user
// @Tags Auth
// @Accept json
// @Param id path int true "user_id"
// @Produce json
// @Success 200 {object} model.User
// @Router /api/v1/auth/{id} [put]
func (d *Auth) Update(c *fiber.Ctx) error {
	ctx := c.UserContext()

	uID, err := uuid.Parse(c.Params(UserIDParam))
	if err != nil {
		return errs.ErrInvalidReq
	}

	user := &model.User{}
	user.ID = uID

	if err = util.ReadRequest(c, user); err != nil {
		return errs.ErrInvalidReq
	}

	updatedUser, err := d.UserUsecase.Update(ctx, user)
	if err != nil {
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
	ctx := c.UserContext()

	uID, err := uuid.Parse(c.Params(UserIDParam))
	if err != nil {
		return errs.ErrInvalidReq
	}

	if err = d.UserUsecase.Delete(ctx, uID); err != nil {
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
// @Success 200 {object} model.UsersList
// @Failure 500 {object} httpErrors.RestError
// @Router /auth/search [get]
func (d *Auth) FindByUsername(c *fiber.Ctx) error {
	ctx := c.UserContext()

	if c.Query("username") == "" {
		return ErrMissingUsernameQuery
	}

	paginationQuery, err := util.GetPaginationFromCtx(c)
	if err != nil {
		return errs.ErrInvalidReq
	}

	response, err := d.UserUsecase.FindByUsername(ctx, c.Query("username"), paginationQuery)
	if err != nil {
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
// @Success 200 {object} model.UsersList
// @Failure 500 {object} httpErrors.RestError
// @Router /auth/all [get]
func (d *Auth) GetUsers(c *fiber.Ctx) error {
	ctx := c.UserContext()

	paginationQuery, err := util.GetPaginationFromCtx(c)
	if err != nil {
		return errs.ErrPaginationQueryMissing
	}

	usersList, err := d.UserUsecase.GetUsers(ctx, paginationQuery)
	if err != nil {
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
// @Success 200 {object} model.User
// @Failure 500 {object} httpErrors.RestError
// @Router /auth/me [get]
func (d *Auth) GetMe(c *fiber.Ctx) error {
	user, ok := c.Locals(UserKey).(*model.User)
	if !ok {
		return errs.ErrUnauthorized
	}

	return c.Status(fiber.StatusOK).JSON(user)
}
