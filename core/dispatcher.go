// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

// IDispatcher is the interface for event dispatchers.
type IDispatcher[T any] interface {
	Subscribe(cb Callback[T]) SubscribeID
	Unsubscribe(id SubscribeID)
	Dispatch(ev T) int
}

// SubscribeID identifies a specific subscription to a dispatcher.
type SubscribeID int32

// SubscribeIDNone is the default SubscribeID, meaning there is no subscription.
const SubscribeIDNone = SubscribeID(0)

// Callback is a function that will be called when an event is dispatched.
type Callback[T any] func(T) bool

// subscription holds the SubscribeID and Callback for a subscriber.
type subscription[T any] struct {
	id SubscribeID
	cb Callback[T]
}

var _ IDispatcher[int] = &Dispatcher[int]{}

// Dispatcher is an event dispatcher.
type Dispatcher[T any] struct {
	subs   []subscription[T]
	nextId SubscribeID
}

// Subscribe adds the callback to the list invoked when an event is dispatched.
// Returns a SubscribeID that can be used to Unsubscribe later.
func (e *Dispatcher[T]) Subscribe(cb Callback[T]) SubscribeID {
	e.nextId++
	e.subs = append(e.subs, subscription[T]{id: e.nextId, cb: cb})
	return e.nextId
}

// Unsubscribe removes the callback with the given ID from the list invoked when an event is dispatched.
func (e *Dispatcher[T]) Unsubscribe(id SubscribeID) {
	for i := 0; i < len(e.subs); i++ {
		if e.subs[i].id == id {
			if i == len(e.subs)-1 {
				e.subs[i].cb = nil
			} else {
				e.subs[i] = e.subs[len(e.subs)-1]
				e.subs[len(e.subs)-1].cb = nil
			}
			e.subs = e.subs[:len(e.subs)-1]
			i--
		}
	}
}

// Dispatch invokes all subscribed callbacks with the given event.
// Returns the number of subscribers that consumed the event.
func (e *Dispatcher[T]) Dispatch(ev T) int {
	count := 0
	for _, sub := range e.subs {
		if sub.cb(ev) {
			count++
		}
	}
	return count
}
