/*
// ******************************************************************
// Purpose: exported public functions that handles did functions
// from libindy
// Author:  alexandru.leonte@siemens.com
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
	"github.com/joyride9999/IndySdkGoBindings/did"
	"encoding/json"
	"errors"
	"unsafe"
)

// CreateAndStoreDID creates and DID with keys ... nothing is written to blockchain
// returns did, verkey, error
func CreateAndStoreDID(walletHandle int, seed string) (string, string, error) {

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

	didcfg, err := json.Marshal(didjs)
	if err != nil {
		return "", "", errors.New("cant read json")
	}

	upDid := unsafe.Pointer(C.CString(string(didcfg)))
	defer C.free(upDid)

	channel := did.CreateAndStoreMyDid(walletHandle, upDid)
	result := <-channel
	if result.Error != nil {
		return "", "", result.Error
	}
	return result.Results[0].(string), result.Results[1].(string), result.Error
}

// ReplaceKeyStart generates temporary key for an existing DID.
func ReplaceKeyStart(walletHandle int, Did string, identityJson string) (string, error) {

	upDid := unsafe.Pointer(C.CString(Did))
	defer C.free(upDid)
	upIdentityJson := unsafe.Pointer(C.CString(identityJson))
	defer C.free(upIdentityJson)

	channel := did.ReplaceKeyStart(walletHandle, upDid, upIdentityJson)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// ReplaceKeyApply applies temporary keys as main for existing DID
func ReplaceKeyApply(walletHandle int, Did string) error {
	upDid := unsafe.Pointer(C.CString(Did))
	defer C.free(upDid)
	channel := did.ReplaceKeyApply(walletHandle, upDid)
	result := <-channel
	return result.Error
}

// StoreTheirDid saves DID for a pairwise connection in a secured wallet to verify transaction.
func StoreTheirDid(walletHandle int, identityJson string) error {
	upIdentityJson := unsafe.Pointer(C.CString(identityJson))
	defer C.free(upIdentityJson)
	channel := did.StoreTheirDid(walletHandle, upIdentityJson)
	result := <-channel
	return result.Error
}

// KeyForDid returns ver key for DID.
func KeyForDid(poolHandle int, walletHandle int, Did string) (string, error) {
	upDid := unsafe.Pointer(C.CString(Did))
	defer C.free(upDid)
	channel := did.KeyForDid(poolHandle, walletHandle, upDid)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// KeyForLocalDID gets the key for the local DID.
// returns key, error
func KeyForLocalDID(walletHandle int, Did string) (string, error) {
	upDid := unsafe.Pointer(C.CString(Did))
	defer C.free(upDid)
	channel := did.KeyForLocalDid(walletHandle, upDid)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// SetEndPointForDid set/replaces endpoint information for the given DID
func SetEndPointForDid(walletHandle int, Did string, address string, transportKey string) error {
	upDid := unsafe.Pointer(C.CString(Did))
	defer C.free(upDid)
	upAddress := unsafe.Pointer(C.CString(address))
	defer C.free(upAddress)
	upTransportKey := unsafe.Pointer(C.CString(transportKey))
	defer C.free(upTransportKey)
	channel := did.SetEndPointForDid(walletHandle, upDid, upAddress, upTransportKey)
	result := <-channel
	return result.Error
}

// GetEndPointForDid returns endpoint information for the given DID
func GetEndPointForDid(walletHandle int, poolHandle int, Did string) (string, string, error) {
	upDid := unsafe.Pointer(C.CString(Did))
	defer C.free(upDid)
	channel := did.GetEndPointForDid(walletHandle, poolHandle, upDid)
	result := <-channel
	if result.Error != nil {
		return "", "", result.Error
	}
	return result.Results[0].(string), result.Results[1].(string), result.Error
}

// SetDidMetadata saves/replaces meta information for the given DID.
func SetDidMetadata(walletHandle int, Did string, metadata string) error {
	upDid := unsafe.Pointer(C.CString(Did))
	defer C.free(upDid)
	upMetadata := unsafe.Pointer(C.CString(metadata))
	defer C.free(upMetadata)
	channel := did.SetDidMetadata(walletHandle, upDid, upMetadata)
	result := <-channel
	return result.Error
}

// GetDidMetadata retrieves meta information for the given DID.
func GetDidMetadata(walletHandle int, Did string) (string, error) {
	upDid := unsafe.Pointer(C.CString(Did))
	defer C.free(upDid)
	channel := did.GetDidMetadata(walletHandle, upDid)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// GetDidWithMetadata retrieves DID, metadata and verkey stored in the wallet.
func GetDidWithMetadata(walletHandle int, Did string) (string, error) {
	upDid := unsafe.Pointer(C.CString(Did))
	defer C.free(upDid)
	channel := did.GetDidWithMetadata(walletHandle, upDid)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// ListDidsWithMeta lists DIDs and metadata stored in the wallet.
func ListDidsWithMeta(walletHandle int) (string, error) {
	channel := did.ListDidsWithMeta(walletHandle)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// AbbreviateVerKey retrieves abbreviated key if exists, otherwise returns full verkey.
func AbbreviateVerKey(Did string, verKey string) (string, error) {
	upDid := unsafe.Pointer(C.CString(Did))
	defer C.free(upDid)
	upVerKey := unsafe.Pointer(C.CString(verKey))
	defer C.free(upVerKey)

	channel := did.AbbreviateVerKey(upDid, upVerKey)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// QualifyDid updates DID related entities stored in the wallet.
func QualifyDid(walletHandle int, Did string, method string) (string, error) {

	upDid := unsafe.Pointer(C.CString(Did))
	defer C.free(upDid)
	upMethod := unsafe.Pointer(C.CString(method))
	defer C.free(upMethod)

	channel := did.QualifyDid(walletHandle, upDid, upMethod)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}