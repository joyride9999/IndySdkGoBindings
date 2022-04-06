/*
// ******************************************************************
// Purpose: Wrapper to call libindy, imports functions from indy_ledger.h
// Author:  alexandru.leonte@siemens.com
// Notes:
// Copyright (c): Siemens SRL
// This work is licensed under the terms of the Apache License Version 2.0.  See
// the LICENSE.txt file in the top-level directory.
// ******************************************************************
*/

package ledger

/*
#cgo CFLAGS: -I ../include
#cgo LDFLAGS: -L${SRCDIR}/../lib -lindy
#include <indy_core.h>
typedef void (*cb_buildRequest)(indy_handle_t, indy_error_t, char*);
extern void buildRequestCB(indy_handle_t, indy_error_t, char*);

typedef void (*cb_signAndSubmitRequest)(indy_handle_t, indy_error_t, char*);
extern void signAndSubmitRequestCB(indy_handle_t, indy_error_t, char*);

typedef void (*cb_parseGetResponse)(indy_handle_t, indy_error_t, char*, char*);
extern void parseGetResponseCB(indy_handle_t, indy_error_t, char*, char*);

typedef void (*cb_parseGetResponseDelta)(indy_handle_t, indy_error_t, char*, char*, unsigned long long);
extern void parseGetResponseDeltaCB(indy_handle_t, indy_error_t, char*, char*, unsigned long long);

typedef void (*cb_parseGetNymResponse)(indy_handle_t, indy_error_t, char*);
extern void parseGetNymResponseCB(indy_handle_t, indy_error_t, char*);

typedef void (*cb_appendRequestEndorser)(indy_handle_t, indy_error_t, char*);
extern void appendRequestEndorserCB(indy_handle_t, indy_error_t, char*);

typedef void (*cb_multiSignRequest)(indy_handle_t, indy_error_t, char*);
extern void multiSignRequestCB(indy_handle_t, indy_error_t, char*);

*/
import "C"
import (
	"errors"
	"github.com/joyride9999/IndySdkGoBindings/indyUtils"
	"unsafe"
)

//export multiSignRequestCB
func multiSignRequestCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, request *C.char) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle),
			indyUtils.IndyResult{Error: nil,
				Results: []interface{}{
					string(C.GoString(request)),
				}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// MultiSignRequest    Multi signs request message.
func MultiSignRequest(wh int, did string, request string) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
	   Multi signs request message.

	   Adds submitter information to passed request json, signs it with submitter
	   sign key (see wallet_sign).

	   :param wallet_handle: wallet handle (created by open_wallet).
	   :param submitter_did: Id of Identity stored in secured Wallet.
	   :param request_json: Request data json.
	   :return: Signed request json.
	*/

	// Call to indy function
	res := C.indy_multi_sign_request(commandHandle,
		(C.indy_handle_t)(wh),
		C.CString(did),
		C.CString(request),
		(C.cb_multiSignRequest)(unsafe.Pointer(C.multiSignRequestCB)))

	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export appendRequestEndorserCB
func appendRequestEndorserCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, request *C.char) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle),
			indyUtils.IndyResult{Error: nil,
				Results: []interface{}{
					string(C.GoString(request)),
				}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// AppendRequestEndorser   Append Endorser to an existing request.
func AppendRequestEndorser(request string, endorserDid string) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
	  Append Endorser to an existing request.

	    An author of request still is a `DID` used as a `submitter_did` parameter for the building of the request.
	    But it is expecting that the transaction will be sent by the specified Endorser.

	    Note: Both Transaction Author and Endorser must sign output request after that.

	    More about Transaction Endorser: https://github.com/hyperledger/indy-node/blob/master/design/transaction_endorser.md
	                                     https://github.com/hyperledger/indy-sdk/blob/master/docs/configuration.md

	    :param request_json: original request data json.
	    :param endorser_did: DID of the Endorser that will submit the transaction.
	                         The Endorser's DID must be present on the ledger.

	    :return: Updated request result as json.
	*/

	// Call to indy function
	res := C.indy_append_request_endorser(commandHandle,
		C.CString(request),
		C.CString(endorserDid),
		(C.cb_appendRequestEndorser)(unsafe.Pointer(C.appendRequestEndorserCB)))

	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export parseGetResponseDeltaCB
func parseGetResponseDeltaCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, objId *C.char, objJSON *C.char, timestamp C.ulonglong) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle),
			indyUtils.IndyResult{Error: nil,
				Results: []interface{}{
					string(C.GoString(objId)),
					string(C.GoString(objJSON)),
					uint64(timestamp),
				}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// ParseGetRevocRegResponse Parse a GET_REVOC_REG response to get Revocation Registry in the format compatible with Anoncreds API.
func ParseGetRevocRegResponse(getRevRegResp string) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
	   :param get_revoc_reg_response: response of GET_REVOC_REG request.
	    :return: Revocation Registry Definition Id, Revocation Registry json and Timestamp.
	      {
	          "value": Registry-specific data {
	              "accum": string - current accumulator value.
	          },
	          "ver": string - version revocation registry json
	      }
	*/

	// Call to indy function
	res := C.indy_parse_get_revoc_reg_response(commandHandle,
		C.CString(getRevRegResp),
		(C.cb_parseGetResponseDelta)(unsafe.Pointer(C.parseGetResponseDeltaCB)))

	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

// ParseGetRevocRegDeltaResponse Parse a GET_REVOC_REG_DELTA response to get Revocation Registry Delta in the format compatible with Anoncreds API.
func ParseGetRevocRegDeltaResponse(getRevRegDeltaResp string) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
	   :param get_revoc_reg_delta_response: response of GET_REVOC_REG_DELTA request.
	    :return: Revocation Registry Definition Id, Revocation Registry Delta json and Timestamp.
	      {
	          "value": Registry-specific data {
	              prevAccum: string - previous accumulator value.
	              accum: string - current accumulator value.
	              issued: array<number> - an array of issued indices.
	              revoked: array<number> an array of revoked indices.
	          },
	          "ver": string
	      }
	*/

	// Call to indy function
	res := C.indy_parse_get_revoc_reg_delta_response(commandHandle,
		C.CString(getRevRegDeltaResp),
		(C.cb_parseGetResponseDelta)(unsafe.Pointer(C.parseGetResponseDeltaCB)))

	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export parseGetResponseCB
func parseGetResponseCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, objId *C.char, objJSON *C.char) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle),
			indyUtils.IndyResult{Error: nil,
				Results: []interface{}{
					string(C.GoString(objId)),
					string(C.GoString(objJSON)),
				}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// ParseGetRevocRegDefResponse Parse a GET_REVOC_REG_DEF response to get Revocation Registry Definition in the format compatible with Anoncreds API.
func ParseGetRevocRegDefResponse(getRevocRegDefResponse string) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
	  :param get_revoc_ref_def_response: response of GET_REVOC_REG_DEF request.
	    :return: Revocation Registry Definition Id and Revocation Registry Definition json.
	      {
	          "id": string - ID of the Revocation Registry,
	          "revocDefType": string - Revocation Registry type (only CL_ACCUM is supported for now),
	          "tag": string - Unique descriptive ID of the Registry,
	          "credDefId": string - ID of the corresponding CredentialDefinition,
	          "value": Registry-specific data {
	              "issuanceType": string - Type of Issuance(ISSUANCE_BY_DEFAULT or ISSUANCE_ON_DEMAND),
	              "maxCredNum": number - Maximum number of credentials the Registry can serve.
	              "tailsHash": string - Hash of tails.
	              "tailsLocation": string - Location of tails file.
	              "publicKeys": <public_keys> - Registry's public key.
	          },
	          "ver": string - version of revocation registry definition json.
	      }
	*/

	// Call to indy function
	res := C.indy_parse_get_revoc_reg_def_response(commandHandle,
		C.CString(getRevocRegDefResponse),
		(C.cb_parseGetResponse)(unsafe.Pointer(C.parseGetResponseCB)))

	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

// ParseGetCredDefResponse Parse a GET_CRED_DEF response to get Credential Definition in the format compatible with Anoncreds API.
func ParseGetCredDefResponse(getCredDefResp string) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
	   :param get_cred_def_response: response of GET_CRED_DEF request.
	   :return: Credential Definition Id and Credential Definition json.
	     {
	         id: string - identifier of credential definition
	         schemaId: string - identifier of stored in ledger schema
	         type: string - type of the credential definition. CL is the only supported type now.
	         tag: string - allows to distinct between credential definitions for the same issuer and schema
	         value: Dictionary with Credential Definition's data: {
	             primary: primary credential public key,
	             Optional<revocation>: revocation credential public key
	         },
	         ver: Version of the Credential Definition json
	     }
	*/

	// Call to indy function
	res := C.indy_parse_get_cred_def_response(commandHandle,
		C.CString(getCredDefResp),
		(C.cb_parseGetResponse)(unsafe.Pointer(C.parseGetResponseCB)))

	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

