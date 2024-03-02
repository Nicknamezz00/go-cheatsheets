/*
 * MIT License
 *
 * Copyright (c) 2023 Runze Wu
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 */

package circuitbreaker

import (
	"context"
	"errors"
	"time"
)

type State int

const (
	UnknownState = iota
	FailureState
	SuccessState
)

type Counter interface {
	Count(State)
	ConsecutiveFailures() uint32
	LastActivity() time.Time
	Reset()
}

func NewCounter() Counter {
	// implement
	return nil
}

type Circuit func(context.Context) error

func Breaker(circuit Circuit, failureThreshlold uint32) Circuit {
	cnt := NewCounter()
	return func(ctx context.Context) error {
		if cnt.ConsecutiveFailures() >= failureThreshlold {
			canRetry := func(cnt Counter) bool {
				backoffLevel := cnt.ConsecutiveFailures() - failureThreshlold
				// Calculates when should the circuit breaker resume propagating requests to the service
				shouldRetryAt := cnt.LastActivity().Add(time.Second * 2 << backoffLevel)
				return time.Now().After(shouldRetryAt)
			}
			if !canRetry(cnt) {
				return errors.New("error service unavailable")
			}
		}
		if err := circuit(ctx); err != nil {
			cnt.Count(FailureState)
			return err
		}
		cnt.Count(SuccessState)
		return nil
	}
}
