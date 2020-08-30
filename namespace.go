package unikv

// Namespace is the namespace type
type Namespace struct {
	Name   string
	Prefix string
}

// NewNamespace creates a new namespace
func NewNamespace(name string) *Namespace {
	conf, ok := configure.Namespaces[name]
	prefix := name
	if ok {
		if conf.Prefix != "" {
			prefix = conf.Prefix
		}
		if conf.Prefix == "$!empty" {
			prefix = ""
		}
	}
	return &Namespace{
		Name:   name,
		Prefix: prefix,
	}
}

// NewBucket creates a new bucket on a namespace
func (ns *Namespace) NewBucket() {

}
