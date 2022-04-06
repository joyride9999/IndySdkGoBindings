package indySDK

import (
	"fmt"
	"indySDK/wallet"
	"testing"
)

type StorageRecord struct {
	Id    string `json:"name,omitempty"`
	Type  string `json:"type,omitempty"`
	Value []byte `json:"value,omitempty"`
	Tags  string `json:"tags,omitempty"`
}

type InMemoryStorage struct {
	StoredMetadata  map[string]string
	StoredRecords   []StorageRecord
	MetadataHandles map[int]string
	StorageHandles  map[int]string
	RecordHandles   map[int]StorageRecord
}

//TODO: move this to separate folder
//TODO: test it more...looks buggy

func (e *InMemoryStorage) Create(storageName string, storageConfig string, credentialsJson string, metadata string) (int, error) {
	e.StoredMetadata[storageName] = metadata

	return 0, nil
}

func (e *InMemoryStorage) Open(storageName string, storageConfig string, credentialsJson string) (int, int, error) {
	nextStorageHandle := len(e.StorageHandles) + 1
	e.StorageHandles[nextStorageHandle] = storageName

	return nextStorageHandle, 0, nil
}

func (e *InMemoryStorage) Close(storageHandle int) error {
	delete(e.StorageHandles, storageHandle)

	return nil
}

func (e *InMemoryStorage) Delete(storageName string, storageConfig string, credentialsJson string) (int, error) {
	delete(e.StoredMetadata, storageName)

	return 0, nil
}

func (e *InMemoryStorage) AddRecord(storageHandle int, recordType string, recordId string, recordValue []byte, tagsJson string) (int, error) {
	var record StorageRecord
	record.Id = recordId
	record.Type = recordType
	record.Value = recordValue
	record.Tags = tagsJson

	e.StoredRecords = append(e.StoredRecords, record)

	return 0, nil
}

//TODO: this looks buggy
func (e *InMemoryStorage) UpdateRecordValue(storageHandle int, recordType string, recordId string, recordValue []byte) (int, error) {
	var record StorageRecord

	for i := range e.StoredRecords {
		storedRecord := e.StoredRecords[i]
		if storedRecord.Id == recordId && storedRecord.Type == recordType {
			record = storedRecord
			break
		}
	}
	record.Value = recordValue

	return 0, nil
}

func (e *InMemoryStorage) UpdateRecordTags(storageHandle int, recordType string, recordId string, tagsJson string) (int, error) {
	var record StorageRecord

	for i := range e.StoredRecords {
		storedRecord := e.StoredRecords[i]
		if storedRecord.Id == recordId && storedRecord.Type == recordType {
			record = storedRecord
			break
		}
	}
	record.Tags = tagsJson

	return 0, nil
}

func (e *InMemoryStorage) AddRecordTags(storageHandle int, recordType string, recordId string, tagsJson string) (int, error) {
	e.UpdateRecordTags(storageHandle, recordType, recordId, tagsJson)

	return 0, nil
}

func (e *InMemoryStorage) DeleteRecordTags(storageHandle int, recordType string, recordId string, tagsJson string) (int, error) {
	tagsJson = ""
	e.UpdateRecordTags(storageHandle, recordType, recordId, tagsJson)

	return 0, nil
}

