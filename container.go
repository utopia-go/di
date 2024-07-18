package di

import (
	"errors"
	"fmt"
)

type Container struct {
	dependencies map[string]Dependency
	instances    map[string]interface{}
}

func NewContainer() *Container {
	container := &Container{}

	di := Dependency{name: "di", callback: func(params ...interface{}) interface{} { return container }}

	container.dependencies[di.name] = di

	return container
}

func (container *Container) Set(dependency Dependency) *Container {
	if dependency.name != "di" {
		if container.instances[dependency.name] != "" {
			delete(container.instances, dependency.name)
		}

		container.dependencies[dependency.name] = dependency
	}

	return container
}

func (container *Container) Inject(dependency Dependency, fresh bool) (interface{}, error) {
	if container.instances[dependency.name] != "" && !fresh {
		return container.instances[dependency.name], nil
	}

	var arguments []interface{}

	for i := 0; i < len(dependency.dependencies); i++ {
		if _, ok := container.instances[dependency.dependencies[i]]; ok {
			arguments = append(arguments, container.instances[dependency.dependencies[i]])
			continue
		}

		if _, ok := container.dependencies[dependency.dependencies[i]]; !ok {
			return container.instances[dependency.dependencies[i]], errors.New(fmt.Sprintf("Failed to find dependency %s", dependency.dependencies[i]))
		}
		d, _ := container.Get(dependency.dependencies[i])
		arguments = append(arguments, d)
	}

	results := dependency.callback(arguments...)

	container.instances[dependency.name] = results

	return results, nil
}

func (container *Container) Get(name string) (interface{}, error) {
	if _, ok := container.dependencies[name]; !ok {
		return nil, errors.New(fmt.Sprintf("Failed to find dependency %s", name))
	}

	return container.Inject(container.dependencies[name], false)
}

func (container *Container) Refresh(name string) *Container {
	if _, ok := container.dependencies[name]; ok {
		delete(container.dependencies, name)
	}
	return container
}

func (container *Container) Has(name string) bool {
	_, ok := container.dependencies[name]
	return ok
}
