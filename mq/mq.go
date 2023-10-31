package mq

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/go-amqp"

	"github.com/fsufitch/bounce-paste/common"
)

type MQ struct {
	config               Config
	log                  common.Logger
	reconnectRequestedCh chan any
	connectionStateCh    chan ConnectionState
	asyncManageConnLock  common.Lock
	connectionDialLock   common.Lock
}

func New(config Config, log common.Logger) MQ {
	return MQ{
		config:               config,
		log:                  log,
		reconnectRequestedCh: make(chan any),
		connectionStateCh:    make(chan ConnectionState),
	}
}

func ProvideMQ(config Config, log common.Logger) MQ {
	return New(config, log)
}

type ConnectionState struct {
	Conn  *amqp.Conn
	Error error
}

func (mq MQ) dial(ctx context.Context) (*amqp.Conn, error) {
	releaseLock, err := mq.connectionDialLock.Acquire()
	if err != nil {
		return nil, fmt.Errorf("failed to acquire dial lock: %w", err)
	}
	defer releaseLock()

	var conn *amqp.Conn
	for attempt := 1; attempt <= mq.config.ConnectAttempts; attempt++ {
		if attempt > 1 {
			mq.log.Debugf("sleep %s before next MQ connect attempt", mq.config.ConnectDelay.String())
			select {
			case <-time.After(mq.config.ConnectDelay):
			case <-ctx.Done():
				mq.log.Debugf("sleep interrupted")
				return nil, fmt.Errorf("context canceled")
			}
		}
		mq.log.Infof("MQ connect attempt; host=%v user=%v attempt=#%d/%d", mq.config.Host, mq.config.User, attempt, mq.config.ConnectAttempts)
		conn, err = amqp.Dial(ctx,
			fmt.Sprintf("amqp://%s", mq.config.Host), &amqp.ConnOptions{
				SASLType: amqp.SASLTypePlain(mq.config.User, mq.config.Password),
			})
		if err != nil {
			mq.log.Infof("MQ connect failed; host=%v user=%v attempt=#%d/%d err=%v", mq.config.Host, mq.config.User, attempt, mq.config.ConnectAttempts, err)
			continue
		}

		mq.log.Infof("MQ connect success; host=%v user=%v attempt=#%d/%d", mq.config.Host, mq.config.User, attempt, mq.config.ConnectAttempts)
	}

	if err != nil {
		err = fmt.Errorf("exhausted %d MQ connection attempts; %w", mq.config.ConnectAttempts, err)
		mq.log.Errorf("%v", err)
	}
	return conn, err
}

func (mq MQ) ManageConnectionAsync(ctx context.Context) (err error) {
	releaseLock, err := mq.asyncManageConnLock.Acquire()
	if err != nil {
		return fmt.Errorf("could not acquire connection management lock: %w", err)
	}
	defer releaseLock()

	dialSignal := make(chan bool, 1)

	// Consumer for reconnection requests, which only triggers a reconnection when there isn't a connection attempt already happening
	filterWorker := func() {
		mq.log.Debugf("MQ filter worker started")
		for {
			select {
			case <-mq.reconnectRequestedCh:
			case <-ctx.Done():
				mq.log.Infof("filter worker canceled while receiving reconnect requests")
				return
			}
			mq.log.Debugf("MQ filter worker received signal")
			if mq.connectionDialLock.IsLocked() {
				mq.log.Infof("received MQ reconnect signal, but reconnect already in progress")
				return
			}
			select {
			case dialSignal <- true:
			case <-ctx.Done():
				mq.log.Infof("filter worker canceled while sending dial signal")
				return
			}
		}
	}

	// Main worker; responds to dial calls, and feeds the current value through the output
	dialWorker := func() {
		mq.log.Debugf("MQ dial worker started")
		var conn *amqp.Conn
		var err error = fmt.Errorf("no connection yet attempted")

		for {
			select {
			case <-dialSignal:
				mq.log.Debugf("MQ dial worker received signal")
				conn, err = mq.dial(ctx)
				mq.log.Infof("MQ connection state: conn=%v err=%v", conn, err)
			case mq.connectionStateCh <- ConnectionState{Conn: conn, Error: err}:
			case <-ctx.Done():
				mq.log.Infof("dial worker canceled")
				return
			}
		}
	}

	go filterWorker()
	go dialWorker()

	mq.log.Infof("MQ connection management started")

	return nil
}

func (mq MQ) Reconnect(ctx context.Context) {
	var oldState, newState ConnectionState
	select {
	case oldState = <-mq.connectionStateCh:
	case <-ctx.Done():
		return
	}

	mq.log.Debugf("send MQ reconnect request")

	select {
	case mq.reconnectRequestedCh <- nil:
	case <-ctx.Done():
		return
	}

	for {
		select {
		case newState = <-mq.connectionStateCh:
		case <-ctx.Done():
			return
		}
		if &newState != &oldState {
			break
		}
	}
}

func (mq MQ) GetConnection(ctx context.Context) (*amqp.Conn, error) {
	var connState ConnectionState
	select {
	case connState = <-mq.connectionStateCh:
	case <-ctx.Done():
		return nil, fmt.Errorf("canceled")
	}
	if connState.Conn != nil && connState.Error == nil {
		return connState.Conn, nil
	}

	mq.log.Debugf("tried to get MQ connection but found bad state; reconnecting (err=%v)", connState.Error)
	mq.Reconnect(ctx)

	select {
	case connState = <-mq.connectionStateCh:
	case <-ctx.Done():
		return nil, fmt.Errorf("canceled")
	}
	if connState.Conn != nil && connState.Error == nil {
		return connState.Conn, nil
	}

	err := fmt.Errorf("connection state bad even after reconnection; %w", connState.Error)
	mq.log.Errorf("%v", err)
	return nil, err
}

func (mq MQ) GetSession(ctx context.Context) (*amqp.Session, error) {
	conn, err := mq.GetConnection(ctx)
	if err != nil {
		return nil, err
	}
	return conn.NewSession(ctx, nil)
}
