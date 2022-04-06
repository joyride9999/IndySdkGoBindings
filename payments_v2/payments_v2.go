/*
// ******************************************************************
// Purpose: Wrapper to call libindy, imports functions from indy_payments.h
// Author:  adrian.toader@siemens.com
// Notes:
// Copyright (c): Siemens SRL
// This work is licensed under the terms of the Apache License Version 2.0.  See
// the LICENSE.txt file in the top-level directory.
// ******************************************************************
*/

package payments_v2


/*
#cgo CFLAGS: -I ../include
#cgo LDFLAGS: -L${SRCDIR}/../lib -lindy
#include <indy_core.h>
typedef void (*cb_buildGetPaymentSourcesWithFromRequest)(indy_handle_t, indy_error_t, char*, char*);
extern void buildGetPaymentSourcesWithFromRequestCB(indy_handle_t, indy_error_t, char*, char*);

typedef void (*cb_parseGetPaymentSourcesWithFromResponse)(indy_handle_t, indy_error_t, char*, indy_u64_t);
extern void parseGetPaymentSourcesWithFromResponseCB(indy_handle_t, indy_error_t, char*, indy_u64_t);
*/
import "C"

import (
	"errors"
	"indySDK/indyUtils"
	"unsafe"
)

//export buildGetPaymentSourcesWithFromRequestCB
func buildGetPaymentSourcesWithFromRequestCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, getSourcesTxnJs *C.char, paymentMethod *C.char) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil, Results: []interface{}{string(C.GoString(getSourcesTxnJs)), string(C.GoString(paymentMethod))}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// BuildGetPaymentSourcesWithFromRequest builds Indy request for getting sources list for payment address
func BuildGetPaymentSourcesWithFromRequest(wh int, submitterDID string, paymentAddress string, from int64) chan indyUtils.IndyResult {
	handle, future := indyUtils.NewFutureCommand()

	commandHandle := (C.indy_handle_t)(handle)

	/*
		Builds Indy request for getting sources list for payment address.

		:param wallet_handle: wallet handle (created by open_wallet).
		:param submitter_did: (Optional) DID of request sender
		:param payment_address: target payment address
		:param from: shift to the next slice of payment sources

		:return: get_sources_txn_json - Indy request for getting sources list for payment address
				payment_method - used payment method
	*/

	// Call indy_build_get_payment_sources_with_from_request
	res := C.indy_build_get_payment_sources_with_from_request(commandHandle,
		(C.indy_handle_t)(wh),
		C.CString(submitterDID),
		C.CString(paymentAddress),
		C.ulonglong(from),
		(C.cb_buildGetPaymentSourcesWithFromRequest)(unsafe.Pointer(C.buildGetPaymentSourcesWithFromRequestCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export parseGetPaymentSourcesWithFromResponseCB
func parseGetPaymentSourcesWithFromResponseCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, sourcesJs *C.char, next C.ulonglong) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil, Results: []interface{}{string(C.GoString(sourcesJs)), int64(next)}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// ParseGetPaymentSourcesWithFromResponse parses response for Indy request for getting sources list
func ParseGetPaymentSourcesWithFromResponse(paymentMethod string, respJs string) chan indyUtils.IndyResult {
	handle, future := indyUtils.NewFutureCommand()

	commandHandle := (C.indy_handle_t)(handle)

	/*
		Parses response for Indy request for getting sources list.

		:param wallet_handle: wallet handle (created by open_wallet).
		:param payment_method: payment method to use.
		:param resp_json: response for Indy request for getting sources list

		:return: next - pointer to the next slice of payment sources
				sources_json - parsed (payment method and node version agnostic) sources info as json:
				  [{
				     source: <str>, // source input
				     paymentAddress: <str>, //payment address for this source
				     amount: <int>, // amount
				     extra: <str>, // optional data from payment transaction
				  }]
	*/

	// Call indy_parse_get_payment_sources_with_from_response
	res := C.indy_parse_get_payment_sources_with_from_response(commandHandle,
		C.CString(paymentMethod),
		C.CString(respJs),
		(C.cb_parseGetPaymentSourcesWithFromResponse)(unsafe.Pointer(C.parseGetPaymentSourcesWithFromResponseCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}
