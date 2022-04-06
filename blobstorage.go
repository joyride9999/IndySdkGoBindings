/*
// ******************************************************************
// Purpose: exported public functions that handles blobstorage functions
// from libindy
// Author:  alexandru.leonte@siemens.com
// Notes:
// Copyright (c): Siemens SRL
// This work is licensed under the terms of the Apache License Version 2.0.  See
// the LICENSE.txt file in the top-level directory.
// ******************************************************************
*/

package indySDK

import "indySDK/blobstorage"

// IndyOpenBlobStorageReader opens blob reader
func IndyOpenBlobStorageReader(blobStorageType string, config string) (blobHandle int, err error) {
	channel := blobstorage.OpenBlobStorageReader(blobStorageType, config)
	result := <-channel
	if result.Error != nil {
		return -1, result.Error
	}

	return result.Results[0].(int), result.Error
}

// IndyOpenBlobStorageWriter opens blob writer
func IndyOpenBlobStorageWriter(blobStorageType string, config string) (blobHandle int, err error) {
	channel := blobstorage.OpenBlobStorageWriter(blobStorageType, config)
	result := <-channel
	if result.Error != nil {
		return -1, result.Error
	}

	return result.Results[0].(int), result.Error
}
