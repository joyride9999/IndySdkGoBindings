/*
// ******************************************************************
// Purpose: wallet unit testing
// Author: adrian.toader@siemens.com
// Notes:
// Copyright (c): Siemens SRL
// This work is licensed under the terms of the Apache License Version 2.0.  See
// the LICENSE.txt file in the top-level directory.
// ******************************************************************
*/


package indySDK

import (
	"github.com/joyride9999/IndySdkGoBindings/wallet"
	"fmt"
	"testing"
)

func TestCreateWallet(t *testing.T) {
	// Prepare and run test cases
	type args struct {
		walletCfg  wallet.Config
		walletCred wallet.Credential
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"create-nonexisting-wallet",
			args{walletCfg: wallet.Config{
				ID:            "a1235",
				StorageType:   "default",
				StorageConfig: wallet.StorageConfig{Path: ".\\out\\wallet"},
			},
				walletCred: wallet.Credential{
					Key: "123",
				}},
			false},
		{"create-existing-wallet",
			args{walletCfg: wallet.Config{
				ID:            "a1235",
				StorageType:   "default",
				StorageConfig: wallet.StorageConfig{Path: ".\\out\\wallet"},
			},
				walletCred: wallet.Credential{
					Key: "123",
				}},
			true},
	}
	fmt.Println("Test")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CreateWallet(tt.args.walletCfg, tt.args.walletCred)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateWallet() error = '%v', wantErr = '%v'", err, tt.wantErr)
				return
			}
		})
	}
}

func TestDeleteWallet(t *testing.T) {
	// Prepare and run test cases
	type args struct {
		walletCfg  wallet.Config
		walletCred wallet.Credential
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"delete-existing-wallet",
			args{walletCfg: wallet.Config{
				ID:            "a1235",
				StorageType:   "default",
				StorageConfig: wallet.StorageConfig{Path: ".\\out\\wallet"},
			},
				walletCred: wallet.Credential{
					Key: "123",
				}},
			false},
		{"delete-nonexisting-wallet",
			args{walletCfg: wallet.Config{
				ID:            "a1235",
				StorageType:   "default",
				StorageConfig: wallet.StorageConfig{Path: ".\\out\\wallet"},
			},
				walletCred: wallet.Credential{
					Key: "123",
				}},
			true},
	}
	fmt.Println("Test")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//err := DeleteWallet(tt.args.walletCfg.Id)
			err := DeleteWallet(tt.args.walletCfg, tt.args.walletCred)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteWallet() error = '%v', wantErr = '%v'", err, tt.wantErr)
				return
			}
		})
	}
}

func TestOpenWallet(t *testing.T) {
	// Prepare and run test cases
	type args struct {
		walletCfg  wallet.Config
		walletCred wallet.Credential
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"open-existing-wallet",
			args{walletCfg: wallet.Config{
				ID:            "a1235",
				StorageType:   "default",
				StorageConfig: wallet.StorageConfig{Path: ".\\out\\wallet"},
			},
				walletCred: wallet.Credential{
					Key: "123",
				}},
			false},
	}
	fmt.Println("Test OpenWallet")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println("Opening wallet...")
			handler, errOpen := OpenWallet(tt.args.walletCfg, tt.args.walletCred)
			if (errOpen != nil) != tt.wantErr {
				t.Errorf("OpenWallet() error = '%v', wantErr = '%v'", errOpen, tt.wantErr)
				return
			}
			handler = handler
		})
	}
}

func TestExportWallet(t *testing.T) {
	// Prepare and run test cases
	type args struct {
		walletCfg    wallet.Config
		walletCred   wallet.Credential
		walletExport wallet.ExportConfig
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"export-existing-wallet",
			args{walletCfg: wallet.Config{
				ID:            "a1235",
				StorageType:   "default",
				StorageConfig: wallet.StorageConfig{Path: ".\\out\\wallet"},
			},
				walletCred: wallet.Credential{
					Key: "123",
				},
				walletExport: wallet.ExportConfig{
					Path: ".\\out\\wallet\\export",
					Key:  "123",
				}},
			false},
	}
	fmt.Println("Test ExportWallet")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println("Opening wallet...")
			handler, errOpen := OpenWallet(tt.args.walletCfg, tt.args.walletCred)
			if (errOpen != nil) != tt.wantErr {
				t.Errorf("OpenWallet() error = '%v', wantErr = '%v'", errOpen, tt.wantErr)
				return
			}
			fmt.Println("Exporting wallet...")
			errExport := ExportWallet(handler, tt.args.walletExport)
			if (errExport != nil) != tt.wantErr {
				t.Errorf("ExportWallet() error = '%v', wantErr = '%v'", errExport, tt.wantErr)
				return
			}
		})
	}
}

func TestImportWallet(t *testing.T) {
	// Prepare and run test cases
	type args struct {
		walletCfg    wallet.Config
		walletCred   wallet.Credential
		walletImport wallet.ImportConfig
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"import-wallet",
			args{walletCfg: wallet.Config{
				ID:            "a1239",
				StorageType:   "default",
				StorageConfig: wallet.StorageConfig{Path: ".\\out\\wallet"},
			},
				walletCred: wallet.Credential{
					Key: "123",
				},
				walletImport: wallet.ImportConfig{
					Path: ".\\out\\wallet\\export",
					Key:  "123",
				}},
			false},
	}
	fmt.Println("Test ImportWallet")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println("Importing wallet...")
			errImport := ImportWallet(tt.args.walletCfg, tt.args.walletCred, tt.args.walletImport)
			if (errImport != nil) != tt.wantErr {
				t.Errorf("ImportWallet() error = '%v', wantErr = '%v'", errImport, tt.wantErr)
			}
		})
	}
}

func TestGenerateWalletKey(t *testing.T) {
	// Prepare and run test cases
	type args struct {
		walletCfg  wallet.Config
		walletCred wallet.Credential
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"generate-wallet-key",
			args{walletCfg: wallet.Config{
				ID:            "a1237",
				StorageType:   "default",
				StorageConfig: wallet.StorageConfig{Path: ".\\out\\wallet"},
			}},
			false},
	}
	fmt.Println("Test GenerateWalletKey")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println("Generating wallet key...")
			errGenerate := GenerateWalletKey(tt.args.walletCfg)
			if (errGenerate != nil) != tt.wantErr {
				t.Errorf("GenerateWalletKey() error = '%v', wantErr = '%v'", errGenerate, tt.wantErr)
			}
		})
	}
}
