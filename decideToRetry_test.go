package retro

import (
	"errors"
	"testing"
	"time"
)

func TestTryFnStrategyMaxAttemptWithShouldRetryFunction(t *testing.T) {
	currentAttempt := 0
	currentShouldRetry := 0
	errExpected := errors.New("expected error")

	strategy := Retry{
		MaxAttempts: 5,
		Delay:       time.Millisecond * 100,
		Execute: func() error {
			currentAttempt++
			return errExpected
		},
		ShouldRetry: func(err error) bool {
			currentShouldRetry++
			return true
		},
	}

	errActual := strategy.Run()

	if errActual != errExpected {
		t.Error("fn returns an unexpected error")
	}
	if strategy.Error == nil {
		t.Error("strategy must hold an error")
	}
	if strategy.MaxAttempts != currentAttempt {
		t.Error("currentAttempt differs from MaxAttempts")
	}
	if strategy.MaxAttempts-1 != currentShouldRetry {
		t.Error("currentShouldRetry differs from MaxAttempts", currentShouldRetry)
	}
}

func TestTryFnStrategyMaxAttemptWithShouldRetryFunctionTilSecondAttempt(t *testing.T) {
	currentAttempt := 0
	currentShouldRetry := 0
	errExpected := errors.New("expected error")

	strategy := Retry{
		MaxAttempts: 5,
		Delay:       time.Millisecond * 100,
		Execute: func() error {
			currentAttempt++
			return errExpected
		},
		ShouldRetry: func(err error) bool {
			currentShouldRetry++
			if currentShouldRetry == 2 {
				return false
			}
			return true
		},
	}

	errActual := strategy.Run()

	if errActual != errExpected {
		t.Error("fn returns an unexpected error")
	}
	if strategy.Error == nil {
		t.Error("strategy must hold an error")
	}
	if strategy.Error != ErrShouldRetryFunctionError {
		t.Error("strategy.error differs from ErrShouldRetryFunctionError")
	}
	if currentAttempt != 2 {
		t.Error("currentAttempt differs from expected", currentAttempt)
	}
	if currentShouldRetry != 2 {
		t.Error("currentShouldRetry differs from expected", currentShouldRetry)
	}
}
