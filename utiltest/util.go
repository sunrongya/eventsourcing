package utiltest

import (
	"github.com/stretchr/testify/assert"
	es "github.com/sunrongya/eventsourcing"
	"reflect"
	"testing"
	"time"
)

type mockService struct {
	command es.Command
}

func (m *mockService) HandleCommands() {
}

func (m *mockService) CommandChannel() chan<- es.Command {
	return nil
}

func (m *mockService) PublishCommand(command es.Command) error {
	if err := checkCommand(command); err != nil {
		return err
	}
	m.command = command
	return nil
}

func (m *mockService) RestoreAggregate(guid es.Guid) es.Aggregate {
	return nil
}

func TestServicePublishCommand(t *testing.T, fn func(es.Service) es.Command) {
	mockS := new(mockService)
	command := fn(mockS)
	assert.NotNil(t, mockS.command, "")
	assert.Equal(t, command, mockS.command)
}

type CommandFieldError struct {
	Field string
}

func (c CommandFieldError) Error() string {
	return "missing field: " + c.Field
}

func checkCommand(command es.Command) error {
	rv := reflect.Indirect(reflect.ValueOf(command))
	rt := rv.Type()

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		if field.PkgPath != "" {
			continue // Skip private field.
		}

		tag := field.Tag.Get("eh")
		if tag == "optional" {
			continue // Optional field.
		}

		if isZero(rv.Field(i)) {
			return CommandFieldError{field.Name}
		}
	}
	return nil
}

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Func, reflect.Map, reflect.Slice:
		return v.IsNil()
	case reflect.Struct:
		// Special case to get zero values by method.
		switch obj := v.Interface().(type) {
		case time.Time:
			return obj.IsZero()
		}

		// Check public fields for zero values.
		z := true
		for i := 0; i < v.NumField(); i++ {
			if v.Type().Field(i).PkgPath != "" {
				continue // Skip private fields.
			}
			z = z && isZero(v.Field(i))
		}
		return z
	case reflect.Bool:
		return false
	}

	// Compare other types directly:
	z := reflect.Zero(v.Type())
	return v.Interface() == z.Interface()
}
