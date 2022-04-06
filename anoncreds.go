/*
// ******************************************************************
// Purpose: exported public functions that handles anoncreds functions
// from libindy
// Author:  alexandru.leonte@siemens.com
// Notes:
// Copyright (c): Siemens SRL
// This work is licensed under the terms of the Apache License Version 2.0.  See
// the LICENSE.txt file in the top-level directory.
// ******************************************************************
*/

package indySDK

import "indySDK/anoncreds"

// CreateRevocationState Create revocation state for a credential in the particular time moment
func CreateRevocationState(blobReaderHandle int, revRegDefJson string, revRegDeltaJson string, timestamp uint64, credRevId string) (revStateJson string, err error) {
	channel := anoncreds.CreateRevocationState(blobReaderHandle, revRegDefJson, revRegDeltaJson, timestamp, credRevId)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// IssuerRevokeCredential   Revoke a credential identified by a cred_revoc_id (returned by issuer_create_credential).
func IssuerRevokeCredential(issuerHandle int, blobReaderHandle int, revRegId string, credRevId string) (revRegDeltaJson string, err error) {
	channel := anoncreds.IssuerRevokeCredential(issuerHandle, blobReaderHandle, revRegId, credRevId)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

func IssuerCreateAndStoreRevocReg(wh int,
	issuerDid string,
	revocDefType string,
	tag string,
	credDefId string,
	configJson string,
	blobHandle int) (revocRegId string, revocRegDefJson string, revocRegEntryJson string, err error) {
	channel := anoncreds.CreateAndStoreRevocReg(wh, issuerDid, revocDefType, tag, credDefId, configJson, blobHandle)
	result := <-channel
	if result.Error != nil {
		return "", "", "", result.Error
	}
	return result.Results[0].(string), result.Results[1].(string), result.Results[2].(string), result.Error
}

func IssuerCreateSchema(submitterDid string, name string, version string, attrs string) (schemaId string, schemaJson string, err error) {
	channel := anoncreds.IssuerCreateSchema(submitterDid, name, version, attrs)
	result := <-channel
	if result.Error != nil {
		return "", "", result.Error
	}
	return result.Results[0].(string), result.Results[1].(string), result.Error
}

func IssuerCreateAndStoreCredentialDefinition(wh int, did string, schema string, tag string, signatureType string, configJs string) (credDefId string, credDefJson string, err error) {
	channel := anoncreds.IssuerCreateAndStoreCredentialDef(wh, did, schema, tag, signatureType, configJs)
	result := <-channel
	if result.Error != nil {
		return "", "", result.Error
	}
	return result.Results[0].(string), result.Results[1].(string), result.Error
}

func IssuerRotateCredentialDefStart(walletHandle int, credDefID string, configJson string) (string, error) {
	channel := anoncreds.IssuerRotateCredentialDefStart(walletHandle, credDefID, configJson)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

func IssuerRotateCredentialDefApply(walletHandle int, credDefID string) error {
	channel := anoncreds.IssuerRotateCredentialDefApply(walletHandle, credDefID)
	result := <-channel
	return result.Error
}

// IssuerCreateCredentialOffer Create credential offer
func IssuerCreateCredentialOffer(wh int, credDefId string) (credOffer string, err error) {
	channel := anoncreds.IssuerCreateCredentialOffer(wh, credDefId)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// ProverCreateMasterSecret creates a master secret with a given name and stores it in the wallet.
func ProverCreateMasterSecret(wh int, masterSecretName string) (idMasterSecret string, err error) {
	channel := anoncreds.ProverCreateMasterSecret(wh, masterSecretName)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// ProverCreateCredentialRequest Creates a credential request for the given credential offer.
func ProverCreateCredentialRequest(wh int, proverDID string,
	credOfferJSON string,
	credDefinitionJSON string,
	masterSecretID string) (credentialRequest string, credentialRequestMetadata string, err error) {
	channel := anoncreds.ProverCreateCredentialRequest(wh, proverDID, credOfferJSON, credDefinitionJSON, masterSecretID)
	result := <-channel
	if result.Error != nil {
		return "", "", result.Error
	}
	return result.Results[0].(string), result.Results[1].(string), result.Error
}

// IssuerCreateCredential Creates a credential
func IssuerCreateCredential(whIssuer int,
	credOfferJson string,
	credRequestJson string,
	credValueJson string,
	revocRegistryId string,
	blobHandle int) (credentialJson string, credentialRevocationId string, revocationRegistryDeltaJson string, err error) {
	channel := anoncreds.IssuerCreateCredential(whIssuer, credOfferJson, credRequestJson, credValueJson, revocRegistryId, blobHandle)
	result := <-channel
	if result.Error != nil {
		return "", "", "", result.Error
	}
	return result.Results[0].(string), result.Results[1].(string), result.Results[2].(string), result.Error
}

// ProverStoreCredential stores the credential in the wallet
func ProverStoreCredential(whProver int, credentialIdOptional string, credRequestMetadataJson string, credJson string, credDefJson string, revocRegDefJsonOptional string) (credentialId string, err error) {
	channel := anoncreds.ProverStoreCredential(whProver, credentialIdOptional, credRequestMetadataJson, credJson, credDefJson, revocRegDefJsonOptional)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// ProverDeleteCredential deletes identified credential from wallet
func ProverDeleteCredential(walletHandle int, credentialID string) error {
	channel := anoncreds.ProverDeleteCredential(walletHandle, credentialID)
	result := <-channel
	return result.Error
}

func ProverGetCredentials(walletHandle int, filterJson string) (string, error) {
	channel := anoncreds.ProverGetCredentials(walletHandle, filterJson)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// GenerateNonce nonce
func GenerateNonce() (nonce string, err error) {
	channel := anoncreds.GenerateNonce()
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// ProverSearchForCredentialForProofReq search for credential and returns a search handle
func ProverSearchForCredentialForProofReq(wh int, proofRequestJson, extraQueryJson string) (searchHandle int, err error) {
	channel := anoncreds.ProverSearchForCredentialsForProofReq(wh, proofRequestJson, extraQueryJson)
	result := <-channel
	if result.Error != nil {
		return -1, result.Error
	}
	return result.Results[0].(int), result.Error
}

// ProverCloseCredentialsSearchForProofReq close handle
func ProverCloseCredentialsSearchForProofReq(searchHandle int) (err error) {
	channel := anoncreds.ProverCloseCredentialsSearchForProofReq(searchHandle)
	result := <-channel
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// ProverFetchCredentialsForProofReq - gets credential out of a search handle
func ProverFetchCredentialsForProofReq(sh int, itemReferent string, count int) (credentialJson string, err error) {
	channel := anoncreds.ProverFetchCredentialsForProofReq(sh, itemReferent, count)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

func ProverCreateProof(wh int, proofRequestJson string, requestedCredentialsJson string, masterSecretId string, schemasForAttrsJson string, credentialDefsForAttrsJson string, revStatesJson string) (proofJson string, err error) {
	channel := anoncreds.ProverCreateProof(wh, proofRequestJson, requestedCredentialsJson, masterSecretId, schemasForAttrsJson, credentialDefsForAttrsJson, revStatesJson)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

func VerifierVerifyProof(proofRequestJson string, proofJson string, schemasJson string, credDefsJson string, revRegDefsJson string, revRegsJson string) (valid bool, err error) {
	channel := anoncreds.VerifierVerifyProof(proofRequestJson, proofJson, schemasJson, credDefsJson, revRegDefsJson, revRegsJson)
	result := <-channel
	if result.Error != nil {
		return false, result.Error
	}

	b := result.Results[0].(bool)
	if b {
		return true, result.Error
	} else {
		return false, result.Error
	}

}

func ProverGetCredential(wh int, credentialId string) (credentialJson string, err error) {
	channel := anoncreds.ProverGetCredential(wh, credentialId)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}

	return result.Results[0].(string), result.Error
}

func ProverGetCredentialsForProofRequest(walletHandle int, proofReqJson string) (credentialJson string, err error) {
	channel := anoncreds.ProverGetCredentialsForProofReq(walletHandle, proofReqJson)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

func ProverSearchCredentials(walletHandle int, queryJson string) (searchHandle int, totalCount int, err error) {
	channel := anoncreds.ProverSearchCredentials(walletHandle, queryJson)
	result := <-channel
	if result.Error != nil {
		return 0, 0, result.Error
	}
	return result.Results[0].(int), result.Results[1].(int), result.Error
}

func ProverFetchCredentials(searchHandle int, totalCount int) (credentialsJson string, err error) {
	channel := anoncreds.ProverFetchCredentials(searchHandle, totalCount)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

func ToUnqualified(entity string) (res string, err error) {
	channel := anoncreds.ToUnqualified(entity)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}
