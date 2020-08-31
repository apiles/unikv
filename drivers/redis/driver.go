package redisdriver

import (
	"fmt"

	"github.com/apiles/unikv"
	"github.com/gomodule/redigo/redis"
)

// Driver is the memorydriver
type Driver struct {
	conn   redis.Conn
	prefix string
}

// Get gets data
func (d *Driver) Get(key string) (string, error) {
	key = unikv.ConcatPrefix(d.prefix, key)
	data, err := d.conn.Do("GET", key)
	if err != nil {
		if err == redis.ErrNil {
			return "", unikv.ErrNotFound
		}
		return "", err
	}
	if _, ok := data.([]byte); !ok {
		if data == nil {
			return "", unikv.ErrNotFound
		} else {
			return "", fmt.Errorf("Unknown result %v", data)
		}
	}
	return string(data.([]byte)), nil
}

// Put puts data
func (d *Driver) Put(key string, value string) error {
	key = unikv.ConcatPrefix(d.prefix, key)
	_, err := d.conn.Do("SET", key, value)
	return err
}

// Unset unsets data
func (d *Driver) Unset(key string) error {
	key = unikv.ConcatPrefix(d.prefix, key)
	_, err := d.conn.Do("SET", key, "EX", "1")
	return err
}

// Close closes driver
func (d *Driver) Close() error {
	return d.conn.Close()
}

// NewDriver creates a driver
func NewDriver(prefix string, ctx *DriverContext) (*Driver, error) {
	drv := &Driver{
		prefix: prefix,
	}
	var err error
	drv.conn, err = redis.DialURL(ctx.Server, ctx.Options...)
	if err != nil {
		return nil, err
	}
	return drv, nil
}
