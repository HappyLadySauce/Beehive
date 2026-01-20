package options

import (
	"encoding/json"
	"fmt"

	"github.com/HappyLadySauce/Beehive/internal/pkg/options"
	"github.com/spf13/pflag"
	"k8s.io/component-base/cli/flag"
)

type Options struct {
	Name      string
	Log       *options.LogOptions   `json:"log" mapstructure:"log"`
	Grpc      *options.GrpcOptions  `json:"grpc" mapstructure:"grpc"`
	Etcd      *options.EtcdOptions  `json:"etcd" mapstructure:"etcd"`
	InsecureServing    *options.InsecureServingOptions 	`json:"insecure-serving" mapstructure:"insecure-serving"`
	WebSocket *options.WebSocketOptions     `json:"websocket" mapstructure:"websocket"`
	Services  *MicroServicesOptions `json:"services" mapstructure:"services"`
}


// MicroServicesOptions 微服务地址配置
type MicroServicesOptions struct {
	AuthServiceAddr     string `json:"auth-service-addr" mapstructure:"auth-service-addr"`
	UserServiceAddr     string `json:"user-service-addr" mapstructure:"user-service-addr"`
	MessageServiceAddr  string `json:"message-service-addr" mapstructure:"message-service-addr"`
	PresenceServiceAddr string `json:"presence-service-addr" mapstructure:"presence-service-addr"`
}

// NewMicroServicesOptions 创建微服务配置选项
func NewMicroServicesOptions() *MicroServicesOptions {
	return &MicroServicesOptions{
		AuthServiceAddr:     "etcd://beehive-auth",
		UserServiceAddr:     "etcd://beehive-user",
		MessageServiceAddr:  "etcd://beehive-message",
		PresenceServiceAddr: "etcd://beehive-presence",
	}
}

// Validate 验证微服务配置
func (o *MicroServicesOptions) Validate() []error {
	var errs []error
	if o.AuthServiceAddr == "" {
		errs = append(errs, fmt.Errorf("services.auth_service_addr cannot be empty"))
	}
	if o.UserServiceAddr == "" {
		errs = append(errs, fmt.Errorf("services.user_service_addr cannot be empty"))
	}
	if o.MessageServiceAddr == "" {
		errs = append(errs, fmt.Errorf("services.message_service_addr cannot be empty"))
	}
	if o.PresenceServiceAddr == "" {
		errs = append(errs, fmt.Errorf("services.presence_service_addr cannot be empty"))
	}
	return errs
}

// AddFlags 添加命令行标志
func (o *MicroServicesOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.AuthServiceAddr, "services.auth-service-addr", o.AuthServiceAddr, "Auth Service gRPC address (supports etcd:// prefix)")
	fs.StringVar(&o.UserServiceAddr, "services.user-service-addr", o.UserServiceAddr, "User Service gRPC address (supports etcd:// prefix)")
	fs.StringVar(&o.MessageServiceAddr, "services.message-service-addr", o.MessageServiceAddr, "Message Service gRPC address (supports etcd:// prefix)")
	fs.StringVar(&o.PresenceServiceAddr, "services.presence-service-addr", o.PresenceServiceAddr, "Presence Service gRPC address (supports etcd:// prefix)")
}

func NewOptions(basename string) *Options {
	return &Options{
		Name:      basename,
		Log:       options.NewLogOptions(),
		Grpc:      options.NewGrpcOptions(),
		Etcd:      options.NewEtcdOptions(),
		InsecureServing:    options.NewInsecureServingOptions(),
		WebSocket: options.NewWebSocketOptions(),
		Services:  NewMicroServicesOptions(),
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

	// add insecure serving flags to the NamedFlagSets
	insecureServingFlagSet := nfs.FlagSet("Insecure Serving")
	o.InsecureServing.AddFlags(insecureServingFlagSet)

	// add websocket flags to the NamedFlagSets
	websocketFlagSet := nfs.FlagSet("WebSocket")
	o.WebSocket.AddFlags(websocketFlagSet)

	// add services flags to the NamedFlagSets
	servicesFlagSet := nfs.FlagSet("Micro Services")
	o.Services.AddFlags(servicesFlagSet)

	return nfs
}

func (o *Options) String() string {
	data, _ := json.Marshal(o)

	return string(data)
}
