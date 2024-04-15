package app

import (
	"context"
	"os"
	"os/signal"
	"single-window/config"
	v1 "single-window/internal/controller/http/v1"
	"single-window/internal/controller/http/v1/middlewares"
	"single-window/internal/entity"
	"single-window/pkg/accessor"
	"single-window/pkg/httpserver"
	"single-window/pkg/jwt"
	"single-window/pkg/logger"
	"single-window/pkg/orm"
	"single-window/pkg/postgres"
	"single-window/pkg/s3"
	"single-window/pkg/swcontext"
	"syscall"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	middleware "github.com/oapi-codegen/gin-middleware"
)

func Run(cfg *config.Config, l *logger.Logger) {
	db, conn := initDB(cfg, l)
	defer conn.Close()

	// TODO: Будет ли потом для чего-то исопльзоваться context для s3?
	minio, err := s3.New(context.Background(), &cfg.MINIO, l)
	if err != nil {
		l.Fatal("failed to create s3, err: %w", err)
	}

	swagger, err := v1.GetSwagger()
	if err != nil {
		l.Fatal("failed to load swagger spec, err: %v", err)
	}

	enforcer, err := accessor.NewEnforcerFromGORM(db, cfg.Casbin.ConfigPath)
	if err != nil {
		l.Fatal("failed to create casbin enforcer, err: %v", err)
	}

	parser, err := jwt.NewParser([]byte(cfg.JWT.PublicKey))
	if err != nil {
		l.Fatal("failed to create jwt parser, err: %v", err)
	}
	swCtxAdapter := swcontext.NewGinContextCookieAdapter[entity.Claims]()
	auth := middlewares.NewAuthenticator(parser, swCtxAdapter)

	routes, err := v1.NewRoutes(db, minio, cfg, l, swCtxAdapter, accessor.NewCasbinAccessor(enforcer))
	if err != nil {
		l.Fatal("failed to create routes, err: %v", err)
	}

	r := gin.Default()
	if cfg.App.Mode == "local" {
		// TODO: добавить возможность конфигурировать cors через конфиги или энвы
		corsCfg := cors.DefaultConfig()
		corsCfg.AllowAllOrigins = false
		corsCfg.AllowOrigins = append(corsCfg.AllowOrigins, "http://localhost:3000")
		corsCfg.AddAllowHeaders("deviceId", "X-Real-IP")
		corsCfg.AllowCredentials = true
		r.Use(cors.New(corsCfg))
	}
	r.Use(middleware.OapiRequestValidatorWithOptions(swagger, &middleware.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: auth.OAPIAuthenticate,
		},
	}))
	v1.RegisterHandlers(r, routes)

	httpServer := httpserver.New(r, httpserver.WithPort(cfg.HTTP.Port))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %s", s.String())
	case err = <-httpServer.Notify():
		l.Error("app - Run - httpServer.Notify: %v", err)
	}

	err = httpServer.Shutdown()
	if err != nil {
		l.Error("app - Run - httpServer.Shutdown: %v", err)
	}
}

func initDB(cfg *config.Config, l *logger.Logger) (*gorm.DB, *postgres.Postgres) {
	conn, err := postgres.New(cfg.PG.URL, postgres.WithMaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal("app - Run - postgres.New: %v", err)
	}

	db, err := orm.NewPostgresORM(conn)
	if err != nil {
		l.Fatal("failed to initialize gorm, err: %v", err)
	}

	return db, conn
}
