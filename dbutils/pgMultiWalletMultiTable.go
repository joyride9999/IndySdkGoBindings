/*
// ******************************************************************
// Purpose: Implements custom wallet storage using pg database
// Author:  alexandru.leonte@siemens.com
// Notes:
// Copyright (c): Siemens SRL
// This work is licensed under the terms of the Apache License Version 2.0.  See
// the LICENSE.txt file in the top-level directory.
// ******************************************************************
*/

package dbutils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"github.com/jackc/pgconn"
	cmap "github.com/orcaman/concurrent-map"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"github.com/joyride9999/IndySdkGoBindings/indyUtils"
	"github.com/joyride9999/IndySdkGoBindings/wallet"
	"regexp"
	"strconv"
)

func checkSchemaName(schema string) bool {
	//Check for allowed chars
	var AllowedChars = regexp.MustCompile(`^[a-zA-Z0-9\-_]+$`).MatchString
	if !AllowedChars(schema) {
		return false
	}
	return true
}

// AddSchemaTable Add schema prefix to the queries
func AddSchemaTable(schema string, table string) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		// bad schema name
		if !checkSchemaName(schema) {
			return nil
		}
		if len(schema) > 0 {
			tn := fmt.Sprintf("%s.%s", schema, table)
			return tx.Table(tn)
		}

		return tx.Table(table)
	}
}

func NewPgMultiSchemaStorage() *pgMultiSchemaStorage {
	storage := new(pgMultiSchemaStorage)
	storage.MetadataHandles = cmap.New()
	storage.StorageHandles = cmap.New()
	storage.HandlesToDb = cmap.New()
	storage.WalletIdsToDb = cmap.New()
	storage.SearchHandles = cmap.New()
	storage.SearchHandlesIterator = cmap.New()
	return storage
}

// Storage implementation for postgre
type pgMultiSchemaStorage struct {
	MetadataHandles       cmap.ConcurrentMap //handle to metadata
	StorageHandlesCounter indyUtils.Counter
	StorageHandles        cmap.ConcurrentMap // handle to walletid(pk) (int to string)
	WalletIdsToDb         cmap.ConcurrentMap //wallet id to db connection (string to *gorm.db)
	HandlesToDb           cmap.ConcurrentMap //storage handles to db connection(int to *gorm.db)
	SearchHandles         cmap.ConcurrentMap
	SearchHandlesIterator cmap.ConcurrentMap
	SearchHandleCounter   indyUtils.Counter
}

func (e *pgMultiSchemaStorage) GetWalletIdFromHandle(storageHandle int) (walletId string, er error) {
	defer func() {
		if r := recover(); r != nil {
			walletId = ""
			er = errors.New("can't get the wallet id")
		}
	}()
	handleKey := strconv.Itoa(storageHandle)
	tmp, ok := e.StorageHandles.Get(handleKey)

	if !ok {
		return "", errors.New("can't find the handle")
	}

	walletId = tmp.(string)
	return walletId, nil

}

func (e *pgMultiSchemaStorage) GetDbFromWalletId(walletId string) (db *gorm.DB, er error) {
	defer func() {
		if r := recover(); r != nil {
			db = nil
			er = errors.New("can't get the db handle")
		}
	}()

	tmp, ok := e.WalletIdsToDb.Get(walletId)

	if !ok {
		return nil, errors.New("can't find the handle")
	}

	db = tmp.(*gorm.DB)
	return db, nil
}

func (e *pgMultiSchemaStorage) GetDbFromHandle(storageHandle int) (db *gorm.DB, er error) {
	defer func() {
		if r := recover(); r != nil {
			db = nil
			er = errors.New("can't get the db handle")
		}
	}()
	handleKey := strconv.Itoa(storageHandle)
	tmp, ok := e.HandlesToDb.Get(handleKey)

	if !ok {
		return nil, errors.New("can't find the handle")
	}

	db = tmp.(*gorm.DB)
	return db, nil

}

func (e *pgMultiSchemaStorage) OpenDb(dsn string, walletID string, logLlv int) (*gorm.DB, error) {
	// Connect to database
	db, errOpenDB := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.LogLevel(logLlv)),
		NamingStrategy: schema.NamingStrategy{TablePrefix: walletID + ".",
			SingularTable: false},
	})
	if errOpenDB != nil {
		return nil, errOpenDB // WalletStorageError: Storage error occurred during wallet operation
	}

	return db, nil
}

func (e *pgMultiSchemaStorage) CloseDb(db *gorm.DB) {
	if db != nil {
		sqlDB, _ := db.DB()
		_ = sqlDB.Close()
	}
}

func (e *pgMultiSchemaStorage) CreateSchema(db *gorm.DB, schemaName string) error {

	if !checkSchemaName(schemaName) {
		return errors.New("characters not allowed")
	}

	sql := fmt.Sprintf("CREATE SCHEMA \"%s\"", schemaName)
	dbg := db.Exec(sql)
	if dbg.Error != nil {
		return dbg.Error
	}

	return nil
}

