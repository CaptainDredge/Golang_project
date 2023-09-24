package main

import (
	"context"
	"encoding/json"
	"fmt"
	"gin-rest-api-example/internal/account"
	accountDB "gin-rest-api-example/internal/account/database"
	"gin-rest-api-example/internal/article"
	articleDB "gin-rest-api-example/internal/article/database"
	"gin-rest-api-example/internal/cache"
	"gin-rest-api-example/internal/config"
	"gin-rest-api-example/internal/database"
	"gin-rest-api-example/internal/game"
	gameDB "gin-rest-api-example/internal/game/database"
	"gin-rest-api-example/internal/metric"
	"gin-rest-api-example/internal/middleware"
	"gin-rest-api-example/pkg/logging"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var serverCmd = &cobra.Command{
	Use: "server",
	Run: func(cmd *cobra.Command, args []string) {
		runApplication()
	},
}

func runApplication() {
	// load configs and sets default logger configs.
	conf, err := config.Load(configFile)
	if err != nil {
		log.Fatal(err)
	}
	logging.SetConfig(&logging.Config{
		Encoding:    conf.LoggingConfig.Encoding,
		Level:       zapcore.Level(conf.LoggingConfig.Level),
		Development: conf.LoggingConfig.Development,
	})
	defer logging.DefaultLogger().Sync()

	// setup application(di + run server)
	app := fx.New(
		fx.Supply(conf),
		fx.Supply(logging.DefaultLogger().Desugar()),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log.Named("fx")}
		}),
		fx.StopTimeout(conf.ServerConfig.GracefulShutdown+time.Second),
		fx.Invoke(
			printAppInfo,
		),
		fx.Provide(
			metric.NewMetricsProvider,
			// setup database
			database.NewDatabase,
			// setup cache
			cache.NewCacher,
			// setup account packages
			accountDB.NewAccountDB,
			account.NewAuthMiddleware,
			account.NewHandler,
			// setup article packages
			articleDB.NewArticleDB,
			article.NewHandler,
			// setup game packages
			gameDB.NewGameDB,
			game.NewHandler,
			// server
			newServer,
		),
		fx.Invoke(
			account.RouteV1,
			article.RouteV1,
			game.RouteV1,
			func(r *gin.Engine) {},
		),
	)
	app.Run()
}

func newServer(lc fx.Lifecycle, cfg *config.Config, mp *metric.MetricsProvider) *gin.Engine {
	gin.SetMode(gin.DebugMode)
	r := gin.New()
	r.Use(middleware.LoggingMiddleware("/metric"), gin.Recovery())

	metric.Route(r)
	r.Use(metric.MetricsMiddleware(mp))

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.ServerConfig.Port),
		Handler:      r,
		ReadTimeout:  cfg.ServerConfig.ReadTimeout,
		WriteTimeout: cfg.ServerConfig.WriteTimeout,
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logging.FromContext(ctx).Infof("Start to rest api server :%d", cfg.ServerConfig.Port)
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logging.DefaultLogger().Errorw("failed to close http server", "err", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logging.FromContext(ctx).Info("Stopped rest api server")
			return srv.Shutdown(ctx)
		},
	})
	return r
}

func printAppInfo(cfg *config.Config) {
	b, _ := json.MarshalIndent(&cfg, "", "  ")
	logging.DefaultLogger().Infof("application information\n%s", string(b))
}
