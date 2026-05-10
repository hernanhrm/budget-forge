package di

import (
	"github.com/samber/do/v2"
)

type Scope = do.Scope
type Injector = do.Injector

func New() *do.RootScope {
	return do.New()
}

func NewWithOpts(opts *do.InjectorOpts) *do.RootScope {
	return do.NewWithOpts(opts)
}

func Provide[T any](i do.Injector, provider do.Provider[T]) {
	do.Provide(i, provider)
}

func ProvideNamed[T any](i do.Injector, name string, provider do.Provider[T]) {
	do.ProvideNamed(i, name, provider)
}

func ProvideValue[T any](i do.Injector, value T) {
	do.ProvideValue(i, value)
}

func ProvideNamedValue[T any](i do.Injector, name string, value T) {
	do.ProvideNamedValue(i, name, value)
}

func ProvideTransient[T any](i do.Injector, provider do.Provider[T]) {
	do.ProvideTransient(i, provider)
}

func ProvideNamedTransient[T any](i do.Injector, name string, provider do.Provider[T]) {
	do.ProvideNamedTransient(i, name, provider)
}

func Invoke[T any](i do.Injector) (T, error) {
	return do.Invoke[T](i)
}

func MustInvoke[T any](i do.Injector) T {
	return do.MustInvoke[T](i)
}

func InvokeNamed[T any](i do.Injector, name string) (T, error) {
	return do.InvokeNamed[T](i, name)
}

func MustInvokeNamed[T any](i do.Injector, name string) T {
	return do.MustInvokeNamed[T](i, name)
}

func Override[T any](i do.Injector, provider do.Provider[T]) {
	do.Override(i, provider)
}

func OverrideNamed[T any](i do.Injector, name string, provider do.Provider[T]) {
	do.OverrideNamed(i, name, provider)
}

func OverrideValue[T any](i do.Injector, value T) {
	do.OverrideValue(i, value)
}

func OverrideNamedValue[T any](i do.Injector, name string, value T) {
	do.OverrideNamedValue(i, name, value)
}

func Shutdown(i do.Injector) error {
	return i.Shutdown()
}

func HealthCheck(i do.Injector) map[string]error {
	return i.HealthCheck()
}
