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

package payments

/*
#cgo CFLAGS: -I ../include
#cgo LDFLAGS: -L${SRCDIR}/../lib -lindy
#include <indy_core.h>
typedef void (*cb_createPaymentAddress)(indy_handle_t, indy_error_t, char*);
typedef void (*cb_listPaymentAddress)(indy_handle_t, indy_error_t, char*);
typedef void (*cb_addRequestFees)(indy_handle_t, indy_error_t, char*, char*);
typedef void (*cb_parseResponseWithFees)(indy_handle_t, indy_error_t, char*);
typedef void (*cb_buildPaymentReq)(indy_handle_t, indy_error_t, char*);
typedef void (*cb_parsePaymentResponse)(indy_handle_t, indy_error_t, char*);
typedef void (*cb_buildMintReq)(indy_handle_t, indy_error_t, char*, char*);
typedef void (*cb_buildSetTxnFeesReq)(indy_handle_t, indy_error_t, char*);
typedef void (*cb_buildGetTxnFeesReq)(indy_handle_t, indy_error_t, char*);
typedef void (*cb_parseGetTxnFeesResponse)(indy_handle_t, indy_error_t, char*);
typedef void (*cb_buildVerifyPaymentReq)(indy_handle_t, indy_error_t, char*, char*);
typedef void (*cb_parseVerifyPaymentResponse)(indy_handle_t, indy_error_t, char*);
typedef void (*cb_signWithAddress)(indy_handle_t, indy_error_t, indy_u8_t, indy_u32_t);
typedef void (*cb_verifyWithAddress)(indy_handle_t, indy_error_t, bool);

extern void createPaymentAddressCB(indy_handle_t, indy_error_t, char*);
extern void listPaymentAddressCB(indy_handle_t, indy_error_t, char*);
extern void addRequestFeesCB(indy_handle_t, indy_error_t, char*, char*);
extern void parseResponseWithFeesCB(indy_handle_t, indy_error_t, char*);
extern void buildPaymentReqCB(indy_handle_t, indy_error_t, char*);
extern void parsePaymentResponseCB(indy_handle_t, indy_error_t, char*);
extern void buildMintReqCB(indy_handle_t, indy_error_t, char*, char*);
extern void buildSetTxnFeesReqCB(indy_handle_t, indy_error_t, char*);
extern void buildGetTxnFeesReqCB(indy_handle_t, indy_error_t, char*);
extern void parseGetTxnFeesResponseCB(indy_handle_t, indy_error_t, char*);
extern void buildVerifyPaymentReqCB(indy_handle_t, indy_error_t, char*, char*);
extern void parseVerifyPaymentResponseCB(indy_handle_t, indy_error_t, char*, char*);
extern void signWithAddressCB(indy_handle_t, indy_error_t, indy_u8_t, indy_u32_t);
extern void verifyWithAddressCB(indy_handle_t, indy_error_t, bool);
*/
import "C"
import (
	"github.com/joyride9999/IndySdkGoBindings/indyUtils"
	"encoding/json"
	"errors"
	"unsafe"
)

