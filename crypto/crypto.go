/*
// ******************************************************************
// Purpose: Wrapper to call libindy, imports functions from indy_crypto.h
// Author:  adrian.toader@siemens.com
// Notes:
// Copyright (c): Siemens SRL
// This work is licensed under the terms of the Apache License Version 2.0.  See
// the LICENSE.txt file in the top-level directory.
// ******************************************************************
*/

package crypto

/*
#cgo CFLAGS: -I ../include
#cgo LDFLAGS: -L${SRCDIR}/../lib -lindy
#include <indy_core.h>
typedef void (*cb_createKey)(indy_handle_t, indy_error_t, char*);
typedef void (*cb_setKeyMetadata)(indy_handle_t, indy_error_t);
typedef void (*cb_getKeyMetadata)(indy_handle_t, indy_error_t, char*);
typedef void (*cb_sign)(indy_handle_t, indy_error_t, indy_u8_t*, indy_u32_t);
typedef void (*cb_verify)(indy_handle_t, indy_error_t, bool);
typedef void (*cb_anonCrypt)(indy_handle_t, indy_error_t, indy_u8_t*, indy_u32_t);
typedef void (*cb_anonDecrypt)(indy_handle_t, indy_error_t, indy_u8_t*, indy_u32_t);
typedef void (*cb_packMsg)(indy_handle_t, indy_error_t, indy_u8_t*, indy_u32_t);
typedef void (*cb_unpackMsg)(indy_handle_t, indy_error_t, indy_u8_t*, indy_u32_t);
//typedef void (*cb_authCrypt)(indy_handle_t, indy_error_t, indy_u8_t, indy_u32_t);
//typedef void (*cb_authDecrypt)(indy_handle_t, indy_error_t, char*, indy_u8_t, indy_u32_t);

extern void createKeyCB(indy_handle_t, indy_error_t, char*);
extern void setKeyMetadataCB(indy_handle_t, indy_error_t);
extern void getKeyMetadataCB(indy_handle_t, indy_error_t, char*);
extern void signCB(indy_handle_t, indy_error_t, indy_u8_t*, indy_u32_t);
extern void verifyCB(indy_handle_t, indy_error_t, bool);
extern void anonCryptCB(indy_handle_t, indy_error_t, indy_u8_t*, indy_u32_t);
extern void anonDecryptCB(indy_handle_t, indy_error_t, indy_u8_t*, indy_u32_t);
extern void packMsgCB(indy_handle_t, indy_error_t, indy_u8_t*, indy_u32_t);
extern void unpackMsgCB(indy_handle_t, indy_error_t, indy_u8_t*, indy_u32_t);
//extern void authCryptCB(indy_handle_t, indy_error_t, indy_u8_t, indy_u32_t);
//extern void authDecryptCB(indy_handle_t, indy_error_t, char*, indy_u8_t, indy_u32_t);
*/
import "C"

import (
	"github.com/joyride9999/IndySdkGoBindings/indyUtils"
	"errors"
	"unsafe"
)

