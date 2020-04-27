package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/ais_auth"
	"github.com/WeCodingNow/AIS_SUG_backend/auth"

	authhttp "github.com/WeCodingNow/AIS_SUG_backend/auth/delivery/http"
	authpostgres "github.com/WeCodingNow/AIS_SUG_backend/auth/repository/postgres"
	authusecase "github.com/WeCodingNow/AIS_SUG_backend/auth/usecase"

	aishttp "github.com/WeCodingNow/AIS_SUG_backend/ais/delivery/http"
	aispostgres "github.com/WeCodingNow/AIS_SUG_backend/ais/repository/postgres"
	aisusecase "github.com/WeCodingNow/AIS_SUG_backend/ais/usecase"

	aisauthusecase "github.com/WeCodingNow/AIS_SUG_backend/ais_auth/usecase"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"

	echoLog "github.com/labstack/gommon/log"

	_ "github.com/lib/pq"
)

type App struct {
	e *echo.Echo

	authUC    auth.UseCase
	aisUC     ais.UseCase
	aisAuthUC ais_auth.UseCase
}

const (
	apiURL = "*"
)

func (a *App) Run(port string) error {
	e := echo.New()
	e.Debug = true

	e.Logger.SetLevel(echoLog.DEBUG)
	e.Pre(middleware.RemoveTrailingSlash())
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
	aishttp.RegisterHTTPEndpoints(e, a.aisUC)

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
	authUC := authusecase.NewAuthUseCase(
		userRepo,
		viper.GetString("auth.hash_salt"),
		[]byte(viper.GetString("auth.signing_key")),
		viper.GetDuration("auth.token_ttl"),
	)

	aisRepo := aispostgres.NewAisRepository(db)
	aisUC := aisusecase.NewAisUseCase(aisRepo)

	aisAuthUC := aisauthusecase.NewAisAuthUseCase(aisUC, authUC)

	return &App{
		authUC:    authUC,
		aisUC:     aisUC,
		aisAuthUC: aisAuthUC,
	}
}
