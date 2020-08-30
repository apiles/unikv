package memorydriver

import "github.com/apiles/unikv"

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
	unikv.RegisterDriver(Descriptor)
}