// ParseGetSchemaResponse Parse a GET_SCHEMA response to get Schema in the format compatible with Anoncreds API.
func ParseGetSchemaResponse(schemaResponse string) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
			 :param get_schema_response: response of GET_SCHEMA request.
		    :return: Schema Id and Schema json.
		     {
		         id: identifier of schema
		         attrNames: array of attribute name strings
		         name: Schema's name string
		         version: Schema's version string
		         ver: Version of the Schema json
		     }
	*/

	// Call to indy function
	res := C.indy_parse_get_schema_response(commandHandle,
		C.CString(schemaResponse),
		(C.cb_parseGetResponse)(unsafe.Pointer(C.parseGetResponseCB)))

	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export parseGetNymResponseCB
func parseGetNymResponseCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, nymJson *C.char) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle),
			indyUtils.IndyResult{Error: nil,
				Results: []interface{}{
					string(C.GoString(nymJson)),
				}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// ParseGetNymResponse Parse a GET_NYM response to get NYM data.
func ParseGetNymResponse(getNymResponse string) chan indyUtils.IndyResult {

	// Prepare the call parameters.
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)
	/*
	   Parse a GET_NYM response to get NYM data.

	   :param response: response on GET_NYM request.
	   :return: NYM data
	   {
	       did: DID as base58-encoded string for 16 or 32 bit DID value.
	       verkey: verification key as base58-encoded string.
	       role: Role associated number
	                               null (common USER)
	                               0 - TRUSTEE
	                               2 - STEWARD
	                               101 - TRUST_ANCHOR
	                               101 - ENDORSER - equal to TRUST_ANCHOR that will be removed soon
	                               201 - NETWORK_MONITOR
	   }
	 */

	// Call to indy function
	res := C.indy_parse_get_nym_response(commandHandle,
		C.CString(getNymResponse),
		(C.cb_parseGetNymResponse)(unsafe.Pointer(C.parseGetNymResponseCB)))

	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}
