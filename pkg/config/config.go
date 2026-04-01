// Package config provides simple functions for wrapping Viper package.
package config

import "github.com/spf13/viper"

func init() {
	viper.AddConfigPath(".")
}

// AddPaths adds paths for Viper to look in where finding file.
func AddPaths(paths ...string) {
	for _, path := range paths {
		viper.AddConfigPath(path)
	}
}

// LoadEnv receives filename (without extension) and loads content of .env file into Viper.
func LoadEnv(filename string) error {
	viper.SetConfigName(filename)
	viper.SetConfigType("env")
	return viper.MergeInConfig()
}
