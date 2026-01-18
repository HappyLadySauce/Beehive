package options

func (o *Options) Validate() []error {
	var errs []error

	errs = append(errs, o.Log.Validate()...)
	errs = append(errs, o.Grpc.Validate()...)
	errs = append(errs, o.Postgresql.Validate()...)
	errs = append(errs, o.Etcd.Validate()...)

	return errs
}
