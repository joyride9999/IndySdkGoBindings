/*
// ******************************************************************
// Purpose: crypto unit testing
// Author: angel.draghici@siemens.com, adrian.toader@siemens.com
// Notes:
// Copyright (c): Siemens SRL
// This work is licensed under the terms of the Apache License Version 2.0.  See
// the LICENSE.txt file in the top-level directory.
// ******************************************************************
*/

package indySDK

import (
	"encoding/json"
	"fmt"
	"indySDK/crypto"
	"indySDK/indyUtils"
	"testing"
)

func TestAnonCrypt(t *testing.T) {
	// Prepare wallet for tests
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	// Get did for wallet
	_, verKey, errDid := CreateAndStoreDID(walletHandle, "")
	if errDid != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDid)
		return
	}

	// Message to be encrypted
	message := []uint8("{\"reqId\":1496822211362017764}")

	type args struct {
		Message []byte
		VerKey  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"anon-crypt-works", args{Message: message, VerKey: verKey}, false},
		{"anon-crpy-empty-message", args{Message: []byte(""), VerKey: verKey}, true},
		{"anon-crypt-invalid-ver-key", args{Message: message, VerKey: "invalid-ver-key"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypted, errAnon := AnonCrypt(tt.args.VerKey, tt.args.Message, uint32(len(tt.args.Message)))
			hasError := errAnon != nil
			if hasError != tt.wantErr {
				t.Errorf("AnonCrypt() error = '%v'", errAnon)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errAnon)
				return
			} else {
				fmt.Println(encrypted)
			}
		})
	}
	return
}

func TestAnonDecrypt(t *testing.T) {
	// Prepare wallet for tests
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	// Get did for wallet
	_, verKey, errDid := CreateAndStoreDID(walletHandle, "")
	if errDid != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDid)
		return
	}

	// Message to be encrypted
	message := []uint8("{\"reqId\":1496822211362017764}")

	// Message encryption
	encrypted, errCrypt := AnonCrypt(verKey, message, uint32(len(message)))
	if errCrypt != nil {
		t.Errorf("AnonCrypt() error = '%v'", errCrypt)
		return
	}

	type args struct {
		Message []byte
		VerKey  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"anon-decrypt-works", args{Message: encrypted, VerKey: verKey}, false},
		{"anon-decrpy-invalid-message", args{Message: []byte("invalid-message"), VerKey: verKey}, true},
		{"anon-decrypt-invalid-ver-key", args{Message: encrypted, VerKey: "invalid-ver-key"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decrypted, errDecrypt := AnonDecrypt(walletHandle, tt.args.VerKey, tt.args.Message, uint32(len(tt.args.Message)))
			hasError := errDecrypt != nil
			if hasError != tt.wantErr {
				t.Errorf("AnonDecrypt() error = '%v'", errDecrypt)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errDecrypt)
				return
			} else {
				fmt.Println(string(decrypted))
			}
		})
	}
	return
}

func TestCreateKey(t *testing.T) {
	// Prepare wallet for tests
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	// Create key information
	seed := crypto.Key{Seed: "00000000000000000000000011111111", CryptoType: "ed25519"}

	type args struct {
		WalletHandle int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"create-key-works", args{WalletHandle: walletHandle}, false},
		{"create-key-invalid-wallet-handle", args{WalletHandle: walletHandle + 100}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verKey, errCreateKey := CreateKey(tt.args.WalletHandle, seed)
			hasError := errCreateKey != nil
			if hasError != tt.wantErr {
				t.Errorf("CreateKey() error = '%v'", errCreateKey)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errCreateKey)
				return
			} else {
				fmt.Println(verKey)
			}
		})
	}
	return
}

