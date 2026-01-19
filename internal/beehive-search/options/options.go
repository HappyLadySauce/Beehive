package options

import (
	"encoding/json"

	"github.com/HappyLadySauce/Beehive/internal/pkg/options"
	"github.com/spf13/pflag"
	"k8s.io/component-base/cli/flag"
)

type Options struct {
	Name string
	Log  *options.LogOptions  `json:"log" mapstructure:"log"`
	Grpc *options.GrpcOptions `json:"grpc" mapstructure:"grpc"`
	Etcd *options.EtcdOptions `json:"etcd" mapstructure:"etcd"`
}

func NewOptions(basename string) *Options {
	return &Options{
		Name: basename,
		Log:  options.NewLogOptions(),
		Grpc: options.NewGrpcOptions(),
		Etcd: options.NewEtcdOptions(),
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

	// add grpc flags to the NamedFlagSets
	grpcFlagSet := nfs.FlagSet("gRPC")
	o.Grpc.AddFlags(grpcFlagSet)

	// add etcd flags to the NamedFlagSets
	etcdFlagSet := nfs.FlagSet("Etcd")
	o.Etcd.AddFlags(etcdFlagSet)

	return nfs
}

func (o *Options) String() string {
	data, _ := json.Marshal(o)

	return string(data)
}
