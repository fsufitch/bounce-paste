package common

import (
	"fmt"
	"io"
	"os"

	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	EnvironFromProcessEnviron,
)

type LoggingWireSetFactory struct {
	Writer          io.Writer
	Prefix          *string
	DisableWarnings bool
}

func (lwsc LoggingWireSetFactory) BuildWireSet() wire.ProviderSet {
	warn := func(format string, args ...any) {
		if lwsc.DisableWarnings {
			return
		}
		fmt.Fprintf(os.Stderr, "(logging init) "+format+"\n", args...)
	}

	writer := lwsc.Writer
	if lwsc.Writer == nil {
		warn("using default logging destination (stdout)")
		writer = os.Stdout
	}

	prefix := ""
	if lwsc.Prefix == nil {
		prefix = fmt.Sprintf("%s(%d)", os.Args[0], os.Getpid())
		warn("using default logging prefix (%s)", prefix)
	}

	return wire.NewSet(
		wire.Value(LogPrefix(prefix)),
		wire.InterfaceValue(new(LogWriter), writer),
	)
}
