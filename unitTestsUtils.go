/*
// ******************************************************************
// Purpose: unit test utils
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
	"github.com/Jeffail/gabs/v2"
	"github.com/joyride9999/IndySdkGoBindings/indyUtils"
	"github.com/joyride9999/IndySdkGoBindings/pool"
	"github.com/joyride9999/IndySdkGoBindings/wallet"
)

const tag = "tag0"
const credValuesJson = `{
				"name" : { "raw": "testName", "encoded": "1"},
				"age" : { "raw": "22", "encoded": "2"},
				"location" : { "raw": "testLocation", "encoded": "3"}
}`
const schemaAttributes = `["name", "age", "location"]`
const seedTrustee1 = "000000000000000000000000Trustee1"
const seedSteward1 = "000000000000000000000000Steward1"
const seedMy1 = "00000000000000000000000000000My1"
const didTrustee = "V4SGRU86Z58d6TV7PBUe6f"
const didMy1 = "VsKV7grR1BUE29mG2Fm2kX"
const endPoint = "127.0.0.1:9700"
const metadata = "ed25519"
const poolGenesisTxn = "pool.txn"
const recordId1 = "recordId1"
const recordId2 = "recordId2"
const recordType = "testType"
const recordValue1 = "recordValue"
const recordValue2 = "recordValue2"
const recordTags1 = `{"tagName1":"str1","tagName2":"5","tagName3":"12"}`
const recordTags2 = `{"tagName1":"str2","tagName2":"pre_str3","tagName3":"2"}`
const recordOptions = `{"retrieveType": true, "retrieveValue": true, "retrieveTags": true}`

/*
	anoncreds_test.go
*/

// createAndStoreCredential contains prerequisites for creating and storing a credential.
func createAndStoreCredential(whIssuer int, didIssuer string, whHolder int, didHolder string) (string, string, string, string, string, string, error) {
	// Create a schema
	schemaId, schemaJson, errSchema := IssuerCreateSchema(didIssuer, "gvt", "1.0", schemaAttributes)
	if errSchema != nil {
		return "", "", "", "", "", "", errSchema
	}

	// Create credential definition
	credentialDefId, credentialDefJson, errCredential := IssuerCreateAndStoreCredentialDefinition(whIssuer, didIssuer, schemaJson, tag, "CL", `{"support-revocation": false}`)
	if errCredential != nil {
		return "", "", "", "", "", "", errCredential
	}

	// Create credential offer
	credOffer, errCredOffer := IssuerCreateCredentialOffer(whIssuer, credentialDefId)
	if errCredOffer != nil {
		return "", "", "", "", "", "", errCredOffer
	}

	// Credential request from offer
	masterSecret, errMaster := ProverCreateMasterSecret(whHolder, "")
	if errMaster != nil {
		return "", "", "", "", "", "", errMaster
	}

	credentialRequest, credentialRequestMetadata, errRequest := ProverCreateCredentialRequest(whHolder, didHolder, credOffer, credentialDefJson, masterSecret)
	if errRequest != nil {
		return "", "", "", "", "", "", errRequest
	}

	// Create the credential
	credentialJson, _, _, errCreateCred := IssuerCreateCredential(whIssuer, credOffer, credentialRequest, credValuesJson, "", 0)
	if errCreateCred != nil {
		return "", "", "", "", "", "", errCreateCred
	}

	// Store the credential into holder wallet
	credentialID, errStore := ProverStoreCredential(whHolder, "", credentialRequestMetadata, credentialJson, credentialDefJson, "")
	if errStore != nil {
		return "", "", "", "", "", "", errStore
	}

	return credentialID, schemaId, schemaJson, credentialDefId, credentialDefJson, masterSecret, nil
}

// searchAndFetchCredForProofReq contains the process of searching for the referents from proof request.
func searchAndFetchCredForProofReq(whHolder int, proofRequest string) (string, string, error) {
	// Search for data
	searchHandle, errSearch := ProverSearchForCredentialForProofReq(whHolder, proofRequest, "")
	if errSearch != nil {
		return "", "", errSearch
	}
	// Close search
	defer ProverCloseCredentialsSearchForProofReq(searchHandle)

	// Get the data for attr1_referent
	attrCreds, errAttrFetch := ProverFetchCredentialsForProofReq(searchHandle, `attr1_referent`, 10)
	if errAttrFetch != nil {
		return "", "", errAttrFetch
	}

	// Get the data for predicate1_referent
	predCreds, errPredFetch := ProverFetchCredentialsForProofReq(searchHandle, `predicate1_referent`, 10)
	if errPredFetch != nil {
		return "", "", errPredFetch
	}

	return attrCreds, predCreds, nil
}

/*
	did_test.go
*/

func issuerCredentials() wallet.Credential {
	res := wallet.Credential{Key: "8dvfYSt5d1taSd6yJdpjq4emkwsPDDLYxkNFysFD2cZY"}
	return res
}