func TestSign(t *testing.T) {
	// Prepare wallet for tests
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	// Create key information
	seed := crypto.Key{Seed: "00000000000000000000000011111111", CryptoType: "ed25519"}
	verKey, errCreateKey := CreateKey(walletHandle, seed)
	if errCreateKey != nil {
		t.Errorf("CreateKey() error = '%v'", errCreateKey)
		return
	}

	// Message to be signed
	message := []uint8("{\"reqId\":1496822211362017764}")

	type args struct {
		VerKey  string
		Message []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"sign-works", args{VerKey: verKey, Message: message}, false},
		{"sign-with-invalid-ver-key", args{VerKey: "invalid-ver-key", Message: message}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signature, errSign := Sign(walletHandle, tt.args.VerKey, tt.args.Message, uint32(len(tt.args.Message)))
			hasError := errSign != nil
			if hasError != tt.wantErr {
				t.Errorf("Sign() error = '%v'", errSign)
				return
			}
			if tt.wantErr {
				t.Log("Error expected: ", errSign)
				return
			} else {
				fmt.Println(signature)
			}
		})
	}

	return
}

func TestVerify(t *testing.T) {
	// Prepare wallet for tests
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	// Create key information
	seed := crypto.Key{Seed: "00000000000000000000000011111111", CryptoType: "ed25519"}
	verKey, errCreateKey := CreateKey(walletHandle, seed)
	if errCreateKey != nil {
		t.Errorf("CreateKey() error = '%v'", errCreateKey)
		return
	}

	// Sign the message
	message := []uint8("{\"reqId\":1496822211362017764}")
	signature, errSign := Sign(walletHandle, verKey, message, uint32(len(message)))
	if errSign != nil {
		t.Errorf("Sign() error = '%v'", errSign)
		return
	}

	type args struct {
		VerKey    string
		Signature []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"verify-works", args{VerKey: verKey, Signature: signature}, false},
		{"verify-invalid-signature", args{VerKey: verKey, Signature: []byte("invalid-signature")}, false},
		{"verify-invalid-ver-key", args{VerKey: "invalid-ver-key", Signature: signature}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verified, errVerify := Verify(tt.args.VerKey, message, uint32(len(message)), tt.args.Signature, uint32(len(signature)))
			hasError := errVerify != nil
			if hasError != tt.wantErr {
				t.Errorf("Verify() error = '%v'", errVerify)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errVerify)
				return
			} else {
				fmt.Println(verified)
			}
		})
	}

	return
}

func TestSetKeyMetadata(t *testing.T) {
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	// Create key information
	seed := crypto.Key{Seed: "00000000000000000000000011111111", CryptoType: "ed25519"}
	verKey, errCreateKey := CreateKey(walletHandle, seed)
	if errCreateKey != nil {
		t.Errorf("CreateKey() error = '%v'", errCreateKey)
		return
	}

	type args struct {
		VerKey string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"set-key-metadata-works", args{VerKey: verKey}, false},
		{"set-key-metadata-invalid-ver-key", args{VerKey: "invalid-ver-key"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errSet := SetKeyMetadata(walletHandle, tt.args.VerKey, metadata)
			hasError := errSet != nil
			if hasError != tt.wantErr {
				t.Errorf("SetKeyMetadata() error = '%v'", errSet)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errSet)
				return
			}
		})
	}

	return
}

func TestGetKeyMetadata(t *testing.T) {
	// Prepare wallet for tests
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	// Create key information
	seed := crypto.Key{Seed: "00000000000000000000000011111111", CryptoType: "ed25519"}
	verKey, errCreateKey := CreateKey(walletHandle, seed)
	if errCreateKey != nil {
		t.Errorf("CreateKey() error = '%v'", errCreateKey)
		return
	}

	// Set key metadata
	errSet := SetKeyMetadata(walletHandle, verKey, metadata)
	if errSet != nil {
		t.Errorf("SetKeyMetadata() error = '%v'", errSet)
		return
	}

	type args struct {
		VerKey string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"get-key-metadata-works", args{VerKey: verKey}, false},
		{"get-key-metadata-invalid-ver-key", args{VerKey: "invalid-ver-key"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keyMetadata, errGet := GetKeyMetadata(walletHandle, tt.args.VerKey)
			hasError := errGet != nil
			if hasError != tt.wantErr {
				t.Errorf("GetKeyMetadata() error = '%v'", errGet)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errGet)
				return
			} else {
				fmt.Println(keyMetadata)
			}
		})
	}

	return
}

