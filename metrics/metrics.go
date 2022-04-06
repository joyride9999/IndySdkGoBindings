/*
// ******************************************************************
// Purpose: Wrapper to call libindy, imports functions from indy_metrics.h
// Author:  adrian.toader@siemens.com
// Notes:
// Copyright (c): Siemens SRL
// This work is licensed under the terms of the Apache License Version 2.0.  See
// the LICENSE.txt file in the top-level directory.
// ******************************************************************
*/

package metrics

/*
#cgo CFLAGS: -I ../include
#cgo LDFLAGS: -L${SRCDIR}/../lib -lindy
#include <indy_core.h>
typedef void (*cb_collect)(indy_handle_t, indy_error_t, char*);

extern void collectCB(indy_handle_t, indy_error_t, char*);
*/
import "C"

import (
	"errors"
	"github.com/joyride9999/IndySdkGoBindings/indyUtils"
)

//export collectCB
func collectCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, js *C.char) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil, Results: []interface{}{string(C.GoString(js))}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// Collect collect metrics
func Collect() chan indyUtils.IndyResult {
	return nil
	//handle, future := indyUtils.NewFutureCommand()
	//
	//commandHandle := (C.indy_handle_t)(handle)
	//
	///*
	//	Collect metrics.
	//
	//	:return: Map in the JSON format. Where keys are names of metrics.
	//
	//*/
	//
	//// Call indy_collect_metrics
	//res := C.indy_collect_metrics(commandHandle,
	//	(C.cb_collect)(unsafe.Pointer(C.collectCB)))
	//if res != 0 {
	//	errMsg := indyUtils.GetIndyError(int(res))
	//	go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
	//	return future
	//}
	//
	//return future
}