//Create - creates a wallet... this adds a record to the metadata table,
//if the wallet id does not exist already (wallet id is pk)
func (e *pgMultiSchemaStorage) Create(walletId string, storageConfig string, credentialsJson string, metadata string) (int, error) {

	storageCfg := wallet.StorageConfig{}
	errU := json.Unmarshal([]byte(storageConfig), &storageCfg)
	if errU != nil {
		return 210, errU // WalletStorageError: Storage error occurred during wallet operation
	}

	db, errOpenDb := e.OpenDb(storageCfg.Dsn, walletId, storageCfg.LogSql)
	if errOpenDb != nil {
		return 210, errOpenDb // WalletStorageError: Storage error occurred during wallet operation
	}
	defer e.CloseDb(db)
	tx := db.Begin()
	errSchema := e.CreateSchema(tx, walletId)
	if errSchema != nil {
		tx.Rollback()
		return 203, errSchema // WalletAlreadyExistsError
	}
	// auto-migrate GORM struct
	mtDb := MetadataDB{}
	errM := tx.Scopes(AddSchemaTable(walletId, mtDb.TableName())).AutoMigrate(mtDb)
	if errM != nil {
		tx.Rollback()
		return 210, errM // WalletStorageError: Storage error occurred during wallet operation
	}
	itDb := ItemsDB{}
	errM = tx.Scopes(AddSchemaTable(walletId, itDb.TableName())).AutoMigrate(itDb)
	if errM != nil {
		tx.Rollback()
		return 210, errM // WalletStorageError: Storage error occurred during wallet operation
	}

	teDb := TagsEncryptedDB{}
	errM = tx.Scopes(AddSchemaTable(walletId, teDb.TableName())).AutoMigrate(teDb)
	if errM != nil {
		tx.Rollback()
		return 210, errM // WalletStorageError: Storage error occurred during wallet operation
	}

	tpDb := TagsPlaintextDB{}
	errM = tx.Scopes(AddSchemaTable(walletId, tpDb.TableName())).AutoMigrate(tpDb)
	if errM != nil {
		tx.Rollback()
		return 210, errM // WalletStorageError: Storage error occurred during wallet operation
	}

	errCon := pgAddConstraintsMultiTableMultiSchema(tx, walletId)
	if errCon != nil {
		tx.Rollback()
		return 210, errM // WalletStorageError: Storage error occurred during wallet operation
	}

	mtDb = MetadataDB{WalletId: walletId, Value: metadata}
	txa := tx.Scopes(AddSchemaTable(walletId, mtDb.TableName())).Create(&mtDb)
	if txa.Error != nil {
		tx.Rollback()
		if txa.Error.(*pgconn.PgError).Code == "23505" {
			return 203, txa.Error // WalletAlreadyExistsError
		}

		return 210, txa.Error // WalletStorageError: Storage error occurred during wallet operation
	}
	errC := tx.Commit().Error
	if errC != nil {
		return 210, errC // WalletStorageError: Storage error occurred during wallet operation
	}

	return 0, nil
}

// Open - returns an internal handle to a wallet id. Wallet must exist in the database
func (e *pgMultiSchemaStorage) Open(walletId string, storageConfig string, credentialsJson string) (int, int, error) {
	var metadata MetadataDB
	var nCount int64

	storageCfg := wallet.StorageConfig{}
	errU := json.Unmarshal([]byte(storageConfig), &storageCfg)
	if errU != nil {
		return 0, 210, errU // WalletStorageError: Storage error occurred during wallet operation
	}

	db, errOpenDb := e.OpenDb(storageCfg.Dsn, walletId, storageCfg.LogSql)
	if errOpenDb != nil {
		return 0, 210, errOpenDb // WalletStorageError: Storage error occurred during wallet operation
	}

	tx := db.Scopes(AddSchemaTable(walletId, metadata.TableName())).Find(&metadata, "wallet_id = ?", walletId).Count(&nCount)
	if tx.Error != nil {
		return 0, 210, tx.Error //WalletStorageError: Storage error occurred during wallet operation
	}
	if nCount != 1 {
		return 0, 200, errors.New("wallet id doesnt exist") //"WalletInvalidHandle: Caller passed invalid wallet handle"
	}

	nextStorageHandle, handleKey := e.StorageHandlesCounter.Get()
	e.StorageHandles.Set(handleKey, walletId)
	e.WalletIdsToDb.Set(walletId, db)
	e.HandlesToDb.Set(handleKey, db)
	e.MetadataHandles.Set(handleKey, metadata)

	return int(nextStorageHandle), 0, nil
}

//Close - removes the internal handle from the cache map for a wallet id
func (e *pgMultiSchemaStorage) Close(storageHandle int) error {
	keyStorageHandle := strconv.Itoa(storageHandle)
	db, err := e.GetDbFromHandle(storageHandle)
	if err != nil {
		return err
	}
	e.CloseDb(db)

	e.StorageHandles.Remove(keyStorageHandle)
	e.HandlesToDb.Remove(keyStorageHandle)

	tmp, ok := e.StorageHandles.Get(keyStorageHandle)
	if !ok {
		return nil
	}
	walletId := tmp.(string)
	e.WalletIdsToDb.Remove(walletId)

	return nil
}

