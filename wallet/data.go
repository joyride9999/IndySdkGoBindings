/*
// ******************************************************************
// Purpose: Wrapper to call libindy, imports functions from indy_wallet.h
// Author:  alexandru.leonte@siemens.com
// Notes:
// Copyright (c): Siemens SRL
// This work is licensed under the terms of the Apache License Version 2.0.  See
// the LICENSE.txt file in the top-level directory.
// ******************************************************************
*/

package wallet

/*
#cgo CFLAGS: -I ../include

#cgo LDFLAGS: -L${SRCDIR}/../lib -lindy
#include <indy_core.h>
#include <stdlib.h>
*/
import "C"
import "unsafe"

// StorageConfig represents Indy wallet storage config
type StorageConfig struct {
	Path   string `json:"path"`
	Dsn    string `json:"dsn"`    // Used with custom pg storage
	LogSql int    `json:"logsql"` // Used with custom pg storage
}

// Config represents Indy wallet config
type Config struct {
	ID            string        `json:"id"`
	StorageType   string        `json:"storage_type"`
	StorageConfig StorageConfig `json:"storage_config"`
}

// Credential represents Indy wallet credential config
type Credential struct {
	Key                   string `json:"key"`
	Rekey                 string `json:"rekey,omitempty"`
	StorageCredentials    string `json:"storage_credentials"`
	KeyDerivationMethod   string `json:"key_derivation_method,omitempty"`
	ReKeyDerivationMethod string `json:"rekey_derivation_method,omitempty"`
}

// ExportConfig represents Indy wallet export config
type ExportConfig struct {
	Path                string `json:"path"`
	Key                 string `json:"key"`
	KeyDerivationMethod string `json:"key_derivation_method,omitempty"`
}

// ImportConfig represents Indy wallet import config
type ImportConfig ExportConfig

// InfoWallet helper struct to hold information about the wallet
type InfoWallet struct {
	Name     string
	Key      string
	Did      string
	Seed     string
	Verkey   string
	IdWallet string
	Handle   int
}

type RecordValue struct {
	Len   int
	Value unsafe.Pointer
}

type IWalletStorage interface {
	Create(name string, config string, credentialsJson string, metadata string) (int, error)
	Open(name string, config string, credentials string) (int, int, error)
	Close(handle int) error
	Delete(name string, config string, credentials string) (int, error)
	AddRecord(handle int, type_ string, id string, value []byte, tagsJson string) (int, error)
	UpdateRecordValue(handle int, type_ string, id string, value []byte) (int, error)
	UpdateRecordTags(handle int, type_ string, id string, tagsJson string) (int, error)
	AddRecordTags(handle int, type_ string, id string, tagsJson string) (int, error)
	DeleteRecordTags(handle int, type_ string, id string, tagsJson string) (int, error)
	DeleteRecord(handle int, type_ string, id string) (int, error)
	GetRecordHandle(handle int, type_ string, id string, optionsJson string) (int, int, error)
	GetRecordId(handle int, recordHandle int) (unsafe.Pointer, int, error)
	GetRecordType(handle int, recordHandle int) (unsafe.Pointer, int, error)
	GetRecordValue(handle int, recordHandle int) (RecordValue, int, error)
	GetRecordTags(handle int, recordHandle int) (unsafe.Pointer, int, error)
	FreeRecord(handle int, recordHandle int) error
	GetStorageMetadata(handle int) (unsafe.Pointer, int, int, error)
	SetStorageMetadata(handle int, metadata string) (int, error)
	FreeStorageMetadata(handle int, metadataHandle int) error
	OpenSearch(handle int, type_ string, query string, options string) (int, int, error)
	OpenSearchAll(handle int) (int, int, error)
	GetSearchTotalCount(handle int, searchHandle int) (int, int, error)
	FetchSearchNext(handle int, searchHandle int) (int, int, error)
	FreeSearch(handle int, searchHandle int) error
}
