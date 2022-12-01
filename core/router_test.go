// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package core

import (
	"reflect"
	"testing"
	"time"
)

type mockWriter struct {
	msgs []*Message
}

func (w *mockWriter) Start() error { return nil }
func (w *mockWriter) Stop() error  { return nil }

func (w *mockWriter) ResolveMessage(msg *Message) bool {
	w.msgs = append(w.msgs, msg)
	return true
}

func TestRouter(t *testing.T) {
	tLog := NewLog()
	router := NewRouter(tLog)

	chain0 := &mockWriter{msgs: *new([]*Message)}
	router.Listen(RSymbol("0"), chain0)

	chain1 := &mockWriter{msgs: *new([]*Message)}
	router.Listen(RSymbol("1"), chain1)

	msg0To1 := &Message{
		Source:      RSymbol("0"),
		Destination: RSymbol("1"),
	}

	msg1To0 := &Message{
		Source:      RSymbol("1"),
		Destination: RSymbol("0"),
	}

	err := router.Send(msg0To1)
	if err != nil {
		t.Fatal(err)
	}
	err = router.Send(msg1To0)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second)

	if !reflect.DeepEqual(chain0.msgs[0], msg1To0) {
		t.Error("Unexpected message")
	}

	if !reflect.DeepEqual(chain1.msgs[0], msg0To1) {
		t.Error("Unexpected message")
	}
}