// Delete - a wallet from the database
// WARN - everything stored for this wallet will be deleted
func (e *pgMultiSchemaStorage) Delete(walletID string, storageConfig string, credentialsJson string) (int, error) {

	storageCfg := wallet.StorageConfig{}
	errU := json.Unmarshal([]byte(storageConfig), &storageCfg)
	if errU != nil {
		return 210, errU // WalletStorageError: Storage error occurred during wallet operation
	}
	db, errOpenDb := e.OpenDb(storageCfg.Dsn, walletID, storageCfg.LogSql)
	if errOpenDb != nil {
		return 210, errOpenDb // WalletStorageError: Storage error occurred during wallet operation
	}

	defer e.CloseDb(db)
	tx := db.Begin()

	var metaData MetadataDB
	var item ItemsDB
	var tagsE TagsEncryptedDB
	var tagsP TagsPlaintextDB

	dM := tx.Scopes(AddSchemaTable(walletID, metaData.TableName())).Where("wallet_id = ?", walletID).Delete(&metaData).Error
	if dM != nil {
		tx.Rollback()
		return 210, dM // WalletStorageError: Storage error occurred during wallet operation
	}

	dI := tx.Scopes(AddSchemaTable(walletID, item.TableName())).Where("wallet_id = ?", walletID).Delete(&item).Error
	if dI != nil {
		tx.Rollback()
		return 210, dI // WalletStorageError: Storage error occurred during wallet operation
	}

	dTE := tx.Scopes(AddSchemaTable(walletID, tagsE.TableName())).Where("wallet_id = ?", walletID).Delete(&tagsE).Error
	if dTE != nil {
		tx.Rollback()
		return 210, dTE // WalletStorageError: Storage error occurred during wallet operation
	}

	dTP := tx.Scopes(AddSchemaTable(walletID, tagsP.TableName())).Where("wallet_id = ?", walletID).Delete(&tagsP).Error
	if dTP != nil {
		tx.Rollback()
		return 210, dTP // WalletStorageError: Storage error occurred during wallet operation
	}

	if tx.Commit().Error != nil {
		return 210, dTP // WalletStorageError: Storage error occurred during wallet operation
	}

	return 0, nil
}

//AddRecord - adds encrypted data to items table.
func (e *pgMultiSchemaStorage) AddRecord(storageHandle int, recordType string, recordId string, recordValue []byte, tagsJson string) (int, error) {

	walletID, errW := e.GetWalletIdFromHandle(storageHandle)
	if errW != nil {
		return 200, errors.New("WalletInvalidHandle: Caller passed invalid wallet handle")
	}

	db, err := e.GetDbFromHandle(storageHandle)
	if err != nil {
		return 200, errors.New("WalletInvalidHandle: Caller passed invalid wallet handle")
	}

	// start transaction
	tx := db.Begin()
	itemDB := ItemsDB{
		WalletId: walletID,
		Type:     recordType,
		Name:     recordId,
		Value:    recordValue,
	}

	dbI := tx.Scopes(AddSchemaTable(walletID, itemDB.TableName())).Create(&itemDB)
	if dbI.Error != nil {
		tx.Rollback()
		return 210, dbI.Error //WalletStorageError: Storage error occurred during wallet operation
	}

	tags, errParse := gabs.ParseJSON([]byte(tagsJson))
	if errParse != nil {
		tx.Rollback()
		return 104, errors.New(indyUtils.GetIndyError(104)) //104: "CommonInvalidParam5: Caller passed invalid value as param 5 (null, invalid json and etc..)",
	}

	children := tags.ChildrenMap()
	for k, child := range children {
		var dbT *gorm.DB
		tagValue, okC := child.Data().(string)
		if !okC {
			tx.Rollback()
			return 104, errors.New(indyUtils.GetIndyError(104)) //104: "CommonInvalidParam5: Caller passed invalid value as param 5 (null, invalid json and etc..)",
		}
		if k[0:1] == "~" { // tags unencrypted
			tp := TagsPlaintextDB{
				WalletId: walletID,
				Name:     k,
				Value:    tagValue,
				ItemId:   itemDB.Id,
			}
			dbT = tx.Scopes(AddSchemaTable(walletID, tp.TableName())).Create(&tp)
		} else { // tags encrypted
			te := TagsEncryptedDB{
				WalletId: walletID,
				Name:     k,
				Value:    tagValue,
				ItemId:   itemDB.Id,
			}
			dbT = tx.Scopes(AddSchemaTable(walletID, te.TableName())).Create(&te)
		}

		if dbT.Error != nil {
			tx.Rollback()
			return 210, dbT.Error //WalletStorageError: Storage error occurred during wallet operation
		}
	}

	errC := tx.Commit().Error
	if errC != nil {
		return 210, errC //WalletStorageError: Storage error occurred during wallet operation
	}
	return 0, nil
}

//UpdateRecordValue - updates record value.
//TODO: test me!!!
func (e *pgMultiSchemaStorage) UpdateRecordValue(storageHandle int, recordType string, recordId string, recordValue []byte) (int, error) {

	walletID, errW := e.GetWalletIdFromHandle(storageHandle)
	if errW != nil {
		return 200, errors.New(indyUtils.GetIndyError(200)) //"WalletInvalidHandle: Caller passed invalid wallet handle"
	}

	db, err := e.GetDbFromHandle(storageHandle)
	if err != nil {
		return 200, errors.New(indyUtils.GetIndyError(200)) //"WalletInvalidHandle: Caller passed invalid wallet handle"
	}

	//start transaction
	tx := db.Begin()

	item := ItemsDB{}
	tu := tx.Scopes(AddSchemaTable(walletID, item.TableName())).Where("name = ? AND type = ?", recordId, recordType).Update("value", recordValue)
	if tu.Error != nil {
		tx.Rollback()
		return 210, tu.Error // storage error
	}
	switch tu.RowsAffected {
	case 0:
		tx.Rollback()
		return 212, errors.New(indyUtils.GetIndyError(212)) // "WalletItemNotFound: Requested wallet item not found"
	case 1:
		tx.Commit()
		return 0, nil
	default:
		tx.Rollback()
		return 112, errors.New(indyUtils.GetIndyError(112)) // 	112: "CommonInvalidState: Invalid library state was detected in runtime. It signals library bug",
	}

}

