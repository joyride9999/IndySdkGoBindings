/*
// ******************************************************************
// Purpose: Wrapper to call libindy, imports functions from indy_wallet.h
// Author:  alexandru.leonte@siemens.com
// Notes:
// Copyright (c): Siemens SRL
// This work is licensed under the terms of the Apache License Version 2.0.  See
// the LICENSE.txt file in the top-level directory.
// ******************************************************************
*/

package wallet

/*
#cgo CFLAGS: -I ../include
#cgo LDFLAGS: -L${SRCDIR}/../lib -lindy
#include <indy_core.h>
#include <stdlib.h>
#include <string.h>

typedef void (*cb_createWallet)(indy_handle_t, indy_error_t);
extern void createWalletCB(indy_handle_t, indy_error_t);

typedef void (*cb_openWallet)(indy_handle_t, indy_error_t, indy_handle_t);
extern void openWalletCB(indy_handle_t, indy_error_t, indy_handle_t);

typedef void (*cb_closeWallet)(indy_handle_t, indy_error_t);
extern void closeWalletCB(indy_handle_t, indy_error_t);

typedef void (*cb_deleteWallet)(indy_handle_t, indy_error_t);
extern void deleteWalletCB(indy_handle_t, indy_error_t);

typedef void (*cb_generateWalletKey)(indy_handle_t, indy_error_t);
extern void generateWalletKeyCB(indy_handle_t, indy_error_t);

typedef void (*cb_importWallet)(indy_handle_t, indy_error_t);
extern void importWalletCB(indy_handle_t, indy_error_t);

typedef void (*cb_exportWallet)(indy_handle_t, indy_error_t);
extern void exportWalletCB(indy_handle_t, indy_error_t);

typedef void (*cb_registerWalletStorage)(indy_handle_t, indy_error_t);
extern void registerWalletStorageCB(indy_handle_t, indy_error_t);

typedef indy_error_t (*cb_createWalletCustom)(char*, char*, char*, char*);
extern indy_error_t createWalletCustomCB(char*, char*, char*, char*);

typedef indy_error_t (*cb_openWalletCustom)(char*, char*, char*, indy_handle_t*);
extern indy_error_t openWalletCustomCB(char*, char*, char*, indy_handle_t*);

typedef indy_error_t (*cb_closeWalletCustom)(indy_handle_t);
extern indy_error_t closeWalletCustomCB(indy_handle_t);

typedef indy_error_t (*cb_deleteWalletCustom)(char*, char*, char*);
extern indy_error_t deleteWalletCustomCB(char*, char*, char*);

typedef indy_error_t (*cb_addRecordWallet)(indy_handle_t, char*, char*, indy_u8_t*, indy_u32_t, char*);
extern indy_error_t addRecordWalletCB(indy_handle_t, char*, char*, indy_u8_t*, indy_u32_t, char*);

typedef indy_error_t (*cb_updateRecordValueWallet)(indy_handle_t, char*, char*, indy_u8_t*, indy_u32_t);
extern indy_error_t updateRecordValueWalletCB(indy_handle_t, char*, char*, indy_u8_t*, indy_u32_t);

typedef indy_error_t (*cb_updateRecordTagsWallet)(indy_handle_t, char*, char*, char*);
extern indy_error_t updateRecordTagsWalletCB(indy_handle_t, char*, char*, char*);

typedef indy_error_t (*cb_addRecordTagsWallet)(indy_handle_t, char*, char*, char*);
extern indy_error_t addRecordTagsWalletCB(indy_handle_t, char*, char*, char*);

typedef indy_error_t (*cb_deleteRecordTagsWallet)(indy_handle_t, char*, char*, char*);
extern indy_error_t deleteRecordTagsWalletCB(indy_handle_t, char*, char*, char*);

typedef indy_error_t (*cb_deleteRecordWallet)(indy_handle_t, char*, char*);
extern indy_error_t deleteRecordWalletCB(indy_handle_t, char*, char*);

typedef indy_error_t (*cb_getRecordHandleWallet)(indy_handle_t, char*, char*, char*, int32_t*);
extern indy_error_t getRecordHandleWalletCB(indy_handle_t, char*, char*, char*, int32_t*);

typedef indy_error_t (*cb_getRecordIdWallet)(indy_handle_t, indy_handle_t, char**);
extern indy_error_t getRecordIdWalletCB(indy_handle_t, indy_handle_t, char**);

typedef indy_error_t (*cb_getRecordTypeWallet)(indy_handle_t, indy_handle_t, char**);
extern indy_error_t getRecordTypeWalletCB(indy_handle_t, indy_handle_t, char**);

typedef indy_error_t (*cb_getRecordValueWallet)(indy_handle_t, indy_handle_t, indy_u8_t**, indy_u32_t*);
extern indy_error_t getRecordValueWalletCB(indy_handle_t, indy_handle_t, indy_u8_t**, indy_u32_t*);

typedef indy_error_t (*cb_getRecordTagsWallet)(indy_handle_t, indy_handle_t, char**);
extern indy_error_t getRecordTagsWalletCB(indy_handle_t, indy_handle_t, char**);

typedef indy_error_t (*cb_freeRecordWallet)(indy_handle_t, indy_handle_t);
extern indy_error_t freeRecordWalletCB(indy_handle_t, indy_handle_t);

typedef indy_error_t (*cb_getStorageMetadataWallet)(indy_handle_t, char**, indy_handle_t*);
extern indy_error_t getStorageMetadataWalletCB(indy_handle_t, char**, indy_handle_t*);

typedef indy_error_t (*cb_setStorageMetadataWallet)(indy_handle_t, char*);
extern indy_error_t setStorageMetadataWalletCB(indy_handle_t, char*);

typedef indy_error_t (*cb_freeStorageMetadataWallet)(indy_handle_t, indy_handle_t);
extern indy_error_t freeStorageMetadataWalletCB(indy_handle_t, indy_handle_t);

typedef indy_error_t (*cb_openSearchWallet)(indy_handle_t, char*, char*, char*, int32_t*);
extern indy_error_t openSearchWalletCB(indy_handle_t, char*, char*, char*, int32_t*);

typedef indy_error_t (*cb_openSearchAllWallet)(indy_handle_t, indy_handle_t*);
extern indy_error_t openSearchAllWalletCB(indy_handle_t, indy_handle_t*);

typedef indy_error_t (*cb_getSearchTotalCountWallet)(indy_handle_t, indy_handle_t, indy_u32_t*);
extern indy_error_t getSearchTotalCountWalletCB(indy_handle_t, indy_handle_t, indy_u32_t*);

typedef indy_error_t (*cb_fetchSearchNextRecordsWallet)(indy_handle_t, indy_handle_t, indy_handle_t*);
extern indy_error_t fetchSearchNextRecordsWalletCB(indy_handle_t, indy_handle_t, indy_handle_t*);

typedef indy_error_t (*cb_freeSearchWallet)(indy_handle_t, indy_handle_t);
extern indy_error_t freeSearchWalletCB(indy_handle_t, indy_handle_t);
*/
import "C"

