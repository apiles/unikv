package persistentmemorydriver

import (
	"encoding/gob"
	"os"
	"reflect"
	"sync"

	"github.com/apiles/unikv"
)

// Driver is the memorydriver
type Driver struct {
	data          map[string]string
	lock          *sync.Mutex
	filename      string
	filemode      int
	commitWhenPut bool
	loadWhenGet   bool
	prefix        string
}

// Commit commits changes to file
func (d *Driver) Commit() error {
	f, err := os.OpenFile(d.filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.FileMode(d.filemode))
	if err != nil {
		return err
	}
	enc := gob.NewEncoder(f)
	return enc.Encode(&d.data)
}

// Load loads data from file
func (d *Driver) Load() error {
	f, err := os.Open(d.filename)
	if err != nil {
		return err
	}
	dec := gob.NewDecoder(f)
	return dec.Decode(&d.data)
}

// Get gets data
func (d *Driver) Get(key string) (string, error) {
	key = unikv.ConcatPrefix(d.prefix, key)
	if d.loadWhenGet {
		err := d.Load()
		if err != nil {
			return "", err
		}
	}
	data, ok := d.data[key]
	if !ok {
		return "", unikv.ErrNotFound
	}
	return data, nil
}

// Put puts data
func (d *Driver) Put(key string, value string) error {
	key = unikv.ConcatPrefix(d.prefix, key)
	d.data[key] = value
	if d.commitWhenPut {
		return d.Commit()
	}
	return nil
}

// List lists the keys
func (d *Driver) List() (interface{}, error) {
	return reflect.ValueOf(d.data).MapKeys(), nil
}

// Unset unsets data
func (d *Driver) Unset(key string) error {
	key = unikv.ConcatPrefix(d.prefix, key)
	delete(d.data, key)
	if d.commitWhenPut {
		return d.Commit()
	}
	return nil
}

// Close closes driver
func (d *Driver) Close() error {
	return d.Commit()
}

// NewDriver creates a driver
func NewDriver(prefix string, ctx *DriverContext) (*Driver, error) {
	f, err := os.Open(ctx.Filename)
	if os.IsNotExist(err) {
		f, err = os.Create(ctx.Filename)
	}
	if err != nil {
		return nil, err
	}
	f.Close()
	drv := &Driver{
		data:          make(map[string]string),
		lock:          &sync.Mutex{},
		filename:      ctx.Filename,
		filemode:      ctx.Filemode,
		commitWhenPut: ctx.CommitWhenPut,
		loadWhenGet:   ctx.LoadWhenGet,
		prefix:        prefix,
	}
	drv.Load()
	return drv, nil
}
