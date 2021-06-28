// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package core

import (
	"fmt"
	"sync"

	log "github.com/ChainSafe/log15"
)

const msgLimit = 4096

// Writer consumes a message and makes the requried on-chain interactions.
type Writer interface {
	ResolveMessage(msg *Message) bool
}

// Router forwards messages from their source to their destination
type Router struct {
	registry  map[RSymbol]Writer
	lock      *sync.RWMutex
	log       log.Logger
	msgChan   chan *Message
	maticChan chan *Message
	stop      chan int
}

func NewRouter(log log.Logger) *Router {
	return &Router{
		registry:  make(map[RSymbol]Writer),
		lock:      &sync.RWMutex{},
		log:       log,
		msgChan:   make(chan *Message, msgLimit),
		maticChan: make(chan *Message, msgLimit),
		stop:      make(chan int),
	}
}

// Send passes a message to the destination Writer if it exists
func (r *Router) Send(msg *Message) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if msg.Reason != NewEra {
		r.log.Trace("Routing message", "source", msg.Source, "dest", msg.Destination, "Reason", msg.Reason)
	}

	w := r.registry[msg.Destination]
	if w == nil {
		return fmt.Errorf("unknown destination symbol: %s", msg.Destination)
	}

	if msg.Destination == RFIS {
		r.QueueMsg(msg)
	} else if msg.Destination == RMATIC {
		r.QueueMaticMsg(msg)
	} else {
		go w.ResolveMessage(msg)
	}
	return nil
}

// Listen registers a Writer with a ChainId which Router.Send can then use to propagate messages
func (r *Router) Listen(symbol RSymbol, w Writer) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.log.Debug("Registering new chain in router", "symbol", symbol)
	r.registry[symbol] = w
}

func (r *Router) QueueMsg(m *Message) {
	r.msgChan <- m
}

func (r *Router) QueueMaticMsg(m *Message) {
	r.maticChan <- m
}

func (r *Router) MsgHandler() {
	r.lock.Lock()
	rfis := r.registry[RFIS]
	if rfis == nil {
		panic("RFIS writer not exist")
	}
	rmatic := r.registry[RMATIC]
	if rmatic == nil {
		panic("RMATIC writer not exist")
	}
	r.lock.Unlock()

out:
	for {
		select {
		case <-r.stop:
			r.log.Info("RFIS msgHandler stop")
			break out
		case msg := <-r.msgChan:
			rfis.ResolveMessage(msg)
		case msg := <-r.maticChan:
			rmatic.ResolveMessage(msg)
		}
	}
}

func (r *Router) StopMsgHandler() {
	close(r.stop)
	close(r.msgChan)
	close(r.maticChan)
}
