/*
// ******************************************************************
// Purpose: in memory wallet unit testing
// Author:  angel.draghici@siemens.com
// Notes:
// Copyright (c): Siemens SRL
// This work is licensed under the terms of the Apache License Version 2.0.  See
// the LICENSE.txt file in the top-level directory.
// ******************************************************************
*/

package indySDK

import (
	"fmt"
	"github.com/joyride9999/IndySdkGoBindings/inMemUtils"
	"github.com/joyride9999/IndySdkGoBindings/indyUtils"
	"github.com/joyride9999/IndySdkGoBindings/wallet"
	"testing"
)

func TestRegisterWalletMemoryStorage(t *testing.T) {
	// Initialize in-memory storage.
	customStorage := inMemUtils.NewInMemoryStorage()

	type args struct {
		StorageType string
		Storage     *inMemUtils.InMemoryStorage
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test-register-wallet-storage-in-memory-custom-type", args{StorageType: "in-mem-storage", Storage: customStorage}, false},
		{"test-register-already-existing-storage", args{StorageType: "in-mem-storage", Storage: customStorage}, true},
		{"test-register-wallet-storage-in-memory-invalid-type", args{StorageType: ""}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errStorage := RegisterWalletStorage(tt.args.StorageType, tt.args.Storage)
			hasError := errStorage != nil
			if hasError != tt.wantErr {
				t.Errorf("RegisterWalletStorage() error = '%v'", errStorage)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errStorage)
				return
			}

			// Create wallet with the registered in-memory storage type.
			walletConfig := wallet.Config{
				ID:            "custom",
				StorageType:   tt.args.StorageType,
				StorageConfig: wallet.StorageConfig{Path: ".\\out\\wallets", Dsn: "", LogSql: 0},
			}
			walletCredential := wallet.Credential{
				Key: "123",
			}
			errCreate := CreateWallet(walletConfig, walletCredential)
			if errCreate != nil {
				t.Errorf("CreateWallet() error = '%v'", errCreate)
				return
			}

			walletHandle, errOpen := OpenWallet(walletConfig, walletCredential)
			if errOpen != nil {
				t.Errorf("OpenWallet() error = '%v'", errOpen)
				return
			}

			defer walletCleanup(walletHandle, walletConfig, walletCredential)

			// Both StoredMetadata and WalletHandles shouldn't be null if create and open wallet succeeded.
			if len(customStorage.MetadataHandles) == 0 || len(customStorage.WalletHandles) == 0 {
				t.Errorf("Test failed")
			}

		})
	}
}

