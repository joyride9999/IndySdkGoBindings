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

import (
	"indySDK/wallet"
)

// CreateWallet creates a new secure wallet with the given unique name
func CreateWallet(config wallet.Config, credential wallet.Credential) error {
	channel := wallet.CreateWallet(config, credential)
	result := <-channel
	return result.Error
}

// OpenWallet opens an existing wallet
func OpenWallet(config wallet.Config, credential wallet.Credential) (int, error) {
	channel := wallet.OpenWallet(config, credential)
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
	channel := wallet.DeleteWallet(config, credentials)
	result := <-channel
	return result.Error
}

// GenerateWalletKey generate wallet master key
func GenerateWalletKey(config wallet.Config) error {
	channel := wallet.GenerateWalletKey(config)
	result := <-channel
	return result.Error
}

// ExportWallet exports opened wallet
func ExportWallet(wh int, config wallet.ExportConfig) error {
	channel := wallet.ExportWallet(wh, config)
	result := <-channel
	return result.Error
}

// ImportWallet creates new secure wallet and imports its content
func ImportWallet(config wallet.Config, credentials wallet.Credential, import_config wallet.ImportConfig) error {
	channel := wallet.ImportWallet(config, credentials, import_config)
	result := <-channel
	return result.Error
}

// RegisterWalletStorage registers new wallet type
func RegisterWalletStorage(storageType string, storage wallet.IWalletStorage) error {
	channel := wallet.RegisterWalletStorage(storageType, storage)
	result := <-channel
	return result.Error
}
