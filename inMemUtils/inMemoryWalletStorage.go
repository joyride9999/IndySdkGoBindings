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

type StorageRecord struct {
	WalletId string `json:"wallet_id"`
	Id       string `json:"name,omitempty"`
	Type     string `json:"type,omitempty"`
	Value    []byte `json:"value,omitempty"`
	Tags     string `json:"tags,omitempty"`
}

type Metadata struct {
	WalletId string `json:"wallet_id"`
	Value    string `json:"value"`
}

func NewInMemoryStorage() *InMemoryStorage {
	customStorage := new(InMemoryStorage)
	customStorage.MetadataHandles = cmap.New()
	customStorage.StorageHandles = cmap.New()
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
	MetadataCounter       indyUtils.Counter
	StoredRecords         []StorageRecord    // Stored wallet records
	StorageHandles        cmap.ConcurrentMap // Storage handles
	StorageHandlesCounter indyUtils.Counter
	SearchHandles         cmap.ConcurrentMap // Search handles
	SearchHandlesName     cmap.ConcurrentMap
	SearchHandlesValue    cmap.ConcurrentMap
	SearchHandlesType     cmap.ConcurrentMap
	SearchHandlesTags     cmap.ConcurrentMap
	SearchHandlesIterator cmap.ConcurrentMap
	SearchHandleCounter   indyUtils.Counter
}

func (e *InMemoryStorage) Create(storageName string, storageConfig string, credentialsJson string, metadata string) (int, error) {
	_, isFound := e.MetadataHandles.Get(storageName)
	if isFound == true {
		return 203, errors.New("WalletAlreadyExistsError: Attempt to create wallet with name used for another exists wallet")
	}

	metadataW := Metadata{WalletId: storageName, Value: metadata}
	_, handleKey := e.MetadataCounter.Get()
	e.MetadataHandles.Set(handleKey, metadataW)

	return 0, nil
}

func (e *InMemoryStorage) Open(storageName string, storageConfig string, credentialsJson string) (int, int, error) {
	nextStorageHandle, handleKey := e.StorageHandlesCounter.Get()
	e.StorageHandles.Set(handleKey, storageName)

	return int(nextStorageHandle), 0, nil
}

func (e *InMemoryStorage) Close(storageHandle int) error {
	keyStorageHandle := strconv.Itoa(storageHandle)

	e.StorageHandles.Remove(keyStorageHandle)

	return nil
}

func (e *InMemoryStorage) Delete(storageName string, storageConfig string, credentialsJson string) (int, error) {
	e.MetadataHandles.Remove(storageName)
	e.StorageHandles.Remove(storageName)

	return 0, nil
}

func (e *InMemoryStorage) AddRecord(storageHandle int, recordType string, recordId string, recordValue []byte, tagsJson string) (int, error) {
	handleKey := strconv.Itoa(storageHandle)
	walletId, ok := e.StorageHandles.Get(handleKey)
	if !ok {
		return 200, errors.New("WalletInvalidHandle: Caller passed invalid wallet handle")
	}

	record := StorageRecord{WalletId: walletId.(string), Id: recordId, Type: recordType, Value: recordValue, Tags: tagsJson}
	e.StoredRecords = append(e.StoredRecords, record)

	return 0, nil
}

func (e *InMemoryStorage) UpdateRecordValue(storageHandle int, recordType string, recordId string, recordValue []byte) (int, error) {
	for i := range e.StoredRecords {
		if e.StoredRecords[i].Id == recordId && e.StoredRecords[i].Type == recordType {
			e.StoredRecords[i].Value = recordValue
			break
		}
	}

	return 0, nil
}

func (e *InMemoryStorage) UpdateRecordTags(storageHandle int, recordType string, recordId string, tagsJson string) (int, error) {
	for i := range e.StoredRecords {
		if e.StoredRecords[i].Id == recordId && e.StoredRecords[i].Type == recordType {
			e.StoredRecords[i].Tags = tagsJson
			break
		}
	}

	return 0, nil
}

func (e *InMemoryStorage) AddRecordTags(storageHandle int, recordType string, recordId string, tagsJson string) (int, error) {
	addTags, _ := gabs.ParseJSON([]byte(tagsJson))

	for i := range e.StoredRecords {
		if e.StoredRecords[i].Id == recordId && e.StoredRecords[i].Type == recordType {
			storedTags, _ := gabs.ParseJSON([]byte(e.StoredRecords[i].Tags))
			storedTags.Merge(addTags)

			e.StoredRecords[i].Tags = storedTags.String()
			break
		}
	}

	return 0, nil
}

