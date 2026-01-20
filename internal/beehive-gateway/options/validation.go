package options

func (o *Options) Validate() []error {
	var errs []error

	errs = append(errs, o.Log.Validate()...)
	errs = append(errs, o.Grpc.Validate()...)
	errs = append(errs, o.Etcd.Validate()...)
	errs = append(errs, o.InsecureServing.Validate()...)
	errs = append(errs, o.WebSocket.Validate()...)
	errs = append(errs, o.Services.Validate()...)

	return errs
}
