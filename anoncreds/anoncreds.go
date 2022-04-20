/*
// ******************************************************************
// Purpose: Wrapper to call libindy, imports functions from indy_anoncreds.h
// Author:  alexandru.leonte@siemens.com
// Notes:
// Copyright (c): Siemens SRL
// This work is licensed under the terms of the Apache License Version 2.0.  See
// the LICENSE.txt file in the top-level directory.
// ******************************************************************
*/

package anoncreds

/*
#cgo CFLAGS: -I ../include
#cgo LDFLAGS: -L${SRCDIR}/../lib -lindy
#include <indy_core.h>
typedef void (*cb_issuerCreateSchema)(indy_handle_t, indy_error_t, char*, char*);
extern void issuerCreateSchemaCB(indy_handle_t, indy_error_t, char*, char*);

typedef void (*cb_issuerCreateAndStoreCredentialDef)(indy_handle_t, indy_error_t, char*, char*);
extern void issuerCreateAndStoreCredentialDefCB(indy_handle_t, indy_error_t, char*, char*);

typedef void (*cb_issuerRotateCredentialDefStart)(indy_handle_t, indy_error_t, char*);
extern void issuerRotateCredentialDefStartCB(indy_handle_t, indy_error_t, char*);

typedef void (*cb_issuerRotateCredentialDefApply)(indy_handle_t, indy_error_t);
extern void issuerRotateCredentialDefApplyCB(indy_handle_t, indy_error_t);

typedef void (*cb_issuerCreateCredentialOffer)(indy_handle_t, indy_error_t, char*);
extern void issuerCreateCredentialOfferCB(indy_handle_t, indy_error_t, char*);

typedef void (*cb_proverCreateCredentialRequest)(indy_handle_t, indy_error_t, char*, char*);
extern void proverCreateCredentialRequestCB(indy_handle_t, indy_error_t, char*, char*);

typedef void (*cb_proverCreateMasterSecret)(indy_handle_t, indy_error_t, char*);
extern void proverCreateMasterSecretCB(indy_handle_t, indy_error_t, char*);

typedef void (*cb_issuerCreateCredential)(indy_handle_t, indy_error_t, char*, char*, char*);
extern void issuerCreateCredentialCB(indy_handle_t, indy_error_t, char*, char*, char*);

typedef void (*cb_proverStoreCredential)(indy_handle_t, indy_error_t, char*);
extern void proverStoreCredentialCB(indy_handle_t, indy_error_t, char*);

typedef void (*cb_generateNonce)(indy_handle_t, indy_error_t, char*);
extern void generateNonceCB(indy_handle_t, indy_error_t, char*);

typedef void (*cb_proverSearchForCredentialsForProofReq)(indy_handle_t, indy_error_t, indy_handle_t);
extern void proverSearchForCredentialsForProofReqCB(indy_handle_t, indy_error_t, indy_handle_t);

typedef void (*cb_proverFetchCredentialsForProofReq)(indy_handle_t, indy_error_t, char*);
extern void proverFetchCredentialsForProofReqCB(indy_handle_t, indy_error_t, char*);

typedef void (*cb_proverCreateProof)(indy_handle_t, indy_error_t, char*);
extern void proverCreateProofCB(indy_handle_t, indy_error_t, char*);

typedef void (*cb_verifierVerifyProof)(indy_handle_t, indy_error_t, indy_bool_t);
extern void verifierVerifyProofCB(indy_handle_t, indy_error_t, indy_bool_t);

typedef void (*cb_proverGetCredential)(indy_handle_t, indy_error_t, char*);
extern void proverGetCredentialCB(indy_handle_t, indy_error_t, char*);

typedef void (*cb_createAndStoreRevocReg)(indy_handle_t, indy_error_t, char*, char*, char*);
extern void createAndStoreRevocRegCB(indy_handle_t, indy_error_t, char*, char*, char*);

typedef void (*cb_createRevocationState)(indy_handle_t, indy_error_t, char*);
extern void createRevocationStateCB(indy_handle_t, indy_error_t, char*);

typedef void (*cb_proverCloseCredentialsSearchForProofReq)(indy_handle_t, indy_error_t);
extern void proverCloseCredentialsSearchForProofReqCB(indy_handle_t, indy_error_t);

typedef void (*cb_issuerRevokeCredential)(indy_handle_t, indy_error_t, char*);
extern void issuerRevokeCredentialCB(indy_handle_t, indy_error_t, char*);

typedef void (*cb_proverDeleteCredential)(indy_handle_t, indy_error_t);
extern void proverDeleteCredentialCB(indy_handle_t, indy_error_t);

typedef void (*cb_proverGetCredentials)(indy_handle_t, indy_error_t, char*);
extern void proverGetCredentialsCB(indy_handle_t, indy_error_t, char*);

typedef void (*cb_proverGetCredentialsForProofReq)(indy_handle_t, indy_error_t, char*);
extern void proverGetCredentialsForProofReqCB(indy_handle_t, indy_error_t, char*);

typedef void (*cb_proverSearchCredentials)(indy_handle_t, indy_error_t, indy_handle_t, indy_u32_t);
extern void proverSearchCredentialsCB(indy_handle_t, indy_error_t, indy_handle_t, indy_u32_t);

typedef void (*cb_proverFetchCredentials)(indy_handle_t, indy_error_t, char*);
extern void proverFetchCredentialsCB(indy_handle_t, indy_error_t, char*);

typedef void (*cb_toUnqualified)(indy_handle_t, indy_error_t, char*);
extern void toUnqualifiedCB(indy_handle_t, indy_error_t, char*);
*/
import "C"
import (
	"github.com/joyride9999/IndySdkGoBindings/indyUtils"
	"errors"
	"unsafe"
)

//export issuerRevokeCredentialCB
func issuerRevokeCredentialCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, revocRegDeltaJson *C.char) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle),
			indyUtils.IndyResult{Error: nil,
				Results: []interface{}{
					string(C.GoString(revocRegDeltaJson)),
				}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// IssuerRevokeCredential      Revoke a credential identified by a cred_revoc_id (returned by issuer_create_credential).
//
//    The corresponding credential definition and revocation registry must be already
//    created an stored into the wallet.
//
//    This call returns revoc registry delta as json file intended to be shared as REVOC_REG_ENTRY transaction.
//    Note that it is possible to accumulate deltas to reduce ledger load.
func IssuerRevokeCredential(wh int, bh int, revRegId unsafe.Pointer, credRevId unsafe.Pointer) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
	   :param wallet_handle: wallet handle (created by open_wallet).
	   :param blob_storage_reader_handle: pre-configured blob storage reader instance handle that will allow
	   to read revocation tails
	   :param rev_reg_id: id of revocation registry stored in wallet
	   :param cred_revoc_id: local id for revocation info
	   :return: Revocation registry delta json with a revoked credential.
	*/
	res := C.indy_issuer_revoke_credential(commandHandle,
		(C.indy_handle_t)(wh),
		(C.indy_handle_t)(bh),
		(*C.char)(revRegId),
		(*C.char)(credRevId),
		(C.cb_issuerRevokeCredential)(unsafe.Pointer(C.issuerRevokeCredentialCB)))

	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export proverCloseCredentialsSearchForProofReqCB
func proverCloseCredentialsSearchForProofReqCB(commandHandle C.indy_handle_t, indyError C.indy_error_t) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle),
			indyUtils.IndyResult{Error: nil})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// ProverCloseCredentialsSearchForProofReq       close search handle
func ProverCloseCredentialsSearchForProofReq(searchHandle int) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
	   :param search_handle: Search handle (created by prover_search_credentials_for_proof_req)
	   :return: None
	*/
	res := C.indy_prover_close_credentials_search_for_proof_req(commandHandle,
		(C.indy_handle_t)(searchHandle),
		(C.cb_proverCloseCredentialsSearchForProofReq)(unsafe.Pointer(C.proverCloseCredentialsSearchForProofReqCB)))

	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export createRevocationStateCB
func createRevocationStateCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, revocStateJson *C.char) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle),
			indyUtils.IndyResult{Error: nil,
				Results: []interface{}{
					string(C.GoString(revocStateJson)),
				}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// CreateRevocationState        Create revocation state for a credential in the particular time moment.
func CreateRevocationState(blobReaderHandle int, revRegDefJson unsafe.Pointer, revRegDeltaJson unsafe.Pointer, timestamp uint64, credRevId unsafe.Pointer) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
			:param blob_storage_reader_handle: configuration of blob storage reader handle that will allow to read revocation tails
		    :param rev_reg_def_json: revocation registry definition json
		    :param rev_reg_delta_json: revocation registry definition delta json
		    :param timestamp: time represented as a total number of seconds from Unix Epoch
		    :param cred_rev_id: user credential revocation id in revocation registry
		    :return: revocation state json {
		         "rev_reg": <revocation registry>,
		         "witness": <witness>,
		         "timestamp" : integer
		    }
	*/
	res := C.indy_create_revocation_state(commandHandle,
		(C.indy_handle_t)(blobReaderHandle),
		(*C.char)(revRegDefJson),
		(*C.char)(revRegDeltaJson),
		C.ulonglong(timestamp),
		(*C.char)(credRevId),
		(C.cb_createRevocationState)(unsafe.Pointer(C.createRevocationStateCB)))

	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export createAndStoreRevocRegCB
func createAndStoreRevocRegCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, revocRegId *C.char, revocRegDefJson *C.char, revocRegEntryJson *C.char) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle),
			indyUtils.IndyResult{Error: nil,
				Results: []interface{}{
					string(C.GoString(revocRegId)),
					string(C.GoString(revocRegDefJson)),
					string(C.GoString(revocRegEntryJson)),
				}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// CreateAndStoreRevocReg       Create a new revocation registry for the given credential definition as tuple of entities:
//    - Revocation registry definition that encapsulates credentials definition reference, revocation type specific configuration and
//      secrets used for credentials revocation
//    - Revocation registry state that stores the information about revoked entities in a non-disclosing way. The state can be
//      represented as ordered list of revocation registry entries were each entry represents the list of revocation or issuance operations.
//
//    Revocation registry definition entity contains private and public parts. Private part will be stored in the wallet. Public part
//    will be returned as json intended to be shared with all anoncreds workflow actors usually by publishing REVOC_REG_DEF transaction
//    to Indy distributed ledger.
//
//    Revocation registry state is stored on the wallet and also intended to be shared as the ordered list of REVOC_REG_ENTRY transactions.
//    This call initializes the state in the wallet and returns the initial entry.
//
//    Some revocation registry types (for example, 'CL_ACCUM') can require generation of binary blob called tails used to hide information about revoked credentials in public
//    revocation registry and intended to be distributed out of leger (REVOC_REG_DEF transaction will still contain uri and hash of tails).
//    This call requires access to pre-configured blob storage writer instance handle that will allow to write generated tails.
func CreateAndStoreRevocReg(wh int, issuerDid unsafe.Pointer, revocDefType unsafe.Pointer, tag unsafe.Pointer, credDefId unsafe.Pointer, configJson unsafe.Pointer, blobHandle int) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
			:param wallet_handle: wallet handle (created by open_wallet).
		    :param issuer_did: a DID of the issuer signing transaction to the Ledger
		    :param revoc_def_type: revocation registry type (optional, default value depends on credential definition type). Supported types are:
		                - 'CL_ACCUM': Type-3 pairing based accumulator implemented according to the algorithm in this paper:
		                                  https://github.com/hyperledger/ursa/blob/master/libursa/docs/AnonCred.pdf
		                              This type is default for 'CL' credential definition type.
		    :param tag: allows to distinct between revocation registries for the same issuer and credential definition
		    :param cred_def_id: id of stored in ledger credential definition
		    :param config_json: type-specific configuration of revocation registry as json:
		        - 'CL_ACCUM':
		            "issuance_type": (optional) type of issuance. Currently supported:
		                1) ISSUANCE_BY_DEFAULT: all indices are assumed to be issued and initial accumulator is calculated over all indices;
		                   Revocation Registry is updated only during revocation.
		                2) ISSUANCE_ON_DEMAND: nothing is issued initially accumulator is 1 (used by default);
		            "max_cred_num": maximum number of credentials the new registry can process (optional, default 100000)
		        }
		    :param tails_writer_handle: handle of blob storage to store tails

		    NOTE:
		        Recursive creation of folder for Default Tails Writer (correspondent to `tails_writer_handle`)
		        in the system-wide temporary directory may fail in some setup due to permissions: `IO error: Permission denied`.
		        In this case use `TMPDIR` environment variable to define temporary directory specific for an application.

		    :return:
		        revoc_reg_id: identifier of created revocation registry definition
		        revoc_reg_def_json: public part of revocation registry definition
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
		                    "publicKeys": <public_keys> - Registry's public key (opaque type that contains data structures internal to Ursa.
		                                                                         It should not be parsed and are likely to change in future versions).
		                },
		                "ver": string - version of revocation registry definition json.
		            }
		        revoc_reg_entry_json: revocation registry entry that defines initial state of revocation registry
		            {
		                value: {
		                    prevAccum: string - previous accumulator value.
		                    accum: string - current accumulator value.
		                    issued: array<number> - an array of issued indices.
		                    revoked: array<number> an array of revoked indices.
		                },
		                ver: string - version revocation registry entry json
		            }

	*/
	res := C.indy_issuer_create_and_store_revoc_reg(commandHandle,
		(C.indy_handle_t)(wh),
		(*C.char)(issuerDid),
		(*C.char)(revocDefType),
		(*C.char)(tag),
		(*C.char)(credDefId),
		(*C.char)(configJson),
		(C.indy_handle_t)(blobHandle),
		(C.cb_createAndStoreRevocReg)(unsafe.Pointer(C.createAndStoreRevocRegCB)))

	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export proverGetCredentialCB
func proverGetCredentialCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, credential *C.char) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle),
			indyUtils.IndyResult{Error: nil,
				Results: []interface{}{
					string(C.GoString(credential)),
				}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// ProverGetCredential   Gets human readable credential by the given id.
func ProverGetCredential(wh int, credentialId unsafe.Pointer) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
			:param wallet_handle: wallet handle (created by open_wallet).
		    :param cred_id: Identifier by which requested credential is stored in the wallet
		    :return:  credential json
		     {
		         "referent": string, - id of credential in the wallet
		         "attrs": {"key1":"raw_value1", "key2":"raw_value2"}, - credential attributes
		         "schema_id": string, - identifier of schema
		         "cred_def_id": string, - identifier of credential definition
		         "rev_reg_id": Optional<string>, - identifier of revocation registry definition
		         "cred_rev_id": Optional<string> - identifier of credential in the revocation registry definition
		     }

	*/
	res := C.indy_prover_get_credential(commandHandle,
		(C.indy_handle_t)(wh),
		(*C.char)(credentialId),
		(C.cb_proverGetCredential)(unsafe.Pointer(C.proverGetCredentialCB)))

	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export verifierVerifyProofCB
func verifierVerifyProofCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, valid C.indy_bool_t) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle),
			indyUtils.IndyResult{Error: nil,
				Results: []interface{}{
					(bool)(valid),
				}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// VerifierVerifyProof   Verifies a proof (of multiple credential).
//    All required schemas, public keys and revocation registries must be provided.
//
//    IMPORTANT: You must use *_id's (`schema_id`, `cred_def_id`, `rev_reg_id`) listed in `proof[identifiers]`
//        as the keys for corresponding `schemas_json`, `credential_defs_json`, `rev_reg_defs_json`, `rev_regs_json` objects.
func VerifierVerifyProof(proofRequestJson, proofJson, schemasJson, credDefsJson, revRegDefsJson, revRegsJson unsafe.Pointer) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
			  :param proof_request_json:
		        {
		            "name": string,
		            "version": string,
		            "nonce": string, - a decimal number represented as a string (use `generate_nonce` function to generate 80-bit number)
		            "requested_attributes": { // set of requested attributes
		                 "<attr_referent>": <attr_info>, // see below
		                 ...,
		            },
		            "requested_predicates": { // set of requested predicates
		                 "<predicate_referent>": <predicate_info>, // see below
		                 ...,
		             },
		            "non_revoked": Optional<<non_revoc_interval>>, // see below,
		                           // If specified prover must proof non-revocation
		                           // for date in this interval for each attribute
		                           // (can be overridden on attribute level)
		            "ver": Optional<str>  - proof request version:
		                - omit to use unqualified identifiers for restrictions
		                - "1.0" to use unqualified identifiers for restrictions
		                - "2.0" to use fully qualified identifiers for restrictions
		        }
		    :param proof_json: created for request proof json
		        {
		            "requested_proof": {
		                "revealed_attrs": {
		                    "requested_attr1_id": {sub_proof_index: number, raw: string, encoded: string}, // NOTE: check that `encoded` value match to `raw` value on application level
		                    "requested_attr4_id": {sub_proof_index: number: string, encoded: string}, // NOTE: check that `encoded` value match to `raw` value on application level
		                },
		                "revealed_attr_groups": {
		                    "requested_attr5_id": {
		                        "sub_proof_index": number,
		                        "values": {
		                            "attribute_name": {
		                                "raw": string,
		                                "encoded": string
		                            }
		                        }, // NOTE: check that `encoded` value match to `raw` value on application level
		                    }
		                },
		                "unrevealed_attrs": {
		                    "requested_attr3_id": {sub_proof_index: number}
		                },
		                "self_attested_attrs": {
		                    "requested_attr2_id": self_attested_value,
		                },
		                "requested_predicates": {
		                    "requested_predicate_1_referent": {sub_proof_index: int},
		                    "requested_predicate_2_referent": {sub_proof_index: int},
		                }
		            }
		            "proof": {
		                "proofs": [ <credential_proof>, <credential_proof>, <credential_proof> ],
		                "aggregated_proof": <aggregated_proof>
		            }
		            "identifiers": [{schema_id, cred_def_id, Optional<rev_reg_id>, Optional<timestamp>}]
		        }
		    :param schemas_json: all schema jsons participating in the proof
		         {
		             <schema1_id>: <schema1>,
		             <schema2_id>: <schema2>,
		             <schema3_id>: <schema3>,
		         }
		    :param credential_defs_json: all credential definitions json participating in the proof
		         {
		             "cred_def1_id": <credential_def1>,
		             "cred_def2_id": <credential_def2>,
		             "cred_def3_id": <credential_def3>,
		         }
		    :param rev_reg_defs_json: all revocation registry definitions json participating in the proof
		         {
		             "rev_reg_def1_id": <rev_reg_def1>,
		             "rev_reg_def2_id": <rev_reg_def2>,
		             "rev_reg_def3_id": <rev_reg_def3>,
		         }
		    :param rev_regs_json: all revocation registries json participating in the proof
		         {
		             "rev_reg_def1_id": {
		                 "timestamp1": <rev_reg1>,
		                 "timestamp2": <rev_reg2>,
		             },
		             "rev_reg_def2_id": {
		                 "timestamp3": <rev_reg3>
		             },
		             "rev_reg_def3_id": {
		                 "timestamp4": <rev_reg4>
		             },
		         }
		    :return: valid: true - if signature is valid, false - otherwise
	*/
	res := C.indy_verifier_verify_proof(commandHandle,
		(*C.char)(proofRequestJson),
		(*C.char)(proofJson),
		(*C.char)(schemasJson),
		(*C.char)(credDefsJson),
		(*C.char)(revRegDefsJson),
		(*C.char)(revRegsJson),
		(C.cb_verifierVerifyProof)(unsafe.Pointer(C.verifierVerifyProofCB)))

	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export proverCreateProofCB
func proverCreateProofCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, proofJson *C.char) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle),
			indyUtils.IndyResult{Error: nil,
				Results: []interface{}{
					string(C.GoString(proofJson)),
				}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// ProverCreateProof   Creates a proof according to the given proof request
//    Either a corresponding credential with optionally revealed attributes or self-attested attribute must be provided
//    for each requested attribute (see indy_prover_get_credentials_for_pool_req).
//    A proof request may request multiple credentials from different schemas and different issuers.
//    All required schemas, public keys and revocation registries must be provided.
//    The proof request also contains nonce.
//    The proof contains either proof or self-attested attribute value for each requested attribute.
func ProverCreateProof(wh int, proofRequestJson, requestedCredentialsJson, masterSecretId, schemasForAttrsJson, credentialDefsForAttrsJson, revStatesJson unsafe.Pointer) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
			:param wallet_handle: wallet handle (created by open_wallet).
		    :param proof_req_json: proof request json
		        {
		            "name": string,
		            "version": string,
		            "nonce": string, - a decimal number represented as a string (use `generate_nonce` function to generate 80-bit number)
		            "requested_attributes": { // set of requested attributes
		                 "<attr_referent>": <attr_info>, // see below
		                 ...,
		            },
		            "requested_predicates": { // set of requested predicates
		                 "<predicate_referent>": <predicate_info>, // see below
		                 ...,
		             },
		            "non_revoked": Optional<<non_revoc_interval>>, // see below,
		                           // If specified prover must proof non-revocation
		                           // for date in this interval for each attribute
		                           // (applies to every attribute and predicate but can be overridden on attribute level)
		                           // (can be overridden on attribute level)
		            "ver": Optional<str>  - proof request version:
		                - omit to use unqualified identifiers for restrictions
		                - "1.0" to use unqualified identifiers for restrictions
		                - "2.0" to use fully qualified identifiers for restrictions
		        }
		    :param requested_credentials_json: either a credential or self-attested attribute for each requested attribute
		        {
		            "self_attested_attributes": {
		                "self_attested_attribute_referent": string
		            },
		            "requested_attributes": {
		                "requested_attribute_referent_1": {"cred_id": string, "timestamp": Optional<number>, revealed: <bool> }},
		                "requested_attribute_referent_2": {"cred_id": string, "timestamp": Optional<number>, revealed: <bool> }}
		            },
		            "requested_predicates": {
		                "requested_predicates_referent_1": {"cred_id": string, "timestamp": Optional<number> }},
		            }
		        }
		    :param master_secret_name: the id of the master secret stored in the wallet
		    :param schemas_json: all schemas json participating in the proof request
		          {
		              <schema1_id>: <schema1>,
		              <schema2_id>: <schema2>,
		              <schema3_id>: <schema3>,
		          }
		    :param credential_defs_json: all credential definitions json participating in the proof request
		          {
		              "cred_def1_id": <credential_def1>,
		              "cred_def2_id": <credential_def2>,
		              "cred_def3_id": <credential_def3>,
		          }
		    :param rev_states_json: all revocation states json participating in the proof request
		          {
		              "rev_reg_def1_id or credential_1_id": {
		                  "timestamp1": <rev_state1>,
		                  "timestamp2": <rev_state2>,
		              },
		              "rev_reg_def2_id or credential_2_id": {
		                  "timestamp3": <rev_state3>
		              },
		              "rev_reg_def3_id or credential_3_id": {
		                  "timestamp4": <rev_state4>
		              },
		          } - Note: use credential_id instead rev_reg_id in case proving several credentials from the same revocation registry.

		    where
		        attr_referent: Proof-request local identifier of requested attribute
		        attr_info: Describes requested attribute
		            {
		                "name": Optional<string>, // attribute name, (case insensitive and ignore spaces)
		                "names": Optional<[string, string]>, // attribute names, (case insensitive and ignore spaces)
		                                                     // NOTE: should either be "name" or "names", not both and not none of them.
		                                                     // Use "names" to specify several attributes that have to match a single credential.
		                "restrictions": Optional<filter_json>, // see below
		                "non_revoked": Optional<<non_revoc_interval>>, // see below,
		                               // If specified prover must proof non-revocation
		                               // for date in this interval this attribute
		                           // (overrides proof level interval)
		            }
		        predicate_referent: Proof-request local identifier of requested attribute predicate
		        predicate_info: Describes requested attribute predicate
		            {
		                "name": attribute name, (case insensitive and ignore spaces)
		                "p_type": predicate type (">=", ">", "<=", "<")
		                "p_value": predicate value
		                "restrictions": Optional<wql query>, // see below
		                "non_revoked": Optional<<non_revoc_interval>>, // see below,
		                               // If specified prover must proof non-revocation
		                               // for date in this interval this attribute
		                               // (overrides proof level interval)
		            }
		        non_revoc_interval: Defines non-revocation interval
		            {
		                "from": Optional<int>, // timestamp of interval beginning
		                "to": Optional<int>, // timestamp of interval ending
		            }
		        where wql query: indy-sdk/docs/design/011-wallet-query-language/README.md
		            The list of allowed fields:
		                "schema_id": <credential schema id>,
		                "schema_issuer_did": <credential schema issuer did>,
		                "schema_name": <credential schema name>,
		                "schema_version": <credential schema version>,
		                "issuer_did": <credential issuer did>,
		                "cred_def_id": <credential definition id>,
		                "rev_reg_id": <credential revocation registry id>, // "None" as string if not present

		    :return: Proof json
		      For each requested attribute either a proof (with optionally revealed attribute value) or
		      self-attested attribute value is provided.
		      Each proof is associated with a credential and corresponding schema_id, cred_def_id, rev_reg_id and timestamp.
		      There is also aggregated proof part common for all credential proofs.
		            {
		                "requested_proof": {
		                    "revealed_attrs": {
		                        "requested_attr1_id": {sub_proof_index: number, raw: string, encoded: string},
		                        "requested_attr4_id": {sub_proof_index: number: string, encoded: string},
		                    },
		                    "revealed_attr_groups": {
		                        "requested_attr5_id": {
		                            "sub_proof_index": number,
		                            "values": {
		                                "attribute_name": {
		                                    "raw": string,
		                                    "encoded": string
		                                }
		                            },
		                        }
		                    },
		                    "unrevealed_attrs": {
		                        "requested_attr3_id": {sub_proof_index: number}
		                    },
		                    "self_attested_attrs": {
		                        "requested_attr2_id": self_attested_value,
		                    },
		                    "predicates": {
		                        "requested_predicate_1_referent": {sub_proof_index: int},
		                        "requested_predicate_2_referent": {sub_proof_index: int},
		                    }
		                }
		                "proof": {
		                    "proofs": [ <credential_proof>, <credential_proof>, <credential_proof> ],
		                    "aggregated_proof": <aggregated_proof>
		                } (opaque type that contains data structures internal to Ursa.
		                  It should not be parsed and are likely to change in future versions).
		                "identifiers": [{schema_id, cred_def_id, Optional<rev_reg_id>, Optional<timestamp>}]
		            }
	*/

	// Call C.indy_prover_create_proof
	res := C.indy_prover_create_proof(commandHandle,
		(C.indy_handle_t)(wh),
		(*C.char)(proofRequestJson),
		(*C.char)(requestedCredentialsJson),
		(*C.char)(masterSecretId),
		(*C.char)(schemasForAttrsJson),
		(*C.char)(credentialDefsForAttrsJson),
		(*C.char)(revStatesJson),
		(C.cb_proverCreateProof)(unsafe.Pointer(C.proverCreateProofCB)))

	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export proverFetchCredentialsForProofReqCB
func proverFetchCredentialsForProofReqCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, credentialsJson *C.char) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle),
			indyUtils.IndyResult{Error: nil,
				Results: []interface{}{
					string(C.GoString(credentialsJson)),
				}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// ProverFetchCredentialsForProofReq   Fetch next records for the requested item using proof request search handle (created by prover_search_credentials_for_proof_req).
func ProverFetchCredentialsForProofReq(searchHandle int, itemReferent unsafe.Pointer, count int) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
			:param search_handle: Search handle (created by prover_search_credentials_for_proof_req)
		    :param item_referent: Referent of attribute/predicate in the proof request
		    :param count: Count of records to fetch
		    :return: credentials_json: List of credentials for the given proof request.
		        [{
		            cred_info: <credential_info>,
		            interval: Optional<non_revoc_interval>
		        }]
		    where credential_info is
		        {
		            "referent": string, - id of credential in the wallet
		            "attrs": {"key1":"raw_value1", "key2":"raw_value2"}, - credential attributes
		            "schema_id": string, - identifier of schema
		            "cred_def_id": string, - identifier of credential definition
		            "rev_reg_id": Optional<string>, - identifier of revocation registry definition
		            "cred_rev_id": Optional<string> - identifier of credential in the revocation registry definition
		        }
		    NOTE: The list of length less than the requested count means that search iterator correspondent to the requested <item_referent> is completed.
	*/
	res := C.indy_prover_fetch_credentials_for_proof_req(commandHandle,
		(C.indy_handle_t)(searchHandle),
		(*C.char)(itemReferent),
		(C.indy_u32_t)(count),
		(C.cb_proverFetchCredentialsForProofReq)(unsafe.Pointer(C.proverFetchCredentialsForProofReqCB)))

	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export proverSearchForCredentialsForProofReqCB
func proverSearchForCredentialsForProofReqCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, searchHandle C.indy_handle_t) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle),
			indyUtils.IndyResult{Error: nil,
				Results: []interface{}{
					(int)(searchHandle),
				}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// ProverSearchForCredentialsForProofReq   Search for credentials matching the given proof request.
//
//    Instead of immediately returning of fetched credentials this call returns search_handle that can be used later
//    to fetch records by small batches (with prover_fetch_credentials_for_proof_req).
func ProverSearchForCredentialsForProofReq(wh int, proofRequestJson, extraQueryJson unsafe.Pointer) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
	  param wallet_handle: wallet handle (created by open_wallet).
	    :param proof_request_json: proof request json
	        {
	            "name": string,
	            "version": string,
	            "nonce": string, - a decimal number represented as a string (use `indy_generate_nonce` function to generate 80-bit number)
	            "requested_attributes": { // set of requested attributes
	                 "<attr_referent>": <attr_info>, // see below
	                 ...,
	            },
	            "requested_predicates": { // set of requested predicates
	                 "<predicate_referent>": <predicate_info>, // see below
	                 ...,
	             },
	            "non_revoked": Optional<<non_revoc_interval>>, // see below,
	                           // If specified prover must proof non-revocation
	                           // for date in this interval for each attribute
	                           // (applies to every attribute and predicate but can be overridden on attribute level)
	                           // (can be overridden on attribute level)
	            "ver": Optional<str>  - proof request version:
	                - omit to use unqualified identifiers for restrictions
	                - "1.0" to use unqualified identifiers for restrictions
	                - "2.0" to use fully qualified identifiers for restrictions
	        }
	    :param extra_query_json:(Optional) List of extra queries that will be applied to correspondent attribute/predicate:
	        {
	            "<attr_referent>": <wql query>,
	            "<predicate_referent>": <wql query>,
	        }


	    where
	    attr_info: Describes requested attribute
	        {
	            "name": Optional<string>, // attribute name, (case insensitive and ignore spaces)
	            "names": Optional<[string, string]>, // attribute names, (case insensitive and ignore spaces)
	                                                 // NOTE: should either be "name" or "names", not both and not none of them.
	                                                 // Use "names" to specify several attributes that have to match a single credential.
	            "restrictions": Optional<filter_json>, // see below
	            "non_revoked": Optional<<non_revoc_interval>>, // see below,
	                           // If specified prover must proof non-revocation
	                           // for date in this interval this attribute
	                           // (overrides proof level interval)
	        }
	    predicate_referent: Proof-request local identifier of requested attribute predicate
	    predicate_info: Describes requested attribute predicate
	        {
	            "name": attribute name, (case insensitive and ignore spaces)
	            "p_type": predicate type (">=", ">", "<=", "<")
	            "p_value": predicate value
	            "restrictions": Optional<wql query>, // see below
	            "non_revoked": Optional<<non_revoc_interval>>, // see below,
	                           // If specified prover must proof non-revocation
	                           // for date in this interval this attribute
	                           // (overrides proof level interval)
	        }
	    non_revoc_interval: Defines non-revocation interval
	        {
	            "from": Optional<int>, // timestamp of interval beginning
	            "to": Optional<int>, // timestamp of interval ending
	        }
	    extra_query_json:(Optional) List of extra queries that will be applied to correspondent attribute/predicate:
	        {
	            "<attr_referent>": <wql query>,
	            "<predicate_referent>": <wql query>,
	        }
	    where wql query: indy-sdk/docs/design/011-wallet-query-language/README.md
	        The list of allowed fields:
	            "schema_id": <credential schema id>,
	            "schema_issuer_did": <credential schema issuer did>,
	            "schema_name": <credential schema name>,
	            "schema_version": <credential schema version>,
	            "issuer_did": <credential issuer did>,
	            "cred_def_id": <credential definition id>,
	            "rev_reg_id": <credential revocation registry id>, // "None" as string if not present

	    :return: search_handle: Search handle that can be used later to fetch records by small batches (with prover_fetch_credentials_for_proof_req)
	*/
	res := C.indy_prover_search_credentials_for_proof_req(commandHandle,
		(C.indy_handle_t)(wh),
		(*C.char)(proofRequestJson),
		(*C.char)(extraQueryJson),
		(C.cb_proverSearchForCredentialsForProofReq)(unsafe.Pointer(C.proverSearchForCredentialsForProofReqCB)))

	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export generateNonceCB
func generateNonceCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, nonce *C.char) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle),
			indyUtils.IndyResult{Error: nil,
				Results: []interface{}{
					string(C.GoString(nonce)),
				}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// GenerateNonce Generates 80-bit numbers that can be used as a nonce for proof request
// The name must be unique.
func GenerateNonce() chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
	  :return: nonce: generated number as a string
	*/
	res := C.indy_generate_nonce(commandHandle,
		(C.cb_generateNonce)(unsafe.Pointer(C.generateNonceCB)))

	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export proverCreateMasterSecretCB
func proverCreateMasterSecretCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, idMasterSecret *C.char) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle),
			indyUtils.IndyResult{Error: nil,
				Results: []interface{}{
					string(C.GoString(idMasterSecret)),
				}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// ProverCreateMasterSecret Creates a master secret with a given name and stores it in the wallet.
// The name must be unique.
func ProverCreateMasterSecret(wh int, secretName unsafe.Pointer) chan indyUtils.IndyResult {
	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
			 :param wallet_handle: wallet handle (created by open_wallet).
		    :param master_secret_name: (optional, if not present random one will be generated) new master id
	*/
	res := C.indy_prover_create_master_secret(commandHandle,
		(C.indy_handle_t)(wh),
		(*C.char)(secretName),
		(C.cb_proverCreateMasterSecret)(unsafe.Pointer(C.proverCreateMasterSecretCB)))

	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export issuerCreateCredentialOfferCB
func issuerCreateCredentialOfferCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, credentialOfferJson *C.char) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle),
			indyUtils.IndyResult{Error: nil,
				Results: []interface{}{
					string(C.GoString(credentialOfferJson)),
				}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// IssuerCreateCredentialOffer Create credential offer that will be used by Prover for
//    credential request creation. Offer includes nonce and key correctness proof
//    for authentication between protocol steps and integrity checking.
func IssuerCreateCredentialOffer(wh int, credDefinitionId unsafe.Pointer) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
	   :param wallet_handle: wallet handle (created by open_wallet).
	   :param cred_def_id: id of credential definition stored in the wallet
	   :return:credential offer json:
	    {
	        "schema_id": string, - identifier of schema
	        "cred_def_id": string, - identifier of credential definition
	        // Fields below can depend on Cred Def type
	        "nonce": string,
	        "key_correctness_proof" : key correctness proof for credential definition correspondent to cred_def_id
	                                  (opaque type that contains data structures internal to Ursa.
	                                  It should not be parsed and are likely to change in future versions).
	    }
	*/
	res := C.indy_issuer_create_credential_offer(commandHandle,
		(C.indy_handle_t)(wh),
		(*C.char)(credDefinitionId),
		(C.cb_issuerCreateCredentialOffer)(unsafe.Pointer(C.issuerCreateCredentialOfferCB)))

	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export proverCreateCredentialRequestCB
func proverCreateCredentialRequestCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, credentialRequestJSON *C.char, credentialRequestMetadataJSON *C.char) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle),
			indyUtils.IndyResult{Error: nil,
				Results: []interface{}{
					string(C.GoString(credentialRequestJSON)),
					string(C.GoString(credentialRequestMetadataJSON)),
				}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// ProverCreateCredentialRequest  Creates a credential request for the given credential offer.
//
//    The method creates a blinded master secret for a master secret identified by a provided name.
//    The master secret identified by the name must be already stored in the secure wallet (see prover_create_master_secret)
//    The blinded master secret is a part of the credential request.
func ProverCreateCredentialRequest(wh int, proverDID, credentialOfferJSON, credentialDefinitionJSON, masterSecretId unsafe.Pointer) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
	 :param wallet_handle: wallet handle (created by open_wallet).
	    :param prover_did: a DID of the prover
	    :param cred_offer_json: credential offer as a json containing information about the issuer and a credential
	        {
	            "schema_id": string, - identifier of schema
	            "cred_def_id": string, - identifier of credential definition
	             ...
	            Other fields that contains data structures internal to Ursa.
	            These fields should not be parsed and are likely to change in future versions.
	        }
	    :param cred_def_json: credential definition json related to <cred_def_id> in <cred_offer_json>
	    :param master_secret_id: the id of the master secret stored in the wallet
	    :return:
	     cred_req_json: Credential request json for creation of credential by Issuer
	     {
	      "prover_did" : string,
	      "cred_def_id" : string,
	         // Fields below can depend on Cred Def type
	      "blinded_ms" : <blinded_master_secret>,
	                     (opaque type that contains data structures internal to Ursa.
	                      It should not be parsed and are likely to change in future versions).
	      "blinded_ms_correctness_proof" : <blinded_ms_correctness_proof>,
	                     (opaque type that contains data structures internal to Ursa.
	                      It should not be parsed and are likely to change in future versions).
	      "nonce": string
	    }
	     cred_req_metadata_json:  Credential request metadata json for further processing of received form Issuer credential.
	                              Credential request metadata contains data structures internal to Ursa.
	                              Credential request metadata mustn't be shared with Issuer.
	*/
	res := C.indy_prover_create_credential_req(commandHandle,
		(C.indy_handle_t)(wh),
		(*C.char)(proverDID),
		(*C.char)(credentialOfferJSON),
		(*C.char)(credentialDefinitionJSON),
		(*C.char)(masterSecretId),
		(C.cb_proverCreateCredentialRequest)(unsafe.Pointer(C.proverCreateCredentialRequestCB)))

	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export issuerCreateSchemaCB
func issuerCreateSchemaCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, id *C.char, schemaJson *C.char) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle),
			indyUtils.IndyResult{Error: nil,
				Results: []interface{}{
					string(C.GoString(id)),
					string(C.GoString(schemaJson)),
				}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// IssuerCreateSchema creates a credential schema
func IssuerCreateSchema(submitterDid unsafe.Pointer, name unsafe.Pointer, version unsafe.Pointer, attrs unsafe.Pointer) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
	   :param issuer_did: DID of schema issuer
	   :param name: a name the schema
	   :param version: a version of the schema
	   :param attrs: a list of schema attributes descriptions (the number of attributes should be less or equal than 125)
	                     `["attr1", "attr2"]`
	   :return:
	       schema_id: identifier of created schema
	       schema_json: schema as json
	       {
	           id: identifier of schema
	           attrNames: array of attribute name strings
	           name: schema's name string
	           version: schema's version string,
	           ver: version of the Schema json
	       }
	*/
	res := C.indy_issuer_create_schema(commandHandle,
		(*C.char)(submitterDid),
		(*C.char)(name),
		(*C.char)(version),
		(*C.char)(attrs),
		(C.cb_issuerCreateSchema)(unsafe.Pointer(C.issuerCreateSchemaCB)))

	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export issuerCreateAndStoreCredentialDefCB
func issuerCreateAndStoreCredentialDefCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, credDefId *C.char, credDefJson *C.char) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle),
			indyUtils.IndyResult{Error: nil,
				Results: []interface{}{
					string(C.GoString(credDefId)),
					string(C.GoString(credDefJson)),
				}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// IssuerCreateAndStoreCredentialDef Create credential definition entity that encapsulates credentials issuer DID, credential schema, secrets used for
//    signing credentials and secrets used for credentials revocation.
//    Credential definition entity contains private and public parts. Private part will be stored in the wallet.
//    Public part will be returned as json intended to be shared with all anoncreds workflow actors usually by
//    publishing CRED_DEF transaction to Indy distributed ledger.
//
//    It is IMPORTANT for current version GET Schema from Ledger with correct seq_no to save compatibility with Ledger.
//
//    Note: Use combination of `issuer_rotate_credential_def_start` and `issuer_rotate_credential_def_apply` functions
//    to generate new keys for an existing credential definition.
func IssuerCreateAndStoreCredentialDef(wh int, issuerDid unsafe.Pointer, schemaJson unsafe.Pointer, tag unsafe.Pointer,
	signatureType unsafe.Pointer, configJson unsafe.Pointer) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
	   :param wallet_handle: wallet handle (created by open_wallet).
	    :param issuer_did: a DID of the issuer signing cred_def transaction to the Ledger
	    :param schema_json: credential schema as a json
	        {
	            id: identifier of schema
	            attrNames: array of attribute name strings
	            name: schema's name string
	            version: schema's version string,
	            seqNo: (Optional) schema's sequence number on the ledger,
	            ver: version of the Schema json
	        }
	    :param tag: allows to distinct between credential definitions for the same issuer and schema
	    :param signature_type: credential definition type (optional, 'CL' by default) that defines credentials signature and revocation math.
	    Supported types are:
	        - 'CL': Camenisch-Lysyanskaya credential signature type that is implemented according to the algorithm in this paper:
	                    https://github.com/hyperledger/ursa/blob/master/libursa/docs/AnonCred.pdf
	                And is documented in this HIPE:
	                    https://github.com/hyperledger/indy-hipe/blob/c761c583b1e01c1e9d3ceda2b03b35336fdc8cc1/text/anoncreds-protocol/README.md
	    :param  config_json: (optional) type-specific configuration of credential definition as json:
	        - 'CL':
	            {
	                "support_revocation" - bool (optional, default false) whether to request non-revocation credential
	            }
	    :return:
	        cred_def_id: identifier of created credential definition
	        cred_def_json: public part of created credential definition
	            {
	                id: string - identifier of credential definition
	                schemaId: string - identifier of stored in ledger schema
	                type: string - type of the credential definition. CL is the only supported type now.
	                tag: string - allows to distinct between credential definitions for the same issuer and schema
	                value: Dictionary with Credential Definition's data is depended on the signature type: {
	                    primary: primary credential public key,
	                    Optional<revocation>: revocation credential public key
	                },
	                ver: Version of the CredDef json
	            }
	*/
	res := C.indy_issuer_create_and_store_credential_def(commandHandle,
		(C.indy_handle_t)(wh),
		(*C.char)(issuerDid),
		(*C.char)(schemaJson),
		(*C.char)(tag),
		(*C.char)(signatureType),
		(*C.char)(configJson),
		(C.cb_issuerCreateAndStoreCredentialDef)(unsafe.Pointer(C.issuerCreateAndStoreCredentialDefCB)))

	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export issuerRotateCredentialDefStartCB
func issuerRotateCredentialDefStartCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, credDefJson *C.char) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle),
			indyUtils.IndyResult{Error: nil,
				Results: []interface{}{
					string(C.GoString(credDefJson))}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

func IssuerRotateCredentialDefStart(walletHandle int, credDefID unsafe.Pointer, configJson unsafe.Pointer) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
			Generate temporary credential definitional keys for an existing one (owned by the caller of the library).

		    Use `issuer_rotate_credential_def_apply` function to set generated temporary keys as the main.
		    WARNING: Rotating the credential definitional keys will result in making all credentials issued under the previous keys unverifiable.
		    :param wallet_handle: wallet handle (created by open_wallet).
		    :param cred_def_id: an identifier of created credential definition stored in the wallet
		    :param  config_json: (optional) type-specific configuration of credential definition as json:
		        - 'CL':
		            {
		                "support_revocation" - bool (optional, default false) whether to request non-revocation credential
		            }
		    :return:
		        cred_def_json: public part of temporary created credential definition
	*/

	// Call C.indy_issuer_rotate_credential_def_start
	res := C.indy_issuer_rotate_credential_def_start(commandHandle,
		C.indy_handle_t(walletHandle),
		(*C.char)(credDefID),
		(*C.char)(configJson),
		(C.cb_issuerRotateCredentialDefStart)(unsafe.Pointer(C.issuerRotateCredentialDefStartCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export issuerRotateCredentialDefApplyCB
func issuerRotateCredentialDefApplyCB(commandHandle C.indy_handle_t, indyError C.indy_error_t) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle),
			indyUtils.IndyResult{Error: nil})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

func IssuerRotateCredentialDefApply(walletHandle int, credDefID unsafe.Pointer) chan indyUtils.IndyResult {

	// Prepare the call parameters.
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
			Apply temporary keys as main for an existing Credential Definition (owned by the caller of the library).

		    WARNING: Rotating the credential definitional keys will result in making all credentials issued under the previous keys unverifiable.
		    :param wallet_handle: wallet handle (created by open_wallet).
		    :param cred_def_id: an identifier of created credential definition stored in the wallet
	*/

	// Call C.indy_issuer_rotate_credential_def_start
	res := C.indy_issuer_rotate_credential_def_apply(commandHandle,
		C.indy_handle_t(walletHandle),
		(*C.char)(credDefID),
		(C.cb_issuerRotateCredentialDefApply)(unsafe.Pointer(C.issuerRotateCredentialDefApplyCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export issuerCreateCredentialCB
func issuerCreateCredentialCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, credentialJson *C.char, credRevocId *C.char, revocRegDeltaJson *C.char) {
	if indyError == 0 {

		revocId := ""
		if credRevocId != nil {
			revocId = string(C.GoString(credRevocId))
		}

		deltaJson := ""
		if revocRegDeltaJson != nil {
			deltaJson = string(C.GoString(revocRegDeltaJson))
		}
		indyUtils.RemoveFuture((int)(commandHandle),
			indyUtils.IndyResult{Error: nil,
				Results: []interface{}{
					string(C.GoString(credentialJson)),
					revocId,
					deltaJson,
				}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// IssuerCreateCredential creates a credential
func IssuerCreateCredential(wh int, credOfferJson, credRequestJson, credentialTranscript, revocRegistryId unsafe.Pointer, blobHandle int) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
	  Check Cred Request for the given Cred Offer and issue Credential for the given Cred Request.

	    Cred Request must match Cred Offer. The credential definition and revocation registry definition
	    referenced in Cred Offer and Cred Request must be already created and stored into the wallet.

	    Information for this credential revocation will be store in the wallet as part of revocation registry under
	    generated cred_revoc_id local for this wallet.

	    This call returns revoc registry delta as json file intended to be shared as REVOC_REG_ENTRY transaction.
	    Note that it is possible to accumulate deltas to reduce ledger load.

	    :param wallet_handle: wallet handle (created by open_wallet).
	    :param cred_offer_json: a cred offer created by issuer_create_credential_offer
	    :param cred_req_json: a credential request created by prover_create_credential_req
	    :param cred_values_json: a credential containing attribute values for each of requested attribute names.
	     Example:
	     {
	      "attr1" : {"raw": "value1", "encoded": "value1_as_int" },
	      "attr2" : {"raw": "value1", "encoded": "value1_as_int" }
	     }
	     If you want to use empty value for some credential field, you should set "raw" to "" and "encoded" should not be empty
	    :param rev_reg_id: (Optional) id of revocation registry definition stored in the wallet
	    :param blob_storage_reader_handle: pre-configured blob storage reader instance handle that
	    will allow to read revocation tails
	    :return:
	     cred_json: Credential json containing signed credential values
	     {
	         "schema_id": string,
	         "cred_def_id": string,
	         "rev_reg_def_id", Optional<string>,
	         "values": <see cred_values_json above>,
	         // Fields below can depend on Cred Def type
	         "signature": <credential signature>,
	                       (opaque type that contains data structures internal to Ursa.
	                        It should not be parsed and are likely to change in future versions).
	         "signature_correctness_proof": credential signature correctness proof
	                                         (opaque type that contains data structures internal to Ursa.
	                                          It should not be parsed and are likely to change in future versions).
	         "rev_reg" - (Optional) revocation registry accumulator value on the issuing moment.
	                     (opaque type that contains data structures internal to Ursa.
	                      It should not be parsed and are likely to change in future versions).
	         "witness" - (Optional) revocation related data
	                     (opaque type that contains data structures internal to Ursa.
	                      It should not be parsed and are likely to change in future versions).
	     }
	     cred_revoc_id: local id for revocation info (Can be used for revocation of this cred)
	     revoc_reg_delta_json: Revocation registry delta json with a newly issued credential
	*/

	res := C.indy_issuer_create_credential(commandHandle,
		(C.indy_handle_t)(wh),
		(*C.char)(credOfferJson),
		(*C.char)(credRequestJson),
		(*C.char)(credentialTranscript),
		(*C.char)(revocRegistryId),
		(C.indy_handle_t)(blobHandle),
		(C.cb_issuerCreateCredential)(unsafe.Pointer(C.issuerCreateCredentialCB)))

	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export proverStoreCredentialCB
func proverStoreCredentialCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, credentialId *C.char) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle),
			indyUtils.IndyResult{Error: nil,
				Results: []interface{}{
					string(C.GoString(credentialId)),
				}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// ProverStoreCredential stores credential. Check credential provided by Issuer for the given credential request,
//    updates the credential by a master secret and stores in a secure wallet.
//
//    To support efficient search the following tags will be created for stored credential:
//        {
//            "schema_id": <credential schema id>,
//            "schema_issuer_did": <credential schema issuer did>,
//            "schema_name": <credential schema name>,
//            "schema_version": <credential schema version>,
//            "issuer_did": <credential issuer did>,
//            "cred_def_id": <credential definition id>,
//            "rev_reg_id": <credential revocation registry id>, # "None" as string if not present
//            // for every attribute in <credential values> that credential attribute tagging policy marks taggable
//            "attr::<attribute name>::marker": "1",
//            "attr::<attribute name>::value": <attribute raw value>,
//        }
func ProverStoreCredential(wh int, credentialId, credRequestMetadataJson, credJson, credDefJson, revocRegDefJson unsafe.Pointer) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
	 :param wallet_handle: wallet handle (created by open_wallet).
	    :param cred_id: (optional, default is a random one) identifier by which credential will be stored in the wallet
	    :param cred_req_metadata_json: a credential request metadata created by prover_create_credential_req
	    :param cred_json: credential json received from issuer
	    :param cred_def_json: credential definition json related to <cred_def_id> in <cred_json>
	    :param rev_reg_def_json: revocation registry definition json related to <rev_reg_def_id> in <cred_json>
	    :return: cred_id: identifier by which credential is stored in the wallet
	*/

	res := C.indy_prover_store_credential(commandHandle,
		(C.indy_handle_t)(wh),
		(*C.char)(credentialId),
		(*C.char)(credRequestMetadataJson),
		(*C.char)(credJson),
		(*C.char)(credDefJson),
		(*C.char)(revocRegDefJson),
		(C.cb_proverStoreCredential)(unsafe.Pointer(C.proverStoreCredentialCB)))

	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export proverDeleteCredentialCB
func proverDeleteCredentialCB(commandHandle C.indy_handle_t, indyError C.indy_error_t) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

func ProverDeleteCredential(walletHandle int, credentialID unsafe.Pointer) chan indyUtils.IndyResult {

	// Prepare the call parameters
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
			Delete identified credential from wallet.

		    :param wallet_handle: wallet handle (created by open_wallet).
		    :param cred_id: identifier by which wallet stores credential to delete
	*/

	// Call C.indy_prover_delete_credential
	res := C.indy_prover_delete_credential(commandHandle,
		C.indy_handle_t(walletHandle),
		(*C.char)(credentialID),
		(C.cb_proverDeleteCredential)(unsafe.Pointer(C.proverDeleteCredentialCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}
	return future
}

//export proverGetCredentialsCB
func proverGetCredentialsCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, credentialsJson *C.char) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle),
			indyUtils.IndyResult{Error: nil,
				Results: []interface{}{
					string(C.GoString(credentialsJson)),
				}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

func ProverGetCredentials(walletHandle int, filterJson unsafe.Pointer) chan indyUtils.IndyResult {

	// Prepare the call parameters.
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
			Gets human readable credentials according to the filter.
		    If filter is NULL, then all credentials are returned. Credentials can be filtered by tags created during saving of credential.

			NOTE: This method is deprecated because immediately returns all fetched credentials.
		    Use <prover_search_credentials> to fetch records by small batches.
		    :param wallet_handle: wallet handle (created by open_wallet).
		    :param filter_json: filter for credentials
		        {
		            "schema_id": string, (Optional)
		            "schema_issuer_did": string, (Optional)
		            "schema_name": string, (Optional)
		            "schema_version": string, (Optional)
		            "issuer_did": string, (Optional)
		            "cred_def_id": string, (Optional)
		        }
		    :return:  credentials json
		     [{
		         "referent": string, - id of credential in the wallet
		         "attrs": {"key1":"raw_value1", "key2":"raw_value2"}, - credential attributes
		         "schema_id": string, - identifier of schema
		         "cred_def_id": string, - identifier of credential definition
		         "rev_reg_id": Optional<string>, - identifier of revocation registry definition
		         "cred_rev_id": Optional<string> - identifier of credential in the revocation registry definition
		     }]
	*/

	// Call C.indy_prover_get_credentials
	res := C.indy_prover_get_credentials(commandHandle,
		C.indy_handle_t(walletHandle),
		(*C.char)(filterJson),
		(C.cb_proverGetCredentials)(unsafe.Pointer(C.proverGetCredentialsCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}
	return future
}

//export proverGetCredentialsForProofReqCB
func proverGetCredentialsForProofReqCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, credentialsJson *C.char) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle),
			indyUtils.IndyResult{Error: nil,
				Results: []interface{}{
					string(C.GoString(credentialsJson)),
				}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

func ProverGetCredentialsForProofReq(walletHandle int, proofReqJson unsafe.Pointer) chan indyUtils.IndyResult {

	// Prepare the call parameters.
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
			Gets human readable credentials matching the given proof request.

			NOTE: This method is deprecated because immediately returns all fetched credentials.
		    Use <prover_search_credentials_for_proof_req> to fetch records by small batches.

			:param wallet_handle: wallet handle (created by open_wallet).
		    :param proof_request_json: proof request json
		        {
		            "name": string,
		            "version": string,
		            "nonce": string, - a decimal number represented as a string (use `indy_generate_nonce` function to generate 80-bit number)
		            "requested_attributes": { // set of requested attributes
		                 "<attr_referent>": <attr_info>, // see below
		                 ...,
		            },
		            "requested_predicates": { // set of requested predicates
		                 "<predicate_referent>": <predicate_info>, // see below
		                 ...,
		             },
		            "non_revoked": Optional<<non_revoc_interval>>, // see below,
		                           // If specified prover must proof non-revocation
		                           // for date in this interval for each attribute
		                           // (applies to every attribute and predicate but can be overridden on attribute level)
		            "ver": Optional<str>  - proof request version:
		                - omit to use unqualified identifiers for restrictions
		                - "1.0" to use unqualified identifiers for restrictions
		                - "2.0" to use fully qualified identifiers for restrictions
		        }
		    where
		    attr_referent: Proof-request local identifier of requested attribute
		    attr_info: Describes requested attribute
		        {
		            "name": Optional<string>, // attribute name, (case insensitive and ignore spaces)
		            "names": Optional<[string, string]>, // attribute names, (case insensitive and ignore spaces)
		                                                 // NOTE: should either be "name" or "names", not both and not none of them.
		                                                 // Use "names" to specify several attributes that have to match a single credential.
		            "restrictions": Optional<filter_json>, // see below
		            "non_revoked": Optional<<non_revoc_interval>>, // see below,
		                           // If specified prover must proof non-revocation
		                           // for date in this interval this attribute
		                           // (overrides proof level interval)
		        }
		    predicate_referent: Proof-request local identifier of requested attribute predicate
		    predicate_info: Describes requested attribute predicate
		        {
		            "name": attribute name, (case insensitive and ignore spaces)
		            "p_type": predicate type (">=", ">", "<=", "<")
		            "p_value": int predicate value
		            "restrictions": Optional<filter_json>, // see below
		            "non_revoked": Optional<<non_revoc_interval>>, // see below,
		                           // If specified prover must proof non-revocation
		                           // for date in this interval this attribute
		                           // (overrides proof level interval)
		        }
		    non_revoc_interval: Defines non-revocation interval
		        {
		            "from": Optional<int>, // timestamp of interval beginning
		            "to": Optional<int>, // timestamp of interval ending
		        }
		     filter_json:
		        {
		           "schema_id": string, (Optional)
		           "schema_issuer_did": string, (Optional)
		           "schema_name": string, (Optional)
		           "schema_version": string, (Optional)
		           "issuer_did": string, (Optional)
		           "cred_def_id": string, (Optional)
		        }

		    :return: json with credentials for the given proof request.
		        {
		            "attrs": {
		                "<attr_referent>": [{ cred_info: <credential_info>, interval: Optional<non_revoc_interval> }],
		                ...,
		            },
		            "predicates": {
		                "requested_predicates": [{ cred_info: <credential_info>, timestamp: Optional<integer> }, { cred_info: <credential_2_info>, timestamp: Optional<integer> }],
		                "requested_predicate_2_referent": [{ cred_info: <credential_2_info>, timestamp: Optional<integer> }]
		            }
		        }, where <credential_info> is
		        {
		            "referent": string, - id of credential in the wallet
		            "attrs": {"key1":"raw_value1", "key2":"raw_value2"}, - credential attributes
		            "schema_id": string, - identifier of schema
		            "cred_def_id": string, - identifier of credential definition
		            "rev_reg_id": Optional<string>, - identifier of revocation registry definition
		            "cred_rev_id": Optional<string> - identifier of credential in the revocation registry definition
		        }
	*/

	// Call C.indy_prover_get_credentials_for_proof_req
	res := C.indy_prover_get_credentials_for_proof_req(commandHandle,
		C.indy_handle_t(walletHandle),
		(*C.char)(proofReqJson),
		(C.cb_proverGetCredentialsForProofReq)(unsafe.Pointer(C.proverGetCredentialsForProofReqCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}
	return future
}

//export proverSearchCredentialsCB
func proverSearchCredentialsCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, searchHandle C.indy_handle_t, totalCount C.indy_u32_t) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle),
			indyUtils.IndyResult{Error: nil,
				Results: []interface{}{
					int(searchHandle),
					int(totalCount),
				}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

func ProverSearchCredentials(walletHandle int, queryJson unsafe.Pointer) chan indyUtils.IndyResult {

	// Prepare the call parameters.
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
			Search for credentials stored in wallet.
		    Credentials can be filtered by tags created during saving of credential.

		    Instead of immediately returning of fetched credentials this call returns search_handle that can be used later
			to fetch records by small batches (with prover_credentials_search_fetch_records).

			:param wallet_handle: wallet handle (created by open_wallet).
		    :param query_json: wql style filter for credentials searching based on tags.
		        where wql query: indy-sdk/docs/design/011-wallet-query-language/README.md
		    :return:
		        search_handle: Search handle that can be used later to fetch records by small batches
		            (with prover_credentials_search_fetch_records)
		        total_count: Total count of records
	*/

	// Call C.indy_prover_search_credentials
	res := C.indy_prover_search_credentials(commandHandle,
		C.indy_handle_t(walletHandle),
		(*C.char)(queryJson),
		(C.cb_proverSearchCredentials)(unsafe.Pointer(C.proverSearchCredentialsCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}
	return future
}

//export proverFetchCredentialsCB
func proverFetchCredentialsCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, credentialsJson *C.char) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle),
			indyUtils.IndyResult{Error: nil,
				Results: []interface{}{
					string(C.GoString(credentialsJson)),
				}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

func ProverFetchCredentials(searchHandle int, totalCount int) chan indyUtils.IndyResult {

	// Prepare the call parameters.
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
	   Fetch next credentials for search.
	   :param search_handle: Search handle (created by prover_open_credentials_search)
	   :param count: Count of records to fetch
	   :return: credentials_json: List of credentials:
	   [{
	        "referent": string, - id of credential in the wallet
	        "attrs": {"key1":"raw_value1", "key2":"raw_value2"}, - credential attributes
	        "schema_id": string, - identifier of schema
	        "cred_def_id": string, - identifier of credential definition
	        "rev_reg_id": Optional<string>, - identifier of revocation registry definition
	        "cred_rev_id": Optional<string> - identifier of credential in the revocation registry definition
	   }]
	   NOTE: The list of length less than the requested count means credentials search iterator is completed.
	*/

	// Call C.indy_prover_fetch_credentials
	res := C.indy_prover_fetch_credentials(commandHandle,
		C.indy_handle_t(searchHandle),
		C.indy_u32_t(totalCount),
		(C.cb_proverFetchCredentials)(unsafe.Pointer(C.proverFetchCredentialsCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}
	return future
}

//export toUnqualifiedCB
func toUnqualifiedCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, res *C.char) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle),
			indyUtils.IndyResult{Error: nil,
				Results: []interface{}{
					string(C.GoString(res)),
				}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

func ToUnqualified(entity unsafe.Pointer) chan indyUtils.IndyResult {

	// Prepare the call parameters.
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
			Get unqualified form (short form without method) of a fully qualified entity like DID.

		    This function should be used to the proper casting of fully qualified entity to unqualified form in the following cases:
		        - Issuer, which works with fully qualified identifiers, creates a Credential Offer for Prover, which doesn't support fully qualified identifiers.
		        - Verifier prepares a Proof Request based on fully qualified identifiers or Prover, which doesn't support fully qualified identifiers.
		        - another case when casting to unqualified form needed

			:param entity: target entity to disqualify. Can be one of:
		                Did
		                SchemaId
		                CredentialDefinitionId
		                RevocationRegistryId
		                Schema
		                CredentialDefinition
		                RevocationRegistryDefinition
		                CredentialOffer
		                CredentialRequest
		                ProofRequest
		    :return: entity either in unqualified form or original if casting isn't possible
	*/

	// Call C.indy_to_unqualified
	res := C.indy_to_unqualified(commandHandle,
		(*C.char)(entity),
		(C.cb_proverGetCredentials)(unsafe.Pointer(C.proverGetCredentialsCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}
	return future
}
