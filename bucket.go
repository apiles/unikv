package unikv

import (
	"encoding/gob"
	"strconv"
	"strings"
)

// Bucket is a unikv store bucket
type Bucket struct {
	Name          string
	Prefix        string
	NamespaceName string
	Driver        Driver
}

// GetString gets string value
func (b *Bucket) GetString(key interface{}) (string, error) {
	return b.Driver.Get(concatPrefix(b.Prefix, NewKey(key).String()))
}

// PutString puts string value
func (b *Bucket) PutString(key interface{}, str string) error {
	return b.Driver.Put(concatPrefix(b.Prefix, NewKey(key).String()), str)
}

// GetInt gets int value
func (b *Bucket) GetInt(key interface{}) (int, error) {
	str, err := b.GetString(key)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(str)
}

// PutInt puts int value
func (b *Bucket) PutInt(key interface{}, value int) error {
	return b.PutString(key, strconv.Itoa(value))
}

// Get gets value into dest
func (b *Bucket) Get(key interface{}, dest interface{}) error {
	str, err := b.GetString(key)
	if err != nil {
		return err
	}
	dec := gob.NewDecoder(strings.NewReader(str))
	return dec.Decode(dest)
}

// Put puts value
func (b *Bucket) Put(key interface{}, value interface{}) error {
	tmpWriter := new(temporaryStringWriter)
	enc := gob.NewEncoder(tmpWriter)
	err := enc.Encode(value)
	if err != nil {
		return err
	}
	return b.PutString(key, tmpWriter.buffer)
}

// Close closes a bucket
func (b *Bucket) Close() error {
	return b.Driver.Close()
}
