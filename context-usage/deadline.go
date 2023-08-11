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

package main

import "time"

var (
	deadline time.Time
)

func main() {
	deadline = time.Now().Add(1 * time.Minute)
	// goroutine 1
	go func() {
		for {
			select {
			case <-time.After(1 * time.Second):
				println("goroutine 1: doing some work that costs one second")
			default:
				// Check context
				if time.Now().After(deadline) {
					println("goroutine 1: cancel, time after deadline")
					return
				}
			}
		}
	}()
	// goroutine 2
	go func() {
		for {
			select {
			case <-time.After(1 * time.Second):
				println("goroutine 2: doing some work that costs one second")
			default:
				// Check context
				if time.Now().After(deadline) {
					println("goroutine 2: cancel, time after deadline")
					return
				}
			}
		}
	}()
	time.Sleep(2 * time.Second)
	println("main: context canceled")
	deadline = time.Now()
	time.Sleep(1 * time.Second)
	// do other things...
}
