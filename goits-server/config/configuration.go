// Package config is just a thin wrapper around viper package for configuration initialization.
package config

import (
	"log"

	"github.com/spf13/viper"
)

// initializes config provider(viper) and reads the configuration.
//
// by default, looks for goits.cfg.yaml under etc sub directory of current directory.
func InitConfig() {
	viper.SetConfigName("goits.cfg")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./etc")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Could not find configuration", err)
	}
}

// deserialize the configuration under given key into given target data.
func ReadInto(key string, target interface{}) {
	viper.UnmarshalKey(key, target)
}
