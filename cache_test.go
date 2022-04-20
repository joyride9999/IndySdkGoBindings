/*
// ******************************************************************
// Purpose: cache unit testing
// Author: angel.draghici@siemens.com
// Notes:
// Copyright (c): Siemens SRL
// This work is licensed under the terms of the Apache License Version 2.0.  See
// the LICENSE.txt file in the top-level directory.
// ******************************************************************
*/

package indySDK

import (
	"github.com/joyride9999/IndySdkGoBindings/indyUtils"
	"fmt"
	"testing"
)

func TestCacheGetSchema(t *testing.T) {
	// Create pool for issuer
	poolHandle, errPool := getPoolLedger("pool")
	if errPool != nil {
		t.Errorf("getPoolLedger() error = '%v'", errPool)
		return
	}
	defer ClosePoolHandle(poolHandle)

	// Create wallet for issuer
	whIssuer, errCreate := createWallet(issuerConfig(), issuerCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(whIssuer, issuerConfig(), issuerCredentials())

	// Get did for issuer
	didIssuer, _, errDid := CreateAndStoreDID(whIssuer, seedSteward1)
	if errDid != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDid)
		return
	}

	// Create schema for issuer
	schemaId, _, errSchema := IssuerCreateSchema(didIssuer, "gvt", "1.0", schemaAttributes)
	if errSchema != nil {
		t.Errorf("IssuerCreateSchema() error = '%v'", errSchema)
		return
	}

	type args struct {
		SchemaId string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"get-schema-cache", args{SchemaId: schemaId}, false},
		{"get-schema-cache-invalid-id", args{SchemaId: "invalid-schema"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schemaRetrieved, errGet := GetCacheSchema(poolHandle, whIssuer, didIssuer, tt.args.SchemaId, `{}`)
			hasError := errGet != nil
			if hasError != tt.wantErr {
				if errGet.Error() == indyUtils.GetIndyError(309) {
					t.Log(errGet)
				} else {
					t.Errorf("GetCacheSchema() error = '%v'", errGet)
					return
				}
			}

			if tt.wantErr {
				t.Log(errGet)
			} else {
				fmt.Println(schemaRetrieved)
			}
		})
	}

	return
}