//UpdateRecordTags - updates tags
//TODO: test me
func (e *pgMultiSchemaStorage) UpdateRecordTags(storageHandle int, recordType string, recordId string, tagsJson string) (int, error) {
	walletId, errW := e.GetWalletIdFromHandle(storageHandle)
	if errW != nil {
		return 200, errors.New(indyUtils.GetIndyError(200)) //"WalletInvalidHandle: Caller passed invalid wallet handle"
	}

	db, err := e.GetDbFromHandle(storageHandle)
	if err != nil {
		return 200, errors.New(indyUtils.GetIndyError(200)) //"WalletInvalidHandle: Caller passed invalid wallet handle"
	}

	//start transaction
	tx := db.Begin()
	item := ItemsDB{}

	tS := tx.Scopes(AddSchemaTable(walletId, item.TableName())).Find(&item, "wallet_id=? and type=? and name=?", walletId, recordType, recordId)
	if tS.Error != nil {
		tx.Rollback()
		return 210, tS.Error // storage error
	}

	if tS.RowsAffected != 1 {
		tx.Rollback()
		return 112, errors.New(indyUtils.GetIndyError(112)) // 	112: "CommonInvalidState: Invalid library state was detected in runtime. It signals library bug",
	}

	tagsE := TagsEncryptedDB{}
	tagsP := TagsPlaintextDB{}

	// delete old tags
	errDTE := tx.Scopes(AddSchemaTable(walletId, tagsE.TableName())).Where("item_id = ?", item.Id).Delete(&tagsE).Error
	if errDTE != nil {
		tx.Rollback()
		return 210, errDTE // storage error
	}

	errDTP := tx.Scopes(AddSchemaTable(walletId, tagsP.TableName())).Where("item_id = ?", item.Id).Delete(&tagsP).Error
	if errDTP != nil {
		tx.Rollback()
		return 210, errDTP // storage error
	}

	// New tags
	if len(tagsJson) > 0 {
		tags, errParse := gabs.ParseJSON([]byte(tagsJson))
		if errParse != nil {
			tx.Rollback()
			return 103, errors.New(indyUtils.GetIndyError(103)) //103: "CommonInvalidParam4: Caller passed invalid value as param 4 (null, invalid json and etc..)",
		}

		children := tags.ChildrenMap()
		for key, child := range children {
			var errTags error
			tagValue, okC := child.Data().(string)
			if !okC {
				tx.Rollback()
				return 103, errors.New(indyUtils.GetIndyError(103)) //103: "CommonInvalidParam4: Caller passed invalid value as param 4 (null, invalid json and etc..)",
			}
			if key[0:1] == "~" { // tags unencrypted
				tagsP = TagsPlaintextDB{
					WalletId: walletId,
					Name:     key,
					Value:    tagValue,
					ItemId:   item.Id,
				}
				errTags = tx.Scopes(AddSchemaTable(walletId, tagsP.TableName())).Create(&tagsP).Error
			} else { // tags encrypted
				tagsE = TagsEncryptedDB{
					WalletId: walletId,
					Name:     key,
					Value:    tagValue,
					ItemId:   item.Id,
				}
				errTags = tx.Scopes(AddSchemaTable(walletId, tagsE.TableName())).Create(&tagsE).Error
			}

			if errTags != nil {
				tx.Rollback()
				return 210, errTags //WalletStorageError: Storage error occurred during wallet operation
			}
		}
	}

	errC := tx.Commit().Error
	if errC != nil {
		return 210, errC //WalletStorageError: Storage error occurred during wallet operation
	}

	return 0, nil
}

//AddRecordTags - add a tag to a record.
//TODO: test me
func (e *pgMultiSchemaStorage) AddRecordTags(storageHandle int, recordType string, recordId string, tagsJson string) (int, error) {
	walletId, errW := e.GetWalletIdFromHandle(storageHandle)
	if errW != nil {
		return 200, errors.New(indyUtils.GetIndyError(200)) //"WalletInvalidHandle: Caller passed invalid wallet handle"
	}

	db, err := e.GetDbFromHandle(storageHandle)
	if err != nil {
		return 200, errors.New(indyUtils.GetIndyError(200)) //"WalletInvalidHandle: Caller passed invalid wallet handle"
	}

	//start transaction
	tx := db.Begin()
	item := ItemsDB{}

	tS := tx.Scopes(AddSchemaTable(walletId, item.TableName())).Find(&item, "wallet_id=? and type=? and name=?", walletId, recordType, recordId)
	if tS.Error != nil {
		tx.Rollback()
		return 210, tS.Error // storage error
	}

	if tS.RowsAffected != 1 {
		tx.Rollback()
		return 112, errors.New(indyUtils.GetIndyError(112)) // 	112: "CommonInvalidState: Invalid library state was detected in runtime. It signals library bug",
	}

	tagsE := TagsEncryptedDB{}
	tagsP := TagsPlaintextDB{}

	// New tags
	if len(tagsJson) > 0 {
		tags, errParse := gabs.ParseJSON([]byte(tagsJson))
		if errParse != nil {
			tx.Rollback()
			return 103, errors.New(indyUtils.GetIndyError(103)) //103: "CommonInvalidParam4: Caller passed invalid value as param 4 (null, invalid json and etc..)",
		}

		children := tags.ChildrenMap()
		for key, child := range children {
			var errTags error
			tagValue, okC := child.Data().(string)
			if !okC {
				tx.Rollback()
				return 103, errors.New(indyUtils.GetIndyError(103)) //103: "CommonInvalidParam4: Caller passed invalid value as param 4 (null, invalid json and etc..)",
			}
			if key[0:1] == "~" { // tags unencrypted
				tagsP = TagsPlaintextDB{
					WalletId: walletId,
					Name:     key,
					Value:    tagValue,
					ItemId:   item.Id,
				}
				errTags = tx.Scopes(AddSchemaTable(walletId, tagsP.TableName())).Create(&tagsP).Error
			} else { // tags encrypted
				tagsE = TagsEncryptedDB{
					WalletId: walletId,
					Name:     key,
					Value:    tagValue,
					ItemId:   item.Id,
				}
				errTags = tx.Scopes(AddSchemaTable(walletId, tagsE.TableName())).Create(&tagsE).Error
			}

			if errTags != nil {
				tx.Rollback()
				return 210, errTags //WalletStorageError: Storage error occurred during wallet operation
			}
		}
	}

	errC := tx.Commit().Error
	if errC != nil {
		return 210, errC //WalletStorageError: Storage error occurred during wallet operation
	}

	return 0, nil
}

