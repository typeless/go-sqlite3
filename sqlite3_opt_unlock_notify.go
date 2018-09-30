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
	"unsafe"
)

type unlockNotification struct {
	notify chan struct{}
	lock   sync.Mutex
}

//export unlock_notify_callback
func unlock_notify_callback(pargv unsafe.Pointer, argc C.int) {
	argv := *(*uintptr)(pargv)
	v := (*[1 << 30]uintptr)(unsafe.Pointer(argv))
	for i := 0; i < int(argc); i++ {
		un := lookupHandle(v[i]).(unlockNotification)
		un.notify <- struct{}{}
	}
}

var notifyMutex sync.Mutex

//export unlock_notify_wait
func unlock_notify_wait(db *C.sqlite3) C.int {
	var un unlockNotification
	un.notify = make(chan struct{})
	defer close(un.notify)

	argv := [1]uintptr{newHandle(nil, un)}
	if rv := C.sqlite3_unlock_notify(db, (*[0]byte)(C._unlock_notify_callback), unsafe.Pointer(&argv)); rv != C.SQLITE_OK {
		return rv
	}
	<-un.notify
	return C.SQLITE_OK
}
