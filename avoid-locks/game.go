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

package avoid_locks

import "fmt"

type Player struct {
	Name string
}

type GameState struct {
	//lock    sync.RWMutex
	players []*Player

	msgCh chan any
}

func (g *GameState) Receive(msg any) {
	g.msgCh <- msg
}

func (g *GameState) loop() {
	for msg := range g.msgCh {
		// Each time we receive a message
		g.handleMessage(msg)
	}
}

func (g *GameState) handleMessage(message any) {
	switch msg := message.(type) {
	case *Player:
		g.addPlayer(msg)
	default:
		panic("invalid message received")
	}
}

func (g *GameState) addPlayer(p *Player) {
	//g.lock.Lock()
	g.players = append(g.players, p)
	//g.lock.Unlock()

	fmt.Println("add player:", p.Name)
}

func NewGameState() *GameState {
	g := &GameState{
		players: []*Player{},
		msgCh:   make(chan any, 10), // some size of the mailbox
	}
	go g.loop()
	return g
}