func TestCacheGetCredDef(t *testing.T) {
	// Create pool
	poolHandle, errPool := getPoolLedger("pool")
	if errPool != nil {
		t.Errorf("getPoolLedger() error = '%v'", errPool)
		return
	}
	defer ClosePoolHandle(poolHandle)

	// Create wallet for issuer
	whIssuer, errCreate := createWallet(issuerConfig(), issuerCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(whIssuer, issuerConfig(), issuerCredentials())

	// Get did for issuer
	didIssuer, verKeyIssuer, errDid := CreateAndStoreDID(whIssuer, "")
	if errDid != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDid)
		return
	}

	// Create wallet for trustee
	whTrustee, errCreate2 := createWallet(holderConfig(), holderCredentials())
	if errCreate2 != nil && errCreate2.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate2)
		return
	}
	defer walletCleanup(whTrustee, holderConfig(), holderCredentials())

	// Get did for issuer
	didTrustee, _, errDid2 := CreateAndStoreDID(whTrustee, seedTrustee1)
	if errDid2 != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDid2)
		return
	}

	// Nym request for issuer
	nymRequest, errNym := BuildNymRequest(didTrustee, didIssuer, verKeyIssuer, "", "TRUSTEE")
	if errNym != nil {
		t.Errorf("BuildNymRequest() error = '%v'", errNym)
		return
	}

	// Sign and submit request
	_, errSign := SignAndSubmitRequest(poolHandle, whTrustee, didTrustee, nymRequest)
	if errSign != nil {
		t.Errorf("SignAndSubmitRequest() error = '%v'", errSign)
		return
	}

	// Create schema for issuer
	schemaId, schemaJson, errSchema := IssuerCreateSchema(didIssuer, "schema-trustee", "1.2", schemaAttributes)
	if errSchema != nil {
		t.Errorf("IssuerCreateSchema() error = '%v'", errSchema)
		return
	}

	// Build schema request
	schemaRequest, errBuildSchema := BuildSchemaRequest(didIssuer, schemaJson)
	if errBuildSchema != nil {
		t.Errorf("BuildSchemaRequest() error = '%v'", errBuildSchema)
		return
	}

	// Sign and submit schema request
	_, errSign2 := SignAndSubmitRequest(poolHandle, whIssuer, didIssuer, schemaRequest)
	if errSign2 != nil {
		t.Errorf("SignAndSubmitRequest() error = '%v'", errSign2)
		return
	}

	// Request to retrieve schema JSON from the ledger
	getSchemaRequest, errBuildSchemaReq := BuildGetSchemaRequest(didIssuer, schemaId)
	if errBuildSchemaReq != nil {
		t.Errorf("BuildGetSchemaRequest() error = '%v'", errBuildSchemaReq)
		return
	}

	// Submit the request
	getSchemaResponse, errSubmit := SubmitRequest(poolHandle, getSchemaRequest)
	if errSubmit != nil {
		t.Errorf("SubmitRequest() error = '%v'", errSubmit)
		return
	}

	// Parse the response
	_, parsedSchemaJson, errParse := ParseGetSchemaResponse(getSchemaResponse)
	if errParse != nil {
		t.Errorf("ParseGetSchemaResponse() error = '%v'", errParse)
		return
	}

	// Create and store credential definition in issuer's wallet
	credDefId, credDefJson, errCredDef := IssuerCreateAndStoreCredentialDefinition(whIssuer, didIssuer, parsedSchemaJson, tag, "CL", `{"support-revocation": false}`)
	if errCredDef != nil {
		t.Errorf("IssuerCreateAndStoreCredentialDefinition() error = '%v'", errCredDef)
		return
	}

	// Build credential definition request
	buildCredRequest, errBuildCred := BuildCredentialDefinitionRequest(didIssuer, credDefJson)
	if errBuildCred != nil {
		t.Errorf("BuildCredentialDefinitionRequest() error = '%v'", errBuildCred)
		return
	}

	// Sign and submit credential definition request
	_, errSign3 := SignAndSubmitRequest(poolHandle, whIssuer, didIssuer, buildCredRequest)
	if errSign3 != nil {
		t.Errorf("SignAndSubmitRequest() error = '%v'", errSign3)
		return
	}

	// Json options
	optionsJson := `{ "noCache": false, "noUpdate": false, "noStore": false, "minFresh": -1 }`

	type args struct {
		CredDefId string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"get-cred-def-cache", args{CredDefId: credDefId}, false},
		{"get-cred-def-invalid-id", args{CredDefId: "invalid-cred"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			credDefRetrieved, errGet := GetCacheCredDef(poolHandle, whIssuer, didIssuer, tt.args.CredDefId, optionsJson)
			hasError := errGet != nil
			if hasError != tt.wantErr {
				if errGet.Error() == indyUtils.GetIndyError(309) {
					t.Log(errGet)
				} else {
					t.Errorf("GetCacheCredDef() error = '%v'", errGet)
					return
				}
			}
			if tt.wantErr {
				t.Log(errGet)
			} else {
				fmt.Println(credDefRetrieved)
			}
		})
	}

	return
}

func TestPurgeSchemaCache(t *testing.T) {
	// Create pool for issuer
	poolHandle, errPool := getPoolLedger("pool")
	if errPool != nil {
		t.Errorf("getPoolLedger() error = '%v'", errPool)
		return
	}
	defer ClosePoolHandle(poolHandle)

	// Create wallet for issuer
	whIssuer, errCreate := createWallet(issuerConfig(), issuerCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(whIssuer, issuerConfig(), issuerCredentials())

	// Get did for issuer
	didIssuer, _, errDid := CreateAndStoreDID(whIssuer, seedSteward1)
	if errDid != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDid)
		return
	}

	// Create schema for issuer
	schemaId, _, errSchema := IssuerCreateSchema(didIssuer, "gvt", "1.0", schemaAttributes)
	if errSchema != nil {
		t.Errorf("IssuerCreateSchema() error = '%v'", errSchema)
		return
	}

	_, errGetCache := GetCacheSchema(poolHandle, whIssuer, didIssuer, schemaId, "{}")
	if errGetCache != nil {
		t.Errorf("GetCacheSchema() error = '%v'", errGetCache)
		return
	}

	options := `{"maxAge": -1}`
	type args struct {
		WalletHandle int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"purge-schema-cache", args{WalletHandle: whIssuer}, false},
		{"purge-schema-cache-invalid-wallet-handle", args{WalletHandle: whIssuer + 100}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errPurge := PurgeSchemaCache(tt.args.WalletHandle, options)
			hasError := errPurge != nil
			if hasError != tt.wantErr {
				t.Errorf("PurgeSchemaCache() error = '%v'", errPurge)
				return
			}
			if tt.wantErr {
				t.Log(errPurge)
			}
		})
	}

	return
}

