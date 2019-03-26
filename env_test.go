package config

import "testing"

func TestAllowed(t *testing.T) {
	tt := []struct {
		name     string
		search   string
		in       []string
		expected bool
	}{
		{
			name: "valid case in single member slice",
			search: "test",
			in: []string{"test"},
			expected: true,
		},
		{
			name: "valid case in multi member slice",
			search: "test",
			in: []string{"test", "foo", "bar", "baz"},
			expected: true,
		},
		{
			name: "empty search returns false",
			search: "test",
			in: []string{},
			expected: true,
		},
		{
			name: "invalid case in multi member slice",
			search: "test",
			in: []string{"foo", "bar", "baz"},
			expected: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func (t *testing.T) {
			actual := allowed(tc.in, tc.search)

			if actual != tc.expected {
				t.Error("invalid result for notIn")
			}
		})
	}
}
