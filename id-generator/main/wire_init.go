//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"

	idgenerator "github.com/fsufitch/bounce-paste/id-generator"
)

// The build tag makes sure the stub is not built in the final build.

func InitializeServer() (*idgenerator.Server, func(), error) {
	panic(wire.Build(IdGeneratorServerWireSet))
}

func init() {
	initializerFunc = InitializeServer
}
