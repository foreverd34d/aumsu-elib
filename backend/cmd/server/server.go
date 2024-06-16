/*
Server запускает сервер электронной библиотеки.

# Конфигурация

Перед запуском сервер читает файл конфигурации configs/config.yml и
файл .env с переменными окружения. Из файла конфигурации читаются
настройки подключения к базе данных и порт сервера.
Если порт в конфигурации не указан, сервер слушает порт 8080. 
Из переменных окружения сервер читает ключ подписи jwt токенов
и пароль к базе данных, если таковой имеется.
*/
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

// Порт по умолчанию.
const defaultPort = "8080"

func main() {
	// Инициализация конфигурационных файлов
	if err := godotenv.Load(); err != nil {
		log.Printf("Couldn't read the .env file: %v. Using env variables.\n", err)
	}
	if err := initViperConfig(); err != nil {
		log.Fatalf("Couldn't read config file: %v\n", err)
	}

	// Получение ключа подписи jwt токенов
	tokenSigningKey := os.Getenv("TOKEN_SIGNING_KEY")
	if tokenSigningKey == "" {
		log.Fatalln("Couldn't get signing key: TOKEN_SIGNING_KEY is not defined")
	}

	// Подключение к базе данных
	dbCtx, dbCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer dbCancel()
	db, err := postgres.NewDB(dbCtx, getDBConfig())
	if err != nil {
		log.Fatalf("Couldn't connect to db: %v\n", err)
	}
	defer db.Close()

	// Инициализация всех путей и middleware
	handler := initHandler(db, tokenSigningKey)
	app := app.NewApp(handler, tokenSigningKey)

	// Получение порта, если есть
	port := viper.GetString("server.port")
	if port == "" {
		port = defaultPort
	}

	app.Logger.SetLevel(echolog.INFO)
	runApp(app, port)
}

// getDBConfig загружает конфигурацию базы данных из файла configs/config.yml.
func getDBConfig() postgres.Config {
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

// initViperConfig считывает файлы конфигурации из директории configs.
func initViperConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

// initHandler инициализирует все сервисы и репозитории для хэндлера.
func initHandler(db *sqlx.DB, tokenSigningKey string) *handler.Handler {
	userRepo := postgres.NewUserRepo(db)
	userService := service.NewUserService(userRepo)

	tokenRepo := postgres.NewSessionRepo(db)
	sessionService := service.NewSessionService(userRepo, tokenRepo, []byte(tokenSigningKey))

	groupRepo := postgres.NewGroupRepo(db)
	groupService := service.NewGroupService(groupRepo)

	specialtyRepo := postgres.NewSpecialtyRepo(db)
	specialtyService := service.NewSpecialtyService(specialtyRepo)

	departmentRepo := postgres.NewDepartmentRepo(db)
	departmentService := service.NewDepartmentService(departmentRepo)

	return &handler.Handler{
		User:       userService,
		Session:    sessionService,
		Group:      groupService,
		Specialty:  specialtyService,
		Department: departmentService,
	}
}

// runApp запускает сервер и корректно завершает его работу при получении SIGINT.
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
