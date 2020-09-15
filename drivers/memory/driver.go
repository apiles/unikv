package memorydriver

import (
	"reflect"
	"sync"

	"github.com/apiles/unikv"
)

// Driver is the memorydriver
type Driver struct {
	data *sync.Map
}

// Get gets data
func (d *Driver) Get(key string) (string, error) {
	data, ok := d.data.Load(key)
	if !ok {
		return "", unikv.ErrNotFound
	}
	return data.(string), nil
}

// Put puts data
func (d *Driver) Put(key string, value string) error {
	d.data.Store(key, value)
	return nil
}

// Unset unsets data
func (d *Driver) Unset(key string) error {
	d.data.Delete(key)
	return nil
}

// List lists the keys
func (d *Driver) List() (interface{}, error) {
	return reflect.ValueOf(d.data).MapKeys(), nil
}

// Close closes driver
func (d *Driver) Close() error {
	return nil
}

// NewDriver creates a driver
func NewDriver() *Driver {
	return &Driver{
		data: new(sync.Map),
	}
}
