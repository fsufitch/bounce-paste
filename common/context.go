package common

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var interruptLock Lock

func ContextWithInterruptCancel(log Logger, parent context.Context) (context.Context, error) {
	log.Debugf("binding signal interrupts to context")

	releaseLock, err := interruptLock.Acquire()
	if err != nil {
		err = fmt.Errorf("could not acquire lock to listen for interrupts: %w", err)
		log.Errorf("%v", err)
		return nil, err
	}

	ctx, cancelCtx := context.WithCancel(parent)

	go func() {
		defer cancelCtx()
		defer releaseLock()
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigs
		log.Infof("received signal: %s", sig.String())
	}()

	return ctx, nil
}
