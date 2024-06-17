// Пакет app предоставляет middleware и функции для создания маршрутов приложения.
package app

import (
	"errors"
	"fmt"

	"github.com/foreverd34d/aumsu-elib/internal/errs"
	"github.com/foreverd34d/aumsu-elib/internal/handler"
	"github.com/foreverd34d/aumsu-elib/internal/model"
	"github.com/go-playground/validator/v10"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// NewApp создает новый экземпляр [echo.Echo] с настроенными middleware и маршрутами для хэндлеров.
// tokenSigningKey используется для подписи и проверки jwt токенов.
func NewApp(h *handler.Handler, tokenSigningKey string) *echo.Echo {
	app := echo.New()
	app.Use(middleware.Logger())
	app.Use(middleware.Recover())
	app.Use(mapErrors)

	app.Validator = &BindValidator{validator: validator.New(validator.WithRequiredStructEnabled())}

	auth := app.Group("/auth")
	{
		auth.POST("/session", h.CreateSession)
		auth.PUT("/session", h.UpdateSession)
		auth.DELETE("/session", h.DeleteSession)
	}

	jwtConfig := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(model.JWTClaims)
		},
		SigningKey: []byte(tokenSigningKey),
	}
	api := app.Group("/api", echojwt.WithConfig(jwtConfig))
	{
		users := api.Group("/users", checkRole(model.AdminRole))
		{
			users.POST("", h.CreateUser)
			users.GET("", h.GetAllUsers)
			users.GET("/:id", h.GetUser)
			users.PUT("/:id", h.UpdateUser)
			users.DELETE("/:id", h.DeleteUser)
		}
		groups := api.Group("/groups", checkRole(model.AdminRole))
		{
			groups.POST("", h.CreateGroup)
			groups.GET("", h.GetAllGroups)
			groups.GET("/:id", h.GetGroup)
			groups.PUT("/:id", h.UpdateGroup)
			groups.DELETE("/:id", h.DeleteGroup)
		}
		specialties := api.Group("/specialties", checkRole(model.AdminRole))
		{
			specialties.POST("", h.CreateSpecialty)
			specialties.GET("", h.GetAllSpecialties)
			specialties.GET("/:id", h.GetSpecialty)
			specialties.PUT("/:id", h.UpdateSpecialty)
			specialties.DELETE(":/id", h.DeleteSpecialty)
		}
	}

	return app
}

// BindValidator представляет валидатор структур для хэндлеров.
type BindValidator struct {
	validator *validator.Validate
}

// Validate реализует метод [echo.Validator] и валидирует входную структуру.
func (bv *BindValidator) Validate(i any) error {
	return bv.validator.Struct(i)
}

// checkRole предоставляет middleware для проверки возможности пользователя с его ролью
// обращаться по защищенному маршруту. В случае отказа возвращается ошибка [echo.ErrForbidden]
func checkRole(role model.UserRole) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, err := extractUser(c)
			if err != nil {
				return err
			}
			if user.Role < role {
				return echo.ErrForbidden.WithInternal(fmt.Errorf("user has no permission to access protected route"))
			}
			return next(c)
		}
	}
}

// extractUser возвращает полезную нагрузку jwt токена пользователя.
// В случае неудачи возвращается ошибка [echo.ErrUnauthorized]
func extractUser(c echo.Context) (*model.JWTClaims, error) {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return nil, echo.ErrUnauthorized.WithInternal(fmt.Errorf("no user token"))
	}
	userClaims, ok := token.Claims.(*model.JWTClaims)
	if !ok {
		return nil, echo.ErrUnauthorized.WithInternal(fmt.Errorf("invalid user claims"))
	}
	return userClaims, nil
}

// mapErrors предоставляет middleware для преобразования ошибок из пакета [errs] в http-ошибки [echo].
func mapErrors(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err != nil {
			if errors.Is(err, errs.NotFound) {
				return echo.ErrNotFound.WithInternal(err)
			}
			if errors.Is(err, errs.RefreshExpired) ||
				errors.Is(err, errs.InvalidPassword) ||
				errors.Is(err, errs.InvalidLogin) {
				return echo.ErrUnauthorized.WithInternal(err)
			}
		}
		return err
	}
}
