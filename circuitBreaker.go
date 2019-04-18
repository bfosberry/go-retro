package retro

import (
	"time"
)

// CircuitBreaker strategy
type CircuitBreaker struct {
	Error     error
	lastError error

	lastTry           time.Time
	currentAttempt    int
	currentBanTimeout time.Duration
	BanTimeout        time.Duration
	MaxAttempts       int

	Execute     func() error
	ShouldRetry func(error) bool
}

func (r *CircuitBreaker) setExecute(execute func() error) {
	r.Execute = execute
}

func (r *CircuitBreaker) getExecute() (func() error, error) {
	if r.Execute == nil {
		return nil, ErrExecuteFunctionNil
	}
	return r.Execute, nil
}

func (r *CircuitBreaker) increaseAttempt() error {
	r.currentAttempt++
	if r.currentAttempt >= r.MaxAttempts {
		return ErrMaxAttemptsReached
	}
	return nil
}

func (r *CircuitBreaker) increaseBanTimeout() (time.Duration, error) {
	r.currentBanTimeout = r.BanTimeout
	if r.BanTimeout == 0 {
		return 0, ErrBanTimeoutIsZero
	}
	return r.currentBanTimeout, nil
}

func (r *CircuitBreaker) setError(err error) {
	r.Error = err
}

func (r *CircuitBreaker) setLastError(err error) {
	r.lastError = err
}

func (r *CircuitBreaker) getLastError() error {
	return r.lastError
}

func (r *CircuitBreaker) resetState() {
	r.currentAttempt = 0
	r.currentBanTimeout = r.BanTimeout
	r.setError(nil)
}

func (r *CircuitBreaker) handleTick() error {
	if r.MaxAttempts == 0 {
		return ErrMaxAttemptsIsZero
	}
	currentBanTimeout, err := r.increaseBanTimeout()
	if err != nil {
		return err
	}
	if !(r.lastTry.IsZero()) {
		now := time.Now()
		if now.Sub(r.lastTry) > currentBanTimeout {
			r.resetState()
		}
	}
	r.lastTry = time.Now()
	if r.currentAttempt >= r.MaxAttempts {
		return ErrBanAttemptsReached
	}
	return nil
}

// Run whatever
func (r *CircuitBreaker) Run() (err error) {
	var execute func() error
	execute, err = r.getExecute()
	if err != nil {
		r.setError(err)
		return r.getLastError()
	}
	err = r.handleTick()
	if err != nil {
		r.setError(err)
		return r.getLastError()
	}
	err = r.increaseAttempt()
	if err != nil {
		r.setError(err)
	}
	err = execute()
	if err == nil {
		r.resetState()
	}
	r.setLastError(err)
	return r.getLastError()
}
