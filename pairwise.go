/*
// ******************************************************************
// Purpose: exported public functions that handles pairwise functions
// from libindy
// Author:  adrian.toader@siemens.com
// Notes:
// Copyright (c): Siemens SRL
// This work is licensed under the terms of the Apache License Version 2.0.  See
// the LICENSE.txt file in the top-level directory.
// ******************************************************************
*/

package indySDK

import "github.com/joyride9999/IndySdkGoBindings/pairwise"

// IsPairwiseExists purge credential definition cache
func IsPairwiseExists(wh int, theirDID string) (bool, error) {
	channel := pairwise.IsPairwiseExists(wh, theirDID)
	result := <-channel
	if result.Error != nil {
		return false, result.Error
	}
	return result.Results[0].(bool), result.Error
}

// CreatePairwise creates pairwise
func CreatePairwise(wh int, theirDID string, myDID, meta string) error {
	channel := pairwise.CreatePairwise(wh, theirDID, myDID, meta)
	result := <-channel
	return result.Error
}

// ListPairwise get list of saved pairwise.
func ListPairwise(wh int) (string, error) {
	channel := pairwise.ListPairwise(wh)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// GetPairwise gets pairwise information for specific their_did
func GetPairwise(wh int, theirDID string) (string, error) {
	channel := pairwise.GetPairwise(wh, theirDID)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// SetPairwiseMetadata get list of saved pairwise
func SetPairwiseMetadata(wh int, theirDID string, meta string) error {
	channel := pairwise.SetPairwiseMetadata(wh, theirDID, meta)
	result := <-channel
	return result.Error
}
