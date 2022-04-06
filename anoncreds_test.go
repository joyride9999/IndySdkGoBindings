/*
// ******************************************************************
// Purpose: anoncreds unit testing
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
	"github.com/Jeffail/gabs/v2"
	"github.com/google/uuid"
	"github.com/joyride9999/IndySdkGoBindings/indyUtils"
	"testing"
)

func TestIssuerCreateAndStoreCredentialDefinition(t *testing.T) {
	walletHandle, errCreate := createWallet(issuerConfig(), issuerCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, issuerConfig(), issuerCredentials())

	did, _, errDid := CreateAndStoreDID(walletHandle, "")
	if errDid != nil {
		t.Errorf("CreateAndStoreDID() error = '%v'", errDid)
		return
	}

	_, schemaJson, errSchema := IssuerCreateSchema(did, "gvt", "1.0", "[\"len\"]")
	if errSchema != nil {
		t.Errorf("IssuerCreateSchema() error = '%v'", errSchema)
		return
	}

	type args struct {
		SubmitterDid string
		SchemaJson   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"create-non-existing-cred-definition", args{SubmitterDid: did, SchemaJson: schemaJson}, false},
		{"create-without-schema", args{SubmitterDid: did, SchemaJson: ""}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			credentialDefID, credentialDefJson, errCredential := IssuerCreateAndStoreCredentialDefinition(walletHandle, tt.args.SubmitterDid, tt.args.SchemaJson, tag, "CL", `{"support-revocation": false}`)
			hasError := (errCredential != nil)
			if hasError != tt.wantErr {
				t.Errorf("IssuerCreateAndStoreCredentialDefinition() error = '%v', wantErr = '%v'", errCredential, tt.wantErr)
				return
			}
			if tt.wantErr {
				fmt.Println("Expected error: ", errCredential)
			}
			fmt.Println(credentialDefID + "\n" + credentialDefJson)
		})
	}

	return
}

func TestIssuerCreateCredential(t *testing.T) {
	// Create and open issuer wallet
	whIssuer, errCreate := createWallet(issuerConfig(), issuerCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(whIssuer, issuerConfig(), issuerCredentials())

	// Get did for issuer
	didIssuer, _, errDidIssuer := CreateAndStoreDID(whIssuer, "")
	if errDidIssuer != nil {
		t.Errorf("CreateAndStoreDID() error = '%v'", errDidIssuer)
		return
	}

	// Create and open holder wallet
	whHolder, errCreate2 := createWallet(holderConfig(), holderCredentials())
	if errCreate2 != nil && errCreate2.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate2)
		return
	}
	defer walletCleanup(whHolder, holderConfig(), holderCredentials())

	// Get did for holder
	didHolder, _, errDidHolder := CreateAndStoreDID(whHolder, "")
	if errDidHolder != nil {
		t.Errorf("CreateAndStoreDID() error = '%v'", errDidHolder)
		return
	}

	_, schemaJson, errSchema := IssuerCreateSchema(didIssuer, "gvt", "1.0", "[\"len\"]")
	if errSchema != nil {
		t.Errorf("IssuerCreateSchema() error = '%v'", errSchema)
		return
	}

	masterSecret, errMaster := ProverCreateMasterSecret(whHolder, "")
	if errMaster != nil {
		t.Errorf("ProverCreateMasterSecret() error = '%v'", errMaster)
		return
	}

	credentialDefID, credentialDefJson, errCredential := IssuerCreateAndStoreCredentialDefinition(whIssuer, didIssuer, schemaJson, tag, "CL", `{"support-revocation": false}`)
	if errCredential != nil {
		t.Errorf("IssuerCreateAndStoreCredentialDefinition() error = '%v'", errCredential)
		return
	}
	credentialOffer, errOffer := IssuerCreateCredentialOffer(whIssuer, credentialDefID)
	if errOffer != nil {
		t.Errorf("IssuerCreateCredentialOffer() error = '%v'", errOffer)
		return
	}
	credentialRequest, _, errRequest := ProverCreateCredentialRequest(whHolder, didHolder, credentialOffer, credentialDefJson, masterSecret)
	if errRequest != nil {
		t.Errorf("ProverCreateCredentialRequest() error = '%v'", errRequest)
		return
	}

	type args struct {
		SubmitterDid  string
		CredValueJson string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"create-non-existing-cred-definition", args{SubmitterDid: didIssuer, CredValueJson: `{"len": {"raw": "42", "encoded": "42"}}`}, false},
		{"create-with-invalid-credential-values", args{SubmitterDid: didIssuer, CredValueJson: `{"age": {"raw": "18", "encoded": "18"`}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			credential, _, _, errCreateCred := IssuerCreateCredential(whIssuer, credentialOffer, credentialRequest, tt.args.CredValueJson, "", 0)
			hasError := errCreateCred != nil
			if hasError != tt.wantErr {
				t.Errorf("IssuerCreateCredential() error = '%v', wantErr ='%v'", errCreateCred, tt.wantErr)
			}
			if tt.wantErr {
				fmt.Println("Expected error: ", errCreateCred)
			}
			fmt.Println(credential)
		})
	}
	return
}

func TestIssuerRotateCredentialDef(t *testing.T) {
	walletHandle, errCreate := createWallet(issuerConfig(), issuerCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, issuerConfig(), issuerCredentials())

	did, _, errDid := CreateAndStoreDID(walletHandle, "")
	if errDid != nil {
		t.Errorf("CreateAndStoreDID() error = '%v'", errDid)
		return
	}

	_, schemaJson, errSchema := IssuerCreateSchema(did, "gvt", "1.0", "[\"len\"]")
	if errSchema != nil {
		t.Errorf("IssuerCreateSchema() error = '%v'", errSchema)
		return
	}

	credentialDefID, credentialDefJson, errCredential := IssuerCreateAndStoreCredentialDefinition(walletHandle, did, schemaJson, tag, "CL", `{"support-revocation": false}`)
	if errCredential != nil {
		t.Errorf("IssuerCreateAndStoreCredentialDefinition() error = '%v'", errCredential)
		return
	}

	type args struct {
		WalletHandle int
		ConfigJson   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"rotate-credential-definition", args{WalletHandle: walletHandle, ConfigJson: `{"support-revocation": false}`}, false},
		{"rotate-without-config", args{WalletHandle: walletHandle, ConfigJson: `{}`}, false},
		{"invalid-wallet-handle", args{WalletHandle: walletHandle, ConfigJson: `{}`}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Condition for "invalid-wallet-handle" test case.
			if tt.wantErr {
				tt.args.WalletHandle = -1
			}
			tempCred, errRotate := IssuerRotateCredentialDefStart(tt.args.WalletHandle, credentialDefID, tt.args.ConfigJson)
			hasError := errRotate != nil
			if hasError != tt.wantErr {
				t.Errorf("IssuerRotateCredentialDefStart() error = '%v'", errRotate)
				return
			}
			if tempCred != "" && tempCred != credentialDefJson {
				fmt.Println("IssuerRotateCredentialDefStart(): successful")
			}
			if tt.wantErr {
				fmt.Println("Error expected: ", errRotate)
			} else {
				errApply := IssuerRotateCredentialDefApply(tt.args.WalletHandle, credentialDefID)
				if errApply != nil {
					t.Errorf("IssuerRotateCredentialDefApply() error = '%v'", errApply)
					return
				}
			}
		})
	}

	return
}

func TestProverCreateProof(t *testing.T) {
	// Create and open issuer wallet
	whIssuer, errCreate := createWallet(issuerConfig(), issuerCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(whIssuer, issuerConfig(), issuerCredentials())

	// Get did for issuer
	didIssuer, _, errDidIssuer := CreateAndStoreDID(whIssuer, "")
	if errDidIssuer != nil {
		t.Errorf("CreateAndStoreDID() error = '%v'", errDidIssuer)
		return
	}

	// Create and open holder wallet
	whHolder, errCreate2 := createWallet(holderConfig(), holderCredentials())
	if errCreate2 != nil && errCreate2.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate2)
		return
	}
	defer walletCleanup(whHolder, holderConfig(), holderCredentials())

	// Get did for holder
	didHolder, _, errDidHolder := CreateAndStoreDID(whHolder, "")
	if errDidHolder != nil {
		t.Errorf("CreateAndStoreDID() error = '%v'", errDidHolder)
		return
	}

	nonce, _ := GenerateNonce()

	_, schemaId, schemaJson, credentialDefId, credentialDefJson, masterSecret, errCredential := createAndStoreCredential(whIssuer, didIssuer, whHolder, didHolder)
	if errCredential != nil {
		t.Errorf("createAndStoreCredential() error = '%v'", errCredential)
		return
	}

	// Proof request and search operation for referents.
	proofRequest := fmt.Sprintf(`{"nonce": "%s", "name": "proofRequest", "ver" : "1.0", "version": "0.1",
	"requested_attributes": {"attr1_referent": {"name": "name"}},
	"requested_predicates": {"predicate1_referent": {"name": "age", "p_type": ">=", "p_value" : 2}} }`, nonce)
	attrCreds, predCreds, errSearch := searchAndFetchCredForProofReq(whHolder, proofRequest)
	if errSearch != nil {
		t.Errorf("searchAndFetchCredForProofReq() error = '%v'", errSearch)
		return
	}

	// Read credential definition id from fetched data
	credInfoAttrJson, errGabs := gabs.ParseJSON([]byte(attrCreds))
	if errGabs != nil {
		t.Errorf("Gabs ParseJSON() error = '%v'", errGabs)
		return
	}
	attrCredId := credInfoAttrJson.Path("0.cred_info.referent").String()

	credInfoPredJson, errGabs := gabs.ParseJSON([]byte(predCreds))
	if errGabs != nil {
		t.Errorf("Gabs ParseJSON() error = '%v'", errGabs)
		return
	}
	predCredId := credInfoPredJson.Path("0.cred_info.referent").String()

	requestedCredJson := fmt.Sprintf(`{"self_attested_attributes": {},
		"requested_attributes": {"attr1_referent": {"cred_id": %s, "revealed": true}},
		"requested_predicates": {"predicate1_referent": {"cred_id": %s}}}`,
		attrCredId, predCredId)
	schemasJson := fmt.Sprintf(`{"%s":%s}`, schemaId, schemaJson)
	credDefsJson := fmt.Sprintf(`{"%s":%s}`, credentialDefId, credentialDefJson)

	type args struct {
		ProofRequest string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"create-proof-works", args{ProofRequest: proofRequest}, false},
		{"create-proof-with-invalid-proof-request", args{ProofRequest: ""}, true}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proofJson, errProof := ProverCreateProof(whHolder, tt.args.ProofRequest, requestedCredJson, masterSecret, schemasJson, credDefsJson, "{}")
			hasError := errProof != nil
			if hasError != tt.wantErr {
				t.Errorf("ProverCreateProof() error = '%v', wantErr = '%v'", errProof, tt.wantErr)
				return
			}
			if tt.wantErr {
				fmt.Println("Expected error: ", errProof)
			}
			proofJson = proofJson
		})
	}

	return
}

func TestProverCreateCredentialRequest(t *testing.T) {
	// Create and open issuer wallet
	whIssuer, errCreate := createWallet(issuerConfig(), issuerCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(whIssuer, issuerConfig(), issuerCredentials())

	// Get did for issuer
	didIssuer, _, errDidIssuer := CreateAndStoreDID(whIssuer, "")
	if errDidIssuer != nil {
		t.Errorf("CreateAndStoreDID() error = '%v'", errDidIssuer)
		return
	}

	// Create and open holder wallet
	whHolder, errCreate := createWallet(holderConfig(), holderCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(whHolder, holderConfig(), holderCredentials())

	// Get did for holder
	didHolder, _, errDidHolder := CreateAndStoreDID(whHolder, "")
	if errDidHolder != nil {
		t.Errorf("CreateAndStoreDID() error = '%v'", errDidHolder)
		return
	}

	// Create a schema
	_, schemaJson, errSchema := IssuerCreateSchema(didIssuer, "gvt", "1.0", `["name", "age", "location"]`)
	if errSchema != nil {
		t.Errorf("IssuerCreateSchema() error = '%v'", errSchema)
		return
	}

	// Create credential definition
	credentialDefId, credentialDefJson, errCredential := IssuerCreateAndStoreCredentialDefinition(whIssuer, didIssuer, schemaJson, tag, "CL", `{"support-revocation": false}`)
	if errCredential != nil {
		t.Errorf("IssuerCreateAndStoreCredentialDefinition() error = '%v'", errCredential)
		return
	}

	// Create credential offer
	credOffer, errCredOffer := IssuerCreateCredentialOffer(whIssuer, credentialDefId)
	if errCredOffer != nil {
		t.Errorf("IssuerCreateCredentialOffer() error = '%v'", errCredOffer)
	}

	// Credential request from offer
	masterSecret, errMaster := ProverCreateMasterSecret(whHolder, "")
	if errMaster != nil {
		t.Errorf("ProverCreateMasterSecret() error = '%v'", errMaster)
		return
	}

	type args struct {
		CredOffer string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"create-credential-request-works", args{CredOffer: credOffer}, false},
		{"create-credential-request-without-cred-offer", args{CredOffer: ""}, true}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			credentialRequest, credentialRequestMetadata, errRequest := ProverCreateCredentialRequest(whHolder, didHolder, tt.args.CredOffer, credentialDefJson, masterSecret)
			hasError := errRequest != nil
			if hasError != tt.wantErr {
				t.Errorf("ProverCreateCredentialRequest() error = '%v', wantErr = '%v'", errRequest, tt.wantErr)
				return
			}
			if tt.wantErr {
				fmt.Println("Expected error: ", errRequest)
			}
			fmt.Println(credentialRequest, credentialRequestMetadata)
		})
	}
	return
}

func TestProverCreateMasterSecret(t *testing.T) {
	walletHandle, errCreate := createWallet(issuerConfig(), issuerCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, issuerConfig(), issuerCredentials())

	uid, _ := uuid.NewUUID()
	secretName := uid.String()
	type args struct {
		MasterSecretName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"indy-generates-secret-name", args{""}, false},
		{"named-secret", args{secretName}, false},
		{"existing-named-secret", args{secretName}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			masterSecretId, errMaster := ProverCreateMasterSecret(walletHandle, tt.args.MasterSecretName)
			hasError := errMaster != nil
			if hasError != tt.wantErr {
				t.Errorf("ProverCreateMasterSecret() error = '%v', wantErr = '%v'", errMaster, tt.wantErr)
				return
			}
			if tt.wantErr {
				fmt.Println("Expected error ", errMaster)
			} else {
				fmt.Println("New master secret id ", masterSecretId)
			}

		})
	}

	return

}

func TestProverDeleteCredential(t *testing.T) {
	// Create and open issuer wallet
	whIssuer, errCreate := createWallet(issuerConfig(), issuerCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(whIssuer, issuerConfig(), issuerCredentials())

	// Get did for issuer
	didIssuer, _, errDidIssuer := CreateAndStoreDID(whIssuer, "")
	if errDidIssuer != nil {
		t.Errorf("CreateAndStoreDID() error = '%v'", errDidIssuer)
		return
	}

	// Create and open holder wallet
	whHolder, errCreate2 := createWallet(holderConfig(), holderCredentials())
	if errCreate2 != nil && errCreate2.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate2)
		return
	}
	defer walletCleanup(whHolder, holderConfig(), holderCredentials())

	// Get did for holder
	didHolder, _, errDidHolder := CreateAndStoreDID(whHolder, "")
	if errDidHolder != nil {
		t.Errorf("CreateAndStoreDID() error = '%v'", errDidHolder)
		return
	}

	credentialId, _, _, _, _, _, errCredential := createAndStoreCredential(whIssuer, didIssuer, whHolder, didHolder)
	if errCredential != nil {
		t.Errorf("createAndStoreCredential() error = '%v'", errCredential)
		return
	}

	type args struct {
		WalletHandler int
		CredentialId  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"delete-credential-works", args{WalletHandler: whHolder, CredentialId: credentialId}, false},
		{"delete-invalid-credential", args{WalletHandler: whHolder, CredentialId: "test-invalid-credential"}, true}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errDelete := ProverDeleteCredential(tt.args.WalletHandler, tt.args.CredentialId)
			hasError := errDelete != nil
			if hasError != tt.wantErr {
				t.Errorf("ProverDeleteCredential() error = '%v', wantErr = '%v'", errDelete, tt.wantErr)
				return
			}
			if tt.wantErr {
				fmt.Println("Expected error: ", errDelete)
			} else {
				// Check after credential has been successfully deleted
				_, errGet := ProverGetCredential(tt.args.WalletHandler, tt.args.CredentialId)
				if errGet != nil && errGet.Error() != indyUtils.GetIndyError(212) {
					t.Errorf("ProverGetCredential() error = '%v'", errGet)
					return
				} else if errGet.Error() == indyUtils.GetIndyError(212) {
					fmt.Println("Credential deleted")
				}
			}
		})
	}

	return
}

func TestProverGetCredential(t *testing.T) {
	// Create and open issuer wallet
	whIssuer, errCreate := createWallet(issuerConfig(), issuerCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(whIssuer, issuerConfig(), issuerCredentials())

	// Get did for issuer
	didIssuer, _, errDidIssuer := CreateAndStoreDID(whIssuer, "")
	if errDidIssuer != nil {
		t.Errorf("CreateAndStoreDID() error = '%v'", errDidIssuer)
		return
	}

	// Create and open holder wallet
	whHolder, errCreate2 := createWallet(holderConfig(), holderCredentials())
	if errCreate2 != nil && errCreate2.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate2)
		return
	}
	defer walletCleanup(whHolder, holderConfig(), holderCredentials())

	// Get did for holder
	didHolder, _, errDidHolder := CreateAndStoreDID(whHolder, "")
	if errDidHolder != nil {
		t.Errorf("CreateAndStoreDID() error = '%v'", errDidHolder)
		return
	}

	credentialId, _, _, _, _, _, errCredential := createAndStoreCredential(whIssuer, didIssuer, whHolder, didHolder)
	if errCredential != nil {
		t.Errorf("createAndStoreCredential() error = '%v'", errCredential)
		return
	}

	type args struct {
		WalletHandler int
		CredentialId  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{

		{"get-credential-works", args{WalletHandler: whHolder, CredentialId: credentialId}, false},
		{"get-invalid-credential", args{WalletHandler: whHolder, CredentialId: "test-invalid-credential"}, true}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			credential, errGet := ProverGetCredential(tt.args.WalletHandler, tt.args.CredentialId)
			hasError := errGet != nil
			if hasError != tt.wantErr {
				t.Errorf("ProverGetCredential() error = '%v', wantErr = '%v'", errGet, tt.wantErr)
				return
			}
			if tt.wantErr {
				fmt.Println("Expected error: ", errGet)
			} else {
				fmt.Println(credential)
			}
		})
	}
	return
}

func TestProverGetCredentials(t *testing.T) {
	// Create and open issuer wallet
	whIssuer, errCreate := createWallet(issuerConfig(), issuerCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(whIssuer, issuerConfig(), issuerCredentials())

	// Get did for issuer
	didIssuer, _, errDidIssuer := CreateAndStoreDID(whIssuer, "")
	if errDidIssuer != nil {
		t.Errorf("CreateAndStoreDID() error = '%v'", errDidIssuer)
		return
	}

	didTrustee, _, errDidTrustee := CreateAndStoreDID(whIssuer, "")
	if errDidTrustee != nil {
		t.Errorf("CreateAndStoreDID() error = '%v'", errDidTrustee)
		return
	}

	// Create and open holder wallet
	whHolder, errCreate2 := createWallet(holderConfig(), holderCredentials())
	if errCreate2 != nil && errCreate2.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate2)
		return
	}
	defer walletCleanup(whHolder, holderConfig(), holderCredentials())

	// Get did for holder
	didHolder, _, errDidHolder := CreateAndStoreDID(whHolder, "")
	if errDidHolder != nil {
		t.Errorf("CreateAndStoreDID() error = '%v'", errDidHolder)
		return
	}

	_, _, _, _, _, _, errCredential := createAndStoreCredential(whIssuer, didIssuer, whHolder, didHolder)
	if errCredential != nil {
		t.Errorf("createAndStoreCredential() error = '%v'", errCredential)
		return
	}

	_, _, _, _, _, _, errCredential2 := createAndStoreCredential(whIssuer, didTrustee, whHolder, didHolder)
	if errCredential2 != nil {
		t.Errorf("createAndStoreCredential() error = '%v'", errCredential2)
		return
	}

	filterIssuerDidJson := fmt.Sprintf(`{"issuer_did": "%s"}`, didTrustee)
	type args struct {
		WalletHandler int
		FilterJson    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{

		{"get-credentials-works-without-filter", args{WalletHandler: whHolder, FilterJson: "{}"}, false},
		{"get-credentials-works-with-issuer-did-filter", args{WalletHandler: whHolder, FilterJson: filterIssuerDidJson}, false},
		{"get-credentials-invalid-filter", args{WalletHandler: whHolder, FilterJson: `{"cred_id": onPurpose}`}, true}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			credentials, errGetCredentials := ProverGetCredentials(tt.args.WalletHandler, tt.args.FilterJson)
			hasError := errGetCredentials != nil
			if hasError != tt.wantErr {
				t.Errorf("ProverGetCredentials() error = '%v', wantErr = '%v'", errGetCredentials, tt.wantErr)
				return
			}
			if tt.wantErr {
				fmt.Println("Expected error: ", errGetCredentials)
			} else {
				fmt.Println(credentials)
			}
		})
	}
	return
}

func TestProverGetCredentialsForProofRequest(t *testing.T) {
	// Create and open issuer wallet
	whIssuer, errCreate := createWallet(issuerConfig(), issuerCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(whIssuer, issuerConfig(), issuerCredentials())

	// Get did for issuer
	didIssuer, _, errDidIssuer := CreateAndStoreDID(whIssuer, "")
	if errDidIssuer != nil {
		t.Errorf("CreateAndStoreDID() error = '%v'", errDidIssuer)
		return
	}

	// Create and open holder wallet
	whHolder, errCreate2 := createWallet(holderConfig(), holderCredentials())
	if errCreate2 != nil && errCreate2.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate2)
		return
	}
	defer walletCleanup(whHolder, holderConfig(), holderCredentials())

	// Get did for holder
	didHolder, _, errDidHolder := CreateAndStoreDID(whHolder, "")
	if errDidHolder != nil {
		t.Errorf("CreateAndStoreDID() error = '%v'", errDidHolder)
		return
	}

	_, _, _, _, _, _, errCredential := createAndStoreCredential(whIssuer, didIssuer, whHolder, didHolder)
	if errCredential != nil {
		t.Errorf("createAndStoreCredential() error = '%v'", errCredential)
		return
	}

	nonce, _ := GenerateNonce()
	// Proof request and search operation for referents.
	proofRequest := fmt.Sprintf(`{"nonce": "%s", "name": "proofRequest", "ver" : "1.0", "version": "0.1",
	"requested_attributes": {"attr1_referent": {"name": "name"}},
	"requested_predicates": {"predicate1_referent": {"name": "age", "p_type": ">=", "p_value" : 2}} }`, nonce)

	type args struct {
		ProofRequest string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{

		{"get-credentials-for-proof-req-works", args{ProofRequest: proofRequest}, false},
		{"get-credentials-for-proof-req-without-proof", args{ProofRequest: ""}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			credentials, errGetCredentials := ProverGetCredentialsForProofRequest(whHolder, tt.args.ProofRequest)
			hasError := errGetCredentials != nil
			if hasError != tt.wantErr {
				t.Errorf("ProverGetCredentialsForProofRequest() error = '%v', wantErr = '%v'", errGetCredentials, tt.wantErr)
				return
			}
			if tt.wantErr {
				fmt.Println("Expected error: ", errGetCredentials)
				return
			} else {
				fmt.Println(credentials)
			}
			return
		})
	}
	return
}

func TestProverSearchCredentials(t *testing.T) {
	// Create and open issuer wallet
	whIssuer, errCreate := createWallet(issuerConfig(), issuerCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(whIssuer, issuerConfig(), issuerCredentials())

	// Get did for issuer
	didIssuer, _, errDidIssuer := CreateAndStoreDID(whIssuer, "")
	if errDidIssuer != nil {
		t.Errorf("CreateAndStoreDID() error = '%v'", errDidIssuer)
		return
	}

	// Create and open holder wallet
	whHolder, errCreate2 := createWallet(holderConfig(), holderCredentials())
	if errCreate2 != nil && errCreate2.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate2)
		return
	}
	// TODO: CommonIOError: IO Error (?) in defer
	defer walletCleanup(whHolder, holderConfig(), holderCredentials())

	// Get did for holder
	didHolder, _, errDidHolder := CreateAndStoreDID(whHolder, "")
	if errDidHolder != nil {
		t.Errorf("CreateAndStoreDID() error = '%v'", errDidHolder)
		return
	}

	_, _, _, _, _, _, errCredential := createAndStoreCredential(whIssuer, didIssuer, whHolder, didHolder)
	if errCredential != nil {
		t.Errorf("createAndStoreCredential() error = '%v'", errCredential)
		return
	}

	type args struct {
		ProofRequest string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{

		{"search-credentials-with-empty-proof", args{ProofRequest: "{}"}, false},
		{"search-credentials-with-invalid-proof", args{ProofRequest: "invalid-test"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			searchHandle, totalCount, errSearch := ProverSearchCredentials(whHolder, tt.args.ProofRequest)
			hasError := errSearch != nil
			if hasError != tt.wantErr {
				t.Errorf("ProverSearchCredentials() error = '%v', wantErr = '%v'", errSearch, tt.wantErr)
				return
			}
			if tt.wantErr {
				fmt.Println("Expected error: ", errSearch)
			} else {
				// Check if variables are valid
				_, errFetch := ProverFetchCredentials(searchHandle, totalCount)
				if errFetch != nil {
					t.Errorf("ProverFetchCredentials() error = '%v'", errFetch)
					return
				}
			}

		})
	}
	return
}

func TestProverSearchForCredentialForProofReq(t *testing.T) {
	// Create and open issuer wallet
	whIssuer, errCreate := createWallet(issuerConfig(), issuerCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(whIssuer, issuerConfig(), issuerCredentials())

	// Get did for issuer
	didIssuer, _, errDidIssuer := CreateAndStoreDID(whIssuer, "")
	if errDidIssuer != nil {
		t.Errorf("CreateAndStoreDID() error = '%v'", errDidIssuer)
		return
	}

	// Create and open holder wallet
	whHolder, errCreate2 := createWallet(holderConfig(), holderCredentials())
	if errCreate2 != nil && errCreate2.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate2)
		return
	}
	defer walletCleanup(whHolder, holderConfig(), holderCredentials())

	// Get did for holder
	didHolder, _, errDidHolder := CreateAndStoreDID(whHolder, "")
	if errDidHolder != nil {
		t.Errorf("CreateAndStoreDID() error = '%v'", errDidHolder)
		return
	}

	_, _, _, _, _, _, errCredential := createAndStoreCredential(whIssuer, didIssuer, whHolder, didHolder)
	if errCredential != nil {
		t.Errorf("createAndStoreCredential() error = '%v'", errCredential)
		return
	}

	nonce, _ := GenerateNonce()
	// Proof request and search operation for referents.
	proofRequest := fmt.Sprintf(`{"nonce": "%s", "name": "proofRequest", "ver" : "1.0", "version": "0.1",
	"requested_attributes": {"attr1_referent": {"name": "name"}},
	"requested_predicates": {"predicate1_referent": {"name": "age", "p_type": ">=", "p_value" : 2}} }`, nonce)

	type args struct {
		ProofRequest string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{

		{"search-credentials-for-proof-req-works", args{ProofRequest: proofRequest}, false},
		{"search-credentials-for-proof-req-with-invalid-proof", args{ProofRequest: "invalid-test"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			searchHandle, errSearch := ProverSearchForCredentialForProofReq(whHolder, tt.args.ProofRequest, "")
			hasError := errSearch != nil
			if hasError != tt.wantErr {
				t.Errorf("ProverSearchForCredentialForProofReq() error = '%v', wantErr = '%v'", errSearch, tt.wantErr)
				return
			}
			if tt.wantErr {
				fmt.Println("Expected error: ", errSearch)
			} else {
				fmt.Println(searchHandle)
			}
		})
	}
	return
}

func TestProverStoreCredential(t *testing.T) {
	// Create and open issuer wallet
	whIssuer, errCreate := createWallet(issuerConfig(), issuerCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(whIssuer, issuerConfig(), issuerCredentials())

	// Get did for issuer
	didIssuer, _, errDidIssuer := CreateAndStoreDID(whIssuer, "")
	if errDidIssuer != nil {
		t.Errorf("CreateAndStoreDID() error = '%v'", errDidIssuer)
		return
	}

	// Create and open holder wallet
	whHolder, errCreate2 := createWallet(holderConfig(), holderCredentials())
	if errCreate2 != nil && errCreate2.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate2)
		return
	}
	defer walletCleanup(whHolder, holderConfig(), holderCredentials())

	// Get did for holder
	didHolder, _, errDidHolder := CreateAndStoreDID(whHolder, "")
	if errDidHolder != nil {
		t.Errorf("CreateAndStoreDID() error = '%v'", errDidHolder)
		return
	}

	_, schemaJson, errSchema := IssuerCreateSchema(didIssuer, "gvt", "1.0", schemaAttributes)
	if errSchema != nil {
		t.Errorf("IssuerCreateSchema() error = '%v'", errSchema)
		return
	}

	masterSecret, errMaster := ProverCreateMasterSecret(whHolder, "")
	if errMaster != nil {
		t.Errorf("ProverCreateMasterSecret() error = '%v'", errMaster)
		return
	}

	credentialDefID, credentialDefJson, errCredential := IssuerCreateAndStoreCredentialDefinition(whIssuer, didIssuer, schemaJson, tag, "CL", `{"support-revocation": false}`)
	if errCredential != nil {
		t.Errorf("IssuerCreateAndStoreCredentialDefinition() error = '%v'", errCredential)
		return
	}
	credentialOffer, errOffer := IssuerCreateCredentialOffer(whIssuer, credentialDefID)
	if errOffer != nil {
		t.Errorf("IssuerCreateCredentialOffer() error = '%v'", errOffer)
		return
	}
	credentialRequest, credentialRequestMetadata, errRequest := ProverCreateCredentialRequest(whHolder, didHolder, credentialOffer, credentialDefJson, masterSecret)
	if errRequest != nil {
		t.Errorf("ProverCreateCredentialRequest() error = '%v'", errRequest)
		return
	}

	credentialJson, _, _, errCreateCred := IssuerCreateCredential(whIssuer, credentialOffer, credentialRequest, credValuesJson, "", 0)
	if errCreateCred != nil {
		t.Errorf("IssuerCreateCredential() error = '%v'", errCreateCred)
		return
	}

	type args struct {
		CredentialJson string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"create-non-existing-cred-definition", args{CredentialJson: credentialJson}, false},
		{"create-with-invalid-credential-values", args{CredentialJson: "invalid-credential-json"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			credentialId, errStore := ProverStoreCredential(whHolder, "", credentialRequestMetadata, tt.args.CredentialJson, credentialDefJson, "")
			hasError := errStore != nil
			if hasError != tt.wantErr {
				t.Errorf("IssuerCreateCredential() error = '%v', wantErr ='%v'", errStore, tt.wantErr)
				return
			}
			if tt.wantErr {
				fmt.Println("Expected error: ", errStore)
			} else {
				fmt.Println(credentialId)
			}
		})
	}
	return

}

func TestToUnqualified(t *testing.T) {
	qualified := "did:sov:NcYxiDXkpYi6ov5FcYDi1e"
	unqualified := "NcYxiDXkpYi6ov5FcYDi1e"
	type args struct {
		Entity string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{

		{"to-unqualified-works", args{Entity: qualified}, false},
		{"to-unqualified-unqualified-entity", args{Entity: unqualified}, false},
		{"to-unqualified-empty-entity", args{Entity: ""}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, errQualify := ToUnqualified(tt.args.Entity)
			hasError := errQualify != nil
			if hasError != tt.wantErr {
				t.Errorf("ToUnqualified() error = '%v', wantErr = '%v'", errQualify, tt.wantErr)
				return
			}
			if tt.wantErr {
				fmt.Println("Error expected: ", errQualify)
			} else {
				fmt.Println(result)
			}
		})
	}
	return
}

func TestVerifierVerifyProof(t *testing.T) {
	// Create and open issuer wallet
	whIssuer, errCreate := createWallet(issuerConfig(), issuerCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(whIssuer, issuerConfig(), issuerCredentials())

	// Get did for issuer
	didIssuer, _, errDidIssuer := CreateAndStoreDID(whIssuer, "")
	if errDidIssuer != nil {
		t.Errorf("CreateAndStoreDID() error = '%v'", errDidIssuer)
		return
	}

	// Create and open holder wallet
	whHolder, errCreate2 := createWallet(holderConfig(), holderCredentials())
	if errCreate2 != nil && errCreate2.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("CreateWallet() error = '%v'", errCreate2)
		return
	}
	defer walletCleanup(whHolder, holderConfig(), holderCredentials())

	// Get did for holder
	didHolder, _, errDidHolder := CreateAndStoreDID(whHolder, "")
	if errDidHolder != nil {
		t.Errorf("CreateAndStoreDID() error = '%v'", errDidHolder)
		return
	}

	nonce, _ := GenerateNonce()

	_, schemaId, schemaJson, credentialDefId, credentialDefJson, masterSecret, errCredential := createAndStoreCredential(whIssuer, didIssuer, whHolder, didHolder)
	if errCredential != nil {
		t.Errorf("createAndStoreCredential() error = '%v'", errCredential)
		return
	}

	// Proof request and search operation for referents.
	proofRequest := fmt.Sprintf(`{"nonce": "%s", "name": "proofRequest", "ver" : "1.0", "version": "0.1",
	"requested_attributes": {"attr1_referent": {"name": "name"}},
	"requested_predicates": {"predicate1_referent": {"name": "age", "p_type": ">=", "p_value" : 2}} }`, nonce)
	attrCreds, predCreds, errSearch := searchAndFetchCredForProofReq(whHolder, proofRequest)
	if errSearch != nil {
		t.Errorf("searchAndFetchCredForProofReq() error = '%v'", errSearch)
		return
	}

	// Read credential definition id from fetched data
	credInfoAttrJson, errGabs := gabs.ParseJSON([]byte(attrCreds))
	if errGabs != nil {
		t.Errorf("Gabs ParseJSON() error = '%v'", errGabs)
		return
	}
	attrCredId := credInfoAttrJson.Path("0.cred_info.referent").String()

	credInfoPredJson, errGabs := gabs.ParseJSON([]byte(predCreds))
	if errGabs != nil {
		t.Errorf("Gabs ParseJSON() error = '%v'", errGabs)
		return
	}
	predCredId := credInfoPredJson.Path("0.cred_info.referent").String()

	requestedCredJson := fmt.Sprintf(`{"self_attested_attributes": {},
		"requested_attributes": {"attr1_referent": {"cred_id": %s, "revealed": true}},
		"requested_predicates": {"predicate1_referent": {"cred_id": %s}}}`,
		attrCredId, predCredId)
	schemasJson := fmt.Sprintf(`{"%s":%s}`, schemaId, schemaJson)
	credDefsJson := fmt.Sprintf(`{"%s":%s}`, credentialDefId, credentialDefJson)

	proofJson, errProof := ProverCreateProof(whHolder, proofRequest, requestedCredJson, masterSecret, schemasJson, credDefsJson, "{}")
	if errProof != nil {
		t.Errorf("ProverCreateProof() error = '%v'", errProof)
		return
	}

	type args struct {
		ProofRequest string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"create-proof-works", args{ProofRequest: proofRequest}, false},
		{"create-proof-with-invalid-proof-request", args{ProofRequest: ""}, true}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, errVerify := VerifierVerifyProof(tt.args.ProofRequest, proofJson, schemasJson, credDefsJson, "", "")
			hasError := errVerify != nil
			if hasError != tt.wantErr {
				t.Errorf("VerifierVerifyProof() error = '%v', wantErr = '%v'", errVerify, tt.wantErr)
				return
			}
			if tt.wantErr {
				fmt.Println("Expected error: ", errVerify)
			}
			fmt.Println(valid)
		})
	}

	return
}