//export buildRequestCB
func buildRequestCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, request *C.char) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle),
			indyUtils.IndyResult{Error: nil,
				Results: []interface{}{
					string(C.GoString(request)),
				}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// BuildGetRevocRegDeltaRequest       Builds a GET_REVOC_REG_DELTA request. Request to get the delta of the accumulated state of the Revocation Registry.
//    The Delta is defined by from and to timestamp fields.
//    If from is not specified, then the whole state till to will be returned.
func BuildGetRevocRegDeltaRequest(submitterDid string, revocRegDefId string, from int64, to int64) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	var did *C.char
	if len(submitterDid) > 0 {
		did = C.CString(submitterDid)
	} else {
		did = nil
	}

	/*
	   Builds a GET_REVOC_REG_DELTA request. Request to get the delta of the accumulated state of the Revocation Registry.
	    The Delta is defined by from and to timestamp fields.
	    If from is not specified, then the whole state till to will be returned.

	    :param submitter_did: (Optional) DID of the read request sender (if not provided then default Libindy DID will be used).
	    :param revoc_reg_def_id:  ID of the corresponding Revocation Registry Definition in ledger.
	    :param from_: Requested time represented as a total number of seconds from Unix Epoch
	    :param to: Requested time represented as a total number of seconds from Unix Epoch
	    :return: Request result as json.
	*/

	// Call to indy function
	res := C.indy_build_get_revoc_reg_delta_request(commandHandle,
		did,
		C.CString(revocRegDefId),
		C.longlong(from),
		C.longlong(to),
		(C.cb_buildRequest)(unsafe.Pointer(C.buildRequestCB)))

	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

// BuildRevocRegEntryRequest     Builds a REVOC_REG_ENTRY request.  Request to add the RevocReg entry containing
//    the new accumulator value and issued/revoked indices.
//    This is just a delta of indices, not the whole list. So, it can be sent each time a new credential is issued/revoked.
func BuildRevocRegEntryRequest(submitterDid string, revocRegDefId string, revDefType string, value string) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
	   :param submitter_did: Identifier (DID) of the transaction author as base58-encoded string.
	                          Actual request sender may differ if Endorser is used (look at `append_request_endorser`)
	    :param revoc_reg_def_id:  ID of the corresponding RevocRegDef.
	    :param rev_def_type:  Revocation Registry type (only CL_ACCUM is supported for now).
	    :param value: Registry-specific data:
	       {
	           value: {
	               prevAccum: string - previous accumulator value.
	               accum: string - current accumulator value.
	               issued: array<number> - an array of issued indices.
	               revoked: array<number> an array of revoked indices.
	           },
	           ver: string - version revocation registry entry json

	       }
	    :return: Request result as json.
	*/

	// Call to indy function
	res := C.indy_build_revoc_reg_entry_request(commandHandle,
		C.CString(submitterDid),
		C.CString(revocRegDefId),
		C.CString(revDefType),
		C.CString(value),
		(C.cb_buildRequest)(unsafe.Pointer(C.buildRequestCB)))

	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

// BuildRevocRegDefRequest     Builds a REVOC_REG_DEF request. Request to add the definition of revocation registry
//    to an exists credential definition.
func BuildRevocRegDefRequest(submitterDid string, data string) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
	  :param submitter_did: Identifier (DID) of the transaction author as base58-encoded string.
	                          Actual request sender may differ if Endorser is used (look at `append_request_endorser`)
	    :param data: Revocation Registry data:
	      {
	          "id": string - ID of the Revocation Registry,
	          "revocDefType": string - Revocation Registry type (only CL_ACCUM is supported for now),
	          "tag": string - Unique descriptive ID of the Registry,
	          "credDefId": string - ID of the corresponding CredentialDefinition,
	          "value": Registry-specific data {
	              "issuanceType": string - Type of Issuance(ISSUANCE_BY_DEFAULT or ISSUANCE_ON_DEMAND),
	              "maxCredNum": number - Maximum number of credentials the Registry can serve.
	              "tailsHash": string - Hash of tails.
	              "tailsLocation": string - Location of tails file.
	              "publicKeys": <public_keys> - Registry's public key.
	          },
	          "ver": string - version of revocation registry definition json.
	      }

	    :return: Request result as json.
	*/

	// Call to indy function
	res := C.indy_build_revoc_reg_def_request(commandHandle,
		C.CString(submitterDid),
		C.CString(data),
		(C.cb_buildRequest)(unsafe.Pointer(C.buildRequestCB)))

	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

// BuildGetRevocRegRequest       Builds a GET_REVOC_REG request. Request to get the accumulated state of the Revocation Registry
//    by ID. The state is defined by the given timestamp.
func BuildGetRevocRegRequest(submitterDid string, revRegDefId string, timeStamp int64) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	var did *C.char
	if len(submitterDid) > 0 {
		did = C.CString(submitterDid)
	} else {
		did = nil
	}

	/*
		:param submitter_did: (Optional) DID of the read request sender (if not provided then default Libindy DID will be used).
		    :param revoc_reg_def_id:  ID of the corresponding Revocation Registry Definition in ledger.
		    :param timestamp: Requested time represented as a total number of seconds from Unix Epoch
		    :return: Request result as json.
	*/

	// Call to indy function
	res := C.indy_build_get_revoc_reg_request(commandHandle,
		did,
		C.CString(revRegDefId),
		C.longlong(timeStamp),
		(C.cb_buildRequest)(unsafe.Pointer(C.buildRequestCB)))

	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

// BuildGetRevocRegDefRequest    Builds a GET_REVOC_REG_DEF request. Request to get a revocation registry definition,
//    that Issuer creates for a particular Credential Definition.
func BuildGetRevocRegDefRequest(submitterDid string, revRegDefId string) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	var did *C.char
	if len(submitterDid) > 0 {
		did = C.CString(submitterDid)
	} else {
		did = nil
	}

	/*
		  	:param submitter_did: (Optional) DID of the read request sender (if not provided then default Libindy DID will be used).
			:param rev_reg_def_id: ID of Revocation Registry Definition in ledger.

		    :return: Request result as json.
	*/

	// Call to indy function
	res := C.indy_build_get_revoc_reg_def_request(commandHandle,
		did,
		C.CString(revRegDefId),
		(C.cb_buildRequest)(unsafe.Pointer(C.buildRequestCB)))

	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

// BuildGetCredDefRequest    Builds a GET_CRED_DEF request. Request to get a credential definition (in particular, public key),
//   that Issuer creates for a particular Credential Schema.
func BuildGetCredDefRequest(submitterDid string, credDefId string) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	var did *C.char
	if len(submitterDid) > 0 {
		did = C.CString(submitterDid)
	} else {
		did = nil
	}

	/*
	   :param submitter_did: (Optional) DID of the read request sender (if not provided then default Libindy DID will be used).
	   :param id_: Credential Definition Id in ledger.
	   :return: Request result as json.
	*/

	// Call to indy function
	res := C.indy_build_get_cred_def_request(commandHandle,
		did,
		C.CString(credDefId),
		(C.cb_buildRequest)(unsafe.Pointer(C.buildRequestCB)))

	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

// BuildSchemaRequest Builds a SCHEMA request. Request to add Credential's schema.
func BuildSchemaRequest(submitterDid string, schema string) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
			:param submitter_did: Identifier (DID) of the transaction author as base58-encoded string.
		                          Actual request sender may differ if Endorser is used (look at `append_request_endorser`)
		    :param data: Credential schema.
		                 {
		                     id: identifier of schema
		                     attrNames: array of attribute name strings (the number of attributes should be less or equal than 125)
		                     name: Schema's name string
		                     version: Schema's version string,
		                     ver: Version of the Schema json
		                 }
		    :return: Request result as json.
	*/

	// Call to indy function
	res := C.indy_build_schema_request(commandHandle,
		C.CString(submitterDid),
		C.CString(schema),
		(C.cb_buildRequest)(unsafe.Pointer(C.buildRequestCB)))

	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

// BuildGetSchemaRequest builds a request to get crendential's schema
func BuildGetSchemaRequest(submitterDid string, schemaId string) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	// If we get empty value we send nil
	var res C.indy_error_t

	/*
			:param submitter_did: (Optional) DID of the read request sender (if not provided then default Libindy DID will be used).
		    :param id_: Schema Id in ledger
		    :return: Request result as json.
	*/

	// Call to indy function
	if len(submitterDid) != 0 {
		res = C.indy_build_get_schema_request(commandHandle,
			C.CString(submitterDid),
			C.CString(schemaId),
			(C.cb_buildRequest)(unsafe.Pointer(C.buildRequestCB)))
	} else {
		res = C.indy_build_get_schema_request(commandHandle,
			nil,
			C.CString(schemaId),
			(C.cb_buildRequest)(unsafe.Pointer(C.buildRequestCB)))
	}

	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

// 	BuildCredentialDefinitionRequest Builds an CRED_DEF request. Request to add a credential definition (in particular, public key),
//  that Issuer creates for a particular Credential Schema.
func BuildCredentialDefinitionRequest(submitterDid string, credDefJson string) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
			 :param submitter_did: Identifier (DID) of the transaction author as base58-encoded string.
		                          Actual request sender may differ if Endorser is used (look at `append_request_endorser`)
		    :param data: credential definition json
		                 {
		                     id: string - identifier of credential definition
		                     schemaId: string - identifier of stored in ledger schema
		                     type: string - type of the credential definition. CL is the only supported type now.
		                     tag: string - allows to distinct between credential definitions for the same issuer and schema
		                     value: Dictionary with Credential Definition's data: {
		                         primary: primary credential public key,
		                         Optional<revocation>: revocation credential public key
		                     },
		                     ver: Version of the CredDef json
		                 }
		    :return: Request result as json.
	*/

	// Call to indy function
	res := C.indy_build_cred_def_request(commandHandle,
		C.CString(submitterDid),
		C.CString(credDefJson),
		(C.cb_buildRequest)(unsafe.Pointer(C.buildRequestCB)))

	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