func TestPurgeCredDefCache(t *testing.T) {
	// Create pool for issuer
	poolHandle, errPool := getPoolLedger("pool")
	if errPool != nil {
		t.Errorf("getPoolLedger() error = '%v'", errPool)
		return
	}
	defer ClosePoolHandle(poolHandle)

	// Create wallet for issuer
	whIssuer, errCreate := createWallet(issuerConfig(), issuerCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(whIssuer, issuerConfig(), issuerCredentials())

	// Get did for issuer
	didIssuer, verKeyIssuer, errDid := CreateAndStoreDID(whIssuer, "")
	if errDid != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDid)
		return
	}

	// Create wallet for trustee
	whTrustee, errCreate2 := createWallet(holderConfig(), holderCredentials())
	if errCreate2 != nil && errCreate2.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate2)
		return
	}
	defer walletCleanup(whTrustee, holderConfig(), holderCredentials())

	// Get did for issuer
	didTrustee, _, errDid2 := CreateAndStoreDID(whTrustee, seedTrustee1)
	if errDid2 != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDid2)
		return
	}

	// Nym request for issuer
	nymRequest, errNym := BuildNymRequest(didTrustee, didIssuer, verKeyIssuer, "", "TRUSTEE")
	if errNym != nil {
		t.Errorf("BuildNymRequest() error = '%v'", errNym)
		return
	}

	// Sign and submit request
	_, errSign := SignAndSubmitRequest(poolHandle, whTrustee, didTrustee, nymRequest)
	if errSign != nil {
		t.Errorf("SignAndSubmitRequest() error = '%v'", errSign)
		return
	}

	// Create schema for issuer
	schemaId, schemaJson, errSchema := IssuerCreateSchema(didIssuer, "schema-trustee", "1.2", schemaAttributes)
	if errSchema != nil {
		t.Errorf("IssuerCreateSchema() error = '%v'", errSchema)
		return
	}

	// Build schema request
	schemaRequest, errBuildSchema := BuildSchemaRequest(didIssuer, schemaJson)
	if errBuildSchema != nil {
		t.Errorf("BuildSchemaRequest() error = '%v'", errBuildSchema)
		return
	}

	// Sign and submit schema request
	_, errSign2 := SignAndSubmitRequest(poolHandle, whIssuer, didIssuer, schemaRequest)
	if errSign2 != nil {
		t.Errorf("SignAndSubmitRequest() error = '%v'", errSign2)
		return
	}

	// Request to retrieve schema JSON from the ledger
	getSchemaRequest, errBuildSchemaReq := BuildGetSchemaRequest(didIssuer, schemaId)
	if errBuildSchemaReq != nil {
		t.Errorf("BuildGetSchemaRequest() error = '%v'", errBuildSchemaReq)
		return
	}

	// Submit the request
	getSchemaResponse, errSubmit := SubmitRequest(poolHandle, getSchemaRequest)
	if errSubmit != nil {
		t.Errorf("SubmitRequest() error = '%v'", errSubmit)
		return
	}

	// Parse the response
	_, parsedSchemaJson, errParse := ParseGetSchemaResponse(getSchemaResponse)
	if errParse != nil {
		t.Errorf("ParseGetSchemaResponse() error = '%v'", errParse)
		return
	}

	// Create and store credential definition in issuer's wallet
	credDefId, credDefJson, errCredDef := IssuerCreateAndStoreCredentialDefinition(whIssuer, didIssuer, parsedSchemaJson, tag, "CL", `{"support-revocation": false}`)
	if errCredDef != nil {
		t.Errorf("IssuerCreateAndStoreCredentialDefinition() error = '%v'", errCredDef)
		return
	}

	// Build credential definition request
	buildCredRequest, errBuildCred := BuildCredentialDefinitionRequest(didIssuer, credDefJson)
	if errBuildCred != nil {
		t.Errorf("BuildCredentialDefinitionRequest() error = '%v'", errBuildCred)
		return
	}

	// Sign and submit credential definition request
	_, errSign3 := SignAndSubmitRequest(poolHandle, whIssuer, didIssuer, buildCredRequest)
	if errSign3 != nil {
		t.Errorf("SignAndSubmitRequest() error = '%v'", errSign3)
		return
	}

	_, errGetCacheCred := GetCacheCredDef(poolHandle, whIssuer, didIssuer, credDefId, `{}`)
	if errGetCacheCred != nil {
		t.Errorf("GetCacheCredDef() error = '%v'", errGetCacheCred)
		return
	}

	options := `{"maxAge": -1}`
	type args struct {
		WalletHandle int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"purge-credential-definition-cache", args{WalletHandle: whIssuer}, false},
		{"purge-credential-definition-cache-invalid-wallet-handle", args{WalletHandle: whIssuer + 100}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errPurge := PurgeCredDefCache(tt.args.WalletHandle, options)
			hasError := errPurge != nil
			if hasError != tt.wantErr {
				t.Errorf("PurgeCredDefCache() error = '%v'", errPurge)
				return
			}
			if tt.wantErr {
				t.Log(errPurge)
			}
		})
	}

	return
}
