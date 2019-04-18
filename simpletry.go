package retro

import "time"

// Retry strategy
type Retry struct {
	currentAttempt int
	currentDelay   time.Duration

	// Error holds the error ocurred in the Retry' scope
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

func (r *Retry) setError(err error) {
	r.Error = err
}

func (r *Retry) increaseDelay() (time.Duration, error) {
	r.currentDelay = r.Delay
	return r.currentDelay, nil
}

func (r *Retry) increaseAttempt() error {
	r.currentAttempt++
	if r.currentAttempt >= r.MaxAttempts {
		return ErrMaxAttemptsReached
	}
	return nil
}

func (r *Retry) setExecute(execute func() error) {
	r.Execute = execute
}

func (r *Retry) getExecute() (func() error, error) {
	if r.Execute == nil {
		return nil, ErrExecuteFunctionNil
	}
	return r.Execute, nil
}

func (r *Retry) getShouldRetry() func(error) bool {
	return r.ShouldRetry
}

// Run the execute function and behave according to the strategy
func (r *Retry) Run() error {
	return launcStrategy(r)
}