// BuildGetDdoRequest Builds a request to get a DDO.
func BuildGetDdoRequest(submitterDid string, targetDid string) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
		Builds a request to get a DDO.

	    :param submitter_did: (Optional) DID of the read request sender (if not provided then default Libindy DID will be used).
	    :param target_did: Id of Identity stored in secured Wallet.

	    :return: Request result as json.
	 */

	// Call to indy function
	res := C.indy_build_get_ddo_request(commandHandle,
		C.CString(submitterDid),
		C.CString(targetDid),
		(C.cb_buildRequest)(unsafe.Pointer(C.buildRequestCB)))

	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

// BuildNymRequest creates a request to make an DID public on the blockchain
func BuildNymRequest(submitterDid string, targetDid string, verkey string, alias string, role string) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	// If we get empty value we send nil
	var res C.indy_error_t

	/*
		Builds a NYM request.

		:param submitter_did: Identifier (DID) of the transaction author as base58-encoded string.
			Actual request sender may differ if Endorser is used (look at `append_request_endorser`)
		:param target_did: Target DID as base58-encoded string for 16 or 32 bit DID value.
		:param ver_key: Target identity verification key as base58-encoded string.
		:param alias: NYM's alias.
		:param role: Role of a user NYM record:
		null (common USER)
		TRUSTEE
		STEWARD
		TRUST_ANCHOR
		ENDORSER - equal to TRUST_ANCHOR that will be removed soon
		NETWORK_MONITOR
		empty string to reset role
		:return: Request result as json.*/

	// Call to indy function
	if len(alias) != 0 {
		res = C.indy_build_nym_request(commandHandle,
			C.CString(submitterDid),
			C.CString(targetDid),
			C.CString(verkey),
			C.CString(alias),
			C.CString(role),
			(C.cb_buildRequest)(unsafe.Pointer(C.buildRequestCB)))
	} else {
		res = C.indy_build_nym_request(commandHandle,
			C.CString(submitterDid),
			C.CString(targetDid),
			C.CString(verkey),
			nil,
			C.CString(role),
			(C.cb_buildRequest)(unsafe.Pointer(C.buildRequestCB)))
	}

	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

// BuildAttribRequest Builds an ATTRIB request. Request to add attribute to a NYM record.
func BuildAttribRequest(submitterDid string, targetDid string, hash string, raw string, enc string) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	var hashData *C.char
	if len(hash) > 0 {
		hashData = C.CString(hash)
	} else {
		hashData = nil
	}

	var rawData *C.char
	if len(raw) > 0 {
		rawData = C.CString(raw)
	} else {
		rawData = nil
	}

	var encData *C.char
	if len(enc) > 0 {
		encData = C.CString(enc)
	} else {
		encData = nil
	}
	/*
		Builds an ATTRIB request. Request to add attribute to a NYM record.

	    :param submitter_did: Identifier (DID) of the transaction author as base58-encoded string.
	                          Actual request sender may differ if Endorser is used (look at `append_request_endorser`)
	    :param target_did: Target DID as base58-encoded string for 16 or 32 bit DID value.
	    :param hash: (Optional) Hash of attribute data.
	    :param raw: (Optional) Json, where key is attribute name and value is attribute value.
	    :param enc: (Optional) Encrypted value attribute data.
	    :return: Request result as json.
	 */

	// Call to indy function
	res := C.indy_build_attrib_request(commandHandle,
		C.CString(submitterDid),
		C.CString(targetDid),
		hashData,
		rawData,
		encData,
		(C.cb_buildRequest)(unsafe.Pointer(C.buildRequestCB)))

	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

// BuildGetAttribRequest Builds a GET_ATTRIB request. Request to get information about an Attribute for the specified DID.
func BuildGetAttribRequest(submitterDid string, targetDid string, hash string, raw string, enc string) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	var did *C.char
	if len(submitterDid) > 0 {
		did = C.CString(submitterDid)
	} else {
		did = nil
	}

	var hashData *C.char
	if len(hash) > 0 {
		hashData = C.CString(hash)
	} else {
		hashData = nil
	}

	var rawData *C.char
	if len(raw) > 0 {
		rawData = C.CString(raw)
	} else {
		rawData = nil
	}

	var encData *C.char
	if len(enc) > 0 {
		encData = C.CString(enc)
	} else {
		encData = nil
	}
	/*
		Builds a GET_ATTRIB request. Request to get information about an Attribute for the specified DID.

		:param submitter_did: (Optional) DID of the read request sender (if not provided then default Libindy DID will be used).
	    :param target_did: Target DID as base58-encoded string for 16 or 32 bit DID value.
	    :param xhash: (Optional) Requested attribute name.
	    :param raw: (Optional) Requested attribute hash.
	    :param enc: (Optional) Requested attribute encrypted value.
	    :return: Request result as json.
	*/

	// Call to indy function
	res := C.indy_build_get_attrib_request(commandHandle,
		did,
		C.CString(targetDid),
		hashData,
		rawData,
		encData,
		(C.cb_buildRequest)(unsafe.Pointer(C.buildRequestCB)))

	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

