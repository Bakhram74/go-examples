package postgres

import "time"

type Option func(*Postgres)

func WithMaxPoolSize(size int) Option {
	return func(c *Postgres) {
		c.maxPoolSize = size
	}
}

func WithConnAttempts(attempts int) Option {
	return func(c *Postgres) {
		c.connAttempts = attempts
	}
}

func WithConnTimeout(timeout time.Duration) Option {
	return func(c *Postgres) {
		c.connTimeout = timeout
	}
}
