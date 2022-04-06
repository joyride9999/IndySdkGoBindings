/*
// ******************************************************************
// Purpose: Wrapper to call libindy, imports functions from indy_logger.h
// Author:  alexandru.leonte@siemens.com
// Notes:
// Copyright (c): Siemens SRL
// This work is licensed under the terms of the Apache License Version 2.0.  See
// the LICENSE.txt file in the top-level directory.
// ******************************************************************
*/

package logger

/*
#cgo CFLAGS: -I ../include
#cgo LDFLAGS: -L${SRCDIR}/../lib -lindy
#include <indy_core.h>

typedef void (*cb_indySetLogger)(void*, indy_u32_t, char*, char*, char*, char*, indy_u32_t);
extern void indySetLoggerCB(void*, indy_u32_t, char*, char*, char*, char*, indy_u32_t);

typedef indy_bool_t (*cb_enableLog)(void*, indy_u32_t, char*);
extern indy_bool_t enableLogCB(void*, indy_u32_t, char*);
*/
import "C"

import (
	"fmt"
	"time"
	"unsafe"
)

//export enableLogCB
func enableLogCB(context *C.void, level C.indy_u32_t, target *C.char) C.indy_bool_t {
	return true
}

//export indySetLoggerCB
func indySetLoggerCB(context *C.void, level C.indy_u32_t, target *C.char, message *C.char, modulePath *C.char, file *C.char, line C.indy_u32_t) {

	now := time.Now()
	fmt.Printf("LOG - %s\t%s:%d | %s %s \n", now.Format(time.RFC850), string(C.GoString(file)), uint32(line), string(C.GoString(target)), string(C.GoString(message)))

}

// IndySetLogger sets the logger
func IndySetLogger() {

	res := C.indy_set_logger(nil,
		(C.cb_enableLog)(unsafe.Pointer(C.enableLogCB)),
		(C.cb_indySetLogger)(unsafe.Pointer(C.indySetLoggerCB)),
		nil)

	res = res

	return
}
