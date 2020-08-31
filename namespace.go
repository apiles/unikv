package unikv

import "fmt"

// Namespace is the namespace type
type Namespace struct {
	Name    string
	Prefix  string
	buckets map[string]*Bucket
}

// NewNamespace creates a new namespace
func NewNamespace(name string) *Namespace {
	conf, ok := GetConfigure().Namespaces[name]
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
		Name:    name,
		Prefix:  prefix,
		buckets: make(map[string]*Bucket),
	}
}

// NewBucket creates a new bucket on a namespace
func (ns *Namespace) NewBucket(name string) (*Bucket, error) {
	if _, ok := ns.buckets[name]; ok {
		return ns.buckets[name], nil
	}
	bckconf, ok := configure.Namespaces[ns.Name]
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
		namespace:     ns,
	}
	if _, ok := drivers[conf.Driver]; !ok {
		return nil, fmt.Errorf("Unknow driver %s", conf.Driver)
	}
	drv, err := drivers[conf.Driver].Construct(prefix, conf.Context)
	if err != nil {
		return nil, err
	}
	rslt.Driver = drv
	ns.buckets[name] = rslt
	return rslt, nil
}

// Close closes all the buckets in a namespace
func (ns *Namespace) Close() {
	for k := range ns.buckets {
		ns.buckets[k].Close()
	}
}
