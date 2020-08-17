package configs

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

func ReadConfigFile(file string) (err error) {
	var yamlFile []byte
	_, err = os.Stat(file)
	if os.IsNotExist(err) && err != nil {
		file, err = filepath.EvalSymlinks(file)
		if err != nil {
			return err
		}
		yamlFile, err = ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		err = yaml.Unmarshal(yamlFile, &Conf)
		return err
	}

	yamlFile, err = ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, &Conf)
	return err
}
