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

/*
#include <stdlib.h>
*/
import "C"
import (
	"github.com/joyride9999/IndySdkGoBindings/pairwise"
	"unsafe"
)

// IsPairwiseExists purge credential definition cache
func IsPairwiseExists(wh int, theirDID string) (bool, error) {

	upTheirDid := unsafe.Pointer(C.CString(theirDID))
	defer C.free(upTheirDid)

	channel := pairwise.IsPairwiseExists(wh, upTheirDid)
	result := <-channel
	if result.Error != nil {
		return false, result.Error
	}
	return result.Results[0].(bool), result.Error
}

// CreatePairwise creates pairwise
func CreatePairwise(wh int, theirDID, myDID, meta string) error {

	upTheirDid := unsafe.Pointer(C.CString(theirDID))
	defer C.free(upTheirDid)
	upMyDid := unsafe.Pointer(C.CString(myDID))
	defer C.free(upMyDid)
	upMeta := unsafe.Pointer(C.CString(meta))
	defer C.free(upMeta)

	channel := pairwise.CreatePairwise(wh, upTheirDid, upMyDid, upMeta)
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
	upTheirDid := unsafe.Pointer(C.CString(theirDID))
	defer C.free(upTheirDid)

	channel := pairwise.GetPairwise(wh, upTheirDid)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// SetPairwiseMetadata get list of saved pairwise
func SetPairwiseMetadata(wh int, theirDID string, meta string) error {
	upTheirDid := unsafe.Pointer(C.CString(theirDID))
	defer C.free(upTheirDid)
	upMeta := unsafe.Pointer(C.CString(meta))
	defer C.free(upMeta)

	channel := pairwise.SetPairwiseMetadata(wh, upTheirDid, upMeta)
	result := <-channel
	return result.Error
}
