package unikv

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

func init() {
	loadConfigure()
}
