package config

import "github.com/yaml"

type ConfigEntry struct{
	ConfigMap map[string]interface{}
}

func InitConfig() (map[string]interface{},error){

	configEntry := &ConfigEntry{}
	_,err := yaml.Marshal(configEntry)

	return configEntry.ConfigMap,err
}
