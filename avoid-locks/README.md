## Usage
Avoid locks (Mutex) with actor model

## Test
```shell
go test ./... --race -v
```
Everything should be fine.

Of course, lock is always fine, but we don't want to see locks all over you system.

Most of the time you will not notice the problem of this lock, somebody will have a lock, somebody else want
to adjust a state, but they can't, because it's still locked, and they have to wait. 
It could be such a mess, and you will have actually no clue what is going on and your program is running slow, 
and how are you
going to debug that well?
