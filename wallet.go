/*
// ******************************************************************
// Purpose: exported public functions that handles wallet functions
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
	"github.com/joyride9999/IndySdkGoBindings/wallet"
	"encoding/json"
	"errors"
	"unsafe"
)

// CreateWallet creates a new secure wallet with the given unique name
func CreateWallet(config wallet.Config, credential wallet.Credential) error {

	jsonConfig, err := json.Marshal(config)
	if err != nil {
		return errors.New("cant read json")
	}
	jsonCredential, err := json.Marshal(credential)
	if err != nil {
		return errors.New("cant read json")
	}

	upWalletCfg := unsafe.Pointer(C.CString(string(jsonConfig)))
	defer C.free(upWalletCfg)
	upWalletCredential := unsafe.Pointer(C.CString(string(jsonCredential)))
	defer C.free(upWalletCredential)

	channel := wallet.CreateWallet(upWalletCfg, upWalletCredential)
	result := <-channel
	return result.Error
}

// OpenWallet opens an existing wallet
func OpenWallet(config wallet.Config, credential wallet.Credential) (int, error) {
	jsonConfig, err := json.Marshal(config)
	if err != nil {
		return 0, errors.New("cant read json")
	}
	jsonCredential, err := json.Marshal(credential)
	if err != nil {
		return 0, errors.New("cant read json")
	}

	upWalletCfg := unsafe.Pointer(C.CString(string(jsonConfig)))
	defer C.free(upWalletCfg)
	upWalletCredential := unsafe.Pointer(C.CString(string(jsonCredential)))
	defer C.free(upWalletCredential)

	channel := wallet.OpenWallet(upWalletCfg, upWalletCredential)
	result := <-channel
	if result.Error != nil {
		return 0, result.Error
	}
	return result.Results[0].(int), result.Error
}

// CloseWallet creates a new secure wallet with the given unique name
func CloseWallet(wh int) error {
	channel := wallet.CloseWallet(wh)
	result := <-channel
	return result.Error
}

// DeleteWallet deletes a secure wallet
func DeleteWallet(config wallet.Config, credentials wallet.Credential) error {

	jsonConfig, err := json.Marshal(config)
	if err != nil {
		return errors.New("cant read json")
	}
	jsonCredential, err := json.Marshal(credentials)
	if err != nil {
		return errors.New("cant read json")
	}

	upWalletCfg := unsafe.Pointer(C.CString(string(jsonConfig)))
	defer C.free(upWalletCfg)
	upWalletCredential := unsafe.Pointer(C.CString(string(jsonCredential)))
	defer C.free(upWalletCredential)

	channel := wallet.DeleteWallet(upWalletCfg, upWalletCredential)
	result := <-channel
	return result.Error
}

// GenerateWalletKey generate wallet master key
func GenerateWalletKey(config wallet.Config) error {
	jsonConfig, err := json.Marshal(config)
	if err != nil {
		return errors.New("cant read json")
	}
	upWalletCfg := unsafe.Pointer(C.CString(string(jsonConfig)))
	defer C.free(upWalletCfg)

	channel := wallet.GenerateWalletKey(upWalletCfg)
	result := <-channel
	return result.Error
}

// ExportWallet exports opened wallet
func ExportWallet(wh int, config wallet.ExportConfig) error {
	jsonConfig, err := json.Marshal(config)
	if err != nil {
		return errors.New("cant read json")
	}
	upWalletCfg := unsafe.Pointer(C.CString(string(jsonConfig)))
	defer C.free(upWalletCfg)

	channel := wallet.ExportWallet(wh, upWalletCfg)
	result := <-channel
	return result.Error
}

// ImportWallet creates new secure wallet and imports its content
func ImportWallet(config wallet.Config, credentials wallet.Credential, importConfig wallet.ImportConfig) error {

	jsonConfig, err := json.Marshal(config)
	if err != nil {
		return errors.New("cant read json")
	}
	jsonCredential, err := json.Marshal(credentials)
	if err != nil {
		return errors.New("cant read json")
	}
	jsonImportCfg, err := json.Marshal(importConfig)
	if err != nil {
		return errors.New("cant read json")
	}

	upWalletCfg := unsafe.Pointer(C.CString(string(jsonConfig)))
	defer C.free(upWalletCfg)
	upWalletCredential := unsafe.Pointer(C.CString(string(jsonCredential)))
	defer C.free(upWalletCredential)
	upWalletImportCfg := unsafe.Pointer(C.CString(string(jsonImportCfg)))
	defer C.free(upWalletImportCfg)

	channel := wallet.ImportWallet(upWalletCfg, upWalletCredential, upWalletImportCfg)
	result := <-channel
	return result.Error
}

// RegisterWalletStorage registers new wallet type
func RegisterWalletStorage(storageType string, storage wallet.IWalletStorage) error {

	upStorageType := unsafe.Pointer(C.CString(storageType))
	defer C.free(upStorageType)

	channel := wallet.RegisterWalletStorage(upStorageType, storage)
	result := <-channel
	return result.Error
}
