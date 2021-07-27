// Package config is just a thin wrapper around viper package for configuration initialization.
package config

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var appRootDir string
var configRootDir string

var configNotFoundErr error = errors.New("unknown config key")

var configRegistry map[string]*viper.Viper

func init() {
	appRootDir = os.Getenv("GOITS_HOME")
	configRootDir = filepath.Join(appRootDir, "etc")
	configRegistry = make(map[string]*viper.Viper)
}

// InitConfig initializes config provider and reads the default configuration file into memory.
//
// By default, it looks for goits.cfg.yaml under etc sub directory of current directory.
func InitConfig() {
	viper.SetConfigName("goits.cfg")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configRootDir)

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Could not find configuration", err)
	}
}

// ReadInto deserializes the configuration under given key into given target struct.
func ReadInto(key string, target interface{}) {
	viper.UnmarshalKey(key, target)
}

// ReadFromConfigInto reads config section from a config file previously registered by using NewConfig.
func ReadFromConfigInto(configKey string, key string, target interface{}) error {
	v, ok := configRegistry[configKey]
	if ok {
		return v.UnmarshalKey(key, target)
	}
	return configNotFoundErr
}

// ReadConfig reads and registers a specific config file for later use.
func ReadConfig(configKey string, rootDir string, baseName string, configStruct interface{}) error {
	v := viper.New()

	v.SetConfigName(baseName)
	v.SetConfigType("yaml")
	v.AddConfigPath(ConfigPath(rootDir))

	err := v.ReadInConfig()
	if err == nil {
		configRegistry[configKey] = v
		if configStruct != nil {
			v.Unmarshal(configStruct)
		}
	}
	return err
}

// AppPath returns the full path for the relative config dir.
func ConfigPath(dir string) string {
	return filepath.Join(configRootDir, dir)
}
