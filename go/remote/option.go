package remote

type loadOptions struct {
	insecure bool
	username string
	password string
}

// LoadOption is the option that can be used for loading an
// remote.Image object.
type LoadOption func(*loadOptions)

func WithAuth(username, password string) LoadOption {
	return func(o *loadOptions) {
		o.username = username
		o.password = password
	}
}

func WithInsecure() LoadOption {
	return func(o *loadOptions) {
		o.insecure = true
	}
}
