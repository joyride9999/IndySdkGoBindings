/*
// ******************************************************************
// Purpose: nonsecrets unit testing
// Author: angel.draghici@siemens.com, adrian.toader@siemens.com
// Notes:
// Copyright (c): Siemens SRL
// This work is licensed under the terms of the Apache License Version 2.0.  See
// the LICENSE.txt file in the top-level directory.
// ******************************************************************
*/

package indySDK

import (
	"github.com/Jeffail/gabs/v2"
	"github.com/joyride9999/IndySdkGoBindings/indyUtils"
	"testing"
)

func TestIndyAddWalletRecord(t *testing.T) {
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("createWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	type args struct {
		RecordId string
		Tags string
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"test-add-wallet-record-works", args{RecordId: recordId1, Tags: recordTags1}, false},
		{"test-add-wallet-record-without-tags", args{RecordId: recordId2, Tags: ""}, false},
		{"test-add-wallet-record-empty-id", args{RecordId: ""}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errAdd := IndyAddWalletRecord(walletHandle, recordType, tt.args.RecordId, recordValue1, tt.args.Tags)
			hasError := errAdd != nil
			if hasError != tt.wantErr {
				t.Errorf("IndyAddWalletRecord() error = '%v'", errAdd)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errAdd)
				return
			}

			// Test if non-secret record was added.
			_, errGet := IndyGetWalletRecord(walletHandle, recordType, tt.args.RecordId,
				recordOptions)
			if errGet != nil {
				t.Errorf("IndyGetWalletRecord() error = '%v'", errGet)
				return
			}
		})
	}

	return
}

func TestIndyAddWalletRecordTags(t *testing.T) {
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("createWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	errAddRecord := IndyAddWalletRecord(walletHandle, recordType, recordId1, recordValue1, ""); if errAddRecord != nil {
		t.Errorf("IndyAddWalletRecord() error = '%v'", errAddRecord)
		return
	}

	ok, errCheck := checkRecordTags(walletHandle, recordType, recordId1, recordOptions, "{}"); if errCheck != nil {
		t.Errorf("checkRecordTags() error = '%v'", errCheck)
		return
	}
	if ok == false {
		t.Error("Invalid values")
		return
	}

	type args struct {
		RecordId string
		Tags string
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"test-indy-add-wallet-record-tags-works", args{RecordId: recordId1, Tags: recordTags1}, false},
		{"test-indy-add-wallet-record-tags-works-add-another-tag", args{RecordId: recordId1, Tags: `{"tagName4": "str4"}`}, false},
		{"test-indy-add-wallet-record-tags-record-not-found", args{RecordId: "notFoundRecord"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errAddTags := IndyAddWalletRecordTags(walletHandle, recordType, tt.args.RecordId,
				tt.args.Tags)
			hasError := errAddTags != nil
			if hasError != tt.wantErr {
				t.Errorf("IndyAddWalletRecordTags() error = '%v'", errAddTags)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errAddTags)
				return
			}
			if tt.args.Tags == recordTags1 {
				ok, errCheck2 := checkRecordTags(walletHandle, recordType, recordId1, recordOptions, recordTags1)
				if errCheck2 != nil {
					t.Errorf("checkRecordTags() error = '%v'", errCheck2)
					return
				}
				if ok == false {
					t.Errorf("Values are not correct")
				}
			}
			if tt.args.Tags != recordTags1 {
				record, errGet := IndyGetWalletRecord(walletHandle, recordType, recordId1, recordOptions)
				if errGet != nil {
					t.Errorf("IndyGetWalletRecord() error = '%v'", errGet)
					return
				}
				recordParsed, _ := gabs.ParseJSON([]byte(record))
				if recordParsed.Path("tags").String() == recordTags1 {
					t.Errorf("Test failed")
					return
				}
			}
		})
	}

	return
}

