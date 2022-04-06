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

/*
#cgo CFLAGS: -I ../include
#cgo LDFLAGS: -L${SRCDIR}/../lib -lindy
#include <indy_core.h>
typedef void (*cb_getSchema)(indy_handle_t, indy_error_t, char*);
extern void getSchemaCB(indy_handle_t, indy_error_t, char*);

typedef void (*cb_getCredDef)(indy_handle_t, indy_error_t, char*);
extern void getCredDefCB(indy_handle_t, indy_error_t, char*);

typedef void (*cb_purgeCredDefCache)(indy_handle_t, indy_error_t);
extern void purgeCredDefCacheCB(indy_handle_t, indy_error_t);

typedef void (*cb_purgeSchemaCache)(indy_handle_t, indy_error_t);
extern void purgeSchemaCacheCB(indy_handle_t, indy_error_t);
*/
import "C"

import (
	"errors"
	"github.com/joyride9999/IndySdkGoBindings/indyUtils"
	"unsafe"
)

//export getSchemaCB
func getSchemaCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, schemaJson *C.char) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle),
			indyUtils.IndyResult{Error: nil,
				Results: []interface{}{
					string(C.GoString(schemaJson)),
				}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

func GetSchema(poolHandle int, walletHandle int, submitterDid string, schemaID string, optionsJson string) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
			Gets schema json data for specified schema id.

		    If data is present inside of cache, cached data is returned.
		    Otherwise data is fetched from the ledger and stored inside of cache for future use.

		    EXPERIMENTAL

		    :param pool_handle: pool handle (created by open_pool_ledger).
		    :param wallet_handle: wallet handle (created by open_wallet).
		    :param submitter_did: DID of the submitter stored in secured Wallet.
		    :param id: identifier of schema.
		    :param options_json:
		    {
		        noCache: (bool, optional, false by default) Skip usage of cache,
		        noUpdate: (bool, optional, false by default) Use only cached data, do not try to update.
		        noStore: (bool, optional, false by default) Skip storing fresh data if updated,
		        minFresh: (int, optional, -1 by default) Return cached data if not older than this many seconds. -1 means do not check age.
		    }
		    :return: Schema json.
		    {
		        id: identifier of schema
		        attrNames: array of attribute name strings
		        name: Schema's name string
		        version: Schema's version string
		        ver: Version of the Schema json
		    }
	*/

	// Call C.indy_get_schema
	res := C.indy_get_schema(commandHandle,
		C.indy_handle_t(poolHandle),
		C.indy_handle_t(walletHandle),
		C.CString(submitterDid),
		C.CString(schemaID),
		C.CString(optionsJson),
		(C.cb_getSchema)(unsafe.Pointer(C.getSchemaCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export getCredDefCB
func getCredDefCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, credDefJson *C.char) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle),
			indyUtils.IndyResult{Error: nil,
				Results: []interface{}{
					string(C.GoString(credDefJson)),
				}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

func GetCredDef(poolHandle int, walletHandle int, submitterDid string, credDefId string, optionsJson string) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*

	   Gets credential definition json data for specified credential definition id.

	   If data is present inside of cache, cached data is returned.
	   Otherwise data is fetched from the ledger and stored inside of cache for future use.

	   EXPERIMENTAL

	   :param pool_handle: pool handle (created by open_pool_ledger).
	   :param wallet_handle: wallet handle (created by open_wallet).
	   :param submitter_did: DID of the submitter stored in secured Wallet.
	   :param id: identifier of credential definition.
	   :param options_json:
	   {
	       noCache: (bool, optional, false by default) Skip usage of cache,
	       noUpdate: (bool, optional, false by default) Use only cached data, do not try to update.
	       noStore: (bool, optional, false by default) Skip storing fresh data if updated,
	       minFresh: (int, optional, -1 by default) Return cached data if not older than this many seconds. -1 means do not check age.
	   }
	   :return: Credential Definition json.
	   {
	       id: string - identifier of credential definition
	       schemaId: string - identifier of stored in ledger schema
	       type: string - type of the credential definition. CL is the only supported type now.
	       tag: string - allows to distinct between credential definitions for the same issuer and schema
	       value: Dictionary with Credential Definition's data: {
	           primary: primary credential public key,
	           Optional<revocation>: revocation credential public key
	       },
	       ver: Version of the Credential Definition json
	   }
	*/

	// Call C.indy_get_cred_def
	res := C.indy_get_cred_def(commandHandle,
		C.indy_handle_t(poolHandle),
		C.indy_handle_t(walletHandle),
		C.CString(submitterDid),
		C.CString(credDefId),
		C.CString(optionsJson),
		(C.cb_getCredDef)(unsafe.Pointer(C.getCredDefCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export purgeCredDefCacheCB
func purgeCredDefCacheCB(commandHandle C.indy_handle_t, indyError C.indy_error_t) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// PurgeCredDefCache purge credential definition cache
func PurgeCredDefCache(wh int, options string) chan indyUtils.IndyResult {
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
			Purge credential definition cache.

		    EXPERIMENTAL

		    :param wallet_handle: wallet handle (used for cache)
		    :param options_json:
		    {
		        maxAge: (int, optional, -1 by default) Purge cached data if older than this many seconds. -1 means purge all.
		    }

		    :return: None
	*/

	// Call indy_purge_cred_def_cache
	res := C.indy_purge_cred_def_cache(commandHandle,
		(C.indy_handle_t)(wh),
		C.CString(options),
		(C.cb_purgeCredDefCache)(unsafe.Pointer(C.purgeCredDefCacheCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export purgeSchemaCacheCB
func purgeSchemaCacheCB(commandHandle C.indy_handle_t, indyError C.indy_error_t) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// PurgeSchemaCache purge credential definition cache
func PurgeSchemaCache(wh int, options string) chan indyUtils.IndyResult {
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
			Purge schema cache.

		    EXPERIMENTAL

		    :param wallet_handle: wallet handle (used for cache)
		    :param options_json:
		    {
		        maxAge: (int, optional, -1 by default) Purge cached data if older than this many seconds. -1 means purge all.
		    }
		    :return: None
	*/

	// Call indy_purge_cred_def_cache
	res := C.indy_purge_schema_cache(commandHandle,
		(C.indy_handle_t)(wh),
		C.CString(options),
		(C.cb_purgeSchemaCache)(unsafe.Pointer(C.purgeSchemaCacheCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}
