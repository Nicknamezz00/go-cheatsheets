package main

import (
	"fmt"
	"sync"
)

const n = 5

var (
	wg   sync.WaitGroup
	mu   sync.Mutex
	cond = sync.NewCond(&mu)
	cnt  int
)

func print(id int) {
	defer wg.Done()
	for {
		mu.Lock()
		if cnt > 100 {
			mu.Unlock()
			break
		}
		for cnt%n != id {
			cond.Wait()
		}
		if cnt <= 100 {
			fmt.Printf("goroutine %d prints out %d\n", id+1, cnt)
		}
		cnt++
		cond.Broadcast()
		mu.Unlock()
	}
}

func main() {
	for i := 0; i < n; i++ {
		wg.Add(1)
		go print(i)
	}
	wg.Wait()
	fmt.Println("all done")
}
