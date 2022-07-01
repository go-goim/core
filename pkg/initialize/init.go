package initialize

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-goim/core/pkg/errors"
	"github.com/go-goim/core/pkg/log"
)

type Initializer interface {
	// BeforeInit is called before application init.
	BeforeInit(context.Context) error
	// BeforeRun is called before application run.
	BeforeRun(context.Context) error
}

type InitializerFunc func(context.Context) error

var (
	// defaultTimeout is the default timeout for initializer.
	defaultTimeout = time.Second * 5
	initializers   []Initializer
)

func Register(initializer Initializer) {
	initializers = append(initializers, initializer)
}

func BeforeInit(ctx context.Context) error {
	log.Info("start run BeforeInit")
	var cancel context.CancelFunc
	if ctx == nil {
		ctx = context.Background()
		ctx, cancel = context.WithTimeout(ctx, defaultTimeout)
		defer cancel()
	}

	var (
		done = make(chan struct{})
		errs = make(errors.ErrorSet, 0)
		wg   sync.WaitGroup
	)

	for _, i := range initializers {
		wg.Add(1)
		go func(f InitializerFunc) {
			defer wg.Done()
			if err := f(ctx); err != nil {
				errs = append(errs, err)
			}
		}(i.BeforeInit)
	}

	go func() {
		wg.Wait()
		done <- struct{}{}
	}()

	select {
	case <-ctx.Done(): // timeout
		return ctx.Err()
	case <-done: // shutdown complete
		return errs.Err()
	}
}

func BeforeRun(ctx context.Context) error {
	log.Info("start run BeforeRun")
	var cancel context.CancelFunc
	if ctx == nil {
		ctx = context.Background()
		ctx, cancel = context.WithTimeout(ctx, defaultTimeout)
		defer cancel()
	}

	var (
		done = make(chan struct{})
		errs = make(errors.ErrorSet, 0)
		wg   sync.WaitGroup
	)

	for _, i := range initializers {
		wg.Add(1)
		go func(f InitializerFunc) {
			defer wg.Done()
			if err := f(ctx); err != nil {
				errs = append(errs, err)
			}
		}(i.BeforeRun)
	}

	go func() {
		wg.Wait()
		done <- struct{}{}
	}()

	select {
	case <-ctx.Done(): // timeout
		return ctx.Err()
	case <-done: // shutdown complete
		return errs.Err()
	}
}

type basicInitializer struct {
	name       string
	beforeInit InitializerFunc
	beforeRun  InitializerFunc
}

func (i *basicInitializer) BeforeInit(ctx context.Context) error {
	if i.beforeInit != nil {
		if err := i.beforeInit(ctx); err != nil {
			return fmt.Errorf("%s before init error: %w", i.name, err)
		}
	}
	return nil
}

func (i *basicInitializer) BeforeRun(ctx context.Context) error {
	if i.beforeRun != nil {
		if err := i.beforeRun(ctx); err != nil {
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
