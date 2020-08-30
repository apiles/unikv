package memorydriver

import (
	"fmt"

	"github.com/apiles/unikv"
)

// Constructor constructs new memory drivers
func Constructor(prefix string, ctx unikv.DriverContextRaw) (unikv.Driver, error) {
	return NewDriver(), nil
}

// Descriptor describes memory driver
var Descriptor = &unikv.DriverDescriptor{
	Name:        "memory",
	Constructor: Constructor,
}

func init() {
	if unikv.Version != Version {
		panic(fmt.Errorf("Unmatched version with unikv, expected APIv%d", Version))
	}
	unikv.RegisterDriver(Descriptor)
}
