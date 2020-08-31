package persistentmemorydriver

import (
	"fmt"

	"github.com/apiles/unikv"
)

// DriverContext defines the driver context
type DriverContext struct {
	Filename      string `json:"filename"`
	Filemode      int    `json:"filemode"`
	CommitWhenPut bool   `json:"commit_when_put"`
	LoadWhenGet   bool   `json:"load_when_get"`
}

// Constructor constructs new memory drivers
func Constructor(prefix string, ctx unikv.DriverContextRaw) (unikv.Driver, error) {
	context := &DriverContext{
		Filemode:      0644,
		CommitWhenPut: true,
		LoadWhenGet:   false,
	}
	err := unikv.DriverLoadContext(ctx, context)
	if err != nil {
		return nil, err
	}
	return NewDriver(prefix, context)
}

// Descriptor describes memory driver
var Descriptor = &unikv.DriverDescriptor{
	Name:        "persistent-memory",
	Constructor: Constructor,
}

func init() {
	if unikv.Version != Version {
		panic(fmt.Errorf("Unmatched version with unikv, expected APIv%d", Version))
	}
	unikv.RegisterDriver(Descriptor)
}
