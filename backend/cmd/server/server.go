package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/foreverd34d/aumsu-elib/internal/app"
	"github.com/foreverd34d/aumsu-elib/internal/handler"
	"github.com/foreverd34d/aumsu-elib/internal/repo/postgres"
	"github.com/foreverd34d/aumsu-elib/internal/service"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echolog "github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

const defaultPort = "8080"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Couldn't read the .env file: %v\n", err)
	}
	if err := initViperConfig(); err != nil {
		log.Fatalf("Couldn't read config file: %v\n", err)
	}

	tokenSigningKey := os.Getenv("TOKEN_SIGNING_KEY")
	if tokenSigningKey == "" {
		log.Fatalln("Couldn't get signing key: TOKEN_SIGNING_KEY is not defined")
	}

	dbCtx, dbCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer dbCancel()
	db, err := postgres.NewDB(dbCtx, GetDBConfig())
	if err != nil {
		log.Fatalf("Couldn't connect to db: %v\n", err)
	}
	defer db.Close()

	handler := initHandler(db, tokenSigningKey)
	app := app.NewApp(handler, tokenSigningKey)

	port := viper.GetString("server.port")
	if port == "" {
		port = defaultPort
	}

	app.Logger.SetLevel(echolog.INFO)
	runApp(app, port)
}

func GetDBConfig() postgres.Config {
	dbPort, _ := strconv.Atoi(viper.GetString("database.port"))
	cfg := postgres.Config{
		Host:    viper.GetString("database.host"),
		Port:    dbPort,
		User:    viper.GetString("database.user"),
		DBName:  viper.GetString("database.dbname"),
		SSLMode: viper.GetString("database.sslmode"),
	}
	password := os.Getenv("DB_PASSWORD")
	if password != "" {
		cfg.Password = &password
	}
	return cfg
}

func initViperConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func initHandler(db *sqlx.DB, tokenSigningKey string) *handler.Handler {
	userRepo := postgres.NewUserPostgresRepo(db)
	userService := service.NewUserService(userRepo)

	sessionRepo := postgres.NewSessionPostgesRepo(db)
	sessionService := service.NewSessionService(userRepo, sessionRepo, []byte(tokenSigningKey))

	groupRepo := postgres.NewGroupPostgresRepo(db)
	groupService := service.NewGroupService(groupRepo)

	specialtyRepo := postgres.NewSpecialtyPostgresRepo(db)
	specialtyService := service.NewSpecialtyService(specialtyRepo)

	return &handler.Handler{
		User:      userService,
		Session:   sessionService,
		Group:     groupService,
		Specialty: specialtyService,
	}
}

func runApp(app *echo.Echo, port string) {
	appCtx, appStop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer appStop()
	go func() {
		if err := app.Start(":" + port); err != nil && err != http.ErrServerClosed {
			app.Logger.Fatalf("Error while running server: %v", err)
		}
	}()
	<-appCtx.Done()

	app.Logger.Info("Interrupt signal received, shutting down...")
	stopCtx, stopCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer stopCancel()
	if err := app.Shutdown(stopCtx); err != nil {
		app.Logger.Fatalf("Couldn't shutdown gracefully: %v", err)
	}
}
