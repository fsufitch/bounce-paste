package common

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var ErrBadEnvironment = errors.New("bad environment")
var ErrEnvironmentKeyMissing = errors.New("missing env var")
var ErrEnvironmentBadValue = errors.New("bad env value")

type Environ map[string]string

func ProvideEnvironFromProcessEnvironment() (Environ, error) {
	env := Environ{}
	for _, envStr := range os.Environ() {
		parts := strings.SplitN(envStr, "=", 2)
		if len(parts) < 2 {
			return nil, fmt.Errorf("%w: %s", ErrBadEnvironment, envStr)
		}
		key := parts[0]
		value := parts[1]
		env[key] = value
	}
	return env, nil
}

func (environ Environ) GetString(key string) (string, error) {
	if value, ok := environ[key]; ok {
		return value, nil
	}
	return "", fmt.Errorf("%w: %s", ErrEnvironmentKeyMissing, key)
}

func (environ Environ) GetInt(key string) (int, error) {
	strValue, ok := environ[key]
	if !ok {
		return 0, fmt.Errorf("%w: %s", ErrEnvironmentKeyMissing, key)
	}
	value, err := strconv.ParseInt(strValue, 10, 0)
	if err == nil && (value < 0 || value > 65535) {
		err = fmt.Errorf("invalid port number: %d", value)
	}
	if err != nil {
		return 0, fmt.Errorf("%w: key=%v value=%v", ErrEnvironmentBadValue, key, strValue)
	}
	return int(value), nil
}

func (environ Environ) GetBool(key string) bool {
	value, ok := environ[key]
	return ok && value != ""
}

type DebugMode bool

func ProvideDebugMode(environ Environ) DebugMode {
	return DebugMode(environ.GetBool("DEBUG"))
}
