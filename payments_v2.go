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

import "indySDK/payments_v2"

// BuildGetPaymentSourcesWithFromRequest purge credential definition cache
func BuildGetPaymentSourcesWithFromRequest(wh int, submitterDID string, paymentAddress string, from int64) (string, string, error) {
	channel := payments_v2.BuildGetPaymentSourcesWithFromRequest(wh, submitterDID, paymentAddress, from)
	result := <-channel
	if result.Error != nil {
		return "", "", result.Error
	}
	return result.Results[0].(string), result.Results[1].(string), result.Error
}

// ParseGetPaymentSourcesWithFromResponse parses response for Indy request for getting sources list
func ParseGetPaymentSourcesWithFromResponse(paymentMethod string, respJs string) (int, string, error) {
	channel := payments_v2.ParseGetPaymentSourcesWithFromResponse(paymentMethod, respJs)
	result := <-channel
	if result.Error != nil {
		return 0, "", result.Error
	}
	return result.Results[0].(int), result.Results[1].(string), result.Error
}
