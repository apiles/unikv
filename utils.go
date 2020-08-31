package unikv

import "encoding/json"

// ConcatPrefix concats two prefixes together
func ConcatPrefix(prefix string, str string) string {
	if prefix != "" {
		return prefix + configure.Separator + str
	}
	return str
}

func concatPrefix(prefix string, str string) string {
	return ConcatPrefix(prefix, str)
}

type temporaryStringWriter struct {
	buffer string
}

func (tsw *temporaryStringWriter) Write(p []byte) (n int, err error) {
	tsw.buffer += string(p)
	return len(p), nil
}

// DriverLoadContext loads raw context into a structure
func DriverLoadContext(ctx DriverContextRaw, dest interface{}) error {
	byt, err := json.Marshal(ctx)
	if err != nil {
		return err
	}
	return json.Unmarshal(byt, dest)
}
