package httpdelivery

import (
	"github.com/daronenko/backend-template/internal/app/config"
	"github.com/daronenko/backend-template/internal/models"
	user "github.com/daronenko/backend-template/internal/modules/auth/usecase/v1"
	session "github.com/daronenko/backend-template/internal/modules/session/usecase/v1"
	"github.com/daronenko/backend-template/internal/pkg/errs"
	"github.com/daronenko/backend-template/internal/server/svr"
	"github.com/daronenko/backend-template/pkg/logger"
	"github.com/daronenko/backend-template/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

type Auth struct {
	fx.In

	Conf   *config.Config
	Logger logger.Logger

	UserUsecase    user.Usecase
	SessionUsecase session.Usecase
}

func RegisterAuth(v1 *svr.V1, d Auth) {
	result := v1.Group("/user")

	result.Post("/register", d.Register)
	result.Post("/login", d.Login)
}

// Register godoc
// @Summary Register new user
// @Description register new user, returns user and token
// @Tags Auth
// @Accept json
// @Produce json
// @Success 201 {object} models.User
// @Router /auth/register [post]
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
// @Router /auth/login [post]
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
