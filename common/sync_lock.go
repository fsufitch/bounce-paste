package common

import "errors"

var ErrLocked = errors.New("mutex is locked")
var ErrUnlocked = errors.New("mutex is unlocked")

type Lock chan any

func (m Lock) IsLocked() bool {
	if m == nil {
		return false
	}
	select {
	// Try to "peek" the value in the channel; if it works, then there's a lock in there
	case m <- <-m:
		return true
	default:
		return false
	}
}

func (m *Lock) Acquire() (func() error, error) {
	if *m == nil {
		*m = make(chan any, 1)
	}
	select {
	case *m <- nil:
		return func() error {
			select {
			case <-*m:
				return nil
			default:
				return ErrUnlocked
			}

		}, nil
	default:
		return nil, ErrLocked
	}
}
