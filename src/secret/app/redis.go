package app

import "github.com/go-redis/redis"

// Redis for redis client connection details
type Redis struct {
	client   *redis.Client
	address  string
	password string
	db       int
}

// NewRedis creates an instance of redis which will be used
func NewRedis(addr, password string, db int) Redis {
	return Redis{
		address:  addr,
		password: password,
		db:       db,
	}
}

// RedisHandlerFunc for handling request to redis
type RedisHandlerFunc func(*redis.Client) error

// Redis ensures that connection to redis exists and
// calls callback function with redis client
func (a *App) Redis(f RedisHandlerFunc) error {
	if a.redis.client == nil {
		a.redis.connect()
	}
	// the redis client can be abstracted and
	// function could receive just interface (list of available function)
	// but in that case we need to implement a wrapper for the client
	// which I think is not important for the task
	return f(a.redis.client)
}

func (r *Redis) connect() {
	r.client = redis.NewClient(&redis.Options{
		Addr:     r.address,
		Password: r.password,
		DB:       r.db,
	})
}
