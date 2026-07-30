package initiator

import (
	"context"
	"fmt"
	"time"

	"github.com/alazarbeyenenew2/cves_backend/platform/logger"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

func initRedis(redisUrl string, log logger.Logger) *redis.Client {
	// Parse the Redis URL to get connection options
	opt, err := redis.ParseURL(redisUrl)
	if err != nil {
		log.Error(context.Background(), "unable to parse redis config string")
		// Fatal will terminate the application if the config is invalid
		log.Fatal(context.Background(), err.Error())
	}

	// 1. Set connection pool size (Max connections)
	// viper.GetInt is used as it typically represents a count
	poolSize := viper.GetInt("redis.pool_size")
	if poolSize == 0 {
		// Default to a reasonable pool size, e.g., 10 or 20
		poolSize = 10
	}
	opt.PoolSize = poolSize

	// 2. Set Idle timeout (Max connection idle time)
	// viper.GetDuration is used for time settings
	idleTimeout := viper.GetDuration("redis.idle_timeout")
	if idleTimeout == 0 {
		// Default to 5 minutes, after which idle connections are closed
		idleTimeout = 5 * time.Minute
	}
	opt.IdleTimeout = idleTimeout

	// Initialize the Redis client
	rdb := redis.NewClient(opt)

	// Ping the server to verify the connection is established
	// A small timeout is used for the connection test
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	status := rdb.Ping(ctx)
	if status.Err() != nil {
		// Terminate the application if the connection test fails
		log.Fatal(context.Background(), fmt.Sprintf("Failed to connect to Redis: %v", status.Err()))
	}

	log.Info(context.Background(), "Successfully connected to Redis")
	return rdb
}
