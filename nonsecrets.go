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

/*
#include <stdlib.h>
*/
import "C"
import (
	"github.com/joyride9999/IndySdkGoBindings/nonsecrets"
	"unsafe"
)

// IndyAddWalletRecord Create a new non-secret record in the wallet.
func IndyAddWalletRecord(wh int, recordType string, recordId string, recordValue string, tagsJson string) (err error) {

	upRecordType := unsafe.Pointer(C.CString(recordType))
	defer C.free(upRecordType)
	upRecordId := unsafe.Pointer(C.CString(recordId))
	defer C.free(upRecordId)
	upRecordValue := unsafe.Pointer(C.CString(recordValue))
	defer C.free(upRecordValue)
	upRecordTag := unsafe.Pointer(GetOptionalValue(tagsJson))
	defer C.free(upRecordTag)

	channel := nonsecrets.IndyAddWalletRecord(wh, upRecordType, upRecordId, upRecordValue, upRecordTag)
	result := <-channel
	return result.Error
}

// IndyAddWalletRecordTags Add new tags to the wallet record.
func IndyAddWalletRecordTags(wh int, recordType string, recordId string, tagsJson string) (err error) {
	upRecordType := unsafe.Pointer(C.CString(recordType))
	defer C.free(upRecordType)
	upRecordId := unsafe.Pointer(C.CString(recordId))
	defer C.free(upRecordId)
	upRecordTag := unsafe.Pointer(GetOptionalValue(tagsJson))
	defer C.free(upRecordTag)

	channel := nonsecrets.IndyAddWalletRecordTags(wh, upRecordType, upRecordId, upRecordTag)
	result := <-channel
	return result.Error
}

// IndyGetWalletRecord Create a new non-secret record in the wallet.
func IndyGetWalletRecord(wh int, recordType string, recordId string, options string) (recordJson string, err error) {

	upRecordType := unsafe.Pointer(C.CString(recordType))
	defer C.free(upRecordType)
	upRecordId := unsafe.Pointer(C.CString(recordId))
	defer C.free(upRecordId)

	if len(options) == 0 {
		options = "{}"
	}

	upOptions := unsafe.Pointer(C.CString(options))
	defer C.free(upOptions)

	channel := nonsecrets.IndyGetWalletRecord(wh, upRecordType, upRecordId, upOptions)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// IndyDeleteWalletRecord Delete an existing wallet record in the wallet.
func IndyDeleteWalletRecord(wh int, recordType string, recordId string) (err error) {
	upRecordType := unsafe.Pointer(C.CString(recordType))
	defer C.free(upRecordType)
	upRecordId := unsafe.Pointer(C.CString(recordId))
	defer C.free(upRecordId)

	channel := nonsecrets.IndyDeleteWalletRecord(wh, upRecordType, upRecordId)
	result := <-channel
	return result.Error
}

// IndyDeleteWalletRecordTags Delete tags from the wallet record.
func IndyDeleteWalletRecordTags(wh int, recordType string, recordId string, tagNames string) (err error) {
	upRecordType := unsafe.Pointer(C.CString(recordType))
	defer C.free(upRecordType)
	upRecordId := unsafe.Pointer(C.CString(recordId))
	defer C.free(upRecordId)
	upTagNames := unsafe.Pointer(GetOptionalValue(tagNames))
	defer C.free(upTagNames)

	channel := nonsecrets.IndyDeleteWalletRecordTags(wh, upRecordType, upRecordId, upTagNames)
	result := <-channel
	return result.Error
}

// IndyUpdateWalletRecordValue Update a non-secret wallet record value.
func IndyUpdateWalletRecordValue(wh int, recordType string, recordId string, recordValue string) (err error) {

	upRecordType := unsafe.Pointer(C.CString(recordType))
	defer C.free(upRecordType)
	upRecordId := unsafe.Pointer(C.CString(recordId))
	defer C.free(upRecordId)
	upRecordValue := unsafe.Pointer(C.CString(recordValue))
	defer C.free(upRecordValue)

	channel := nonsecrets.IndyUpdateWalletRecordValue(wh, upRecordType, upRecordId, upRecordValue)
	result := <-channel
	return result.Error
}

// IndyUpdateWalletRecordTags Update a non-secret wallet record value.
func IndyUpdateWalletRecordTags(wh int, recordType string, recordId string, recordTags string) (err error) {
	upRecordType := unsafe.Pointer(C.CString(recordType))
	defer C.free(upRecordType)
	upRecordId := unsafe.Pointer(C.CString(recordId))
	defer C.free(upRecordId)
	upRecordTags := unsafe.Pointer(C.CString(recordTags))
	defer C.free(upRecordTags)

	channel := nonsecrets.IndyUpdateWalletRecordTags(wh, upRecordType, upRecordId, upRecordTags)
	result := <-channel
	return result.Error
}

// IndyOpenWalletSearch Search for wallet records.
func IndyOpenWalletSearch(wh int, recordType string, query string, options string) (searchHandle int, err error) {
	upRecordType := unsafe.Pointer(C.CString(recordType))
	defer C.free(upRecordType)
	upQuery := unsafe.Pointer(GetOptionalValue(query))
	defer C.free(upQuery)
	upOptions := unsafe.Pointer(GetOptionalValue(options))
	defer C.free(upOptions)

	channel := nonsecrets.IndyOpenWalletSearch(wh, upRecordType, upQuery, upOptions)
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
