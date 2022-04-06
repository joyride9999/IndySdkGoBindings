/*
// ******************************************************************
// Purpose: Wrapper to call libindy, imports functions from indy_crypto.h
// Author:  adrian.toader@siemens.com
// Notes:
// Copyright (c): Siemens SRL
// This work is licensed under the terms of the Apache License Version 2.0.  See
// the LICENSE.txt file in the top-level directory.
// ******************************************************************
*/

package crypto

// Key represents key information as json
type Key struct {
	Seed 		string `json:"seed,omitempty"`
	CryptoType 	string `json:"crypto_type,omitempty"`
}
