package tests

import (
	"fmt"
	"github.com/utopia-go/di"
	"testing"
)

func TestDi(t *testing.T) {
	container := di.NewContainer()

	user := di.Dependency{
		Name: "user",
		Callback: func(params ...interface{}) interface{} {
			return fmt.Sprintf("John doe is %d years old", params[0])
		},
	}

	user.Inject("age")

	container.Set(user)
	container.Set(di.NewDependency("age", func(params ...interface{}) interface{} {
		return 25
	}))

	if userObject, err := container.Get("user"); err == nil {
		got := userObject
		want := "John doe is 25 years old"

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}
