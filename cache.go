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

/*
#include <stdlib.h>
*/
import "C"
import (
	"github.com/joyride9999/IndySdkGoBindings/cache"
	"unsafe"
)

// GetCacheCredDef gets credential definition json data for specified credential definition id
func GetCacheCredDef(ph int, wh int, sdid string, credDefId string, options string) (string, error) {

	upSDid := unsafe.Pointer(C.CString(sdid))
	defer C.free(upSDid)
	upCredDefId := unsafe.Pointer(C.CString(credDefId))
	defer C.free(upCredDefId)
	upOptions := unsafe.Pointer(C.CString(options))
	defer C.free(upOptions)

	channel := cache.GetCredDef(ph, wh, upSDid, upCredDefId, upOptions)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// GetCacheSchema gets schema json data for specified schema id
func GetCacheSchema(ph int, wh int, sdid string, schemaId string, options string) (string, error) {

	upSDid := unsafe.Pointer(C.CString(sdid))
	defer C.free(upSDid)
	upSchemaId := unsafe.Pointer(C.CString(schemaId))
	defer C.free(upSchemaId)
	upOptions := unsafe.Pointer(C.CString(options))
	defer C.free(upOptions)

	channel := cache.GetSchema(ph, wh, upSDid, upSchemaId, upOptions)
	result := <-channel
	if result.Error != nil {
		return "", result.Error
	}
	return result.Results[0].(string), result.Error
}

// PurgeCredDefCache purge credential definition cache
func PurgeCredDefCache(wh int, options string) error {

	upOptions := unsafe.Pointer(C.CString(options))
	defer C.free(upOptions)
	channel := cache.PurgeCredDefCache(wh, upOptions)
	result := <-channel
	return result.Error
}

// PurgeSchemaCache Purge schema cache
func PurgeSchemaCache(wh int, options string) error {
	upOptions := unsafe.Pointer(C.CString(options))
	defer C.free(upOptions)
	channel := cache.PurgeSchemaCache(wh, upOptions)
	result := <-channel
	return result.Error
}
