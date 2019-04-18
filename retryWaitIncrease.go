package retro

import "time"

// RetryWaitIncrease strategy
type RetryWaitIncrease struct {
	currentAttempt int
	currentDelay   time.Duration

	// Error holds the error ocurred in the RetryWaitIncrease' scope
	Error error

	// MaxAttempts defines the maximum number of attempts
	MaxAttempts int

	// Delay defines the amount of time between attempts
	Delay time.Duration

	// Execute defines the code to be wrapped under this strategy
	Execute func() error

	// ShouldRetry defines the function that states if execute the next attempt or not
	ShouldRetry func(error) bool
}

func (r *RetryWaitIncrease) setError(err error) {
	r.Error = err
}

func (r *RetryWaitIncrease) increaseDelay() (time.Duration, error) {
	r.currentDelay = r.Delay * time.Duration(r.currentAttempt)
	if r.Delay == 0 {
		return 0, ErrDelayIsZero
	}
	return r.currentDelay, nil
}

func (r *RetryWaitIncrease) increaseAttempt() error {
	r.currentAttempt++
	if r.currentAttempt >= r.MaxAttempts {
		return ErrMaxAttemptsReached
	}
	return nil
}

func (r *RetryWaitIncrease) setExecute(execute func() error) {
	r.Execute = execute
}

func (r *RetryWaitIncrease) getExecute() (func() error, error) {
	if r.Execute == nil {
		return nil, ErrExecuteFunctionNil
	}
	return r.Execute, nil
}

func (r *RetryWaitIncrease) getShouldRetry() func(error) bool {
	return r.ShouldRetry
}

// Run the execute function and behave according to the strategy
func (r *RetryWaitIncrease) Run() error {
	return launcStrategy(r)
}
