package unikv

// Key is the type of unikv key
type Key interface {
	// Convert to string
	String() string
}

// KeyString is string key
type KeyString string

// String converts KeyString to string
func (ks KeyString) String() string {
	return string(ks)
}

// KeyNil is blanck key
var KeyNil = KeyString("(unikv_nil)")

// NewKey creates a new key
func NewKey(key interface{}) Key {
	switch key.(type) {
	case string:
		return KeyString(key.(string))
	case Key:
		return key.(Key)
	}
	return KeyNil
}
