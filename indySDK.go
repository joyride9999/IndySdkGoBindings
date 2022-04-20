/*
// ******************************************************************
// Purpose: helper functions
// from libindy
// Author:  alexandru.leonte@siemens.com
// Notes:
// Copyright (c): Siemens SRL
// This work is licensed under the terms of the Apache License Version 2.0.  See
// the LICENSE.txt file in the top-level directory.
// ******************************************************************
*/

package indySDK

import "C"
import (
	"github.com/joyride9999/IndySdkGoBindings/blobstorage"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"io"
	"math"
	"math/big"
	"path/filepath"
	"strconv"
)

func jsonObjectToString(obj interface{}) string {
	credConfigB, _ := json.Marshal(obj)
	return string(credConfigB)
}

// GetRevRegDef - gets rev reg defs
func GetRevRegDef(poolHandle int, verifierDid string, revRegId string, timeStamp int64) (revRegDefJson string, revRegJson string, ts uint64, err error) {
	getRevocRegDefRequest, errGetRevRegDefReq := BuildGetRevRegDefRequest(verifierDid, revRegId)
	if errGetRevRegDefReq != nil {
		return "", "", 0, errGetRevRegDefReq
	}

	getRevocRegDefResponse, errGetRevocRegDefResponse := SubmitRequest(poolHandle, getRevocRegDefRequest)
	if errGetRevocRegDefResponse != nil {
		return "", "", 0, errGetRevocRegDefResponse
	}

	rvRegId, revRegDefJson, errParse := ParseGetRevocRegDefResponse(getRevocRegDefResponse)
	if errParse != nil {
		return "", "", 0, errParse
	}
	rvRegId = rvRegId

	if timeStamp > 0 {
		getRevocRegRequest, errGetRevRegReq := BuildGetRevocRegRequest(verifierDid, revRegId, timeStamp)

		if errGetRevRegReq != nil {
			return "", "", 0, errGetRevRegReq
		}

		getRevocRegResponse, errGetRevocRegResponse := SubmitRequest(poolHandle, getRevocRegRequest)
		if errGetRevocRegResponse != nil {
			return "", "", 0, errGetRevocRegResponse
		}

		_, revRegJsonResp, timeStamp2, errParseRevRegResp := ParseGetRevocRegResponse(getRevocRegResponse)
		if errParseRevRegResp != nil {
			return "", "", 0, errParseRevRegResp
		}

		return revRegDefJson, revRegJsonResp, timeStamp2, nil

	}

	return "", "", 0, errors.New("should not reach here")

}

// GetRevState  gets rev states
func GetRevState(poolHandle int, subjectDid string, revRegId string, credRevId string, from, to int64) (string, uint64, error) {
	getRevocRegDefRequest, errGetRevRegDefReq := BuildGetRevRegDefRequest(subjectDid, revRegId)
	if errGetRevRegDefReq != nil {
		return "", 0, errGetRevRegDefReq
	}

	getRevocRegDefResponse, errGetRevocRegDefResponse := SubmitRequest(poolHandle, getRevocRegDefRequest)
	if errGetRevocRegDefResponse != nil {
		return "", 0, errGetRevocRegDefResponse
	}

	_, revRegDefJson, errParse := ParseGetRevocRegDefResponse(getRevocRegDefResponse)
	if errParse != nil {
		return "", 0, errParse
	}

	getRevRegDeltaRequest, errGetDelta := BuildGetRevocRegDeltaRequest(subjectDid, revRegId, from, to)
	if errGetDelta != nil {
		return "", 0, errGetDelta
	}

	getRevRegDeltaResponse, errGetRevRegDeltaResp := SubmitRequest(poolHandle, getRevRegDeltaRequest)
	if errGetRevRegDeltaResp != nil {
		return "", 0, errGetRevRegDeltaResp
	}

	_, revRegDeltaJson, timeStamp, errParseDelta := ParseGetRevocRegDeltaResponse(getRevRegDeltaResponse)
	if errParseDelta != nil {
		return "", 0, errParseDelta
	}

	revRegDefObj, errParseJson := gabs.ParseJSON([]byte(revRegDefJson))
	if errParseJson != nil {
		return "", 0, errParseJson
	}

	tailsLocation := revRegDefObj.Path("value.tailsLocation").Data().(string)
	dir := filepath.Dir(tailsLocation)

	config := blobstorage.ConfigBlobStorage{
		BaseDir:    dir,
		UriPattern: "",
	}
	configStr := jsonObjectToString(&config)

	blobReaderHandle, errBlobHandle := IndyOpenBlobStorageReader("default", configStr)
	if errBlobHandle != nil {
		return "", 0, errBlobHandle
	}

	revStateJson, errRevState := CreateRevocationState(blobReaderHandle, revRegDefJson, revRegDeltaJson, timeStamp, credRevId)
	if errRevState != nil {
		return "", 0, errRevState
	}

	return revStateJson, timeStamp, nil
}

