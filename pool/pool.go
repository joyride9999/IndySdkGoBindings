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

/*
#cgo CFLAGS: -I ../include
#cgo LDFLAGS: -L${SRCDIR}/../lib -lindy
#include <indy_core.h>

typedef void (*cb_setProtocolVersion)(indy_handle_t, indy_error_t);
extern void setProtocolVersionCB(indy_handle_t, indy_error_t);

typedef void (*cb_createPoolLedgerConfig)(indy_handle_t, indy_error_t);
extern void createPoolLedgerConfigCB(indy_handle_t, indy_error_t);

typedef void (*cb_openPoolLedger)(indy_handle_t, indy_error_t, indy_handle_t);
extern void openPoolLedgerCB(indy_handle_t, indy_error_t, indy_handle_t);

typedef void (*cb_closePoolLedger)(indy_handle_t, indy_error_t);
extern void closePoolLedgerCB(indy_handle_t, indy_error_t);

*/
import "C"

import (
	"encoding/json"
	"errors"
	"indySDK/indyUtils"
	"unsafe"
)

//export closePoolLedgerCB
func closePoolLedgerCB(commandHandle C.indy_handle_t, indyError C.indy_error_t) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// IndyClosePoolHandle closes pool handle
func IndyClosePoolHandle(ph int) chan indyUtils.IndyResult {

	handle, future := indyUtils.NewFutureCommand()

	commandHandle := (C.indy_handle_t)(handle)
	res := C.indy_close_pool_ledger(commandHandle,
		(C.indy_handle_t)(ph),
		(C.cb_closePoolLedger)(unsafe.Pointer(C.closePoolLedgerCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export setProtocolVersionCB
func setProtocolVersionCB(commandHandle C.indy_handle_t, indyError C.indy_error_t) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// IndySetProtocolVersion sets the communication protocol with the pool
func IndySetProtocolVersion(pv uint64) chan indyUtils.IndyResult {

	handle, future := indyUtils.NewFutureCommand()

	commandHandle := (C.indy_handle_t)(handle)
	protocolVersion := (C.indy_u64_t)(pv)

	res := C.indy_set_protocol_version(commandHandle,
		protocolVersion,
		(C.cb_setProtocolVersion)(unsafe.Pointer(C.setProtocolVersionCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export createPoolLedgerConfigCB
func createPoolLedgerConfigCB(commandHandle C.indy_handle_t, indyError C.indy_error_t) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// IndyCreatePoolLedgerConfig creates a pool configuration out of a txn file
func IndyCreatePoolLedgerConfig(pool Pool) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)
	poolConfig, err := json.Marshal(pool)
	if err != nil {
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: err}) }()
		return future
	}
	poolCfg := string(poolConfig)

	// Call indy_create_pool_ledger_config
	res := C.indy_create_pool_ledger_config(commandHandle,
		C.CString(pool.Name),
		C.CString(poolCfg),
		(C.cb_setProtocolVersion)(unsafe.Pointer(C.createPoolLedgerConfigCB)))

	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export openPoolLedgerCB
func openPoolLedgerCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, retHandle C.indy_handle_t) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil, Results: []interface{}{int(retHandle)}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// IndyOpenPoolLedger opens a pool
func IndyOpenPoolLedger(pool Pool) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	type cfg struct {
		Timeout int `json:"timeout"`
	}

	t := cfg{
		Timeout: 10,
	}

	jsonConfig, err := json.Marshal(t)
	if err != nil {
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: err}) }()
		return future
	}

	/*
	   Opens pool ledger and performs connecting to pool nodes.

	   Pool ledger configuration with corresponded name must be previously created
	   with indy_create_pool_ledger_config method.
	   It is impossible to open pool with the same name more than once.

	   :param config_name: Name of the pool ledger configuration.
	   :param config: (optional) Runtime pool configuration json.
	    if NULL, then default config will be used. Example:
	       {
	           "timeout": int (optional), timeout for network request (in sec).
	           "extended_timeout": int (optional), extended timeout for network request (in sec).
	           "preordered_nodes": array<string> -  (optional), names of nodes which will have a priority during request sending:
	               ["name_of_1st_prior_node",  "name_of_2nd_prior_node", .... ]
	               This can be useful if a user prefers querying specific nodes.
	               Assume that `Node1` and `Node2` nodes reply faster.
	               If you pass them Libindy always sends a read request to these nodes first and only then (if not enough) to others.
	               Note: Nodes not specified will be placed randomly.
	           "number_read_nodes": int (optional) - the number of nodes to send read requests (2 by default)
	               By default Libindy sends a read requests to 2 nodes in the pool.
	               If response isn't received or `state proof` is invalid Libindy sends the request again but to 2 (`number_read_nodes`) * 2 = 4 nodes and so far until completion.
	       }
	   :return: Handle to opened pool to use in methods that require pool connection.
	*/
	// Call indy_open_pool_ledger
	res := C.indy_open_pool_ledger(commandHandle,
		C.CString(pool.Name),
		C.CString(string(jsonConfig)),
		(C.cb_setProtocolVersion)(unsafe.Pointer(C.openPoolLedgerCB)))

	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}
