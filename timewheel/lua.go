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

const (
	// If the pending task is in delete set, delete it.
	// Then add it by score.
	LuaAddTask = `
		local zset_key = KEYS[1]
		local delete_set_key = KEYS[2]
		local score = ARGV[1]
		local task = ARGV[2]
		local task_key = ARGV[3]
		redis.call('srem', delete_set_key, task_key)
		return redis.call('zadd', zset_key, score, task)
	`

	LuaDeleteTask = `
		local delete_set_key = KEYS[1]
		local task_key = ARGV[1]
		redis.call('sadd', delete_set_key, task_key)
		local scnt = redis.call('scard', delete_set_key)
		if (tonumber(scnt) == 1)
		then
			redis.call('expire', delete_set_key, 120)
		end
		return scnt
	`

	LuaZRangeTasks = `
		local zset_key = KEYS[1]
		local delete_set_key = KEYS[2]
		local score1 = ARGV[1]
		local score2 = ARGV[2]
		local delete_set = redis.call('smembers', delete_set_key)
		local targets = redis.call('zrange', zset_key, score1, score2, 'byscore')
		redis.call('zremrangebyscore', zset_key, score1, score2)
		local reply = {}
		reply[1] = delete_set
		for i, v in ipairs(targets) do
			reply[#reply+1] = v
		end
		return reply
	`
)
