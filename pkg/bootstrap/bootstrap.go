package bootstrap

import (
	"fmt"
	"go.uber.org/dig"
)

var Container *dig.Container

func init() {
	Container = dig.New()
}

func Provide(function interface{}, opts ...dig.ProvideOption) {
	err := Container.Provide(function, opts...)
	if err != nil {
		fmt.Println(err)
	}
}
