package config

import (
	"errors"
	"fmt"
	"log"
	"reflect"
)

// Settable is an interface over all conf types supporting Set
type Settable interface {
	Set(variable []string, value interface{}) error
}

// SetField sets field of c with path within c to given value.
func SetField(c Settable, varpath []string, value interface{}) error {
	// c must be a pointer to a struct
	rv := reflect.ValueOf(c)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		return errors.New("c must be pointer to struct")
	}

	// Dereference pointer
	rv = rv.Elem()
	fv := rv.FieldByName(varpath[0])
	for _, name := range varpath[:len(varpath)] {
		// Lookup field by name
		fv = rv.FieldByName(name)
		if !fv.IsValid() {
			return fmt.Errorf("not a field name: %s", name)
		}

		// Field must be exported
		if !fv.CanSet() {
			return fmt.Errorf("cannot set field %s", name)
		}

		rv = fv
	}

	defer func() { // Catch possible panic generated by fv.Set
		if err := recover(); err != nil {
			log.Print("Error: could not convert value to type of variable")
		}
	}()

	fv.Set(reflect.ValueOf(value).Convert(fv.Type()))

	return nil
}