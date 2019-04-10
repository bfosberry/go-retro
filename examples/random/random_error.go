package main

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/codeship/go-retro"
)

const maxAttemps = 2

var (
	// ErrValueBelowThreshold is returned when the value is under the threshold.
	// This error is retryable and will retry 5 times without sleeping
	ErrValueBelowThreshold = retro.NewStaticRetryableError(errors.New("error: value is under the threshold"), maxAttemps, 0)
	errorThreshhold        = 15
	maxValue               = 13
)

// Stores all ints from 0-99 under random unique keys in a map
func main() {
	for i := 0; i < 100; i++ {
		var generatedValue int
		err := retro.DoWithRetry(func() (err error) {
			defer (func() {
				msg := fmt.Sprintf("try %dth chance.", i)
				if err != nil {
					msg = fmt.Sprintf("%s Value %d is below the given threshold", msg, generatedValue)
					if i == maxAttemps {
						msg = fmt.Sprintf("%s and reached maxAttemps", msg)
					}
				} else {
					msg = fmt.Sprintf("%s Value %d above the threshold - OK", msg, generatedValue)
				}
				fmt.Println(msg)
			})()
			generatedValue, err = genValue()
			return err
		})
		if err != nil {
			fmt.Printf("FATAL: Failed to execute genValue() %d: %s\n", i, err.Error())
		} else {
			fmt.Printf("OK: genValue() %d: %d\n", i, generatedValue)
		}
	}
}

func genValue() (int, error) {
	value := rand.Intn(errorThreshhold)
	if value < maxValue {
		return value, ErrValueBelowThreshold
	}
	return value, nil
}
