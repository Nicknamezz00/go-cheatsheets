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

package messaging

// Split splits a channel into n channels that receive message in
// round-robin style.
func Split(bigChan <-chan int, n int) []chan<- int {
	cs := make([]chan<- int, n)
	for i := 0; i < n; i++ {
		cs = append(cs, make(chan int))
	}
	// Distribute the work in a round robin style among the stated number
	// of channels until the main channel has been closed. In that case, close
	// all channels and return.
	distributeToChannels := func(ch <-chan int, cs []chan<- int) {
		// Close all channels when we're done.
		defer func(cs []chan<- int) {
			for _, c := range cs {
				close(c)
			}
		}(cs)
		for {
			for _, c := range cs {
				v, ok := <-bigChan
				if !ok {
					return
				}
				c <- v
			}
		}
	}
	go distributeToChannels(bigChan, cs)
	return cs
}
