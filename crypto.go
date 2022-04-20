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

/*
#include <stdlib.h>
*/
import "C"
import (
	"github.com/joyride9999/IndySdkGoBindings/crypto"
	"encoding/json"
	"errors"
	"unsafe"
)

// CreateKey creates keys pair and stores in the wallet
func CreateKey(wh int, key crypto.Key) (string, error) {

	jsonKey, err := json.Marshal(key)
	if err != nil {
		return "", errors.New("cant read json")
	}
	keyString := string(jsonKey)
	upKey := unsafe.Pointer(C.CString(keyString))
	defer C.free(upKey)

	channel := crypto.CreateKey(wh, upKey)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// SetKeyMetadata saves/replaces the meta information for the giving key in the wallet
func SetKeyMetadata(wh int, verkey string, metadata string) error {

	upVerKey := unsafe.Pointer(C.CString(verkey))
	defer C.free(upVerKey)
	upMetadata := unsafe.Pointer(C.CString(metadata))
	defer C.free(upMetadata)

	channel := crypto.SetKeyMetadata(wh, upVerKey, upMetadata)
	result := <-channel
	return result.Error
}

// GetKeyMetadata retrieves the meta information for the giving key in the wallet
func GetKeyMetadata(wh int, verkey string) (string, error) {

	upVerKey := unsafe.Pointer(C.CString(verkey))
	defer C.free(upVerKey)
	channel := crypto.GetKeyMetadata(wh, upVerKey)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// Sign signs a message with a key
func Sign(wh int, signerVK string, messageRaw []uint8, messageLen uint32) ([]uint8, error) {

	upSignerVK := unsafe.Pointer(C.CString(signerVK))
	defer C.free(upSignerVK)
	upMessageRaw := unsafe.Pointer(C.CBytes(messageRaw))
	defer C.free(upMessageRaw)

	channel := crypto.Sign(wh, upSignerVK, upMessageRaw, messageLen)
	result := <-channel
	if result.Error != nil {
		return []uint8(""), result.Error
	}
	return result.Results[0].([]uint8), result.Error
}

// Verify signs a message with a key
func Verify(signerVK string, messageRaw []uint8, messageLen uint32, signatureRaw []uint8, signatureLen uint32) (bool, error) {

	upSignerVK := unsafe.Pointer(C.CString(signerVK))
	defer C.free(upSignerVK)
	upMessageRaw := unsafe.Pointer(C.CBytes(messageRaw))
	defer C.free(upMessageRaw)
	upSignatureRaw := unsafe.Pointer(C.CBytes(signatureRaw))
	defer C.free(upSignatureRaw)

	channel := crypto.Verify(upSignerVK, upMessageRaw, messageLen, upSignatureRaw, signatureLen)
	result := <-channel
	if result.Error != nil {
		return false, result.Error
	}
	return result.Results[0].(bool), result.Error
}

// AnonCrypt encrypts a message by anonymous-encryption scheme
func AnonCrypt(recipientVK string, messageRaw []uint8, messageLen uint32) ([]uint8, error) {

	upRecipientVK := unsafe.Pointer(C.CString(recipientVK))
	defer C.free(upRecipientVK)
	upMessageRaw := unsafe.Pointer(C.CBytes(messageRaw))
	defer C.free(upMessageRaw)

	channel := crypto.AnonCrypt(upRecipientVK, upMessageRaw, messageLen)
	result := <-channel
	if result.Error != nil {
		return []uint8(""), result.Error
	}
	return result.Results[0].([]uint8), result.Error
}

// AnonDecrypt decrypts a message by anonymous-encryption scheme
func AnonDecrypt(wh int, recipientVK string, messageRaw []uint8, messageLen uint32) ([]uint8, error) {

	upRecipientVK := unsafe.Pointer(C.CString(recipientVK))
	defer C.free(upRecipientVK)
	upMessageRaw := unsafe.Pointer(C.CBytes(messageRaw))
	defer C.free(upMessageRaw)

	channel := crypto.AnonDecrypt(wh, upRecipientVK, upMessageRaw, messageLen)
	result := <-channel
	if result.Error != nil {
		return []uint8(""), result.Error
	}
	return result.Results[0].([]uint8), result.Error
}

// PackMsg packs a message by encrypting the message and serializes it in a JWE-like format
func PackMsg(wh int, messageRaw []uint8, messageLen uint32, receiverKeys string, sender string) ([]uint8, error) {

	upMessageRaw := unsafe.Pointer(C.CBytes(messageRaw))
	defer C.free(upMessageRaw)
	upReceiverKeys := unsafe.Pointer(C.CString(receiverKeys))
	defer C.free(upReceiverKeys)
	upSender := unsafe.Pointer(C.CString(sender))
	defer C.free(upSender)

	channel := crypto.PackMsg(wh, upMessageRaw, messageLen, upReceiverKeys, upSender)
	result := <-channel
	if result.Error != nil {
		return []uint8(""), result.Error
	}
	return result.Results[0].([]uint8), result.Error
}

// UnpackMsg packs a message by encrypting the message and serializes it in a JWE-like format
func UnpackMsg(wh int, messageRaw []uint8, messageLen uint32) ([]uint8, error) {

	upMessageRaw := unsafe.Pointer(C.CBytes(messageRaw))
	defer C.free(upMessageRaw)

	channel := crypto.UnpackMsg(wh, upMessageRaw, messageLen)
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
