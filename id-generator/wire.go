package idgenerator

import (
	"github.com/google/wire"

	"github.com/fsufitch/bounce-paste/mq"
)

var WireSet = wire.NewSet(
	wire.Struct(new(Server), "*"),
	mq.WireSet,
)
