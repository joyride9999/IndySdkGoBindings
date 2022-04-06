/*
// ******************************************************************
// Purpose: did unit testing
// Author: angel.draghici@siemens.com, adrian.toader@siemens.com
// Notes:
// Copyright (c): Siemens SRL
// This work is licensed under the terms of the Apache License Version 2.0.  See
// the LICENSE.txt file in the top-level directory.
// ******************************************************************
*/

package indySDK

import (
	"fmt"
	"indySDK/indyUtils"
	"testing"
)

func TestCreateAndStoreDid(t *testing.T) {
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	type args struct {
		WalletHandle int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"create-and-store-did-works", args{WalletHandle: walletHandle}, false},
		{"create-and-store-did-invalid-wallet-handle", args{WalletHandle: walletHandle + 100}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			did, verKey, errCreateDid := CreateAndStoreDID(tt.args.WalletHandle, "")
			hasError := errCreateDid != nil
			if hasError != tt.wantErr {
				t.Errorf("CreateAndStoreDid() error = '%v'", errCreateDid)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errCreateDid)
				return
			} else {
				fmt.Println(did + "\n" + verKey)
			}
		})
	}

	return
}

func TestReplaceKeyStart(t *testing.T) {
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	did, verKey, errCreateDid := CreateAndStoreDID(walletHandle, "")
	if errCreateDid != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errCreateDid)
		return
	}

	type args struct {
		Did string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"replace-key-start-works", args{Did: did}, false},
		{"replace-key-start-invalid-did", args{Did: "invalid-did"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newVerKey, errReplace := ReplaceKeyStart(walletHandle, tt.args.Did, "{}")
			hasError := errReplace != nil
			if hasError != tt.wantErr {
				t.Errorf("ReplaceKeyStart() error = '%v'", errReplace)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errReplace)
				return
			} else {
				if verKey != newVerKey {
					fmt.Println(newVerKey)
				}
			}
		})
	}

	return
}

func TestReplaceKeyApply(t *testing.T) {
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	did, _, errCreateDid := CreateAndStoreDID(walletHandle, "")
	if errCreateDid != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errCreateDid)
		return
	}

	_, errReplaceStart := ReplaceKeyStart(walletHandle, did, "{}")
	if errReplaceStart != nil {
		t.Errorf("ReplaceKeyStart() error = '%v'", errReplaceStart)
		return
	}

	type args struct {
		Did string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"replace-key-apply-works", args{Did: did}, false},
		{"replace-key-apply-invalid-did", args{Did: "invalid-did"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errReplaceApply := ReplaceKeyApply(walletHandle, tt.args.Did)
			hasError := errReplaceApply != nil
			if hasError != tt.wantErr {
				t.Errorf("ReplaceKeyApply() error = '%v'", errReplaceApply)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errReplaceApply)
				return
			}
		})
	}

	return
}

func TestStoreTheirDid(t *testing.T) {
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	identity := fmt.Sprintf(`{"did": "%s", "verkey": "GjZWsBLgZCR18aL468JAT7w9CZRiBnpxUPPgyQxh4voa"}`, didMy1)
	type args struct {
		WalletHandle int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"store-their-did-works", args{WalletHandle: walletHandle}, false},
		{"store-their-did-invalid-wallet-handle", args{WalletHandle: walletHandle + 100}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errStore := StoreTheirDid(tt.args.WalletHandle, identity)
			hasError := errStore != nil
			if hasError != tt.wantErr {
				t.Errorf("StoreTheirDid() error = '%v'", errStore)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errStore)
				return
			}
		})
	}

	return
}

func TestKeyForDid(t *testing.T) {
	poolHandle, errPool := getPoolLedger("pool")
	if errPool != nil {
		t.Errorf("getPoolLedger() error = '%v'", errPool)
		return
	}
	defer ClosePoolHandle(poolHandle)

	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	did, _, errCreateDid := CreateAndStoreDID(walletHandle, seedTrustee1)
	if errCreateDid != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errCreateDid)
		return
	}

	type args struct {
		Did string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"key-for-did-works", args{Did: did}, false},
		{"key-for-did-invalid-did", args{Did: "invalid-did"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resKey, errKey := KeyForDid(poolHandle, walletHandle, tt.args.Did)
			hasError := errKey != nil
			if hasError != tt.wantErr {
				t.Errorf("KeyForDid() error = '%v'", errKey)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errKey)
				return
			} else {
				fmt.Println(resKey)
			}
		})
	}

	return
}

func TestKeyForLocalDID(t *testing.T) {
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	did, _, errCreateDid := CreateAndStoreDID(walletHandle, seedTrustee1)
	if errCreateDid != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errCreateDid)
		return
	}

	type args struct {
		Did string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"key-for-local-did-works", args{Did: did}, false},
		{"key-for-local-did-invalid-did", args{Did: "invalid-did"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resKey, errKey := KeyForLocalDID(walletHandle, tt.args.Did)
			hasError := errKey != nil
			if hasError != tt.wantErr {
				t.Errorf("KeyForLocalDid() error = '%v'", errKey)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errKey)
				return
			} else {
				fmt.Println(resKey)
			}
		})
	}

	return
}

