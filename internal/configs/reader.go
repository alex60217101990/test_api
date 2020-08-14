package configs

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func ReadConfigFile(file string) error {
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, &Conf)
	if err != nil {
		return err
	}

	return nil
}
