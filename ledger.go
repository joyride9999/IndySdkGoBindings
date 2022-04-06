/*
// ******************************************************************
// Purpose: exported public functions that handles ledger functions
// from libindy
// Author:  alexandru.leonte@siemens.com
// Notes:
// Copyright (c): Siemens SRL
// This work is licensed under the terms of the Apache License Version 2.0.  See
// the LICENSE.txt file in the top-level directory.
// ******************************************************************
*/

package indySDK

import (
	"github.com/viney-shih/go-lock"
	"indySDK/ledger"
	"sync"
	"time"
)

type IndyRequest struct {
	sync.Mutex
}

var indyRequest = lock.NewCASMutex()

// BuildRevocRegEntryRequest Builds a REVOC_REG_ENTRY request. Request to add the definition of revocation registry  to an exists credential definition.
func BuildRevocRegEntryRequest(submitterDid string, revocRegDefId string, revDefType string, value string) (string, error) {
	channel := ledger.BuildRevocRegEntryRequest(submitterDid, revocRegDefId, revDefType, value)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildGetRevocRegDeltaRequest Builds a GET_REVOC_REG_DELTA request. Request to get the delta of the accumulated state of the Revocation Registry.
//    The Delta is defined by from and to timestamp fields.
//    If from is not specified, then the whole state till to will be returned.
func BuildGetRevocRegDeltaRequest(submitterDid string, revocRegDefId string, from int64, to int64) (string, error) {
	channel := ledger.BuildGetRevocRegDeltaRequest(submitterDid, revocRegDefId, from, to)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildRevocRegDefRequest Builds a REVOC_REG_DEF request. Request to add the definition of revocation registry
//    to an exists credential definition.
func BuildRevocRegDefRequest(submitterDid string, revocRegDef string) (string, error) {
	channel := ledger.BuildRevocRegDefRequest(submitterDid, revocRegDef)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildGetDdoRequest creates a request to get DDO
func BuildGetDdoRequest(submitterDid string, targetDid string) (string, error) {
	channel := ledger.BuildGetDdoRequest(submitterDid, targetDid)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildNymRequest creates a nym request (to create an identity ont he blockchain) and returns it
func BuildNymRequest(submitterDid string, targetDid string, targetVerkey string, alias string, role string) (string, error) {
	channel := ledger.BuildNymRequest(submitterDid, targetDid, targetVerkey, alias, role)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildAttribRequest Builds an ATTRIB request. Request to add attribute to a NYM record.
func BuildAttribRequest(submitterDid string, targetDid, hash string, raw string, encrypted string) (string, error) {
	channel := ledger.BuildAttribRequest(submitterDid, targetDid, hash, raw, encrypted)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildGetAttribRequest Builds a GET_ATTRIB request. Request to get information about an Attribute for the specified DID.
func BuildGetAttribRequest(submitterDid string, targetDid, hash string, raw string, encrypted string) (string, error) {
	channel := ledger.BuildGetAttribRequest(submitterDid, targetDid, hash, raw, encrypted)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildGetNymRequest Builds a GET_NYM request. Request to get information about a DID (NYM).
func BuildGetNymRequest(submitterDid string, targetDid string) (string, error) {
	channel := ledger.BuildGetNymRequest(submitterDid, targetDid)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildNodeRequest Builds a NODE request. Request to add a new node to the pool, or updates existing in the pool.
func BuildNodeRequest(submitterDid string, targetDid string, data string) (string, error) {
	channel := ledger.BuildNodeRequest(submitterDid, targetDid, data)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildGetValidatorInfoRequest Builds a GET_VALIDATOR_INFO request.
func BuildGetValidatorInfoRequest(submitterDid string) (string, error) {
	channel := ledger.BuildGetValidatorInfoRequest(submitterDid)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildGetTxnRequest Builds a GET_TXN request. Request to get any transaction by its seq_no.
func BuildGetTxnRequest(submitterDid string, ledgerType string, seqNo int) (string, error) {
	channel := ledger.BuildGetTxnRequest(submitterDid, ledgerType, seqNo)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildPoolConfigRequest Builds a POOL_CONFIG request. Request to change Pool's configuration.
func BuildPoolConfigRequest(submitterDid string, writes bool, force bool) (string, error) {
	channel := ledger.BuildPoolConfigRequest(submitterDid, writes, force)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildPoolRestartRequest Builds a POOL_RESTART request.
func BuildPoolRestartRequest(submitterDid string, action string, dateTime string) (string, error) {
	channel := ledger.BuildPoolRestartRequest(submitterDid, action, dateTime)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildPoolUpgradeRequest Builds a POOL_UPGRADE request.
func BuildPoolUpgradeRequest(submitterDid string, name string, version string, action string, sha256 string, timeOut int32, schedule string,
	justification string, reinstall bool, force bool, package_ string) (string, error) {
	channel := ledger.BuildPoolUpgradeRequest(submitterDid, name, version, action, sha256, timeOut, schedule, justification, reinstall, force, package_)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildGetSchemaRequest creates a schema request
func BuildGetSchemaRequest(submitterDid string, schemaId string) (schemaRequest string, err error) {
	channel := ledger.BuildGetSchemaRequest(submitterDid, schemaId)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildSchemaRequest creates a schema request
func BuildSchemaRequest(submitterDid string, schema string) (request string, err error) {
	channel := ledger.BuildSchemaRequest(submitterDid, schema)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildCredentialDefinitionRequest Builds an CRED_DEF request.
func BuildCredentialDefinitionRequest(submitterDid string, credDefinition string) (request string, err error) {
	channel := ledger.BuildCredentialDefinitionRequest(submitterDid, credDefinition)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildGetRevocRegRequest Builds a GET_REVOC_REG request
func BuildGetRevocRegRequest(submitterDid string, revRegDefId string, timeStamp int64) (request string, err error) {
	channel := ledger.BuildGetRevocRegRequest(submitterDid, revRegDefId, timeStamp)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildGetRevRegDefRequest Builds an GET_REVOC_REG_DEF request.
func BuildGetRevRegDefRequest(submitterDid string, revRegDefId string) (request string, err error) {
	channel := ledger.BuildGetRevocRegDefRequest(submitterDid, revRegDefId)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildGetCredentialDefinitionRequest Builds an GET_CRED_DEF request.
func BuildGetCredentialDefinitionRequest(submitterDid string, credDefinition string) (request string, err error) {
	channel := ledger.BuildGetCredDefRequest(submitterDid, credDefinition)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildAuthRuleRequest Builds a AUTH_RULE request.
func BuildAuthRuleRequest(submitterDid string, txnType string, action string, field string, oldValue string, newValue string, constraint string) (string, error) {
	channel := ledger.BuildAuthRuleRequest(submitterDid, txnType, action, field, oldValue, newValue, constraint)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildAuthRulesRequest Builds a AUTH_RULES request.
func BuildAuthRulesRequest(submitterDid string, data string) (string, error) {
	channel := ledger.BuildAuthRulesRequest(submitterDid, data)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildGetAuthRuleRequest Builds a GET_AUTH_RULE request. Request to get authentication rules for a ledger transaction.
func BuildGetAuthRuleRequest(submitterDid string, txnType string, action string, field string, oldValue string, newValue string) (string, error) {
	channel := ledger.BuildGetAuthRuleRequest(submitterDid, txnType, action, field, oldValue, newValue)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildTxnAuthorAgreementRequest Builds a TXN_AUTHR_AGRMT request. Request to add a new version of Transaction Author Agreement to the ledger.
func BuildTxnAuthorAgreementRequest(submitterDid string, text string, version string, ratificationTs int64, retirementTs int64) (string, error) {
	channel := ledger.BuildTxnAuthorAgreementRequest(submitterDid, text, version, ratificationTs, retirementTs)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildDisableAllTxnAuthorAgreementsRequest Builds a DISABLE_ALL_TXN_AUTHR_AGRMTS request. Request to disable all Transaction Author Agreement on the ledger.
func BuildDisableAllTxnAuthorAgreementsRequest(submitterDid string) (string, error) {
	channel := ledger.BuildDisableAllTxnAuthorAgreementsRequest(submitterDid)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildGetTxnAuthorAgreementRequest  Builds a GET_TXN_AUTHR_AGRMT request. Request to get a specific Transaction Author Agreement from the ledger.
func BuildGetTxnAuthorAgreementRequest(submitterDid string, data string) (string, error) {
	channel := ledger.BuildGetTxnAuthorAgreementRequest(submitterDid, data)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildAcceptanceMechanismsRequest Builds a SET_TXN_AUTHR_AGRMT_AML request. Request to add a new list of acceptance mechanisms for transaction author agreement.
func BuildAcceptanceMechanismsRequest(submitterDid string, aml string, version string, amlContext string) (string, error) {
	channel := ledger.BuildAcceptanceMechanismsRequest(submitterDid, aml, version, amlContext)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildGetAcceptanceMechanismsRequest Builds a GET_TXN_AUTHR_AGRMT_AML request.
func BuildGetAcceptanceMechanismsRequest(submitterDid string, timestamp int64, version string) (string, error) {
	channel := ledger.BuildGetAcceptanceMechanismsRequest(submitterDid, timestamp, version)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// ParseGetRevocRegResponse Parse a GET_REVOC_REG response to get Revocation Registry in the format compatible with Anoncreds API.
func ParseGetRevocRegResponse(getRevRegResp string) (revRegId string, revRegistryDeltaJson string, timestamp uint64, err error) {
	channel := ledger.ParseGetRevocRegResponse(getRevRegResp)
	result := <-channel
	if result.Error != nil {
		return "", "", 0, result.Error
	}
	return result.Results[0].(string), result.Results[1].(string), result.Results[2].(uint64), result.Error
}

// ParseGetRevocRegDeltaResponse Parse a GET_REVOC_REG_DELTA response to get Revocation Registry Delta in the format compatible with Anoncreds API.
func ParseGetRevocRegDeltaResponse(getRevRegDeltaResp string) (revRegId string, revRegistryDeltaJson string, timestamp uint64, err error) {
	channel := ledger.ParseGetRevocRegDeltaResponse(getRevRegDeltaResp)
	result := <-channel
	if result.Error != nil {
		return "", "", 0, result.Error
	}
	return result.Results[0].(string), result.Results[1].(string), result.Results[2].(uint64), result.Error
}

// ParseGetRevocRegDefResponse - parse a rev reg def response
func ParseGetRevocRegDefResponse(getRevocRegDefResponse string) (revRegId string, revRegistryDefJson string, err error) {
	channel := ledger.ParseGetRevocRegDefResponse(getRevocRegDefResponse)
	result := <-channel
	if result.Error != nil {
		return "", "", result.Error
	}
	return result.Results[0].(string), result.Results[1].(string), result.Error
}

// ParseGetSchemaResponse - parse a schema response
func ParseGetSchemaResponse(schemaResponse string) (schemaId string, schemaJson string, err error) {
	channel := ledger.ParseGetSchemaResponse(schemaResponse)
	result := <-channel
	if result.Error != nil {
		return "", "", result.Error
	}
	return result.Results[0].(string), result.Results[1].(string), result.Error
}

// ParseGetNymResponse Parse a GET_NYM response to get NYM data.
func ParseGetNymResponse(nymResponse string) (nymData string, err error) {
	channel := ledger.ParseGetNymResponse(nymResponse)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// ParseGetCredDefResponse - parse a GET_CRED_DEF response
func ParseGetCredDefResponse(getCredDefResp string) (credDefId string, credDefJson string, err error) {
	channel := ledger.ParseGetCredDefResponse(getCredDefResp)
	result := <-channel
	if result.Error != nil {
		return "", "", result.Error
	}
	return result.Results[0].(string), result.Results[1].(string), result.Error
}

// SignAndSubmitRequest sends a request to the blockchain and returns the result
func SignAndSubmitRequest(ph int, wh int, did string, request string) (response string, err error) {
	indyRequest.TryLockWithTimeout(60 * time.Second)
	defer indyRequest.Unlock()
	channel := ledger.SignAndSubmitRequest(ph, wh, did, request)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}

	return result.Results[0].(string), result.Error
}

// SubmitRequest sends a request to the blockchain and returns the result
func SubmitRequest(ph int, request string) (response string, err error) {
	indyRequest.TryLockWithTimeout(60 * time.Second)
	defer indyRequest.Unlock()
	channel := ledger.SubmitRequest(ph, request)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// SignRequest signs request message
func SignRequest(wh int, did string, request string) (response string, err error) {
	// indyRequest.TryLockWithTimeout(60 * time.Second)
	channel := ledger.SignRequest(wh, did, request)
	result := <-channel
	if result.Error != nil {
		return "", err
	}
	return result.Results[0].(string), result.Error
}

// AppendRequestEndorser append an endorser to the request
func AppendRequestEndorser(request, endorserDID string) (response string, err error) {
	channel := ledger.AppendRequestEndorser(request, endorserDID)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// MultiSignRequest signs a request
func MultiSignRequest(wh int, did string, request string) (response string, err error) {
	channel := ledger.MultiSignRequest(wh, did, request)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// GetResponseMetadata Parse transaction response to fetch metadata.
func GetResponseMetadata(response string) (metadataResponse string, err error) {
	channel := ledger.GetResponseMetadata(response)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// AppendTxnAuthorAgreementAcceptanceToRequest Append transaction author agreement acceptance data to a request.
func AppendTxnAuthorAgreementAcceptanceToRequest(requestJson string, text string, version string, taaDigest string, mechanism string, time int64) (string, error) {
	channel := ledger.AppendTxnAuthorAgreementAcceptanceToRequest(requestJson, text, version, taaDigest, mechanism, time)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}