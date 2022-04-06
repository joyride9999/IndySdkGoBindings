/*
// ******************************************************************
// Purpose: provides callbacks functionality from indy to go
// Author:  alexandru.leonte@siemens.com
// Notes:
// Copyright (c): Siemens SRL
// This work is licensed under the terms of the Apache License Version 2.0.  See
// the LICENSE.txt file in the top-level directory.
// ******************************************************************
*/

package indyUtils

/*
#cgo CFLAGS: -I ../include
#cgo LDFLAGS: -L${SRCDIR}/../lib -lindy
#include <indy_core.h>
*/
import "C"

import (
	cmap "github.com/orcaman/concurrent-map"
	"strconv"
	"sync"
)

// IndyResult represents callback result from C-call to libindy
type IndyResult struct {
	Error   error
	Results []interface{}
}

// futures Concurrent channel map
var futures = cmap.New()

// Counter Concurrent counter
type Counter struct {
	sync.Mutex
	Count int32
}

func (c *Counter) Increment() {
	c.Count++
}

func (c *Counter) Get() (int32, string) {
	c.Lock()
	defer c.Unlock()
	c.Increment()
	return c.Count, strconv.Itoa(int(c.Count))
}

var count Counter

// NewFutureCommand creates a new future command
func NewFutureCommand() (C.indy_handle_t, chan IndyResult) {
	commandHandle, futuresKey := count.Get()
	//fmt.Println("command handle %d", commandHandle)
	future := make(chan IndyResult)
	// Save to the map our handle
	futures.Set(futuresKey, future)
	return (C.indy_handle_t)(commandHandle), future
}

// RemoveFuture removes a future from the futures map
func RemoveFuture(commandHandle int, result IndyResult) chan IndyResult {
	//fmt.Println("remove command handle %d", commandHandle)
	futuresKey := strconv.Itoa(commandHandle)
	future, _ := futures.Get(futuresKey)

	future.(chan IndyResult) <- result
	futures.Remove(futuresKey)
	return future.(chan IndyResult)
}