func holderCredentials() wallet.Credential {
	res := wallet.Credential{Key: "19sxnzusdn923mytoskdp9219DMNWIA0x9amDdsAs"}
	return res
}

func trusteeCredentials() wallet.Credential {
	res := wallet.Credential{Key: "0AmLK29JzsDmNQ921ysokMAnzqewPIuOS"}
	return res
}

func testCredentials() wallet.Credential {
	res := wallet.Credential{Key: "18wKLMuAYknwyaiAOX823mNA082Amlz3}"}
	return res
}

func issuerConfig() wallet.Config {
	res := wallet.Config{ID: "issuer1", StorageType: "default", StorageConfig: wallet.StorageConfig{Path: ".\\out\\wallets"}}
	return res
}

func holderConfig() wallet.Config {
	res := wallet.Config{ID: "holder1", StorageType: "default", StorageConfig: wallet.StorageConfig{Path: ".\\out\\wallets"}}
	return res
}

func trusteeConfig() wallet.Config {
	res := wallet.Config{ID: "trustee1", StorageType: "default", StorageConfig: wallet.StorageConfig{Path: ".\\out\\wallets"}}
	return res
}

func testConfig() wallet.Config {
	res := wallet.Config{ID: "wallet_test", StorageType: "default", StorageConfig: wallet.StorageConfig{Path: ".\\out\\wallets"}}
	return res
}

func createWallet(config wallet.Config, credentials wallet.Credential) (int, error) {
	errCreate := CreateWallet(config, credentials)
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		return 0, errCreate
	}

	handler, errOpen := OpenWallet(config, credentials)
	if errOpen != nil {
		return 0, errOpen
	}

	return handler, nil
}

func walletCleanup(handler int, config wallet.Config, credentials wallet.Credential) error {
	errClose := CloseWallet(handler)
	if errClose != nil {
		return errClose
	}

	errDelete := DeleteWallet(config, credentials)
	if errDelete != nil {
		return errDelete
	}
	return nil
}

func identityTrustee1(handle int, seed string) (string, string, error) {
	trusteeDid, trusteeVerKey, err := CreateAndStoreDID(handle, seed)
	if err != nil {
		fmt.Sprintf("CreateAndStoreDid() error = '%v'", err)
		return "", "", err
	}
	return trusteeDid, trusteeVerKey, nil
}

/*
	pool_test.go
*/

func getPoolLedger(poolName string) (int, error) {
	var poolLedger pool.Pool
	poolLedger.Name = poolName
	poolLedger.GenesisTxn = poolGenesisTxn

	errSP := SetPoolProtocolVersion(2)
	if errSP != nil {
		return 0, errSP
	}

	errCreatePool := CreatePoolLedgerConfig(poolLedger)
	if errCreatePool != nil && errCreatePool.Error() != indyUtils.GetIndyError(306) {
		return 0, errCreatePool
	}

	poolHandle, errOp := OpenPoolLedgerConfig(poolLedger)
	if errOp != nil {
		return 0, errOp
	}
	return poolHandle, nil
}

// prepareGetNymReq Builds and sends NYM request. Returns GET_NYM request to be sent to ledger.
func prepareGetNymReq(poolHandle int, walletHandle int, submitterDid string, targetDid string, targetVerKey string, role string) (string, error) {
	nymRequest, errNym := BuildNymRequest(submitterDid, targetDid, targetVerKey, "", role); if errNym != nil {
		return "", errNym
	}
	_, errSign := SignAndSubmitRequest(poolHandle, walletHandle, submitterDid, nymRequest); if errSign != nil {
		return "", errSign
	}

	getNymRequest, errGetNym := BuildGetNymRequest(targetDid, targetDid); if errGetNym != nil {
		return "", errGetNym
	}
	return getNymRequest, nil
}

// prepareGetAttribReq Builds and sends ATTRIB request. Returns GET_ATTRIB request to be sent to ledger.
func prepareGetAttribReq(poolHandle int, walletHandle int, submitterDid string, targetDid string, hash string, raw string, encrypted string) (string, error) {
	attribRequest, errAttrib := BuildAttribRequest(targetDid, targetDid, hash, raw, encrypted); if errAttrib != nil {
		return "", errAttrib
	}
	_, errSign := SignAndSubmitRequest(poolHandle, walletHandle, targetDid, attribRequest); if errSign != nil {
		return "", errSign
	}
	getAttribRequest, errGetAttrib := BuildGetAttribRequest(submitterDid, targetDid, raw, "", ""); if errGetAttrib != nil {
		return "", errGetAttrib
	}
	return getAttribRequest, nil
}

