package di

import (
	"errors"
	"fmt"
	"slices"
)

type Dependency struct {
	dependencies []string
	name         string
	callback     func(params ...interface{}) interface{}
}

func (injection *Dependency) inject(name string) (*Dependency, error) {
	if slices.Contains(injection.dependencies, name) {
		return injection, errors.New(fmt.Sprintf("Dependency already declared for %s", name))
	}

	injection.dependencies = append(injection.dependencies, name)

	return injection, nil
}
