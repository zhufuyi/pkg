package tracer

// alias, for other structs, the following code does not need to change the names of the options
type options = resourceConfig

// Option modifying struct field values by means of an interface
type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (o optionFunc) apply(cfg *options) {
	o(cfg)
}

// set obj fields value
func apply(obj *options, opts ...Option) {
	for _, opt := range opts {
		opt.apply(obj)
	}
}

// WithServiceName set service name
func WithServiceName(name string) Option {
	return optionFunc(func(o *options) {
		o.serviceName = name
	})
}

// WithServiceVersion set service version
func WithServiceVersion(version string) Option {
	return optionFunc(func(o *options) {
		o.serviceVersion = version
	})
}

// WithEnvironment set service environment
func WithEnvironment(environment string) Option {
	return optionFunc(func(o *options) {
		o.environment = environment
	})
}

// WithAttributes set service attributes
func WithAttributes(attributes map[string]string) Option {
	return optionFunc(func(o *options) {
		o.attributes = attributes
	})
}