//export createPaymentAddressCB
func createPaymentAddressCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, paymentAddress *C.char) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil, Results: []interface{}{string(C.GoString(paymentAddress))}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// CreatePaymentAddress creates the payment address for specified payment method
func CreatePaymentAddress (wh int, paymentMethod string, config Config) chan indyUtils.IndyResult {
	handle, future := indyUtils.NewFutureCommand()

	jsonConfig, err := json.Marshal(config)
	if err != nil {
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: err}) }()
		return future
	}

	configString := string(jsonConfig)
	commandHandle := (C.indy_handle_t)(handle)

	/*
		Create the payment address for specified payment method.
		This method generates private part of payment address and stores it in a secure place. Ideally it should be
		secret in libindy wallet (see crypto module).

		Note that payment method should be able to resolve this secret by fully resolvable payment address format.

		:param wallet_handle: wallet handle where to save new address
		:param payment_method: payment method to use (for example, 'sov')
		:param config: payment address config as json:
		  {
		    seed: <str>, // allows deterministic creation of payment address
		  }

		:return: payment_address - public identifier of payment address in fully resolvable payment address format
	*/

	// Call indy_create_payment_address
	res := C.indy_create_payment_address(commandHandle,
		(C.indy_handle_t)(wh),
		C.CString(paymentMethod),
		C.CString(configString),
		(C.cb_createPaymentAddress)(unsafe.Pointer(C.createPaymentAddressCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export listPaymentAddressCB
func listPaymentAddressCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, paymentAddresses *C.char) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil, Results: []interface{}{string(C.GoString(paymentAddresses))}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// ListPaymentAddress lists all payment addresses that are stored in the wallet
func ListPaymentAddress (wh int) chan indyUtils.IndyResult {
	handle, future := indyUtils.NewFutureCommand()

	commandHandle := (C.indy_handle_t)(handle)

	/*
		Lists all payment addresses that are stored in the wallet.

		:param wallet_handle: wallet handle (created by open_wallet).

		:return: payment_addresses_json - json array of string with json addresses
	*/

	// Call indy_list_payment_addresses
	res := C.indy_list_payment_addresses(commandHandle,
		(C.indy_handle_t)(wh),
		(C.cb_listPaymentAddress)(unsafe.Pointer(C.listPaymentAddressCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export addRequestFeesCB
func addRequestFeesCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, reqWithFees *C.char, paymentMethod *C.char) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil, Results: []interface{}{string(C.GoString(reqWithFees)), string(C.GoString(paymentMethod))}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// AddRequestFees modifies Indy request by adding information how to pay fees for this transaction according to this payment method
func AddRequestFees (wh int, submitterDID string, req string, inputs string, outputs string, extra string) chan indyUtils.IndyResult {
	handle, future := indyUtils.NewFutureCommand()

	commandHandle := (C.indy_handle_t)(handle)

	/*
		Modifies Indy request by adding information how to pay fees for this transaction according to this payment method.
		This method consumes set of inputs and outputs. The difference between inputs balance
		and outputs balance is the fee for this transaction.

		Not that this method also produces correct fee signatures.

		Format of inputs is specific for payment method. Usually it should reference payment transaction
		with at least one output that corresponds to payment address that user owns.

		:param wallet_handle: wallet handle (created by open_wallet).
		:param submitter_did: (Optional) DID of request sender
		:param req_json: initial transaction request as json
		:param inputs_json: The list of payment sources as json array:
		  ["source1", ...]
		    - each input should reference paymentAddress
		    - this param will be used to determine payment_method
		:param outputs_json: The list of outputs as json array:
		  [{
		    recipient: <str>, // payment address of recipient
		    amount: <int>, // amount
		  }]
		:param extra: // optional information for payment operation

		:return: req_with_fees_json - modified Indy request with added fees info
				payment_method - used payment method
	*/

	// Call indy_add_request_fees
	res := C.indy_add_request_fees(commandHandle,
		(C.indy_handle_t)(wh),
		C.CString(submitterDID),
		C.CString(req),
		C.CString(inputs),
		C.CString(outputs),
		C.CString(extra),
		(C.cb_addRequestFees)(unsafe.Pointer(C.addRequestFeesCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

////export parseResponseWithFeesCB
//func parseResponseWithFeesCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, receipts *C.char) {
//	if indyError == 0 {
//		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil, Results: []interface{}{string(receipts)}})
//	} else {
//		errMsg := indyUtils.GetIndyError(int(indyError))
//		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
//	}
//}
//
//// ParseResponseWithFees parses response for Indy request with fees
//func ParseResponseWithFees (wh int, paymentMethod string, resp string) chan indyUtils.IndyResult {
//	handle, future := indyUtils.NewFutureCommand()
//
//	commandHandle := (C.indy_handle_t)(handle)
//
//	/*
//		Parses response for Indy request with fees.
//
//		:param wallet_handle: wallet handle (created by open_wallet).
//		:param payment_method: payment method to use
//		:param resp_json: response for Indy request with fees
//
//		:return: receipts_json - parsed (payment method and node version agnostic) receipts info as json:
//				  [{
//				     receipt: <str>, // receipt that can be used for payment referencing and verification
//				     recipient: <str>, //payment address of recipient
//				     amount: <int>, // amount
//				     extra: <str>, // optional data from payment transaction
//				  }]
//	*/
//
//	// Call indy_parse_response_with_fees
//	res := C.indy_parse_response_with_fees(commandHandle,
//		(C.indy_handle_t)(wh),
//		C.CString(paymentMethod),
//		C.CString(resp),
//		(C.cb_parseResponseWithFees)(unsafe.Pointer(C.parseResponseWithFeesCB)))
//	if res != 0 {
//		errMsg := indyUtils.GetIndyError(int(res))
//		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
//		return future
//	}
//
//	return future
//}
//
////export buildPaymentReqCB
//func buildPaymentReqCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, paymentReq *C.char, paymentMethod *C.char) {
//	if indyError == 0 {
//		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil, Results: []interface{}{string(paymentReq), string(paymentMethod)}})
//	} else {
//		errMsg := indyUtils.GetIndyError(int(indyError))
//		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
//	}
//}
//
//// BuildPaymentReq builds Indy request for doing payment according to this payment method
//func BuildPaymentReq (wh int, submitterDID string, req string, inputs string, outputs string, extra string) chan indyUtils.IndyResult {
//	handle, future := indyUtils.NewFutureCommand()
//
//	commandHandle := (C.indy_handle_t)(handle)
//
//	/*
//		Parses response for Indy request with fees.
//
//		:param wallet_handle: wallet handle (created by open_wallet).
//		:param submitter_did: (Optional) DID of request sender
//		:param req_json: initial transaction request as json
//		:param inputs_json: The list of payment sources as json array:
//		  ["source1", ...]
//		    - each input should reference paymentAddress
//		    - this param will be used to determine payment_method
//		:param outputs_json: The list of outputs as json array:
//		  [{
//		    recipient: <str>, // payment address of recipient
//		    amount: <int>, // amount
//		  }]
//		:param extra: // optional information for payment operation
//
//		:return: payment_req_json - Indy request for doing payment
//				payment_method - used payment method
//	*/
//
//	// Call indy_build_get_payment_sources_request
//	res := C.indy_build_payment_req(commandHandle,
//		(C.indy_handle_t)(wh),
//		C.CString(submitterDID),
//		C.CString(req),
//		C.CString(inputs),
//		C.CString(outputs),
//		C.CString(extra),
//		(C.cb_buildPaymentReq)(unsafe.Pointer(C.buildPaymentReqCB)))
//	if res != 0 {
//		errMsg := indyUtils.GetIndyError(int(res))
//		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
//		return future
//	}
//
//	return future
//}
//
////export parsePaymentResponseCB
//func parsePaymentResponseCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, receipts *C.char) {
//	if indyError == 0 {
//		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil, Results: []interface{}{string(receipts)}})
//	} else {
//		errMsg := indyUtils.GetIndyError(int(indyError))
//		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
//	}
//}
//
//// ParsePaymentResponse parses response for Indy request for payment txn.
//func ParsePaymentResponse (wh int, paymentMethod string, resp string) chan indyUtils.IndyResult {
//	handle, future := indyUtils.NewFutureCommand()
//
//	commandHandle := (C.indy_handle_t)(handle)
//
//	/*
//		Parses response for Indy request for payment txn.
//
//		:param wallet_handle: wallet handle (created by open_wallet).
//		:param payment_method: payment method to use
//		:param resp_json: response for Indy request with fees
//
//		:return: receipts_json - parsed (payment method and node version agnostic) receipts info as json:
//				  [{
//				     receipt: <str>, // receipt that can be used for payment referencing and verification
//				     recipient: <str>, //payment address of recipient
//				     amount: <int>, // amount
//				     extra: <str>, // optional data from payment transaction
//				  }]
//	*/
//
//	// Call indy_parse_payment_response
//	res := C.indy_parse_payment_response(commandHandle,
//		(C.indy_handle_t)(wh),
//		C.CString(paymentMethod),
//		C.CString(resp),
//		(C.cb_parseResponseWithFees)(unsafe.Pointer(C.parseResponseWithFeesCB)))
//	if res != 0 {
//		errMsg := indyUtils.GetIndyError(int(res))
//		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
//		return future
//	}
//
//	return future
//}
//
////export buildMintReqCB
//func buildMintReqCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, mintReq *C.char, paymentMethod *C.char) {
//	if indyError == 0 {
//		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil, Results: []interface{}{string(mintReq), string(paymentMethod)}})
//	} else {
//		errMsg := indyUtils.GetIndyError(int(indyError))
//		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
//	}
//}
//
//// BuildMintReq builds Indy request for doing minting according to this payment method
//func BuildMintReq (wh int, submitterDID string, outputs string, extra string) chan indyUtils.IndyResult {
//	handle, future := indyUtils.NewFutureCommand()
//
//	commandHandle := (C.indy_handle_t)(handle)
//
//	/*
//		Builds Indy request for doing minting according to this payment method.
//
//		:param wallet_handle: wallet handle
//		:param submitter_did: (Optional) DID of request sender
//		:param outputs_json: The list of outputs as json array:
//		  [{
//		    recipient: <str>, // payment address of recipient
//		    amount: <int>, // amount
//		  }]
//		:param extra: // optional information for payment operation
//
//		:return: mint_req_json - Indy request for doing minting
//				payment_method - used payment method
//	*/
//
//	// Call indy_build_mint_req
//	res := C.indy_build_mint_req(commandHandle,
//		(C.indy_handle_t)(wh),
//		C.CString(submitterDID),
//		C.CString(outputs),
//		C.CString(extra),
//		(C.cb_buildMintReq)(unsafe.Pointer(C.buildMintReqCB)))
//	if res != 0 {
//		errMsg := indyUtils.GetIndyError(int(res))
//		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
//		return future
//	}
//
//	return future
//}
//
////export buildSetTxnFeesReqCB
//func buildSetTxnFeesReqCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, setTxnFees *C.char) {
//	if indyError == 0 {
//		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil, Results: []interface{}{string(setTxnFees)}})
//	} else {
//		errMsg := indyUtils.GetIndyError(int(indyError))
//		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
//	}
//}
//
//// BuildSetTxnFeesReq builds Indy request for setting fees for transactions in the ledger
//func BuildSetTxnFeesReq (wh int, submitterDID string, paymentMethod string, fees string) chan indyUtils.IndyResult {
//	handle, future := indyUtils.NewFutureCommand()
//
//	commandHandle := (C.indy_handle_t)(handle)
//
//	/*
//		Builds Indy request for setting fees for transactions in the ledger.
//
//		:param wallet_handle: wallet handle
//		:param submitter_did: (Optional) DID of request sender
//		:param payment_method: payment method to use
//		fees_json {
//		  txnType1: amount1,
//		  txnType2: amount2,
//		  .................
//		  txnTypeN: amountN,
//		}
//
//		:return: set_txn_fees_json - Indy request for setting fees for transactions in the ledger
//	*/
//
//	// Call indy_build_set_txn_fees_req
//	res := C.indy_build_set_txn_fees_req(commandHandle,
//		(C.indy_handle_t)(wh),
//		C.CString(submitterDID),
//		C.CString(paymentMethod),
//		C.CString(fees),
//		(C.cb_buildSetTxnFeesReq)(unsafe.Pointer(C.buildSetTxnFeesReqCB)))
//	if res != 0 {
//		errMsg := indyUtils.GetIndyError(int(res))
//		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
//		return future
//	}
//
//	return future
//}
//
////export buildGetTxnFeesReqCB
//func buildGetTxnFeesReqCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, getTxnFees *C.char) {
//	if indyError == 0 {
//		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil, Results: []interface{}{string(getTxnFees)}})
//	} else {
//		errMsg := indyUtils.GetIndyError(int(indyError))
//		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
//	}
//}
//
//// BuildGetTxnFeesReq builds Indy request for getting fees for transactions in the ledger
//func BuildGetTxnFeesReq (wh int, submitterDID string, paymentMethod string) chan indyUtils.IndyResult {
//	handle, future := indyUtils.NewFutureCommand()
//
//	commandHandle := (C.indy_handle_t)(handle)
//
//	/*
//		Builds Indy request for getting fees for transactions in the ledger.
//
//		:param wallet_handle: wallet handle
//		:param submitter_did: (Optional) DID of request sender
//		:param payment_method: payment method to use
//
//		:return: get_txn_fees_json - Indy request for getting fees for transactions in the ledger
//	*/
//
//	// Call indy_build_get_txn_fees_req
//	res := C.indy_build_get_txn_fees_req(commandHandle,
//		(C.indy_handle_t)(wh),
//		C.CString(submitterDID),
//		C.CString(paymentMethod),
//		(C.cb_buildGetTxnFeesReq)(unsafe.Pointer(C.buildGetTxnFeesReqCB)))
//	if res != 0 {
//		errMsg := indyUtils.GetIndyError(int(res))
//		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
//		return future
//	}
//
//	return future
//}
//
////export parseGetTxnFeesResponseCB
//func parseGetTxnFeesResponseCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, fees *C.char) {
//	if indyError == 0 {
//		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil, Results: []interface{}{string(fees)}})
//	} else {
//		errMsg := indyUtils.GetIndyError(int(indyError))
//		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
//	}
//}
//
//// ParseGetTxnFeesResponse parses response for Indy request for getting fees
//func ParseGetTxnFeesResponse (paymentMethod string, response string) chan indyUtils.IndyResult {
//	handle, future := indyUtils.NewFutureCommand()
//
//	commandHandle := (C.indy_handle_t)(handle)
//
//	/*
//		Parses response for Indy request for getting fees.
//
//		:param payment_method: payment method to use
//		:param resp_json: response for Indy request for getting fees
//
//		:return: fees_json {
//		  txnType1: amount1,
//		  txnType2: amount2,
//		  .................
//		  txnTypeN: amountN,
//		}
//	*/
//
//	// Call indy_parse_get_txn_fees_response
//	res := C.indy_parse_get_txn_fees_response(commandHandle,
//		C.CString(paymentMethod),
//		C.CString(response),
//		(C.cb_parseGetTxnFeesResponse)(unsafe.Pointer(C.parseGetTxnFeesResponseCB)))
//	if res != 0 {
//		errMsg := indyUtils.GetIndyError(int(res))
//		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
//		return future
//	}
//
//	return future
//}
//
////export buildVerifyPaymentReqCB
//func buildVerifyPaymentReqCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, verifyTxn *C.char, paymentMethod *C.char) {
//	if indyError == 0 {
//		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil, Results: []interface{}{string(verifyTxn), string(paymentMethod)}})
//	} else {
//		errMsg := indyUtils.GetIndyError(int(indyError))
//		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
//	}
//}
//
//// BuildVerifyPaymentReq builds Indy request for information to verify the payment receipt
//func BuildVerifyPaymentReq (wh int, submitterDID string, receipt string) chan indyUtils.IndyResult {
//	handle, future := indyUtils.NewFutureCommand()
//
//	commandHandle := (C.indy_handle_t)(handle)
//
//	/*
//		Builds Indy request for information to verify the payment receipt.
//
//		:param wallet_handle: wallet handle
//		:param submitter_did: (Optional) DID of request sender
//		:param receipt: payment receipt to verify
//
//		:return: verify_txn_json: Indy request for verification receipt
//				payment_method: used payment method
//	*/
//
//	// Call indy_build_verify_payment_req
//	res := C.indy_build_verify_payment_req(commandHandle,
//		(C.indy_handle_t)(wh),
//		C.CString(submitterDID),
//		C.CString(receipt),
//		(C.cb_buildVerifyPaymentReq)(unsafe.Pointer(C.buildVerifyPaymentReqCB)))
//	if res != 0 {
//		errMsg := indyUtils.GetIndyError(int(res))
//		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
//		return future
//	}
//
//	return future
//}
//
////export parseVerifyPaymentResponseCB
//func parseVerifyPaymentResponseCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, txn *C.char) {
//	if indyError == 0 {
//		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil, Results: []interface{}{string(txn)}})
//	} else {
//		errMsg := indyUtils.GetIndyError(int(indyError))
//		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
//	}
//}
//
//// ParseVerifyPaymentResponse parses Indy response with information to verify receipt
//func ParseVerifyPaymentResponse (submitterDID string, receipt string) chan indyUtils.IndyResult {
//	handle, future := indyUtils.NewFutureCommand()
//
//	commandHandle := (C.indy_handle_t)(handle)
//
//	/*
//		Parses Indy response with information to verify receipt.
//
//		:param payment_method: payment method to use
//		:param resp_json: response of the ledger for verify txn
//
//		:return: txn_json: {
//				    sources: [<str>, ]
//				    receipts: [ {
//				        recipient: <str>, // payment address of recipient
//				        receipt: <str>, // receipt that can be used for payment referencing and verification
//				        amount: <int>, // amount
//				    } ],
//				    extra: <str>, //optional data
//				}
//	*/
//
//	// Call indy_parse_verify_payment_response
//	res := C.indy_parse_verify_payment_response(commandHandle,
//		C.CString(submitterDID),
//		C.CString(receipt),
//		(C.cb_parseVerifyPaymentResponse)(unsafe.Pointer(C.parseVerifyPaymentResponseCB)))
//	if res != 0 {
//		errMsg := indyUtils.GetIndyError(int(res))
//		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
//		return future
//	}
//
//	return future
//}
//
////export signWithAddressCB
//func signWithAddressCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, signatureRaw *C.indy_u8_t, signatureLen C.indy_u32_t) {
//	if indyError == 0 {
//		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil, Results: []interface{}{uint8(signatureRaw), uint32(signatureLen)}})
//	} else {
//		errMsg := indyUtils.GetIndyError(int(indyError))
//		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
//	}
//}
//
//// SignWithAddress signs a message with a payment address
//func SignWithAddress (wh int, address string, messageRaw uint8, messageLen uint32) chan indyUtils.IndyResult {
//	handle, future := indyUtils.NewFutureCommand()
//
//	commandHandle := (C.indy_handle_t)(handle)
//
//	/*
//		Signs a message with a payment address.
//
//		:param wallet_handle: wallet handle
//		:param address: payment address of message signer. The key must be created by calling indy_create_address
//		:param message_raw: a pointer to first byte of message to be signed
//		:param message_len: a message length
//
//		:return: a signature string
//	*/
//
//	// Call indy_sign_with_address
//	res := C.indy_sign_with_address(commandHandle,
//		(C.indy_handle_t)(wh),
//		C.CString(address),
//		C.indy_u8_t(messageRaw),
//		C.indy_u32_t(messageLen),
//		(C.cb_parseVerifyPaymentResponse)(unsafe.Pointer(C.parseVerifyPaymentResponseCB)))
//	if res != 0 {
//		errMsg := indyUtils.GetIndyError(int(res))
//		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
//		return future
//	}
//
//	return future
//}
//
////export verifyWithAddressCB
//func verifyWithAddressCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, result C.indy_bool_t) {
//	if indyError == 0 {
//		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil, Results: []interface{}{bool(result)}})
//	} else {
//		errMsg := indyUtils.GetIndyError(int(indyError))
//		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
//	}
//}
//
//// VerifyWithAddress verify a signature with a payment address
//func VerifyWithAddress (address string, messageRaw uint8, messageLen uint32, signatureRaw uint8, signatureLen uint32) chan indyUtils.IndyResult {
//	handle, future := indyUtils.NewFutureCommand()
//
//	commandHandle := (C.indy_handle_t)(handle)
//
//	/*
//		Verify a signature with a payment address.
//
//		:param address: payment address of the message signer
//		:param message_raw: a pointer to first byte of message that has been signed
//		:param message_len: a message length
//		:param signature_raw: a pointer to first byte of signature to be verified
//		:param signature_len: a signature length
//
//		:return: valid: true - if signature is valid, false - otherwise
//	*/
//
//	// Call indy_verify_with_address
//	res := C.indy_verify_with_address(commandHandle,
//		C.CString(address),
//		C.indy_u8_t(messageRaw),
//		C.indy_u32_t(messageLen),
//		C.indy_u8_t(signatureRaw),
//		C.indy_u32_t(signatureLen),
//		(C.cb_parseVerifyPaymentResponse)(unsafe.Pointer(C.parseVerifyPaymentResponseCB)))
//	if res != 0 {
//		errMsg := indyUtils.GetIndyError(int(res))
//		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
//		return future
//	}
//
//	return future
//}

// EXPERIMENTAL:
// indy_prepare_payment_extra_with_acceptance_data
// indy_get_request_info

// Deprecated:
// indy_build_get_payment_sources_request
// indy_parse_get_payment_sources_response
