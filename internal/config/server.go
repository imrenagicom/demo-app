package config

import (
	"bytes"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func NewServer(path, envPrefix string) (Server, error) {
	fang := viper.New()
	fang.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	fang.AutomaticEnv()
	fang.SetEnvPrefix(envPrefix)
	fang.SetConfigType("yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return Server{}, err
	}
	if err := fang.ReadConfig(bytes.NewBuffer(data)); err != nil {
		return Server{}, err
	}
	// Load configuration
	s := Server{}
	if err = fang.Unmarshal(&s); err != nil {
		return Server{}, err
	}
	return s, nil
}