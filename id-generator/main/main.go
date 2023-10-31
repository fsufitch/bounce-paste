package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/google/wire"

	"github.com/fsufitch/bounce-paste/common"
	idgenerator "github.com/fsufitch/bounce-paste/id-generator"
)

func defaultInitializer() (*idgenerator.Server, func(), error) {
	return nil, func() {}, errors.New("binary was built without wire dependency injection")
}

var initializerFunc func() (*idgenerator.Server, func(), error) = defaultInitializer

func main() {
	server, cleanup, err := initializerFunc()

	if err != nil {
		fmt.Fprintf(os.Stderr, "initialization error: %s", err)
		cleanup()
		os.Exit(1)
	}

	err = server.Run()

	if err != nil {
		fmt.Fprintf(os.Stderr, "server exited with error: %s", err)
		cleanup()
		os.Exit(1)
	}

	cleanup()
	os.Exit(0)
}

func ProvideServerContext(log common.Logger) (idgenerator.ServerContext, error) {
	return common.ContextWithInterruptCancel(log, context.Background())
}

var IdGeneratorServerWireSet = wire.NewSet(
	ProvideServerContext,
	idgenerator.WireSet,
	common.WireSet,
	wire.InterfaceValue(new(common.LogWriter), os.Stdout),
	wire.Value(common.LogPrefix("id-generator")),
)
