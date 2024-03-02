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

import (
	"errors"
	"time"
)

type Message struct {
	content string
}

type Subscription struct {
	Inbox chan Message
	ch    chan Message
}

func (s *Subscription) Publish(msg Message) error {
	if _, ok := <-s.ch; !ok {
		return errors.New("topic has been closed")
	}
	s.ch <- msg
}

type Topic struct {
	Subscribers    []Session
	MessageHistory []Message
}

func (t *Topic) Subscribe(uid uint64) (Subscription, error) {
	// Get session and create one if it's the first
	// Add session to the Topic & MessageHistory
	// Create a subscription
	return Subscription{nil, nil}, nil
}

func (t *Topic) Unsubscribe(Subscription) error {
	// implementation
	return nil
}

func (t *Topic) Delete() error {
	// implementation
	return nil
}

type User struct {
	ID   uint64
	Name string
}

type Session struct {
	User      User
	Timestamp time.Time
}
