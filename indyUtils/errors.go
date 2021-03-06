/*
// ******************************************************************
// Purpose: maps indy error codes to messages
// Author:  alexandru.leonte@siemens.com
// Notes:
// Copyright (c): Siemens SRL
// This work is licensed under the terms of the Apache License Version 2.0.  See
// the LICENSE.txt file in the top-level directory.
// ******************************************************************
*/

package indyUtils

var errorsIndySDK = map[int]string{
	0:   "Success",
	100: "CommonInvalidParam1: Caller passed invalid value as param 1 (null, invalid json and etc..)",
	101: "CommonInvalidParam2: Caller passed invalid value as param 2 (null, invalid json and etc..)",
	102: "CommonInvalidParam3: Caller passed invalid value as param 3 (null, invalid json and etc..)",
	103: "CommonInvalidParam4: Caller passed invalid value as param 4 (null, invalid json and etc..)",
	104: "CommonInvalidParam5: Caller passed invalid value as param 5 (null, invalid json and etc..)",
	105: "CommonInvalidParam6: Caller passed invalid value as param 6 (null, invalid json and etc..)",
	106: "CommonInvalidParam7: Caller passed invalid value as param 7 (null, invalid json and etc..)",
	107: "CommonInvalidParam8: Caller passed invalid value as param 8 (null, invalid json and etc..)",
	108: "CommonInvalidParam9: Caller passed invalid value as param 9 (null, invalid json and etc..)",
	109: "CommonInvalidParam10: Caller passed invalid value as param 10 (null, invalid json and etc..)",
	110: "CommonInvalidParam11: Caller passed invalid value as param 11 (null, invalid json and etc..)",
	111: "CommonInvalidParam12: Caller passed invalid value as param 12 (null, invalid json and etc..)",
	112: "CommonInvalidState: Invalid library state was detected in runtime. It signals library bug",
	113: "CommonInvalidStructure: Object (json, config, key, credential and etc...) passed by library caller has invalid structure",
	114: "CommonIOError: IO Error",
	200: "WalletInvalidHandle: Caller passed invalid wallet handle",
	201: "WalletUnknownTypeError: Unknown type of wallet was passed on create_wallet",
	202: "WalletTypeAlreadyRegisteredError: Attempt to register already existing wallet type",
	203: "WalletAlreadyExistsError: Attempt to create wallet with name used for another exists wallet",
	204: "WalletNotFoundError: Requested entity id isn't present in wallet",
	205: "WalletIncompatiblePoolError: Trying to use wallet with pool that has different name",
	206: "WalletAlreadyOpenedError: Trying to open wallet that was opened already",
	207: "WalletAccessFailed: Attempt to open encrypted wallet with invalid credentials",
	208: "WalletInputError: Input provided to wallet operations is considered not valid",
	209: "WalletDecodingError: Decoding of wallet data during input/output failed",
	210: "WalletStorageError: Storage error occurred during wallet operation",
	211: "WalletEncryptionError: Error during encryption-related operations",
	212: "WalletItemNotFound: Requested wallet item not found",
	213: "WalletItemAlreadyExists: Returned if wallet's add_record operation is used with record name that already exists",
	214: "WalletQueryError: Returned if provided wallet query is invalid",
	300: "PoolLedgerNotCreatedError: Trying to open pool ledger that wasn't created before",
	301: "PoolLedgerInvalidPoolHandle: Caller passed invalid pool ledger handle",
	302: "PoolLedgerTerminated: Pool ledger terminated",
	303: "LedgerNoConsensusError: No concensus during ledger operation",
	304: "LedgerInvalidTransaction: Attempt to parse invalid transaction response",
	305: "LedgerSecurityError: Attempt to send transaction without the necessary privileges",
	306: "PoolLedgerConfigAlreadyExistsError: Attempt to create pool ledger config with name used for another existing pool",
	307: "PoolLedgerTimeout: Timeout for action",
	308: "PoolIncompatibleProtocolVersion: Attempt to open Pool for witch Genesis Transactions are not compatible with set Protocol version. Call pool.indy_set_protocol_version to set correct Protocol version.",
	309: "LedgerNotFound: Item not found on ledger",
	400: "AnoncredsRevocationRegistryFullError: Revocation registry is full and creation of new registry is necessary",
	401: "AnoncredsInvalidUserRevocId",
	404: "AnoncredsMasterSecretDuplicateNameError: Attempt to generate master secret with dupplicated name",
	405: "AnoncredsProofRejected",
	406: "AnoncredsCredentialRevoked",
	407: "AnoncredsCredDefAlreadyExistsError: Attempt to create credential definition with duplicated did schema pair",
	500: "UnknownCryptoTypeError: Unknown format of DID entity keys",
	600: "DidAlreadyExistsError: Attempt to create duplicate did",
	700: "PaymentUnknownMethodError: Unknown payment method was given",
	701: "PaymentIncompatibleMethodsError: No method were scraped from inputs/outputs or more than one were scraped",
	702: "PaymentInsufficientFundsError: Insufficient funds on inputs",
	703: "PaymentSourceDoesNotExistError: No such source on a ledger",
	704: "PaymentOperationNotSupportedError: Operation is not supported for payment method",
	705: "PaymentExtraFundsError: Extra funds on inputs",
}

// GetIndyError returns the string error of an indy_error_t
func GetIndyError(err int) string {
	return errorsIndySDK[err]
}
