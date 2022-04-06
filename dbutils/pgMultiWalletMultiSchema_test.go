/*
// ******************************************************************
// Purpose: unit testing
// Author:  alexandru.leonte@siemens.com
// Notes:
// Copyright (c): Siemens SRL
// This work is licensed under the terms of the Apache License Version 2.0.  See
// the LICENSE.txt file in the top-level directory.
// ******************************************************************
*/

package dbutils

import (
	"github.com/Jeffail/gabs/v2"
	"testing"
)

//TODO: write more tests
func TestWqlToSql(t *testing.T) {
	dsn := "host=localhost user=wallet password=siemens dbname=wallets port=5432 sslmode=disable"
	testStorage := NewPgMultiSchemaStorage()
	queryTest := `{
		"$and": [
			{
            	"~t1": "v1"
        	},
			{
            	"~t2": "v2"
        	},
			{
				"$and": [ {
					"~t3": "v3"
				}
				]
			}

    	],
		"~t4": "v4"
	}`

	js, errGabs := gabs.ParseJSON([]byte(queryTest))
	if errGabs != nil {
		t.Errorf("Cant parse = '%v'", errGabs)
	}

	//wql to sql
	db, _ := testStorage.OpenDb(dsn, "", 4)
	children := js.ChildrenMap()
	var qparams []interface{}
	q := testStorage.OperatorToSql(db, children, &qparams, "test")

	//do the select
	var items []ItemsDB
	var item ItemsDB
	tx := db.Scopes(AddSchemaTable("test", item.TableName())).Where("type = ?", "ittypes")
	tx.Where(q, qparams...).Find(&items)

}
