package redis

import (
	"time"

	"github.com/go-redis/redis/v8"
	goratelimiter "github.com/ralvarezdev/go-rate-limiter/redis"

	internalloader "github.com/ralvarezdev/uru-mobiles-recipes-api/internal/loader"
)

const (
	// EnvRedisURL is the environment variable for the Redis URL
	EnvRedisURL = "REDIS_URL"

	// EnvRateLimiterMaxRequests is the environment variable for the rate limiter max requests
	EnvRateLimiterMaxRequests = "RATE_LIMITER_MAX_REQUESTS"

	// EnvRateLimiterPeriod is the environment variable for the rate limiter period
	EnvRateLimiterPeriod = "RATE_LIMITER_PERIOD"
)

var (
	// RedisURL is the Redis URL
	RedisURL string

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
	// Get the Redis URL from the environment variables
	if err := internalloader.Loader.LoadVariable(
		EnvRedisURL,
		&RedisURL,
	); err != nil {
		panic(err)
	}

	// Create the redis client
	Client = redis.NewClient(
		&redis.Options{
			Addr: RedisURL,
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
