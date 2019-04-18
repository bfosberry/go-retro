package retro

import (
	"errors"
	"time"
)

// ErrBanAttemptsReached occurs when the execute function was executed BanAttempts times
var ErrBanAttemptsReached = errors.New("BanAttempts reached")

// ErrShouldRetryFunctionError occurs when ShouldRetry function returns false
var ErrShouldRetryFunctionError = errors.New("ShouldRetry returned false")

// ErrRecoveryFunctionError occurs when Recovery function returns an error
var ErrRecoveryFunctionError = errors.New("Recovery function finished with error")

// ErrExecuteFunctionNil occurs when Run() is executed without Execute function defined
var ErrExecuteFunctionNil = errors.New("Execute function has not beed defined")

// ErrMaxAttemptsReached occurs when the execute function was executed MaxAttempts times
var ErrMaxAttemptsReached = errors.New("MaxAttempts reached")

// ErrMaxAttemptsIsZero occurs as Max Attempts is supposed to be > 0
var ErrMaxAttemptsIsZero = errors.New("Max Attempts is zero")

// ErrDelayIsZero occurs as Delay is supposed to be > 0
var ErrDelayIsZero = errors.New("Delay is zero")

// ErrBanTimeoutIsZero occurs as Ban Timeout is supposed to be > 0
var ErrBanTimeoutIsZero = errors.New("Ban Timeout is zero")

type strategyRetry interface {
	getExecute() (func() error, error)
	getShouldRetry() func(error) bool
	increaseAttempt() error
	increaseDelay() (time.Duration, error)
	setError(error)
}

func launcStrategy(s strategyRetry) error {
	execute, err := s.getExecute()
	if err != nil {
		s.setError(err)
		return err
	}
	for {
		err := s.increaseAttempt()
		errExecute := execute()
		if errExecute != nil {
			if err != nil {
				s.setError(err)
				return errExecute
			}

			duration, err := s.increaseDelay()
			if err != nil {
				s.setError(err)
				return errExecute
			}
			shouldRetry := s.getShouldRetry()
			if shouldRetry != nil {
				if shouldRetry(errExecute) == false {
					s.setError(ErrShouldRetryFunctionError)
					return errExecute
				}
			}
			if duration > 0 {
				time.Sleep(duration)
			}
			continue
		}
		break
	}
	return nil
}

// Strategy whatever
type Strategy interface {
	setExecute(execute func() error)
	Run() error
}

// Try whatever
func Try(strategy Strategy, execute func() error) error {
	strategy.setExecute(execute)
	return strategy.Run()
}
