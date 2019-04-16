package retro

import (
	"errors"
	"testing"
	"time"
)

func TestTryFnStrategyCircuitBreaker(t *testing.T) {
	currentAttempt := 0
	errExpected := errors.New("expected error")

	strategy := CircuitBreaker{
		MaxAttempts: 2,
		BanTimeout:  time.Second * 1,
		Execute: func() error {
			currentAttempt++
			return errExpected
		},
	}

	i := 0
	for {
		i++

		// What about if run the same strategy lot of times?
		err := strategy.Run()

		// first time it should return errExpected and strategy.Error = Nil
		if i == 1 {
			if err != errExpected {
				t.Error("err != errExpected")
			}
			if strategy.Error != nil {
				t.Error("strategy.Error != nil")
			}
			if currentAttempt != 1 {
				t.Error("currentAttempt != 1")
			}
			continue
		}
		// second time it should return errExpected and strategy.Error = Nil
		if i == 2 {
			if err != errExpected {
				t.Error("err != errExpected")
			}
			if strategy.Error != ErrorMaxAttemptsReached {
				t.Error("strategy.Error != ErrorMaxAttemptsReached")
			}
			if currentAttempt != 2 {
				t.Error("currentAttempt != 2")
			}
			continue
		}

		if i > 5 {
			break
		}

		// fourth, fifth and so on it should return errExpected strategy.ErrorBanAttemptsReached
		// and ErrorMaxAttemptsReached as the strategy is Banned
		if err != errExpected {
			t.Error("err != errExpected")
		}
		if strategy.Error != ErrorBanAttemptsReached {
			t.Error("strategy.Error != ErrorBanAttemptsReached")
		}
		if currentAttempt != 2 {
			t.Error("currentAttempt != 2")
		}
	}
}

func TestTryFnStrategyCircuitBreakerRunAfterBan(t *testing.T) {
	currentAttempt := 0
	errExpected := errors.New("expected error")

	strategy := CircuitBreaker{
		MaxAttempts: 2,
		BanTimeout:  time.Millisecond * 100,
		Execute: func() error {
			currentAttempt++
			return errExpected
		},
	}

	i := 0
	for {
		i++

		// What about if run the same strategy lot of times?
		err := strategy.Run()

		// first time it should return errExpected and strategy.Error = Nil
		if i == 1 {
			if err != errExpected {
				t.Error("err != errExpected")
			}
			if strategy.Error != nil {
				t.Error("strategy.Error != nil")
			}
			if currentAttempt != 1 {
				t.Error("currentAttempt != 1")
			}
			continue
		}
		// second time it should return errExpected and strategy.Error = Nil
		if i == 2 {
			if err != errExpected {
				t.Error("err != errExpected")
			}
			if strategy.Error != ErrorMaxAttemptsReached {
				t.Error("strategy.Error != ErrorMaxAttemptsReached")
			}
			if currentAttempt != 2 {
				t.Error("currentAttempt != 2")
			}
			continue
		}

		// fourth, fifth and so on it should return errExpected strategy.ErrorBanAttemptsReached
		// and ErrorMaxAttemptsReached as the strategy is Banned
		if err != errExpected {
			t.Error("err != errExpected")
		}
		if strategy.Error != ErrorBanAttemptsReached {
			t.Error("strategy.Error != ErrorBanAttemptsReached")
		}
		if currentAttempt != 2 {
			t.Error("currentAttempt != 2")
		}

		time.Sleep(strategy.BanTimeout)
		time.Sleep(time.Millisecond * 10)

		err = strategy.Run()
		if err != errExpected {
			t.Error("err != errExpected")
		}
		if strategy.Error != nil {
			t.Error("strategy.Error != nil")
		}
		if currentAttempt != 3 {
			t.Error("currentAttempt != 3")
		}

		break
	}
}

func TestTryFnStrategyCircuitBreakerWithoutExecuteFn(t *testing.T) {
	strategy := CircuitBreaker{
		MaxAttempts: 2,
		BanTimeout:  time.Second * 1,
	}

	err := strategy.Run()
	if err != nil {
		t.Error("Expecting no error")
	}
	if strategy.Error != ErrorExecuteFunctionNil {
		t.Error("strategy.Error expected to be ErrorExecuteFunctionNil")
	}
}

func TestTryFnStrategyCircuitBreakerWithoutBanTimeout(t *testing.T) {
	strategy := CircuitBreaker{
		MaxAttempts: 2,
		Execute: func() error {
			return nil
		},
	}

	err := strategy.Run()
	if err != nil {
		t.Error("Expecting no error")
	}
	if strategy.Error != ErrorBanTimeoutIsZero {
		t.Error("strategy.Error expected to be ErrorBanTimeoutIsZero")
	}
}

