package initiator

import (
	"context"
	"fmt"
	"time"

	"github.com/alazarbeyenenew2/cves_backend/platform/logger"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/viper"
)

func initDB(log logger.Logger) *pgxpool.Pool {
	var (
		config *pgxpool.Config
		err    error
	)
	dbUser := viper.GetString("db.user")
	dbPassword := viper.GetString("db.password")
	dbhost := viper.GetString("db.host")
	dbport := viper.GetInt("db.port")
	dbsslMode := viper.GetString("db.ssl_mode")
	dbName := viper.GetString("db.name")
	dbURL := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s", dbUser, dbPassword, dbhost, dbport, dbName, dbsslMode)
	config, err = pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Error(context.Background(), "unable to parse pgxpool config string for cves")
		log.Fatal(context.Background(), err.Error())
	}

	idleConnTimeout := viper.GetDuration("database.idle_conn_timeout")
	if idleConnTimeout == 0 {
		idleConnTimeout = 4 * time.Minute
	}
	config.MaxConnIdleTime = idleConnTimeout

	conn, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatal(context.Background(), fmt.Sprintf("failed to connect to database (%s): %v", dbURL, err))
	}

	log.Info(context.Background(), fmt.Sprintf("connected to %s database successfully", dbURL))
	return conn
}