func TestPackMsg(t *testing.T) {
	// Prepare wallet for tests
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	// Get verkey
	_, verKey, errCreateDid := CreateAndStoreDID(walletHandle, seedMy1)
	if errCreateDid != nil {
		t.Errorf("CreateAndStoreDID() error = '%v'", errCreateDid)
		return
	}

	// Create second wallet
	walletHandle2, errCreate2 := createWallet(holderConfig(), holderCredentials())
	if errCreate2 != nil && errCreate2.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle2, holderConfig(), holderCredentials())

	// Get verkey
	_, verKey2, errCreateDid2 := CreateAndStoreDID(walletHandle2, seedSteward1)
	if errCreateDid2 != nil {
		t.Errorf("CreateAndStoreDID() error = '%v'", errCreateDid2)
		return
	}

	message := []uint8("{\"reqId\":1496822211362017764}")
	receivedKeysJSON := []string{
		verKey2,
	}
	receivedKeys, _ := json.Marshal(receivedKeysJSON)

	type args struct {
		WalletHandle int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"pack-message-works", args{WalletHandle: walletHandle}, false},
		{"pack-message-invalid-wallet-handle", args{WalletHandle: walletHandle + 100}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packed, errPack := PackMsg(tt.args.WalletHandle, message, uint32(len(message)), string(receivedKeys), verKey)
			hasError := errPack != nil
			if hasError != tt.wantErr {
				t.Errorf("PackMsg() error = '%v'", errPack)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errPack)
				return
			} else {
				fmt.Println(packed)
			}
		})
	}

	return
}

func TestUnpackMsg(t *testing.T) {
	// Prepare wallet for tests
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	// Get verkey
	_, verKey, errCreateDid := CreateAndStoreDID(walletHandle, seedMy1)
	if errCreateDid != nil {
		t.Errorf("CreateAndStoreDID() error = '%v'", errCreateDid)
		return
	}

	// Create second wallet
	walletHandle2, errCreate2 := createWallet(holderConfig(), holderCredentials())
	if errCreate2 != nil && errCreate2.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle2, holderConfig(), holderCredentials())

	// Get verkey
	_, verKey2, errCreateDid2 := CreateAndStoreDID(walletHandle2, seedSteward1)
	if errCreateDid2 != nil {
		t.Errorf("CreateAndStoreDID() error = '%v'", errCreateDid2)
		return
	}

	message := []uint8("{\"reqId\":1496822211362017764}")
	receivedKeysJSON := []string{
		verKey2,
	}

	receivedKeys, _ := json.Marshal(receivedKeysJSON)
	keys := string(receivedKeys)
	packedMsg, errPack := PackMsg(walletHandle, message, uint32(len(message)), keys, verKey)
	if errPack != nil {
		t.Errorf("PackMsg() error = '%v'", errPack)
		return
	}

	type args struct {
		WalletHandle int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"unpack-message-works", args{WalletHandle: walletHandle2}, false},
		{"unpack-message-invalid-packed-message", args{WalletHandle: walletHandle2 + 100}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			unpackedMsg, errUnpack := UnpackMsg(tt.args.WalletHandle, packedMsg, uint32(len(packedMsg)))
			hasError := errUnpack != nil
			if hasError != tt.wantErr {
				t.Errorf("UnpackMsg() error = '%v'", errUnpack)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errUnpack)
				return
			} else {
				fmt.Println(string(unpackedMsg))
			}
		})
	}

	return
}
