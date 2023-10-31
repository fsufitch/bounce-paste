package mq

import (
	"fmt"
	"time"

	"github.com/fsufitch/bounce-paste/common"
)

type Config struct {
	Host            string
	User            string
	Password        string
	ConnectAttempts int
	ConnectDelay    time.Duration
}

const defaultConnectAttempts = 3
const defaultConnectDelay = 5 * time.Second

func ProvideConfigFromEnviron(log common.Logger, environ common.Environ) (conf Config, err error) {
	log.Debugf("loading MQ config")

	varToStorageMap := map[string]*string{
		"MQ_HOST":     &conf.Host,
		"MQ_USER":     &conf.User,
		"MQ_PASSWORD": &conf.Password,
	}

	for varName, dest := range varToStorageMap {
		*dest, err = environ.GetString(varName)
		if err == nil && *dest == "" {
			err = fmt.Errorf("value may not be blank")
		}
		if err != nil {
			err = fmt.Errorf("error getting env var: %w", err)
			log.Errorf("%s", err)
			return conf, err
		}
		log.Debugf("successfully read var: %s", varName)
	}

	if conf.ConnectAttempts, err = environ.GetInt("MQ_CONN_ATTEMPTS"); err != nil {
		log.Warningf("MQ_CONN_ATTEMPTS unset, using default: %v", defaultConnectAttempts)
		conf.ConnectAttempts = defaultConnectAttempts
	}
	log.Debugf("times to attempt MQ connection: %v", conf.ConnectAttempts)

	if conf.ConnectAttempts, err = environ.GetInt("MQ_CONN_ATTEMPTS"); err != nil {
		log.Warningf("MQ_CONN_ATTEMPTS unset, using default: %v", defaultConnectAttempts)
		conf.ConnectAttempts = defaultConnectAttempts
	}

	var seconds int
	if seconds, err = environ.GetInt("MQ_CONN_DELAY"); err != nil {
		seconds = int(defaultConnectDelay.Seconds())
		log.Warningf("MQ_CONN_DELAY unset, using default: %v", seconds)
	}
	conf.ConnectDelay = time.Duration(seconds) * time.Second
	log.Debugf("delay between MQ connection attempts: %s", conf.ConnectDelay.String())

	return conf, nil
}
