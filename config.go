// Package config parses out a config struct based on environment variables.
//
// It has support for the following types:
// - string
// - time.Duration (using time.ParseDuration)
// - bool (using strconv.ParseBool)
//
// Config is used externally as a simple global function: Parse()
// Adding values to config involves modifying this package directly; it is for this reason that this package is internal.
//
// For more information on struct fields, see Config.
package config

import (
	"github.com/pkg/errors"
	"reflect"
)

// Parse takes a struct pointer and parses out env tags in order to populate it from the environment.
func Parse(i interface{}) error {
	v := reflect.ValueOf(i)
	if !valid(v) {
		return errors.Errorf("config: value given should be a struct pointer, %s given", v.Type().String())
	}

	// grab all fields from the given struct along with their required env tags
	f := fields(v.Elem().Type())

	// next, grab all environment tags so we can fully parse the tree.
	// we deliberately do multiple passes here, as this is a startup cost that people are opting into.
	// multiple passes also means the logic here remains readable.
	e := getenv(env(f))

	// next, we go through and set values for every field; validation comes later.
	return e.populate(v.Elem(), f, false)
}
