/*
// ******************************************************************
// Purpose: Wrapper to call libindy, imports functions from indy_pairwise.h
// Author:  adrian.toader@siemens.com
// Notes:
// Copyright (c): Siemens SRL
// This work is licensed under the terms of the Apache License Version 2.0.  See
// the LICENSE.txt file in the top-level directory.
// ******************************************************************
*/

package pairwise

/*
#cgo CFLAGS: -I ../include
#cgo LDFLAGS: -L${SRCDIR}/../lib -lindy
#include <indy_core.h>
typedef void (*cb_isPairwiseExists)(indy_handle_t, indy_error_t, indy_bool_t);
extern void isPairwiseExistsCB(indy_handle_t, indy_error_t, indy_bool_t);

typedef void (*cb_createPairwise)(indy_handle_t, indy_error_t);
extern void createPairwiseCB(indy_handle_t, indy_error_t);

typedef void (*cb_listPairwise)(indy_handle_t, indy_error_t, char*);
extern void listPairwiseCB(indy_handle_t, indy_error_t, char*);

typedef void (*cb_getPairwise)(indy_handle_t, indy_error_t, char*);
extern void getPairwiseCB(indy_handle_t, indy_error_t, char*);

typedef void (*cb_setPairwise)(indy_handle_t, indy_error_t);
extern void setPairwiseCB(indy_handle_t, indy_error_t);
*/
import "C"
import (
	"errors"
	"indySDK/indyUtils"
	"unsafe"
)

//export isPairwiseExistsCB
func isPairwiseExistsCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, exists C.indy_bool_t) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil, Results: []interface{}{bool(exists)}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// IsPairwiseExists Check if pairwise is exists
func IsPairwiseExists(wh int, theirDID string) chan indyUtils.IndyResult {
	handle, future := indyUtils.NewFutureCommand()

	commandHandle := (C.indy_handle_t)(handle)

	/*
		Check if pairwise is exists.

		:param wallet_handle: wallet handle (created by open_wallet).
		:param their_did: encrypted DID

		:return: Error code
	*/

	// Call indy_is_pairwise_exists
	res := C.indy_is_pairwise_exists(commandHandle,
		(C.indy_handle_t)(wh),
		C.CString(theirDID),
		(C.cb_isPairwiseExists)(unsafe.Pointer(C.isPairwiseExistsCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export createPairwiseCB
func createPairwiseCB(commandHandle C.indy_handle_t, indyError C.indy_error_t) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// CreatePairwise creates pairwise
func CreatePairwise(wh int, theirDID string, myDID string, meta string) chan indyUtils.IndyResult {
	handle, future := indyUtils.NewFutureCommand()

	commandHandle := (C.indy_handle_t)(handle)

	/*
		Creates pairwise.

		:param wallet_handle: wallet handle (created by open_wallet).
		:param their_did: encrypted DID
		:param my_did: encrypted DID
		:param metadata Optional: extra information for pairwise

		:return: Error code
	*/

	// Call indy_create_pairwise
	res := C.indy_create_pairwise(commandHandle,
		(C.indy_handle_t)(wh),
		C.CString(theirDID),
		C.CString(myDID),
		C.CString(meta),
		(C.cb_createPairwise)(unsafe.Pointer(C.createPairwiseCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export listPairwiseCB
func listPairwiseCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, listPairwise *C.char) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil, Results: []interface{}{string(C.GoString(listPairwise))}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// ListPairwise Get list of saved pairwise.
func ListPairwise(wh int) chan indyUtils.IndyResult {
	handle, future := indyUtils.NewFutureCommand()

	commandHandle := (C.indy_handle_t)(handle)

	/*
		Get list of saved pairwise.

		:param wallet_handle: wallet handle (created by open_wallet).

		:return: pairwise_list: list of saved pairwise.
	*/

	// Call indy_list_pairwise
	res := C.indy_list_pairwise(commandHandle,
		(C.indy_handle_t)(wh),
		(C.cb_listPairwise)(unsafe.Pointer(C.listPairwiseCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export getPairwiseCB
func getPairwiseCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, pairwiseInfo *C.char) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil, Results: []interface{}{string(C.GoString(pairwiseInfo))}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// GetPairwise gets pairwise information for specific their_did
func GetPairwise(wh int, theirDID string) chan indyUtils.IndyResult {
	handle, future := indyUtils.NewFutureCommand()

	commandHandle := (C.indy_handle_t)(handle)

	/*
		Gets pairwise information for specific their_did.

		:param wallet_handle: wallet handle (created by open_wallet).
		:param their_did: encrypted DID

		:return: Error code
	*/

	// Call indy_get_pairwise
	res := C.indy_get_pairwise(commandHandle,
		(C.indy_handle_t)(wh),
		C.CString(theirDID),
		(C.cb_getPairwise)(unsafe.Pointer(C.getPairwiseCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export setPairwiseCB
func setPairwiseCB(commandHandle C.indy_handle_t, indyError C.indy_error_t) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// SetPairwiseMetadata save some data in the Wallet for pairwise associated with Did
func SetPairwiseMetadata(wh int, theirDID string, meta string) chan indyUtils.IndyResult {
	handle, future := indyUtils.NewFutureCommand()

	commandHandle := (C.indy_handle_t)(handle)

	/*
		Save some data in the Wallet for pairwise associated with Did.

		:param wallet_handle: wallet handle (created by open_wallet).
		:param their_did: encrypted DID
		:param metadata: some extra information for pairwise

		:return: Error code
	*/

	// Call indy_set_pairwise_metadata
	res := C.indy_set_pairwise_metadata(commandHandle,
		(C.indy_handle_t)(wh),
		C.CString(theirDID),
		C.CString(meta),
		(C.cb_setPairwise)(unsafe.Pointer(C.setPairwiseCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}