func TestIndyDeleteWalletRecord(t *testing.T) {
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("createWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	errAdd := IndyAddWalletRecord(walletHandle, recordType, recordId1, recordValue1, recordTags1)
	if errAdd != nil {
		t.Errorf("IndyAddWalletRecord() error = '%v'", errAdd)
		return
	}

	type args struct {
		RecordId string
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"test-delete-wallet-record-works", args{RecordId: recordId1}, false},
		{"test-delete-wallet-record-not-found-record", args{RecordId: "notFoundRecord"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errDelete := IndyDeleteWalletRecord(walletHandle, recordType, tt.args.RecordId)
			hasError := errDelete != nil
			if hasError != tt.wantErr {
				t.Errorf("IndyDeleteWalletRecord() error = '%v'", errDelete)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errDelete)
				return
			}
			_, errGet := IndyGetWalletRecord(walletHandle, recordType, tt.args.RecordId, recordOptions)
			if errGet != nil && errGet.Error() != indyUtils.GetIndyError(212) {
				t.Errorf("IndyGetWalletRecord() error = '%v'", errGet)
				return
			}
		})
	}

	return
}

func TestIndyGetWalletRecord(t *testing.T) {
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("createWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	errAdd := IndyAddWalletRecord(walletHandle, recordType, recordId1, recordValue1, recordTags1)
	if errAdd != nil {
		t.Errorf("IndyAddWalletRecord() error = '%v'", errAdd)
		return
	}

	expected1 := `{"id": "recordId1", "value": "recordValue", "tags": null, "type": null}`
	expectedRecord, errGabs := gabs.ParseJSON([]byte(expected1)); if errGabs != nil {
		t.Errorf("Gabs Parse error = '%v'", errGabs)
		return
	}

	expected2 := `{"id": "recordId1", "value": "recordValue", "tags": {"tagName1":"str1","tagName2":"5","tagName3":"12"}, "type": "testType"}`
	expectedRecord2, errGabs2 := gabs.ParseJSON([]byte(expected2)); if errGabs2 != nil {
		t.Errorf("Gabs Parse error = '%v'", errGabs2)
		return
	}

	type args struct {
		RecordId string
		Options string
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"test-get-wallet-record-works", args{RecordId: recordId1, Options: "{}"}, false},
		{"test-get-wallet-record-full-data", args{RecordId: recordId1, Options: recordOptions}, false},
		{"test-get-wallet-record-not-found-record", args{RecordId: recordId2, Options: "{}"}, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			record, errGet := IndyGetWalletRecord(walletHandle, recordType, tt.args.RecordId, tt.args.Options)
			hasError := errGet != nil
			if hasError != tt.wantErr {
				t.Errorf("IndyGetWalletRecord() error = '%v'", errGet)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errGet)
				return
			}

			returnedRecord, errGabs2 := gabs.ParseJSON([]byte(record)); if errGabs2 != nil {
				t.Errorf("Gabs Parse error = '%v'", errGabs2)
				return
			}

			if tt.args.Options == "{}" {
				if returnedRecord.Path("tags").String() != expectedRecord.Path("tags").String() ||
					returnedRecord.Path("type").String() != expectedRecord.Path("type").String() {
					t.Errorf("Test failed")
				}
			} else if tt.args.Options == recordOptions {
				if returnedRecord.Path("tags").String() != expectedRecord2.Path("tags").String() ||
					returnedRecord.Path("type").String() != expectedRecord2.Path("type").String() {
					t.Errorf("Test failed")
				}
			}
		})
	}

	return
}

func TestIndyDeleteWalletRecordTags(t *testing.T) {
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("createWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	errAddRecord := IndyAddWalletRecord(walletHandle, recordType, recordId1, recordValue1, recordTags1); if errAddRecord != nil {
		t.Errorf("IndyAddWalletRecord() error = '%v'", errAddRecord)
		return
	}

	ok, errCheck := checkRecordTags(walletHandle, recordType, recordId1, recordOptions, recordTags1); if errCheck != nil {
		t.Errorf("checkRecordTags() error = '%v'", errCheck)
		return
	}
	if ok == false {
		t.Error("Invalid values")
		return
	}
	expectedTags := `{"tagName2": "5", "tagName3": "12"}`

	type args struct {
		RecordId string
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"test-delete-wallet-record-tags-works", args{RecordId: recordId1}, false},
		{"test-delete-wallet-record-not-found-record", args{RecordId: recordId2}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errDelete := IndyDeleteWalletRecordTags(walletHandle, recordType, tt.args.RecordId, `["tagName1"]`)
			hasError := errDelete != nil
			if hasError != tt.wantErr {
				t.Errorf("IndyDeleteWalletRecordTags() error = '%v'", errDelete)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errDelete)
				return
			}

			ok, errCheck = checkRecordTags(walletHandle, recordType, recordId1, recordOptions, expectedTags); if errCheck != nil {
				t.Errorf("checkRecordTags() error = '%v'", errCheck)
				return
			}
			if ok == false {
				t.Errorf("Test failed")
			}
		})
	}
}

