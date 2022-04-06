/*
// ******************************************************************
// Purpose: exported public functions that handles crypto functions
// from libindy
// Author:  angel.draghici@siemens.com, adrian.toader@siemens.com
// Notes:
// Copyright (c): Siemens SRL
// This work is licensed under the terms of the Apache License Version 2.0.  See
// the LICENSE.txt file in the top-level directory.
// ******************************************************************
*/

package indySDK

import "C"
import "indySDK/crypto"

// CreateKey creates keys pair and stores in the wallet
func CreateKey(wh int, key crypto.Key) (string, error) {
	channel := crypto.CreateKey(wh, key)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// SetKeyMetadata saves/replaces the meta information for the giving key in the wallet
func SetKeyMetadata(wh int, verkey string, metadata string) error {
	channel := crypto.SetKeyMetadata(wh, verkey, metadata)
	result := <-channel
	return result.Error
}

// GetKeyMetadata retrieves the meta information for the giving key in the wallet
func GetKeyMetadata(wh int, verkey string) (string, error) {
	channel := crypto.GetKeyMetadata(wh, verkey)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// Sign signs a message with a key
func Sign(wh int, signerVK string, messageRaw []uint8, messageLen uint32) ([]uint8, error) {
	channel := crypto.Sign(wh, signerVK, messageRaw, messageLen)
	result := <-channel
	if result.Error != nil {
		return []uint8(""), result.Error
	}
	return result.Results[0].([]uint8), result.Error
}

// Verify signs a message with a key
func Verify(signerVK string, messageRaw []uint8, messageLen uint32, signatureRaw []uint8, signatureLen uint32) (bool, error) {
	channel := crypto.Verify(signerVK, messageRaw, messageLen, signatureRaw, signatureLen)
	result := <-channel
	if result.Error != nil {
		return false, result.Error
	}
	return result.Results[0].(bool), result.Error
}

// AnonCrypt encrypts a message by anonymous-encryption scheme
func AnonCrypt(recipientVK string, messageRaw []uint8, messageLen uint32) ([]uint8, error) {
	channel := crypto.AnonCrypt(recipientVK, messageRaw, messageLen)
	result := <-channel
	if result.Error != nil {
		return []uint8(""), result.Error
	}
	return result.Results[0].([]uint8), result.Error
}

// AnonDecrypt decrypts a message by anonymous-encryption scheme
func AnonDecrypt(wh int, recipientVK string, messageRaw []uint8, messageLen uint32) ([]uint8, error) {
	channel := crypto.AnonDecrypt(wh, recipientVK, messageRaw, messageLen)
	result := <-channel
	if result.Error != nil {
		return []uint8(""), result.Error
	}
	return result.Results[0].([]uint8), result.Error
}

// PackMsg packs a message by encrypting the message and serializes it in a JWE-like format
func PackMsg(wh int, messageRaw []uint8, messageLen uint32, receiverKeys string, sender string) ([]uint8, error) {
	channel := crypto.PackMsg(wh, messageRaw, messageLen, receiverKeys, sender)
	result := <-channel
	if result.Error != nil {
		return []uint8(""), result.Error
	}
	return result.Results[0].([]uint8), result.Error
}

// UnpackMsg packs a message by encrypting the message and serializes it in a JWE-like format
func UnpackMsg(wh int, messageRaw []uint8, messageLen uint32) ([]uint8, error) {
	channel := crypto.UnpackMsg(wh, messageRaw, messageLen)
	result := <-channel
	if result.Error != nil {
		return []uint8(""), result.Error
	}
	return result.Results[0].([]uint8), result.Error
}

////indy_crypto_auth_crypt	/// **** THIS FUNCTION WILL BE DEPRECATED USE indy_pack_message() INSTEAD ****
//// AuthCrypt dncrypt a message by authenticated-encryption scheme
//func AuthCrypt(wh int, senderVK string, recipientVK string, messageRaw *uint8, messageLen uint32) (string, error) {
//	channel := crypto.AuthCrypt(wh, senderVK, recipientVK, messageRaw, messageLen)
//	result := <-channel
//	if result.Error != nil {
//		return "", result.Error
//	}
//	return result.Results[0].(string), result.Error
//}
//
////indy_crypto_auth_decrypt	/// **** THIS FUNCTION WILL BE DEPRECATED USE indy_unpack_message() INSTEAD ****
//// AuthDecrypt dncrypt a message by authenticated-encryption scheme
//func AuthDecrypt(wh int, recipientVK string, messageRaw *uint8, messageLen uint32) (string, string, error) {
//	channel := crypto.AuthDecrypt(wh, recipientVK, messageRaw, messageLen)
//	result := <-channel
//	if result.Error != nil {
//		return "", "", result.Error
//	}
//	return result.Results[0].(string), result.Results[1].(string), result.Error
//}
