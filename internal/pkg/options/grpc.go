package options

import (
	"fmt"

	"github.com/spf13/pflag"
)

// GrpcOptions is the options for the gRPC server.
type GrpcOptions struct {
	BindAddress string `json:"bind-address" mapstructure:"bind-address"`
	BindPort    int    `json:"bind-port" mapstructure:"bind-port"`
	MaxMsgSize  int    `json:"max-msg-size" mapstructure:"max-msg-size"`
}

// NewGrpcOptions creates a new GrpcOptions.
func NewGrpcOptions() *GrpcOptions {
	return &GrpcOptions{
		BindAddress: "0.0.0.0",
		BindPort:    50050, // 默认端口，可以在使用时覆盖
		MaxMsgSize:  4 * 1024 * 1024,
	}
}

// Validate validates the GrpcOptions.
func (o *GrpcOptions) Validate() []error {
	var errors []error

	if o.BindPort < 0 || o.BindPort > 65535 {
		errors = append(
			errors,
			fmt.Errorf(
				"--insecure-port %v must be between 0 and 65535, inclusive. 0 for turning off insecure (HTTP) port",
				o.BindPort,
			),
		)
	}
	return errors
}

// AddFlags adds the flags to the specified FlagSet.
func (o *GrpcOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.BindAddress, "grpc.bind-address", o.BindAddress, ""+
		"The IP address on which to serve the --grpc.bind-port(set to 0.0.0.0 for all IPv4 interfaces and :: for all IPv6 interfaces).")

	fs.IntVar(&o.BindPort, "grpc.bind-port", o.BindPort, ""+
		"The port on which to serve unsecured, unauthenticated grpc access. It is assumed "+
		"that firewall rules are set up such that this port is not reachable from outside of "+
		"the deployed machine and that port 443 on the iam public address is proxied to this "+
		"port. This is performed by nginx in the default setup. Set to zero to disable.")

	fs.IntVar(&o.MaxMsgSize, "grpc.max-msg-size", o.MaxMsgSize, "gRPC max message size.")
}
