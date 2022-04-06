/*
// ******************************************************************
// Purpose: pairwise unit testing
// Author: angel.draghici@siemens.com
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

func TestCreatePairwise(t *testing.T) {
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("createWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	theirIdentity := fmt.Sprintf(`{"did": "%s", "verkey": "GJ1SzoWzavQYfNL9XkaJdrQejfztN4XqdsiV4ct3LXKL"}`, didTrustee)
	errStore := StoreTheirDid(walletHandle, theirIdentity)
	if errStore != nil {
		t.Errorf("StoreTheirDid() error = '%v'", errStore)
		return
	}

	myDid, _, errDid := CreateAndStoreDID(walletHandle, "")
	if errDid != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDid)
		return
	}

	type args struct {
		Did string
		Metadata string
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"test-create-pairwise", args{Did: myDid}, false},
		{"test-create-pairwise-invalid-did", args{Did: "invalid-did"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errPairwise := CreatePairwise(walletHandle, didTrustee, tt.args.Did, metadata)
			hasError := errPairwise != nil
			if hasError != tt.wantErr {
				t.Errorf("CreatePairwise() error = '%v'", errPairwise)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errPairwise)
				return
			}
		})
	}
	return
}

func TestGetPairwise(t *testing.T) {
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("createWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	theirIdentity := fmt.Sprintf(`{"did": "%s", "verkey": "GJ1SzoWzavQYfNL9XkaJdrQejfztN4XqdsiV4ct3LXKL"}`, didTrustee)
	errStore := StoreTheirDid(walletHandle, theirIdentity)
	if errStore != nil {
		t.Errorf("StoreTheirDid() error = '%v'", errStore)
		return
	}

	theirIdentity2 := fmt.Sprintf(`{"did": "%s", "verkey": "GjZWsBLgZCR18aL468JAT7w9CZRiBnpxUPPgyQxh4voa"}`, didMy1)
	errStore = StoreTheirDid(walletHandle, theirIdentity2)
	if errStore != nil {
		t.Errorf("StoreTheirDid() error = '%v'", errStore)
		return
	}

	myDid, _, errDid := CreateAndStoreDID(walletHandle, "")
	if errDid != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDid)
		return
	}
	myDid2, _, errDid := CreateAndStoreDID(walletHandle, "")
	if errDid != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDid)
		return
	}

	errCreatePairwise := CreatePairwise(walletHandle, didTrustee, myDid, "")
	if errCreatePairwise != nil {
		t.Errorf("CreatePairwise() error = '%v'", errCreatePairwise)
		return
	}
	errCreatePairwise = CreatePairwise(walletHandle, didMy1, myDid2, metadata)
	if errCreatePairwise != nil {
		t.Errorf("CreatePairwise() error = '%v'", errCreatePairwise)
		return
	}

	type args struct {
		TheirDid string
		Did string
		Metadata string
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"test-get-pairwise-works", args{TheirDid: didTrustee, Did: myDid, Metadata: ""}, false},
		{"test-get-pairwise-with-metadata", args{TheirDid: didMy1, Did: myDid2, Metadata: metadata}, false},
		{"test-get-pairwise-invalid-did", args{TheirDid: "invalid-did"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			returnedDid, errGet := GetPairwise(walletHandle, tt.args.TheirDid)
			hasError := errGet != nil
			if hasError != tt.wantErr {
				t.Errorf("GetPairwise() error = '%v'", errGet)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errGet)
				return
			}

			if returnedDid != fmt.Sprintf(`{"my_did":"%s","metadata":"%s"}`, tt.args.Did, tt.args.Metadata) {
				t.Errorf("Test failed")
				return
			}
		})
	}

	return
}

func TestListPairwise(t *testing.T) {
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("createWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	theirIdentity := fmt.Sprintf(`{"did": "%s", "verkey": "GJ1SzoWzavQYfNL9XkaJdrQejfztN4XqdsiV4ct3LXKL"}`, didTrustee)
	errStore := StoreTheirDid(walletHandle, theirIdentity)
	if errStore != nil {
		t.Errorf("StoreTheirDid() error = '%v'", errStore)
		return
	}

	myDid, _, errDid := CreateAndStoreDID(walletHandle, "")
	if errDid != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDid)
		return
	}

	errCreatePairwise := CreatePairwise(walletHandle, didTrustee, myDid, "")
	if errCreatePairwise != nil {
		t.Errorf("CreatePairwise() error = '%v'", errCreatePairwise)
		return
	}

	expectedList := fmt.Sprintf(`["{\"my_did\":\"%s\",\"their_did\":\"%s\",\"metadata\":\"\"}"]`, myDid, didTrustee)

	type args struct {
		WalletHandle int
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"test-list-pairwise-works", args{WalletHandle: walletHandle}, false},
		{"test-list-pairwise-invalid-wallet-handle", args{WalletHandle: 100}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list, errList := ListPairwise(tt.args.WalletHandle)
			hasError := errList != nil
			if hasError != tt.wantErr {
				t.Errorf("ListPairwise() error = '%v'", errList)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errList)
				return
			}

			if list != expectedList {
				t.Errorf("Test failed")
			}
		})
	}

	return
}

func TestIsPairwiseExists(t *testing.T) {
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("createWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	theirIdentity := fmt.Sprintf(`{"did": "%s", "verkey": "GJ1SzoWzavQYfNL9XkaJdrQejfztN4XqdsiV4ct3LXKL"}`, didTrustee)
	errStore := StoreTheirDid(walletHandle, theirIdentity)
	if errStore != nil {
		t.Errorf("StoreTheirDid() error = '%v'", errStore)
		return
	}

	myDid, _, errDid := CreateAndStoreDID(walletHandle, "")
	if errDid != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDid)
		return
	}

	errCreatePairwise := CreatePairwise(walletHandle, didTrustee, myDid, "")
	if errCreatePairwise != nil {
		t.Errorf("CreatePairwise() error = '%v'", errCreatePairwise)
		return
	}

	type args struct {
		WalletHandle int
		TheirDid string
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"test-is-pairwise-exists-works", args{WalletHandle: walletHandle, TheirDid: didTrustee}, false},
		{"test-is-pairwise-exists-not-created", args{WalletHandle: walletHandle, TheirDid: "CnEDk9HrMnmiHXEV1WFgbVCRteYnPqsJwrTdcZaNhFVW"}, false},
		{"test-is-pairwise-exists-invalid-wallet-handle", args{WalletHandle: 100}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists, errPairwiseExists := IsPairwiseExists(walletHandle, tt.args.TheirDid)
			hasError := errPairwiseExists != nil
			if hasError != tt.wantErr {
				t.Errorf("IsPairwiseExists() error = '%v'", errPairwiseExists)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errPairwiseExists)
				return
			}

			if tt.args.TheirDid == didTrustee {
				if !exists {
					t.Errorf("Test failed")
				}
			} else {
				if exists {
					t.Errorf("Test failed")
				}
			}

		})
	}
	return
}

func TestSetPairwiseMetadata(t *testing.T) {
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("createWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	theirIdentity := fmt.Sprintf(`{"did": "%s", "verkey": "GJ1SzoWzavQYfNL9XkaJdrQejfztN4XqdsiV4ct3LXKL"}`, didTrustee)
	errStore := StoreTheirDid(walletHandle, theirIdentity)
	if errStore != nil {
		t.Errorf("StoreTheirDid() error = '%v'", errStore)
		return
	}

	theirIdentity2 := fmt.Sprintf(`{"did": "%s", "verkey": "GjZWsBLgZCR18aL468JAT7w9CZRiBnpxUPPgyQxh4voa"}`, didMy1)
	errStore = StoreTheirDid(walletHandle, theirIdentity2)
	if errStore != nil {
		t.Errorf("StoreTheirDid() error = '%v'", errStore)
		return
	}

	myDid, _, errDid := CreateAndStoreDID(walletHandle, "")
	if errDid != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDid)
		return
	}

	errCreatePairwise := CreatePairwise(walletHandle, didTrustee, myDid, "")
	if errCreatePairwise != nil {
		t.Errorf("CreatePairwise() error = '%v'", errCreatePairwise)
		return
	}

	pairwiseWithoutMetadata, errGetPairwise := GetPairwise(walletHandle, didTrustee)
	if errGetPairwise != nil {
		t.Errorf("GetPairwise() error = '%v'", errGetPairwise)
		return
	}

	type args struct {
		TheirDid string
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"test-set-pairwise-metadata-works", args{TheirDid: didTrustee}, false},
		{"test-set-pairwise-metadata-not-created-pairwise", args{TheirDid: didMy1}, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errSet := SetPairwiseMetadata(walletHandle, tt.args.TheirDid, metadata)
			hasError := errSet != nil
			if hasError != tt.wantErr {
				t.Errorf("SetPairwiseMetadata() error = '%v'", errSet)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errSet)
				return
			}
			pairwiseWithMetadata, errGet := GetPairwise(walletHandle, tt.args.TheirDid)
			if errGet != nil {
				t.Errorf("GetPairwise() error = '%v'", errGet)
				return
			}

			if pairwiseWithMetadata != pairwiseWithoutMetadata {
				if pairwiseWithMetadata != fmt.Sprintf(`{"my_did":"%s","metadata":"%s"}`, myDid, metadata) {
					t.Errorf("Test failed")
				}
			} else {
				t.Errorf("Test failed")
			}
		})
	}
	
	return
}