func TestMemStorageCredentialSearchProof(t *testing.T) {
	fmt.Println("InMemoryStorage: Credentials Test")

	// Register in memory storage for issuer.
	iMemStorageType := "iMemoryStorage"
	iStorage := inMemUtils.NewInMemoryStorage()

	errStorage := RegisterWalletStorage(iMemStorageType, iStorage)
	if errStorage != nil {
		t.Errorf("RegisterWalletStorage() error = '%v'", errStorage)
		return
	}

	// Create trustee wallet
	issuerCfg := wallet.Config{
		ID:            "issuer",
		StorageType:   iMemStorageType,
		StorageConfig: wallet.StorageConfig{},
	}
	issuerCredential := wallet.Credential{
		Key: "123",
	}
	errCreateWalletIssuer := CreateWallet(issuerCfg, issuerCredential)
	if errCreateWalletIssuer != nil && errCreateWalletIssuer.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreateWalletIssuer)
		return
	}

	whIssuer, errOpenTr := OpenWallet(issuerCfg, issuerCredential)
	if errOpenTr != nil {
		t.Errorf("OpenWallet() error = '%v'", errOpenTr)
		return
	}
	defer walletCleanup(whIssuer, issuerCfg, issuerCredential)

	// Get did for issuer
	didIssuer, _, errDidIssuer := CreateAndStoreDID(whIssuer, "000000000000000000000000Trustee1")
	if errDidIssuer != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDidIssuer)
		return
	}

	// holder wallet
	holderCfg := wallet.Config{
		ID:            "holder",
		StorageType:   iMemStorageType,
		StorageConfig: wallet.StorageConfig{},
	}
	holderCredential := wallet.Credential{
		Key: "123",
	}

	errCreateWalletHolder := CreateWallet(holderCfg, holderCredential)
	if errCreateWalletHolder != nil && errCreateWalletHolder.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreateWalletHolder)
		return
	}

	whHolder, errOpenHolder := OpenWallet(holderCfg, holderCredential)
	if errOpenHolder != nil {
		t.Errorf("OpenWallet() error = '%v'", errOpenHolder)
		return
	}

	defer walletCleanup(whHolder, holderCfg, holderCredential)

	// Get did for holder
	didHolder, _, errDidH := CreateAndStoreDID(whHolder, "")
	if errDidH != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDidH)
		return
	}

	// Creates a schema
	schemaId, schemaJson, errSchema := IssuerCreateSchema(didIssuer, "testSch", "0.1", `["a1", "a2", "a3"]`)
	if errSchema != nil {
		t.Errorf("IssuerCreateSchema() error = '%v'", errSchema)
		return
	}

	// Creates credential definition
	credDefId, credDefJson, errCredDef := IssuerCreateAndStoreCredentialDefinition(whIssuer, didIssuer, schemaJson, "testCredDef", "CL", `{"support-revocation": false}`)
	if errCredDef != nil {
		t.Errorf("IssuerCreateAndStoreCredentialDefinition() error = '%v'", errCredDef)
	}

	// Create credential offer
	credOffer, errCredOffer := IssuerCreateCredentialOffer(whIssuer, credDefId)
	if errCredOffer != nil {
		t.Errorf("IssuerCreateCredentialOffer() error = '%v'", errCredOffer)
	}

	// Credential request from offer
	mSecretId, errMSid := ProverCreateMasterSecret(whHolder, "")
	if errMSid != nil {
		t.Errorf("ProverCreateMasterSecret() error = '%v'", errMSid)
	}

	credentialRequest, credentialRequestMetadata, errCredReq := ProverCreateCredentialRequest(whHolder, didHolder, credOffer, credDefJson, mSecretId)
	if errCredReq != nil {
		t.Errorf("ProverCreateMasterSecret() error = '%v'", errCredReq)
	}

	// Credential data
	credentialValues := `{"a1" : { "raw": "1", "encoded": "1"}, "a2" : { "raw": "2", "encoded": "2"}, "a3" : { "raw": "3", "encoded": "3"}}`

	// Creates the credential
	credJson, _, _, errCred := IssuerCreateCredential(whIssuer, credOffer, credentialRequest, credentialValues, "", 0)
	if errCred != nil {
		t.Errorf("IssuerCreateCredential() error = '%v'", errCred)
	}

	// Store the credential into prover wallet
	credentialId, errCredStore := ProverStoreCredential(whHolder, "", credentialRequestMetadata, credJson, credDefJson, "")
	if errCred != nil {
		t.Errorf("ProverStoreCredential() error = '%v'", errCredStore)
	}

	nonce, _ := GenerateNonce()
	proofRequest := fmt.Sprintf(`{"nonce": "%s", "name": "proofRequest", "ver" : "1.0", "version": "0.1", "requested_attributes" : {"attr1_referent" : {"name": "a1"}}, "requested_predicates" : {"predicate1_referent" : { "name": "a2", "p_type": ">=", "p_value" : 2}} }`,
		nonce)

	// Search for data
	searchHandle, errSearch := ProverSearchForCredentialForProofReq(whHolder, proofRequest, "")
	if errSearch != nil {
		t.Errorf("ProverStoreCredential() error = '%v'", errSearch)
	}

	if searchHandle > 0 {
		defer ProverCloseCredentialsSearchForProofReq(searchHandle)
	}

	// Gets the data for attr1_referent
	credForAttr1, errFetchCredential := ProverFetchCredentialsForProofReq(searchHandle, "attr1_referent", 10)
	if errFetchCredential != nil {
		t.Errorf("ProverFetchCredentialsForProofReq() error = '%v'", errFetchCredential)
	}
	credForAttr1 = credForAttr1

	// Gets the data for predicate1_referent
	credForPred1, errFetchCredential2 := ProverFetchCredentialsForProofReq(searchHandle, "predicate1_referent", 10)
	if errFetchCredential2 != nil {
		t.Errorf("ProverFetchCredentialsForProofReq() error = '%v'", errFetchCredential2)
	}
	credForPred1 = credForPred1

	// TODO: read credential definition id from fetched data
	requestedCredJson := fmt.Sprintf(`{ "self_attested_attributes" : {}, "requested_attributes" : { "attr1_referent" : { "cred_id": "%s", "revealed": true  }   }, "requested_predicates" : { "predicate1_referent" : { "cred_id": "%s" }   }     }`, credentialId, credentialId)
	schemasJson := fmt.Sprintf(`{"%s":%s}`, schemaId, schemaJson)
	credDefsJson := fmt.Sprintf(`{"%s":%s}`, credDefId, credDefJson)

	// Creates the proof
	proofJson, errProof := ProverCreateProof(whHolder, proofRequest, requestedCredJson, mSecretId, schemasJson, credDefsJson, "{}")
	if errProof != nil {
		t.Errorf("createCredentialDefinition() error = '%v'", errProof)
		return
	}

	// Verifies the proof
	valid, errVer := VerifierVerifyProof(proofRequest, proofJson, schemasJson, credDefsJson, "{}", "{}")
	if errVer != nil {
		t.Errorf("createCredentialDefinition() error = '%v'", errVer)
		return
	}
	if valid {
		t.Log("Proof is valid\n")
	} else {
		t.Log("Proof is not valid\n")
	}
	t.Log("Success\n")

	return
}

