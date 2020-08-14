package server

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

type RouteType uint8

const (
	Get RouteType = iota
	Post
	Delete
	Put
)

var (
	_RouteTypeNameToValue = map[string]RouteType{
		"Get":    Get,
		"get":    Get,
		"GET":    Get,
		"Post":   Post,
		"post":   Post,
		"POST":   Post,
		"Delete": Delete,
		"delete": Delete,
		"DELETE": Delete,
		"Put":    Put,
		"put":    Put,
		"PUT":    Put,
	}

	_RouteTypeValueToName = map[RouteType]string{
		Get:  "get",
		Post: "post",
	}
)

func (rt RouteType) MarshalYAML() (interface{}, error) {
	s, ok := _RouteTypeValueToName[rt]
	if !ok {
		return nil, fmt.Errorf("invalid RouteType: %d", rt)
	}
	return s, nil
}

func (rt *RouteType) UnmarshalYAML(value *yaml.Node) error {
	v, ok := _RouteTypeNameToValue[value.Value]
	if !ok {
		return fmt.Errorf("invalid RouteType %q", value.Value)
	}
	*rt = v
	return nil
}

func (rt RouteType) MarshalJSON() ([]byte, error) {
	if s, ok := interface{}(rt).(fmt.Stringer); ok {
		return json.Marshal(s.String())
	}
	s, ok := _RouteTypeValueToName[rt]
	if !ok {
		return nil, fmt.Errorf("invalid RouteType: %d", rt)
	}
	return json.Marshal(s)
}

func (rt *RouteType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("RouteType should be a string, got %s", data)
	}
	v, ok := _RouteTypeNameToValue[s]
	if !ok {
		return fmt.Errorf("invalid RouteType %q", s)
	}
	*rt = v
	return nil
}

func (rt RouteType) Val() uint8 {
	return uint8(rt)
}

func (rt RouteType) String() string {
	return _RouteTypeValueToName[rt]
}
