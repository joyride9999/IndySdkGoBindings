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

// Config represents payment address config
type Config struct {
	Seed	string `json:"seed"`
}