//DeleteRecordTags - delete record tag
//TODO: test me
func (e *pgMultiSchemaStorage) DeleteRecordTags(storageHandle int, recordType string, recordId string, tagsJson string) (int, error) {
	walletId, errW := e.GetWalletIdFromHandle(storageHandle)
	if errW != nil {
		return 200, errors.New(indyUtils.GetIndyError(200)) //"WalletInvalidHandle: Caller passed invalid wallet handle"
	}

	db, err := e.GetDbFromHandle(storageHandle)
	if err != nil {
		return 200, errors.New(indyUtils.GetIndyError(200)) //"WalletInvalidHandle: Caller passed invalid wallet handle"
	}

	//start transaction
	tx := db.Begin()
	item := ItemsDB{}

	tS := tx.Scopes(AddSchemaTable(walletId, item.TableName())).
		Find(&item, "wallet_id=? and type=? and name=?", walletId, recordType, recordId)
	if tS.Error != nil {
		tx.Rollback()
		return 210, tS.Error // storage error
	}

	if tS.RowsAffected != 1 {
		tx.Rollback()
		return 112, errors.New(indyUtils.GetIndyError(112)) // 	112: "CommonInvalidState: Invalid library state was detected in runtime. It signals library bug",
	}

	tagsE := TagsEncryptedDB{}
	tagsP := TagsPlaintextDB{}

	// New tags
	if len(tagsJson) > 0 {
		tags, errParse := gabs.ParseJSON([]byte(tagsJson))
		if errParse != nil {
			tx.Rollback()
			return 103, errors.New(indyUtils.GetIndyError(103)) //103: "CommonInvalidParam4: Caller passed invalid value as param 4 (null, invalid json and etc..)",
		}

		children := tags.ChildrenMap()
		for key, child := range children {
			var errTags error
			tagValue, okC := child.Data().(string)
			if !okC {
				tx.Rollback()
				return 103, errors.New(indyUtils.GetIndyError(103)) //103: "CommonInvalidParam4: Caller passed invalid value as param 4 (null, invalid json and etc..)",
			}
			if key[0:1] == "~" { // tags unencrypted
				tagsP = TagsPlaintextDB{
					WalletId: walletId,
					Name:     key,
					Value:    tagValue,
					ItemId:   item.Id,
				}
				errTags = tx.Scopes(AddSchemaTable(walletId, tagsP.TableName())).Where("item_id = ? AND name = ?", item.Id, item.Key).Delete(&tagsP).Error
			} else { // tags encrypted
				tagsE = TagsEncryptedDB{
					WalletId: walletId,
					Name:     key,
					Value:    tagValue,
					ItemId:   item.Id,
				}
				errTags = tx.Scopes(AddSchemaTable(walletId, tagsE.TableName())).Where("item_id = ? AND name = ?", item.Id, item.Key).Delete(&tagsE).Error
			}

			if errTags != nil {
				tx.Rollback()
				return 210, errTags //WalletStorageError: Storage error occurred during wallet operation
			}
		}
	}

	errC := tx.Commit().Error
	if errC != nil {
		return 210, errC //WalletStorageError: Storage error occurred during wallet operation
	}

	return 0, nil
}

//DeleteRecord - deletes a record
//TODO: test me
func (e *pgMultiSchemaStorage) DeleteRecord(storageHandle int, recordType string, recordId string) (int, error) {
	walletId, errW := e.GetWalletIdFromHandle(storageHandle)
	if errW != nil {
		return 200, errors.New(indyUtils.GetIndyError(200)) //"WalletInvalidHandle: Caller passed invalid wallet handle"
	}

	db, err := e.GetDbFromHandle(storageHandle)
	if err != nil {
		return 200, errors.New(indyUtils.GetIndyError(200)) //"WalletInvalidHandle: Caller passed invalid wallet handle"
	}

	//start transaction
	tx := db.Begin()
	item := ItemsDB{}

	errDel := tx.Scopes(AddSchemaTable(walletId, item.TableName())).
		Where("wallet_id=? and type=? and name=?", walletId, recordType, recordId).Delete(&item).Error
	if errDel != nil {
		tx.Rollback()
		return 210, errDel // storage error
	}

	errC := tx.Commit().Error
	if errC != nil {
		return 210, errC //WalletStorageError: Storage error occurred during wallet operation
	}

	return 0, nil
}

