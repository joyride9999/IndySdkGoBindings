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

/*
#include <stdlib.h>
*/
import "C"
import (
	"github.com/joyride9999/IndySdkGoBindings/ledger"
	"github.com/viney-shih/go-lock"
	"sync"
	"time"
	"unsafe"
)

type IndyRequest struct {
	sync.Mutex
}

var indyRequest = lock.NewCASMutex()

// BuildRevocRegEntryRequest Builds a REVOC_REG_ENTRY request. Request to add the definition of revocation registry  to an exists credential definition.
func BuildRevocRegEntryRequest(submitterDid string, revocRegDefId string, revDefType string, value string) (string, error) {

	upSubmitterDid := unsafe.Pointer(C.CString(submitterDid))
	defer C.free(upSubmitterDid)
	upRevocRegDefId := unsafe.Pointer(C.CString(revocRegDefId))
	defer C.free(upRevocRegDefId)
	upRevDefType := unsafe.Pointer(C.CString(revDefType))
	defer C.free(upRevDefType)
	upValue := unsafe.Pointer(C.CString(value))
	defer C.free(upValue)

	channel := ledger.BuildRevocRegEntryRequest(upSubmitterDid, upRevocRegDefId, upRevDefType, upValue)
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

	upSubmitterDid := unsafe.Pointer(GetOptionalValue(submitterDid))
	defer C.free(upSubmitterDid)
	upRevocRegDefId := unsafe.Pointer(C.CString(revocRegDefId))
	defer C.free(upRevocRegDefId)

	channel := ledger.BuildGetRevocRegDeltaRequest(upSubmitterDid, upRevocRegDefId, from, to)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildRevocRegDefRequest Builds a REVOC_REG_DEF request. Request to add the definition of revocation registry
//    to an exists credential definition.
func BuildRevocRegDefRequest(submitterDid string, revocRegDef string) (string, error) {

	upSubmitterDid := unsafe.Pointer(C.CString(submitterDid))
	defer C.free(upSubmitterDid)
	upRevocRegDef := unsafe.Pointer(C.CString(revocRegDef))
	defer C.free(upRevocRegDef)

	channel := ledger.BuildRevocRegDefRequest(upSubmitterDid, upRevocRegDef)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildGetDdoRequest creates a request to get DDO
func BuildGetDdoRequest(submitterDid string, targetDid string) (string, error) {

	upSubmitterDid := unsafe.Pointer(C.CString(submitterDid))
	defer C.free(upSubmitterDid)
	upTargetDid := unsafe.Pointer(C.CString(targetDid))
	defer C.free(upTargetDid)

	channel := ledger.BuildGetDdoRequest(upSubmitterDid, upTargetDid)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildNymRequest creates a nym request (to create an identity ont he blockchain) and returns it
func BuildNymRequest(submitterDid string, targetDid string, targetVerkey string, alias string, role string) (string, error) {

	upSubmitterDid := unsafe.Pointer(C.CString(submitterDid))
	defer C.free(upSubmitterDid)
	upTargetDid := unsafe.Pointer(C.CString(targetDid))
	defer C.free(upTargetDid)
	upTargetVerkey := unsafe.Pointer(C.CString(targetVerkey))
	defer C.free(upTargetVerkey)
	upAlias := unsafe.Pointer(GetOptionalValue(alias))
	defer C.free(upAlias)
	upRole := unsafe.Pointer(C.CString(role))
	defer C.free(upRole)

	channel := ledger.BuildNymRequest(upSubmitterDid, upTargetDid, upTargetVerkey, upAlias, upRole)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildAttribRequest Builds an ATTRIB request. Request to add attribute to a NYM record.
func BuildAttribRequest(submitterDid string, targetDid, hash string, raw string, encrypted string) (string, error) {

	upSubmitterDid := unsafe.Pointer(C.CString(submitterDid))
	defer C.free(upSubmitterDid)
	upTargetDid := unsafe.Pointer(C.CString(targetDid))
	defer C.free(upTargetDid)
	upHash := unsafe.Pointer(GetOptionalValue(hash))
	defer C.free(upHash)
	upRawData := unsafe.Pointer(GetOptionalValue(raw))
	defer C.free(upRawData)
	upEncData := unsafe.Pointer(GetOptionalValue(encrypted))
	defer C.free(upEncData)

	channel := ledger.BuildAttribRequest(upSubmitterDid, upTargetDid, upHash, upRawData, upEncData)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildGetAttribRequest Builds a GET_ATTRIB request. Request to get information about an Attribute for the specified DID.
func BuildGetAttribRequest(submitterDid string, targetDid, raw string, hash string, encrypted string) (string, error) {

	upSubmitterDid := unsafe.Pointer(GetOptionalValue(submitterDid))
	defer C.free(upSubmitterDid)
	upTargetDid := unsafe.Pointer(C.CString(targetDid))
	defer C.free(upTargetDid)
	upHash := unsafe.Pointer(GetOptionalValue(hash))
	defer C.free(upHash)
	upRawData := unsafe.Pointer(GetOptionalValue(raw))
	defer C.free(upRawData)
	upEncData := unsafe.Pointer(GetOptionalValue(encrypted))
	defer C.free(upEncData)

	channel := ledger.BuildGetAttribRequest(upSubmitterDid, upTargetDid, upRawData, upHash, upEncData)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildGetNymRequest Builds a GET_NYM request. Request to get information about a DID (NYM).
func BuildGetNymRequest(submitterDid string, targetDid string) (string, error) {

	upSubmitterDid := unsafe.Pointer(GetOptionalValue(submitterDid))
	defer C.free(upSubmitterDid)

	upTargetDid := unsafe.Pointer(C.CString(targetDid))
	defer C.free(upTargetDid)

	channel := ledger.BuildGetNymRequest(upSubmitterDid, upTargetDid)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildNodeRequest Builds a NODE request. Request to add a new node to the pool, or updates existing in the pool.
func BuildNodeRequest(submitterDid string, targetDid string, data string) (string, error) {

	upSubmitterDid := unsafe.Pointer(C.CString(submitterDid))
	defer C.free(upSubmitterDid)
	upTargetDid := unsafe.Pointer(C.CString(targetDid))
	defer C.free(upTargetDid)
	upData := unsafe.Pointer(C.CString(data))
	defer C.free(upData)

	channel := ledger.BuildNodeRequest(upSubmitterDid, upTargetDid, upData)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildGetValidatorInfoRequest Builds a GET_VALIDATOR_INFO request.
func BuildGetValidatorInfoRequest(submitterDid string) (string, error) {
	upSubmitterDid := unsafe.Pointer(C.CString(submitterDid))
	defer C.free(upSubmitterDid)
	channel := ledger.BuildGetValidatorInfoRequest(upSubmitterDid)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildGetTxnRequest Builds a GET_TXN request. Request to get any transaction by its seq_no.
func BuildGetTxnRequest(submitterDid string, ledgerType string, seqNo int) (string, error) {
	upSubmitterDid := unsafe.Pointer(GetOptionalValue(submitterDid))
	defer C.free(upSubmitterDid)
	upLedgerType := unsafe.Pointer(GetOptionalValue(ledgerType))
	defer C.free(upLedgerType)

	channel := ledger.BuildGetTxnRequest(upSubmitterDid, upLedgerType, seqNo)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildPoolConfigRequest Builds a POOL_CONFIG request. Request to change Pool's configuration.
func BuildPoolConfigRequest(submitterDid string, writes bool, force bool) (string, error) {
	upSubmitterDid := unsafe.Pointer(C.CString(submitterDid))
	defer C.free(upSubmitterDid)

	channel := ledger.BuildPoolConfigRequest(upSubmitterDid, writes, force)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildPoolRestartRequest Builds a POOL_RESTART request.
func BuildPoolRestartRequest(submitterDid string, action string, dateTime string) (string, error) {
	upSubmitterDid := unsafe.Pointer(C.CString(submitterDid))
	defer C.free(upSubmitterDid)
	upAction := unsafe.Pointer(C.CString(action))
	defer C.free(upAction)
	upDateTime := unsafe.Pointer(C.CString(dateTime))
	defer C.free(upDateTime)

	channel := ledger.BuildPoolRestartRequest(upAction, upAction, upDateTime)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildPoolUpgradeRequest Builds a POOL_UPGRADE request.
func BuildPoolUpgradeRequest(submitterDid string, name string, version string, action string, sha256 string, timeOut int32, schedule string,
	justification string, reinstall bool, force bool, indyPackage string) (string, error) {

	upSubmitterDid := unsafe.Pointer(C.CString(submitterDid))
	defer C.free(upSubmitterDid)
	upName := unsafe.Pointer(C.CString(name))
	defer C.free(upName)
	upVersion := unsafe.Pointer(C.CString(version))
	defer C.free(upVersion)
	upAction := unsafe.Pointer(C.CString(action))
	defer C.free(upAction)
	upSha256 := unsafe.Pointer(C.CString(sha256))
	defer C.free(upSha256)
	upSchedule := unsafe.Pointer(GetOptionalValue(schedule))
	defer C.free(upSchedule)
	upReason := unsafe.Pointer(GetOptionalValue(justification))
	defer C.free(upReason)
	upPackage := unsafe.Pointer(GetOptionalValue(indyPackage))
	defer C.free(upPackage)

	channel := ledger.BuildPoolUpgradeRequest(upSubmitterDid, upName, upName, upAction, upSha256, timeOut,
		upSchedule, upReason, reinstall, force, upPackage)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildGetSchemaRequest creates a schema request
func BuildGetSchemaRequest(submitterDid string, schemaId string) (schemaRequest string, err error) {

	upSubmitterDid := unsafe.Pointer(GetOptionalValue(submitterDid))
	defer C.free(upSubmitterDid)
	upSchemaId := unsafe.Pointer(C.CString(schemaId))
	defer C.free(upSchemaId)

	channel := ledger.BuildGetSchemaRequest(upSubmitterDid, upSchemaId)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildSchemaRequest creates a schema request
func BuildSchemaRequest(submitterDid string, schema string) (request string, err error) {
	upSubmitterDid := unsafe.Pointer(C.CString(submitterDid))
	defer C.free(upSubmitterDid)
	upSchemaId := unsafe.Pointer(C.CString(schema))
	defer C.free(upSchemaId)

	channel := ledger.BuildSchemaRequest(upSubmitterDid, upSchemaId)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildCredentialDefinitionRequest Builds an CRED_DEF request.
func BuildCredentialDefinitionRequest(submitterDid string, credDefinition string) (request string, err error) {
	upSubmitterDid := unsafe.Pointer(C.CString(submitterDid))
	defer C.free(upSubmitterDid)
	upCredDef := unsafe.Pointer(C.CString(credDefinition))
	defer C.free(upCredDef)

	channel := ledger.BuildCredentialDefinitionRequest(upSubmitterDid, upCredDef)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildGetRevocRegRequest Builds a GET_REVOC_REG request
func BuildGetRevocRegRequest(submitterDid string, revRegDefId string, timeStamp int64) (request string, err error) {

	upSubmitterDid := unsafe.Pointer(GetOptionalValue(submitterDid))
	defer C.free(upSubmitterDid)
	upRevRegDefId := unsafe.Pointer(C.CString(revRegDefId))
	defer C.free(upRevRegDefId)

	channel := ledger.BuildGetRevocRegRequest(upSubmitterDid, upRevRegDefId, timeStamp)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildGetRevRegDefRequest Builds an GET_REVOC_REG_DEF request.
func BuildGetRevRegDefRequest(submitterDid string, revRegDefId string) (request string, err error) {

	upSubmitterDid := unsafe.Pointer(GetOptionalValue(submitterDid))
	defer C.free(upSubmitterDid)
	upRevRegDefId := unsafe.Pointer(C.CString(revRegDefId))
	defer C.free(upRevRegDefId)
	channel := ledger.BuildGetRevocRegDefRequest(upSubmitterDid, upRevRegDefId)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildGetCredentialDefinitionRequest Builds an GET_CRED_DEF request.
func BuildGetCredentialDefinitionRequest(submitterDid string, credDefinition string) (request string, err error) {

	upSubmitterDid := unsafe.Pointer(GetOptionalValue(submitterDid))
	defer C.free(upSubmitterDid)
	upCredDefinition := unsafe.Pointer(C.CString(credDefinition))
	defer C.free(upCredDefinition)

	channel := ledger.BuildGetCredDefRequest(upSubmitterDid, upCredDefinition)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildAuthRuleRequest Builds a AUTH_RULE request.
func BuildAuthRuleRequest(submitterDid string, txnType string, action string, field string, oldValue string, newValue string, constraint string) (string, error) {

	upSubmitterDid := unsafe.Pointer(C.CString(submitterDid))
	defer C.free(upSubmitterDid)
	upTxnType := unsafe.Pointer(C.CString(txnType))
	defer C.free(upTxnType)
	upAction := unsafe.Pointer(C.CString(action))
	defer C.free(upAction)
	upField := unsafe.Pointer(C.CString(field))
	defer C.free(upField)
	upOldValue := unsafe.Pointer(GetOptionalValue(oldValue))
	defer C.free(upOldValue)
	upNewValue := unsafe.Pointer(GetOptionalValue(newValue))
	defer C.free(upNewValue)
	upConstraint := unsafe.Pointer(C.CString(constraint))
	defer C.free(upConstraint)

	channel := ledger.BuildAuthRuleRequest(upSubmitterDid, upTxnType, upAction, upField, upOldValue, upNewValue, upConstraint)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildAuthRulesRequest Builds a AUTH_RULES request.
func BuildAuthRulesRequest(submitterDid string, data string) (string, error) {
	upSubmitterDid := unsafe.Pointer(C.CString(submitterDid))
	defer C.free(upSubmitterDid)
	upData := unsafe.Pointer(C.CString(data))
	defer C.free(upSubmitterDid)
	channel := ledger.BuildAuthRulesRequest(upSubmitterDid, upData)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildGetAuthRuleRequest Builds a GET_AUTH_RULE request. Request to get authentication rules for a ledger transaction.
func BuildGetAuthRuleRequest(submitterDid string, txnType string, action string, field string, oldValue string, newValue string) (string, error) {

	upSubmitterDid := unsafe.Pointer(GetOptionalValue(submitterDid))
	defer C.free(upSubmitterDid)
	upTxnType := unsafe.Pointer(GetOptionalValue(txnType))
	defer C.free(upTxnType)
	upAction := unsafe.Pointer(GetOptionalValue(action))
	defer C.free(upAction)
	upField := unsafe.Pointer(GetOptionalValue(field))
	defer C.free(upField)
	upOldValue := unsafe.Pointer(GetOptionalValue(oldValue))
	defer C.free(upOldValue)
	upNewValue := unsafe.Pointer(GetOptionalValue(newValue))
	defer C.free(upNewValue)

	channel := ledger.BuildGetAuthRuleRequest(upSubmitterDid, upTxnType, upAction, upField, upOldValue, upNewValue)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildTxnAuthorAgreementRequest Builds a TXN_AUTHR_AGRMT request. Request to add a new version of Transaction Author Agreement to the ledger.
func BuildTxnAuthorAgreementRequest(submitterDid string, text string, version string, ratificationTs int64, retirementTs int64) (string, error) {

	upSubmitterDid := unsafe.Pointer(C.CString(submitterDid))
	defer C.free(upSubmitterDid)
	upText := unsafe.Pointer(GetOptionalValue(text))
	defer C.free(upText)
	upVersion := unsafe.Pointer(C.CString(version))
	defer C.free(upVersion)

	channel := ledger.BuildTxnAuthorAgreementRequest(upSubmitterDid, upText, upVersion, ratificationTs, retirementTs)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildDisableAllTxnAuthorAgreementsRequest Builds a DISABLE_ALL_TXN_AUTHR_AGRMTS request. Request to disable all Transaction Author Agreement on the ledger.
func BuildDisableAllTxnAuthorAgreementsRequest(submitterDid string) (string, error) {
	upSubmitterDid := unsafe.Pointer(C.CString(submitterDid))
	defer C.free(upSubmitterDid)

	channel := ledger.BuildDisableAllTxnAuthorAgreementsRequest(upSubmitterDid)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildGetTxnAuthorAgreementRequest  Builds a GET_TXN_AUTHR_AGRMT request. Request to get a specific Transaction Author Agreement from the ledger.
func BuildGetTxnAuthorAgreementRequest(submitterDid string, data string) (string, error) {
	upSubmitterDid := unsafe.Pointer(GetOptionalValue(submitterDid))
	defer C.free(upSubmitterDid)
	upData := unsafe.Pointer(GetOptionalValue(data))
	defer C.free(upData)

	channel := ledger.BuildGetTxnAuthorAgreementRequest(upSubmitterDid, upData)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildAcceptanceMechanismsRequest Builds a SET_TXN_AUTHR_AGRMT_AML request. Request to add a new list of acceptance mechanisms for transaction author agreement.
func BuildAcceptanceMechanismsRequest(submitterDid string, aml string, version string, amlContext string) (string, error) {
	upSubmitterDid := unsafe.Pointer(C.CString(submitterDid))
	defer C.free(upSubmitterDid)
	upAml := unsafe.Pointer(C.CString(aml))
	defer C.free(upAml)
	upVersion := unsafe.Pointer(C.CString(version))
	defer C.free(upVersion)
	upAmlContext := unsafe.Pointer(GetOptionalValue(amlContext))
	defer C.free(upAmlContext)

	channel := ledger.BuildAcceptanceMechanismsRequest(upSubmitterDid, upAml, upVersion, upAmlContext)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// BuildGetAcceptanceMechanismsRequest Builds a GET_TXN_AUTHR_AGRMT_AML request.
func BuildGetAcceptanceMechanismsRequest(submitterDid string, timestamp int64, version string) (string, error) {
	upSubmitterDid := unsafe.Pointer(GetOptionalValue(submitterDid))
	defer C.free(upSubmitterDid)
	upVersion := unsafe.Pointer(GetOptionalValue(version))
	defer C.free(upVersion)

	channel := ledger.BuildGetAcceptanceMechanismsRequest(upSubmitterDid, timestamp, upVersion)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// ParseGetRevocRegResponse Parse a GET_REVOC_REG response to get Revocation Registry in the format compatible with Anoncreds API.
func ParseGetRevocRegResponse(getRevRegResp string) (revRegId string, revRegistryDeltaJson string, timestamp uint64, err error) {

	upGetRevRegResp := unsafe.Pointer(C.CString(getRevRegResp))
	defer C.free(upGetRevRegResp)

	channel := ledger.ParseGetRevocRegResponse(upGetRevRegResp)
	result := <-channel
	if result.Error != nil {
		return "", "", 0, result.Error
	}
	return result.Results[0].(string), result.Results[1].(string), result.Results[2].(uint64), result.Error
}

// ParseGetRevocRegDeltaResponse Parse a GET_REVOC_REG_DELTA response to get Revocation Registry Delta in the format compatible with Anoncreds API.
func ParseGetRevocRegDeltaResponse(getRevRegDeltaResp string) (revRegId string, revRegistryDeltaJson string, timestamp uint64, err error) {
	upGetRevRegDeltaResp := unsafe.Pointer(C.CString(getRevRegDeltaResp))
	defer C.free(upGetRevRegDeltaResp)

	channel := ledger.ParseGetRevocRegDeltaResponse(upGetRevRegDeltaResp)
	result := <-channel
	if result.Error != nil {
		return "", "", 0, result.Error
	}
	return result.Results[0].(string), result.Results[1].(string), result.Results[2].(uint64), result.Error
}

// ParseGetRevocRegDefResponse - parse a rev reg def response
func ParseGetRevocRegDefResponse(getRevocRegDefResponse string) (revRegId string, revRegistryDefJson string, err error) {
	upGetRevRegDefResp := unsafe.Pointer(C.CString(getRevocRegDefResponse))
	defer C.free(upGetRevRegDefResp)

	channel := ledger.ParseGetRevocRegDefResponse(upGetRevRegDefResp)
	result := <-channel
	if result.Error != nil {
		return "", "", result.Error
	}
	return result.Results[0].(string), result.Results[1].(string), result.Error
}

// ParseGetSchemaResponse - parse a schema response
func ParseGetSchemaResponse(schemaResponse string) (schemaId string, schemaJson string, err error) {
	upSchemaResp := unsafe.Pointer(C.CString(schemaResponse))
	defer C.free(upSchemaResp)

	channel := ledger.ParseGetSchemaResponse(upSchemaResp)
	result := <-channel
	if result.Error != nil {
		return "", "", result.Error
	}
	return result.Results[0].(string), result.Results[1].(string), result.Error
}

// ParseGetNymResponse Parse a GET_NYM response to get NYM data.
func ParseGetNymResponse(nymResponse string) (nymData string, err error) {
	upNymResp := unsafe.Pointer(C.CString(nymResponse))
	defer C.free(upNymResp)

	channel := ledger.ParseGetNymResponse(upNymResp)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// ParseGetCredDefResponse - parse a GET_CRED_DEF response
func ParseGetCredDefResponse(getCredDefResp string) (credDefId string, credDefJson string, err error) {
	upGetCredDefResp := unsafe.Pointer(C.CString(getCredDefResp))
	defer C.free(upGetCredDefResp)

	channel := ledger.ParseGetCredDefResponse(upGetCredDefResp)
	result := <-channel
	if result.Error != nil {
		return "", "", result.Error
	}
	return result.Results[0].(string), result.Results[1].(string), result.Error
}

// SignAndSubmitRequest sends a request to the blockchain and returns the result
func SignAndSubmitRequest(ph int, wh int, did string, request string) (response string, err error) {

	upDid := unsafe.Pointer(C.CString(did))
	defer C.free(upDid)
	upRequest := unsafe.Pointer(C.CString(request))
	defer C.free(upRequest)

	indyRequest.TryLockWithTimeout(60 * time.Second)
	defer indyRequest.Unlock()

	channel := ledger.SignAndSubmitRequest(ph, wh, upDid, upRequest)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}

	return result.Results[0].(string), result.Error
}

// SubmitRequest sends a request to the blockchain and returns the result
func SubmitRequest(ph int, request string) (response string, err error) {
	upRequest := unsafe.Pointer(C.CString(request))
	defer C.free(upRequest)

	indyRequest.TryLockWithTimeout(60 * time.Second)
	defer indyRequest.Unlock()

	channel := ledger.SubmitRequest(ph, upRequest)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// SignRequest signs request message
func SignRequest(wh int, did string, request string) (response string, err error) {
	upDid := unsafe.Pointer(C.CString(did))
	defer C.free(upDid)
	upRequest := unsafe.Pointer(C.CString(request))
	defer C.free(upRequest)

	channel := ledger.SignRequest(wh, upDid, upRequest)
	result := <-channel
	if result.Error != nil {
		return "", err
	}
	return result.Results[0].(string), result.Error
}

// AppendRequestEndorser append an endorser to the request
func AppendRequestEndorser(request, endorserDID string) (response string, err error) {
	upRequest := unsafe.Pointer(C.CString(request))
	defer C.free(upRequest)
	upEndorserDid := unsafe.Pointer(C.CString(endorserDID))
	defer C.free(upEndorserDid)

	channel := ledger.AppendRequestEndorser(upRequest, upEndorserDid)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// MultiSignRequest signs a request
func MultiSignRequest(wh int, did string, request string) (response string, err error) {
	upRequest := unsafe.Pointer(C.CString(request))
	defer C.free(upRequest)
	upDid := unsafe.Pointer(C.CString(did))
	defer C.free(upDid)

	channel := ledger.MultiSignRequest(wh, upDid, upRequest)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// GetResponseMetadata Parse transaction response to fetch metadata.
func GetResponseMetadata(response string) (metadataResponse string, err error) {
	upResponse := unsafe.Pointer(C.CString(response))
	defer C.free(upResponse)

	channel := ledger.GetResponseMetadata(upResponse)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// AppendTxnAuthorAgreementAcceptanceToRequest Append transaction author agreement acceptance data to a request.
func AppendTxnAuthorAgreementAcceptanceToRequest(requestJson string, text string, version string, taaDigest string, mechanism string, time int64) (string, error) {

	upRequest := unsafe.Pointer(C.CString(requestJson))
	defer C.free(upRequest)
	upText := unsafe.Pointer(GetOptionalValue(text))
	defer C.free(upText)
	upVersion := unsafe.Pointer(GetOptionalValue(version))
	defer C.free(upVersion)
	upTaaDigest := unsafe.Pointer(GetOptionalValue(taaDigest))
	defer C.free(upTaaDigest)
	upMech := unsafe.Pointer(GetOptionalValue(mechanism))
	defer C.free(upMech)

	channel := ledger.AppendTxnAuthorAgreementAcceptanceToRequest(upRequest, upText, upVersion, upTaaDigest, upMech, time)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}
