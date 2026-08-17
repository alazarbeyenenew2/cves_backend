package initiator

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alazarbeyenenew2/cves_backend/docs"
	_ "github.com/alazarbeyenenew2/cves_backend/docs"
	"github.com/alazarbeyenenew2/cves_backend/internal/constant/db/dbinterface"
	middleware "github.com/alazarbeyenenew2/cves_backend/internal/handler/middleware"
	"github.com/alazarbeyenenew2/cves_backend/platform/logger"
	"github.com/alazarbeyenenew2/cves_backend/platform/workerpool"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

func Initiate() {
	ctx := context.Background()
	docs.SwaggerInfo.Title = "cves API"
	docs.SwaggerInfo.Description = "API documentation for PulseChannels"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/"

	log, err := zap.NewProduction()
	if err != nil {
		log.Fatal("unable to start logging")
	}
	configName := "config"
	if os.Getenv("CONFIG_NAME") != "" {
		configName = os.Getenv("CONFIG_NAME")
	}

	log.Named("initialization").Info("initializing config")
	err = InitConfig(Config{Names: []string{configName}, Path: "config", Logger: log})
	if err != nil {
		log.Fatal("unable to start config", zap.Error(err))
	}
	log.Named("initialization").Info("initializing config completed")

	log.Named("initialization").Info("initializing logger")
	logger := logger.New(log)
	log.Named("initialization").Info("initializing logger completed")

	log.Named("initialization").Info("initializing database connection")
	pgxPool := initDB(logger)
	log.Named("initialization").Info("initializing database connection completed")

	log.Named("initialization").Info("initializing workpool")
	wp := workerpool.New(viper.GetInt("workerpool.max_workers"), viper.GetInt("workerpool.task_buffer"))
	wp.Start()
	log.Named("initialization").Info("initializing workpool completed")

	log.Named("initialization").Info("initializing redis")
	redisPool := initRedis(viper.GetString("redis.url"), logger)
	log.Named("initialization").Info("initializing redis completed")

	log.Named("initialization").Info("initializing persistence layer")
	persistenceDB := dbinterface.New(pgxPool, logger)
	persistence := initPersistence(&persistenceDB, redisPool, logger)
	log.Named("initialization").Info("initializing persistence layer completed")

	log.Named("initialization").Info("initializing platform layer")
	platform := initPlatform(logger)
	log.Named("initialization").Info("initializing platform completed")

	log.Named("initialization").Info("initializing module layer")
	module := initModule(persistence, logger, wp, platform)
	module.NVD.Start()
	module.NVD.TriggerNow()
	log.Named("initialization").Info("initializing module layer completed")

	log.Named("initialization").Info("initializing handler layer")
	handler := initHandler(module, logger, wp)
	log.Named("initialization").Info("initializing handler layer completed")

	log.Named("initialization").Info("initializing http Engine")
	server := gin.New()
	server.Use(middleware.GinLogger(logger))
	server.Use(middleware.CORS())
	log.Named("initialization").Info("initializing Engine completed")
	ginsrv := server.Group("")
	swaggerUser := viper.GetString("swagger.username")
	swaggerPass := viper.GetString("swagger.password")
	if swaggerUser != "" && swaggerPass != "" {
		swaggerGroup := ginsrv.Group("/swagger", gin.BasicAuth(gin.Accounts{
			swaggerUser: swaggerPass,
		}))
		swaggerGroup.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	logger.Info(ctx, "initializing route")
	initRoute(ginsrv, handler, module, logger, viper.GetString("auth.jwt_secret"), persistence)
	logger.Info(ctx, "done initializing route")

	srv := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", viper.GetString("app.host"), viper.GetInt("app.port")),
		Handler:           server,
		ReadHeaderTimeout: viper.GetDuration("app.timeout"),
		IdleTimeout:       30 * time.Minute,
	}

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		log.Fatal("HTTP server Shutdown")
		wp.Stop()
	}()

	host := fmt.Sprint(viper.GetString("app.host"), ":", viper.GetInt("app.port"))
	logger.Info(ctx, "server listening at port ", zap.Any("link", host))
	err = srv.ListenAndServe()
	if err != nil {
		logger.Fatal(ctx, fmt.Sprintf("Could not start HTTP server: %s", err))
	}

}
