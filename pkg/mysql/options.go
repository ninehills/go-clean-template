package mysql

import "time"

// Option -.
type Option func(*MySQL)

// ConnMaxLifetime -.
func ConnMaxLifetime(seconds int) Option {
	return func(c *MySQL) {
		c.connMaxLifetime = time.Second * time.Duration(seconds)
	}
}

// MaxOpenConns -.
func MaxOpenConns(conns int) Option {
	return func(c *MySQL) {
		c.maxOpenConns = conns
	}
}

// MaxIdleConns -.
func MaxIdleConns(conns int) Option {
	return func(c *MySQL) {
		c.maxIdleConns = conns
	}
}