//GetRecordHandle - gets a record from the items table if exists...else appropiate error
func (e *pgMultiSchemaStorage) GetRecordHandle(storageHandle int, recordType string, recordName string, optionsJson string) (recordId int, indyErrCode int, err error) {
	var item ItemsDB

	walletId, errW := e.GetWalletIdFromHandle(storageHandle)
	if errW != nil {
		return 0, 200, errors.New("WalletInvalidHandle: Caller passed invalid wallet handle")
	}

	db, err := e.GetDbFromHandle(storageHandle)
	if err != nil {
		return 0, 200, errors.New("WalletInvalidHandle: Caller passed invalid wallet handle")
	}

	var nCount int64
	tx := db.Scopes(AddSchemaTable(walletId, item.TableName())).Find(&item, "wallet_id=? and type=? and name=?", walletId, recordType, recordName).Count(&nCount)

	if tx.Error != nil {
		return 0, 210, tx.Error // WalletStorageError: Storage error occurred during wallet operation
	}
	if nCount == 0 {
		return 0, 212, errors.New("WalletItemNotFound: Requested wallet item not found")
	}

	sh, shKey := e.SearchHandleCounter.Get()
	e.SearchHandles.Set(shKey, item)

	return int(sh), 0, nil
}

//GetRecordId - get record id (item.name) ... not to be confused with row id( item.id)
func (e *pgMultiSchemaStorage) GetRecordId(storageHandle int, recordHandle int) (string, int, error) {

	searchHandleKey := strconv.Itoa(recordHandle)
	tmp, ok := e.SearchHandles.Get(searchHandleKey)
	if !ok {
		return "", 200, errors.New("WalletInvalidHandle: Caller passed invalid wallet handle")
	}

	item, okCast := tmp.(ItemsDB)
	if !okCast {
		return "", 208, errors.New(indyUtils.GetIndyError(208)) //WalletInputError: Input provided to wallet operations is considered not valid
	}

	return item.Name, 0, nil
}

//GetRecordType - get record type
func (e *pgMultiSchemaStorage) GetRecordType(storageHandle int, recordHandle int) (string, int, error) {

	searchHandleKey := strconv.Itoa(recordHandle)
	tmp, ok := e.SearchHandles.Get(searchHandleKey)
	if !ok {
		return "", 200, errors.New("WalletInvalidHandle: Caller passed invalid wallet handle")
	}

	item, okCast := tmp.(ItemsDB)
	if !okCast {
		return "", 208, errors.New(indyUtils.GetIndyError(208)) //WalletInputError: Input provided to wallet operations is considered not valid
	}

	return item.Type, 0, nil
}

//GetRecordValue - get record value
func (e *pgMultiSchemaStorage) GetRecordValue(storageHandle int, recordHandle int) ([]byte, int, error) {
	searchHandleKey := strconv.Itoa(recordHandle)
	tmp, ok := e.SearchHandles.Get(searchHandleKey)
	if !ok {
		return nil, 200, errors.New("WalletInvalidHandle: Caller passed invalid wallet handle")
	}

	item, okCast := tmp.(ItemsDB)
	if !okCast {
		return nil, 208, errors.New(indyUtils.GetIndyError(208)) //WalletInputError: Input provided to wallet operations is considered not valid
	}
	return item.Value, 0, nil
}

//GetRecordTags - get record tags.
func (e *pgMultiSchemaStorage) GetRecordTags(storageHandle int, recordHandle int) (string, int, error) {
	searchHandleKey := strconv.Itoa(recordHandle)
	tmp, ok := e.SearchHandles.Get(searchHandleKey)
	if !ok {
		return "", 200, errors.New("WalletInvalidHandle: Caller passed invalid wallet handle")
	}

	db, err := e.GetDbFromHandle(storageHandle)
	if err != nil {
		return "", 200, errors.New(indyUtils.GetIndyError(200)) //"WalletInvalidHandle: Caller passed invalid wallet handle"
	}

	item, okCast := tmp.(ItemsDB)
	if !okCast {
		return "", 208, errors.New(indyUtils.GetIndyError(208)) //WalletInputError: Input provided to wallet operations is considered not valid
	}

	var tp TagsPlaintextDB
	var te TagsEncryptedDB
	var tps []TagsPlaintextDB
	var tes []TagsEncryptedDB

	errTp := db.Scopes(AddSchemaTable(item.WalletId, tp.TableName())).Where("item_id=?", item.Id).Find(&tps).Error
	if errTp != nil {
		return "", 210, errTp // WalletStorageError: Storage error occurred during wallet operation
	}

	errTe := db.Scopes(AddSchemaTable(item.WalletId, te.TableName())).Where("item_id=?", item.Id).Find(&tes).Error
	if errTe != nil {
		return "", 210, errTe // WalletStorageError: Storage error occurred during wallet operation
	}

	jsonObj := gabs.New()
	for _, tagP := range tps {
		jsonObj.Set(tagP.Value, tagP.Name)
	}

	for _, tagE := range tes {
		jsonObj.Set(tagE.Value, tagE.Name)
	}

	tags := jsonObj.String()

	return tags, 0, nil
}

//FreeRecord - free search handle record
func (e *pgMultiSchemaStorage) FreeRecord(storageHandle int, recordHandle int) error {
	searchHandleKey := strconv.Itoa(recordHandle)
	e.SearchHandles.Remove(searchHandleKey)
	return nil
}

//GetStorageMetadata - gets metadata
func (e *pgMultiSchemaStorage) GetStorageMetadata(storageHandle int) (string, int, int, error) { // (string , int, indycode, error)

	handleKey := strconv.Itoa(storageHandle)
	tmp, ok := e.MetadataHandles.Get(handleKey)
	if !ok {
		return "", 0, 212, errors.New(indyUtils.GetIndyError(212)) //"WalletItemNotFound: Requested wallet item not found"
	}

	metadata, okM := tmp.(MetadataDB)
	if !okM {
		return "", 0, 210, errors.New(indyUtils.GetIndyError(210)) //"WalletStorageError: Storage error occurred during wallet operation"
	}

	return metadata.Value, storageHandle, 0, nil
}

