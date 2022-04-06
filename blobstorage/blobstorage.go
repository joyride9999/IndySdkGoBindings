/*
// ******************************************************************
// Purpose: Wrapper to call libindy, imports functions from indy_blob_storage.h
// Author:  alexandru.leonte@siemens.com
// Notes:
// Copyright (c): Siemens SRL
// This work is licensed under the terms of the Apache License Version 2.0.  See
// the LICENSE.txt file in the top-level directory.
// ******************************************************************
*/

package blobstorage

/*
   #cgo CFLAGS: -I ../include
   #cgo LDFLAGS: -L${SRCDIR}/../lib -lindy
   #include <indy_core.h>

   typedef void (*cb_open_blob_storage)(indy_handle_t, indy_error_t, indy_handle_t);
   extern void openBlobStorageReaderCB(indy_handle_t, indy_error_t, indy_handle_t);
   extern void openBlobStorageWriterCB(indy_handle_t, indy_error_t, indy_handle_t);
*/
import "C"
import (
	"errors"
	"github.com/joyride9999/IndySdkGoBindings/indyUtils"
	"unsafe"
)

//export openBlobStorageReaderCB
func openBlobStorageReaderCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, blobHandle C.indy_handle_t) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle),
			indyUtils.IndyResult{Error: nil,
				Results: []interface{}{
					int(blobHandle),
				}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

//export openBlobStorageWriterCB
func openBlobStorageWriterCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, blobHandle C.indy_handle_t) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle),
			indyUtils.IndyResult{Error: nil,
				Results: []interface{}{
					int(blobHandle),
				}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// OpenBlobStorageReader Open a blob storage reader
func OpenBlobStorageReader(blobStorageType string, configJson string) chan indyUtils.IndyResult {

	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)
	/*

	 */

	// Call indy_open_blob_storage_reader
	res := C.indy_open_blob_storage_reader(commandHandle,
		C.CString(blobStorageType),
		C.CString(configJson),
		(C.cb_open_blob_storage)(unsafe.Pointer(C.openBlobStorageReaderCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() {
			indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)})
		}()
		return future
	}

	return future
}

// OpenBlobStorageWriter Open a blob storage writer
func OpenBlobStorageWriter(blobStorageType string, configJson string) chan indyUtils.IndyResult {

	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)
	/*

	 */

	// Call indy_open_blob_storage_writer
	res := C.indy_open_blob_storage_writer(commandHandle,
		C.CString(blobStorageType),
		C.CString(configJson),
		(C.cb_open_blob_storage)(unsafe.Pointer(C.openBlobStorageWriterCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() {
			indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)})
		}()
		return future
	}

	return future
}
