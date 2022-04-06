/*
// ******************************************************************
// Purpose: exported public functions that handles nonsecrets functions
// from libindy
// Author:  adrian.toader@siemens.com
// Notes:
// Copyright (c): Siemens SRL
// This work is licensed under the terms of the Apache License Version 2.0.  See
// the LICENSE.txt file in the top-level directory.
// ******************************************************************
*/

package indySDK

import "github.com/joyride9999/IndySdkGoBindings/nonsecrets"

// IndyAddWalletRecord Create a new non-secret record in the wallet.
func IndyAddWalletRecord(wh int, recordType string, recordId string, recordValue string, tagsJson string) (err error) {
	channel := nonsecrets.IndyAddWalletRecord(wh, recordType, recordId, recordValue, tagsJson)
	result := <-channel
	return result.Error
}

// IndyAddWalletRecordTags Add new tags to the wallet record.
func IndyAddWalletRecordTags(wh int, recordType string, recordId string, tagsJson string) (err error) {
	channel := nonsecrets.IndyAddWalletRecordTags(wh, recordType, recordId, tagsJson)
	result := <-channel
	return result.Error
}

// IndyGetWalletRecord Create a new non-secret record in the wallet.
func IndyGetWalletRecord(wh int, recordType string, recordId string, options string) (recordJson string, err error) {
	channel := nonsecrets.IndyGetWalletRecord(wh, recordType, recordId, options)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// IndyDeleteWalletRecord Delete an existing wallet record in the wallet.
func IndyDeleteWalletRecord(wh int, recordType string, recordId string) (err error) {
	channel := nonsecrets.IndyDeleteWalletRecord(wh, recordType, recordId)
	result := <-channel
	return result.Error
}

// IndyDeleteWalletRecordTags Delete tags from the wallet record.
func IndyDeleteWalletRecordTags(wh int, recordType string, recordId string, tagNames string) (err error) {
	channel := nonsecrets.IndyDeleteWalletRecordTags(wh, recordType, recordId, tagNames)
	result := <-channel
	return result.Error
}

// IndyUpdateWalletRecordValue Update a non-secret wallet record value.
func IndyUpdateWalletRecordValue(wh int, recordType string, recordId string, recordValue string) (err error) {
	channel := nonsecrets.IndyUpdateWalletRecordValue(wh, recordType, recordId, recordValue)
	result := <-channel
	return result.Error
}

// IndyUpdateWalletRecordTags Update a non-secret wallet record value.
func IndyUpdateWalletRecordTags(wh int, recordType string, recordId string, recordTags string) (err error) {
	channel := nonsecrets.IndyUpdateWalletRecordTags(wh, recordType, recordId, recordTags)
	result := <-channel
	return result.Error
}

// IndyOpenWalletSearch Search for wallet records.
func IndyOpenWalletSearch(wh int, recordType string, query string, options string) (searchHandle int, err error) {
	channel := nonsecrets.IndyOpenWalletSearch(wh, recordType, query, options)
	result := <-channel
	if result.Error != nil {
		return 0, result.Error
	}
	return result.Results[0].(int), result.Error
}

// IndyFetchWalletSearchNextRecords Fetch next records for wallet search.
func IndyFetchWalletSearchNextRecords(wh int, sh int, count int32) (recordsJson string, err error) {
	channel := nonsecrets.IndyFetchWalletSearchNextRecords(wh, sh, count)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// IndyCloseWalletSearch Close wallet search (make search handle invalid).
func IndyCloseWalletSearch(sh int) (err error) {
	channel := nonsecrets.IndyCloseWalletSearch(sh)
	result := <-channel
	return result.Error
}
