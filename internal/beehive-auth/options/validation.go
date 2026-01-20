package options

func (o *Options) Validate() []error {
	var errs []error

	errs = append(errs, o.Log.Validate()...)
	errs = append(errs, o.Grpc.Validate()...)
	errs = append(errs, o.JWT.Validate()...)
	errs = append(errs, o.Redis.Validate()...)
	errs = append(errs, o.Etcd.Validate()...)
	errs = append(errs, o.Services.Validate()...)

	return errs
}