// GetSchema - gets schema
func GetSchema(ph int, did string, schemaId string) (string, string, error) {
	getSchemaRequestJson, errGet := BuildGetSchemaRequest(did, schemaId)
	if errGet != nil {
		return "", "", errGet
	}
	response, errSubmit := SubmitRequest(ph, getSchemaRequestJson)
	if errSubmit != nil {
		return "", "", errSubmit
	}

	id, js, errParse := ParseGetSchemaResponse(response)
	if errParse != nil {
		return "", "", errParse
	}

	return id, js, nil
}

// GetCredDef gets cred def
func GetCredDef(ph int, did string, credDefId string) (string, string, uint64, error) {
	getCredDefRequestJson, errGet := BuildGetCredentialDefinitionRequest(did, credDefId)
	if errGet != nil {
		return "", "", 0, errGet
	}
	response, errSubmit := SubmitRequest(ph, getCredDefRequestJson)
	if errSubmit != nil {
		return "", "", 0, errSubmit
	}

	responseJson, errParse := gabs.ParseJSON([]byte(response))
	if errParse != nil {
		return "", "", 0, errSubmit
	}

	creationTS := uint64(responseJson.Path("result.txnTime").Data().(float64))

	id, js, errParse := ParseGetCredDefResponse(response)
	if errParse != nil {
		return "", "", 0, errParse
	}

	return id, js, creationTS, nil
}

// EncodeValue - helper function to encode the raw value ...
// see implementation in aca-py https://github.com/hyperledger/aries-cloudagent-python/blob/main/aries_cloudagent/messaging/util.py
func EncodeValue(value interface{}) string {

	i32BoundUp := math.MaxInt32
	i32BoundMin := math.MinInt32

	if value == nil {
		value = "None"
	}

	var s string
	switch value.(type) {
	case bool:
		if value.(bool) {
			return "1"
		} else {
			return "0"
		}
	case int:
		v := value.(int)
		if v >= i32BoundMin && v <= i32BoundUp {
			return strconv.Itoa(v)
		}
		s = fmt.Sprintf("%v", value)
	case int64:
		v := value.(int64)
		if v >= int64(i32BoundMin) && v <= int64(i32BoundUp) {
			return strconv.Itoa(int(v))
		}

		s = strconv.FormatInt(value.(int64), 10)
	case uint64:
		v := value.(int64)
		if v >= 0 && v <= int64(i32BoundUp) {
			return strconv.Itoa(int(v))
		}
		s = strconv.FormatUint(value.(uint64), 10)
	case string:
		s = value.(string)
	case float32:
		s = fmt.Sprintf("%g", value.(float32))
	case float64:
		s = fmt.Sprintf("%g", value.(float64))
	default:
		fmt.Printf("Conversion error")
	}

	h := sha256.New()
	_, err := io.WriteString(h, s)
	if err != nil {
		return ""
	}
	digest := h.Sum(nil)

	bigint := big.NewInt(0).SetBytes(digest).String()
	return bigint
}

func GetOptionalValue(val string) *C.char {
	var ret *C.char

	if len(val) > 0 {
		ret = C.CString(val)
	} else {
		ret = nil
	}
	return ret
}
