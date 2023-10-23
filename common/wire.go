package common

import (
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	ProvideEnvironFromProcessEnvironment,
	ProvideDebugMode,
	ProvideLogger,
	ProvideLogLevel,

	// Note: provide LogWriter and LogPrefix elsewhere
)
