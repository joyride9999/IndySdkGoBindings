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

package indySDK

/*
#include <stdlib.h>
*/
import "C"
import (
	"github.com/joyride9999/IndySdkGoBindings/mod"
	"encoding/json"
	"errors"
	"unsafe"
)

// SetRuntimeConfig set libindy runtime configuration
func SetRuntimeConfig(config mod.Config) error {

	jsonConfig, err := json.Marshal(config)
	if err != nil {
		return errors.New("cant read json")
	}

	upCfgJson := unsafe.Pointer(C.CString(string(jsonConfig)))
	defer C.free(upCfgJson)

	channel := mod.SetRuntimeConfig(upCfgJson)
	result := <-channel
	return result.Error
}
