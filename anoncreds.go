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

/*
#include <stdlib.h>
*/
import "C"
import (
	"github.com/joyride9999/IndySdkGoBindings/anoncreds"
	"unsafe"
)

// CreateRevocationState Create revocation state for a credential in the particular time moment
func CreateRevocationState(blobReaderHandle int, revRegDefJson string, revRegDeltaJson string, timestamp uint64, credRevId string) (revStateJson string, err error) {

	upRevRegDefJson := unsafe.Pointer(C.CString(revRegDefJson))
	upRevRegDeltaJson := unsafe.Pointer(C.CString(revRegDeltaJson))
	upCredRevId := unsafe.Pointer(C.CString(credRevId))
	defer C.free(upRevRegDefJson)
	defer C.free(upRevRegDeltaJson)
	defer C.free(upCredRevId)

	channel := anoncreds.CreateRevocationState(blobReaderHandle, upRevRegDefJson, upRevRegDeltaJson, timestamp, upCredRevId)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// IssuerRevokeCredential   Revoke a credential identified by a cred_revoc_id (returned by issuer_create_credential).
func IssuerRevokeCredential(issuerHandle int, blobReaderHandle int, revRegId string, credRevId string) (revRegDeltaJson string, err error) {

	upRevRegId := unsafe.Pointer(C.CString(revRegId))
	upCredRevId := unsafe.Pointer(C.CString(credRevId))
	defer C.free(upRevRegId)
	defer C.free(upCredRevId)

	channel := anoncreds.IssuerRevokeCredential(issuerHandle, blobReaderHandle, upRevRegId, upCredRevId)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

func IssuerCreateAndStoreRevocReg(wh int, issuerDid string, revocDefType string, tag string, credDefId string,
	configJson string, blobHandle int) (revocRegId string, revocRegDefJson string, revocRegEntryJson string, err error) {

	upIssuerDid := unsafe.Pointer(C.CString(issuerDid))
	defer C.free(upIssuerDid)
	upRevocDefType := unsafe.Pointer(C.CString(revocDefType))
	defer C.free(upRevocDefType)
	upTag := unsafe.Pointer(C.CString(tag))
	defer C.free(upTag)
	upCredDefId := unsafe.Pointer(C.CString(credDefId))
	defer C.free(upCredDefId)
	upConfigJson := unsafe.Pointer(C.CString(configJson))
	defer C.free(upConfigJson)

	channel := anoncreds.CreateAndStoreRevocReg(wh, upIssuerDid, upRevocDefType, upTag, upCredDefId, upConfigJson, blobHandle)
	result := <-channel
	if result.Error != nil {
		return "", "", "", result.Error
	}
	return result.Results[0].(string), result.Results[1].(string), result.Results[2].(string), result.Error
}

func IssuerCreateSchema(submitterDid string, name string, version string, attrs string) (schemaId string, schemaJson string, err error) {

	upSubmitterDid := unsafe.Pointer(C.CString(submitterDid))
	defer C.free(upSubmitterDid)
	upName := unsafe.Pointer(C.CString(name))
	defer C.free(upName)
	upVersion := unsafe.Pointer(C.CString(version))
	defer C.free(upVersion)
	upAttrs := unsafe.Pointer(C.CString(attrs))
	defer C.free(upAttrs)

	channel := anoncreds.IssuerCreateSchema(upSubmitterDid, upName, upVersion, upAttrs)
	result := <-channel
	if result.Error != nil {
		return "", "", result.Error
	}
	return result.Results[0].(string), result.Results[1].(string), result.Error
}

func IssuerCreateAndStoreCredentialDefinition(wh int, did string, schema string, tag string, signatureType string, configJs string) (credDefId string, credDefJson string, err error) {

	upDid := unsafe.Pointer(C.CString(did))
	defer C.free(upDid)
	upSchema := unsafe.Pointer(C.CString(schema))
	defer C.free(upSchema)
	upTag := unsafe.Pointer(C.CString(tag))
	defer C.free(upTag)
	upSignatureType := unsafe.Pointer(C.CString(signatureType))
	defer C.free(upSignatureType)
	upConfigJson := unsafe.Pointer(C.CString(configJs))
	defer C.free(upConfigJson)

	channel := anoncreds.IssuerCreateAndStoreCredentialDef(wh, upDid, upSchema, upTag, upSignatureType, upConfigJson)
	result := <-channel
	if result.Error != nil {
		return "", "", result.Error
	}
	return result.Results[0].(string), result.Results[1].(string), result.Error
}

func IssuerRotateCredentialDefStart(walletHandle int, credDefID string, configJson string) (string, error) {

	upCredDefId := unsafe.Pointer(C.CString(credDefID))
	defer C.free(upCredDefId)
	upConfigJson := unsafe.Pointer(C.CString(configJson))
	defer C.free(upConfigJson)

	channel := anoncreds.IssuerRotateCredentialDefStart(walletHandle, upCredDefId, upConfigJson)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

func IssuerRotateCredentialDefApply(walletHandle int, credDefID string) error {
	upCredDefId := unsafe.Pointer(C.CString(credDefID))
	defer C.free(upCredDefId)

	channel := anoncreds.IssuerRotateCredentialDefApply(walletHandle, upCredDefId)
	result := <-channel
	return result.Error
}

// IssuerCreateCredentialOffer Create credential offer
func IssuerCreateCredentialOffer(wh int, credDefId string) (credOffer string, err error) {
	upCredDefId := unsafe.Pointer(C.CString(credDefId))
	defer C.free(upCredDefId)

	channel := anoncreds.IssuerCreateCredentialOffer(wh, upCredDefId)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// ProverCreateMasterSecret creates a master secret with a given name and stores it in the wallet.
func ProverCreateMasterSecret(wh int, masterSecretName string) (idMasterSecret string, err error) {

	upSecretName := unsafe.Pointer(GetOptionalValue(masterSecretName))
	defer C.free(upSecretName)

	channel := anoncreds.ProverCreateMasterSecret(wh, upSecretName)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// ProverCreateCredentialRequest Creates a credential request for the given credential offer.
func ProverCreateCredentialRequest(wh int, proverDID string, credOfferJSON string, credDefinitionJSON string,
	masterSecretID string) (credentialRequest string, credentialRequestMetadata string, err error) {

	upProverDid := unsafe.Pointer(C.CString(proverDID))
	defer C.free(upProverDid)
	upCredOfferJson := unsafe.Pointer(C.CString(credOfferJSON))
	defer C.free(upCredOfferJson)
	upCredDefJson := unsafe.Pointer(C.CString(credDefinitionJSON))
	defer C.free(upCredDefJson)
	upMasterSecretId := unsafe.Pointer(C.CString(masterSecretID))
	defer C.free(upMasterSecretId)

	channel := anoncreds.ProverCreateCredentialRequest(wh, upProverDid, upCredOfferJson, upCredDefJson, upMasterSecretId)
	result := <-channel
	if result.Error != nil {
		return "", "", result.Error
	}
	return result.Results[0].(string), result.Results[1].(string), result.Error
}

// IssuerCreateCredential Creates a credential
func IssuerCreateCredential(whIssuer int, credOfferJson, credRequestJson, credValueJson, revocRegistryId string,
	blobHandle int) (credentialJson string, credentialRevocationId string, revocationRegistryDeltaJson string, err error) {

	upCredOfferJson := unsafe.Pointer(C.CString(credOfferJson))
	defer C.free(upCredOfferJson)
	upCredRequestJson := unsafe.Pointer(C.CString(credRequestJson))
	defer C.free(upCredRequestJson)
	upCredValueJson := unsafe.Pointer(C.CString(credValueJson))
	defer C.free(upCredValueJson)
	upRevRegId := unsafe.Pointer(GetOptionalValue(revocRegistryId))
	defer C.free(upRevRegId)

	channel := anoncreds.IssuerCreateCredential(whIssuer, upCredOfferJson, upCredRequestJson, upCredValueJson, upRevRegId, blobHandle)
	result := <-channel
	if result.Error != nil {
		return "", "", "", result.Error
	}
	return result.Results[0].(string), result.Results[1].(string), result.Results[2].(string), result.Error
}

// ProverStoreCredential stores the credential in the wallet
func ProverStoreCredential(whProver int, credentialIdOptional, credRequestMetadataJson, credJson, credDefJson, revocRegDefJsonOptional string) (credentialId string, err error) {

	upCredId := unsafe.Pointer(GetOptionalValue(credentialIdOptional))
	defer C.free(upCredId)
	upCredRequestMetadataJson := unsafe.Pointer(C.CString(credRequestMetadataJson))
	defer C.free(upCredRequestMetadataJson)
	upCredentialJson := unsafe.Pointer(C.CString(credJson))
	defer C.free(upCredentialJson)
	upCredDefJson := unsafe.Pointer(C.CString(credDefJson))
	defer C.free(upCredDefJson)
	upRevRegDef := unsafe.Pointer(GetOptionalValue(revocRegDefJsonOptional))
	defer C.free(upRevRegDef)

	channel := anoncreds.ProverStoreCredential(whProver, upCredId, upCredRequestMetadataJson, upCredentialJson, upCredDefJson, upRevRegDef)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// ProverDeleteCredential deletes identified credential from wallet
func ProverDeleteCredential(walletHandle int, credentialID string) error {

	upCredentialId := unsafe.Pointer(C.CString(credentialID))
	defer C.free(upCredentialId)

	channel := anoncreds.ProverDeleteCredential(walletHandle, upCredentialId)
	result := <-channel
	return result.Error
}

func ProverGetCredentials(walletHandle int, filterJson string) (string, error) {

	upFilter := unsafe.Pointer(C.CString(filterJson))
	defer C.free(upFilter)

	channel := anoncreds.ProverGetCredentials(walletHandle, upFilter)
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

	upProofRequest := unsafe.Pointer(C.CString(proofRequestJson))
	defer C.free(upProofRequest)

	upExtraQuery := unsafe.Pointer(GetOptionalValue(extraQueryJson))
	defer C.free(upExtraQuery)

	channel := anoncreds.ProverSearchForCredentialsForProofReq(wh, upProofRequest, upExtraQuery)
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

	upItemReferent := unsafe.Pointer(C.CString(itemReferent))
	defer C.free(upItemReferent)

	channel := anoncreds.ProverFetchCredentialsForProofReq(sh, upItemReferent, count)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

func ProverCreateProof(wh int, proofRequestJson, requestedCredentialsJson, masterSecretId, schemasForAttrsJson, credentialDefsForAttrsJson, revStatesJson string) (proofJson string, err error) {

	upProofRequest := unsafe.Pointer(C.CString(proofRequestJson))
	defer C.free(upProofRequest)
	upRequestedCredentials := unsafe.Pointer(C.CString(requestedCredentialsJson))
	defer C.free(upRequestedCredentials)
	upMasterSecretId := unsafe.Pointer(C.CString(masterSecretId))
	defer C.free(upMasterSecretId)
	upSchemasForAttrs := unsafe.Pointer(C.CString(schemasForAttrsJson))
	defer C.free(upSchemasForAttrs)
	upCredentialDefsForAttrs := unsafe.Pointer(C.CString(credentialDefsForAttrsJson))
	defer C.free(upCredentialDefsForAttrs)

	// Library needs this even if empty ... so is not optional
	if len(revStatesJson) == 0 {
		revStatesJson = "{}"
	}

	upRevStates := unsafe.Pointer(C.CString(revStatesJson))
	defer C.free(upRevStates)

	channel := anoncreds.ProverCreateProof(wh, upProofRequest, upRequestedCredentials, upMasterSecretId, upSchemasForAttrs, upCredentialDefsForAttrs, upRevStates)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

func VerifierVerifyProof(proofRequestJson, proofJson, schemasJson, credDefsJson, revRegDefsJson, revRegsJson string) (valid bool, err error) {

	upProofRequest := unsafe.Pointer(C.CString(proofRequestJson))
	defer C.free(upProofRequest)
	upProof := unsafe.Pointer(C.CString(proofJson))
	defer C.free(upProof)
	upSchemas := unsafe.Pointer(C.CString(schemasJson))
	defer C.free(upSchemas)
	upCredDefs := unsafe.Pointer(C.CString(credDefsJson))
	defer C.free(upCredDefs)

	// Library needs this even if empty ... so is not optional
	if len(revRegDefsJson) == 0 {
		revRegDefsJson = "{}"
	}
	upRevRegDefs := unsafe.Pointer(C.CString(revRegDefsJson))
	defer C.free(upRevRegDefs)

	if len(revRegsJson) == 0 {
		revRegsJson = "{}"
	}
	upRevRegs := unsafe.Pointer(C.CString(revRegsJson))
	defer C.free(upRevRegs)

	channel := anoncreds.VerifierVerifyProof(upProofRequest, upProof, upSchemas, upCredDefs, upRevRegDefs, upRevRegs)
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

	upCredentialId := unsafe.Pointer(C.CString(credentialId))
	defer C.free(upCredentialId)

	channel := anoncreds.ProverGetCredential(wh, upCredentialId)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}

	return result.Results[0].(string), result.Error
}

func ProverGetCredentialsForProofRequest(walletHandle int, proofReqJson string) (credentialJson string, err error) {

	upProofReq := unsafe.Pointer(C.CString(proofReqJson))
	defer C.free(upProofReq)

	channel := anoncreds.ProverGetCredentialsForProofReq(walletHandle, upProofReq)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

func ProverSearchCredentials(walletHandle int, queryJson string) (searchHandle int, totalCount int, err error) {

	upQuery := unsafe.Pointer(C.CString(queryJson))
	defer C.free(upQuery)

	channel := anoncreds.ProverSearchCredentials(walletHandle, upQuery)
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

	upEntity := unsafe.Pointer(C.CString(entity))
	defer C.free(upEntity)
	channel := anoncreds.ToUnqualified(upEntity)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}