func (e *InMemoryStorage) DeleteRecordTags(storageHandle int, recordType string, recordId string, tagsJson string) (int, error) {
	for i := range e.StoredRecords {
		if e.StoredRecords[i].Id == recordId && e.StoredRecords[i].Type == recordType {
			e.StoredRecords[i].Tags = ""
		}
	}

	return 0, nil
}

func (e *InMemoryStorage) DeleteRecord(storageHandle int, recordType string, recordId string) (int, error) {
	var index int

	for i := range e.StoredRecords {
		if e.StoredRecords[i].Id == recordId && e.StoredRecords[i].Type == recordType {
			index = i
			break
		}
	}

	e.StoredRecords[index] = e.StoredRecords[len(e.StoredRecords)-1]
	e.StoredRecords = e.StoredRecords[:len(e.StoredRecords)-1]

	return 0, nil
}

func (e *InMemoryStorage) GetRecordHandle(storageHandle int, recordType string, recordId string, optionsJson string) (int, int, error) {
	var record StorageRecord
	count := 0

	for i := range e.StoredRecords {
		if e.StoredRecords[i].Id == recordId && e.StoredRecords[i].Type == recordType {
			count++
			record = e.StoredRecords[i]
			break
		}
	}

	if count == 0 {
		return 0, 212, errors.New(indyUtils.GetIndyError(212))
	}

	sh, shKey := e.SearchHandleCounter.Get()
	e.SearchHandles.Set(shKey, record)
	e.SearchHandlesName.Set(shKey, C.CString(record.Id))
	e.SearchHandlesValue.Set(shKey, wallet.RecordValue{
		Len:   len(record.Value),
		Value: C.CBytes(record.Value),
	})
	e.SearchHandlesType.Set(shKey, C.CString(record.Type))
	e.SearchHandlesTags.Set(shKey, C.CString(record.Tags))

	return int(sh), 0, nil
}

func (e *InMemoryStorage) GetRecordId(storageHandle int, recordHandle int) (unsafe.Pointer, int, error) {
	searchHandleKey := strconv.Itoa(recordHandle)
	tmp, ok := e.SearchHandlesName.Get(searchHandleKey)
	if !ok {
		return nil, 200, errors.New("WalletInvalidHandle: Caller passed invalid wallet handle")
	}

	item, okCast := tmp.(*C.char)
	if !okCast {
		return nil, 208, errors.New(indyUtils.GetIndyError(208)) //WalletInputError: Input provided to wallet operations is considered not valid
	}

	return unsafe.Pointer(item), 0, nil
}

func (e *InMemoryStorage) GetRecordType(storageHandle int, recordHandle int) (unsafe.Pointer, int, error) {
	searchHandleKey := strconv.Itoa(recordHandle)
	tmp, ok := e.SearchHandlesType.Get(searchHandleKey)
	if !ok {
		return nil, 200, errors.New("WalletInvalidHandle: Caller passed invalid wallet handle")
	}

	item, okCast := tmp.(*C.char)
	if !okCast {
		return nil, 208, errors.New(indyUtils.GetIndyError(208)) //WalletInputError: Input provided to wallet operations is considered not valid
	}

	return unsafe.Pointer(item), 0, nil
}

func (e *InMemoryStorage) GetRecordValue(storageHandle int, recordHandle int) (wallet.RecordValue, int, error) {
	searchHandleKey := strconv.Itoa(recordHandle)
	tmp, ok := e.SearchHandlesValue.Get(searchHandleKey)
	if !ok {
		return wallet.RecordValue{}, 200, errors.New("WalletInvalidHandle: Caller passed invalid wallet handle")
	}

	recordValue := tmp.(wallet.RecordValue)

	return recordValue, 0, nil
}

func (e *InMemoryStorage) GetRecordTags(storageHandle int, recordHandle int) (unsafe.Pointer, int, error) {
	searchHandleKey := strconv.Itoa(recordHandle)
	tmp, ok := e.SearchHandlesTags.Get(searchHandleKey)
	if !ok {
		return nil, 200, errors.New("WalletInvalidHandle: Caller passed invalid wallet handle")
	}

	item, okCast := tmp.(*C.char)
	if !okCast {
		return nil, 208, errors.New(indyUtils.GetIndyError(208)) //WalletInputError: Input provided to wallet operations is considered not valid
	}

	return unsafe.Pointer(item), 0, nil
}