func TestTryFnStrategyCircuitBreakerWithoutMaxAttemps(t *testing.T) {
	strategy := CircuitBreaker{
		BanTimeout: time.Second,
		Execute: func() error {
			return nil
		},
	}

	err := strategy.Run()
	if err != nil {
		t.Error("Expecting no error")
	}
	if strategy.Error != ErrorMaxAttemptsIsZero {
		t.Error("strategy.Error expected to be ErrorMaxAttempsIsZero")
	}
}

func TestTryFnStrategyCircuitBreakerRunMaxAttemptsEqualsTo1(t *testing.T) {
	currentAttempt := 0
	errExpected := errors.New("expected error")

	strategy := CircuitBreaker{
		MaxAttempts: 1,
		BanTimeout:  time.Millisecond * 100,
		Execute: func() error {
			currentAttempt++
			return errExpected
		},
	}

	i := 0
	for {
		i++

		// What about if run the same strategy lot of times?
		err := strategy.Run()

		// first time it should return errExpected and strategy.Error = Nil
		if i == 1 {
			if err != errExpected {
				t.Error("err != errExpected")
			}
			if strategy.Error != ErrorMaxAttemptsReached {
				t.Error("strategy.Error != nil")
			}
			if currentAttempt != 1 {
				t.Error("currentAttempt != 1")
			}
			continue
		}

		// fourth, fifth and so on it should return errExpected strategy.ErrorBanAttemptsReached
		// and ErrorMaxAttemptsReached as the strategy is Banned
		if err != errExpected {
			t.Error("err != errExpected")
		}
		if strategy.Error != ErrorBanAttemptsReached {
			t.Error("strategy.Error != ErrorBanAttemptsReached")
		}
		if currentAttempt != 1 {
			t.Error("currentAttempt != 1")
		}

		time.Sleep(strategy.BanTimeout)
		time.Sleep(time.Millisecond * 10)

		err = strategy.Run()
		if err != errExpected {
			t.Error("err != errExpected")
		}
		if strategy.Error != ErrorMaxAttemptsReached {
			t.Error("strategy.Error != ErrorMaxAttemptsReached")
		}
		if currentAttempt != 2 {
			t.Error("currentAttempt != 2")
		}

		err = strategy.Run()
		if err != errExpected {
			t.Error("err != errExpected")
		}
		if strategy.Error != ErrorBanAttemptsReached {
			t.Error("strategy.Error != ErrorBanAttemptsReached")
		}
		if currentAttempt != 2 {
			t.Error("currentAttempt != 2")
		}

		break
	}
}

func TestTryFnStrategyCircuitBreakerRunResetAfterSuccess(t *testing.T) {
	currentAttempt := 0
	errExpected := errors.New("expected error")

	strategy := CircuitBreaker{
		MaxAttempts: 3,
		BanTimeout:  time.Second * 100,
		Execute: func() error {
			currentAttempt++
			if currentAttempt == 2 {
				return nil
			}
			return errExpected
		},
	}

	err := strategy.Run()
	if err != errExpected {
		t.Error("err != errExpected")
	}
	if strategy.Error != nil {
		t.Error("strategy.Error != nil")
	}
	if currentAttempt != 1 {
		t.Error("currentAttempt != 1")
	}

	err = strategy.Run()
	if err != nil {
		t.Error("err != errExpected")
	}
	if strategy.Error != nil {
		t.Error("strategy.Error != nil")
	}
	if currentAttempt != 2 {
		t.Error("currentAttempt != 2")
	}

	err = strategy.Run()
	if err != errExpected {
		t.Error("err != nil")
	}
	if strategy.Error != nil {
		t.Error("strategy.Error != nil")
	}
	if currentAttempt != 3 {
		t.Error("currentAttempt != 3")
	}

	err = strategy.Run()
	if err != errExpected {
		t.Error("err != nil")
	}
	if strategy.Error != nil {
		t.Error("strategy.Error != nil")
	}
	if currentAttempt != 4 {
		t.Error("currentAttempt != 4")
	}

	err = strategy.Run()
	if err != errExpected {
		t.Error("err != nil")
	}
	if strategy.Error != ErrorMaxAttemptsReached {
		t.Error("strategy.Error != nil")
	}
	if currentAttempt != 5 {
		t.Error("currentAttempt != 5")
	}

	err = strategy.Run()
	if err != errExpected {
		t.Error("err != nil")
	}
	if strategy.Error != ErrorBanAttemptsReached {
		t.Error("strategy.Error != nil")
	}
	if currentAttempt != 5 {
		t.Error("currentAttempt != 5")
	}

	err = strategy.Run()
	if err != errExpected {
		t.Error("err != nil")
	}
	if strategy.Error != ErrorBanAttemptsReached {
		t.Error("strategy.Error != nil")
	}
	if currentAttempt != 5 {
		t.Error("currentAttempt != 5")
	}
}
