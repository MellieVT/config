package config

import (
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// field is a holder for internal state while parsing.
type field struct {
	f        reflect.StructField
	env      string
	required string
	in       []string

	// subfields are used for nested structs.
	subfields []*field
}

// valid checks that a reflected struct is usable (aka it's a struct pointer)
func valid(v reflect.Value) bool {
	if v.IsNil() || v.Kind() != reflect.Ptr {
		return false
	}

	return v.Elem().Kind() == reflect.Struct
}

// fields pulls out all fields with their associated tags from a reflected type.
// If v is not a struct, this will panic.
func fields(v reflect.Type) []*field {
	j := v.NumField()
	r := make([]*field, j)

	for i := 0; i < j; i++ {
		f := v.Field(i)

		// if we have a nested struct, we ignore the "env" value and only care about the required value and subfields.
		if f.Type.Kind() == reflect.Struct {
			r[i] = &field{
				f:        f,
				required: f.Tag.Get("required"),

				// TODO(amelia): guard against infinite recursion here?
				subfields: fields(f.Type),
			}
			continue
		}

		r[i] = &field{
			f:        f,
			env:      f.Tag.Get("env"),
			required: f.Tag.Get("required"),
			in:       strings.Split(f.Tag.Get("in"), ","),
		}
	}

	return r
}

func set(value reflect.Value, str, key string) error {
	switch value.Kind() {
	case reflect.Bool:
		b, err := strconv.ParseBool(str)
		if err != nil {
			return errors.Wrapf(err, "config: failed to parse bool from env variable '%s' (contained: '%s')", key, str)
		}

		value.SetBool(b)

	case reflect.String:
		value.SetString(str)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// we are explicitly a time.Duration here.
		// hint: time.Time and time.Duration are int64s.
		// this is the neatest way to check for a specific type in this case!
		if t := value.Type(); t.PkgPath() == "time" && t.Name() == "Duration" {
			d, err := time.ParseDuration(str)
			if err != nil {
				return errors.Wrapf(err, "config: failed to parse time.Duration from env variable '%s' (contained: '%s')", key, str)
			}

			value.SetInt(int64(d))

			return nil
		}

		bits := value.Type().Bits()
		i, err := strconv.ParseInt(str, 10, bits)
		if err != nil {
			return errors.Wrapf(err, "config: failed to parse int from env variable '%s' (contained: '%s', bitsize: %d)", key, str, bits)
		}

		value.SetInt(i)

	case reflect.Float32, reflect.Float64:
		bits := value.Type().Bits()
		i, err := strconv.ParseFloat(str, bits)
		if err != nil {
			return errors.Wrapf(err, "config: failed to parse float from env variable '%s' (contained: '%s', bitsize: %d)", key, str, bits)
		}

		value.SetFloat(i)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		bits := value.Type().Bits()
		i, err := strconv.ParseUint(str, 0, bits)
		if err != nil {
			return errors.Wrapf(err, "config: failed to parse uint from env variable '%s' (contained: '%s', bitsize: %d)", key, str, bits)
		}

		value.SetUint(i)
	}

	return nil
}
