package unikv

import "encoding/json"

func concatPrefix(prefix string, str string) string {
	if prefix != "" {
		return prefix + PrefixSeparator + str
	}
	return str
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
