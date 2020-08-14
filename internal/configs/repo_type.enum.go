package configs

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

type RepoType uint8

const (
	RepoPostgres RepoType = iota
)

var (
	_RepoTypeNameToValue = map[string]RepoType{
		"postgres": RepoPostgres,
		"Postgres": RepoPostgres,
	}

	_RepoTypeValueToName = map[RepoType]string{
		RepoPostgres: "postgres",
	}
)

func (rt RepoType) MarshalYAML() (interface{}, error) {
	s, ok := _RepoTypeValueToName[rt]
	if !ok {
		return nil, fmt.Errorf("invalid RepoType: %d", rt)
	}
	return s, nil
}

func (rt *RepoType) UnmarshalYAML(value *yaml.Node) error {
	v, ok := _RepoTypeNameToValue[value.Value]
	if !ok {
		return fmt.Errorf("invalid RepoType %q", value.Value)
	}
	*rt = v
	return nil
}

func (rt RepoType) MarshalJSON() ([]byte, error) {
	if s, ok := interface{}(rt).(fmt.Stringer); ok {
		return json.Marshal(s.String())
	}
	s, ok := _RepoTypeValueToName[rt]
	if !ok {
		return nil, fmt.Errorf("invalid RepoType: %d", rt)
	}
	return json.Marshal(s)
}

func (rt *RepoType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("RepoType should be a string, got %s", data)
	}
	v, ok := _RepoTypeNameToValue[s]
	if !ok {
		return fmt.Errorf("invalid RepoType %q", s)
	}
	*rt = v
	return nil
}

func (rt RepoType) Val() uint8 {
	return uint8(rt)
}

// it's for using with flag package
func (rt *RepoType) Set(val string) error {
	if at, ok := _RepoTypeNameToValue[val]; ok {
		*rt = at
		return nil
	}
	return fmt.Errorf("invalid repository type: %v", val)
}

func (rt RepoType) String() string {
	return _RepoTypeValueToName[rt]
}
