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

import "indySDK/did"

// CreateAndStoreDID creates and DID with keys ... nothing is written to blockchain
// returns did, verkey, error
func CreateAndStoreDID(walletHandle int, seed string) (string, string, error) {
	channel := did.CreateAndStoreMyDid(walletHandle, seed)
	result := <-channel
	if result.Error != nil {
		return "", "", result.Error
	}
	return result.Results[0].(string), result.Results[1].(string), result.Error
}

// ReplaceKeyStart generates temporary key for an existing DID.
func ReplaceKeyStart(walletHandle int, Did string, identityJson string) (string, error) {
	channel := did.ReplaceKeyStart(walletHandle, Did, identityJson)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// ReplaceKeyApply applies temporary keys as main for existing DID
func ReplaceKeyApply(walletHandle int, Did string) error {
	channel := did.ReplaceKeyApply(walletHandle, Did)
	result := <-channel
	return result.Error
}

// StoreTheirDid saves DID for a pairwise connection in a secured wallet to verify transaction.
func StoreTheirDid(walletHandle int, identityJson string) error {
	channel := did.StoreTheirDid(walletHandle, identityJson)
	result := <-channel
	return result.Error
}

// KeyForDid returns ver key for DID.
func KeyForDid(poolHandle int, walletHandle int, Did string) (string, error) {
	channel := did.KeyForDid(poolHandle, walletHandle, Did)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// KeyForLocalDID gets the key for the local DID.
// returns key, error
func KeyForLocalDID(walletHandle int, Did string) (string, error) {
	channel := did.KeyForLocalDid(walletHandle, Did)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// SetEndPointForDid set/replaces endpoint information for the given DID
func SetEndPointForDid(walletHandle int, Did string, address string, transportKey string) error {
	channel := did.SetEndPointForDid(walletHandle, Did, address, transportKey)
	result := <-channel
	return result.Error
}

// GetEndPointForDid returns endpoint information for the given DID
func GetEndPointForDid(walletHandle int, poolHandle int, Did string) (string, string, error) {
	channel := did.GetEndPointForDid(walletHandle, poolHandle, Did)
	result := <-channel
	if result.Error != nil {
		return "", "", result.Error
	}
	return result.Results[0].(string), result.Results[1].(string), result.Error
}

// SetDidMetadata saves/replaces meta information for the given DID.
func SetDidMetadata(walletHandle int, Did string, metadata string) error {
	channel := did.SetDidMetadata(walletHandle, Did, metadata)
	result := <-channel
	return result.Error
}

// GetDidMetadata retrieves meta information for the given DID.
func GetDidMetadata(walletHandle int, Did string) (string, error) {
	channel := did.GetDidMetadata(walletHandle, Did)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// GetDidWithMetadata retrieves DID, metadata and verkey stored in the wallet.
func GetDidWithMetadata(walletHandle int, Did string) (string, error) {
	channel := did.GetDidWithMetadata(walletHandle, Did)
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
	channel := did.AbbreviateVerKey(Did, verKey)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// QualifyDid updates DID related entities stored in the wallet.
func QualifyDid(walletHandle int, Did string, method string) (string, error) {
	channel := did.QualifyDid(walletHandle, Did, method)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}