// BuildGetNymRequest Builds a GET_NYM request. Request to get information about a DID (NYM).
func BuildGetNymRequest(submitterDid string, targetDid string) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	var did *C.char
	if len(submitterDid) > 0 {
		did = C.CString(submitterDid)
	} else {
		did = nil
	}
	/*
		Builds a GET_NYM request. Request to get information about a DID (NYM).

	    :param submitter_did: (Optional) DID of the read request sender (if not provided then default Libindy DID will be used).
	    :param target_did: Target DID as base58-encoded string for 16 or 32 bit DID value.
	    :return: Request result as json.
	 */

	// Call to indy function
	res := C.indy_build_get_nym_request(commandHandle,
		did,
		C.CString(targetDid),
		(C.cb_buildRequest)(unsafe.Pointer(C.buildRequestCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

// BuildNodeRequest Builds a NODE request. Request to add a new node to the pool, or updates existing in the pool.
func BuildNodeRequest(submitterDid string, targetDid string, data string) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)
	/*
		Builds a NODE request. Request to add a new node to the pool, or updates existing in the pool.

	    :param submitter_did: Identifier (DID) of the transaction author as base58-encoded string.
	                          Actual request sender may differ if Endorser is used (look at `append_request_endorser`)
	    :param target_did: Target Node's DID.  It differs from submitter_did field.
	    :param data: Data associated with the Node:
	      {
	          alias: string - Node's alias
	          blskey: string - (Optional) BLS multi-signature key as base58-encoded string.
	          blskey_pop: string - (Optional) BLS key proof of possession as base58-encoded string.
	          client_ip: string - (Optional) Node's client listener IP address.
	          client_port: string - (Optional) Node's client listener port.
	          node_ip: string - (Optional) The IP address other Nodes use to communicate with this Node.
	          node_port: string - (Optional) The port other Nodes use to communicate with this Node.
	          services: array<string> - (Optional) The service of the Node. VALIDATOR is the only supported one now.
	      }
	    :return: Request result as json.
	 */

	// Call to indy function
	res := C.indy_build_node_request(commandHandle,
		C.CString(submitterDid),
		C.CString(targetDid),
		C.CString(data),
		(C.cb_buildRequest)(unsafe.Pointer(C.buildRequestCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

// BuildGetValidatorInfoRequest Builds a GET_VALIDATOR_INFO request.
func BuildGetValidatorInfoRequest(submitterDid string) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
		Builds a GET_VALIDATOR_INFO request.

	    :param submitter_did: Id of Identity stored in secured Wallet.
	    :return: Request result as json.
	 */

	// Call to indy function
	res := C.indy_build_get_validator_info_request(commandHandle,
		C.CString(submitterDid),
		(C.cb_buildRequest)(unsafe.Pointer(C.buildRequestCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

// BuildGetTxnRequest Builds a GET_TXN request. Request to get any transaction by its seq_no.
func BuildGetTxnRequest(submitterDid string, ledgerType string, seqNo int) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	var did, _type *C.char
	if len(submitterDid) > 0 {
		did = C.CString(submitterDid)
	} else {
		did = nil
	}

	if len(ledgerType) > 0 {
		_type = C.CString(ledgerType)
	} else {
		_type = nil
	}

	/*
		Builds a GET_TXN request. Request to get any transaction by its seq_no.
	    :param submitter_did: (Optional) DID of the read request sender (if not provided then default Libindy DID will be used).
	    :param ledger_type: (Optional) type of the ledger the requested transaction belongs to:
	        DOMAIN - used default,
	        POOL,
	        CONFIG
	        any number
	    :param seq_no: requested transaction sequence number as it's stored on Ledger.
	    :return: Request result as json.
	 */

	// Call to indy function
	res := C.indy_build_get_txn_request(commandHandle,
		did,
		_type,
		(C.indy_i32_t)(seqNo),
		(C.cb_buildRequest)(unsafe.Pointer(C.buildRequestCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}
	return future
}

// BuildPoolConfigRequest Builds a POOL_CONFIG request. Request to change Pool's configuration.
func BuildPoolConfigRequest(submitterDid string, writes bool, force bool) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
		Builds a POOL_CONFIG request. Request to change Pool's configuration.

	    :param submitter_did: Identifier (DID) of the transaction author as base58-encoded string.
	                          Actual request sender may differ if Endorser is used (look at `append_request_endorser`)
	    :param writes: Whether any write requests can be processed by the pool
	                   (if false, then pool goes to read-only state). True by default.
	    :param force: Whether we should apply transaction (for example, move pool to read-only state)
	                  without waiting for consensus of this transaction
	    :return: Request result as json.
	 */

	// Call to indy function
	res := C.indy_build_pool_config_request(commandHandle,
		C.CString(submitterDid),
		(C.indy_bool_t)(writes),
		(C.indy_bool_t)(force),
		(C.cb_buildRequest)(unsafe.Pointer(C.buildRequestCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}
	return future
}

// BuildPoolRestartRequest Builds a POOL_RESTART request.
func BuildPoolRestartRequest(submitterDid string, action string, dateTime string) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
		Builds a POOL_RESTART request
	    :param submitter_did: Identifier (DID) of the transaction author as base58-encoded string.
	                          Actual request sender may differ if Endorser is used (look at `append_request_endorser`)
	    :param action       : Action that pool has to do after received transaction.
	                          Can be "start" or "cancel"
	    :param datetime     : Time when pool must be restarted.
	*/

	// Call to indy function
	res := C.indy_build_pool_restart_request(commandHandle,
		C.CString(submitterDid),
		C.CString(action),
		C.CString(dateTime),
		(C.cb_buildRequest)(unsafe.Pointer(C.buildRequestCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}
	return future
}

// BuildPoolUpgradeRequest Builds a POOL_UPGRADE request. Request to upgrade the Pool (sent by Trustee).
func BuildPoolUpgradeRequest(submitterDid string, name string, version string, action string, sha256 string, timeOut int32, schedule string,
	justification string, reinstall bool, force bool, package_ string) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	var indyTimeOut C.indy_i32_t
	if timeOut != 0 {
		indyTimeOut = C.indy_i32_t(timeOut)
	} else {
		indyTimeOut = 0
	}

	var indySchedule *C.char
	if len(schedule) > 0 {
		indySchedule = C.CString(schedule)
	} else {
		indySchedule = nil
	}

	var indyJustification *C.char
	if len(justification) > 0 {
		indyJustification = C.CString(justification)
	} else {
		indySchedule = nil
	}

	var indyPackage *C.char
	if len(package_) > 0 {
		indyPackage = C.CString(package_)
	} else {
		indyPackage = nil
	}

	/*
			Builds a POOL_UPGRADE request. Request to upgrade the Pool (sent by Trustee).
			It upgrades the specified Nodes (either all nodes in the Pool, or some specific ones).

			:param submitter_did: Identifier (DID) of the transaction author as base58-encoded string.
								  Actual request sender may differ if Endorser is used (look at `append_request_endorser`)
			:param name: Human-readable name for the upgrade.
			:param version: The version of indy-node package we perform upgrade to.
							Must be greater than existing one (or equal if reinstall flag is True).
			:param action: Either start or cancel.
			:param _sha256: sha256 hash of the package.
			:param _timeout: (Optional) Limits upgrade time on each Node.
			:param schedule: (Optional) Schedule of when to perform upgrade on each node. Map Node DIDs to upgrade time.
			:param justification: (Optional) justification string for this particular Upgrade.
			:param reinstall: Whether it's allowed to re-install the same version. False by default.
			:param force: Whether we should apply transaction (schedule Upgrade) without waiting
						  for consensus of this transaction.
			:param package: (Optional) Package to be upgraded.
			:return: Request result as json.
	*/

	// Call to indy function
	res := C.indy_build_pool_upgrade_request(commandHandle,
		C.CString(submitterDid),
		C.CString(name),
		C.CString(version),
		C.CString(action),
		C.CString(sha256),
		indyTimeOut,
		indySchedule,
		indyJustification,
		(C.indy_bool_t)(reinstall),
		(C.indy_bool_t)(force),
		indyPackage,
		(C.cb_buildRequest)(unsafe.Pointer(C.buildRequestCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}
	return future
}

// BuildAuthRuleRequest Builds a AUTH_RULE request.
func BuildAuthRuleRequest(submitterDid string, txnType string, action string, field string, oldValue string, newValue string, constraint string) chan indyUtils.IndyResult {

	// Prepare the call parameters.
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	var indyOldValue *C.char
	if len(oldValue) > 0 {
		indyOldValue = C.CString(oldValue)
	} else {
		indyOldValue = nil
	}

	var indyNewValue *C.char
	if len(newValue) > 0 {
		indyNewValue = C.CString(newValue)
	} else {
		indyNewValue = nil
	}

	/*
			Builds a AUTH_RULE request. Request to change authentication rules for a ledger transaction.

		    :param submitter_did: Identifier (DID) of the transaction author as base58-encoded string.
		                          Actual request sender may differ if Endorser is used (look at `append_request_endorser`)
		    :param txn_type: ledger transaction alias or associated value.
		    :param action: type of an action.
		       Can be either "ADD" (to add a new rule) or "EDIT" (to edit an existing one).
		    :param field: transaction field.
		    :param old_value: (Optional) old value of a field, which can be changed to a new_value (mandatory for EDIT action).
		    :param new_value: (Optional) new value that can be used to fill the field.
		    :param constraint: set of constraints required for execution of an action in the following format:
		        {
		            constraint_id - <string> type of a constraint.
		                Can be either "ROLE" to specify final constraint or  "AND"/"OR" to combine constraints.
		            role - <string> (optional) role of a user which satisfy to constrain.
		            sig_count - <u32> the number of signatures required to execution action.
		            need_to_be_owner - <bool> (optional) if user must be an owner of transaction (false by default).
		            off_ledger_signature - <bool> (optional) allow signature of unknow for ledger did (false by default).
		            metadata - <object> (optional) additional parameters of the constraint.
		        }
		      can be combined by
		        {
		            'constraint_id': <"AND" or "OR">
		            'auth_constraints': [<constraint_1>, <constraint_2>]
		        }
		    Default ledger auth rules: https://github.com/hyperledger/indy-node/blob/master/docs/source/auth_rules.md
		    More about AUTH_RULE request: https://github.com/hyperledger/indy-node/blob/master/docs/source/requests.md#auth_rule
		    :return: Request result as json.
	*/

	// Call to indy function
	res := C.indy_build_auth_rule_request(commandHandle,
		C.CString(submitterDid),
		C.CString(txnType),
		C.CString(action),
		C.CString(field),
		indyOldValue,
		indyNewValue,
		C.CString(constraint),
		(C.cb_buildRequest)(unsafe.Pointer(C.buildRequestCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}
	return future
}

// BuildAuthRulesRequest Builds a AUTH_RULES request.
func BuildAuthRulesRequest(submitterDid string, data string) chan indyUtils.IndyResult {

	// Prepare the call parameters.
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
		 Builds a AUTH_RULES request. Request to change multiple authentication rules for a ledger transaction.
	    :param submitter_did: Identifier (DID) of the transaction author as base58-encoded string.
	                          Actual request sender may differ if Endorser is used (look at `append_request_endorser`)
	    :param data: a list of auth rules: [
	        {
	            "auth_type": ledger transaction alias or associated value,
	            "auth_action": type of an action,
	            "field": transaction field,
	            "old_value": (Optional) old value of a field, which can be changed to a new_value (mandatory for EDIT action),
	            "new_value": (Optional) new value that can be used to fill the field,
	            "constraint": set of constraints required for execution of an action in the format described above for `build_auth_rule_request` function.
	        }
	    ]
	    Default ledger auth rules: https://github.com/hyperledger/indy-node/blob/master/docs/source/auth_rules.md
	    More about AUTH_RULE request: https://github.com/hyperledger/indy-node/blob/master/docs/source/requests.md#auth_rules
	    :return: Request result as json.
	 */

	// Call to indy function
	res := C.indy_build_auth_rules_request(commandHandle,
		C.CString(submitterDid),
		C.CString(data),
		(C.cb_buildRequest)(unsafe.Pointer(C.buildRequestCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}
	return future
}

// BuildGetAuthRuleRequest Builds a GET_AUTH_RULE request. Request to get authentication rules for a ledger transaction.
func BuildGetAuthRuleRequest(submitterDid string, txnType string, action string, field string, oldValue string, newValue string) chan indyUtils.IndyResult {

	// Prepare the call parameters.
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	var did *C.char
	if len(submitterDid) > 0 {
		did = C.CString(submitterDid)
	} else {
		did = nil
	}

	var indyTxnType *C.char
	if len(txnType) > 0 {
		indyTxnType = C.CString(txnType)
	} else {
		indyTxnType = nil
	}

	var indyAction *C.char
	if len(action) > 0 {
		indyAction = C.CString(action)
	} else {
		indyAction = nil
	}

	var indyField *C.char
	if len(field) > 0 {
		indyField = C.CString(field)
	} else {
		indyField = nil
	}

 	var indyOldValue *C.char
	if len(oldValue) > 0 {
		indyOldValue = C.CString(oldValue)
	} else {
		indyOldValue = nil
	}

	var indyNewValue *C.char
	if len(newValue) > 0 {
		indyNewValue = C.CString(newValue)
	} else {
		indyNewValue = nil
	}

	/*
			 Builds a GET_AUTH_RULE request. Request to get authentication rules for a ledger transaction.
	   		 NOTE: Either none or all transaction related parameters must be specified (`old_value` can be skipped for `ADD` action).
				* none - to get all authentication rules for all ledger transactions
				* all - to get authentication rules for specific action (`old_value` can be skipped for `ADD` action)

			:param submitter_did: (Optional) DID of the read request sender (if not provided then default Libindy DID will be used).
			:param txn_type: (Optional) target ledger transaction alias or associated value.
			:param action: (Optional) target action type. Can be either "ADD" or "EDIT".
			:param field: (Optional) target transaction field.
			:param old_value: (Optional) old value of field, which can be changed to a new_value (must be specified for EDIT action).
			:param new_value: (Optional) new value that can be used to fill the field.
			:return: Request result as json.
	*/

	// Call to indy function
	res := C.indy_build_get_auth_rule_request(commandHandle,
		did,
		indyTxnType,
		indyAction,
		indyField,
		indyOldValue,
		indyNewValue,
		(C.cb_buildRequest)(unsafe.Pointer(C.buildRequestCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}
	return future
}

// BuildTxnAuthorAgreementRequest Builds a TXN_AUTHR_AGRMT request. Request to add a new version of Transaction Author Agreement to the ledger.
func BuildTxnAuthorAgreementRequest(submitterDid string, text string, version string, ratificationTs int64, retirementTs int64) chan indyUtils.IndyResult {

	// Prepare the call parameter.
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	var indyText *C.char
	if len(text) > 0 {
		indyText = C.CString(text)
	} else {
		indyText = nil
	}

	var indyRatificationTs C.longlong
	if ratificationTs > 0 {
		indyRatificationTs = C.longlong(ratificationTs)
	} else {
		indyRatificationTs = -1
	}

	var indyRetirementTs C.longlong
	if retirementTs > 0 {
		indyRetirementTs = C.longlong(retirementTs)
	} else {
		indyRetirementTs = -1
	}

	/*
		Builds a TXN_AUTHR_AGRMT request. Request to add a new version of Transaction Author Agreement to the ledger.
	    EXPERIMENTAL

	    :param submitter_did: Identifier (DID) of the transaction author as base58-encoded string.
	                          Actual request sender may differ if Endorser is used (look at `append_request_endorser`)
	    :param text: (Optional) a content of the TTA.
	                          Mandatory in case of adding a new TAA. An existing TAA text can not be changed.
	                          for Indy Node version <= 1.12.0:
	                              Use empty string to reset TAA on the ledger
	                          for Indy Node version > 1.12.0
	                              Should be omitted in case of updating an existing TAA (setting `retirement_ts`)
	    :param version: a version of the TTA (unique UTF-8 string).
	    :param ratification_ts: (Optional) the date (timestamp) of TAA ratification by network government.
	                          for Indy Node version <= 1.12.0:
	                             Must be omitted
	                          for Indy Node version > 1.12.0:
	                             Must be specified in case of adding a new TAA
	                             Can be omitted in case of updating an existing TAA
	    :param retirement_ts: (Optional) the date (timestamp) of TAA retirement.
	                          for Indy Node version <= 1.12.0:
	                              Must be omitted
	                          for Indy Node version > 1.12.0:
	                              Must be omitted in case of adding a new (latest) TAA.
	                              Should be used for updating (deactivating) non-latest TAA on the ledger.
	    Note: Use `build_disable_all_txn_author_agreements_request` to disable all TAA's on the ledger.
	    :return: Request result as json.
	 */

	// Call to indy function
	res := C.indy_build_txn_author_agreement_request(commandHandle,
		C.CString(submitterDid),
		indyText,
		C.CString(version),
		indyRatificationTs,
		indyRetirementTs,
		(C.cb_buildRequest)(unsafe.Pointer(C.buildRequestCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}
	return future
}

// BuildDisableAllTxnAuthorAgreementsRequest Builds a DISABLE_ALL_TXN_AUTHR_AGRMTS request. Request to disable all Transaction Author Agreement on the ledger.
func BuildDisableAllTxnAuthorAgreementsRequest(submitterDid string) chan indyUtils.IndyResult {

	// Prepare the call parameters.
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
		Builds a DISABLE_ALL_TXN_AUTHR_AGRMTS request. Request to disable all Transaction Author Agreement on the ledger.
	    EXPERIMENTAL

	    :param submitter_did: Identifier (DID) of the transaction author as base58-encoded string.
	                          Actual request sender may differ if Endorser is used (look at `append_request_endorser`)
	    :return: Request result as json.
	 */

	// Call to indy function.
	res := C.indy_build_disable_all_txn_author_agreements_request(commandHandle,
		C.CString(submitterDid),
		(C.cb_buildRequest)(unsafe.Pointer(C.buildRequestCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}
	return future
}

// BuildGetTxnAuthorAgreementRequest  Builds a GET_TXN_AUTHR_AGRMT request. Request to get a specific Transaction Author Agreement from the ledger.
func BuildGetTxnAuthorAgreementRequest(submitterDid string, data string) chan indyUtils.IndyResult {

	// Prepare the call parameters.
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	var did *C.char
	if len(submitterDid) > 0 {
		did = C.CString(submitterDid)
	} else {
		did = nil
	}

	var indyData *C.char
	if len(data) > 0 {
		indyData = C.CString(data)
	} else {
		indyData = nil
	}

	/*
	   Builds a GET_TXN_AUTHR_AGRMT request. Request to get a specific Transaction Author Agreement from the ledger.
	   EXPERIMENTAL

	   :param submitter_did: (Optional) DID of the read request sender (if not provided then default Libindy DID will be used).
	   :param data: (Optional) specifies a condition for getting specific TAA.
	    Contains 3 mutually exclusive optional fields:
	    {
	        hash: Optional<str> - hash of requested TAA,
	        version: Optional<str> - version of requested TAA.
	        timestamp: Optional<i64> - ledger will return TAA valid at requested timestamp.
	    }
	    Null data or empty JSON are acceptable here. In this case, ledger will return the latest version of TAA.

	   :return: Request result as json.
	 */

	// Call to indy function.
	res := C.indy_build_get_txn_author_agreement_request(commandHandle,
		did,
		indyData,
		(C.cb_buildRequest)(unsafe.Pointer(C.buildRequestCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}
	return future
}

// BuildAcceptanceMechanismsRequest Builds a SET_TXN_AUTHR_AGRMT_AML request. Request to add a new list of acceptance mechanisms for transaction author agreement.
func BuildAcceptanceMechanismsRequest(submitterDid string, aml string, version string, amlContext string) chan indyUtils.IndyResult {

	// Prepare the call parameters.
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	var indyAmlContext *C.char
	if len(amlContext) > 0 {
		indyAmlContext = C.CString(amlContext)
	} else {
		indyAmlContext = nil
	}

	/*
	 	Builds a SET_TXN_AUTHR_AGRMT_AML request. Request to add a new list of acceptance mechanisms for transaction author agreement.
	    Acceptance Mechanism is a description of the ways how the user may accept a transaction author agreement.
	    EXPERIMENTAL

	    :param submitter_did: Identifier (DID) of the transaction author as base58-encoded string.
	                          Actual request sender may differ if Endorser is used (look at `append_request_endorser`)
	    :param aml: a set of new acceptance mechanisms:
	    {
	        “<acceptance mechanism label 1>”: { acceptance mechanism description 1},
	        “<acceptance mechanism label 2>”: { acceptance mechanism description 2},
	        ...
	    }
	    :param version: a version of new acceptance mechanisms. (Note: unique on the Ledger)
	    :param aml_context: (Optional) common context information about acceptance mechanisms (may be a URL to external resource).
	    :return: Request result as json.
	 */

	// Call to indy function.
	res := C.indy_build_acceptance_mechanisms_request(commandHandle,
		C.CString(submitterDid),
		C.CString(aml),
		C.CString(version),
		indyAmlContext,
		(C.cb_buildRequest)(unsafe.Pointer(C.buildRequestCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}
	return future
}

// BuildGetAcceptanceMechanismsRequest Builds a GET_TXN_AUTHR_AGRMT_AML request.
func BuildGetAcceptanceMechanismsRequest(submitterDid string, timestamp int64, version string) chan indyUtils.IndyResult {

	// Prepare the call parameters.
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	var did *C.char
	if len(submitterDid) > 0 {
		did = C.CString(submitterDid)
	} else {
		did = nil
	}

	var indyVersion *C.char
	if len(version) > 0 {
		indyVersion = C.CString(version)
	} else {
		indyVersion = nil
	}

	/*
	   Builds a GET_TXN_AUTHR_AGRMT_AML request. Request to get a list of  acceptance mechanisms from the ledger
	   valid for specified time or the latest one.
	   EXPERIMENTAL

	   :param submitter_did: (Optional) DID of the read request sender (if not provided then default Libindy DID will be used).
	   :param timestamp: time to get an active acceptance mechanisms. Pass -1 to get the latest one.
	   :param version: (Optional) version of acceptance mechanisms.

	   NOTE: timestamp and version cannot be specified together.
	   :return: Request result as json.
	 */

	// Call to indy function.
	res := C.indy_build_get_acceptance_mechanisms_request(commandHandle,
		did,
		C.longlong(timestamp),
		indyVersion,
		(C.cb_buildRequest)(unsafe.Pointer(C.buildRequestCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}
	return future
}

//export signAndSubmitRequestCB
func signAndSubmitRequestCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, response *C.char) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle),
			indyUtils.IndyResult{Error: nil,
				Results: []interface{}{
					string(C.GoString(response)),
				}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// SignAndSubmitRequest sends a request to the blockchain
func SignAndSubmitRequest(poolHandle int, walletHandle int, submitterDid string, requestJson string) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
	   Signs and submits request message to validator pool.

	   Adds submitter information to passed request json, signs it with submitter
	   sign key (see wallet_sign), and sends signed request message
	   to validator pool (see write_request).

	   :param pool_handle: pool handle (created by open_pool_ledger).
	   :param wallet_handle: wallet handle (created by open_wallet).
	   :param submitter_did: Id of Identity stored in secured Wallet.
	   :param request_json: Request data json.
	   :return: Request result as json.
	*/

	// Call to indy function
	res := C.indy_sign_and_submit_request(commandHandle,
		(C.indy_handle_t)(poolHandle),
		(C.indy_handle_t)(walletHandle),
		C.CString(submitterDid),
		C.CString(requestJson),
		(C.cb_signAndSubmitRequest)(unsafe.Pointer(C.signAndSubmitRequestCB)))

	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

// SubmitRequest Publishes request message to validator pool (no signing, unlike sign_and_submit_request).
// The request is sent to the validator pool as is. It's assumed that it's already prepared.
func SubmitRequest(poolHandle int, requestJson string) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
	   :param pool_handle: pool handle (created by open_pool_ledger).
	    :param request_json: Request data json.
	    :return: Request result as json.
	*/

	// Call to indy function
	res := C.indy_submit_request(commandHandle,
		(C.indy_handle_t)(poolHandle),
		C.CString(requestJson),
		(C.cb_signAndSubmitRequest)(unsafe.Pointer(C.signAndSubmitRequestCB)))

	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

// SignRequest Signs request message. Adds submitter information to passed request json, signs it with submitter sign key
func SignRequest(walletHandle int, submitterDid string, requestJson string) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
	 	Signs request message.
	    Adds submitter information to passed request json, signs it with submitter
	    sign key (see wallet_sign).

	    :param wallet_handle: wallet handle (created by open_wallet).
	    :param submitter_did: Id of Identity stored in secured Wallet.
	    :param request_json: Request data json.
	    :return: Signed request json.
	*/

	// Call C.indy_sign_request
	res := C.indy_sign_request(commandHandle,
		(C.indy_handle_t)(walletHandle),
		C.CString(submitterDid),
		C.CString(requestJson),
		(C.cb_signAndSubmitRequest)(unsafe.Pointer(C.signAndSubmitRequestCB)))

	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

// GetResponseMetadata Parse transaction response to fetch metadata.
func GetResponseMetadata(response string) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)
	/*
		 Parse transaction response to fetch metadata.
	     The important use case for this method is validation of Node's response freshens.
	     Distributed Ledgers can reply with outdated information for consequence read request after write.

	     To reduce pool load libindy sends read requests to one random node in the pool.
	     Consensus validation is performed based on validation of nodes multi signature for current ledger Merkle Trie root.
	     This multi signature contains information about the latest ldeger's transaction ordering time and sequence number that this method returns.
	     If node that returned response for some reason is out of consensus and has outdated ledger
	     it can be caught by analysis of the returned latest ledger's transaction ordering time and sequence number.

	     There are two ways to filter outdated responses:
	         1) based on "seqNo" - sender knows the sequence number of transaction that he consider as a fresh enough.
	         2) based on "txnTime" - sender knows the timestamp that he consider as a fresh enough.

	     Note: response of GET_VALIDATOR_INFO request isn't supported
	    :param response: response of write or get request.
	    :return: Response Metadata.
	    {
	        "seqNo": Option<u64> - transaction sequence number,
	        "txnTime": Option<u64> - transaction ordering time,
	        "lastSeqNo": Option<u64> - the latest transaction seqNo for particular Node,
	        "lastTxnTime": Option<u64> - the latest transaction ordering time for particular Node
	    }
	 */

	// Call to indy function
	res := C.indy_get_response_metadata(commandHandle,
		C.CString(response),
		(C.cb_buildRequest)(unsafe.Pointer(C.buildRequestCB)))

	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

// AppendTxnAuthorAgreementAcceptanceToRequest Append transaction author agreement acceptance data to a request.
func AppendTxnAuthorAgreementAcceptanceToRequest(requestJson string, text string, version string, taaDigest string, mechanism string, time int64) chan indyUtils.IndyResult {

	// Prepare the call parameters.
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	var indyText *C.char
	if len(text) > 0 {
		indyText = C.CString(text)
	} else {
		indyText = nil
	}

	var indyVersion *C.char
	if len(version) > 0 {
		indyVersion = C.CString(version)
	} else {
		indyVersion = nil
	}

	var indyTaaDigest *C.char
	if len(taaDigest) > 0 {
		indyTaaDigest = C.CString(taaDigest)
	} else {
		indyTaaDigest = nil
	}

	/*
	   Append transaction author agreement acceptance data to a request.
	   This function should be called before signing and sending a request if there is any transaction author agreement set on the Ledger.

	   EXPERIMENTAL

	   This function may calculate hash by itself or consume it as a parameter.
	   If all text, version and taa_digest parameters are specified, a check integrity of them will be done.
	   :param request_json: original request data json.
	   :param text and version: (Optional) raw data about TAA from ledger.
	              These parameters should be passed together.
	              These parameters are required if taa_digest parameter is omitted.
	   :param taa_digest: (Optional) digest on text and version.
	                     Digest is sha256 hash calculated on concatenated strings: version || text.
	                     This parameter is required if text and version parameters are omitted.
	   :param mechanism: mechanism how user has accepted the TAA
	   :param time: UTC timestamp when user has accepted the TAA. Note that the time portion will be discarded to avoid a privacy risk.

	   :return: Updated request result as json.
	 */

	// Call to indy function
	res := C.indy_append_txn_author_agreement_acceptance_to_request(commandHandle,
		C.CString(requestJson),
		indyText,
		indyVersion,
		indyTaaDigest,
		C.CString(mechanism),
		C.ulonglong(time),
		(C.cb_buildRequest)(unsafe.Pointer(C.buildRequestCB)))

	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}