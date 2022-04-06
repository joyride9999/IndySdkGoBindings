/*
// ******************************************************************
// Purpose: exported public functions that handles payments functions
// from libindy
// Author:  adrian.toader@siemens.com
// Notes:
// Copyright (c): Siemens SRL
// This work is licensed under the terms of the Apache License Version 2.0.  See
// the LICENSE.txt file in the top-level directory.
// ******************************************************************
*/

package indySDK

import "github.com/joyride9999/IndySdkGoBindings/payments"

// CreatePaymentAddress creates the payment address for specified payment method
func CreatePaymentAddress(wh int, paymentMethod string, options payments.Config) (string, error) {
	channel := payments.CreatePaymentAddress(wh, paymentMethod, options)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// ListPaymentAddress lists all payment addresses that are stored in the wallet
func ListPaymentAddress(wh int) (string, error) {
	channel := payments.ListPaymentAddress(wh)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// AddRequestFees lists all payment addresses that are stored in the wallet
func AddRequestFees(wh int, submitterDID string, req string, inputs string, outputs string, extra string) (string, string, error) {
	channel := payments.AddRequestFees(wh, submitterDID, req, inputs, outputs, extra)
	result := <-channel
	if result.Error != nil {
		return "", "", result.Error
	}
	return result.Results[0].(string), result.Results[1].(string), result.Error
}

//// ParseResponseWithFees parses response for Indy request with fees
//func ParseResponseWithFees(wh int, paymentMethod string, resp string) (string, error) {
//	channel := payments.ParseResponseWithFees(wh, paymentMethod, resp)
//	result := <-channel
//	if result.Error != nil {
//		return "", result.Error
//	}
//	return result.Results[0].(string), result.Error
//}
//
//// BuildPaymentReq builds Indy request for doing payment according to this payment method
//func BuildPaymentReq(wh int, submitterDID string, req string, inputs string, outputs string, extra string) (string, string, error) {
//	channel := payments.BuildPaymentReq(wh, submitterDID, req, inputs, outputs, extra)
//	result := <-channel
//	if result.Error != nil {
//		return "", "", result.Error
//	}
//	return result.Results[0].(string), result.Results[1].(string), result.Error
//}
//
//// ParsePaymentResponse parses response for Indy request for payment txn
//func ParsePaymentResponse(wh int, paymentMethod string, resp string) (string, error) {
//	channel := payments.ParsePaymentResponse(wh, paymentMethod, resp)
//	result := <-channel
//	if result.Error != nil {
//		return "", result.Error
//	}
//	return result.Results[0].(string), result.Error
//}
//
//// BuildMintReq builds Indy request for doing minting according to this payment method
//func BuildMintReq(wh int, submitterDID string, outputs string, extra string) (string, string, error) {
//	channel := payments.BuildMintReq(wh, submitterDID, outputs, extra)
//	result := <-channel
//	if result.Error != nil {
//		return "", "", result.Error
//	}
//	return result.Results[0].(string), result.Results[1].(string), result.Error
//}
//
//// BuildSetTxnFeesReq builds Indy request for setting fees for transactions in the ledger
//func BuildSetTxnFeesReq(wh int, submitterDID string, paymentMethod string, fees string) (string, error) {
//	channel := payments.BuildSetTxnFeesReq(wh, submitterDID, paymentMethod, fees)
//	result := <-channel
//	if result.Error != nil {
//		return "", result.Error
//	}
//	return result.Results[0].(string), result.Error
//}
//
//// BuildGetTxnFeesReq builds Indy request for getting fees for transactions in the ledger
//func BuildGetTxnFeesReq(wh int, submitterDID string, paymentMethod string) (string, error) {
//	channel := payments.BuildGetTxnFeesReq(wh, submitterDID, paymentMethod)
//	result := <-channel
//	if result.Error != nil {
//		return "", result.Error
//	}
//	return result.Results[0].(string), result.Error
//}
//
//// ParseGetTxnFeesResponse parses response for Indy request for getting fees
//func ParseGetTxnFeesResponse (paymentMethod string, response string) (string, error) {
//	channel := payments.ParseGetTxnFeesResponse(paymentMethod, response)
//	result := <-channel
//	if result.Error != nil {
//		return "", result.Error
//	}
//	return result.Results[0].(string), result.Error
//}
//
//// BuildVerifyPaymentReq builds Indy request for information to verify the payment receipt
//func BuildVerifyPaymentReq (wh int, submitterDID string, receipt string) (string, string, error) {
//	channel := payments.BuildVerifyPaymentReq(wh, submitterDID, receipt)
//	result := <-channel
//	if result.Error != nil {
//		return "", "", result.Error
//	}
//	return result.Results[0].(string), result.Results[1].(string), result.Error
//}
//
//// ParseVerifyPaymentResponse parses Indy response with information to verify receipt
//func ParseVerifyPaymentResponse (submitterDID string, receipt string) (string, error) {
//	channel := payments.ParseVerifyPaymentResponse(submitterDID, receipt)
//	result := <-channel
//	if result.Error != nil {
//		return "", result.Error
//	}
//	return result.Results[0].(string), result.Error
//}
//
//// SignWithAddress signs a message with a payment address
//func SignWithAddress (wh int, address string, messageRaw uint8, messageLen uint32) (string, error) {
//	channel := payments.SignWithAddress(wh, address, messageRaw, messageLen)
//	result := <-channel
//	if result.Error != nil {
//		return "", result.Error
//	}
//	return result.Results[0].(string), result.Error
//}
//
//// VerifyWithAddress verify a signature with a payment address
//func VerifyWithAddress (address string, messageRaw uint8, messageLen uint32, signatureRaw uint8, signatureLen uint32) (bool, error) {
//	channel := payments.VerifyWithAddress(address, messageRaw, messageLen, signatureRaw, signatureLen)
//	result := <-channel
//	if result.Error != nil {
//		return false, result.Error
//	}
//	return result.Results[0].(bool), result.Error
//}
