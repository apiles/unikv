package redisdriver

import (
	"fmt"

	"github.com/apiles/unikv"
	"github.com/gomodule/redigo/redis"
)

// Driver is the memorydriver
type Driver struct {
	conn        redis.Conn
	prefix      string
	ctx         *DriverContext
	reconnected bool
}

// Get gets data
func (d *Driver) Get(key string) (string, error) {
	key = unikv.ConcatPrefix(d.prefix, key)
	data, err := d.conn.Do("GET", key)
	if err != nil {
		if err == redis.ErrNil {
			return "", unikv.ErrNotFound
		}
		if !d.reconnected {
			d.connect()
			d.reconnected = true
			return d.Get(key)
		}
		return "", err
	}
	if d.reconnected {
		d.reconnected = false
	}
	if _, ok := data.([]byte); !ok {
		if data == nil {
			return "", unikv.ErrNotFound
		}
		return "", fmt.Errorf("Unknown result %v", data)
	}
	return string(data.([]byte)), nil
}

// Put puts data
func (d *Driver) Put(key string, value string) error {
	key = unikv.ConcatPrefix(d.prefix, key)
	_, err := d.conn.Do("SET", key, value)
	if err != nil {
		if !d.reconnected {
			d.connect()
			d.reconnected = true
			return d.Put(key, value)
		}
	}
	if d.reconnected {
		d.reconnected = false
	}
	return err
}

// Unset unsets data
func (d *Driver) Unset(key string) error {
	key = unikv.ConcatPrefix(d.prefix, key)
	_, err := d.conn.Do("SET", key, "EX", "1")
	if err != nil {
		if !d.reconnected {
			d.connect()
			d.reconnected = true
			return d.Unset(key)
		}
	}
	if d.reconnected {
		d.reconnected = false
	}
	return err
}

// List lists the keys
func (d *Driver) List() (interface{}, error) {
	data, err := d.conn.Do("KEYS", d.prefix+"*")
	if err != nil {
		if err == redis.ErrNil {
			return "", unikv.ErrNotFound
		}
		if !d.reconnected {
			d.connect()
			d.reconnected = true
			return d.List()
		}
		return "", err
	}
	if d.reconnected {
		d.reconnected = false
	}
	dat, ok := data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("Invalid response from redis server")
	}
	var rslt []string = make([]string, len(dat))
	for k, v := range dat {
		rslt[k] = unikv.TrimPrefix(d.prefix, string(v.([]byte)))
	}
	return rslt, err
}

// Close closes driver
func (d *Driver) Close() error {
	return d.conn.Close()
}

func (d *Driver) connect() error {
	if d.conn != nil {
		d.conn.Close()
	}
	var err error
	d.conn, err = redis.DialURL(d.ctx.Server, d.ctx.Options...)
	return err
}

// NewDriver creates a driver
func NewDriver(prefix string, ctx *DriverContext) (*Driver, error) {
	drv := &Driver{
		prefix: prefix,
		ctx:    ctx,
	}
	var err error
	err = drv.connect()
	if err != nil {
		return nil, err
	}
	return drv, nil
}
