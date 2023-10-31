package mq

import "github.com/google/wire"

var WireSet = wire.NewSet(
	ProvideConfigFromEnviron,
	ProvideMQ,
)
