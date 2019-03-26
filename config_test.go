package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type embedded1 struct {
	String  string  `env:"EMBEDDED_TEST1_STRING"`
	Bool    string  `env:"EMBEDDED_TEST1_BOOL"`
	Float32 float32 `env:"EMBEDDED_TEST1_FLOAT32"`
}

type embedded2 struct {
	String  string  `env:"EMBEDDED_TEST2_STRING"`
	Bool    string  `env:"EMBEDDED_TEST2_BOOL"`
	Float32 float32 `env:"EMBEDDED_TEST2_FLOAT32"`
}

type embedded3 struct {
	String  string  `env:"EMBEDDED_TEST3_STRING"`
	Bool    string  `env:"EMBEDDED_TEST3_BOOL"`
	Float32 float32 `env:"EMBEDDED_TEST3_FLOAT32"`
}

type kitchenSink struct {
	String                 string        `env:"KITCHEN_SINK_STRING"`
	Bool                   bool          `env:"KITCHEN_SINK_BOOL"`
	Int                    int           `env:"KITCHEN_SINK_INT"`
	Int8                   int8          `env:"KITCHEN_SINK_INT8"`
	Int16                  int16         `env:"KITCHEN_SINK_INT16"`
	Int32                  int32         `env:"KITCHEN_SINK_INT32"`
	Int64                  int64         `env:"KITCHEN_SINK_INT64"`
	Uint                   uint          `env:"KITCHEN_SINK_UINT"`
	Uint8                  uint8         `env:"KITCHEN_SINK_UINT8"`
	Uint16                 uint16        `env:"KITCHEN_SINK_UINT16"`
	Uint32                 uint32        `env:"KITCHEN_SINK_UINT32"`
	Uint64                 uint64        `env:"KITCHEN_SINK_UINT64"`
	Float32                float32       `env:"KITCHEN_SINK_FLOAT32"`
	Float64                float64       `env:"KITCHEN_SINK_FLOAT64"`
	Duration               time.Duration `env:"KITCHEN_SINK_DURATION"`
	RequiredBool           bool          `env:"KITCHEN_SINK_REQUIRED_BOOL" required:"true"`
	AllowedString          string        `env:"KITCHEN_SINK_ALLOWED_STRING" in:"foobar,racket,badger"`
	Embedded               embedded1
	EmbeddedRequired       embedded2 `required:"KITCHEN_SINK_REQUIRED=true"`
	EmbeddedSimpleRequired embedded3 `required:"true"`
}

type simple struct {
	String   string        `env:"SIMPLE_STRING"`
	Bool     bool          `env:"SIMPLE_BOOL"`
	Int      int           `env:"SIMPLE_INT"`
	Uint     uint          `env:"SIMPLE_UINT"`
	Float32  float32       `env:"SIMPLE_FLOAT32"`
	Float64  float64       `env:"SIMPLE_FLOAT64"`
	Duration time.Duration `env:"SIMPLE_DURATION"`
}

func TestParse(t *testing.T) {
	vars := map[string]string{
		"SIMPLE_STRING": "test string",
		"SIMPLE_BOOL": "true",
		"SIMPLE_INT": "-123456",
		"SIMPLE_UINT": "5605",
		"SIMPLE_FLOAT32": "14.56",
		"SIMPLE_FLOAT64": "123456.789",
		"SIMPLE_DURATION": "30s",
	}

	for key, val := range vars {
		_ = os.Setenv(key, val)
	}

	var c simple

	err := Parse(&c)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "test string", c.String, "c.String")
	assert.True(t, c.Bool, "c.Bool")
	assert.Equal(t, -123456, c.Int, "c.Int")
	assert.Equal(t, uint(5605), c.Uint, "c.Uint")
	assert.Equal(t, float32(14.56), c.Float32, "c.Float32")
	assert.Equal(t, float64(123456.789), c.Float64, "c.Float64")

	if c.Duration != time.Second*30 {
		t.Errorf("c.Duration should be 30 seconds, %s given", c.Duration.String())
	}
}
