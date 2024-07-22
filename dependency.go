package di

import (
	"errors"
	"fmt"
	"slices"
)

type Dependency struct {
	Dependencies []string
	Name         string
	Callback     func(params ...interface{}) interface{}
}

func NewDependency(name string, callback func(params ...interface{}) interface{}) Dependency {
	return Dependency{Name: name, Callback: callback}
}

func (injection *Dependency) Inject(name string) (*Dependency, error) {
	if slices.Contains(injection.Dependencies, name) {
		return injection, errors.New(fmt.Sprintf("Dependency already declared for %s", name))
	}

	injection.Dependencies = append(injection.Dependencies, name)

	return injection, nil
}
