package server

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/WeCodingNow/AIS_SUG_backend/auth"
	authhttp "github.com/WeCodingNow/AIS_SUG_backend/auth/delivery/http"
	authpostgres "github.com/WeCodingNow/AIS_SUG_backend/auth/repository/postgres"
	authusecase "github.com/WeCodingNow/AIS_SUG_backend/auth/usecase"
	"github.com/WeCodingNow/AIS_SUG_backend/config"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	echoLog "github.com/labstack/gommon/log"

	_ "github.com/lib/pq"
)

type App struct {
	e *echo.Echo

	authUC auth.UseCase
}

const (
	apiURL = "localhost"
)

func (a *App) Run(port string) error {
	e := echo.New()
	e.Debug = true

	// a.authUC.ParseToken()

	e.Logger.SetLevel(echoLog.DEBUG)
	e.Use(
		middleware.Logger(),
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins:     []string{apiURL},
			AllowCredentials: true,
			AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodDelete, http.MethodOptions},
			AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderXCSRFToken},
			ExposeHeaders:    []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderXCSRFToken},
		}),
	)
	authhttp.RegisterHTTPEndpoints(e, a.authUC)

	e.Server = &http.Server{
		Addr:           ":" + port,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	a.e = e

	go func() {
		if err := a.e.Start(":" + port); err != nil {
			log.Fatalf("Failed to start: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.e.Shutdown(ctx)
}

func NewApp() *App {
	db := initDB()

	userRepo := authpostgres.NewUserRepository(db)

	return &App{
		authUC: authusecase.NewAuthUseCase(
			userRepo,
			"123",
			[]byte{'1', '2', '3'},
			24*3600,
		),
	}
}

func initDB() *sql.DB {
	config := config.Config()
	db, err := sql.Open("postgres", config.PostgresConfig.ToString())

	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}