func TestMemStorageNonSecrets(t *testing.T) {
	iMemStorageType := "iMemoryStorage"
	iStorage := inMemUtils.NewInMemoryStorage()

	errStorage := RegisterWalletStorage(iMemStorageType, iStorage)
	if errStorage != nil {
		t.Errorf("RegisterWalletStorage() error = '%v'", errStorage)
		return
	}

	issuerConfig := wallet.Config{
		ID:            "issuer",
		StorageType:   iMemStorageType,
		StorageConfig: wallet.StorageConfig{Path: "indySDK\\out\\wallets", Dsn: "", LogSql: 0},
	}
	issuerCredential := wallet.Credential{
		Key: "123",
	}
	errCreate := CreateWallet(issuerConfig, issuerCredential)
	if errCreate != nil {
		t.Errorf("CreateWallet() error = '%v'", errCreate)
		return
	}

	whIssuer, errOpen := OpenWallet(issuerConfig, issuerCredential)
	if errOpen != nil {
		t.Errorf("OpenWallet() error = '%v'", errOpen)
		return
	}
	defer walletCleanup(whIssuer, issuerConfig, issuerCredential)

	errAdd := IndyAddWalletRecord(whIssuer, "type1", "id1", "value1", `{"name": "Customer"}`)
	if errAdd != nil {
		t.Errorf("IndyAddWalletRecord() error = '%v'", errAdd)
		return
	}

	errAdd = IndyAddWalletRecord(whIssuer, "type2", "id1", "value1", `{"name": "Customer"}`)
	if errAdd != nil {
		t.Logf("EXPECTED IndyAddWalletRecord() error = '%v'", errAdd)
	}

	record, errGet := IndyGetWalletRecord(whIssuer, "type1", "id1", recordOptions)
	if errGet != nil {
		t.Errorf("IndyGetWalletRecord() error = '%v'", errGet)
		return
	}
	fmt.Println("AddWalletRecord(): " + record)

	errAddT := IndyAddWalletRecordTags(whIssuer, "type1", "id1", `{"state": "California"}`)
	if errAddT != nil {
		t.Errorf("IndyAddWalletRecordTags() error = '%v'", errAddT)
		return
	}

	record, errGet = IndyGetWalletRecord(whIssuer, "type1", "id1", recordOptions)
	if errGet != nil {
		t.Errorf("IndyGetWalletRecord() error = '%v'", errGet)
		return
	}
	fmt.Println("AddWalletRecordTags(): " + record)

	errDelT := IndyDeleteWalletRecordTags(whIssuer, "type1", "id1", `["name"]`)
	if errDelT != nil {
		t.Errorf("IndyDeleteWalletRecordTags() error = '%v'", errDelT)
		return
	}

	record, errGet = IndyGetWalletRecord(whIssuer, "type1", "id1", recordOptions)
	if errGet != nil {
		t.Errorf("IndyGetWalletRecord() error = '%v'", errGet)
		return
	}
	fmt.Println("DeleteWalletRecordTags(): " + record)

	errUpdV := IndyUpdateWalletRecordValue(whIssuer, "type1", "id1", "value-edited1")
	if errUpdV != nil {
		t.Errorf("IndyUpdateWalletRecordValue() error = '%v'", errUpdV)
		return
	}

	record, errGet = IndyGetWalletRecord(whIssuer, "type1", "id1", recordOptions)
	if errGet != nil {
		t.Errorf("IndyGetWalletRecord() error = '%v'", errGet)
		return
	}
	fmt.Println("UpdateWalletRecordValue(): " + record)

	errUpdT := IndyUpdateWalletRecordTags(whIssuer, "type1", "id1", `{"name": "Customer", "state": "Arizona"}`)
	if errUpdT != nil {
		t.Errorf("IndyUpdateWalletRecordTags() error = '%v'", errUpdT)
		return
	}

	record, errGet = IndyGetWalletRecord(whIssuer, "type1", "id1", recordOptions)
	if errGet != nil {
		t.Errorf("IndyGetWalletRecord() error = '%v'", errGet)
		return
	}
	fmt.Println("UpdateWalletRecordTags(): " + record)

	errDel := IndyDeleteWalletRecord(whIssuer, "type1", "id1")
	if errDel != nil {
		t.Errorf("IndyDeleteWalletRecord() error = '%v'", errDel)
		return
	}

	record, errGet = IndyGetWalletRecord(whIssuer, "type1", "id1", recordOptions)
	if errGet != nil { // expected to fail
		t.Logf("EXPECTED IndyGetWalletRecord() error = '%v'", errGet)
	}

	errAdd = IndyAddWalletRecord(whIssuer, "grading", "recordId-1", "value", `{"~t1": "1", "~t2": "v2", "~t3": "v3", "~t4": "v4"}`)
	if errAdd != nil {
		t.Errorf("IndyAddWalletRecord() error = '%v'", errAdd)
		return
	}

	errAdd = IndyAddWalletRecord(whIssuer, "grading", "recordId-2", "value", `{"~t1": "1", "~t2": "v2.2", "~t3": "v3", "~t4": "v4"}`)
	if errAdd != nil {
		t.Errorf("IndyAddWalletRecord() error = '%v'", errAdd)
		return
	}

	errAdd = IndyAddWalletRecord(whIssuer, "grading", "recordId-3", "value", `{"~t1": "1", "~t2": "v2", "~t3": "v3.3", "~t4": "v4"}`)
	if errAdd != nil {
		t.Errorf("IndyAddWalletRecord() error = '%v'", errAdd)
		return
	}

	query := `{
		"$and": [
			{
            	"~t1": {"$gte": "1"}
        	},
			{
            	"~t2": "v2"
        	},
        	{
				"$not": {
					"~t3": "v3.3"
				}
    		}
    	]
	}`

	sh, errSearch := IndyOpenWalletSearch(whIssuer, recordType, query, recordOptions)
	if errSearch != nil {
		t.Errorf("IndyOpenWalletSearch() error = '%v'", errSearch)
		return
	}

	recordsJson, errFetch := IndyFetchWalletSearchNextRecords(whIssuer, sh, 3)
	if errFetch != nil {
		t.Errorf("IndyFetchWalletSearchNextRecords() error = '%v'", errFetch)
		return
	}

	errClose := IndyCloseWalletSearch(sh)
	if errClose != nil {
		t.Errorf("IndyCloseWalletSearch() error = '%v'", errClose)
		return
	}

	fmt.Println(recordsJson)
}
