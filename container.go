package di

import (
	"errors"
	"fmt"
)

type Container struct {
	Dependencies map[string]Dependency
	Instances    map[string]interface{}
}

func NewContainer() *Container {
	container := &Container{}
	container.Dependencies = make(map[string]Dependency)
	container.Instances = make(map[string]interface{})

	di := Dependency{Name: "di", Callback: func(params ...interface{}) interface{} { return container }}

	container.Dependencies[di.Name] = di

	return container
}

func (container *Container) Set(dependency Dependency) *Container {
	if dependency.Name != "di" {
		if container.Instances[dependency.Name] != "" {
			delete(container.Instances, dependency.Name)
		}

		container.Dependencies[dependency.Name] = dependency
	}

	return container
}

func (container *Container) Inject(dependency Dependency, fresh bool) (interface{}, error) {
	if _, ok := container.Instances[dependency.Name]; ok && !fresh {
		return container.Instances[dependency.Name], nil
	}

	var arguments []interface{}

	for i := 0; i < len(dependency.Dependencies); i++ {
		if _, ok := container.Instances[dependency.Dependencies[i]]; ok {
			arguments = append(arguments, container.Instances[dependency.Dependencies[i]])
			continue
		}

		if _, ok := container.Dependencies[dependency.Dependencies[i]]; !ok {
			return container.Instances[dependency.Dependencies[i]], errors.New(fmt.Sprintf("Failed to find dependency %s", dependency.Dependencies[i]))
		}
		d, _ := container.Get(dependency.Dependencies[i])
		arguments = append(arguments, d)
	}

	results := dependency.Callback(arguments...)

	container.Instances[dependency.Name] = results

	return results, nil
}

func (container *Container) Get(name string) (interface{}, error) {
	if _, ok := container.Dependencies[name]; !ok {
		return nil, errors.New(fmt.Sprintf("Failed to find dependency %s", name))
	}

	return container.Inject(container.Dependencies[name], false)
}

func (container *Container) Refresh(name string) *Container {
	if _, ok := container.Dependencies[name]; ok {
		delete(container.Dependencies, name)
	}
	return container
}

func (container *Container) Has(name string) bool {
	_, ok := container.Dependencies[name]
	return ok
}
