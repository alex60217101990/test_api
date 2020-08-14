package configs

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

type LoggerType uint8

const (
	Zap LoggerType = iota
	Base
)

var (
	_LoggerTypeNameToValue = map[string]LoggerType{
		"zap":  Zap,
		"ZAP":  Zap,
		"Zap":  Zap,
		"base": Base,
		"Base": Base,
		"BASE": Base,
	}

	_LoggerTypeValueToName = map[LoggerType]string{
		Zap:  "zap",
		Base: "base",
	}
)

func (lt LoggerType) MarshalYAML() (interface{}, error) {
	s, ok := _LoggerTypeValueToName[lt]
	if !ok {
		return nil, fmt.Errorf("invalid LoggerType: %d", lt)
	}
	return s, nil
}

func (lt *LoggerType) UnmarshalYAML(value *yaml.Node) error {
	v, ok := _LoggerTypeNameToValue[value.Value]
	if !ok {
		return fmt.Errorf("invalid LoggerType %q", value.Value)
	}
	*lt = v
	return nil
}

func (lt LoggerType) MarshalJSON() ([]byte, error) {
	if s, ok := interface{}(lt).(fmt.Stringer); ok {
		return json.Marshal(s.String())
	}
	s, ok := _LoggerTypeValueToName[lt]
	if !ok {
		return nil, fmt.Errorf("invalid LoggerType: %d", lt)
	}
	return json.Marshal(s)
}

func (lt *LoggerType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("LoggerType should be a string, got %s", data)
	}
	v, ok := _LoggerTypeNameToValue[s]
	if !ok {
		return fmt.Errorf("invalid LoggerType %q", s)
	}
	*lt = v
	return nil
}

func (lt LoggerType) Val() uint8 {
	return uint8(lt)
}

// it's for using with flag package
func (lt *LoggerType) Set(val string) error {
	if at, ok := _LoggerTypeNameToValue[val]; ok {
		*lt = at
		return nil
	}
	return fmt.Errorf("invalid logger type: %v", val)
}

func (lt LoggerType) String() string {
	return _LoggerTypeValueToName[lt]
}
