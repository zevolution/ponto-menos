package env

import (
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

const (
	DEFAULT_CONFIG_FILE_PATH string = "./ponto-menos.toml"
)

func init() {
	viperCfg()
}

func viperCfg() error {
	dir, file := filepath.Split(DEFAULT_CONFIG_FILE_PATH)
	file = strings.TrimSuffix(file, filepath.Ext(file))
	viper.AddConfigPath(dir)
	viper.SetConfigName(file)
	viper.SetConfigType("toml")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	return nil
}

func GetOrDefault(key string, d interface{}) interface{} {
	if viper.InConfig(key) {
		return viper.Get(key)
	} else {
		return d
	}
}
