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
	"time"
)

func main() {
	start := time.Now()

	userName := fetchUser()

	likes := fetchUserLikes(userName)
	match := fetchUserMatch(userName)

	fmt.Println("likes: ", likes)
	fmt.Println("match: ", match)

	// fetchUser: 100ms
	// fetchUserLikes: 150ms
	// fetchUserMatch: 100ms
	// took ~350 ms
	fmt.Println("took: ", time.Since(start))
}

func fetchUser() string {
	// HTTP round trip
	time.Sleep(time.Millisecond * 100)

	return "BOB"
}

func fetchUserLikes(userName string) int {
	time.Sleep(time.Millisecond * 150)

	return 123
}

func fetchUserMatch(userName string) string {
	time.Sleep(time.Millisecond * 100)

	return "ALICE"
}