func TestSetEndPointForDid(t *testing.T) {
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	did, verKey, errCreateDid := CreateAndStoreDID(walletHandle, seedTrustee1)
	if errCreateDid != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errCreateDid)
		return
	}

	type args struct {
		Did string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"set-endpoint-for-works", args{Did: did}, false},
		{"set-endpoint-for-did-invalid-did", args{Did: "invalid-did"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errSet := SetEndPointForDid(walletHandle, tt.args.Did, endPoint, verKey)
			hasError := errSet != nil
			if hasError != tt.wantErr {
				t.Errorf("SetEndPointForDid() error = '%v'", errSet)
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

func TestGetEndPointForDid(t *testing.T) {
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	did, verKey, errCreateDid := CreateAndStoreDID(walletHandle, seedTrustee1)
	if errCreateDid != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errCreateDid)
		return
	}

	errSet := SetEndPointForDid(walletHandle, did, endPoint, verKey)
	if errSet != nil {
		t.Errorf("SetEndPointForDid() error = '%v'", errSet)
	}

	type args struct {
		Did string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"get-endpoint-for-works", args{Did: did}, false},
		{"get-endpoint-for-did-invalid-did", args{Did: "invalid-did"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			address, key, errGet := GetEndPointForDid(walletHandle, -1, tt.args.Did)
			hasError := errGet != nil
			if hasError != tt.wantErr {
				t.Errorf("GetEndPointForDid() error = '%v'", errGet)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errGet)
				return
			} else {
				fmt.Println(address + "\n" + key)
			}
		})
	}

	return
}

func TestSetDidMetadata(t *testing.T) {
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	did, _, errCreateDid := CreateAndStoreDID(walletHandle, seedTrustee1)
	if errCreateDid != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errCreateDid)
		return
	}

	type args struct {
		Did string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"set-did-metadata-works", args{Did: did}, false},
		{"set-did-metadata-invalid-did", args{Did: "invalid-did"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errSet := SetDidMetadata(walletHandle, tt.args.Did, metadata)
			hasError := errSet != nil
			if hasError != tt.wantErr {
				t.Errorf("SetDidMetadata() error = '%v'", errSet)
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

func TestGetDidMetadata(t *testing.T) {
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	did, _, errCreateDid := CreateAndStoreDID(walletHandle, seedTrustee1)
	if errCreateDid != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errCreateDid)
		return
	}

	errSet := SetDidMetadata(walletHandle, did, metadata)
	if errSet != nil {
		t.Errorf("SetDidMetadata() error = '%v'", errSet)
		return
	}

	type args struct {
		Did string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"get-did-metadata-works", args{Did: did}, false},
		{"get-did-metadata-invalid-did", args{Did: "invalid-did"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resMetadata, errGet := GetDidMetadata(walletHandle, tt.args.Did)
			hasError := errGet != nil
			if hasError != tt.wantErr {
				t.Errorf("GetDidMetadata() error = '%v'", errGet)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errGet)
				return
			} else {
				fmt.Println(resMetadata)
			}
		})
	}

	return
}

func TestGetDidWithMetadata(t *testing.T) {
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	did, _, errCreateDid := CreateAndStoreDID(walletHandle, seedTrustee1)
	if errCreateDid != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errCreateDid)
		return
	}

	errSet := SetDidMetadata(walletHandle, did, metadata)
	if errSet != nil {
		t.Errorf("SetDidMetadata() error = '%v'", errSet)
		return
	}

	type args struct {
		Did string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"get-did-with-metadata-works", args{Did: did}, false},
		{"get-did-with-metadata-invalid-did", args{Did: "invalid-did"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resDid, errGet := GetDidWithMetadata(walletHandle, tt.args.Did)
			hasError := errGet != nil
			if hasError != tt.wantErr {
				t.Errorf("GetDidWithMetadata() error = '%v'", errGet)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errGet)
				return
			} else {
				fmt.Println(resDid)
			}
		})
	}

	return
}

func TestListDidsWithMeta(t *testing.T) {
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	_, _, errCreateDid := CreateAndStoreDID(walletHandle, seedTrustee1)
	if errCreateDid != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errCreateDid)
		return
	}

	_, _, errCreateDid2 := CreateAndStoreDID(walletHandle, seedSteward1)
	if errCreateDid2 != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errCreateDid2)
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
		{"list-dids-with-metadata-works", args{WalletHandle: walletHandle}, false},
		{"list-dids-with-metadata-invalid-wallet-handle", args{WalletHandle: walletHandle + 100}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resList, errList := ListDidsWithMeta(tt.args.WalletHandle)
			hasError := errList != nil
			if hasError != tt.wantErr {
				t.Errorf("ListDidsWithMeta() error = '%v'", errList)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errList)
				return
			} else {
				fmt.Println(resList)
			}
		})
	}

	return
}

func TestAbbreviateVerKey(t *testing.T) {
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	did, verKey, errCreateDid := CreateAndStoreDID(walletHandle, "")
	if errCreateDid != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errCreateDid)
		return
	}

	type args struct {
		Did string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"abbreviate-key-works", args{Did: did}, false},
		{"abbreviate-key-invalid-did", args{Did: "invalid-did"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			abbvrKey, errAbbrv := AbbreviateVerKey(tt.args.Did, verKey)
			hasError := errAbbrv != nil
			if hasError != tt.wantErr {
				t.Errorf("AbbreviateVerKey() error = '%v'", errAbbrv)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errAbbrv)
				return
			} else {
				fmt.Println(abbvrKey)
			}
		})
	}

	return
}

func TestQualifyDid(t *testing.T) {
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	did, _, errCreateDid := CreateAndStoreDID(walletHandle, seedMy1)
	if errCreateDid != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errCreateDid)
		return
	}

	method := "peer"

	type args struct {
		Did string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"qualify-did-works", args{Did: did}, false},
		{"qualify-did-invalid-did", args{Did: "invalid-did"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qualifiedDid, errQualify := QualifyDid(walletHandle, did, method)
			hasError := errQualify != nil
			if hasError != tt.wantErr {
				t.Errorf("QualifyDid() error = '%v'", errQualify)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errQualify)
				return
			} else {
				fmt.Println(qualifiedDid)
			}
		})
	}

	return
}
