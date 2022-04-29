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
	"github.com/Jeffail/gabs/v2"
	"github.com/joyride9999/IndySdkGoBindings/inMemUtils"
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

			// Both StoredMetadata and StorageHandles shouldn't be null if create and open wallet succeeded.
			if len(customStorage.MetadataHandles) == 0 || len(customStorage.StorageHandles) == 0 {
				t.Errorf("Test failed")
			}
		})
	}
}

func TestMemStorageNonSecretsRecords(t *testing.T) {
	fmt.Println("InMemoryStorage: Non Secrets Records")

	// Initialize in memory storage.
	storageType := "in-mem-storage"
	customStorage := inMemUtils.NewInMemoryStorage()

	errStorage := RegisterWalletStorage(storageType, customStorage)
	if errStorage != nil {
		t.Errorf("RegisterWalletStorage() error = '%v'", errStorage)
		return
	}

	// Create wallet.
	walletConfig := wallet.Config{
		ID:            "custom",
		StorageType:   storageType,
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

	// Add wallet record and retrieve it from wallet storage.
	errAdd := IndyAddWalletRecord(walletHandle, recordType, recordId1, recordValue1, recordTags1)
	if errAdd != nil {
		t.Errorf("IndyAddWalletRecord() error = '%v'", errAdd)
		return
	}

	record, errGet := IndyGetWalletRecord(walletHandle, recordType, recordId1, recordOptions)
	if errGet != nil {
		t.Errorf("IndyGetWalletRecord() error = '%v'", errGet)
		return
	}

	expected := `{"id": "recordId1", "value": "recordValue", "tags": {"tagName1":"str1","tagName2":"5","tagName3":"12"}, "type": "testType"}`
	expectedRecord, _ := gabs.ParseJSON([]byte(expected))

	returnedRecord, errGabs2 := gabs.ParseJSON([]byte(record))
	if errGabs2 != nil {
		t.Errorf("Gabs Parse error = '%v'", errGabs2)
		return
	}

	// Check if retrieved record is the expected one.
	if returnedRecord.Path("tags").String() != expectedRecord.Path("tags").String() ||
		returnedRecord.Path("type").String() != expectedRecord.Path("type").String() {
		t.Errorf("Test failed")
	}

	// Add another tag to the record.
	errAddTags := IndyAddWalletRecordTags(walletHandle, recordType, recordId1, `{"tagName4": "str4"}`)
	if errAddTags != nil {
		t.Errorf("IndyAddWalletRecordTags() error = '%v'", errAddTags)
		return
	}

	record, errGet = IndyGetWalletRecord(walletHandle, recordType, recordId1, recordOptions)
	if errGet != nil {
		t.Errorf("IndyGetWalletRecord() error = '%v'", errGet)
		return
	}

	// Check if retrieved record tags are different from the initial one.
	returnedRecord, _ = gabs.ParseJSON([]byte(record))
	if returnedRecord.Path("tags").String() == recordTags1 {
		t.Errorf("Test failed")
		return
	}

	// Add another wallet record.
	errAdd = IndyAddWalletRecord(walletHandle, recordType, recordId2, recordValue2, recordTags1)
	if errAdd != nil {
		t.Errorf("IndyAddWalletRecord() error = '%v'", errAdd)
		return
	}

	// Open search.
	searchHandle, errOpenSearch := IndyOpenWalletSearch(walletHandle, recordType, `{"tagName1": "str1"}`, `{"retrieveRecords": true}`)
	if errOpenSearch != nil {
		t.Errorf("IndyOpenWalletSearch() error = '%v'", errOpenSearch)
		return
	}
	defer IndyCloseWalletSearch(searchHandle)

	// Fetch wallet records.
	searchRecords, errFetch := IndyFetchWalletSearchNextRecords(walletHandle, searchHandle, int32(2))
	if errFetch != nil {
		t.Errorf("IndyFetchWalletSearchNextRecords() error = '%v'", errFetch)
		return
	}

	// Check fetched records to match the expected ones.
	expectedRecords := `{"records": 
		[
			{"id":"recordId1","value":"recordValue","tags": null, "type": null}, 
			{"id":"recordId2","value":"recordValue2","tags": null, "type": null}
		]
	}`
	expectedRecordsParsed, _ := gabs.ParseJSON([]byte(expectedRecords))
	searchRecordsParsed, _ := gabs.ParseJSON([]byte(searchRecords))
	for _, search := range searchRecordsParsed.S("records").Children() {
		ok := false
		for _, expected := range expectedRecordsParsed.S("records").Children() {
			if search.String() == expected.String() {
				ok = true
			}
			if ok == true {
				break
			}
		}
		if !ok {
			t.Errorf("Test failed")
			break
		}
	}

	return
}

func TestMemStorageCredentialSearchProof(t *testing.T) {
	fmt.Println("InMemoryStorage: Credentials Test")

	// Register in memory storage for issuer.
	inMemStorageType := "memoryStorage"
	storage := inMemUtils.NewInMemoryStorage()

	errStorage := RegisterWalletStorage(inMemStorageType, storage)
	if errStorage != nil {
		t.Errorf("RegisterWalletStorage() error = '%v'", errStorage)
		return
	}

	// Create issuer wallet.
	issuerConfig := wallet.Config{
		ID:            "issuer",
		StorageType:   inMemStorageType,
		StorageConfig: wallet.StorageConfig{Path: ".\\out\\wallets", Dsn: "", LogSql: 0},
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

	didIssuer, _, errDidIssuer := CreateAndStoreDID(whIssuer, "000000000000000000000000Trustee1")
	if errDidIssuer != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDidIssuer)
		return
	}

	// Create holder wallet
	holderConfig := wallet.Config{
		ID:            "holder",
		StorageType:   inMemStorageType,
		StorageConfig: wallet.StorageConfig{Path: ".\\out\\wallets", Dsn: "", LogSql: 0},
	}
	holderCredential := wallet.Credential{
		Key: "123",
	}
	errCreate = CreateWallet(holderConfig, holderCredential)
	if errCreate != nil {
		t.Errorf("CreateWallet() error = '%v'", errCreate)
		return
	}

	whHolder, errOpen := OpenWallet(holderConfig, holderCredential)
	if errOpen != nil {
		t.Errorf("OpenWallet() error = '%v'", errOpen)
		return
	}
	defer walletCleanup(whHolder, holderConfig, holderCredential)

	// Get did for holder
	didHolder, _, errDidH := CreateAndStoreDID(whHolder, "")
	if errDidH != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDidH)
		return
	}

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
	_, errCredStore := ProverStoreCredential(whHolder, "", credentialRequestMetadata, credJson, credDefJson, "")
	if errCred != nil {
		t.Errorf("ProverStoreCredential() error = '%v'", errCredStore)
	}

	nonce, _ := GenerateNonce()
	proofRequest := fmt.Sprintf(`{"nonce": "%s", "name": "proofRequest", "ver" : "1.0", "version": "0.1", "requested_attributes" : {"attr1_referent" : {"name": "a1"}}, "requested_predicates" : {"predicate1_referent" : { "name": "a2", "p_type": ">=", "p_value" : 2}} }`,
		nonce)

	// Search for data
	searchHandle, errSearch := ProverSearchForCredentialForProofReq(whHolder, proofRequest, "")
	if errSearch != nil {
		t.Errorf("ProverSearchForCredentialProofReq() error = '%v'", errSearch)
	}

	if searchHandle > 0 {
		defer ProverCloseCredentialsSearchForProofReq(searchHandle)
	}

	// Gets the data for attr1_referent
	credForAttr1, errFetchCredential := ProverFetchCredentialsForProofReq(searchHandle, "attr1_referent", 10)
	if errFetchCredential != nil {
		t.Errorf("ProverFetchCredentialsForProofReq() error = '%v'", errFetchCredential)
	}

	// Gets the data for predicate1_referent
	credForPred1, errFetchCredential2 := ProverFetchCredentialsForProofReq(searchHandle, "predicate1_referent", 10)
	if errFetchCredential2 != nil {
		t.Errorf("ProverFetchCredentialsForProofReq() error = '%v'", errFetchCredential2)
	}

	// Read credential definition id from fetched data
	credInfoAttrJson, errGabs := gabs.ParseJSON([]byte(credForAttr1))
	if errGabs != nil {
		t.Errorf("Gabs ParseJSON() error = '%v'", errGabs)
		return
	}
	attrCredId := credInfoAttrJson.Path("0.cred_info.referent").Data()

	credInfoPredJson, errGabs := gabs.ParseJSON([]byte(credForPred1))
	if errGabs != nil {
		t.Errorf("Gabs ParseJSON() error = '%v'", errGabs)
		return
	}
	predCredId := credInfoPredJson.Path("0.cred_info.referent").Data()

	requestedCredJson := fmt.Sprintf(`{ "self_attested_attributes" : {}, "requested_attributes" : { "attr1_referent" : { "cred_id": "%s", "revealed": true  }   }, "requested_predicates" : { "predicate1_referent" : { "cred_id": "%s" }   }     }`, attrCredId, predCredId)
	schemasJson := fmt.Sprintf(`{"%s":%s}`, schemaId, schemaJson)
	credDefsJson := fmt.Sprintf(`{"%s":%s}`, credDefId, credDefJson)

	// Creates the proof
	proofJson, errProof := ProverCreateProof(whHolder, proofRequest, requestedCredJson, mSecretId, schemasJson, credDefsJson, "{}")
	if errProof != nil {
		t.Errorf("createCredentialDefinition() error = '%v'", errProof)
		return
	}

	if len(proofJson) == 0 {
		t.Errorf("Test failed, incorrect value")
	}

	return
}
