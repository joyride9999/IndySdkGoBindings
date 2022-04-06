/*
// ******************************************************************
// Purpose: exported public functions that handles metrics functions
// from libindy
// Author:  adrian.toader@siemens.com
// Notes:
// Copyright (c): Siemens SRL
// This work is licensed under the terms of the Apache License Version 2.0.  See
// the LICENSE.txt file in the top-level directory.
// ******************************************************************
*/

package indySDK

import "indySDK/metrics"

// Collect collect metrics
func Collect() (string, error) {
	channel := metrics.Collect()
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}
