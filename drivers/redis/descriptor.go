package redisdriver

import (
	"fmt"
	"strconv"
	"time"

	"github.com/apiles/unikv"
	"github.com/gomodule/redigo/redis"
)

// Constructor constructs new redis drivers
func Constructor(prefix string, ctx unikv.DriverContextRaw) (unikv.Driver, error) {
	context := &DriverContext{}
	err := unikv.DriverLoadContext(ctx, context)
	for k, v := range context.RawOptions {
		var ok bool
		switch k {
		case "username":
			if _, ok = v.(string); !ok {
				return nil, fmt.Errorf("wrong context argument: username")
			}
			context.Options = append(context.Options, redis.DialUsername(v.(string)))
		case "password":
			if _, ok = v.(string); !ok {
				return nil, fmt.Errorf("wrong context argument: password")
			}
			context.Options = append(context.Options, redis.DialPassword(v.(string)))
		case "keepalive":
			var str string
			switch v.(type) {
			case string:
				str = v.(string)
			case int:
				str = strconv.Itoa(v.(int))
			}
			t, err := time.ParseDuration(str)
			if err != nil {
				context.Options = append(context.Options, redis.DialKeepAlive(t))
			}
		case "tls":
			if _, ok = v.(bool); !ok {
				return nil, fmt.Errorf("wrong context argument: tls")
			}
			context.Options = append(context.Options, redis.DialUseTLS(v.(bool)))
		case "database":
			if _, ok = v.(int); !ok {
				return nil, fmt.Errorf("wrong context argument: database")
			}
			context.Options = append(context.Options, redis.DialDatabase(v.(int)))
		}
	}
	if err != nil {
		return nil, err
	}
	return NewDriver(context)
}

// Descriptor describes memory driver
var Descriptor = &unikv.DriverDescriptor{
	Name:        "redis",
	Constructor: Constructor,
}

func init() {
	if unikv.Version != Version {
		panic(fmt.Errorf("Unmatched version with unikv, expected APIv%d", Version))
	}
	unikv.RegisterDriver(Descriptor)
}