func TestIndyUpdateWalletRecordValue(t *testing.T) {
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("createWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	errAddRecord := IndyAddWalletRecord(walletHandle, recordType, recordId1, recordValue1, "{}"); if errAddRecord != nil {
		t.Errorf("IndyAddWalletRecord() error = '%v'", errAddRecord)
		return
	}

	ok, errCheck := checkRecordValue(walletHandle, recordType, recordId1, recordOptions, recordValue1); if errCheck != nil {
		t.Errorf("checkRecordValue() error = '%v'", errCheck)
		return
	}
	if ok == false {
		t.Error("Invalid values")
		return
	}

	type args struct {
		RecordId string
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"test-update-wallet-record-value", args{RecordId: recordId1}, false},
		{"test-update-wallet-record-value-not-found-record", args{RecordId: recordId2}, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errUpdate := IndyUpdateWalletRecordValue(walletHandle, recordType, tt.args.RecordId, "recordValue2")
			hasError := errUpdate != nil
			if hasError != tt.wantErr {
				t.Errorf("IndyUpdateWalletRecordValue() error = '%v'", errUpdate)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errUpdate)
				return
			}
			ok, errCheck = checkRecordValue(walletHandle, recordType,tt.args.RecordId, recordOptions, "recordValue2"); if errCheck != nil {
				t.Errorf("checkRecordField() error = '%v'", errCheck)
				return
			}
			if ok != true {
				t.Errorf("Test failed")
			}
		})
	}

	return
}

func TestIndyUpdateWalletRecordTags(t *testing.T) {
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("createWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	errAddRecord := IndyAddWalletRecord(walletHandle, recordType, recordId1, recordValue1, ""); if errAddRecord != nil {
		t.Errorf("IndyAddWalletRecord() error = '%v'", errAddRecord)
		return
	}

	ok, errCheck := checkRecordTags(walletHandle, recordType, recordId1, recordOptions, "{}"); if errCheck != nil {
		t.Errorf("checkRecordTags() error = '%v'", errCheck)
		return
	}
	if ok == false {
		t.Error("Invalid values")
		return
	}

	type args struct {
		RecordId string
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"test-update-wallet-record-tags-works", args{RecordId: recordId1}, false},
		{"test-update-wallet-record-tags-not-found-record", args{RecordId: recordId2}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errUpdate := IndyUpdateWalletRecordTags(walletHandle, recordType, tt.args.RecordId, recordTags1)
			hasError := errUpdate != nil
			if hasError != tt.wantErr {
				t.Errorf("IndyUpdateWalletRecordTags() error = '%v'", errUpdate)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errUpdate)
				return
			}
			ok, errCheck = checkRecordTags(walletHandle, recordType, recordId1, recordOptions, recordTags1); if errCheck != nil {
				t.Errorf("checkRecordField() error = '%v'", errCheck)
				return
			}
			if ok == false {
				t.Error("Test failed")
				return
			}
		})
	}
	return
}

func TestIndyOpenWalletSearch(t *testing.T) {
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("createWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	errAddRecord := IndyAddWalletRecord(walletHandle, recordType, recordId1, recordValue1, recordTags1); if errAddRecord != nil {
		t.Errorf("IndyAddWalletRecord() error = '%v'", errAddRecord)
		return
	}

	errAddRecord = IndyAddWalletRecord(walletHandle, recordType, recordId2, recordValue2, recordTags2); if errAddRecord != nil {
		t.Errorf("IndyAddWalletRecord() error = '%v'", errAddRecord)
		return
	}

	type args struct {
		WalletHandle int
		Query string
		Options string
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"test-open-wallet-search-works", args{WalletHandle: walletHandle, Query: "", Options: ""}, false},
		{"test-open-wallet-search-works-full-params", args{WalletHandle: walletHandle, Query: `{"tagName1": "str2"}`, Options: recordOptions}, false},
		{"test-open-wallet-search-invalid-wallet-handle", args{WalletHandle: walletHandle + 100}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			searchHandle, errOpenSearch := IndyOpenWalletSearch(tt.args.WalletHandle, recordType, tt.args.Query, tt.args.Options)
			hasError := errOpenSearch != nil
			if hasError != tt.wantErr {
				t.Errorf("IndyOpenWalletSearch() error = '%v'", errOpenSearch)
				return
			}
			defer IndyCloseWalletSearch(searchHandle)

			if tt.wantErr {
				t.Log("Expected error: ", errOpenSearch)
				return
			}

			if searchHandle == 0 {
				t.Errorf("Test failed")
				return
			}
		})
	}

	return
}

