package redis

import (
	"time"

	"github.com/go-redis/redis/v8"
	goratelimiter "github.com/ralvarezdev/go-rate-limiter/redis"

	internalloader "github.com/ralvarezdev/uru-mobiles-recipes-api/internal/loader"
)

const (
	// EnvRedisAddress is the environment variable for the Redis address
	EnvRedisAddress = "REDIS_ADDRESS"
	
	// EnvRedisUsername is the environment variable for the Redis username
	EnvRedisUsername = "REDIS_USERNAME"
	
	// EnvRedisPassword is the environment variable for the Redis password
	EnvRedisPassword = "REDIS_PASSWORD"
	
	// EnvRedisDB is the environment variable for the Redis database number
	EnvRedisDB = "REDIS_DB"

	// EnvRateLimiterMaxRequests is the environment variable for the rate limiter max requests
	EnvRateLimiterMaxRequests = "RATE_LIMITER_MAX_REQUESTS"

	// EnvRateLimiterPeriod is the environment variable for the rate limiter period
	EnvRateLimiterPeriod = "RATE_LIMITER_PERIOD"
)

var (
	// RedisAddress is the Redis address
	RedisAddress string
	
	// RedisUsername is the Redis username
	RedisUsername string
	
	// RedisPassword is the Redis password
	RedisPassword string
	
	// RedisDB is the Redis database number
	RedisDB int

	// RateLimiterMaxRequests is the rate limiter max requests
	RateLimiterMaxRequests int

	// RateLimiterPeriod is the rate limiter period in seconds
	RateLimiterPeriod time.Duration

	// Client is the Redis client
	Client *redis.Client

	// RateLimiter is the Redis rate limiter client
	RateLimiter goratelimiter.RateLimiter
)

// Load initializes the Redis client
func Load() {
	// Get the Redis address, username and password from the environment variables
	for env, dest := range map[string]*string{
		EnvRedisAddress:  &RedisAddress,
		EnvRedisUsername: &RedisUsername,
		EnvRedisPassword: &RedisPassword,
	} {
		if err := internalloader.Loader.LoadVariable(
			env,
			dest,
		); err != nil {
			panic(err)
		}
	}

	// Get the Redis database number from the environment variable
	if err := internalloader.Loader.LoadIntVariable(
		EnvRedisDB,
		&RedisDB,
	); err != nil {
		panic(err)
	}

	// Create the redis client
	Client = redis.NewClient(
		&redis.Options{
			Addr: RedisAddress,
			Username: RedisUsername,
			Password: RedisPassword,
			DB:       RedisDB,
		},
	)

	// Load the rate limiter max requests
	if err := internalloader.Loader.LoadIntVariable(
		EnvRateLimiterMaxRequests,
		&RateLimiterMaxRequests,
	); err != nil {
		panic(err)
	}

	// Load the rate limiter period
	if err := internalloader.Loader.LoadDurationVariable(
		EnvRateLimiterPeriod,
		&RateLimiterPeriod,
	); err != nil {
		panic(err)
	}

	// Create the rate limiter
	rateLimiter, err := goratelimiter.NewDefaultRateLimiter(
		Client,
		RateLimiterMaxRequests,
		RateLimiterPeriod,
	)
	if err != nil {
		panic(err)
	}
	RateLimiter = rateLimiter
}
