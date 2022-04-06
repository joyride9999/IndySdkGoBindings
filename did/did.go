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

/*
#cgo CFLAGS: -I ../include
#cgo LDFLAGS: -L${SRCDIR}/../lib -lindy
#include <indy_core.h>

typedef void (*cb_createAndStoreMyDid)(indy_handle_t, indy_error_t, char*, char*);
extern void createAndStoreMyDidCB(indy_handle_t, indy_error_t, char*, char*);

typedef void (*cb_replaceKeyStart)(indy_handle_t, indy_error_t, char*);
extern void replaceKeyStartCB(indy_handle_t, indy_error_t, char*);

typedef void (*cb_replaceKeyApply)(indy_handle_t, indy_error_t);
extern void replaceKeyApplyCB(indy_handle_t, indy_error_t);

typedef void (*cb_storeTheirDid)(indy_handle_t, indy_error_t);
extern void storeTheirDidCB(indy_handle_t, indy_error_t);

typedef void (*cb_keyForDid)(indy_handle_t, indy_error_t, char*);
extern void keyForDidCB(indy_handle_t, indy_error_t, char*);

typedef void (*cb_keyForLocalDid)(indy_handle_t, indy_error_t, char*);
extern void keyForLocalDidCB(indy_handle_t, indy_error_t, char*);

typedef void (*cb_setEndPointForDid)(indy_handle_t, indy_error_t);
extern void setEndPointForDidCB(indy_handle_t, indy_error_t);

typedef void (*cb_getEndPointForDid)(indy_handle_t, indy_error_t, char*, char*);
extern void getEndPointForDidCB(indy_handle_t, indy_error_t, char*, char*);

typedef void (*cb_setDidMetadata)(indy_handle_t, indy_error_t);
extern void setDidMetadataCB(indy_handle_t, indy_error_t);

typedef void (*cb_getDidMetadata)(indy_handle_t, indy_error_t, char*);
extern void getDidMetadataCB(indy_handle_t, indy_error_t, char*);

typedef void (*cb_getDidWithMetadata)(indy_handle_t, indy_error_t, char*);
extern void getDidWithMetadataCB(indy_handle_t, indy_error_t, char*);

typedef void (*cb_listDidsWithMeta)(indy_handle_t, indy_error_t, char*);
extern void listDidsWithMetaCB(indy_handle_t, indy_error_t, char*);

typedef void (*cb_abbreviateVerKey)(indy_handle_t, indy_error_t, char*);
extern void abbreviateVerKeyCB(indy_handle_t, indy_error_t, char*);

typedef void (*cb_qualifyDid)(indy_handle_t, indy_error_t, char*);
extern void qualifyDidCB(indy_handle_t, indy_error_t, char*);
*/
import "C"
import (
	"encoding/json"
	"errors"
	"github.com/joyride9999/IndySdkGoBindings/indyUtils"
	"unsafe"
)

//export createAndStoreMyDidCB
func createAndStoreMyDidCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, did *C.char, verKey *C.char) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle),
			indyUtils.IndyResult{Error: nil,
				Results: []interface{}{
					string(C.GoString(did)),
					string(C.GoString(verKey)),
				}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// CreateAndStoreMyDid creates a did and returns it with verkey. Nothing is written to the blockchain...
func CreateAndStoreMyDid(walletHandle int, seed string) chan indyUtils.IndyResult {
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	type didJson struct {
		Did        string `json:"did,omitempty"`
		Seed       string `json:"seed,omitempty"`
		Crypto     string `json:"crypto_type,omitempty"`
		Cid        string `json:"cid,omitempty"`
		MethodName string `json:"method_name,omitempty"`
	}

	didjs := didJson{
		Seed: seed,
	}

	didCfg, err := json.Marshal(didjs)
	if err != nil {
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: err}) }()
		return future
	}

	/*
		Create keys (signing and encryption keys) for a new
		DID (owned by the caller of the library).
		Identity's DID must be either explicitly provided, or taken as the first 16 bit of verkey.
		Saves the Identity DID with keys in a secured Wallet, so that it can be used to sign
		and encrypt transactions.

		:param wallet_handle: wallet handler (created by open_wallet).
		:param did_json: Identity information as json. Example:
			{
				"did": string, (optional;
						if not provided and cid param is false then the first 16 bit of the verkey will be
						used as a new DID;
						if not provided and cid is true then the full verkey will be used as a new DID;
						if provided, then keys will be replaced - key rotation use case)
				"seed": string, (optional) Seed that allows deterministic key creation (if not set random one will be created).
											Can be UTF-8, base64 or hex string.
				"crypto_type": string, (optional; if not set then ed25519 curve is used;
						  currently only 'ed25519' value is supported for this field)
				"cid": bool, (optional; if not set then false is used;)
				"method_name": string, (optional) method name to create fully qualified did.
			}
		:return: DID and verkey (for verification of signature)
	*/

	// Call indy_create_and_store_my_did
	res := C.indy_create_and_store_my_did(commandHandle,
		(C.indy_handle_t)(walletHandle),
		C.CString(string(didCfg)),
		(C.cb_createAndStoreMyDid)(unsafe.Pointer(C.createAndStoreMyDidCB)))

	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export replaceKeyStartCB
func replaceKeyStartCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, verKey *C.char) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil, Results: []interface{}{string(C.GoString(verKey))}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// ReplaceKeyStart generates temporary keys for an existing DID
func ReplaceKeyStart(walletHandle int, did string, identityJson string) chan indyUtils.IndyResult {
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
			Generated temporary keys (signing and encryption keys) for an existing
		 	DID (owned by the caller of the library).

			:param wallet_handle: wallet handler (created by open_wallet).
			:param command_handle: command handle to map callback to user context.
			:param identity_json: Identity information as json. Example:
			{
			    "seed": string, (optional) Seed that allows deterministic key creation (if not set random one will be created).
			                               Can be UTF-8, base64 or hex string.
			    "crypto_type": string, (optional; if not set then ed25519 curve is used;
			              currently only 'ed25519' value is supported for this field)
			}

			:return verKey
	*/

	// Call indy_replace_keys_start
	res := C.indy_replace_keys_start(commandHandle,
		(C.indy_handle_t)(walletHandle),
		C.CString(did),
		C.CString(identityJson),
		(C.cb_replaceKeyStart)(unsafe.Pointer(C.replaceKeyStartCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export replaceKeyApplyCB
func replaceKeyApplyCB(commandHandle C.indy_handle_t, indyError C.indy_error_t) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// ReplaceKeyApply applies temporary keys as main for existing DID
func ReplaceKeyApply(walletHandle int, did string) chan indyUtils.IndyResult {
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
		Apply temporary keys as main for an existing DID (owned by the caller of the library).

		:param wallet_handle: wallet handler (created by open_wallet).
		:param command_handle: command handle to map callback to user context.
		:param did: DID stored in the wallet

		:return: error code
	*/

	// Call indy_replace_keys_apply
	res := C.indy_replace_keys_apply(commandHandle,
		(C.indy_handle_t)(walletHandle),
		C.CString(did),
		(C.cb_replaceKeyApply)(unsafe.Pointer(C.replaceKeyApplyCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export storeTheirDidCB
func storeTheirDidCB(commandHandle C.indy_handle_t, indyError C.indy_error_t) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// StoreTheirDid saves DID for a pairwise connection in a secured wallet to verify transaction.
func StoreTheirDid(walletHandle int, identityJson string) chan indyUtils.IndyResult {
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
		Saves their DID for a pairwise connection in a secured Wallet,
		so that it can be used to verify transaction.

		:param wallet_handle: wallet handler (created by open_wallet).
		:param identity_json: Identity information as json. Example:
		{
				"did": string, (required)
				"verkey": string
					- optional is case of adding a new DID, and DID is cryptonym: did == verkey,
					- mandatory in case of updating an existing DID
		}
		:return: error code
	*/

	// Call indy_store_their_did
	res := C.indy_store_their_did(commandHandle,
		(C.indy_handle_t)(walletHandle),
		C.CString(identityJson),
		(C.cb_storeTheirDid)(unsafe.Pointer(C.storeTheirDidCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export keyForDidCB
func keyForDidCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, verKey *C.char) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil, Results: []interface{}{string(C.GoString(verKey))}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// KeyForDid returns ver key for the given DID
func KeyForDid(poolHandle int, walletHandle int, did string) chan indyUtils.IndyResult {
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
		Returns ver key (key id) for the given DID.

		"indy_key_for_did" call follow the idea that we resolve information about their DID from
		the ledger with cache in the local wallet. The "indy_open_wallet" call has freshness parameter
		that is used for checking the freshness of cached pool value.

		Note if you don't want to resolve their DID info from the ledger you can use
		"indy_key_for_local_did" call instead that will look only to the local wallet and skip
		freshness checking.

		Note that "indy_create_and_store_my_did" makes similar wallet record as "indy_create_key".
		As result we can use returned ver key in all generic crypto and messaging functions.

		:param wallet_handle: Wallet handle (created by open_wallet).
		:param did: The DID to resolve key.

		:return: error code
	*/

	// Call indy_key_for_did
	res := C.indy_key_for_did(commandHandle,
		C.indy_handle_t(poolHandle),
		C.indy_handle_t(walletHandle),
		C.CString(did),
		(C.cb_keyForDid)(unsafe.Pointer(C.keyForDidCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}
	return future
}

//export keyForLocalDidCB
func keyForLocalDidCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, verKey *C.char) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil, Results: []interface{}{string(C.GoString(verKey))}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// KeyForLocalDid gets the key for the local DID.
func KeyForLocalDid(walletHandle int, did string) chan indyUtils.IndyResult {
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	res := C.indy_key_for_local_did(commandHandle,
		(C.indy_handle_t)(walletHandle),
		C.CString(did),
		(C.cb_keyForLocalDid)(unsafe.Pointer(C.keyForLocalDidCB)))

	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export setEndPointForDidCB
func setEndPointForDidCB(commandHandle C.indy_handle_t, indyError C.indy_error_t) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// SetEndPointForDid set/replaces endpoint information for the given DID
func SetEndPointForDid(walletHandle int, did string, address string, transportKey string) chan indyUtils.IndyResult {
	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := C.indy_handle_t(handle)

	/*
			Set/replaces endpoint information for the given DID.

		    :param wallet_handle: Wallet handle (created by open_wallet).
		    :param did: The DID to resolve endpoint.
		    :param address: The DIDs endpoint address.
		    :param transport_key: The DIDs transport key (ver key, key id).

		    :return: Error code
	*/

	// Call indy_set_endpoint_for_did
	res := C.indy_set_endpoint_for_did(commandHandle,
		(C.indy_handle_t)(walletHandle),
		C.CString(did),
		C.CString(address),
		C.CString(transportKey),
		(C.cb_setEndPointForDid)(unsafe.Pointer(C.setEndPointForDidCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export getEndPointForDidCB
func getEndPointForDidCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, address *C.char, transportKey *C.char) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil, Results: []interface{}{string(C.GoString(address)),
			string(C.GoString(transportKey))}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// GetEndPointForDid returns endpoint information for the given DID
func GetEndPointForDid(walletHandle int, poolHandle int, did string) chan indyUtils.IndyResult {
	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := C.indy_handle_t(handle)

	/*
			Returns endpoint information for the given DID.
		    :param wallet_handle: Wallet handle (created by open_wallet).
		    :param pool_handle: Pool handle (created by open_pool).
		    :param did: The DID to resolve endpoint.

		    :return: (endpoint, transport_vk)
	*/

	// Call indy_get_endpoint_for_did
	res := C.indy_get_endpoint_for_did(commandHandle,
		(C.indy_handle_t)(walletHandle),
		(C.indy_handle_t)(poolHandle),
		C.CString(did),
		(C.cb_getEndPointForDid)(unsafe.Pointer(C.getEndPointForDidCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export setDidMetadataCB
func setDidMetadataCB(commandHandle C.indy_handle_t, indyError C.indy_error_t) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// SetDidMetadata saves/replaces meta information for the given DID
func SetDidMetadata(walletHandle int, did string, metadata string) chan indyUtils.IndyResult {
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
	 	Saves/replaces the meta information for the given DID in the wallet.

	    :param wallet_handle: Wallet handle (created by open_wallet).
	    :param did: the DID to store metadata.
	    :param metadata: the meta information that will be store with the DID.
	    :return: Error code
	*/

	// Call indy_set_did_metadata
	res := C.indy_set_did_metadata(commandHandle,
		C.indy_handle_t(walletHandle),
		C.CString(did),
		C.CString(metadata),
		(C.cb_setDidMetadata)(unsafe.Pointer(C.setDidMetadataCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export getDidMetadataCB
func getDidMetadataCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, metadata *C.char) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil, Results: []interface{}{string(C.GoString(metadata))}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// GetDidMetadata retrieves meta information for the given DID
func GetDidMetadata(walletHandle int, did string) chan indyUtils.IndyResult {
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := C.indy_handle_t(handle)

	/*
			Retrieves the meta information for the given DID in the wallet.

		    :param wallet_handle: Wallet handle (created by open_wallet).
		    :param did: The DID to retrieve metadata.
		    :return: metadata: The meta information stored with the DID; Can be null if no metadata was saved for this DID.
	*/

	// Call indy_get_did_metadata
	res := C.indy_get_did_metadata(commandHandle,
		(C.indy_handle_t)(walletHandle),
		C.CString(did),
		(C.cb_getDidMetadata)(unsafe.Pointer(C.getDidMetadataCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export getDidWithMetadataCB
func getDidWithMetadataCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, did *C.char) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil, Results: []interface{}{
			string(C.GoString(did))}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// GetDidWithMetadata retrieves DID, metadata and verKey stored in the wallet.
func GetDidWithMetadata(walletHandle int, did string) chan indyUtils.IndyResult {
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := C.indy_handle_t(handle)
	/*
			Get DID metadata and verkey stored in the wallet.

			:param wallet_handle: wallet handler (created by open_wallet).
		    :param did: The DID to retrieve metadata.
		    :return: DID with verkey and metadata.
	*/

	// Call indy_get_did_with_metadata
	res := C.indy_get_my_did_with_meta(commandHandle,
		(C.indy_handle_t)(walletHandle),
		C.CString(did),
		(C.cb_getDidWithMetadata)(unsafe.Pointer(C.getDidWithMetadataCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export listDidsWithMetaCB
func listDidsWithMetaCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, dids *C.char) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil, Results: []interface{}{
			string(C.GoString(dids))}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// ListDidsWithMeta lists DIDs and metadata stored in the wallet.
func ListDidsWithMeta(walletHandle int) chan indyUtils.IndyResult {
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
			List DIDs and metadata stored in the wallet.

		    :param wallet_handle: wallet handler (created by open_wallet).
		    :return: List of DIDs with verkeys and meta data.
	*/

	// Call indy_list_my_dids_with_meta
	res := C.indy_list_my_dids_with_meta(commandHandle,
		C.indy_handle_t(walletHandle),
		(C.cb_listDidsWithMeta)(unsafe.Pointer(C.listDidsWithMetaCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export abbreviateVerKeyCB
func abbreviateVerKeyCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, fullVerKey *C.char) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil, Results: []interface{}{
			string(C.GoString(fullVerKey))}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// AbbreviateVerKey retrieves abbreviated key if exists, otherwise returns full verkey.
func AbbreviateVerKey(did string, verKey string) chan indyUtils.IndyResult {
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := C.indy_handle_t(handle)

	/*
	 	Retrieves abbreviated verkey if it is possible otherwise return full verkey.

	    :param did: The DID.
	    :param full_verkey: The DIDs verification key,
	    :return: metadata: Either abbreviated or full verkey.
	*/

	// Call C.indy_abbreviate_verkey
	res := C.indy_abbreviate_verkey(commandHandle,
		C.CString(did),
		C.CString(verKey),
		(C.cb_abbreviateVerKey)(unsafe.Pointer(C.abbreviateVerKeyCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export qualifyDidCB
func qualifyDidCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, qualifiedDid *C.char) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil, Results: []interface{}{
			string(C.GoString(qualifiedDid))}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// QualifyDid updates DID related entities stored in the wallet.
func QualifyDid(walletHandle int, did string, method string) chan indyUtils.IndyResult {
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := C.indy_handle_t(handle)

	/*
			Update DID stored in the wallet to make fully qualified or to do other DID maintenance.
		        - If the DID has no prefix, a prefix will be appended (prepend did:peer to a legacy did)
		        - If the DID has a prefix, a prefix will be updated (migrate did:peer to did:peer-new)
		    Update DID related entities stored in the wallet.

		    :param wallet_handle: wallet handler (created by open_wallet).
		    :param did: target DID stored in the wallet.
		    :param method: method to apply to the DID.
		    :return: fully qualified did
	*/

	// Call indy_qualify_did
	res := C.indy_qualify_did(commandHandle,
		C.indy_handle_t(walletHandle),
		C.CString(did),
		C.CString(method),
		(C.cb_qualifyDid)(unsafe.Pointer(C.qualifyDidCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}
