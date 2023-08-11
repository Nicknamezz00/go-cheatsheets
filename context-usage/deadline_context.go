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

import (
	"context"
	"time"
)

func main() {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	defer cancel()
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				println("goroutine 1: cancel, time after deadline")
				return
			default:
				time.Sleep(1 * time.Second)
				println("goroutine 1: doing some work that costs one second")
			}
		}
	}(ctx)
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				println("goroutine 2: cancel, time after deadline")
				return
			default:
				time.Sleep(1 * time.Second)
				println("goroutine 2: doing some work that costs one second")
			}
		}
	}(ctx)
	time.Sleep(2 * time.Second)
	println("main: context canceled")
	time.Sleep(1 * time.Second)
	// do other things...
}
