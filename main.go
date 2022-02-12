package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/erikjuhani/git-gong/cmd"
	"github.com/erikjuhani/git-gong/fs"
	"github.com/spf13/viper"
)

const (
	configPath = ".gong/config"
	configType = "toml"
)

func LoadConfig() error {
	// Check that .gong directory exists
	if err := fs.EnsureDir(filepath.Dir(configPath)); err != nil {
		return err
	}

	// Setup configuration
	viper.SetConfigFile(configPath)
	viper.SetConfigType(configType)

	// Write a new config file if it does not exist
	if err := viper.SafeWriteConfigAs(configPath); err != nil {
		if _, ok := err.(viper.ConfigFileAlreadyExistsError); ok {
			return nil
		}
		return err
	}

	return viper.ReadInConfig()
}

func main() {
	if err := LoadConfig(); err != nil {
		log.Printf("%v", err)
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.SafeWriteConfigAs(configPath)
		} else {
			panic(fmt.Errorf("fatal error loading config file %w", err))
		}
	}

	cmd.Execute()
}
