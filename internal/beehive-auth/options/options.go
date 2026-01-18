package options

import (
	"encoding/json"
	"github.com/spf13/pflag"
	"k8s.io/component-base/cli/flag"
	"github.com/HappyLadySauce/Beehive/internal/pkg/options"
)

type Options struct {
	Name 	string
	Log		*options.LogOptions `json:"log" mapstructure:"log"`
}


func NewOptions(basename string) *Options {
	return &Options{
		Name: basename,
		Log: options.NewLogOptions(),
	}
}

// AddFlags adds the flags to the specified FlagSet and returns the grouped flag sets.
func (o *Options) AddFlags(fs *pflag.FlagSet) *flag.NamedFlagSets {
	nfs := &flag.NamedFlagSets{}

	// add config flags to the NamedFlagSets
	configFS := nfs.FlagSet("Config")
	options.AddConfigFlag(o.Name, configFS)

	// add log flags to the NamedFlagSets
	logsFlagSet := nfs.FlagSet("Logs")
	o.Log.AddFlags(logsFlagSet)

	return nfs
}

func (o *Options) String() string {
	data, _ := json.Marshal(o)

	return string(data)
}