// prepareGetSchemaReq Builds and sends SCHEMA request. Returns GET_SCHEMA request to be sent to ledger.
func prepareGetSchemaReq(poolHandle int, walletHandle int, submitterDid string, name string, version string, attrs string) (string, string, string, error) {
	schemaId, schemaJson, errSchema := IssuerCreateSchema(submitterDid, name, version, attrs); if errSchema != nil {
		return "", "", "", errSchema
	}
	schemaRequest, errSchemaReq := BuildSchemaRequest(submitterDid, schemaJson); if errSchemaReq != nil {
		return "", "", "", errSchemaReq
	}
	_, errSign := SignAndSubmitRequest(poolHandle, walletHandle, submitterDid, schemaRequest); if errSign != nil {
		return "", "", "", errSign
	}
	getSchemaReq, errGetSchema := BuildGetSchemaRequest(submitterDid, schemaId); if errGetSchema != nil {
		return "", "", "", errGetSchema
	}
	return schemaId, schemaJson, getSchemaReq, nil
}

func isEqual(expected *gabs.Container, resulted *gabs.Container) bool {
	ok := false

	if len(expected.ChildrenMap()) != len(resulted.ChildrenMap()) {
		return ok
	} else {
		// If expected underlying value is an object, the map of children isn't empty.
		if len(expected.ChildrenMap()) != 0 && len(resulted.ChildrenMap()) != 0 {
			// Iterate through expected items.
			for path, element := range expected.ChildrenMap() {
				// Search elements by their path in resulted.
				search := resulted.Path(path)
				switch v := search.Data().(type) {
				case []interface{}:
					// Element is an array.
					ok = true
					if len(element.Children()) != len(search.Children()) {
						return false
					}
					// Iterate through array's components.
					for _, component := range element.Children() {
						exists := false
						for _, item := range search.Children() {
							exists = isEqual(component, item)
							if exists {
								break
							}
						}
						if !exists {
							ok = false
							break
						}
					}

				case map[string]interface{}:
					// Element is a *gabs.Container.
					ok = isIncluded(search, element)
				case nil:
					// Element is not found in resulted JSON.
					ok = false
				default:
					if search.Data() == element.Data() {
						ok = true
					} else {
						ok = false
					}
					v = v
				}
				if ok == false {
					break
				}
			}
		// If underlying value is not an object.
		} else {
			if expected.Data() == resulted.Data() {
				ok = true
			} else {
				ok = false
			}
		}
	}
	return ok
}

func isIncluded(expected *gabs.Container, resulted *gabs.Container) bool {
	ok := false

	// Iterate through expected items.
	for path, element := range expected.ChildrenMap() {
		// Search elements by their path in resulted.
		search := resulted.Path(path)
		switch search.Data().(type) {
		case []interface{}:
			// Element is an array.
			ok = true
			if len(element.Children()) != len(search.Children()) {
				return false
			}
			// Iterate through array's components.
			for _, component := range element.Children() {
				exists := false
				for _, item := range search.Children() {
					exists = isEqual(component, item)
					if exists {
						break
					}
				}
				if !exists {
					ok = false
					break
				}
			}

		case map[string]interface{}:
			// Exception if *gabs.Container is found, but it's empty.
			if len(search.ChildrenMap()) == 0 {
				ok = true
			} else {
				// Element is a *gabs.Container.
				ok = isIncluded(element, search)
			}
		case nil:
			// Element is not found in resulted JSON.
			ok = false
		default:
			if search.Data() == element.Data() {
				ok = true
			} else {
				ok = false
			}
		}
		if ok == false {
			break
		}
	}
	return ok
}

/*
	nonsecrets_test.go
*/

// checkRecordValue checks a record's value field for an expected value.
func checkRecordValue(walletHandle int, recordType string, recordId string, recordOptions string, expectedValue string) (bool, error) {
	ok := false

	record, errGet := IndyGetWalletRecord(walletHandle, recordType, recordId, recordOptions)
	if errGet != nil {
		return ok, errGet
	}

	recordParsed, errParse := gabs.ParseJSON([]byte(record)); if errParse != nil {
		return ok, errParse
	}
	if expectedValue == recordParsed.Path("value").Data() {
		ok = true
	}

	return ok, nil
}

// checkRecordTags checks a record's tags field for an expected value.
func checkRecordTags(walletHandle int, recordType string, recordId string, recordOptions string, expectedTags string) (bool, error) {
	ok := false

	record, errGet := IndyGetWalletRecord(walletHandle, recordType, recordId, recordOptions)
	if errGet != nil {
		return ok, errGet
	}

	recordParsed, errParse := gabs.ParseJSON([]byte(record)); if errParse != nil {
		return ok, errParse
	}
	if expectedTags == `{}` {
		if recordParsed.Path("tags").String() == expectedTags {
			ok = true
		}
	} else {
		expectedTagsParsed, errParse2 := gabs.ParseJSON([]byte(expectedTags))
		if errParse2 != nil {
			return ok, errParse2
		}
		if isEqual(expectedTagsParsed, recordParsed.Path("tags")) {
			ok = true
		}
	}

	return ok, nil
}