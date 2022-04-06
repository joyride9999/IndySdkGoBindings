/*
// ******************************************************************
// Purpose: Wrapper to call libindy, imports functions from indy_cache.h
// Author:  adrian.toader@siemens.com
// Notes:
// Copyright (c): Siemens SRL
// This work is licensed under the terms of the Apache License Version 2.0.  See
// the LICENSE.txt file in the top-level directory.
// ******************************************************************
*/

package cache

// Options represents Indy cache options
type Options struct {
	NoCache bool `json:"noCache"`
	NoUpdate bool `json:"noUpdate"`
	NoStore  bool `json:"nostore"`
	MinFresh int  `json:"minFresh"`
}
