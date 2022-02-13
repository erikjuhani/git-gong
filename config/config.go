package config

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/erikjuhani/git-gong/fs"
	"github.com/spf13/viper"
)

type ConfigKey = string

const (
	AllowedBranchPatternsKey ConfigKey = "rules.allowed_branch_patterns"
)

const (
	configPath = ".gong/config"
	configType = "toml"
)

var initFns = []func(){
	genAllowedBranchPatterns,
}

type Patterns []*regexp.Regexp

func (p Patterns) Match(s string) bool {
	if len(p) == 0 {
		return true
	}

	for _, pattern := range p {
		if pattern.MatchString(s) {
			return true
		}
	}

	return false
}

func (p Patterns) AddPattern(pattern string) {
	p = append(p, regexFromGlobLike("", pattern))
}

var (
	AllowedBranchPatterns = &Patterns{}
)

func Get(key ConfigKey) interface{} {
	return viper.Get(key)
}

func genAllowedBranchPatterns() {
	patterns := viper.GetStringSlice(AllowedBranchPatternsKey)

	for _, pattern := range patterns {
		AllowedBranchPatterns.AddPattern(pattern)
	}
}

var regexReplaceCharMap = []string{
	"/", "\\/",
	"(", "\\(",
	")", "\\)",
}

func regexFromGlobLike(regex string, pattern string) *regexp.Regexp {
	parts := strings.SplitN(pattern, "*", 2)

	replacer := strings.NewReplacer(regexReplaceCharMap...)
	regex = replacer.Replace(parts[0])

	if len(parts) > 1 && parts[1] != "" {
		return regexFromGlobLike(regex, parts[1])
	}

	return regexp.MustCompile(regex)
}

func Load() error {
	if err := loadConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return viper.SafeWriteConfigAs(configPath)
		}
		return fmt.Errorf("fatal error loading config file %w", err)
	}
	return nil
}

func setDefaults() {
	viper.SetDefault(AllowedBranchPatternsKey, make([]string, 0))
}

func loadConfig() error {
	// Check that .gong directory exists
	if err := fs.EnsureDir(filepath.Dir(configPath)); err != nil {
		return err
	}

	// Write a new config file if it does not exist
	err := viper.SafeWriteConfigAs(configPath)

	if _, ok := err.(viper.ConfigFileAlreadyExistsError); ok {
		// Set defaults to handle empty cases
		setDefaults()

		// Read in the config
		err := viper.ReadInConfig()

		for _, fn := range initFns {
			fn()
		}

		return err
	}

	return err
}

func init() {
	viper.SetConfigFile(configPath)
	viper.SetConfigType(configType)
}
