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

type ConfigBlobStorage struct {
	BaseDir    string `json:"base_dir"`
	UriPattern string `json:"uri_pattern"`
}
