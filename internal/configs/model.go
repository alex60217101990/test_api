package configs

import (
	"encoding/json"

	"github.com/pkg/errors"
)

var Conf *Configs

type Configs struct {
	Ver         *string    `yaml:"ver"`
	ServiceName string     `yaml:"service-name" json:"service_name"`
	IsDebug     bool       `yaml:"-" json:"-"`
	LoggerType  LoggerType `yaml:"logger-type" json:"logger_type"`
	APIKey      string     `yaml:"api-key" json:"api_key"`
	Server      *Server    `yaml:"http-server" json:"http_server"`
	DB          *DB        `yaml:"db"`
	Keys        *Keys      `yaml:"keys"`
	Timeouts    *Timeouts  `yaml:"timeouts"`
}

type Server struct {
	Host string `yaml:"server-host" json:"server_host"`
	Port uint16 `yaml:"server-port" json:"server_port"`
}

type DB struct {
	RepoType RepoType `yaml:"db-type" json:"db_type"`
	Host     string   `yaml:"db-host" json:"db_host"`
	UserName string   `yaml:"db-user" json:"db_user"`
	Password string   `yaml:"db-password" json:"db_password"`
	Port     uint16   `yaml:"db-port" json:"db_port"`
	DbName   string   `yaml:"db-name" json:"db_name"`
}

func (db *DB) MarshalJSON() ([]byte, error) {
	type alias struct {
		RepoType string `json:"db_type"`
		Host     string `json:"db_host"`
		UserName string `json:"db_user"`
		Password string `json:"db_password"`
		Port     uint16 `json:"db_port"`
		DbName   string `json:"db_name"`
	}
	if db == nil {
		db = &DB{}
	}
	return json.Marshal(alias{
		RepoType: db.RepoType.String(),
		Host:     db.Host,
		UserName: db.UserName,
		Password: db.Password,
		Port:     db.Port,
		DbName:   db.DbName,
	})
}

func (db *DB) UnmarshalJSON(data []byte) (err error) {
	type alias struct {
		RepoType string `json:"db_type"`
		Host     string `json:"db_host"`
		UserName string `json:"db_user"`
		Password string `json:"db_password"`
		Port     uint16 `json:"db_port"`
		DbName   string `json:"db_name"`
	}
	var tmp alias
	if err = json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	if db == nil {
		db = &DB{}
	}

	err = db.RepoType.Set(tmp.RepoType)
	if err != nil {
		return errors.WithMessagef(err, "failed to parse '%s'", tmp.RepoType)
	}

	db.Host = tmp.Host
	db.UserName = tmp.UserName
	db.Password = tmp.Password
	db.Port = tmp.Port
	db.DbName = tmp.DbName

	return nil
}

func (db *DB) MarshalYAML() (interface{}, error) {
	type alias struct {
		RepoType string `yaml:"db-type"`
		Host     string `yaml:"db-host"`
		UserName string `yaml:"db-user"`
		Password string `yaml:"db-password"`
		Port     uint16 `yaml:"db-port"`
		DbName   string `yaml:"db-name"`
	}
	if db == nil {
		db = &DB{}
	}
	return alias{
		RepoType: db.RepoType.String(),
		Host:     db.Host,
		UserName: db.UserName,
		Password: db.Password,
		Port:     db.Port,
		DbName:   db.DbName,
	}, nil
}

func (db *DB) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type alias struct {
		RepoType string `yaml:"db-type"`
		Host     string `yaml:"db-host"`
		UserName string `yaml:"db-user"`
		Password string `yaml:"db-password"`
		Port     uint16 `yaml:"db-port"`
		DbName   string `yaml:"db-name"`
	}
	var tmp alias
	if err := unmarshal(&tmp); err != nil {
		return err
	}

	if db == nil {
		db = &DB{}
	}

	err := db.RepoType.Set(tmp.RepoType)
	if err != nil {
		return errors.WithMessagef(err, "failed to parse '%s'", tmp.RepoType)
	}

	db.Host = tmp.Host
	db.UserName = tmp.UserName
	db.Password = tmp.Password
	db.Port = tmp.Port
	db.DbName = tmp.DbName

	return nil
}

type Keys struct {
	PrvKeyRepo string `yaml:"prv-key-repo" json:"prv_key_repo"`
	PubKeyRepo string `yaml:"pub-key-repo" json:"pub_key_repo"`
	PrvKeyAuth string `yaml:"prv-key-auth" json:"prv_key_auth"`
	PubKeyAuth string `yaml:"pub-key-auth" json:"pub_key_auth"`
}

type Timeouts struct {
	DefaultTimeout uint8  `yaml:"default-timeout" json:"default_timeout"`
	ExpHours       uint8  `yaml:"exp-hours" json:"exp_hours"`
	ExpSeconds     uint16 `yaml:"exp-seconds" json:"exp_seconds"`
}
