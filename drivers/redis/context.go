package redisdriver

import "github.com/gomodule/redigo/redis"

// DriverContext is the context specificated for
// this driver
type DriverContext struct {
	Server     string `json:"server"`
	Options    []redis.DialOption
	RawOptions map[string]interface{} `json:"options"`
}
