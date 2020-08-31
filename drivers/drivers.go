// Package drivers imports all the available drivers
package drivers

import (
	// memory driver
	_ "github.com/apiles/unikv/drivers/memory"
	// persistent-memory driver
	_ "github.com/apiles/unikv/drivers/persistent-memory"
	// redis driver
	_ "github.com/apiles/unikv/drivers/redis"
)
