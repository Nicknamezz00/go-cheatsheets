To avoid bugs, especially lost wakeup,

whenever use condition variables, follow this pattern:

```go
mu.Lock()
// do something that might affect the condition
cond.Broadcast()
```
---
Another place:
```go
mu.Lock()
while condition == false {
  cond.Wait()
}
// now condition is true, and we still have the lock
mu.Unlock()
```