func TestIndyFetchWalletSearchNextRecords(t *testing.T) {
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("createWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	errAddRecord := IndyAddWalletRecord(walletHandle, recordType, recordId1, recordValue1, recordTags1);
	if errAddRecord != nil {
		t.Errorf("IndyAddWalletRecord() error = '%v'", errAddRecord)
		return
	}

	errAddRecord = IndyAddWalletRecord(walletHandle, recordType, recordId2, recordValue2, recordTags2);
	if errAddRecord != nil {
		t.Errorf("IndyAddWalletRecord() error = '%v'", errAddRecord)
		return
	}
	expectedRecords := `{"records": 
	[
		{"id":"recordId1","value":"recordValue","tags": null, "type": null}, 
		{"id":"recordId2","value":"recordValue2","tags": null, "type": null}
	]
}`
	expectedRecordsParsed, _ := gabs.ParseJSON([]byte(expectedRecords))
	options := `{ 
		"retrieveRecords": true,
        "retrieveTotalCount": false,
        "retrieveType": false,
        "retrieveValue": true,
        "retrieveTags": false
	}`

	type args struct {
		SearchHandle int
		Query string
		Options string
		ExpectedRecords *gabs.Container
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"test-fetch-wallet-search-next-records", args{Query: "", Options: "", ExpectedRecords: expectedRecordsParsed}, false},
		{"test-fetch-wallet-search-next-records-for-options", args{Query: "", Options: options, ExpectedRecords: expectedRecordsParsed}, false},
		{"test-fetch-wallet-search-next-records-for-query", args{Query: `{"tagName1": "str2"}`, Options: "",
			ExpectedRecords: expectedRecordsParsed}, false},
		{"test-fetch-wallet-search-next-records-invalid-search-handle", args{SearchHandle: 100, ExpectedRecords: expectedRecordsParsed}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			searchHandle, errOpenSearch := IndyOpenWalletSearch(walletHandle, recordType, tt.args.Query, tt.args.Options)
			if errOpenSearch != nil {
				t.Errorf("IndyOpenWalletSearch() error = '%v'", errOpenSearch)
				return
			}
			defer IndyCloseWalletSearch(searchHandle)

			if tt.args.SearchHandle != 0 {
				searchHandle = tt.args.SearchHandle
			}
			searchRecords, errFetch := IndyFetchWalletSearchNextRecords(walletHandle, searchHandle, int32(2))
			hasError := errFetch != nil
			if hasError != tt.wantErr {
				t.Errorf("IndyFetchWalletSearchNextRecords() error = '%v'", errFetch)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errFetch)
				return
			}

			searchRecordsParsed, _ := gabs.ParseJSON([]byte(searchRecords))
			for _, search := range searchRecordsParsed.S("records").Children() {
				ok := false
				for _, expected := range expectedRecordsParsed.S("records").Children() {
					if search.String() == expected.String() {
						ok = true
					}
					if ok == true {
						break
					}
				}
				if !ok {
					t.Errorf("Test failed")
					break
				}
			}
		})
	}

	return
}

func TestIndyCloseWalletSearch(t *testing.T) {
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("createWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	errAddRecord := IndyAddWalletRecord(walletHandle, recordType, recordId1, recordValue1, recordTags1);
	if errAddRecord != nil {
		t.Errorf("IndyAddWalletRecord() error = '%v'", errAddRecord)
		return
	}

	searchHandle, errOpen := IndyOpenWalletSearch(walletHandle, recordType, "", "")
	if errOpen != nil {
		t.Errorf("IndyOpenWalletSearch() error = '%v'", errOpen)
		return
	}

	type args struct {
		SearchHandle int
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"test-close-wallet-search-works", args{SearchHandle: searchHandle}, false},
		{"test-close-wallet-search-invalid-handle", args{SearchHandle: 100}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errClose := IndyCloseWalletSearch(tt.args.SearchHandle)
			hasError := errClose != nil
			if hasError != tt.wantErr {
				t.Errorf("IndyCloseWalletSearch() error = '%v'", errClose)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errClose)
			}
		})
	}
	return
}