//export createKeyCB
func createKeyCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, vkey *C.char) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil, Results: []interface{}{string(C.GoString(vkey))}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// CreateKey creates keys pair and stores in the wallet
func CreateKey(wh int, key unsafe.Pointer) chan indyUtils.IndyResult {
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
			Creates keys pair and stores in the wallet.

			:param wallet_handle: Wallet handle (created by open_wallet).
			:param key_json: Key information as json. Example:
			{
			    "seed": string, (optional) Seed that allows deterministic key creation (if not set random one will be created).
			                               Can be UTF-8, base64 or hex string.
			    "crypto_type": string, // Optional (if not set then ed25519 curve is used); Currently only 'ed25519' value is supported for this field.
			}

		    :return: Error code

	*/

	// Call indy_create_key
	res := C.indy_create_key(commandHandle,
		(C.indy_handle_t)(wh),
		(*C.char)(key),
		(C.cb_createKey)(unsafe.Pointer(C.createKeyCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export setKeyMetadataCB
func setKeyMetadataCB(commandHandle C.indy_handle_t, indyError C.indy_error_t) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// SetKeyMetadata saves/replaces the meta information for the giving key in the wallet
func SetKeyMetadata(wh int, verkey, metadata unsafe.Pointer) chan indyUtils.IndyResult {
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
		Saves/replaces the meta information for the giving key in the wallet.

		:param wallet_handle: Wallet handle (created by open_wallet).
		:param verkey: the key (verkey, key id) to store metadata.
		:param metadata: the meta information that will be store with the key.

		:return: Error code
	*/

	// Call indy_set_key_metadata
	res := C.indy_set_key_metadata(commandHandle,
		(C.indy_handle_t)(wh),
		(*C.char)(verkey),
		(*C.char)(metadata),
		(C.cb_setKeyMetadata)(unsafe.Pointer(C.setKeyMetadataCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export getKeyMetadataCB
func getKeyMetadataCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, meta *C.char) {
	// TODO is meta param ok here?
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil, Results: []interface{}{string(C.GoString(meta))}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// GetKeyMetadata retrieves the meta information for the giving key in the wallet
func GetKeyMetadata(wh int, verkey unsafe.Pointer) chan indyUtils.IndyResult {
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
		Retrieves the meta information for the giving key in the wallet.

		:param wallet_handle: Wallet handle (created by open_wallet).
		:param verkey: the key (verkey, key id) to store metadata.

		:return: Error code
	*/

	// Call indy_get_key_metadata
	res := C.indy_get_key_metadata(commandHandle,
		(C.indy_handle_t)(wh),
		(*C.char)(verkey),
		(C.cb_getKeyMetadata)(unsafe.Pointer(C.getKeyMetadataCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export signCB
func signCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, sigRaw *C.indy_u8_t, sigLen C.indy_u32_t) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil, Results: []interface{}{C.GoBytes(unsafe.Pointer(sigRaw), C.int(sigLen))}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// Sign signs a message with a key
func Sign(wh int, signerVK unsafe.Pointer, messageRaw unsafe.Pointer, messageLen uint32) chan indyUtils.IndyResult {
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
		Signs a message with a key.
		Note to use DID keys with this function you can call indy_key_for_did to get key id (verkey)
		for specific DID.

		:param wallet_handle: Wallet handle (created by open_wallet).
		:param signer_vk: id (verkey) of message signer. The key must be created by calling indy_create_key or indy_create_and_store_my_did
		:param message_raw: a pointer to first byte of message to be signed
		:param message_len: a message length

		:return: a signature string.
	*/

	// Call indy_crypto_sign
	res := C.indy_crypto_sign(commandHandle,
		(C.indy_handle_t)(wh),
		(*C.char)(signerVK),
		(*C.indy_u8_t)(messageRaw),
		C.indy_u32_t(messageLen),
		(C.cb_sign)(unsafe.Pointer(C.signCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export verifyCB
func verifyCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, valid C.indy_bool_t) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil, Results: []interface{}{bool(valid)}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// Verify verify a signature with a verkey.
func Verify(signerVK unsafe.Pointer, messageRaw unsafe.Pointer, messageLen uint32, signatureRaw unsafe.Pointer, signatureLen uint32) chan indyUtils.IndyResult {
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
		Verify a signature with a verkey.
		Note to use DID keys with this function you can call indy_key_for_did to get key id (verkey)
		for specific DID.

		:param signer_vk: id (verkey) of message signer. The key must be created by calling indy_create_key or indy_create_and_store_my_did
		:param message_raw: a pointer to first byte of message that has been signed
		:param message_len: a message length
		:param signature_raw: a pointer to first byte of signature to be verified
		:param signature_len: a signature length

		:return:
	*/

	// Call indy_crypto_verify
	res := C.indy_crypto_verify(commandHandle,
		(*C.char)(signerVK),
		(*C.indy_u8_t)(messageRaw),
		C.indy_u32_t(messageLen),
		(*C.indy_u8_t)(signatureRaw),
		C.indy_u32_t(signatureLen),
		(C.cb_verify)(unsafe.Pointer(C.verifyCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export anonCryptCB
func anonCryptCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, msg *C.indy_u8_t, msgLen C.indy_u32_t) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil,
			Results: []interface{}{C.GoBytes(unsafe.Pointer(msg), C.int(uint32(msgLen))), uint32(msgLen)}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// AnonCrypt encrypts a message by anonymous-encryption scheme
func AnonCrypt(recipientVK unsafe.Pointer, messageRaw unsafe.Pointer, messageLen uint32) chan indyUtils.IndyResult {
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
		Encrypts a message by anonymous-encryption scheme.
		Sealed boxes are designed to anonymously send messages to a Recipient given its public key.
		Only the Recipient can decrypt these messages, using its private key.
		While the Recipient can verify the integrity of the message, it cannot verify the identity of the Sender.

		Note to use DID keys with this function you can call indy_key_for_did to get key id (verkey)
		for specific DID.

		Note: use indy_pack_message() function for A2A goals.

		:param recipient_vk: verkey of message recipient
		:param message_raw: a pointer to first byte of message that to be encrypted
		:param message_len: a message length

		:return:
	*/

	// Call indy_crypto_anon_crypt
	res := C.indy_crypto_anon_crypt(commandHandle,
		(*C.char)(recipientVK),
		(*C.indy_u8_t)(messageRaw),
		C.indy_u32_t(messageLen),
		(C.cb_anonCrypt)(unsafe.Pointer(C.anonCryptCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export anonDecryptCB
func anonDecryptCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, msg *C.indy_u8_t, msgLen C.indy_u32_t) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil,
			Results: []interface{}{C.GoBytes(unsafe.Pointer(msg), C.int(msgLen))}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// AnonDecrypt decrypts a message by anonymous-encryption scheme
func AnonDecrypt(wh int, recipientVK, messageRaw unsafe.Pointer, messageLen uint32) chan indyUtils.IndyResult {
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)
	/*
		Dencrypts a message by anonymous-encryption scheme.
		Sealed boxes are designed to anonymously send messages to a Recipient given its public key.
		Only the Recipient can decrypt these messages, using its private key.
		While the Recipient can verify the integrity of the message, it cannot verify the identity of the Sender.

		Note to use DID keys with this function you can call indy_key_for_did to get key id (verkey)
		for specific DID.

		Note: use indy_unpack_message() function for A2A goals.

		:param wallet_handle: wallet handler (created by open_wallet).
		:param recipient_vk: verkey of message recipient
		:param message_raw: a pointer to first byte of message that to be encrypted
		:param message_len: a message length

		:return:
	*/

	// Call indy_crypto_anon_decrypt
	res := C.indy_crypto_anon_decrypt(commandHandle,
		(C.indy_handle_t)(wh),
		(*C.char)(recipientVK),
		(*C.indy_u8_t)(messageRaw),
		C.indy_u32_t(messageLen),
		(C.cb_anonDecrypt)(unsafe.Pointer(C.anonDecryptCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export packMsgCB
func packMsgCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, jwe *C.indy_u8_t, jweLen C.indy_u32_t) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil,
			Results: []interface{}{C.GoBytes(unsafe.Pointer(jwe), C.int(uint32(jweLen))), uint32(jweLen)}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// PackMsg packs a message by encrypting the message and serializes it in a JWE-like format
func PackMsg(wh int, messageRaw unsafe.Pointer, messageLen uint32, receiverKeys, sender unsafe.Pointer) chan indyUtils.IndyResult {
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
		Packs a message by encrypting the message and serializes it in a JWE-like format (Experimental)

		Note to use DID keys with this function you can call indy_key_for_did to get key id (verkey)
		for specific DID.

		:param wallet_handle: wallet handle (created by open_wallet).
		:param message: a pointer to the first byte of the message to be packed
		:param message_len: the length of the message
		:param receivers: a string in the format of a json list which will contain the list of receiver's keys
		               the message is being encrypted for.
		               Example:
		               "[<receiver edge_agent_1 verkey>, <receiver edge_agent_2 verkey>]"
		:param sender: the sender's verkey as a string When null pointer is used in this parameter, anoncrypt is used

		:return:
	*/

	// Call indy_pack_message
	res := C.indy_pack_message(commandHandle,
		(C.indy_handle_t)(wh),
		(*C.indy_u8_t)(messageRaw),
		C.indy_u32_t(messageLen),
		(*C.char)(receiverKeys),
		(*C.char)(sender),
		(C.cb_packMsg)(unsafe.Pointer(C.packMsgCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export unpackMsgCB
func unpackMsgCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, js *C.indy_u8_t, jsLen C.indy_u32_t) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil,
			Results: []interface{}{C.GoBytes(unsafe.Pointer(js), C.int(uint32(jsLen))), uint32(jsLen)}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// UnpackMsg unpacks a JWE-like formatted message outputted by indy_pack_message
func UnpackMsg(wh int, messageRaw unsafe.Pointer, messageLen uint32) chan indyUtils.IndyResult {
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)
	/*
		Unpacks a JWE-like formatted message outputted by indy_pack_message (Experimental)

		:param wallet_handle: wallet handle (created by open_wallet).
		:param jwe_data: a pointer to the first byte of the message to be packed
		:param jwe_len: the length of the message

		:return:
	*/

	// Call indy_unpack_message

	res := C.indy_unpack_message(commandHandle,
		(C.indy_handle_t)(wh),
		(*C.indy_u8_t)(messageRaw),
		C.indy_u32_t(messageLen),
		(C.cb_unpackMsg)(unsafe.Pointer(C.unpackMsgCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

////export authCryptCB
//func authCryptCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, msg C.indy_u8_t, msgLen C.indy_u32_t) {
//	if indyError == 0 {
//		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil, Results: []interface{}{uint8(msg), uint32(msgLen)}})
//	} else {
//		errMsg := indyUtils.GetIndyError(int(indyError))
//		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
//	}
//}
//
//// AuthCrypt encrypt a message by authenticated-encryption scheme
//func AuthCrypt(wh int, senderVK string, recipientVK string, messageRaw *uint8, messageLen uint32) chan indyUtils.IndyResult {
//
//	handle, future := indyUtils.NewFutureCommand()
//
//	commandHandle := (C.indy_handle_t)(handle)
//
//	/*
//		**** THIS FUNCTION WILL BE DEPRECATED USE indy_pack_message() INSTEAD ****
//		Encrypt a message by authenticated-encryption scheme.
//
//		Sender can encrypt a confidential message specifically for Recipient, using Sender's public key.
//		Using Recipient's public key, Sender can compute a shared secret key.
//		Using Sender's public key and his secret key, Recipient can compute the exact same shared secret key.
//		That shared secret key can be used to verify that the encrypted message was not tampered with,
//		before eventually decrypting it.
//
//		Note to use DID keys with this function you can call indy_key_for_did to get key id (verkey)
//		for specific DID.
//
//		:param wallet_handle: wallet handle (created by open_wallet).
//		:param sender_vk: id (verkey) of message sender. The key must be created by calling indy_create_key or indy_create_and_store_my_did
//		:param recipient_vk: id (verkey) of message recipient
//		:param message_raw: a pointer to first byte of message that to be encrypted
//		:param message_len: a message length
//
//		:return:
//	*/
//
//	// Call indy_crypto_auth_crypt
//	res := C.indy_crypto_auth_crypt(commandHandle,
//		(C.indy_handle_t)(wh),
//		C.CString(senderVK),
//		C.CString(recipientVK),
//		(*C.indy_u8_t)(unsafe.Pointer(messageRaw)),
//		C.indy_u32_t(messageLen),
//		(C.cb_authCrypt)(unsafe.Pointer(C.authCryptCB)))
//	if res != 0 {
//		errMsg := indyUtils.GetIndyError(int(res))
//		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
//		return future
//	}
//
//	return future
//}
//
////export authDecryptCB
//func authDecryptCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, senderVK *C.char, msg C.indy_u8_t, msgLen C.indy_u32_t) {
//	if indyError == 0 {
//		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil, Results: []interface{}{string(C.GoString(senderVK)), uint8(msg), uint32(msgLen)}})
//	} else {
//		errMsg := indyUtils.GetIndyError(int(indyError))
//		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
//	}
//}
//
//// AuthDecrypt dencrypt a message by authenticated-encryption scheme
//func AuthDecrypt(wh int, recipientVK string, messageRaw *uint8, messageLen uint32) chan indyUtils.IndyResult {
//
//	handle, future := indyUtils.NewFutureCommand()
//
//	commandHandle := (C.indy_handle_t)(handle)
//
//	/*
//		**** THIS FUNCTION WILL BE DEPRECATED USE indy_unpack_message() INSTEAD ****
//		Decrypt a message by authenticated-encryption scheme.
//
//		Sender can encrypt a confidential message specifically for Recipient, using Sender's public key.
//		Using Recipient's public key, Sender can compute a shared secret key.
//		Using Sender's public key and his secret key, Recipient can compute the exact same shared secret key.
//		That shared secret key can be used to verify that the encrypted message was not tampered with,
//		before eventually decrypting it.
//
//		Note to use DID keys with this function you can call indy_key_for_did to get key id (verkey)
//		for specific DID.
//
//		:param wallet_handle: wallet handle (created by open_wallet).
//		:param recipient_vk: id (verkey) of message recipient
//		:param message_raw: a pointer to first byte of message that to be encrypted
//		:param message_len: a message length
//
//		:return:
//	*/
//
//	// Call indy_crypto_auth_decrypt
//	res := C.indy_crypto_auth_decrypt(commandHandle,
//		(C.indy_handle_t)(wh),
//		C.CString(recipientVK),
//		(*C.indy_u8_t)(unsafe.Pointer(messageRaw)),
//		C.indy_u32_t(messageLen),
//		(C.cb_authDecrypt)(unsafe.Pointer(C.authDecryptCB)))
//	if res != 0 {
//		errMsg := indyUtils.GetIndyError(int(res))
//		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
//		return future
//	}
//
//	return future
//}
