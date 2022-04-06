/*
// ******************************************************************
// Purpose: Wrapper to call libindy, imports functions from indy_did.h
// Author:  alexandru.leonte@siemens.com
// Notes:
// Copyright (c): Siemens SRL
// This work is licensed under the terms of the Apache License Version 2.0.  See
// the LICENSE.txt file in the top-level directory.
// ******************************************************************
*/

package did

type IdentityKey struct {
	Seed		string `json:"seed,omitempty"`
}

type IdentityDID struct {
	Did			string `json:"did"`
	VerKey		string `json:"verkey,omitempty"`
}
