package retro

import (
	"errors"
	"testing"
	"time"
)

func TestTryRetryWaitIncreaseFnStrategyMaxAttemptZeroDelay(t *testing.T) {
	currentAttempt := 0
	errExpected := errors.New("expected error")

	strategy := RetryWaitIncrease{
		MaxAttempts: 5,
		Execute: func() error {
			currentAttempt++
			return errExpected
		},
	}

	errActual := strategy.Run()

	if errActual != errExpected {
		t.Error("fn returns an unexpected error")
	}
	if strategy.Error == nil {
		t.Error("strategy must hold an error")
	}
	if !(strategy.Error == ErrDelayIsZero) {
		t.Error("strategy error should be Delay Is Zero")
	}
}

func TestTryRetryWaitIncreaseFnStrategyMaxAttemptNonZeroDelay(t *testing.T) {
	currentAttempt := 0
	errExpected := errors.New("expected error")
	startTime := time.Now()
	delay := time.Millisecond * 50

	strategy := RetryWaitIncrease{
		MaxAttempts: 5,
		Delay:       delay,
		Execute: func() error {
			currentAttempt++
			endTime := time.Now()
			if endTime.Sub(startTime) < delay*time.Duration(currentAttempt-1) {
				t.Error("Incorrect increase delay")
			}
			startTime = time.Now()
			return errExpected
		},
	}

	errActual := strategy.Run()

	if errActual != errExpected {
		t.Error("fn returns an unexpected error")
	}
	if strategy.Error == nil {
		t.Error("strategy must hold an error")
	}
	if !(strategy.Error == ErrMaxAttemptsReached) {
		t.Error("strategy error should be Delay Is Zero")
	}
}

func TestTryRetryWaitIncreaseFnStrategyMaxAttemptNonZeroDelayViaTry(t *testing.T) {
	currentAttempt := 0
	errExpected := errors.New("expected error")
	startTime := time.Now()
	delay := time.Millisecond * 50

	strategy := RetryWaitIncrease{
		MaxAttempts: 5,
		Delay:       delay,
	}

	errActual := Try(&strategy,
		func() error {
			currentAttempt++
			endTime := time.Now()
			if endTime.Sub(startTime) < delay*time.Duration(currentAttempt-1) {
				t.Error("Incorrect increase delay")
			}
			startTime = time.Now()
			return errExpected
		},
	)

	if errActual != errExpected {
		t.Error("fn returns an unexpected error")
	}
	if strategy.Error == nil {
		t.Error("strategy must hold an error")
	}
	if !(strategy.Error == ErrMaxAttemptsReached) {
		t.Error("strategy error should be Delay Is Zero")
	}
}
