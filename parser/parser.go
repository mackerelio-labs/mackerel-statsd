// Package parser implements statsd protocol.
package parser

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

const (
	TypeCounter = iota
	TypeTimer
)

// Metric represents a metric received via statsd protocol.
// This struct do not have a timestamp because Statsd will append a timestamp to them
// when these metrics is flushing to the backend.
type Metric struct {
	Type         int
	Name         string
	Value        float64
	SamplingRate float64
	Tags         []*Tag
}

func (m *Metric) String() string {
	return fmt.Sprintf("%s:%g|@%.1f|#%v", m.Name, m.Value, m.SamplingRate, m.Tags)
}

type Tag struct {
	Key   string
	Value string
}

func (t *Tag) String() string {
	return fmt.Sprintf("%s:%s", t.Key, t.Value)
}

func toType(c string) (int, error) {
	switch c {
	case "c":
		return TypeCounter, nil
	case "ms":
		return TypeTimer, nil
	default:
		return -1, fmt.Errorf("unknown type: %s", c)
	}
}

var msgSep = []byte{'\n'}

func Parse(msg []byte) ([]*Metric, error) {
	var a []*Metric
	for len(msg) > 0 {
		line, rest, _ := bytes.Cut(msg, msgSep)
		m, err := parse(string(line))
		if err != nil {
			return nil, err
		}
		a = append(a, m)
		msg = rest
	}
	return a, nil
}

func parse(line string) (*Metric, error) {
	// <bucket-name>:<value>|c(|@<sampling-rate>)(|#<key>:<value>,...)
	name, rest, found := strings.Cut(line, ":")
	if !found {
		return nil, fmt.Errorf("'%s' invalid format, expects ':'", line)
	}
	valStr, rest, found := strings.Cut(rest, "|")
	if !found {
		return nil, fmt.Errorf("'%s' invalid format, expects '|'", line)
	}
	val, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return nil, fmt.Errorf("'%s' invalid syntax, parsing '%s'", line, valStr)
	}
	typeStr, rest, _ := strings.Cut(rest, "|")
	typ, err := toType(typeStr)
	if err != nil {
		return nil, fmt.Errorf("'%s' invalid format, %w", line, err)
	}
	m := &Metric{
		Type:         typ,
		Name:         name,
		Value:        val,
		SamplingRate: 1.0,
		Tags:         nil,
	}

	for rest != "" {
		s, t, _ := strings.Cut(rest, "|")
		switch s[0] {
		default:
			return nil, fmt.Errorf("'%s' invalid format, unsupported optional argument '%s'", line, s)
		case '@':
			rate, err := strconv.ParseFloat(s[1:], 64)
			if err != nil {
				return nil, fmt.Errorf("'%s' invalid format, parsing '%s'", line, s[1:])
			}
			m.SamplingRate = rate
		case '#':
			colon, comma := detectDialect(s[1:])
			a := strings.Split(s[1:], comma)
			tags := make([]*Tag, len(a))
			for i, tag := range a {
				k, v, _ := strings.Cut(tag, colon)
				tags[i] = &Tag{Key: k, Value: v}
			}
			m.Tags = tags
		}
		rest = t
	}
	return m, nil
}

func detectDialect(s string) (colon, comma string) {
	colon = ":"
	comma = ","
	i := strings.IndexAny(s, ":=")
	if i >= 0 {
		colon = string(s[i])
		s = s[i+1:]
	}
	i = strings.IndexAny(s, ",;")
	if i >= 0 {
		comma = string(s[i])
	}
	return
}
