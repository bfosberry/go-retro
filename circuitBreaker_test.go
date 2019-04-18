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
			if strategy.Error != ErrMaxAttemptsReached {
				t.Error("strategy.Error != ErrMaxAttemptsReached")
			}
			if currentAttempt != 2 {
				t.Error("currentAttempt != 2")
			}
			continue
		}

		if i > 5 {
			break
		}

		// fourth, fifth and so on it should return errExpected strategy.ErrBanAttemptsReached
		// and ErrMaxAttemptsReached as the strategy is Banned
		if err != errExpected {
			t.Error("err != errExpected")
		}
		if strategy.Error != ErrBanAttemptsReached {
			t.Error("strategy.Error != ErrBanAttemptsReached")
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
			if strategy.Error != ErrMaxAttemptsReached {
				t.Error("strategy.Error != ErrMaxAttemptsReached")
			}
			if currentAttempt != 2 {
				t.Error("currentAttempt != 2")
			}
			continue
		}

		// fourth, fifth and so on it should return errExpected strategy.ErrBanAttemptsReached
		// and ErrMaxAttemptsReached as the strategy is Banned
		if err != errExpected {
			t.Error("err != errExpected")
		}
		if strategy.Error != ErrBanAttemptsReached {
			t.Error("strategy.Error != ErrBanAttemptsReached")
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
	if strategy.Error != ErrExecuteFunctionNil {
		t.Error("strategy.Error expected to be ErrExecuteFunctionNil")
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
	if strategy.Error != ErrBanTimeoutIsZero {
		t.Error("strategy.Error expected to be ErrBanTimeoutIsZero")
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
	if strategy.Error != ErrMaxAttemptsIsZero {
		t.Error("strategy.Error expected to be ErrMaxAttemptsIsZero")
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
			if strategy.Error != ErrMaxAttemptsReached {
				t.Error("strategy.Error != nil")
			}
			if currentAttempt != 1 {
				t.Error("currentAttempt != 1")
			}
			continue
		}

		// fourth, fifth and so on it should return errExpected strategy.ErrBanAttemptsReached
		// and ErrMaxAttemptsReached as the strategy is Banned
		if err != errExpected {
			t.Error("err != errExpected")
		}
		if strategy.Error != ErrBanAttemptsReached {
			t.Error("strategy.Error != ErrBanAttemptsReached")
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
		if strategy.Error != ErrMaxAttemptsReached {
			t.Error("strategy.Error != ErrMaxAttemptsReached")
		}
		if currentAttempt != 2 {
			t.Error("currentAttempt != 2")
		}

		err = strategy.Run()
		if err != errExpected {
			t.Error("err != errExpected")
		}
		if strategy.Error != ErrBanAttemptsReached {
			t.Error("strategy.Error != ErrBanAttemptsReached")
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
	if strategy.Error != ErrMaxAttemptsReached {
		t.Error("strategy.Error != nil")
	}
	if currentAttempt != 5 {
		t.Error("currentAttempt != 5")
	}

	err = strategy.Run()
	if err != errExpected {
		t.Error("err != nil")
	}
	if strategy.Error != ErrBanAttemptsReached {
		t.Error("strategy.Error != nil")
	}
	if currentAttempt != 5 {
		t.Error("currentAttempt != 5")
	}

	err = strategy.Run()
	if err != errExpected {
		t.Error("err != nil")
	}
	if strategy.Error != ErrBanAttemptsReached {
		t.Error("strategy.Error != nil")
	}
	if currentAttempt != 5 {
		t.Error("currentAttempt != 5")
	}
}

func TestTryFnStrategyCircuitBreakerViaTry(t *testing.T) {
	executed := false
	strategy := CircuitBreaker{
		MaxAttempts: 2,
		BanTimeout:  time.Second * 100,
	}

	err := Try(&strategy,
		func() error {
			executed = true
			return nil
		},
	)
	if err != nil {
		t.Error("Expecting no error")
	}
	if executed == false {
		t.Error("executed expected to be true")
	}
}
