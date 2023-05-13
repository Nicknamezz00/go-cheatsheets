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
	"fmt"
	"sync"
	"time"
)

func fetchUser() string {
	// HTTP round trip
	time.Sleep(time.Millisecond * 100)

	return "BOB"
}

func fetchUserLikes(userName string, respCh chan any, wg *sync.WaitGroup) {
	time.Sleep(time.Millisecond * 150)

	respCh <- 11
	wg.Done()
}

func fetchUserMatch(userName string, respCh chan any, wg *sync.WaitGroup) {
	time.Sleep(time.Millisecond * 100)

	respCh <- "ALICE"
	wg.Done()
}

func main() {
	start := time.Now()

	userName := fetchUser()
	respCh := make(chan any, 10) // make it buffered

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go fetchUserLikes(userName, respCh, wg)
	go fetchUserMatch(userName, respCh, wg)
	wg.Wait() // block until we have 2 wg.Done()

	// According to line 73-74, we close the channel.
	// But who guarantees that the two goroutines above are done?
	// WaitGroup guarantees.
	close(respCh)

	// This approach causes deadlock, because channel keeps waiting for inputs,
	// but nothing ever comes, so it just keeps waiting forever.
	for resp := range respCh {
		fmt.Println("response: ", resp)
	}

	// fetchUser: 100ms
	// max(fetchUserLikes, fetchUserMatch): 150ms
	// took ~250 ms
	fmt.Println("took: ", time.Since(start))
}
