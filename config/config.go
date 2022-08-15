package config

import (
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"sort"
	"strings"
)

const (
	DevelopmentMode = "development-mode"
	Port            = "port"
	LogLevel        = "log-level"
	ProtectedRules  = "protected-rules"
)

type Config struct {
	DevelopmentMode bool     `json:"development-mode"`
	Port            string   `json:"port"`
	LogLevel        string   `json:"log-level"`
	ProtectedRules  []string `json:"protected-rules"`
}

func init() {
	// Automatically read configuration options from environment variables.
	// e.g. --development-mode be configurable using ARMOR_DEVELOPMENT_MODE.
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("ARMOR")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))

	// Read configuration file from working directory and/or /etc.
	// File formats supported include JSON, TOML, YAML, HCL, envfile and Java properties config files
	viper.SetConfigType("yaml")
	viper.SetConfigName("." + "armor")

	flag.String(DevelopmentMode, "false",
		"Development mode. If true, the server will not enforce authentication and will allow all requests.")
	flag.String(Port, ":8080", "Port to listen on.")
	flag.String(LogLevel, "debug", "The log level to use.")
	flag.StringSlice(ProtectedRules, []string{"1000", "2147483647"},
		"The default armor rules protected from deletion or update and managed by terraform.")
}

func NewConfig() (*Config, error) {
	var err error
	var cfg Config

	err = viper.ReadInConfig()
	if err != nil {
		if err.(viper.ConfigFileNotFoundError) != err {
			return nil, err
		}
	}

	flag.Parse()

	err = viper.BindPFlags(flag.CommandLine)
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&cfg, decoderHook)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func decoderHook(dc *mapstructure.DecoderConfig) {
	dc.TagName = "json"
	dc.ErrorUnused = true
}

func (c *Config) Validate(required []string) error {
	present := func(key string) bool {
		for _, requiredKey := range required {
			if requiredKey == key {
				return len(viper.GetString(requiredKey)) > 0
			}
		}
		return true
	}
	var keys sort.StringSlice = viper.AllKeys()
	errs := make([]string, 0)

	keys.Sort()
	for _, key := range keys {
		if !present(key) {
			errs = append(errs, key)
		}
	}

	for _, key := range errs {
		log.Printf("required key '%s' not configured", key)
	}
	if len(errs) > 0 {
		return errors.New(fmt.Sprintf("missing configuration values %s", errs))
	}
	return nil
}

func (c *Config) IsProtectedRule(priority string) bool {
	return contains(c.ProtectedRules, priority)
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func SetupConfig() (*Config, error) {
	cfg, err := NewConfig()
	if err != nil {
		return nil, err
	}

	if err = cfg.Validate([]string{
		DevelopmentMode,
		Port,
		LogLevel,
		ProtectedRules,
	}); err != nil {
		return nil, err
	}
	return cfg, nil
}