//SetStorageMetadata - updates metadata for a wallet
//TODO: test me
func (e *pgMultiSchemaStorage) SetStorageMetadata(storageHandle int, metadata string) (int, error) {
	walletId, errW := e.GetWalletIdFromHandle(storageHandle)
	if errW != nil {
		return 200, errors.New(indyUtils.GetIndyError(200)) //"WalletInvalidHandle: Caller passed invalid wallet handle"
	}
	db, err := e.GetDbFromHandle(storageHandle)
	if err != nil {
		return 200, errors.New(indyUtils.GetIndyError(200)) //"WalletInvalidHandle: Caller passed invalid wallet handle"
	}

	mt := MetadataDB{}
	tx := db.Scopes(AddSchemaTable(walletId, mt.TableName())).Where("wallet_id = ?", walletId).Update("value", metadata)
	if tx.Error != nil {
		return 210, tx.Error ////"WalletStorageError: Storage error occurred during wallet operation"
	}

	return 0, nil
}

//FreeStorageMetadata - frees storage ...
func (e *pgMultiSchemaStorage) FreeStorageMetadata(storageHandle int, metadataHandle int) error {
	searchHandleKey := strconv.Itoa(metadataHandle)
	e.MetadataHandles.Remove(searchHandleKey)

	return nil
}

//OpenSearch - search handle
func (e *pgMultiSchemaStorage) OpenSearch(storageHandle int, recordType string, queryJson string, optionsJson string) (int, int, error) {

	walletId, errW := e.GetWalletIdFromHandle(storageHandle)
	if errW != nil {
		return 0, 200, errors.New(indyUtils.GetIndyError(200)) //"WalletInvalidHandle: Caller passed invalid wallet handle"
	}

	db, err := e.GetDbFromHandle(storageHandle)
	if err != nil {
		return 0, 200, errors.New(indyUtils.GetIndyError(200)) //"WalletInvalidHandle: Caller passed invalid wallet handle"
	}

	wqlQuery, errGabs := gabs.ParseJSON([]byte(queryJson))
	if errGabs != nil {
		return 0, 113, errors.New(indyUtils.GetIndyError(113))
	}

	children := wqlQuery.ChildrenMap()
	var qparams []interface{}
	sqlClause := e.OperatorToSql(db, children, &qparams, walletId)

	//do the select
	var items []ItemsDB
	var item ItemsDB
	tx := db.Scopes(AddSchemaTable(walletId, item.TableName())).Where("type = ?", recordType)
	errF := tx.Where(sqlClause, qparams...).Find(&items).Error
	if errF != nil {
		return 0, 208, errF //WalletInputError: Input provided to wallet operations is considered not valid
	}

	sh, shKey := e.SearchHandleCounter.Get()
	e.SearchHandles.Set(shKey, items)
	e.SearchHandlesIterator.Set(shKey, 0)

	return int(sh), 0, nil
}

//OpenSearchAll - search handle
//TODO: testme
func (e *pgMultiSchemaStorage) OpenSearchAll(storageHandle int) (int, int, error) {
	walletId, errW := e.GetWalletIdFromHandle(storageHandle)
	if errW != nil {
		return 0, 200, errors.New(indyUtils.GetIndyError(200)) //"WalletInvalidHandle: Caller passed invalid wallet handle"
	}

	db, err := e.GetDbFromHandle(storageHandle)
	if err != nil {
		return 0, 200, errors.New(indyUtils.GetIndyError(200)) //"WalletInvalidHandle: Caller passed invalid wallet handle"
	}

	//do the select
	var items []ItemsDB
	var item ItemsDB
	errF := db.Scopes(AddSchemaTable(walletId, item.TableName())).Where("wallet_id = ?", walletId).Find(&items).Error

	if errF != nil {
		return 0, 208, errF //WalletInputError: Input provided to wallet operations is considered not valid
	}

	sh, shKey := e.SearchHandleCounter.Get()
	e.SearchHandles.Set(shKey, items)
	e.SearchHandlesIterator.Set(shKey, 0)

	return int(sh), 0, nil
}

//GetSearchTotalCount - gets results count
func (e *pgMultiSchemaStorage) GetSearchTotalCount(storageHandle int, searchHandle int) (int, int, error) {
	searchHandleKey := strconv.Itoa(searchHandle)
	tmp, ok := e.SearchHandles.Get(searchHandleKey)

	if !ok {
		return 0, 208, errors.New(indyUtils.GetIndyError(208)) //WalletInputError: Input provided to wallet operations is considered not valid
	}

	items, okCast := tmp.([]ItemsDB)
	if !okCast {
		return 0, 208, errors.New(indyUtils.GetIndyError(208)) //WalletInputError: Input provided to wallet operations is considered not valid
	}

	return len(items), 0, nil
}

