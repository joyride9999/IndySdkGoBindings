/*
// ******************************************************************
// Purpose: Wrapper to call libindy, imports functions from indy_pool.h
// Author:  alexandru.leonte@siemens.com
// Notes:
// Copyright (c): Siemens SRL
// This work is licensed under the terms of the Apache License Version 2.0.  See
// the LICENSE.txt file in the top-level directory.
// ******************************************************************
*/

package pool

type Pool struct {
	Name       string `json:"name"`
	GenesisTxn string `json:"genesis_txn"`
}
