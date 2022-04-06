/*
// ******************************************************************
// Purpose: Wrapper to call libindy, imports functions from indy_nonsecrets.h
// Author:  alexandru.leonte@siemens.com
// Notes:
// Copyright (c): Siemens SRL
// This work is licensed under the terms of the Apache License Version 2.0.  See
// the LICENSE.txt file in the top-level directory.
// ******************************************************************
*/

package nonsecrets

/*
#cgo CFLAGS: -I ../include
#cgo LDFLAGS: -L${SRCDIR}/../lib -lindy
#include <indy_core.h>

typedef void (*cb_addWalletRecord)(indy_handle_t, indy_error_t);
extern void addWalletRecordCB(indy_handle_t, indy_error_t);

typedef void (*cb_addWalletRecordTags)(indy_handle_t, indy_error_t);
extern void addWalletRecordTagsCB(indy_handle_t, indy_error_t);

typedef void (*cb_getWalletRecord)(indy_handle_t, indy_error_t, char*);
extern void getWalletRecordCB(indy_handle_t, indy_error_t, char*);

typedef void (*cb_deleteWalletRecord)(indy_handle_t, indy_error_t);
extern void deleteWalletRecordCB(indy_handle_t, indy_error_t);

typedef void (*cb_deleteWalletRecordTags)(indy_handle_t, indy_error_t);
extern void deleteWalletRecordTagsCB(indy_handle_t, indy_error_t);

typedef void (*cb_updateWalletRecordValue)(indy_handle_t, indy_error_t);
extern void updateWalletRecordValueCB(indy_handle_t, indy_error_t);

typedef void (*cb_updateWalletRecordTags)(indy_handle_t, indy_error_t);
extern void updateWalletRecordTagsCB(indy_handle_t, indy_error_t);

typedef void (*cb_openWalletSearch)(indy_handle_t, indy_error_t, indy_handle_t);
extern void openWalletSearchCB(indy_handle_t, indy_error_t, indy_handle_t);

typedef void (*cb_fetchWalletSearchNextRecords)(indy_handle_t, indy_error_t, char*);
extern void fetchWalletSearchNextRecordsCB(indy_handle_t, indy_error_t, char*);

typedef void (*cb_closeWalletSearch)(indy_handle_t, indy_error_t);
extern void closeWalletSearchCB(indy_handle_t, indy_error_t);
*/
import "C"
import (
	"errors"
	"github.com/joyride9999/IndySdkGoBindings/indyUtils"
	"unsafe"
)

