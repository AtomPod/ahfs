package setting

import (
	"github.com/spf13/viper"
)

type Storage struct {
	Type string
	Path string
	Name string
}

func (s *Storage) Unmarshal(v interface{}) error {
	path := s.Path
	if len(path) == 0 {
		path = "storage." + s.Name + ".config"
	}

	configStorage := viper.Sub(path)
	if configStorage == nil {
		return nil
	}

	return configStorage.Unmarshal(v)
}

func getStorage(name string) Storage {

	var storage Storage
	storage.Name = name

	keyName := "storage." + name
	typeStorageCfg := viper.Sub(keyName)
	if typeStorageCfg == nil {
		return storage
	}

	storage.Type = typeStorageCfg.GetString("type")
	storage.Path = typeStorageCfg.GetString("path")

	return storage
}
