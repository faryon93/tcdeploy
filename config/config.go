package config

import (
	"github.com/BurntSushi/toml"
)

// ----------------------------------------------------------------------------------
//  constants
// ----------------------------------------------------------------------------------

const (
	PROVIDER_TC = "tc"
)


// ----------------------------------------------------------------------------------
//  types
// ----------------------------------------------------------------------------------

type Config struct {
	Path string
	Provider string `toml:"provider"`

	TcUrl string `toml:"tc_url"`
	TcUser string `toml:"tc_user"`
	TcPassword string `toml:"tc_password"`
	TcBuildConfId string `toml:"tc_build_conf"`
}


// ----------------------------------------------------------------------------------
//  functions
// ----------------------------------------------------------------------------------

func Load(path string) (*Config, error) {
	// decode the config file to struct
	var config Config
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, err
	}

	config.Path = path

	return &config, nil
}