//FetchSearchNext - advance search counter
func (e *pgMultiSchemaStorage) FetchSearchNext(storageHandle int, searchHandle int) (int, int, error) {
	searchHandleKey := strconv.Itoa(searchHandle)
	tmp, ok := e.SearchHandles.Get(searchHandleKey)

	if !ok {
		return 0, 208, errors.New(indyUtils.GetIndyError(208)) //WalletInputError: Input provided to wallet operations is considered not valid
	}

	items, okCast := tmp.([]ItemsDB)
	if !okCast {
		return 0, 208, errors.New(indyUtils.GetIndyError(208)) //WalletInputError: Input provided to wallet operations is considered not valid
	}

	tmp1, okI := e.SearchHandlesIterator.Get(searchHandleKey)
	if !okI {
		return 0, 208, errors.New(indyUtils.GetIndyError(208)) //WalletInputError: Input provided to wallet operations is considered not valid
	}

	counter := tmp1.(int)

	//bounds check
	if counter >= len(items) {
		return 0, 212, errors.New(indyUtils.GetIndyError(212)) //"WalletItemNotFound: Requested wallet item not found"
	}
	e.SearchHandlesIterator.Set(searchHandleKey, counter+1)
	item := items[counter]
	handleId, handleKey := e.SearchHandleCounter.Get()
	e.SearchHandles.Set(handleKey, item)
	return int(handleId), 0, nil
}

//FreeSearch - close search handler
func (e *pgMultiSchemaStorage) FreeSearch(storageHandle int, searchHandle int) error {
	searchHandleKey := strconv.Itoa(searchHandle)
	e.SearchHandles.Remove(searchHandleKey)
	e.SearchHandlesIterator.Remove(searchHandleKey)
	return nil
}

// AndToSql - and WQL to sql
func (e *pgMultiSchemaStorage) AndToSql(db *gorm.DB, jsObject []*gabs.Container, qparams *[]interface{}, schema string, bOp bool) string {

	q := " ( "
	if bOp {
		q = " AND ( "
	}
	for key, child := range jsObject {

		q = q + e.OperatorToSql(db, child.ChildrenMap(), qparams, schema)
		if key < len(jsObject)-1 {
			q = q + " AND "
		}

	}
	q = q + " ) "

	return q
}

// OrToSql - or WQL to sql
func (e *pgMultiSchemaStorage) OrToSql(db *gorm.DB, jsObject []*gabs.Container, qparams *[]interface{}, schema string, bOp bool) string {

	q := " ( "
	if bOp {
		q = " AND ( "
	}
	for key, child := range jsObject {

		q = q + e.OperatorToSql(db, child.ChildrenMap(), qparams, schema)
		if key < len(jsObject)-1 {
			q = q + " OR "
		}

	}
	q = q + " ) "

	return q
}

// NotToSql - not WQL to sql
// TODO : testme!!!
func (e *pgMultiSchemaStorage) NotToSql(db *gorm.DB, jsObject []*gabs.Container, qparams *[]interface{}, schema string, bOp bool) string {

	q := " ( "
	if bOp {
		q = " AND ( "
	}
	for key, child := range jsObject {

		q = q + e.OperatorToSql(db, child.ChildrenMap(), qparams, schema)
		if key < len(jsObject)-1 {
			q = q + " Not "
		}

	}
	q = q + " ) "

	return q
}

// SubOperatorToSql - build sql  from wql
func (e *pgMultiSchemaStorage) SubOperatorToSql(db *gorm.DB, js *gabs.Container, tagName string, schema string, bOp bool) (string, *gorm.DB) {

	s := "(id in (?))"
	if bOp {
		s = "AND (id in (?))"
	}

	// get tags table
	var d *gorm.DB
	if tagName[0:1] == "~" { // tags unencrypted
		tp := TagsPlaintextDB{}
		d = db.Scopes(AddSchemaTable(schema, tp.TableName())).Select("item_id")
	} else { // tags encrypted
		te := TagsEncryptedDB{}
		//d = db.Table(te.TableName()).Select("item_id")
		d = db.Scopes(AddSchemaTable(schema, te.TableName())).Select("item_id")

	}

	// get subquery condition
	children := js.ChildrenMap()
	//equal clause
	if len(children) == 0 {
		value := js.Data().(string)
		return s, d.Where("name = ? AND value = ? ", tagName, value)
	}

	for key, child := range children {
		switch key {
		case "$neq":
			value := child.Data().(string)
			return s, d.Where("name = ? AND value != ? ", key, value)
		case "$gt":
			value := child.Data().(string)
			return s, d.Where("name = ? AND value > ?", key, value)
		case "$gte":
			value := child.Data().(string)
			return s, d.Where("name = ? AND value >= ?", key, value)
		case "$lt":
			value := child.Data().(string)
			return s, d.Where("name = ? AND value < ?", key, value)
		case "$lte":
			value := child.Data().(string)
			return s, d.Where("name = ? AND value <= ?", key, value)
		case "$like":
			value := child.Data().(string)
			return s, d.Where("name = ? AND value LIKE ?", key, value)
		case "$in":
			var inValue []string
			for _, vals := range child.Children() {
				inValue = append(inValue, vals.Data().(string))
			}
			return s, d.Where("name = ? AND value IN (?)", key, inValue)
		}
	}

	// TODO: should not be here
	return s, d
}

// OperatorToSql - build sql  from wql
func (e *pgMultiSchemaStorage) OperatorToSql(db *gorm.DB, jsObject map[string]*gabs.Container, qparams *[]interface{}, schema string) string {
	s := ""
	b := false
	for key, child := range jsObject {

		switch key {
		case "$and":
			s = s + e.AndToSql(db, child.Children(), qparams, schema, b)
		case "$or":
			s = s + e.OrToSql(db, child.Children(), qparams, schema, b)
		case "$not":
			s = s + e.NotToSql(db, child.Children(), qparams, schema, b)
		default:
			sq, d := e.SubOperatorToSql(db, child, key, schema, b)
			s = s + sq
			*qparams = append(*qparams, d)
		}

		b = true

	}
	return s
}
