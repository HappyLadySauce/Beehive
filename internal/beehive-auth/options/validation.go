package options

func (o *Options) Validate() []error {
	var errs []error

	errs = append(errs, o.Log.Validate()...)
	errs = append(errs, o.Grpc.Validate()...)
	errs = append(errs, o.JWT.Validate()...)
	errs = append(errs, o.Redis.Validate()...)
	errs = append(errs, o.Etcd.Validate()...)

	// Check for nil pointer before calling Validate
	if o.UserSvc == nil {
		// If UserSvc is nil (e.g., not in config file), initialize with defaults
		o.UserSvc = NewUserServiceOptions()
	}
	errs = append(errs, o.UserSvc.Validate()...)

	return errs
}
