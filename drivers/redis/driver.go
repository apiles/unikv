package redisdriver

import (
	"github.com/apiles/unikv"
	"github.com/gomodule/redigo/redis"
)

// Driver is the memorydriver
type Driver struct {
	conn redis.Conn
}

// Get gets data
func (d *Driver) Get(key string) (string, error) {
	data, err := d.conn.Do("GET", key)
	if err != nil {
		if err == redis.ErrNil {
			return "", unikv.ErrNotFound
		}
		return "", err
	}
	return data.(string), nil
}

// Put puts data
func (d *Driver) Put(key string, value string) error {
	_, err := d.conn.Do("SET", key, value)
	return err
}

// Unset unsets data
func (d *Driver) Unset(key string) error {
	_, err := d.conn.Do("SET", key, "EX", "1")
	return err
}

// Close closes driver
func (d *Driver) Close() error {
	return nil
}

// NewDriver creates a driver
func NewDriver(ctx *DriverContext) (*Driver, error) {
	drv := &Driver{}
	var err error
	drv.conn, err = redis.DialURL(ctx.Server, ctx.Options...)
	if err != nil {
		return nil, err
	}
	return drv, nil
}
