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
	"container/list"
	"log"
	"sync"
	"time"
)

//reference: https://github.com/HDT3213/godis/blob/master/lib/timewheel/timewheel.go

const (
	DefaultSlotNumber   = 10
	DefaultTimeInterval = time.Second
)

type task struct {
	job      func()
	position int
	cycle    int
	key      string
}

type TimeWheel struct {
	sync.Once
	interval        time.Duration
	ticker          *time.Ticker
	addTaskCh       chan *task
	removeTaskCh    chan string
	slots           []*list.List
	currentPosition int
	stopCh          chan struct{}
	keyToElementMap map[string]*list.Element
}

func NewTimeWheel(slotsNum int, interval time.Duration) *TimeWheel {
	if slotsNum <= 0 {
		slotsNum = DefaultSlotNumber
	}
	if interval <= 0 {
		interval = DefaultTimeInterval
	}
	tw := TimeWheel{
		interval:        interval,
		addTaskCh:       make(chan *task),
		removeTaskCh:    make(chan string),
		slots:           make([]*list.List, 0, slotsNum),
		stopCh:          make(chan struct{}),
		keyToElementMap: make(map[string]*list.Element),
	}
	for i := 0; i < slotsNum; i++ {
		tw.slots = append(tw.slots, list.New())
	}
	return &tw
}

func (t *TimeWheel) Run() {
	t.ticker = time.NewTicker(t.interval)
	go t.run()
}

func (t *TimeWheel) run() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	for {
		select {
		case <-t.stopCh:
			return
		case <-t.ticker.C:
			t.tick()
		case task := <-t.addTaskCh:
			t.addTask(task)
		case key := <-t.removeTaskCh:
			t.RemoveTask(key)
		}
	}
}

func (t *TimeWheel) AddTask(key string, job func(), executionTime time.Time) {
	pos, cycle := t.getPositionAndCycle(executionTime)
	t.addTaskCh <- &task{
		job:      job,
		position: pos,
		cycle:    cycle,
		key:      key,
	}
}

func (t *TimeWheel) RemoveTask(key string) {
	t.removeTaskCh <- key
}

func (t *TimeWheel) Stop() {
	t.Do(func() {
		t.ticker.Stop()
		close(t.stopCh)
	})
}

func (t *TimeWheel) tick() {
	defer func() {
		t.currentPosition = (t.currentPosition + 1) % len(t.slots)
	}()
	l := t.slots[t.currentPosition]
	t.execute(l)
}

func (t *TimeWheel) addTask(task *task) {
	l := t.slots[task.position]
	if _, ok := t.keyToElementMap[task.key]; ok {
		t.removeTask(task.key)
	}
	element := l.PushBack(task)
	t.keyToElementMap[task.key] = element
}

func (t *TimeWheel) removeTask(key string) {
	element, ok := t.keyToElementMap[key]
	if !ok {
		return
	}
	delete(t.keyToElementMap, key)
	task, _ := element.Value.(*task)
	_ = t.slots[task.position].Remove(element)
}

func (t *TimeWheel) execute(l *list.List) {
	for e := l.Front(); e != nil; {
		taskElement := e.Value.(*task)
		if taskElement.cycle > 0 {
			// not yet, skip this
			taskElement.cycle--
			e = e.Next()
			continue
		}

		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.Println(err)
				}
			}()
			taskElement.job()
		}()

		// delete it after we're done
		next := e.Next()
		l.Remove(e)
		delete(t.keyToElementMap, taskElement.key)
		e = next
	}
}

func (t *TimeWheel) getPositionAndCycle(executionTime time.Time) (int, int) {
	delay := int(time.Until(executionTime))
	cycle := delay / (len(t.slots) * int(t.interval))
	pos := (t.currentPosition + delay/int(t.interval)) % len(t.slots)
	return pos, cycle
}
