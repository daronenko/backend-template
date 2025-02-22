package util

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"

	"github.com/daronenko/backend-template/internal/app/config"
	"github.com/daronenko/backend-template/pkg/sanitize"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// Get request ID from Fiber context
func GetRequestID(c *fiber.Ctx) string {
	return c.Get(fiber.HeaderXRequestID)
}

// ReqIDCtxKey is a key used for the Request ID in context
type ReqIDCtxKey struct{}

// Get ctx with timeout and request ID from Fiber context
func GetCtxWithReqID(c *fiber.Ctx) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	ctx = context.WithValue(ctx, ReqIDCtxKey{}, GetRequestID(c))
	return ctx, cancel
}

// Get context with request ID
func GetRequestCtx(c *fiber.Ctx) context.Context {
	return context.WithValue(context.Background(), ReqIDCtxKey{}, GetRequestID(c))
}

// Get config path for local or docker
func GetConfigPath(configPath string) string {
	if configPath == "docker" {
		return "./config/config-docker"
	}
	return "./config/config-local"
}

// Configure JWT cookie
func ConfigureJWTCookie(cfg *config.Config, jwtToken string) *fiber.Cookie {
	return &fiber.Cookie{
		Name:     cfg.Service.Auth.Session.Cookie.Name,
		Value:    jwtToken,
		Path:     "/",
		MaxAge:   cfg.Service.Auth.Session.Cache.Expire,
		Secure:   cfg.Service.Auth.Session.Cookie.Secure,
		HTTPOnly: cfg.Service.Auth.Session.Cookie.HTTPOnly,
	}
}

// Create session cookie
func CreateSessionCookie(cfg *config.Config, session string) *fiber.Cookie {
	return &fiber.Cookie{
		Name:     cfg.Service.Auth.Session.Cookie.Name,
		Value:    session,
		Path:     "/",
		MaxAge:   cfg.Service.Auth.Session.Cache.Expire,
		Secure:   cfg.Service.Auth.Session.Cookie.Secure,
		HTTPOnly: cfg.Service.Auth.Session.Cookie.HTTPOnly,
	}
}

// Delete session cookie
func DeleteSessionCookie(c *fiber.Ctx, sessionName string) {
	c.Cookie(&fiber.Cookie{
		Name:     sessionName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HTTPOnly: true,
	})
}

// UserCtxKey is a key used for the User object in the context
type UserCtxKey struct{}

// Get user from context
// func GetUserFromCtx(ctx context.Context) (*model.User, error) {
// 	user, ok := ctx.Value(UserCtxKey{}).(*model.User)
// 	if !ok {
// 		return nil, httpErrors.Unauthorized
// 	}
// 	return user, nil
// }

// Get user IP address
func GetIPAddress(c *fiber.Ctx) string {
	return c.IP()
}

// Error response with logging error for Fiber context
// func ErrResponseWithLog(c *fiber.Ctx, log logger.Logger, err error) error {
// 	errRespStatus, errResp := httpErrors.ErrorResponse(err)

// 	log.Errorf(
// 		"ErrResponseWithLog, RequestID: %s, IPAddress: %s, Error: %v",
// 		GetRequestID(c),
// 		GetIPAddress(c),
// 		err,
// 	)

// 	return c.Status(errRespStatus).JSON(errResp)
// }

// Log response error for Fiber context
// func LogResponseError(c *fiber.Ctx, log logger.Logger, err error) {
// 	log.Errorf(
// 		"ErrResponseWithLog, RequestID: %s, IPAddress: %s, Error: %s",
// 		GetRequestID(c),
// 		GetIPAddress(c),
// 		err,
// 	)
// }

// Read request body and validate
func ReadRequest(c *fiber.Ctx, request interface{}) error {
	if err := c.BodyParser(request); err != nil {
		return err
	}
	return validate.Struct(request)
}

// Read an image from the request
// func ReadImage(c *fiber.Ctx, field string) (*multipart.FileHeader, error) {
// 	image, err := c.FormFile(field)
// 	if err != nil {
// 		return nil, errors.WithMessage(err, "c.FormFile")
// 	}

// 	// Check content type of image
// 	if err = CheckImageContentType(image); err != nil {
// 		return nil, err
// 	}

// 	return image, nil
// }

// Read, sanitize, and validate request body
func SanitizeRequest(c *fiber.Ctx, request interface{}) error {
	body := c.Body()
	sanBody, err := sanitize.SanitizeJSON(body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("")
	}

	if err = json.Unmarshal(sanBody, request); err != nil {
		return err
	}

	return validate.Struct(request)
}

var allowedImagesContentTypes = map[string]string{
	"image/bmp":                "bmp",
	"image/gif":                "gif",
	"image/png":                "png",
	"image/jpeg":               "jpeg",
	"image/jpg":                "jpg",
	"image/svg+xml":            "svg",
	"image/webp":               "webp",
	"image/tiff":               "tiff",
	"image/vnd.microsoft.icon": "ico",
}

// Check the content type of an image file
func CheckImageFileContentType(fileContent []byte) (string, error) {
	contentType := http.DetectContentType(fileContent)

	extension, ok := allowedImagesContentTypes[contentType]
	if !ok {
		return "", errors.New("this content type is not allowed")
	}

	return extension, nil
}
