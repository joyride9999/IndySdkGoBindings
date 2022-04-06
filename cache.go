/*
// ******************************************************************
// Purpose: exported public functions that handles cache functions
// from libindy
// Author:  angel.draghici@siemens.com
// Notes:
// Copyright (c): Siemens SRL
// This work is licensed under the terms of the Apache License Version 2.0.  See
// the LICENSE.txt file in the top-level directory.
// ******************************************************************
*/

package indySDK

import "indySDK/cache"

// GetCacheCredDef gets credential definition json data for specified credential definition id
func GetCacheCredDef(ph int, wh int, sdid string, id string, options string) (string, error) {
	channel := cache.GetCredDef(ph, wh, sdid, id, options)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// GetCacheSchema gets schema json data for specified schema id
func GetCacheSchema(ph int, wh int, sdid string, id string, options string) (string, error) {
	channel := cache.GetSchema(ph, wh, sdid, id, options)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// PurgeCredDefCache purge credential definition cache
func PurgeCredDefCache(wh int, options string) error {
	channel := cache.PurgeCredDefCache(wh, options)
	result := <-channel
	return result.Error
}

// PurgeSchemaCache Purge schema cache
func PurgeSchemaCache(wh int, options string) error {
	channel := cache.PurgeSchemaCache(wh, options)
	result := <-channel
	return result.Error
}
