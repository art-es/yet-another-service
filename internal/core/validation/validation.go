//go:generate mockgen -source=validation.go -destination=mock/validation.go -package=mock
package validation

type Validator interface {
	Struct(s interface{}) error
	Var(field interface{}, tag string) error
}