//export getWalletRecordCB
func getWalletRecordCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, recordJson *C.char) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle),
			indyUtils.IndyResult{Error: nil,
				Results: []interface{}{
					string(C.GoString(recordJson)),
				}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// IndyGetWalletRecord Create a new non-secret record in the wallet.
func IndyGetWalletRecord(wh int, recordType string, recordId string, options string) chan indyUtils.IndyResult {

	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	var optionalOptions string
	if len(options) > 0 {
		optionalOptions = options
	} else {
		optionalOptions = "{}"
	}

	/*
		Get an wallet record by id

		    :param wallet_handle: wallet handler (created by open_wallet).
		    :param type_: allows to separate different record types collections
		    :param id: the id of record
		    :param options_json:
		      {
		        retrieveType: (optional, false by default) Retrieve record type,
		        retrieveValue: (optional, true by default) Retrieve record value,
		        retrieveTags: (optional, true by default) Retrieve record tags
		      }
		    :return: wallet record json:
		     {
		       id: "Some id",
		       type: "Some type", // present only if retrieveType set to true
		       value: "Some value", // present only if retrieveValue set to true
		       tags: <tags json>, // present only if retrieveTags set to true
		     }
	*/

	res := C.indy_get_wallet_record(commandHandle,
		(C.indy_handle_t)(wh),
		C.CString(recordType),
		C.CString(recordId),
		C.CString(optionalOptions),
		(C.cb_getWalletRecord)(unsafe.Pointer(C.getWalletRecordCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export addWalletRecordCB
func addWalletRecordCB(commandHandle C.indy_handle_t, indyError C.indy_error_t) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// IndyAddWalletRecord Create a new non-secret record in the wallet.
func IndyAddWalletRecord(wh int, recordType string, recordId string, recordValue string, tagsJson string) chan indyUtils.IndyResult {

	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	var optionalTagsJson *C.char
	if len(tagsJson) > 0 {
		optionalTagsJson = C.CString(tagsJson)
	} else {
		optionalTagsJson = nil
	}
	/*
	   :param wallet_handle: wallet handler (created by open_wallet).
	   :param type_: allows to separate different record types collections
	   :param id_: the id of record
	   :param value: the value of record
	   :param tags_json: the record tags used for search and storing meta information as json:
	      {
	        "tagName1": <str>, // string tag (will be stored encrypted)
	        "tagName2": <str>, // string tag (will be stored encrypted)
	        "~tagName3": <str>, // string tag (will be stored un-encrypted)
	        "~tagName4": <str>, // string tag (will be stored un-encrypted)
	      }
	   :return: None

	*/

	res := C.indy_add_wallet_record(commandHandle,
		(C.indy_handle_t)(wh),
		C.CString(recordType),
		C.CString(recordId),
		C.CString(recordValue),
		optionalTagsJson,
		(C.cb_addWalletRecord)(unsafe.Pointer(C.addWalletRecordCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export addWalletRecordTagsCB
func addWalletRecordTagsCB(commandHandle C.indy_handle_t, indyError C.indy_error_t) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// IndyAddWalletRecordTags Add new tags to the wallet record.
func IndyAddWalletRecordTags(wh int, recordType string, recordId string, tagsJson string) chan indyUtils.IndyResult {

	// Prepare the call parameters.
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	var optionalTagsJson *C.char
	if len(tagsJson) > 0 {
		optionalTagsJson = C.CString(tagsJson)
	} else {
		optionalTagsJson = nil
	}

	/*
		Add new tags to the wallet record

	    :param wallet_handle: wallet handler (created by open_wallet).
	    :param type_: allows to separate different record types collections
	    :param id_: the id of record
	    :param tags_json: the record tags used for search and storing meta information as json:
	       {
	         "tagName1": <str>, // string tag (will be stored encrypted)
	         "tagName2": <str>, // string tag (will be stored encrypted)
	         "~tagName3": <str>, // string tag (will be stored un-encrypted)
	         "~tagName4": <str>, // string tag (will be stored un-encrypted)
	       }
	    :return: None
	 */

	// Call to indy function.
	res := C.indy_add_wallet_record_tags(commandHandle,
		(C.indy_handle_t)(wh),
		C.CString(recordType),
		C.CString(recordId),
		optionalTagsJson,
		(C.cb_addWalletRecordTags)(unsafe.Pointer(C.addWalletRecordTagsCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export deleteWalletRecordCB
func deleteWalletRecordCB(commandHandle C.indy_handle_t, indyError C.indy_error_t) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// IndyDeleteWalletRecord Delete an existing wallet record in the wallet.
func IndyDeleteWalletRecord(wh int, recordType string, recordId string) chan indyUtils.IndyResult {

	// Prepare the call parameters.
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
		 Delete an existing wallet record in the wallet

	    :param wallet_handle: wallet handler (created by open_wallet).
	    :param type_: allows to separate different record types collections
	    :param id_: the id of record
	    :return: None
	 */

	// Call to indy function.
	res := C.indy_delete_wallet_record(commandHandle,
		(C.indy_handle_t)(wh),
		C.CString(recordType),
		C.CString(recordId),
		(C.cb_deleteWalletRecord)(unsafe.Pointer(C.deleteWalletRecordCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export deleteWalletRecordTagsCB
func deleteWalletRecordTagsCB(commandHandle C.indy_handle_t, indyError C.indy_error_t) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// IndyDeleteWalletRecordTags Delete tags from the wallet record.
func IndyDeleteWalletRecordTags(wh int, recordType string, recordId string, tagNames string) chan indyUtils.IndyResult {

	// Prepare the call parameters.
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
		Delete tags from the wallet record

		:param wallet_handle: wallet handler (created by open_wallet).
	    :param type_: allows to separate different record types collections
	    :param id_: the id of record
	    :param tag_names_json: the list of tag names to remove from the record as json array: ["tagName1", "tagName2", ...]
	    :return: None
	 */

	// Call to indy function.
	res := C.indy_delete_wallet_record_tags(commandHandle,
		(C.indy_handle_t)(wh),
		C.CString(recordType),
		C.CString(recordId),
		C.CString(tagNames),
		(C.cb_deleteWalletRecordTags)(unsafe.Pointer(C.deleteWalletRecordTagsCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export updateWalletRecordValueCB
func updateWalletRecordValueCB(commandHandle C.indy_handle_t, indyError C.indy_error_t) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// IndyUpdateWalletRecordValue Update a non-secret wallet record value.
func IndyUpdateWalletRecordValue(wh int, recordType string, recordId string, recordValue string) chan indyUtils.IndyResult {

	// Prepare the call parameters.
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
		 Update a non-secret wallet record value

	    :param wallet_handle: wallet handler (created by open_wallet).
	    :param type_: allows to separate different record types collections
	    :param id_: the id of record
	    :param value: the value of record
	    :return: None
	 */

	// Call to indy function
	res := C.indy_update_wallet_record_value(commandHandle,
		(C.indy_handle_t)(wh),
		C.CString(recordType),
		C.CString(recordId),
		C.CString(recordValue),
		(C.cb_updateWalletRecordValue)(unsafe.Pointer(C.updateWalletRecordValueCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export updateWalletRecordTagsCB
func updateWalletRecordTagsCB(commandHandle C.indy_handle_t, indyError C.indy_error_t) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// IndyUpdateWalletRecordTags Update a non-secret wallet record value.
func IndyUpdateWalletRecordTags(wh int, recordType string, recordId string, recordTags string) chan indyUtils.IndyResult {

	// Prepare the call parameters.
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
	   Update a non-secret wallet record value

	   :param wallet_handle: wallet handler (created by open_wallet).
	   :param type_: allows to separate different record types collections
	   :param id_: the id of record
	   :param tags_json: the record tags used for search and storing meta information as json:
	      {
	        "tagName1": <str>, // string tag (will be stored encrypted)
	        "tagName2": <str>, // string tag (will be stored encrypted)
	        "~tagName3": <str>, // string tag (will be stored un-encrypted)
	        "~tagName4": <str>, // string tag (will be stored un-encrypted)
	      }
	   :return: None
	 */

	// Call to indy function
	res := C.indy_update_wallet_record_tags(commandHandle,
		(C.indy_handle_t)(wh),
		C.CString(recordType),
		C.CString(recordId),
		C.CString(recordTags),
		(C.cb_updateWalletRecordTags)(unsafe.Pointer(C.updateWalletRecordTagsCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export openWalletSearchCB
func openWalletSearchCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, searchHandle C.indy_handle_t) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil, Results: []interface{}{int(searchHandle)}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// IndyOpenWalletSearch Search for wallet records.
func IndyOpenWalletSearch(wh int, recordType string, query string, options string) chan indyUtils.IndyResult {

	// Prepare the call parameters.
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	var indyOptions *C.char
	if len(options) > 0 {
		indyOptions = C.CString(options)
	} else {
		indyOptions = C.CString("{}")
	}

	var indyQuery *C.char
	if len(query) > 0 {
		indyQuery = C.CString(query)
	} else {
		indyQuery = C.CString("{}")
	}

	/*
	    Search for wallet records

	    :param wallet_handle: wallet handler (created by open_wallet).
	    :param type_: allows to separate different record types collections
	    :param query_json: MongoDB style query to wallet record tags:
	      {
	        "tagName": "tagValue",
	        $or: {
	          "tagName2": { $regex: 'pattern' },
	          "tagName3": { $gte: '123' },
	        },
	      }
	    :param options_json: //TODO: FIXME: Think about replacing by bitmask
	      {
	        retrieveRecords: (optional, true by default) If false only "counts" will be calculated,
	        retrieveTotalCount: (optional, false by default) Calculate total count,
	        retrieveType: (optional, false by default) Retrieve record type,
	        retrieveValue: (optional, true by default) Retrieve record value,
	        retrieveTags: (optional, false by default) Retrieve record tags,
	      }
	    :return: search_handle: Wallet search handle that can be used later
	             to fetch records by small batches (with fetch_wallet_search_next_records)
	 */

	// Call to indy function
	res := C.indy_open_wallet_search(commandHandle,
		(C.indy_handle_t)(wh),
		C.CString(recordType),
		indyQuery,
		indyOptions,
		(C.cb_openWalletSearch)(unsafe.Pointer(C.openWalletSearchCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export fetchWalletSearchNextRecordsCB
func fetchWalletSearchNextRecordsCB(commandHandle C.indy_handle_t, indyError C.indy_error_t, recordsJson *C.char) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil, Results: []interface{}{string(C.GoString(recordsJson))}})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// IndyFetchWalletSearchNextRecords Fetch next records for wallet search.
func IndyFetchWalletSearchNextRecords(wh int, sh int, count int32) chan indyUtils.IndyResult {

	// Prepare the call parameters.
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
		   Fetch next records for wallet search.

		   :param wallet_handle: wallet handler (created by open_wallet)
		   :param wallet_search_handle: wallet wallet handle (created by open_wallet_search)
		   :param count: Count of records to fetch
		   :return: wallet records json:
			{
			  totalCount: <str>, // present only if retrieveTotalCount set to true
			  records: [{ // present only if retrieveRecords set to true
				  id: "Some id",
				  type: "Some type", // present only if retrieveType set to true
				  value: "Some value", // present only if retrieveValue set to true
				  tags: <tags json>, // present only if retrieveTags set to true
			  }],
			}
	 */

	// Call to indy function
	res := C.indy_fetch_wallet_search_next_records(commandHandle,
		(C.indy_handle_t)(wh),
		(C.indy_handle_t)(sh),
		(C.indy_u32_t)(count),
		(C.cb_fetchWalletSearchNextRecords)(unsafe.Pointer(C.fetchWalletSearchNextRecordsCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}

//export closeWalletSearchCB
func closeWalletSearchCB(commandHandle C.indy_handle_t, indyError C.indy_error_t) {
	if indyError == 0 {
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: nil})
	} else {
		errMsg := indyUtils.GetIndyError(int(indyError))
		indyUtils.RemoveFuture((int)(commandHandle), indyUtils.IndyResult{Error: errors.New(errMsg)})
	}
}

// IndyCloseWalletSearch Close wallet search (make search handle invalid).
func IndyCloseWalletSearch(sh int) chan indyUtils.IndyResult {

	// Prepare the call parameters.
	handle, future := indyUtils.NewFutureCommand()
	commandHandle := (C.indy_handle_t)(handle)

	/*
		Close wallet search (make search handle invalid)

		:param wallet_search_handle: wallet wallet handle (created by open_wallet_search)
	    :return: None
	 */

	// Call to indy function
	res := C.indy_close_wallet_search(commandHandle,
		(C.indy_handle_t)(sh),
		(C.cb_closeWalletSearch)(unsafe.Pointer(C.closeWalletSearchCB)))
	if res != 0 {
		errMsg := indyUtils.GetIndyError(int(res))
		go func() { indyUtils.RemoveFuture((int)(handle), indyUtils.IndyResult{Error: errors.New(errMsg)}) }()
		return future
	}

	return future
}
