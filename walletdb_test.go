/*
// ******************************************************************
// Purpose: wallet db unit testing
// Author: alexandru.leonte@siemens.com
// Notes:
// Copyright (c): Siemens SRL
// This work is licensed under the terms of the Apache License Version 2.0.  See
// the LICENSE.txt file in the top-level directory.
// ******************************************************************
*/

package indySDK

import (
	"fmt"
	"github.com/google/uuid"
	"indySDK/dbUtils"
	"indySDK/indyUtils"
	"indySDK/wallet"
	"testing"
)

func TestPgStorageAddGetRecords(t *testing.T) {

	fmt.Println("PostgreSQL Database: Start")
	dsn := "host=localhost user=user password=password dbname=wallets port=5432 sslmode=disable"
	pgMultiSchemaStorage := dbutils.NewPgMultiSchemaStorage()
	storageType := "postgresdb"
	errStorage := RegisterWalletStorage(storageType, pgMultiSchemaStorage)
	if errStorage != nil {
		t.Errorf("RegisterWalletTypet() error = '%v'", errStorage)
		return
	}

	// Create trustee wallet
	trusteeConfig := wallet.Config{
		ID:            "trustee",
		StorageType:   storageType,
		StorageConfig: wallet.StorageConfig{Dsn: dsn, LogSql: 4},
	}
	trusteeCredential := wallet.Credential{
		Key: "123",
	}
	errCreateWalletTrust := CreateWallet(trusteeConfig, trusteeCredential)
	if errCreateWalletTrust != nil && errCreateWalletTrust.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreateWalletTrust)
		return
	}

	whTrustee, errOpenTr := OpenWallet(trusteeConfig, trusteeCredential)
	if errOpenTr != nil {
		t.Errorf("OpenWallet() error = '%v'", errOpenTr)
		return
	}

	defer CloseWallet(whTrustee)

	// Get did for trustee
	_, _, errDidT := CreateAndStoreDID(whTrustee, "000000000000000000000000Trustee1")
	if errDidT != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDidT)
		return
	}

	recordType := "test"
	recordValue := "testvalue"
	id, _ := uuid.NewUUID()
	errPR := IndyAddWalletRecord(whTrustee, recordType, id.String(), recordValue, "{\"tagName1\": \"t1\", \"~tagName2\": \"t2\"}")
	if errPR != nil {
		t.Errorf("IndyAddWalletRecord() error = '%v'", errPR)
		return
	}

	record, errRecord := IndyGetWalletRecord(whTrustee, recordType, id.String(), "{\"retrieveTags\": true, \"retrieveType\":true}")
	if errRecord != nil {
		t.Errorf("IndyAddWalletRecord() error = '%v'", errRecord)
		return
	}
	//TODO check values are equal
	record = record

	return
}

func TestPgStorageCredentialSearchProof(t *testing.T) {

	fmt.Println("PostgreSQL Database: Start")
	dsn := "host=localhost user=user password=password dbname=wallets port=5432 sslmode=disable"
	pgMultiSchemaStorage := dbutils.NewPgMultiSchemaStorage()
	storageType := "postgresdb"
	errStorage := RegisterWalletStorage(storageType, pgMultiSchemaStorage)
	if errStorage != nil {
		t.Errorf("RegisterWalletTypet() error = '%v'", errStorage)
		return
	}

	// Create trustee wallet
	issuerConfig := wallet.Config{
		ID:            "issuer",
		StorageType:   storageType,
		StorageConfig: wallet.StorageConfig{Dsn: dsn, LogSql: 4},
	}
	issuerCredential := wallet.Credential{
		Key: "123",
	}
	errCreateWalletIssuer := CreateWallet(issuerConfig, issuerCredential)
	if errCreateWalletIssuer != nil && errCreateWalletIssuer.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreateWalletIssuer)
		return
	}

	whIssuer, errOpenTr := OpenWallet(issuerConfig, issuerCredential)
	if errOpenTr != nil {
		t.Errorf("OpenWallet() error = '%v'", errOpenTr)
		return
	}

	defer walletCleanup(whIssuer, issuerConfig, issuerCredential)

	// Get did for issuer
	didIssuer, _, errDidIssuer := CreateAndStoreDID(whIssuer, "000000000000000000000000Trustee1")
	if errDidIssuer != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDidIssuer)
		return
	}

	// holder wallet
	holderConfig := wallet.Config{
		ID:            "holder",
		StorageType:   "postgresdb",
		StorageConfig: wallet.StorageConfig{Dsn: dsn, LogSql: 4},
	}
	holderCredential := wallet.Credential{
		Key: "123",
	}

	errCreateWalletHolder := CreateWallet(holderConfig, holderCredential)
	if errCreateWalletHolder != nil && errCreateWalletHolder.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreateWalletHolder)
		return
	}

	whHolder, errOpenHolder := OpenWallet(holderConfig, holderCredential)
	if errOpenHolder != nil {
		t.Errorf("OpenWallet() error = '%v'", errOpenHolder)
		return
	}

	defer walletCleanup(whHolder, holderConfig, holderCredential)

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
