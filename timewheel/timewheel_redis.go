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
 * The above copyright notice and this permission notice shall be included in all * copies or substantial portions of the Software.
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
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	http2 "github.com/Nicknamezz00/timewheel/pkg/http"
	"github.com/Nicknamezz00/timewheel/pkg/redis"
	"github.com/demdxx/gocast"
)

type RTask struct {
	Key         string            `json:"key"`
	CallbackURL string            `json:"callback_url"`
	Method      string            `json:"method"`
	Req         interface{}       `json:"req"`
	Header      map[string]string `json:"header"`
}

type RTimeWheel struct {
	sync.Once
	redisClient *redis.Client
	httpClient  *http2.Client
	stopCh      chan struct{}
	ticker      time.Ticker
}

func NewRTimeWheel(redisClient *redis.Client, httpClient *http2.Client) *RTimeWheel {
	return &RTimeWheel{
		redisClient: redisClient,
		httpClient:  httpClient,
		stopCh:      make(chan struct{}),
	}
}

func (r *RTimeWheel) Run() {
	r.ticker = *time.NewTicker(time.Second)
	go r.run()
}

func (r *RTimeWheel) Stop() {
	r.Do(func() {
		r.ticker.Stop()
		close(r.stopCh)
	})
}

func (r *RTimeWheel) AddTask(ctx context.Context, key string, task *RTask, executionTime time.Time) error {
	if err := r.checkTask(task); err != nil {
		return err
	}
	task.Key = key
	taskBody, _ := json.Marshal(task)
	_, err := r.redisClient.Eval(ctx, LuaAddTask, 2, []interface{}{
		r.getMinuteSlice(executionTime),  // minute-level zset timewheel slot
		r.getDeleteSetKey(executionTime), // set of tasks to be deleted
		executionTime.Unix(),             // timestamp as score
		string(taskBody),
		key,
	})
	return err
}

func (r *RTimeWheel) RemoveTask(ctx context.Context, key string, executionTime time.Time) error {
	_, err := r.redisClient.Eval(ctx, LuaDeleteTask, 1, []interface{}{
		r.getDeleteSetKey(executionTime),
		key,
	})
	return err
}

func (r *RTimeWheel) run() {
	for {
		select {
		case <-r.stopCh:
			return
		case <-r.ticker.C:
			go r.executeTasks()
		}
	}
}

func (r *RTimeWheel) executeTasks() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic at executeTasks")
		}
	}()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tasks, err := r.getExecutableTasks(ctx)
	if err != nil {
		log.Println("cannot get executable tasks")
		return
	}

	var wg sync.WaitGroup
	for _, task := range tasks {
		wg.Add(1)
		task := task
		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.Println("panic at for-loop tasks")
				}
				wg.Done()
			}()
			if err := r.execute(ctx, task); err != nil {
				log.Println("error at execute task")
			}
		}()
	}

	wg.Wait()
}

func (r *RTimeWheel) execute(ctx context.Context, task *RTask) error {
	return r.httpClient.JSONDo(ctx, task.Method, task.CallbackURL, task.Header, task.Req, nil)
}

func (r *RTimeWheel) checkTask(task *RTask) error {
	if task.Method != http.MethodGet && task.Method != http.MethodPost {
		return fmt.Errorf("invalid method: %s", task.Method)
	}
	if !strings.HasPrefix(task.CallbackURL, "http://") && !strings.HasPrefix(task.CallbackURL, "https://") {
		return fmt.Errorf("invalid url: %s", task.CallbackURL)
	}
	return nil
}

func (r *RTimeWheel) getExecutableTasks(ctx context.Context) ([]*RTask, error) {
	now := time.Now()
	minuteSlice := r.getMinuteSlice(now)
	deleteSetKey := r.getDeleteSetKey(now)
	nowSecond := GetTimeSecond(now)
	score1 := nowSecond.Unix()
	score2 := nowSecond.Add(time.Second).Unix()
	rawReply, err := r.redisClient.Eval(ctx, LuaZRangeTasks, 2, []interface{}{
		minuteSlice, deleteSetKey, score1, score2,
	})
	if err != nil {
		return nil, err
	}
	replies := gocast.ToInterfaceSlice(rawReply)
	if len(replies) == 0 {
		return nil, fmt.Errorf("invalid replies: %v", replies)
	}
	delete := gocast.ToStringSlice(replies[0])
	deletedSet := make(map[string]struct{}, len(delete))

	tasks := make([]*RTask, 0, len(replies)-1)
	for i := 1; i < len(replies); i++ {
		var t RTask
		if err := json.Unmarshal([]byte(gocast.ToString(replies[i])), &t); err != nil {
			log.Printf("error at unmarshal: %v", err)
			continue
		}
		if _, ok := deletedSet[t.Key]; ok {
			continue
		}
		tasks = append(tasks, &t)
	}
	return tasks, nil
}

func GetTimeStr(t time.Time) string {
	return t.Format("2006-01-02-15:04")
}

func GetTimeSecond(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.Local)
}

func (r *RTimeWheel) getMinuteSlice(executionTime time.Time) string {
	return fmt.Sprintf("timewheel_redis_task_{%s}", GetTimeStr(executionTime))
}

func (r *RTimeWheel) getDeleteSetKey(executionTime time.Time) string {
	return fmt.Sprintf("timewheel_redis_delete_set_{%s}", GetTimeStr(executionTime))
}
