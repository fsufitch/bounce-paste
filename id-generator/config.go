package idgenerator

import (
	"errors"
	"fmt"

	"github.com/fsufitch/bounce-paste/common"
)

type Config struct {
	Host     string
	Port     int
	IdPrefix string
}

var ErrConfigFailed = errors.New("id-generator config failed")
var defaultHost = "0.0.0.0"
var defaultPort = 9009
var defaultIdPrefix = "DEV_"

func ConfigFromEnviron(environ common.Environ) (*Config, error) {
	conf := Config{}
	var err error

	if conf.Host, err = environ.GetString("HOST"); errors.Is(err, common.ErrEnvironmentKeyMissing) {
		// XXX: use logging facility
		fmt.Printf("HOST not specified, using default: %s\n", defaultHost)
		conf.Host = defaultHost
	} else if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrConfigFailed, err)
	}

	if conf.Port, err = environ.GetInt("PORT"); errors.Is(err, common.ErrEnvironmentKeyMissing) {
		// XXX: use logging facility
		fmt.Printf("PORT not specified, using default: %v\n", defaultPort)
		conf.Port = defaultPort
	} else if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrConfigFailed, err)
	}

	if conf.IdPrefix, err = environ.GetString("ID_PREFIX"); errors.Is(err, common.ErrEnvironmentKeyMissing) {
		// XXX: use logging facility
		fmt.Printf("ID_PREFIX not specified, using default: %v\n", defaultIdPrefix)
		conf.IdPrefix = defaultIdPrefix
	} else if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrConfigFailed, err)
	}

	return &conf, nil
}
