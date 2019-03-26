package config

import (
	"os"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

// env will grab a slice of all available env tags.
func env(t []*field) []string {
	m := make([]string, 0)

	for _, f := range t {
		if len(f.subfields) > 0 {
			m = append(m, env(f.subfields)...)
			continue
		}

		if f.env != "" {
			m = append(m, f.env)
		}
	}

	return m
}

// getenv will grab every environment key at once based on all struct tags available on the given object.
func getenv(e []string) envMap {
	m := make(map[string]string)

	for _, key := range e {
		m[key] = os.Getenv(key)
	}

	return m
}

// envMap is a map of environment variables to values that supports validation via the required:"" struct tag.
type envMap map[string]string

// required returns true if a value is actually required according to the environment.
func (e envMap) required(req string) bool {
	if req == "" {
		return false
	}

	if req == "true" {
		return true
	}

	r := strings.SplitN(req, "=", 2)
	if len(r) != 2 {
		return false
	}

	key, val := r[0], r[1]
	actual, ok := e[key]
	if !ok {
		return false
	}

	return actual == val
}

func (e envMap) populate(v reflect.Value, f []*field, force bool) error {
	for i, field := range f {
		value := v.Field(i)

		r := e.required(field.required) || force

		if len(field.subfields) > 0 {
			err := e.populate(v.Index(i), field.subfields, r)
			if err != nil {
				return err
			}
		}

		if field.env == "" {
			continue
		}

		val, ok := e[field.env]

		if !ok {
			// if this field is required but not set, error.
			if r {
				return errors.Errorf("config: env variable '%s' is required but not set", field.env)
			}

			// otherwise, skip this field (it's not set and not required, so it remains the nil value).
			continue
		}

		if !allowed(field.in, val) {
			return errors.Errorf("config: env variable '%s' must be one of (%s); %s given", field.env, strings.Join(field.in, ", "), val)
		}

		err := set(value, val, field.env)
		if err != nil {
			return err
		}
	}

	return nil
}

// simple loop over a slice of strings to check if a given search is there.
func allowed(allowed []string, search string) bool {
	if len(allowed) == 0 {
		return true
	}

	for _, a := range allowed {
		if search == a {
			return true
		}
	}

	return false
}
