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
	UserSvc *UserServiceOptions   `json:"user-service" mapstructure:"user-service"`
}

// UserServiceOptions User Service gRPC 客户端配置
type UserServiceOptions struct {
	Addr string `json:"addr" mapstructure:"addr"`
}

// NewUserServiceOptions 创建 User Service 配置选项
func NewUserServiceOptions() *UserServiceOptions {
	return &UserServiceOptions{
		Addr: "localhost:50051",
	}
}

// Validate 验证 User Service 配置
func (o *UserServiceOptions) Validate() []error {
	var errs []error
	if o.Addr == "" {
		errs = append(errs, fmt.Errorf("user-service.addr cannot be empty"))
	}
	return errs
}

// AddFlags 添加命令行标志
func (o *UserServiceOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Addr, "user-service.addr", o.Addr, "User Service gRPC address")
}

func NewOptions(basename string) *Options {
	return &Options{
		Name:    basename,
		Log:     options.NewLogOptions(),
		Grpc:    options.NewGrpcOptions(),
		JWT:     options.NewJWTOptions(),
		Redis:   options.NewRedisOptions(),
		Etcd:    options.NewEtcdOptions(),
		UserSvc: NewUserServiceOptions(),
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

	// add user service flags to the NamedFlagSets
	userSvcFlagSet := nfs.FlagSet("User Service")
	o.UserSvc.AddFlags(userSvcFlagSet)

	return nfs
}

func (o *Options) String() string {
	data, _ := json.Marshal(o)

	return string(data)
}