import (
	"github.com/joyride9999/IndySdkGoBindings/indyUtils"
	"errors"
	"unsafe"
)

//export createWalletCB
func createWalletCB(commandHandle C.indy_handle_t, indyError C.indy_error_t) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// CreateWallet creates a new secure wallet with the given unique id
func CreateWallet(config, credential unsafe.Pointer) chan indyUtils.IndyResult {
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
	 Creates a new secure wallet with the given unique name.

	    :param config: Wallet configuration json.
	     {
	       "id": string, Identifier of the wallet.
	             Configured storage uses this identifier to lookup exact wallet data placement.
	       "storage_type": optional<string>, Type of the wallet storage. Defaults to 'default'.
	                      'Default' storage type allows to store wallet data in the local file.
	                      Custom storage types can be registered with indy_register_wallet_storage call.
	       "storage_config": optional<object>, Storage configuration json. Storage type defines set of supported keys.
	                         Can be optional if storage supports default configuration.
	                          For 'default' storage type configuration is:
	       {
	         "path": optional<string>, Path to the directory with wallet files.
	                 Defaults to $HOME/.indy_client/wallet.
	                 Wallet will be stored in the file {path}/{id}/sqlite.db
	       }
	     }
	    :param credentials: Wallet credentials json
	     {
	       "key": string, Key or passphrase used for wallet key derivation.
	                      Look to key_derivation_method param for information about supported key derivation methods.
	       "storage_credentials": optional<object> Credentials for wallet storage. Storage type defines set of supported keys.
	                              Can be optional if storage supports default configuration.
	                               For 'default' storage type should be empty.
	       "key_derivation_method": optional<string> Algorithm to use for wallet key derivation:
	                                ARGON2I_MOD - derive secured wallet master key (used by default)
	                                ARGON2I_INT - derive secured wallet master key (less secured but faster)
	                                RAW - raw wallet key master provided (skip derivation).
	                                      RAW keys can be generated with generate_wallet_key call
	     }
	    :return: Error code

	*/
	// Call indy_create_wallet
	res := C.indy_create_wallet(commandHandle,
		(*C.char)(config),
		(*C.char)(credential),
		(C.cb_createWallet)(unsafe.Pointer(C.createWalletCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export closeWalletCB
func closeWalletCB(commandHandle C.indy_handle_t, indyError C.indy_error_t) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// CloseWallet  Closes opened wallet and frees allocated resources.
func CloseWallet(wh int) chan indyUtils.IndyResult {
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
	  wallet_handle: wallet handle returned by indy_open_wallet.
	    ///
	    /// #Returns
	    /// Error code
	*/
	// Call indy_create_wallet
	res := C.indy_close_wallet(commandHandle,
		(C.indy_handle_t)(wh),
		(C.cb_closeWallet)(unsafe.Pointer(C.closeWalletCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export openWalletCB
func openWalletCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, retHandle C.indy_handle_t) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil, Results: []interface{}{int(retHandle)}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// OpenWallet opens an existing indy wallet
func OpenWallet(config, credential unsafe.Pointer) chan indyUtils.IndyResult {
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
	   Opens the wallet with specific name.
	   Wallet with corresponded name must be previously created with indy_create_wallet method.
	   It is impossible to open wallet with the same name more than once.

	   :param name: Name of the wallet.
	   :param config: Wallet configuration json.
	   {
	      "id": string, Identifier of the wallet.
	            Configured storage uses this identifier to lookup exact wallet data placement.
	      "storage_type": optional<string>, Type of the wallet storage. Defaults to 'default'.
	                      'Default' storage type allows to store wallet data in the local file.
	                      Custom storage types can be registered with indy_register_wallet_storage call.
	      "storage_config": optional<object>, Storage configuration json. Storage type defines set of supported keys.
	                        Can be optional if storage supports default configuration.
	                         For 'default' storage type configuration is:
	          {
	             "path": optional<string>, Path to the directory with wallet files.
	                     Defaults to $HOME/.indy_client/wallet.
	                     Wallet will be stored in the file {path}/{id}/sqlite.db
	          }

	   }
	   :param credentials: Wallet credentials json
	   {
	      "key": string, Key or passphrase used for wallet key derivation.
	                     Look to key_derivation_method param for information about supported key derivation methods.
	      "rekey": optional<string>, If present, then wallet master key will be rotated to a new one.
	      "storage_credentials": optional<object> Credentials for wallet storage. Storage type defines set of supported keys.
	                             Can be optional if storage supports default configuration.
	                             For 'default' storage type should be empty.
	      "key_derivation_method": optional<string> Algorithm to use for wallet key derivation:
	                              ARGON2I_MOD - derive secured wallet master key (used by default)
	                              ARGON2I_INT - derive secured wallet master key (less secured but faster)
	                              RAW - raw wallet master key provided (skip derivation)
	      "rekey_derivation_method": optional<string> algorithm to use for master rekey derivation:
	                              ARGON2I_MOD - derive secured wallet master rekey (used by default)
	                              ARGON2I_INT - derive secured wallet master rekey (less secured but faster)
	                              RAW - raw wallet rekey master provided (skip derivation).
	                                    RAW keys can be generated with generate_wallet_key call
	   }
	   :return: Handle to opened wallet to use in methods that require wallet access.
	*/

	// Call indy_open_wallet
	res := C.indy_open_wallet(commandHandle,
		(*C.char)(config),
		(*C.char)(credential),
		(C.cb_openWallet)(unsafe.Pointer(C.openWalletCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() {
			indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)})
		}()
		return future
	}

	return future
}

//export deleteWalletCB
func deleteWalletCB(commandHandle C.indy_handle_t, indyError C.indy_error_t) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// DeleteWallet deletes an existing indy wallet
func DeleteWallet(config, credentials unsafe.Pointer) chan indyUtils.IndyResult {
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
		Deletes created wallet.

		:param config: Wallet configuration json.
		{
		  "id": string, Identifier of the wallet.
		        Configured storage uses this identifier to lookup exact wallet data placement.
		  "storage_type": optional<string>, Type of the wallet storage. Defaults to 'default'.
		                 'Default' storage type allows to store wallet data in the local file.
		                 Custom storage types can be registered with indy_register_wallet_storage call.
		  "storage_config": optional<object>, Storage configuration json. Storage type defines set of supported keys.
		                    Can be optional if storage supports default configuration.
		                    For 'default' storage type configuration is:
		  {
		    "path": optional<string>, Path to the directory with wallet files.
		            Defaults to $HOME/.indy_client/wallet.
		            Wallet will be stored in the file {path}/{id}/sqlite.db
		  }
		}
		:param credentials: Wallet credentials json
		{
		  "key": string, Key or passphrase used for wallet key derivation.
		                 Look to key_derivation_method param for information about supported key derivation methods.
		  "storage_credentials": optional<object> Credentials for wallet storage. Storage type defines set of supported keys.
		                         Can be optional if storage supports default configuration.
		                         For 'default' storage type should be empty.
		  "key_derivation_method": optional<string> Algorithm to use for wallet key derivation:
		                            ARGON2I_MOD - derive secured wallet master key (used by default)
		                            ARGON2I_INT - derive secured wallet master key (less secured but faster)
		                            RAW - raw wallet key master provided (skip derivation).
		                               RAW keys can be generated with indy_generate_wallet_key call
		}

		:return: Error code
	*/

	// Call indy_delete_wallet
	res := C.indy_delete_wallet(commandHandle,
		(*C.char)(config),
		(*C.char)(credentials),
		(C.cb_deleteWallet)(unsafe.Pointer(C.deleteWalletCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() {
			indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)})
		}()
		return future
	}

	return future
}

//export generateWalletKeyCB
func generateWalletKeyCB(commandHandle C.indy_handle_t, indyError C.indy_error_t) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// GenerateWalletKey generates wallet master key
func GenerateWalletKey(config unsafe.Pointer) chan indyUtils.IndyResult {
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
			Generate wallet master key.
		    Returned key is compatible with "RAW" key derivation method.
		    It allows to avoid expensive key derivation for use cases when wallet keys can be stored in a secure enclave.

			:param config: (optional) key configuration json.
			{
			  "seed": string, (optional) Seed that allows deterministic key creation (if not set random one will be created).
			                             Can be UTF-8, base64 or hex string.
			}

			:return: Error code
	*/

	// Call indy_generate_wallet_key
	res := C.indy_generate_wallet_key(commandHandle,
		(*C.char)(config),
		(C.cb_generateWalletKey)(unsafe.Pointer(C.generateWalletKeyCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() {
			indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)})
		}()
		return future
	}

	return future
}

//export exportWalletCB
func exportWalletCB(commandHandle C.indy_handle_t, indyError C.indy_error_t) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// ExportWallet exports opened wallet
func ExportWallet(wh int, config unsafe.Pointer) chan indyUtils.IndyResult {
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
		:param wallet_handle: wallet handle returned by indy_open_wallet
		:param export_config: JSON containing settings for input operation.
		  {
		    "path": <string>, Path of the file that contains exported wallet content
		    "key": <string>, Key or passphrase used for wallet export key derivation.
		                    Look to key_derivation_method param for information about supported key derivation methods.
		    "key_derivation_method": optional<string> Algorithm to use for wallet export key derivation:
		                             ARGON2I_MOD - derive secured export key (used by default)
		                             ARGON2I_INT - derive secured export key (less secured but faster)
		                             RAW - raw export key provided (skip derivation).
		                               RAW keys can be generated with indy_generate_wallet_key call
		  }

		:return: Error code
	*/

	// Call indy_export_wallet
	res := C.indy_export_wallet(commandHandle,
		(C.indy_handle_t)(wh),
		(*C.char)(config),
		(C.cb_exportWallet)(unsafe.Pointer(C.exportWalletCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() {
			indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)})
		}()
		return future
	}

	return future
}

//export importWalletCB
func importWalletCB(commandHandle C.indy_handle_t, indyError C.indy_error_t) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// ImportWallet imports opened wallet
func ImportWallet(config, credentials, importConfig unsafe.Pointer) chan indyUtils.IndyResult {
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
		Creates a new secure wallet and then imports its content
		according to fields provided in import_config
		This can be seen as an indy_create_wallet call with additional content import

		:param config: Wallet configuration json.
		{
		  "id": string, Identifier of the wallet.
		        Configured storage uses this identifier to lookup exact wallet data placement.
		  "storage_type": optional<string>, Type of the wallet storage. Defaults to 'default'.
		                 'Default' storage type allows to store wallet data in the local file.
		                 Custom storage types can be registered with indy_register_wallet_storage call.
		  "storage_config": optional<object>, Storage configuration json. Storage type defines set of supported keys.
		                    Can be optional if storage supports default configuration.
		                    For 'default' storage type configuration is:
		  {
		    "path": optional<string>, Path to the directory with wallet files.
		            Defaults to $HOME/.indy_client/wallet.
		            Wallet will be stored in the file {path}/{id}/sqlite.db
		  }
		}
		:param credentials: Wallet credentials json
		{
		  "key": string, Key or passphrase used for wallet key derivation.
		                 Look to key_derivation_method param for information about supported key derivation methods.
		  "storage_credentials": optional<object> Credentials for wallet storage. Storage type defines set of supported keys.
		                         Can be optional if storage supports default configuration.
		                         For 'default' storage type should be empty.
		  "key_derivation_method": optional<string> Algorithm to use for wallet key derivation:
		                            ARGON2I_MOD - derive secured wallet master key (used by default)
		                            ARGON2I_INT - derive secured wallet master key (less secured but faster)
		                            RAW - raw wallet key master provided (skip derivation).
		                               RAW keys can be generated with indy_generate_wallet_key call
		}
		:param import_config: Import settings json.
		{
		  "path": <string>, path of the file that contains exported wallet content
		  "key": <string>, key used for export of the wallet
		}

		:return: Error code
	*/

	// Call indy_import_wallet
	res := C.indy_import_wallet(commandHandle,
		(*C.char)(config),
		(*C.char)(credentials),
		(*C.char)(importConfig),
		(C.cb_importWallet)(unsafe.Pointer(C.importWalletCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() {
			indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)})
		}()
		return future
	}

	return future
}

//export createWalletCustomCB
func createWalletCustomCB(storageName *C.char, storageConfig *C.char, credentialsJson *C.char, metadata *C.char) (err C.indy_error_t) {
	defer func() {
		if r := recover(); r != nil {
			err = (C.indy_error_t)(210) // WalletStorageError: Storage error occurred during wallet operation
			return
		}
	}()

	indyErr, errCreate := walletStorageGlobal.Create(C.GoString(storageName), C.GoString(storageConfig), C.GoString(credentialsJson), C.GoString(metadata))
	if errCreate != nil {
		return (C.indy_error_t)(indyErr)
	}

	return (C.indy_error_t)(0) // Success
}

//export openWalletCustomCB
func openWalletCustomCB(walletID *C.char, storageConfig *C.char, credentialsJson *C.char, handle *C.indy_handle_t) (err C.indy_error_t) {
	defer func() {
		if r := recover(); r != nil {
			err = (C.indy_error_t)(210) // WalletStorageError: Storage error occurred during wallet operation
			return
		}
	}()

	storageHandle, indyCode, errOH := walletStorageGlobal.Open(C.GoString(walletID), C.GoString(storageConfig), C.GoString(credentialsJson))
	if errOH != nil {
		return (C.indy_error_t)(indyCode)
	}

	*handle = C.indy_handle_t(storageHandle)
	return (C.indy_error_t)(0) // Success
}

//export closeWalletCustomCB
func closeWalletCustomCB(storageHandle C.indy_handle_t) (err C.indy_error_t) {
	defer func() {
		if r := recover(); r != nil {
			err = (C.indy_error_t)(210) // WalletStorageError: Storage error occurred during wallet operation
			return
		}
	}()

	walletStorageGlobal.Close(int(storageHandle))
	return (C.indy_error_t)(0) // Success
}

//export deleteWalletCustomCB
func deleteWalletCustomCB(storageName *C.char, storageConfig *C.char, credentialsJson *C.char) (err C.indy_error_t) {
	defer func() {
		if r := recover(); r != nil {
			err = (C.indy_error_t)(210) // WalletStorageError: Storage error occurred during wallet operation
			return
		}
	}()

	indyError, errD := walletStorageGlobal.Delete(C.GoString(storageName), C.GoString(storageConfig), C.GoString(credentialsJson))
	if errD != nil {
		return (C.indy_error_t)(indyError) // error
	}

	return (C.indy_error_t)(0) // Success
}

//export addRecordWalletCB
func addRecordWalletCB(storageHandle C.indy_handle_t, recordType *C.char, recordId *C.char, recordValue *C.indy_u8_t, valueLen C.indy_u32_t, tagsJson *C.char) (err C.indy_error_t) {
	defer func() {
		if r := recover(); r != nil {
			err = (C.indy_error_t)(210) // WalletStorageError: Storage error occurred during wallet operation
			return
		}
	}()

	buffValue := C.GoBytes(unsafe.Pointer(recordValue), C.int(valueLen))
	indyError, errAddRecord := walletStorageGlobal.AddRecord(int(storageHandle), C.GoString(recordType), C.GoString(recordId), buffValue, C.GoString(tagsJson))
	if errAddRecord != nil {
		return (C.indy_error_t)(indyError)
	}

	return (C.indy_error_t)(0) // Success
}

//export updateRecordValueWalletCB
func updateRecordValueWalletCB(storageHandle C.indy_handle_t, recordType *C.char, recordId *C.char, recordValue *C.indy_u8_t, recordValueLen C.indy_u32_t) (err C.indy_error_t) {
	defer func() {
		if r := recover(); r != nil {
			err = (C.indy_error_t)(210) // WalletStorageError: Storage error occurred during wallet operation
			return
		}
	}()

	buffValue := C.GoBytes(unsafe.Pointer(recordValue), C.int(recordValueLen))
	indyError, errUpdateValue := walletStorageGlobal.UpdateRecordValue(int(storageHandle), C.GoString(recordType), C.GoString(recordId), buffValue)
	if errUpdateValue != nil {
		return (C.indy_error_t)(indyError)
	}

	return (C.indy_error_t)(0) // Success
}

//export updateRecordTagsWalletCB
func updateRecordTagsWalletCB(storageHandle C.indy_handle_t, recordType *C.char, recordId *C.char, tagsJson *C.char) (err C.indy_error_t) {
	defer func() {
		if r := recover(); r != nil {
			err = (C.indy_error_t)(210) // WalletStorageError: Storage error occurred during wallet operation
			return
		}
	}()

	indyError, errUpdateTags := walletStorageGlobal.UpdateRecordTags(int(storageHandle), C.GoString(recordType), C.GoString(recordId), C.GoString(tagsJson))
	if errUpdateTags != nil {
		return (C.indy_error_t)(indyError) //error
	}

	return (C.indy_error_t)(0) // Success
}

//export addRecordTagsWalletCB
func addRecordTagsWalletCB(storageHandle C.indy_handle_t, recordType *C.char, recordId *C.char, tagsJson *C.char) (err C.indy_error_t) {
	defer func() {
		if r := recover(); r != nil {
			err = (C.indy_error_t)(210) // WalletStorageError: Storage error occurred during wallet operation
			return
		}
	}()

	indyError, errAddRecordTags := walletStorageGlobal.AddRecordTags(int(storageHandle), C.GoString(recordType), C.GoString(recordId), C.GoString(tagsJson))
	if errAddRecordTags != nil {
		return (C.indy_error_t)(indyError)
	}

	return (C.indy_error_t)(0) // Success
}

//export deleteRecordTagsWalletCB
func deleteRecordTagsWalletCB(storageHandle C.indy_handle_t, recordType *C.char, recordId *C.char, tagsName *C.char) (err C.indy_error_t) {
	defer func() {
		if r := recover(); r != nil {
			err = (C.indy_error_t)(210) // WalletStorageError: Storage error occurred during wallet operation
			return
		}
	}()

	indyError, errDeleteTags := walletStorageGlobal.DeleteRecordTags(int(storageHandle), C.GoString(recordType), C.GoString(recordId), C.GoString(tagsName))
	if errDeleteTags != nil {
		return (C.indy_error_t)(indyError) //error
	}

	return (C.indy_error_t)(0) // Success
}

//export deleteRecordWalletCB
func deleteRecordWalletCB(storageHandle C.indy_handle_t, recordType *C.char, recordId *C.char) (err C.indy_error_t) {
	defer func() {
		if r := recover(); r != nil {
			err = (C.indy_error_t)(210) // WalletStorageError: Storage error occurred during wallet operation
			return
		}
	}()

	indyError, errD := walletStorageGlobal.DeleteRecord(int(storageHandle), C.GoString(recordType), C.GoString(recordId))
	if errD != nil {
		return (C.indy_error_t)(indyError) // error
	}
	return (C.indy_error_t)(0) // Success
}

//export getRecordHandleWalletCB
func getRecordHandleWalletCB(storageHandle C.indy_handle_t, recordType *C.char, recordId *C.char, optionsJson *C.char, recordHandle *C.int32_t) (err C.indy_error_t) {
	defer func() {
		if r := recover(); r != nil {
			err = (C.indy_error_t)(210) // WalletStorageError: Storage error occurred during wallet operation
			return
		}
	}()

	itemId, indyError, errR := walletStorageGlobal.GetRecordHandle(int(storageHandle), C.GoString(recordType), C.GoString(recordId), C.GoString(optionsJson))
	if errR != nil {
		return (C.indy_error_t)(indyError) // error
	}

	indyRecordHandle := C.indy_handle_t(itemId)
	*recordHandle = indyRecordHandle
	return (C.indy_error_t)(indyError)

}

//export getRecordIdWalletCB
func getRecordIdWalletCB(storageHandle C.indy_handle_t, recordHandle C.indy_handle_t, recordID **C.char) (err C.indy_error_t) {
	// returns the item.name field (not to be confused with rowid (item.id)
	defer func() {
		if r := recover(); r != nil {
			err = (C.indy_error_t)(210) // WalletStorageError: Storage error occurred during wallet operation
			return
		}
	}()

	sRecordID, indyError, errR := walletStorageGlobal.GetRecordId(int(storageHandle), int(recordHandle))
	if errR != nil {
		return (C.indy_error_t)(indyError) // error
	}

	*recordID = (*C.char)(sRecordID)
	return (C.indy_error_t)(indyError) // Success
}

//export getRecordTypeWalletCB
func getRecordTypeWalletCB(storageHandle C.indy_handle_t, recordHandle C.indy_handle_t, recordType **C.char) (err C.indy_error_t) {
	defer func() {
		if r := recover(); r != nil {
			err = (C.indy_error_t)(210) // WalletStorageError: Storage error occurred during wallet operation
			return
		}
	}()

	sRecordType, indyError, errT := walletStorageGlobal.GetRecordType(int(storageHandle), int(recordHandle))
	if errT != nil {
		return (C.indy_error_t)(indyError) //error
	}

	*recordType = (*C.char)(sRecordType)
	return (C.indy_error_t)(0) // Success
}

//export getRecordValueWalletCB
func getRecordValueWalletCB(handle C.indy_handle_t, recordHandle C.indy_handle_t, recordValue **C.indy_u8_t, recordValueLen *C.indy_u32_t) (err C.indy_error_t) {
	defer func() {
		if r := recover(); r != nil {
			err = (C.indy_error_t)(210) // WalletStorageError: Storage error occurred during wallet operation
			return
		}
	}()

	bufValue, indyErr, errGet := walletStorageGlobal.GetRecordValue(int(handle), int(recordHandle))
	if errGet != nil {
		return (C.indy_error_t)(indyErr)
	}

	*recordValue = (*C.indy_u8_t)(bufValue.Value)

	bufferLength := uint32(bufValue.Len)
	cBufferSize := (C.uint32_t)(bufferLength)
	*recordValueLen = cBufferSize

	return (C.indy_error_t)(0) // Success
}

//export getRecordTagsWalletCB
func getRecordTagsWalletCB(storageHandle C.indy_handle_t, recordHandle C.indy_handle_t, tagsJson **C.char) (err C.indy_error_t) {
	defer func() {
		if r := recover(); r != nil {
			err = (C.indy_error_t)(210) // WalletStorageError: Storage error occurred during wallet operation
			return
		}
	}()

	recordTags, indyError, errR := walletStorageGlobal.GetRecordTags(int(storageHandle), int(recordHandle))
	if errR != nil {
		return (C.indy_error_t)(indyError) // error
	}

	*tagsJson = (*C.char)(recordTags)
	return (C.indy_error_t)(indyError) // Success
}

//export freeRecordWalletCB
func freeRecordWalletCB(handle C.indy_handle_t, recordHandle C.indy_handle_t) (err C.indy_error_t) {
	defer func() {
		if r := recover(); r != nil {
			err = (C.indy_error_t)(210) // WalletStorageError: Storage error occurred during wallet operation
			return
		}
	}()

	walletStorageGlobal.FreeRecord(int(handle), int(recordHandle))

	return (C.indy_error_t)(0) // Success
}

//export getStorageMetadataWalletCB
func getStorageMetadataWalletCB(storageHandle C.indy_handle_t, metadata **C.char, metadataHandle *C.indy_handle_t) (err C.indy_error_t) {
	defer func() {
		if r := recover(); r != nil {
			err = (C.indy_error_t)(210) // WalletStorageError: Storage error occurred during wallet operation
			return
		}
	}()

	sMetadata, iMetadata, indyError, errG := walletStorageGlobal.GetStorageMetadata(int(storageHandle))
	if errG != nil {
		return (C.indy_error_t)(indyError) // error
	}

	*metadata = (*C.char)(sMetadata)
	*metadataHandle = C.indy_handle_t(iMetadata)

	return (C.indy_error_t)(0) // Success
}

//export setStorageMetadataWalletCB
func setStorageMetadataWalletCB(storageHandle C.indy_handle_t, metadata *C.char) (err C.indy_error_t) {
	defer func() {
		if r := recover(); r != nil {
			err = (C.indy_error_t)(210) // WalletStorageError: Storage error occurred during wallet operation
			return
		}
	}()

	indyErr, errS := walletStorageGlobal.SetStorageMetadata(int(storageHandle), C.GoString(metadata))
	if errS != nil {
		return (C.indy_error_t)(indyErr) // error
	}
	return (C.indy_error_t)(0)
}

//export freeStorageMetadataWalletCB
func freeStorageMetadataWalletCB(storageHandle C.indy_handle_t, metadataHandle C.indy_handle_t) (err C.indy_error_t) {
	defer func() {
		if r := recover(); r != nil {
			err = (C.indy_error_t)(210) // WalletStorageError: Storage error occurred during wallet operation
			return
		}
	}()

	walletStorageGlobal.FreeStorageMetadata(int(storageHandle), int(metadataHandle))

	return (C.indy_error_t)(0) // Success
}

//export openSearchWalletCB
func openSearchWalletCB(storageHandle C.indy_handle_t, recordType *C.char, queryJson *C.char, optionsJson *C.char, searchHandle *C.int32_t) (err C.indy_error_t) {
	defer func() {
		if r := recover(); r != nil {
			err = (C.indy_error_t)(210) // WalletStorageError: Storage error occurred during wallet operation
			return
		}
	}()

	openSearch, indyCode, errO := walletStorageGlobal.OpenSearch(int(storageHandle), C.GoString(recordType), C.GoString(queryJson), C.GoString(optionsJson))
	if errO != nil {
		return (C.indy_error_t)(indyCode) // error
	}

	*searchHandle = C.int(openSearch)
	return (C.indy_error_t)(0) // Success
}

//export openSearchAllWalletCB
func openSearchAllWalletCB(storageHandle C.indy_handle_t, searchHandle *C.indy_handle_t) (err C.indy_error_t) {
	defer func() {
		if r := recover(); r != nil {
			err = (C.indy_error_t)(210) // WalletStorageError: Storage error occurred during wallet operation
			return
		}
	}()

	searchAllHandle, indyErr, errO := walletStorageGlobal.OpenSearchAll(int(storageHandle))
	if errO != nil {
		return (C.indy_error_t)(indyErr)
	}
	*searchHandle = C.int(searchAllHandle)

	return (C.indy_error_t)(0) // Success
}

//export getSearchTotalCountWalletCB
func getSearchTotalCountWalletCB(storageHandle C.indy_handle_t, searchHandle C.indy_handle_t, totalCount *C.indy_u32_t) (err C.indy_error_t) {
	defer func() {
		if r := recover(); r != nil {
			err = (C.indy_error_t)(210) // WalletStorageError: Storage error occurred during wallet operation
			return
		}
	}()

	searchCount, indyErr, errC := walletStorageGlobal.GetSearchTotalCount(int(storageHandle), int(searchHandle))
	if errC != nil {
		return (C.indy_error_t)(indyErr) // error
	}
	*totalCount = C.indy_u32_t(searchCount)

	return (C.indy_error_t)(0) // Success
}

//export fetchSearchNextRecordsWalletCB
func fetchSearchNextRecordsWalletCB(storageHandle C.indy_handle_t, searchHandle C.indy_handle_t, recordId *C.indy_handle_t) (err C.indy_error_t) {
	defer func() {
		if r := recover(); r != nil {
			err = (C.indy_error_t)(210) // WalletStorageError: Storage error occurred during wallet operation
			return
		}
	}()

	itemId, indyCode, errF := walletStorageGlobal.FetchSearchNext(int(storageHandle), int(searchHandle))
	if errF != nil {
		return (C.indy_error_t)(indyCode) // error
	}

	*recordId = C.int(itemId)
	return (C.indy_error_t)(0) // succcess
}

//export freeSearchWalletCB
func freeSearchWalletCB(storageHandle C.indy_handle_t, searchHandle C.indy_handle_t) (err C.indy_error_t) {
	defer func() {
		if r := recover(); r != nil {
			err = (C.indy_error_t)(210) // WalletStorageError: Storage error occurred during wallet operation
			return
		}
	}()

	walletStorageGlobal.FreeSearch(int(storageHandle), int(searchHandle))

	return (C.indy_error_t)(0) // Success
}

//export registerWalletStorageCB
func registerWalletStorageCB(commandHandle C.indy_handle_t, indyError C.indy_error_t) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

var walletStorageGlobal IWalletStorage

func RegisterWalletStorage(storageType unsafe.Pointer, storage IWalletStorage) chan indyUtils.IndyResult {
	walletStorageGlobal = storage

	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)
	res := C.indy_register_wallet_storage(commandHandle,
		(*C.char)(storageType),
		(C.cb_createWalletCustom)(unsafe.Pointer(C.createWalletCustomCB)),
		(C.cb_openWalletCustom)(unsafe.Pointer(C.openWalletCustomCB)),
		(C.cb_closeWalletCustom)(unsafe.Pointer(C.closeWalletCustomCB)),
		(C.cb_deleteWalletCustom)(unsafe.Pointer(C.deleteWalletCustomCB)),
		(C.cb_addRecordWallet)(unsafe.Pointer(C.addRecordWalletCB)),
		(C.cb_updateRecordValueWallet)(unsafe.Pointer(C.updateRecordValueWalletCB)),
		(C.cb_updateRecordTagsWallet)(unsafe.Pointer(C.updateRecordTagsWalletCB)),
		(C.cb_addRecordTagsWallet)(unsafe.Pointer(C.addRecordTagsWalletCB)),
		(C.cb_deleteRecordTagsWallet)(unsafe.Pointer(C.deleteRecordTagsWalletCB)),
		(C.cb_deleteRecordWallet)(unsafe.Pointer(C.deleteRecordWalletCB)),
		(C.cb_getRecordHandleWallet)(unsafe.Pointer(C.getRecordHandleWalletCB)),
		(C.cb_getRecordIdWallet)(unsafe.Pointer(C.getRecordIdWalletCB)),
		(C.cb_getRecordTypeWallet)(unsafe.Pointer(C.getRecordTypeWalletCB)),
		(C.cb_getRecordValueWallet)(unsafe.Pointer(C.getRecordValueWalletCB)),
		(C.cb_getRecordTagsWallet)(unsafe.Pointer(C.getRecordTagsWalletCB)),
		(C.cb_freeRecordWallet)(unsafe.Pointer(C.freeRecordWalletCB)),
		(C.cb_getStorageMetadataWallet)(unsafe.Pointer(C.getStorageMetadataWalletCB)),
		(C.cb_setStorageMetadataWallet)(unsafe.Pointer(C.setStorageMetadataWalletCB)),
		(C.cb_freeStorageMetadataWallet)(unsafe.Pointer(C.freeStorageMetadataWalletCB)),
		(C.cb_openSearchWallet)(unsafe.Pointer(C.openSearchWalletCB)),
		(C.cb_openSearchAllWallet)(unsafe.Pointer(C.openSearchAllWalletCB)),
		(C.cb_getSearchTotalCountWallet)(unsafe.Pointer(C.getSearchTotalCountWalletCB)),
		(C.cb_fetchSearchNextRecordsWallet)(unsafe.Pointer(C.fetchSearchNextRecordsWalletCB)),
		(C.cb_freeSearchWallet)(unsafe.Pointer(C.freeSearchWalletCB)),
		(C.cb_registerWalletStorage)(unsafe.Pointer(C.registerWalletStorageCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() {
			indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)})
		}()
		return future
	}
	return future
}
