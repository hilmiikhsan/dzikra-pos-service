package adapter

// import "codebase-app/internal/pkg/validator"
func WithValidator(v Validator) Option {
	return func(a *Adapter) {
		a.Validator = v
	}
}
