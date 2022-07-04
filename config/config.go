package config

import (
	"errors"
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
)

type Config struct {
	DevelopmentMode bool   `json:"development-mode"`
	Port            string `json:"port"`
}

func init() {
	// Automatically read configuration options from environment variables.
	// e.g. --development-mode be configurable using ARMOR_DEVELOPMENT_MODE.
	viper.SetEnvPrefix("ARMOR")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))

	// Read configuration file from working directory and/or /etc.
	// File formats supported include JSON, TOML, YAML, HCL, envfile and Java properties config files
	viper.SetConfigName("armor")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc")

	flag.String(DevelopmentMode, "false", "Toggle for development mode.")
	flag.String(Port, ":8080", "The address the metric endpoint binds to.")
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

func (c Config) Validate(required []string) error {
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
		return errors.New("missing configuration values")
	}
	return nil
}
