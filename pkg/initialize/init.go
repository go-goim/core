package initialize

import (
	"fmt"
)

type Initializer interface {
	// BeforeInit is called before application init.
	BeforeInit() error
	// BeforeRun is called before application run.
	BeforeRun() error
}

var (
	initializers []Initializer
)

func Register(initializer Initializer) {
	initializers = append(initializers, initializer)
}

func BeforeInit() error {
	for _, i := range initializers {
		if err := i.BeforeInit(); err != nil {
			return err
		}
	}
	return nil
}

func BeforeRun() error {
	for _, i := range initializers {
		if err := i.BeforeRun(); err != nil {
			return err
		}
	}
	return nil
}

type InitializerFunc func() error

type basicInitializer struct {
	name       string
	beforeInit InitializerFunc
	beforeRun  InitializerFunc
}

func (i *basicInitializer) BeforeInit() error {
	if i.beforeInit != nil {
		if err := i.beforeInit(); err != nil {
			return fmt.Errorf("%s before init error: %w", i.name, err)
		}
	}
	return nil
}

func (i *basicInitializer) BeforeRun() error {
	if i.beforeRun != nil {
		if err := i.beforeRun(); err != nil {
			return fmt.Errorf("%s after init error: %w", i.name, err)
		}
	}
	return nil
}

func NewBasicInitializer(name string, beforeInit InitializerFunc, afterInit InitializerFunc) Initializer {
	return &basicInitializer{
		name:       name,
		beforeInit: beforeInit,
		beforeRun:  afterInit,
	}
}
