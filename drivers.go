package unikv

import "fmt"

// Driver declares a unikv storage driver
type Driver interface {
	Get(key string) (string, error)
	Put(key string, data string) error
	Unset(key string) error
	Close() error
}

// DriverContextRaw contains the raw map of driver configuration
// passed in within unikv configure
type DriverContextRaw map[string]interface{}

// DriverDescriptor is used for storing meta information
// of the driver
type DriverDescriptor struct {
	Name        string
	Constructor func(prefix string, ctx DriverContextRaw) (Driver, error)
}

// Construct a new driver instance
func (ds *DriverDescriptor) Construct(prefix string, ctx DriverContextRaw) (Driver, error) {
	drv, err := ds.Constructor(prefix, ctx)
	// TODO: add specific error processing logic
	return drv, err
}

var drivers map[string]*DriverDescriptor = make(map[string]*DriverDescriptor)

// RegisterDriver registers a new driver
func RegisterDriver(driver *DriverDescriptor) {
	drivers[driver.Name] = driver
}

// ErrNotFound is thrown when trying to access
// an unset key
var ErrNotFound = fmt.Errorf("ErrNotFound")

// IsErrNotFound checks if an error is unikv.ErrNotFound
func IsErrNotFound(err error) bool {
	return err == ErrNotFound
}
