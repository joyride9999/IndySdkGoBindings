/*
// ******************************************************************
// Purpose: Implements custom wallet storage in memory
// Author:  angel.draghici@siemens.com
// Notes:
// Copyright (c): Siemens SRL
// This work is licensed under the terms of the Apache License Version 2.0.  See
// the LICENSE.txt file in the top-level directory.
// ******************************************************************
*/

package inMemUtils

/*
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"github.com/Jeffail/gabs/v2"
	"github.com/joyride9999/IndySdkGoBindings/indyUtils"
	"github.com/joyride9999/IndySdkGoBindings/wallet"
	cmap "github.com/orcaman/concurrent-map"
	"strconv"
	"unsafe"
)

func NewInMemoryStorage() *InMemoryStorage {
	customStorage := new(InMemoryStorage)
	customStorage.MetadataHandles = cmap.New()
	customStorage.StoredMetadata = cmap.New()
	customStorage.StoredRecords = cmap.New()
	customStorage.StoredTagsP = cmap.New()
	customStorage.StoredTagsE = cmap.New()
	customStorage.WalletHandles = cmap.New()
	customStorage.SearchHandles = cmap.New()
	customStorage.SearchHandlesIterator = cmap.New()
	customStorage.SearchHandlesName = cmap.New()
	customStorage.SearchHandlesValue = cmap.New()
	customStorage.SearchHandlesType = cmap.New()
	customStorage.SearchHandlesTags = cmap.New()

	return customStorage
}

type InMemoryStorage struct {
	MetadataHandles       cmap.ConcurrentMap // Stored metadata of correspondent wallet.
	StoredMetadata        cmap.ConcurrentMap
	WalletHandlesCounter  indyUtils.Counter
	WalletHandles         cmap.ConcurrentMap // Wallet handles
	StoredRecords         cmap.ConcurrentMap
	RecordCounter         indyUtils.Counter
	StoredTagsE           cmap.ConcurrentMap
	StoredTagsP           cmap.ConcurrentMap
	SearchHandles         cmap.ConcurrentMap // Search handles
	SearchHandlesName     cmap.ConcurrentMap
	SearchHandlesValue    cmap.ConcurrentMap
	SearchHandlesType     cmap.ConcurrentMap
	SearchHandlesTags     cmap.ConcurrentMap
	SearchHandlesIterator cmap.ConcurrentMap
	SearchHandleCounter   indyUtils.Counter
}

func (e *InMemoryStorage) Create(walletId string, config string, credentialsJson string, metadata string) (int, error) {
	_, ok := e.StoredMetadata.Get(walletId)
	if ok == true {
		return 203, errors.New(indyUtils.GetIndyError(203)) //203: "WalletAlreadyExistsError: Attempt to create wallet with name used for another exists wallet"
	}

	metadataW := Metadata{WalletId: walletId, Value: metadata}
	e.StoredMetadata.Set(walletId, metadataW)

	return 0, nil
}

func (e *InMemoryStorage) Open(walletId string, storageConfig string, credentialsJson string) (int, int, error) {
	metadata, ok := e.StoredMetadata.Get(walletId)
	if !ok {
		return 0, 210, errors.New(indyUtils.GetIndyError(210)) //210: "WalletStorageError: Storage error occurred during wallet operation"
	}

	nextWalletHandle, handleKey := e.WalletHandlesCounter.Get()

	e.WalletHandles.Set(handleKey, walletId)
	e.MetadataHandles.Set(handleKey, metadata)

	e.StoredRecords.Set(walletId, cmap.New())
	e.StoredTagsE.Set(walletId, cmap.New())
	e.StoredTagsP.Set(walletId, cmap.New())

	return int(nextWalletHandle), 0, nil
}

func (e *InMemoryStorage) Close(walletHandle int) error {
	handleKey := strconv.Itoa(walletHandle)
	_, ok := e.WalletHandles.Get(handleKey)
	if !ok {
		return errors.New(indyUtils.GetIndyError(200)) //200: "WalletInvalidHandle: Caller passed invalid wallet handle"
	}

	keyStorageHandle := strconv.Itoa(walletHandle)
	e.WalletHandles.Remove(keyStorageHandle)

	return nil
}

func (e *InMemoryStorage) Delete(walletId string, storageConfig string, credentialsJson string) (int, error) {
	_, ok := e.StoredMetadata.Get(walletId)
	if !ok {
		return 210, errors.New(indyUtils.GetIndyError(210)) //210: "WalletStorageError: Storage error occurred during wallet operation"
	}

	e.StoredMetadata.Remove(walletId)
	e.StoredRecords.Remove(walletId)
	e.StoredTagsP.Remove(walletId)
	e.StoredTagsE.Remove(walletId)

	return 0, nil
}

func (e *InMemoryStorage) AddRecord(walletHandle int, recordType string, recordId string, recordValue []byte, tagsJson string) (int, error) {
	handleKey := strconv.Itoa(walletHandle)
	tmp, ok := e.WalletHandles.Get(handleKey)
	if !ok {
		return 200, errors.New(indyUtils.GetIndyError(200)) //200: "WalletInvalidHandle: Caller passed invalid wallet handle"
	}
	walletId := tmp.(string)

	record := StorageRecord{WalletId: walletId, Name: recordId, Type: recordType, Value: recordValue}

	rId, err := e.addRecord(walletId, record)
	if err != 0 {
		return err, errors.New(indyUtils.GetIndyError(err))
	}

	tags, errParse := gabs.ParseJSON([]byte(tagsJson))
	if errParse != nil {
		return 104, errors.New(indyUtils.GetIndyError(104)) //104: "CommonInvalidParam5: Caller passed invalid value as param 5 (null, invalid json and etc..)"
	}

	tagsMap := tags.ChildrenMap()
	for tag, child := range tagsMap {
		tagValue, okC := child.Data().(string)
		if !okC {
			return 104, errors.New(indyUtils.GetIndyError(104)) //104: "CommonInvalidParam5: Caller passed invalid value as param 5 (null, invalid json and etc..)",
		}

		if tag[0:1] == "~" { // plain
			tagP := TagsPlaintext{
				WalletId: walletId,
				Name:     tag,
				Value:    tagValue,
			}
			ok = e.addTagsP(walletId, rId, tagP)
			if !ok {
				return 210, errors.New(indyUtils.GetIndyError(210)) //210: "WalletStorageError: Storage error occurred during wallet operation"
			}
		} else { // encrypted
			tagE := TagsEncrypted{
				WalletId: walletId,
				Name:     tag,
				Value:    tagValue,
			}
			ok = e.addTagsE(walletId, rId, tagE)
			if !ok {
				return 210, errors.New(indyUtils.GetIndyError(210)) //210: "WalletStorageError: Storage error occurred during wallet operation"
			}
		}
	}

	return 0, nil
}

func (e *InMemoryStorage) UpdateRecordValue(walletHandle int, recordType string, recordId string, recordValue []byte) (int, error) {
	var record StorageRecord
	id := ""

	walletId, recordMap, _, _, err := e.getStorageMaps(walletHandle)
	if err != 0 {
		return err, errors.New(indyUtils.GetIndyError(err))
	}

	ok := false
	for t := range recordMap.IterBuffered() {
		sRecord := t.Val.(StorageRecord)

		if sRecord.Name == recordId && sRecord.Type == recordType {
			ok = true
			id = t.Key
			record = sRecord
			break
		}
	}

	if !ok {
		return 212, errors.New(indyUtils.GetIndyError(212)) // "WalletItemNotFound: Requested wallet item not found"
	}

	record.Value = recordValue
	recordMap.Set(id, record)
	e.StoredRecords.Set(walletId, recordMap)

	return 0, nil
}

func (e *InMemoryStorage) UpdateRecordTags(walletHandle int, recordType string, recordId string, tagsJson string) (int, error) {
	id := ""

	walletId, recordMap, tagsMapP, tagsMapE, err := e.getStorageMaps(walletHandle)
	if err != 0 {
		return err, errors.New(indyUtils.GetIndyError(err))
	}

	ok := false
	for t := range recordMap.IterBuffered() {
		sRecord := t.Val.(StorageRecord)

		if sRecord.Name == recordId && sRecord.Type == recordType {
			ok = true
			id = t.Key
			break
		}
	}

	if !ok {
		return 210, errors.New(indyUtils.GetIndyError(212)) //212: "WalletItemNotFound: Requested wallet item not found"
	}

	ok = e.rIdFromTags("", walletId, id, tagsMapP, tagsMapE)
	if !ok {
		return 210, errors.New(indyUtils.GetIndyError(212)) //212: "WalletItemNotFound: Requested wallet item not found"
	}

	if len(tagsJson) > 0 {
		tags, errParse := gabs.ParseJSON([]byte(tagsJson))
		if errParse != nil {
			return 103, errors.New(indyUtils.GetIndyError(103)) //103: "CommonInvalidParam4: Caller passed invalid value as param 4 (null, invalid json and etc..)",
		}

		tagsMap := tags.ChildrenMap()
		for tag, child := range tagsMap {
			tagValue, okC := child.Data().(string)
			if !okC {
				return 104, errors.New(indyUtils.GetIndyError(104)) //104: "CommonInvalidParam5: Caller passed invalid value as param 5 (null, invalid json and etc..)",
			}

			if tag[0:1] == "~" { // plain
				tagP := TagsPlaintext{
					WalletId: walletId,
					Name:     tag,
					Value:    tagValue,
				}
				ok = e.addTagsP(walletId, id, tagP)
				if !ok {
					return 210, errors.New(indyUtils.GetIndyError(210)) //210: "WalletStorageError: Storage error occurred during wallet operation"
				}
			} else { // encrypted
				tagE := TagsEncrypted{
					WalletId: walletId,
					Name:     tag,
					Value:    tagValue,
				}
				ok = e.addTagsE(walletId, id, tagE)
				if !ok {
					return 210, errors.New(indyUtils.GetIndyError(210)) //210: "WalletStorageError: Storage error occurred during wallet operation"
				}
			}
		}
	}

	return 0, nil
}

func (e *InMemoryStorage) AddRecordTags(walletHandle int, recordType string, recordId string, tagsJson string) (int, error) {
	id := ""

	walletId, recordMap, _, _, err := e.getStorageMaps(walletHandle)
	if err != 0 {
		return err, errors.New(indyUtils.GetIndyError(err))
	}

	ok := false
	for t := range recordMap.IterBuffered() {
		sRecord := t.Val.(StorageRecord)

		if sRecord.Name == recordId && sRecord.Type == recordType {
			ok = true
			id = t.Key
			break
		}
	}
	if !ok {
		return 210, errors.New(indyUtils.GetIndyError(212)) //212: "WalletItemNotFound: Requested wallet item not found"
	}

	if len(tagsJson) > 0 {
		tags, errParse := gabs.ParseJSON([]byte(tagsJson))
		if errParse != nil {
			return 103, errors.New(indyUtils.GetIndyError(103)) //103: "CommonInvalidParam4: Caller passed invalid value as param 4 (null, invalid json and etc..)",
		}

		tagsMap := tags.ChildrenMap()
		for tag, child := range tagsMap {
			tagValue, okC := child.Data().(string)
			if !okC {
				return 104, errors.New(indyUtils.GetIndyError(104)) //104: "CommonInvalidParam5: Caller passed invalid value as param 5 (null, invalid json and etc..)",
			}

			if tag[0:1] == "~" { // plain
				tagP := TagsPlaintext{
					WalletId: walletId,
					Name:     tag,
					Value:    tagValue,
				}
				ok = e.addTagsP(walletId, id, tagP)
				if !ok {
					return 210, errors.New(indyUtils.GetIndyError(210)) //210: "WalletStorageError: Storage error occurred during wallet operation"
				}
			} else { // encrypted
				tagE := TagsEncrypted{
					WalletId: walletId,
					Name:     tag,
					Value:    tagValue,
				}
				ok = e.addTagsE(walletId, id, tagE)
				if !ok {
					return 210, errors.New(indyUtils.GetIndyError(210)) //210: "WalletStorageError: Storage error occurred during wallet operation"
				}
			}
		}
	}

	return 0, nil
}

func (e *InMemoryStorage) DeleteRecordTags(walletHandle int, recordType string, recordId string, tagsJson string) (int, error) {
	id := ""

	walletId, recordMap, tagsMapP, tagsMapE, err := e.getStorageMaps(walletHandle)
	if err != 0 {
		return err, errors.New(indyUtils.GetIndyError(err))
	}

	ok := false
	for t := range recordMap.IterBuffered() {
		sRecord := t.Val.(StorageRecord)

		if sRecord.Name == recordId && sRecord.Type == recordType {
			ok = true
			id = t.Key
			break
		}
	}

	if !ok {
		return 210, errors.New(indyUtils.GetIndyError(212)) //212: "WalletItemNotFound: Requested wallet item not found"
	}

	if len(tagsJson) > 0 {
		tags, errParse := gabs.ParseJSON([]byte(tagsJson))
		if errParse != nil {
			return 103, errors.New(indyUtils.GetIndyError(103)) //103: "CommonInvalidParam4: Caller passed invalid value as param 4 (null, invalid json and etc..)",
		}

		for _, tagGabs := range tags.Children() {
			ok = e.rIdFromTags(tagGabs.Data().(string), walletId, id, tagsMapP, tagsMapE)
			if !ok {
				return 210, errors.New(indyUtils.GetIndyError(210)) //210: "WalletStorageError: Storage error occurred during wallet operation"
			}
		}
	}

	return 0, nil
}

func (e *InMemoryStorage) DeleteRecord(walletHandle int, recordType string, recordId string) (int, error) {
	id := ""

	walletId, recordMap, tagsMapP, tagsMapE, err := e.getStorageMaps(walletHandle)
	if err != 0 {
		return err, errors.New(indyUtils.GetIndyError(err))
	}

	ok := false
	for t := range recordMap.IterBuffered() {
		sRecord := t.Val.(StorageRecord)

		if sRecord.Name == recordId && sRecord.Type == recordType {
			ok = true
			id = t.Key
			break
		}
	}

	if !ok {
		return 210, errors.New(indyUtils.GetIndyError(212)) //212: "WalletItemNotFound: Requested wallet item not found"
	}

	ok = e.rIdFromTags("", walletId, id, tagsMapP, tagsMapE)
	if !ok {
		return 210, errors.New(indyUtils.GetIndyError(212)) //212: "WalletItemNotFound: Requested wallet item not found"
	}

	recordMap.Remove(id)
	e.StoredRecords.Set(walletId, recordMap)

	return 0, nil
}

func (e *InMemoryStorage) GetRecordHandle(walletHandle int, recordType string, recordId string, optionsJson string) (int, int, error) {
	var record StorageRecord

	_, recordMap, _, _, err := e.getStorageMaps(walletHandle)
	if err != 0 {
		return 0, err, errors.New(indyUtils.GetIndyError(err))
	}

	ok := false
	for t := range recordMap.IterBuffered() {
		sRecord := t.Val.(StorageRecord)

		if sRecord.Name == recordId && sRecord.Type == recordType {
			ok = true
			record = sRecord
			break
		}
	}
	if !ok {
		return 0, 210, errors.New(indyUtils.GetIndyError(212)) //212: "WalletItemNotFound: Requested wallet item not found"
	}

	sh, shKey := e.SearchHandleCounter.Get()
	e.SearchHandles.Set(shKey, record)
	e.SearchHandlesName.Set(shKey, C.CString(record.Name))
	e.SearchHandlesValue.Set(shKey, wallet.RecordValue{
		Len:   len(record.Value),
		Value: C.CBytes(record.Value),
	})
	e.SearchHandlesType.Set(shKey, C.CString(record.Type))

	return int(sh), 0, nil
}

func (e *InMemoryStorage) GetRecordId(walletHandle int, recordHandle int) (unsafe.Pointer, int, error) {
	searchHandleKey := strconv.Itoa(recordHandle)
	tmp, ok := e.SearchHandlesName.Get(searchHandleKey)
	if !ok {
		return nil, 200, errors.New(indyUtils.GetIndyError(200)) //200: "WalletInvalidHandle: Caller passed invalid wallet handle"
	}

	item, okCast := tmp.(*C.char)
	if !okCast {
		return nil, 208, errors.New(indyUtils.GetIndyError(208)) //208: WalletInputError: Input provided to wallet operations is considered not valid
	}

	return unsafe.Pointer(item), 0, nil
}

func (e *InMemoryStorage) GetRecordType(walletHandle int, recordHandle int) (unsafe.Pointer, int, error) {
	searchHandleKey := strconv.Itoa(recordHandle)
	tmp, ok := e.SearchHandlesType.Get(searchHandleKey)
	if !ok {
		return nil, 200, errors.New(indyUtils.GetIndyError(200)) //200: "WalletInvalidHandle: Caller passed invalid wallet handle"
	}

	item, okCast := tmp.(*C.char)
	if !okCast {
		return nil, 208, errors.New(indyUtils.GetIndyError(208)) //208: WalletInputError: Input provided to wallet operations is considered not valid
	}

	return unsafe.Pointer(item), 0, nil
}

func (e *InMemoryStorage) GetRecordValue(walletHandle int, recordHandle int) (wallet.RecordValue, int, error) {
	searchHandleKey := strconv.Itoa(recordHandle)
	tmp, ok := e.SearchHandlesValue.Get(searchHandleKey)
	if !ok {
		return wallet.RecordValue{}, 200, errors.New(indyUtils.GetIndyError(200)) //200: "WalletInvalidHandle: Caller passed invalid wallet handle"
	}

	recordValue := tmp.(wallet.RecordValue)

	return recordValue, 0, nil
}

func (e *InMemoryStorage) GetRecordTags(walletHandle int, recordHandle int) (unsafe.Pointer, int, error) {
	walletId, _, tagsMapP, tagsMapE, err := e.getStorageMaps(walletHandle)
	if err != 0 {
		return nil, err, errors.New(indyUtils.GetIndyError(err))
	}

	searchHandleKey := strconv.Itoa(recordHandle)
	tmp, ok := e.SearchHandles.Get(searchHandleKey)
	if !ok {
		return nil, 200, errors.New(indyUtils.GetIndyError(200)) //200: "WalletInvalidHandle: Caller passed invalid wallet handle"
	}

	item, okCast := tmp.(StorageRecord)
	if !okCast {
		return nil, 208, errors.New(indyUtils.GetIndyError(208)) //208: WalletInputError: Input provided to wallet operations is considered not valid
	}

	var tps []TagsPlaintext
	var tes []TagsEncrypted

	if !tagsMapP.IsEmpty() {
		tmp, ok = tagsMapP.Get(walletId)
		if !ok {
			return nil, 210, errors.New(indyUtils.GetIndyError(210)) //210: "WalletStorageError: Storage error occurred during wallet operation"
		}
		mapP := tmp.(map[TagsPlaintext][]string)

		for tag, list := range mapP {
			for _, id := range list {
				if id == item.Id {
					tps = append(tps, tag)
				}
			}
		}
	}

	if !tagsMapE.IsEmpty() {
		tmp, ok = tagsMapE.Get(walletId)
		if !ok {
			return nil, 210, errors.New(indyUtils.GetIndyError(210)) //210: "WalletStorageError: Storage error occurred during wallet operation"
		}
		mapE := tmp.(map[TagsEncrypted][]string)

		for tag, list := range mapE {
			for _, id := range list {
				if id == item.Id {
					tes = append(tes, tag)
				}
			}
		}
	}

	jsonObj := gabs.New()
	for _, tagP := range tps {
		jsonObj.Set(tagP.Value, tagP.Name)
	}

	for _, tagE := range tes {
		jsonObj.Set(tagE.Value, tagE.Name)
	}

	tags := jsonObj.String()
	upTags := C.CString(tags)
	e.SearchHandlesTags.Set(searchHandleKey, upTags)

	return unsafe.Pointer(upTags), 0, nil
}

func (e *InMemoryStorage) FreeRecord(walletHandle int, recordHandle int) error {
	searchHandleKey := strconv.Itoa(recordHandle)
	e.SearchHandles.Remove(searchHandleKey)

	pName, okName := e.SearchHandlesName.Get(searchHandleKey)
	if okName {
		C.free(unsafe.Pointer(pName.(*C.char)))
	}
	e.SearchHandlesName.Remove(searchHandleKey)

	pValue, okValue := e.SearchHandlesValue.Get(searchHandleKey)
	if okValue {
		rv, okCast := pValue.(wallet.RecordValue)
		if okCast {
			C.free(rv.Value)
		}
	}
	e.SearchHandlesValue.Remove(searchHandleKey)

	pType, okType := e.SearchHandlesType.Get(searchHandleKey)
	if okType {
		C.free(unsafe.Pointer(pType.(*C.char)))
	}
	e.SearchHandlesType.Remove(searchHandleKey)

	pTags, okTags := e.SearchHandlesTags.Get(searchHandleKey)
	if okTags {
		upTags, okCast := pTags.(unsafe.Pointer)
		if okCast {
			C.free(upTags)
		}
	}
	e.SearchHandlesTags.Remove(searchHandleKey)

	return nil
}

func (e *InMemoryStorage) GetStorageMetadata(walletHandle int) (unsafe.Pointer, int, int, error) {
	handleKey := strconv.Itoa(walletHandle)
	tmp, ok := e.MetadataHandles.Get(handleKey)
	if !ok {
		return nil, 0, 212, errors.New(indyUtils.GetIndyError(212)) //"WalletItemNotFound: Requested wallet item not found"
	}

	metadata, okM := tmp.(Metadata)
	if !okM {
		return nil, 0, 210, errors.New(indyUtils.GetIndyError(210)) //"WalletStorageError: Storage error occurred during wallet operation"
	}
	metadataValue := C.CString(metadata.Value)

	return unsafe.Pointer(metadataValue), walletHandle, 0, nil
}

func (e *InMemoryStorage) SetStorageMetadata(walletHandle int, metadata string) (int, error) {
	handleKey := strconv.Itoa(walletHandle)
	tmp, ok := e.MetadataHandles.Get(handleKey)
	if !ok {
		return 212, errors.New(indyUtils.GetIndyError(212)) //"WalletItemNotFound: Requested wallet item not found"
	}

	newMeta, okM := tmp.(Metadata)
	if !okM {
		return 210, errors.New(indyUtils.GetIndyError(210)) //"WalletStorageError: Storage error occurred during wallet operation"
	}

	newMeta.Value = metadata
	e.MetadataHandles.Set(handleKey, newMeta)

	return 0, nil
}

func (e *InMemoryStorage) FreeStorageMetadata(walletHandle int, metadataHandle int) error {
	searchHandleKey := strconv.Itoa(metadataHandle)

	pMetadata, ok := e.MetadataHandles.Get(searchHandleKey)
	if ok {
		C.free(unsafe.Pointer(C.CString(pMetadata.(Metadata).Value)))
	}
	e.MetadataHandles.Remove(searchHandleKey)

	return nil
}

func (e *InMemoryStorage) OpenSearch(walletHandle int, recordType string, queryJson string, optionsJson string) (int, int, error) {
	walletId, recordMap, tagsMapP, tagsMapE, err := e.getStorageMaps(walletHandle)
	if err != 0 {
		return 0, err, errors.New(indyUtils.GetIndyError(err))
	}

	var searchedRecords []StorageRecord
	tagListE := make(map[TagsEncrypted][]string)
	tagListP := make(map[TagsPlaintext][]string)

	if !tagsMapE.IsEmpty() {
		tmp, ok := tagsMapE.Get(walletId)
		if !ok {
			return 0, 210, errors.New(indyUtils.GetIndyError(210)) //210: "WalletStorageError: Storage error occurred during wallet operation"
		}
		tagListE = tmp.(map[TagsEncrypted][]string)
	}

	if !tagsMapP.IsEmpty() {
		tmp, ok := tagsMapP.Get(walletId)
		if !ok {
			return 0, 210, errors.New(indyUtils.GetIndyError(210)) //210: "WalletStorageError: Storage error occurred during wallet operation"
		}
		tagListP = tmp.(map[TagsPlaintext][]string)
	}

	query, errGabs := gabs.ParseJSON([]byte(queryJson))
	if errGabs != nil {
		return 0, 113, errors.New(indyUtils.GetIndyError(113)) //113: "CommonInvalidStructure: Object (json, config, key, credential and etc...) passed by library caller has invalid structure"
	}

	if query != nil {
		children := query.ChildrenMap()

		for key, subquery := range children {
			switch key {
			case "$and":
				andRecords, err := e.andSearchCase(walletId, subquery.Children(), recordMap, tagListE, tagListP)
				if err != 0 {
					return 0, err, errors.New(indyUtils.GetIndyError(err))
				}

				searchedRecords = append(searchedRecords, andRecords...)
			case "$or":
				orRecords, err := e.orSearchCase(walletId, subquery.Children(), recordMap, tagListE, tagListP)
				if err != 0 {
					return 0, err, errors.New(indyUtils.GetIndyError(err))
				}

				searchedRecords = append(searchedRecords, orRecords...)
			case "$not":
				notRecords, err := e.notSearchCase(walletId, subquery.ChildrenMap(), recordMap, tagListE, tagListP)
				if err != 0 {
					return 0, err, errors.New(indyUtils.GetIndyError(err))
				}

				for i := len(searchedRecords) - 1; i >= 0; i-- {
					for j := range notRecords {
						if searchedRecords[i].Name == notRecords[j].Name && searchedRecords[i].Type == notRecords[j].Type {
							searchedRecords = append(searchedRecords[:i], searchedRecords[i+1:]...)
						}
					}
				}
			default:
				var records []string
				sIds, ok := e.checkTags(walletId, subquery, tagListE, tagListP)
				if ok != 0 {
					return 0, ok, errors.New(indyUtils.GetIndyError(ok)) //212: "WalletItemNotFound: Requested wallet item not found"
				}

				records = append(records, sIds...)

				for i := range records {
					tmp, ok := recordMap.Get(records[i])
					if !ok {
						return 0, 210, errors.New(indyUtils.GetIndyError(210)) //210: "WalletStorageError: Storage error occurred during wallet operation"
					}

					record := tmp.([]StorageRecord)
					searchedRecords = append(searchedRecords, record...)
				}
			}
		}
	}

	sh, shKey := e.SearchHandleCounter.Get()
	e.SearchHandles.Set(shKey, searchedRecords)
	e.SearchHandlesIterator.Set(shKey, 0)

	return int(sh), 0, nil
}

func (e *InMemoryStorage) OpenSearchAll(walletHandle int) (int, int, error) {
	_, recordMap, _, _, err := e.getStorageMaps(walletHandle)
	if err != 0 {
		return 0, err, errors.New(indyUtils.GetIndyError(err))
	}

	var searchedRecords []StorageRecord

	for _, record := range recordMap.Items() {
		records := record.([]StorageRecord)
		searchedRecords = append(searchedRecords, records...)
	}

	sh, shKey := e.SearchHandleCounter.Get()
	e.SearchHandles.Set(shKey, searchedRecords)
	e.SearchHandlesIterator.Set(shKey, 0)

	return int(sh), 0, nil
}

func (e *InMemoryStorage) GetSearchTotalCount(walletHandle int, searchHandle int) (int, int, error) {
	searchHandleKey := strconv.Itoa(searchHandle)
	items, ok := e.SearchHandles.Get(searchHandleKey)
	if !ok {
		return 0, 210, errors.New(indyUtils.GetIndyError(210)) //"WalletStorageError: Storage error occurred during wallet operation"
	}

	records := items.([]StorageRecord)
	return len(records), 0, nil
}

func (e *InMemoryStorage) FetchSearchNext(walletHandle int, searchHandle int) (int, int, error) {
	searchHandleKey := strconv.Itoa(searchHandle)
	tmp, ok := e.SearchHandles.Get(searchHandleKey)
	if !ok {
		return 0, 208, errors.New(indyUtils.GetIndyError(208)) //WalletInputError: Input provided to wallet operations is considered not valid
	}

	records, okCast := tmp.([]StorageRecord)
	if !okCast {
		return 0, 208, errors.New(indyUtils.GetIndyError(208)) //WalletInputError: Input provided to wallet operations is considered not valid
	}

	tmp1, okI := e.SearchHandlesIterator.Get(searchHandleKey)
	if !okI {
		return 0, 208, errors.New(indyUtils.GetIndyError(208)) //WalletInputError: Input provided to wallet operations is considered not valid
	}
	counter := tmp1.(int)

	if counter >= len(records) {
		return 0, 212, errors.New(indyUtils.GetIndyError(212))
	}

	e.SearchHandlesIterator.Set(searchHandleKey, counter+1)
	item := records[counter]
	handleId, handleKey := e.SearchHandleCounter.Get()
	e.SearchHandles.Set(handleKey, item)
	e.SearchHandlesName.Set(handleKey, C.CString(item.Name))
	e.SearchHandlesValue.Set(handleKey, wallet.RecordValue{
		Len:   len(item.Value),
		Value: C.CBytes(item.Value),
	})
	e.SearchHandlesType.Set(handleKey, C.CString(item.Type))

	return int(handleId), 0, nil
}

func (e *InMemoryStorage) FreeSearch(walletHandle int, searchHandle int) error {
	searchHandleKey := strconv.Itoa(searchHandle)
	e.SearchHandles.Remove(searchHandleKey)
	e.SearchHandlesIterator.Remove(searchHandleKey)

	return nil
}