func (e *InMemoryStorage) FreeRecord(storageHandle int, recordHandle int) error {
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

func (e *InMemoryStorage) GetStorageMetadata(storageHandle int) (unsafe.Pointer, int, int, error) {
	handleKey := strconv.Itoa(storageHandle)
	tmp, ok := e.MetadataHandles.Get(handleKey)
	if !ok {
		return nil, 0, 212, errors.New(indyUtils.GetIndyError(212)) //"WalletItemNotFound: Requested wallet item not found"
	}

	metadata, okM := tmp.(Metadata)
	if !okM {
		return nil, 0, 210, errors.New(indyUtils.GetIndyError(210)) //"WalletStorageError: Storage error occurred during wallet operation"
	}
	metadataValue := C.CString(metadata.Value)

	return unsafe.Pointer(metadataValue), storageHandle, 0, nil
}

func (e *InMemoryStorage) SetStorageMetadata(storageHandle int, metadata string) (int, error) {
	handleKey := strconv.Itoa(storageHandle)
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

func (e *InMemoryStorage) FreeStorageMetadata(storageHandle int, metadataHandle int) error {
	searchHandleKey := strconv.Itoa(metadataHandle)

	pMetadata, ok := e.MetadataHandles.Get(searchHandleKey)
	if ok {
		C.free(unsafe.Pointer(C.CString(pMetadata.(Metadata).Value)))
	}
	e.MetadataHandles.Remove(searchHandleKey)

	return nil
}

func (e *InMemoryStorage) OpenSearch(storageHandle int, recordType string, queryJson string, optionsJson string) (int, int, error) {
	var searchedRecords []StorageRecord
	notFound := true

	wqlQuery, errGabs := gabs.ParseJSON([]byte(queryJson))
	if errGabs != nil {
		return 0, 113, errors.New(indyUtils.GetIndyError(113))
	}

	for index := 0; index < len(e.StoredRecords); index++ {
		if e.StoredRecords[index].Type == recordType {
			if wqlQuery != nil {
				tagsParsed, _ := gabs.ParseJSON([]byte(e.StoredRecords[index].Tags))
				if IsIncluded(wqlQuery, tagsParsed) {
					searchedRecords = append(searchedRecords, e.StoredRecords[index])
					notFound = false
				}
			} else {
				searchedRecords = append(searchedRecords, e.StoredRecords[index])
				notFound = false
			}
		}
	}

	if notFound == true {
		return 0, 208, errors.New(indyUtils.GetIndyError(208))
	} else {
		sh, shKey := e.SearchHandleCounter.Get()
		e.SearchHandles.Set(shKey, searchedRecords)
		e.SearchHandlesIterator.Set(shKey, 0)

		return int(sh), 0, nil
	}
}

func (e *InMemoryStorage) OpenSearchAll(storageHandle int) (int, int, error) {
	var searchedRecords []string

	for index := 0; index < len(e.StoredRecords); index++ {
		searchedRecords = append(searchedRecords, e.StoredRecords[index].Id)
	}

	sh, shKey := e.SearchHandleCounter.Get()
	e.SearchHandles.Set(shKey, searchedRecords)
	e.SearchHandlesIterator.Set(shKey, 0)

	return int(sh), 0, nil
}

func (e *InMemoryStorage) GetSearchTotalCount(storageHandle int, searchHandle int) (int, int, error) {
	searchHandleKey := strconv.Itoa(searchHandle)
	items, ok := e.SearchHandles.Get(searchHandleKey)
	if !ok {
		return 0, 210, errors.New(indyUtils.GetIndyError(210)) //"WalletStorageError: Storage error occurred during wallet operation"
	}

	records := items.([]StorageRecord)
	return len(records), 0, nil
}

func (e *InMemoryStorage) FetchSearchNext(storageHandle int, searchHandle int) (int, int, error) {
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
	e.SearchHandlesName.Set(handleKey, C.CString(item.Id))
	e.SearchHandlesValue.Set(handleKey, wallet.RecordValue{
		Len:   len(item.Value),
		Value: C.CBytes(item.Value),
	})
	e.SearchHandlesType.Set(handleKey, C.CString(item.Type))

	return int(handleId), 0, nil
}

func (e *InMemoryStorage) FreeSearch(storageHandle int, searchHandle int) error {
	searchHandleKey := strconv.Itoa(searchHandle)
	e.SearchHandles.Remove(searchHandleKey)
	e.SearchHandlesIterator.Remove(searchHandleKey)

	return nil
}