func (e *InMemoryStorage) DeleteRecord(storageHandle int, recordType string, recordId string) (int, error) {
	var index int

	for i := range e.StoredRecords {
		storedRecord := e.StoredRecords[i]
		if storedRecord.Id == recordId && storedRecord.Type == recordType {
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
	record.Id = recordId
	record.Type = recordType
	result := 0
	//TODO: recheck this
	for i := range e.StoredRecords {
		storedRecord := e.StoredRecords[i]
		if storedRecord.Id == recordId && storedRecord.Type == recordType {
			record = storedRecord
			result = 1
			break
		}
	}
	if result == 0 {
		return 0, 0, nil
	} else {
		var nextRecordHandle = len(e.RecordHandles) + 1
		e.RecordHandles[nextRecordHandle] = record

		return nextRecordHandle, 0, nil
	}
}

func (e *InMemoryStorage) GetRecordId(storageHandle int, recordHandle int) (string, int, error) {
	recordId := e.RecordHandles[recordHandle].Id

	return recordId, 0, nil
}

func (e *InMemoryStorage) GetRecordType(storageHandle int, recordHandle int) (string, int, error) {
	recordType := e.RecordHandles[recordHandle].Type

	return recordType, 0, nil
}

func (e *InMemoryStorage) GetRecordValue(storageHandle int, recordHandle int) ([]byte, int, error) {
	storageValue := e.RecordHandles[recordHandle].Value

	return storageValue, 0, nil
}

func (e *InMemoryStorage) GetRecordTags(storageHandle int, recordHandle int) (string, int, error) {
	recordTags := e.RecordHandles[recordHandle].Tags

	return recordTags, 0, nil
}

func (e *InMemoryStorage) FreeRecord(storageHandle int, recordHandle int) error {
	delete(e.RecordHandles, recordHandle)

	return nil
}

func (e *InMemoryStorage) GetStorageMetadata(storageHandle int) (string, int, int, error) {
	metadata := e.StoredMetadata[e.StorageHandles[storageHandle]]

	nextMetadataHandle := len(e.MetadataHandles) + 1
	e.MetadataHandles[nextMetadataHandle] = metadata

	return metadata, nextMetadataHandle, 0, nil
}

func (e *InMemoryStorage) SetStorageMetadata(storageHandle int, metadata string) (int, error) {
	e.StoredMetadata[e.StorageHandles[storageHandle]] = metadata

	return 0, nil
}

func (e *InMemoryStorage) FreeStorageMetadata(storageHandle int, metadataHandle int) error {
	delete(e.MetadataHandles, metadataHandle)

	return nil
}

func (e *InMemoryStorage) OpenSearch(storageHandle int, recordType string, queryJson string, optionsJson string) (int, int, error) {
	return 0, 0, nil
}

func (e *InMemoryStorage) OpenSearchAll(storageHandle int) (int, int, error) {
	return 0, 0, nil
}

func (e *InMemoryStorage) GetSearchTotalCount(storageHandle int, searchHandle int) (int, int, error) {
	return 0, 0, nil
}

func (e *InMemoryStorage) FetchSearchNext(storageHandle int, searchHandle int) (int, int, error) {
	return 0, 0, nil
}

func (e *InMemoryStorage) FreeSearch(storageHandle int, searchHandle int) error {
	return nil
}

func TestRegisterWallet(t *testing.T) {
	// Prepare and run test cases
	type args struct {
		walletCfg  wallet.Config
		walletCred wallet.Credential
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"register-wallet-type",
			args{walletCfg: wallet.Config{
				ID:            "customWallet",
				StorageType:   "customWallet",
				StorageConfig: wallet.StorageConfig{Path: ".\\out\\wallet"},
			},
				walletCred: wallet.Credential{
					Key: "123",
				}},
			false},
		{"register-wallet-type-default",
			args{walletCfg: wallet.Config{
				ID:            "defaultWallet",
				StorageType:   "default",
				StorageConfig: wallet.StorageConfig{Path: ".\\out\\wallet"},
			},
				walletCred: wallet.Credential{
					Key: "123",
				}},
			false},
	}
	fmt.Println("Test")

	testStorage := new(InMemoryStorage)

	testStorage.StoredMetadata = make(map[string]string)
	testStorage.MetadataHandles = make(map[int]string)
	testStorage.StorageHandles = make(map[int]string)
	testStorage.RecordHandles = make(map[int]StorageRecord)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storageType := tt.args.walletCfg.StorageType
			if storageType != "default" {
				err := RegisterWalletStorage(storageType, testStorage)
				if err != nil {
					t.Errorf("RegisterWalletTypet() error = '%v'", err)
					return
				}
			}

			errCreate := CreateWallet(tt.args.walletCfg, tt.args.walletCred)
			if errCreate != nil != tt.wantErr {
				t.Errorf("CreateWallet() error = '%v', wantErr = '%v'", errCreate, tt.wantErr)
				return
			}

			handler, errOpen := OpenWallet(tt.args.walletCfg, tt.args.walletCred)
			if errOpen != nil != tt.wantErr {
				t.Errorf("OpenWallet() error = '%v', wantErr = '%v'", errOpen, tt.wantErr)
			}

			did, verkey, errDid := CreateAndStoreDID(handler, "")
			if errDid != nil {
				t.Errorf("CreateAndStoreDID() error = '%v'", errDid)
				return
			}
			fmt.Println(fmt.Sprintf("Wallet did: '%v' | verkey: '%v'", did, verkey))

			key, errKey := KeyForLocalDID(handler, did)
			if errKey != nil {
				t.Errorf("KeyForLicalDID() error = '%v'", errKey)
			}
			if key == verkey {
				fmt.Println("Local Key and Verkey are equal.")
			}

			errClose := CloseWallet(handler)
			if errClose != nil != tt.wantErr {
				t.Errorf("CloseWallet() error = '%v', wantErr = '%v'", errClose, tt.wantErr)
			}

			errDelete := DeleteWallet(tt.args.walletCfg, tt.args.walletCred)
			if errDelete != nil != tt.wantErr {
				t.Errorf("DeleteWallet() error = '%v', wantErr = '%v'", errDelete, tt.wantErr)
			}
		})
	}
}
