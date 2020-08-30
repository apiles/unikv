package unikv

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

// Configure is unikv's configure structure
type Configure struct {
	Default    *ConfigureDefault               `yaml:"default"`
	Namespaces map[string]*ConfigureNamespaces `yaml:"namespaces"`
}

// ConfigureDefault is the default segment
// in unikv's configure
type ConfigureDefault struct {
	Driver  string           `yaml:"driver"`
	Context DriverContextRaw `yaml:"context"`
}

// ConfigureNamespaces is the namespace segment
// in unikv's configure
type ConfigureNamespaces struct {
	Prefix  string                       `yaml:"prefix"`
	Buckets map[string]*ConfigureBuckets `yaml:"buckets"`
}

// ConfigureBuckets is the bucket segment
// in unikv's configure
type ConfigureBuckets struct {
	Prefix  string           `yaml:"prefix"`
	Driver  string           `yaml:"driver"`
	Context DriverContextRaw `yaml:"context"`
}

func determineConfigureFilePath() string {
	rslt, ok := os.LookupEnv(ConfigureEnvName)
	if ok == false {
		rslt = DefaultConfigureFile
	}
	return rslt
}

var configure *Configure

func loadConfigure() {
	content, err := ioutil.ReadFile(determineConfigureFilePath())
	if err != nil {
		panic(err)
	}
	configure = &Configure{
		Default: &ConfigureDefault{
			Driver:  DefaultDriver,
			Context: make(DriverContextRaw),
		},
		Namespaces: make(map[string]*ConfigureNamespaces),
	}
	err = yaml.Unmarshal(content, configure)
	if err != nil {
		panic(err)
	}
}
