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

/*
#cgo CFLAGS: -I ../include
#cgo LDFLAGS: -L${SRCDIR}/../lib -lindy
#include <indy_core.h>

*/
import "C"
import (
	"encoding/json"
	"errors"
	"github.com/joyride9999/IndySdkGoBindings/indyUtils"
)

// SetRuntimeConfig set libindy runtime configuration
func SetRuntimeConfig(config Config) chan indyUtils.IndyResult {
	handle, future := indyUtils.NewFutureCommand()

	jsonCoinfig, err := json.Marshal(config)
	if err != nil {
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: err}) }()
		return future
	}

	configString := string(jsonCoinfig)

	/*
		Set libindy runtime configuration. Can be optionally called to change current params.

		:param config: {
		    "crypto_thread_pool_size": Optional<int> - size of thread pool for the most expensive crypto operations. (4 by default)
		    "collect_backtrace": Optional<bool> - whether errors backtrace should be collected.
		        Capturing of backtrace can affect library performance.
		        NOTE: must be set before invocation of any other API functions.
		}
	*/

	// Call indy_set_runtime_config
	res := C.indy_set_runtime_config(C.CString(configString))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}
