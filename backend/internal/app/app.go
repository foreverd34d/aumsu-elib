package app

import (
	"github.com/foreverd34d/aumsu-elib/internal/handler"
	"github.com/foreverd34d/aumsu-elib/internal/model"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewApp(handler *handler.Handler, tokenSigningKey string) *echo.Echo {
	app := echo.New()
	app.Use(middleware.Logger())
	app.Use(middleware.Recover())

	auth := app.Group("/auth")
	{
		auth.POST("/session", handler.CreateSession)
		auth.PUT("/session", handler.UpdateSession)
		auth.DELETE("/session", handler.DeleteSession)
	}

	jwtConfig := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(model.TokenClaims)
		},
		SigningKey: []byte(tokenSigningKey),
	}
	api := app.Group("/api", echojwt.WithConfig(jwtConfig))
	{
		users := api.Group("/users")
		{
			users.GET("", handler.GetAllUsers)
			users.GET("/:id", handler.GetUser)
			users.POST("", handler.CreateUser)
			users.PUT("/:id", handler.UpdateUser)
			users.DELETE("/:id", handler.DeleteUser)
		}
	}

	return app
}
