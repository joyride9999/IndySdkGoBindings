/*
// ******************************************************************
// Purpose: Wrapper to call libindy, imports functions from indy_mod.h
// Author:  adrian.toader@siemens.com
// Notes:
// Copyright (c): Siemens SRL
// This work is licensed under the terms of the Apache License Version 2.0.  See
// the LICENSE.txt file in the top-level directory.
// ******************************************************************
*/

package mod

// Config represents libindy runtime configuration
type Config struct {
	CryptoThreadPoolSize	int `json:"crypto_thread_pool_size"`
	CollectBacktrace		bool `json:"collect_backtrace"`
}
