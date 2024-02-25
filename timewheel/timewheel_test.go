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

package timewheel

import (
	"context"
	"testing"
	"time"

	http2 "github.com/Nicknamezz00/timewheel/pkg/http"
	"github.com/Nicknamezz00/timewheel/pkg/redis"
)

func TestTimeWheel(t *testing.T) {
	tw := NewTimeWheel(10, 500*time.Millisecond)
	tw.Run()
	defer tw.Stop()

	done := make(chan struct{}, 3)

	tw.AddTask("test1", func() {
		t.Logf("executing test1 at %v", time.Now())
		done <- struct{}{}
	}, time.Now().Add(time.Second))
	tw.AddTask("test2", func() {
		t.Logf("executing test2 at %v", time.Now())
		done <- struct{}{}
	}, time.Now().Add(5*time.Second))
	tw.AddTask("test3", func() {
		t.Logf("executing test3 at %v", time.Now())
		done <- struct{}{}
	}, time.Now().Add(3*time.Second))

	<-time.After(6 * time.Second)
	if len(done) != 3 {
		t.Fatalf("%d tasks finished, expect 3", len(done))
	}
}

const (
	network  = "tcp"
	address  = "127.0.0.1:6379"
	password = ""
)

var (
	callbackURL    = ""
	callbackMethod = "POST"
	callbackReq    interface{}
	callbackHeader map[string]string
)

func TestTimeWheelRedis(t *testing.T) {
	rtw := NewRTimeWheel(redis.NewClient(network, address, password), http2.NewClient())
	defer rtw.Stop()
	rtw.Run()

	ctx := context.Background()
	if err := rtw.AddTask(ctx, "test1", &RTask{
		CallbackURL: callbackURL,
		Method:      callbackMethod,
		Req:         callbackReq,
		Header:      callbackHeader,
	}, time.Now().Add(time.Second)); err != nil {
		t.Error(err)
		return
	}
	if err := rtw.AddTask(ctx, "test2", &RTask{
		CallbackURL: callbackURL,
		Method:      callbackMethod,
		Req:         callbackReq,
		Header:      callbackHeader,
	}, time.Now().Add(4*time.Second)); err != nil {
		t.Error(err)
		return
	}

	if err := rtw.RemoveTask(ctx, "test2", time.Now().Add(4*time.Second)); err != nil {
		t.Error(err)
		return
	}
	<-time.After(5 * time.Second)
}
