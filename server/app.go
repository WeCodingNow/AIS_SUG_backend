package server

import (
	"context"
	"database/sql"
	"fmt"
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
	"github.com/spf13/viper"

	echoLog "github.com/labstack/gommon/log"

	_ "github.com/lib/pq"
)

type App struct {
	e *echo.Echo

	authUC auth.UseCase
}

const (
	apiURL = "*"
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
			viper.GetString("auth.hash_salt"),
			[]byte(viper.GetString("auth.signing_key")),
			viper.GetDuration("auth.token_ttl"),
		),
	}
}

func makePostgresString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		viper.GetString("postgres.host"), viper.GetInt("postgres.port"),
		viper.GetString("postgres.user"), viper.GetString("postgres.password"),
		viper.GetString("postgres.dbname"),
	)
}

func initDB() *sql.DB {
	if err := config.Init(); err != nil {
		panic(err)
	}

	db, err := sql.Open("postgres", makePostgresString())

	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}
