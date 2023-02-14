package parser

import (
	"reflect"
	"strings"
	"testing"
)

func TestPrase(t *testing.T) {
	tests := map[string][]*Metric{
		"custom.metric.name:10|c": []*Metric{
			{
				Name:         "custom.metric.name",
				Value:        10.0,
				Type:         TypeCounter,
				SamplingRate: 1.0,
			},
		},

		"custom.metric.name:10|c|@0.1": []*Metric{
			{
				Name:         "custom.metric.name",
				Value:        10.0,
				Type:         TypeCounter,
				SamplingRate: 0.1,
			},
		},

		"custom.metric.name:10|ms|@0.1": []*Metric{
			{
				Name:         "custom.metric.name",
				Value:        10.0,
				Type:         TypeTimer,
				SamplingRate: 0.1,
			},
		},

		"custom.metric.name:10|c|#tag1:v1": []*Metric{
			{
				Name:         "custom.metric.name",
				Value:        10.0,
				Type:         TypeCounter,
				SamplingRate: 1.0,
				Tags: []*Tag{
					{Key: "tag1", Value: "v1"},
				},
			},
		},

		"custom.metric.name:10|c|#tag-only1;tag-only2": []*Metric{
			{
				Name:         "custom.metric.name",
				Value:        10.0,
				Type:         TypeCounter,
				SamplingRate: 1.0,
				Tags: []*Tag{
					{Key: "tag-only1", Value: ""},
					{Key: "tag-only2", Value: ""},
				},
			},
		},

		"custom.metric.name:10|c|@0.1|#tag1:v1,tag2:v2": []*Metric{
			{
				Name:         "custom.metric.name",
				Value:        10.0,
				Type:         TypeCounter,
				SamplingRate: 0.1,
				Tags: []*Tag{
					{Key: "tag1", Value: "v1"},
					{Key: "tag2", Value: "v2"},
				},
			},
		},

		"custom.metric.name:10|c|@0.1|#tag1=v1,tag2=v2": []*Metric{
			{
				Name:         "custom.metric.name",
				Value:        10.0,
				Type:         TypeCounter,
				SamplingRate: 0.1,
				Tags: []*Tag{
					{Key: "tag1", Value: "v1"},
					{Key: "tag2", Value: "v2"},
				},
			},
		},

		"custom.metric.name:10|c|@0.1|#tag1=v1;tag2=v2": []*Metric{
			{
				Name:         "custom.metric.name",
				Value:        10.0,
				Type:         TypeCounter,
				SamplingRate: 0.1,
				Tags: []*Tag{
					{Key: "tag1", Value: "v1"},
					{Key: "tag2", Value: "v2"},
				},
			},
		},

		"custom.cpu.user:10|c\ncustom.cpu.sys:20|c": []*Metric{
			{
				Name:         "custom.cpu.user",
				Value:        10.0,
				Type:         TypeCounter,
				SamplingRate: 1.0,
			},
			{
				Name:         "custom.cpu.sys",
				Value:        20.0,
				Type:         TypeCounter,
				SamplingRate: 1.0,
			},
		},
	}
	for msg, metric := range tests {
		t.Run(msg, func(t *testing.T) {
			a, err := Parse([]byte(msg))
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(a, metric) {
				t.Errorf("Parse(%q) = %v; want %v", msg, a, metric)
			}
			t.Logf("%#v", a)
		})
	}
}

func TestParseErr(t *testing.T) {
	invalidMessages := map[string]string{
		"custom.metric.name":         "invalid format, expects ':'",
		"custom.metric.name:10":      "invalid format, expects '|'",
		"custom.metric.name:10|x":    "invalid format, unknown type: x",
		"custom.metric.name:a|c":     "invalid syntax, parsing 'a'",
		"custom.metric.name:10|c|@r": "invalid format, parsing 'r'",
		"custom.metric.name:10|c|%r": "invalid format, unsupported optional argument '%r'",
	}
	for msg, want := range invalidMessages {
		_, err := Parse([]byte(msg))
		if err == nil {
			t.Errorf("Parse(%q) returns nil; expects /%s/", msg, want)
		}
		if s := err.Error(); !strings.Contains(s, want) {
			t.Errorf("Parse(%q) returns %q; expects /%s/", msg, s, want)
		}
	}
}
