package unikv

func concatPrefix(prefix string, str string) string {
	if prefix != "" {
		return prefix + PrefixSeparator + str
	}
	return str
}
