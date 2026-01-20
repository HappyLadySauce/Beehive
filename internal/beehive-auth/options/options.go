package options

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/pflag"
	"k8s.io/component-base/cli/flag"

	"github.com/HappyLadySauce/Beehive/internal/pkg/options"
)

type Options struct {
	Name    string
	Log     *options.LogOptions   `json:"log" mapstructure:"log"`
	Grpc    *options.GrpcOptions  `json:"grpc" mapstructure:"grpc"`
	JWT     *options.JWTOptions   `json:"jwt" mapstructure:"jwt"`
	Redis   *options.RedisOptions `json:"redis" mapstructure:"redis"`
	Etcd    *options.EtcdOptions  `json:"etcd" mapstructure:"etcd"`
	Services  *MicroServicesOptions `json:"services" mapstructure:"services"`
}


// MicroServicesOptions 微服务地址配置
type MicroServicesOptions struct {
	UserServiceAddr     string `json:"user-service-addr" mapstructure:"user-service-addr"`
}

// NewMicroServicesOptions 创建微服务配置选项
func NewMicroServicesOptions() *MicroServicesOptions {
	return &MicroServicesOptions{
		UserServiceAddr:     "etcd://beehive-user",
	}
}

// Validate 验证微服务配置
func (o *MicroServicesOptions) Validate() []error {
	var errs []error
	if o.UserServiceAddr == "" {
		errs = append(errs, fmt.Errorf("services.user_service_addr cannot be empty"))
	}
	return errs
}

// AddFlags 添加命令行标志
func (o *MicroServicesOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.UserServiceAddr, "services.user-service-addr", o.UserServiceAddr, "User Service gRPC address (supports etcd:// prefix)")
}

func NewOptions(basename string) *Options {
	return &Options{
		Name:    basename,
		Log:     options.NewLogOptions(),
		Grpc:    options.NewGrpcOptions(),
		JWT:     options.NewJWTOptions(),
		Redis:   options.NewRedisOptions(),
		Etcd:    options.NewEtcdOptions(),
		Services: NewMicroServicesOptions(),
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

	// add jwt flags to the NamedFlagSets
	jwtFlagSet := nfs.FlagSet("JWT")
	o.JWT.AddFlags(jwtFlagSet)

	// add redis flags to the NamedFlagSets
	redisFlagSet := nfs.FlagSet("Redis")
	o.Redis.AddFlags(redisFlagSet)

	// add etcd flags to the NamedFlagSets
	etcdFlagSet := nfs.FlagSet("Etcd")
	o.Etcd.AddFlags(etcdFlagSet)

	// add services flags to the NamedFlagSets
	servicesFlagSet := nfs.FlagSet("Services")
	o.Services.AddFlags(servicesFlagSet)

	return nfs
}

func (o *Options) String() string {
	data, _ := json.Marshal(o)

	return string(data)
}
