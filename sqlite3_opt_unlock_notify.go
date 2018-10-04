// Copyright (C) 2018 Yasuhiro Matsumoto <mattn.jp@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// +build cgo
// +build sqlite_unlock_notify

package sqlite3

/*
#cgo CFLAGS: -DSQLITE_ENABLE_UNLOCK_NOTIFY

#include <sqlite3-binding.h>

extern void _unlock_notify_callback(void *arg, int argc);
*/
import "C"
import (
	"sync"
)

var notify_unlock_chan_lock sync.Mutex
var notify_unlock_chan_seqnum uint
var notify_unlock_chan_table map[uint]chan struct{} = make(map[uint]chan struct{})

//export unlock_notify_chan_open
func unlock_notify_chan_open() uint {
	notify_unlock_chan_lock.Lock()
	c := make(chan struct{})
	h := notify_unlock_chan_seqnum
	notify_unlock_chan_table[h] = c
	notify_unlock_chan_seqnum++
	notify_unlock_chan_lock.Unlock()
	return h
}

//export unlock_notify_chan_close
func unlock_notify_chan_close(h uint) {
	notify_unlock_chan_lock.Lock()
	c, ok := notify_unlock_chan_table[h]
	if !ok {
		panic("non-existent notification entry")
	}
	notify_unlock_chan_lock.Unlock()
	close(c)
}

//export unlock_notify_chan_poll
func unlock_notify_chan_poll(h uint) {
	notify_unlock_chan_lock.Lock()
	c, ok := notify_unlock_chan_table[h]
	if !ok {
		panic("non-existent notification entry")
	}
	notify_unlock_chan_lock.Unlock()
	<-c
	delete(notify_unlock_chan_table, h)
}
