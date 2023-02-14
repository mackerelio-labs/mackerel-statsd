package metricname

import (
	"testing"
)

func TestJoin(t *testing.T) {
	tests := map[string][]string{
		"custom.name": []string{"custom", "name"},
		"":            []string{},
		"name":        []string{"name"},
	}
	for want, args := range tests {
		if s := Join(args...); s != want {
			t.Errorf("Join(%v) = %q; want %q", args, s, want)
		}
	}
}

func TestGroup(t *testing.T) {
	tests := map[string]string{
		"name":              "",
		"custom.group.name": "custom.group",
		"custom.name":       "custom",
		"":                  "",
	}
	for name, group := range tests {
		if s := Group(name); s != group {
			t.Errorf("Group(%q) = %q; want %q", name, s, group)
		}
	}
}
