/*
// ******************************************************************
// Purpose: Defines table for storing wallets into the database
// Author:  alexandru.leonte@siemens.com
// Notes:
// Copyright (c): Siemens SRL
// This work is licensed under the terms of the Apache License Version 2.0.  See
// the LICENSE.txt file in the top-level directory.
// ******************************************************************
*/

package dbutils

import (
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

// Tables definitions
type MetadataDB struct {
	WalletId string `gorm:"column:wallet_id;primary_key"`
	Value    string `gorm:"column:value;not null"` // this comes as a base64 encoded json
}

func (m *MetadataDB) TableName() string {
	return "metadata"
}

type ItemsDB struct {
	WalletId string `gorm:"column:wallet_id;primaryKey"`
	Id       int64  `gorm:"column:id;primaryKey;auto_increment:true"`
	Type     string `gorm:"column:type"`
	Name     string `gorm:"column:name"`
	Value    []byte `gorm:"column:value"`
	Key      []byte `gorm:"column:key"` // NOT USED ACTUALLY
}

func (i *ItemsDB) TableName() string {
	return "items"
}

type TagsEncryptedDB struct {
	WalletId string `gorm:"column:wallet_id;primaryKey"`
	Name     string `gorm:"column:name;primaryKey"`
	Value    string `gorm:"column:value"`
	ItemId   int64  `gorm:"column:item_id;primaryKey;auto_increment:false;not null"`
}

func (te *TagsEncryptedDB) TableName() string {
	return "tags_encrypted"
}

type TagsPlaintextDB struct {
	WalletId string `gorm:"column:wallet_id;primaryKey"`
	Name     string `gorm:"column:name;primaryKey;not null"`
	Value    string `gorm:"column:value"`
	ItemId   int64  `gorm:"column:item_id;primaryKey;auto_increment:false;not null"`
}

func (tp *TagsPlaintextDB) TableName() string {
	return "tags_plaintext"
}

// Set constrains on tables
func pgAddConstraintsMultiTableMultiSchema(db *gorm.DB, schema string) error {

	if !checkSchemaName(schema) {
		return errors.New("Characters not allowed in schema name")
	}

	sql := fmt.Sprintf("CREATE UNIQUE INDEX IF NOT EXISTS \"%s_ux_metadata_wallet_id_id\" ON \"%s\".\"metadata\"(wallet_id)", schema, schema)
	errA := db.Exec(sql).Error
	if errA != nil {
		return errA
	}

	sql = fmt.Sprintf("CREATE UNIQUE INDEX IF NOT EXISTS \"%s_ux_metadata_values\" ON \"%s\".\"metadata\"(wallet_id, value)", schema, schema)
	errA = db.Exec(sql).Error
	if errA != nil {
		return errA
	}

	sql = fmt.Sprintf("CREATE UNIQUE INDEX IF NOT EXISTS \"%s_ux_items_wallet_id_id\" ON \"%s\".\"items\"(wallet_id, id)", schema, schema)
	errA = db.Exec(sql).Error
	if errA != nil {
		return errA
	}

	sql = fmt.Sprintf("CREATE UNIQUE INDEX IF NOT EXISTS \"%s_ux_items_type_name\" ON \"%s\".\"items\"(wallet_id, type, name)", schema, schema)
	errA = db.Exec(sql).Error
	if errA != nil {
		return errA
	}

	// TODO : check if constraint already exist and dont add it anymore
	sql = fmt.Sprintf("ALTER TABLE \"%s\".items ADD CONSTRAINT \"%s_fk_items_metadata\" FOREIGN KEY (wallet_id) REFERENCES \"%s\".\"metadata\" (wallet_id) ON DELETE CASCADE ON UPDATE CASCADE", schema, schema, schema)
	errA = db.Exec(sql).Error
	if errA != nil && errA.(*pgconn.PgError).Code != "42710" {
		return errA
	}

	sql = fmt.Sprintf("CREATE INDEX IF NOT EXISTS \"%s_ix_tags_encrypted_name\" ON \"%s\".\"tags_encrypted\"(wallet_id, name)", schema, schema)
	errA = db.Exec(sql).Error
	if errA != nil {
		return errA
	}

	sql = fmt.Sprintf("CREATE INDEX IF NOT EXISTS \"%s_ix_tags_encrypted_value\" ON \"%s\".\"tags_encrypted\"(wallet_id, md5(value))", schema, schema)
	errA = db.Exec(sql).Error
	if errA != nil {
		return errA
	}

	sql = fmt.Sprintf("CREATE INDEX IF NOT EXISTS \"%s_ix_tags_encrypted_wallet_id_item_id\" ON \"%s\".\"tags_encrypted\"(wallet_id, item_id)", schema, schema)
	errA = db.Exec(sql).Error
	if errA != nil {
		return errA
	}

	// TODO : check if constraint already exist and dont add it anymore
	sql = fmt.Sprintf("ALTER TABLE \"%s\".\"tags_encrypted\" ADD CONSTRAINT \"%s_fk_tagse_items\" FOREIGN KEY (wallet_id, item_id) REFERENCES \"%s\".\"items\" (wallet_id,id) ON DELETE CASCADE ON UPDATE CASCADE", schema, schema, schema)
	errA = db.Exec(sql).Error
	if errA != nil && errA.(*pgconn.PgError).Code != "42710" {
		return errA
	}

	sql = fmt.Sprintf("CREATE INDEX IF NOT EXISTS \"%s_ix_tags_plaintext_name\" ON \"%s\".\"tags_plaintext\"(wallet_id, name)", schema, schema)
	errA = db.Exec(sql).Error
	if errA != nil {
		return errA
	}

	sql = fmt.Sprintf("CREATE INDEX IF NOT EXISTS \"%s_ix_tags_plaintext_value\" ON \"%s\".\"tags_plaintext\"(wallet_id, value)", schema, schema)
	errA = db.Exec(sql).Error
	if errA != nil {
		return errA
	}

	sql = fmt.Sprintf("CREATE INDEX IF NOT EXISTS \"%s_ix_tags_plaintext_wallet_id_item_id\" ON \"%s\".\"tags_plaintext\"(wallet_id, item_id)", schema, schema)
	errA = db.Exec(sql).Error
	if errA != nil {
		return errA
	}

	// TODO : check if constraint already exist and dont add it anymore
	sql = fmt.Sprintf("ALTER TABLE \"%s\".\"tags_plaintext\" ADD CONSTRAINT \"%s_fk_tagsp_items\" FOREIGN KEY (wallet_id, item_id) REFERENCES \"%s\".\"items\" (wallet_id,id) ON DELETE CASCADE ON UPDATE CASCADE", schema, schema, schema)
	errA = db.Exec(sql).Error
	if errA != nil && errA.(*pgconn.PgError).Code != "42710" {
		return errA
	}

	return nil
}
