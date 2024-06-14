package app

import (
	"errors"

	"github.com/foreverd34d/aumsu-elib/internal/errs"
	"github.com/foreverd34d/aumsu-elib/internal/handler"
	"github.com/foreverd34d/aumsu-elib/internal/model"
	"github.com/go-playground/validator/v10"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewApp(h *handler.Handler, tokenSigningKey string) *echo.Echo {
	app := echo.New()
	app.Use(middleware.Logger())
	app.Use(middleware.Recover())

	app.Validator = &BindValidator{validator: validator.New()}

	auth := app.Group("/auth", mapErrorsMiddleware)
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
	api := app.Group("/api", echojwt.WithConfig(jwtConfig), mapErrorsMiddleware)
	{
		users := api.Group("/users", checkRoleMiddleware(model.AdminRole))
		{
			users.POST("", h.CreateUser)
			users.GET("", h.GetAllUsers)
			users.GET("/:id", h.GetUser)
			users.PUT("/:id", h.UpdateUser)
			users.DELETE("/:id", h.DeleteUser)
		}
		groups := api.Group("/groups", checkRoleMiddleware(model.AdminRole))
		{
			groups.POST("", h.CreateGroup)
			groups.GET("", h.GetAllGroups)
			groups.GET("/:id", h.GetGroup)
			groups.PUT("/:id", h.UpdateGroup)
			groups.DELETE("/:id", h.DeleteGroup)
		}
		specialties := api.Group("/specialties", checkRoleMiddleware(model.AdminRole))
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

type BindValidator struct {
	validator *validator.Validate
}

func (bv *BindValidator) Validate(i any) error {
	return bv.validator.Struct(i)
}

func checkRoleMiddleware(role model.UserRole) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, err := extractUserFromContext(c)
			if err != nil {
				return echo.ErrUnauthorized
			}
			if user.Role < role {
				return echo.ErrForbidden
			}
			return next(c)
		}
	}
}

func extractUserFromContext(c echo.Context) (*model.JWTClaims, error) {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return nil, echo.ErrUnauthorized
	}
	userClaims, ok := token.Claims.(*model.JWTClaims)
	if !ok {
		return nil, echo.ErrUnauthorized
	}
	return userClaims, nil
}

func mapErrorsMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err != nil {
			if errors.Is(err, errs.NotFound) {
				return echo.ErrNotFound.WithInternal(err)
			}
			if errors.Is(err, errs.RefreshExpired) || errors.Is(err, errs.InvalidPassword) || errors.Is(err, errs.InvalidLogin) {
				return echo.ErrUnauthorized.WithInternal(err)
			}
		}
		return err
	}
}
