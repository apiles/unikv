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
func (ns *Namespace) NewBucket(name string) (*Bucket, error) {
	bckconf, ok := configure.Namespaces[name]
	var conf *ConfigureBuckets
	prefix := concatPrefix(ns.Prefix, name)
	if ok {
		conf, ok = bckconf.Buckets[name]
	}
	if !ok {
		conf = &ConfigureBuckets{
			Prefix:  "",
			Driver:  configure.Default.Driver,
			Context: configure.Default.Context,
		}
	}
	if conf.Prefix != "" {
		prefix = concatPrefix(ns.Prefix, conf.Prefix)
	}
	if conf.Prefix == "$!empty" {
		prefix = ns.Prefix
	}
	rslt := &Bucket{
		Name:          name,
		Prefix:        prefix,
		NamespaceName: ns.Name,
	}
	drv, err := drivers[conf.Driver].Construct(prefix, conf.Context)
	if err != nil {
		return nil, err
	}
	rslt.Driver = drv
	return rslt, nil
}
