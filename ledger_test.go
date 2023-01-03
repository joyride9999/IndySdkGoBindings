/*
// ******************************************************************
// Purpose: ledger unit testing
// Author: angel.draghici@siemens.com, adrian.toader@siemens.com
// Notes:
// Copyright (c): Siemens SRL
// This work is licensed under the terms of the Apache License Version 2.0.  See
// the LICENSE.txt file in the top-level directory.
// ******************************************************************
*/

package indySDK

import (
	"github.com/joyride9999/IndySdkGoBindings/indyUtils"
	"encoding/json"
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestSignRequest(t *testing.T) {
	walletHandle, errCreate := createWallet(trusteeConfig(), trusteeCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("createWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, trusteeConfig(), trusteeCredentials())

	did, _, errDid := CreateAndStoreDID(walletHandle, seedTrustee1)
	if errDid != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDid)
		return
	}

	message := `{
			"reqId": 1496822211362017764, 
			"identifier": "GJ1SzoWzavQYfNL9XkaJdrQejfztN4XqdsiV4ct3LXKL", 
			"operation": { 
				"type": "1",
           		"dest": "VsKV7grR1BUE29mG2Fm2kX",
            	"verkey": "GjZWsBLgZCR18aL468JAT7w9CZRiBnpxUPPgyQxh4voa"
				}
			}`
	expectedSignature := "65hzs4nsdQsTUqLCLy2qisbKLfwYKZSWoyh1C6CU59p5pfG3EHQXGAsjW4Qw4QdwkrvjSgQuyv8qyABcXRBznFKW"

	type args struct {
		Did     string
		Message string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"sign-message-works", args{Did: did, Message: message}, false},
		{"sign-empty-message", args{Did: did, Message: ""}, false},
		{"sign-message-invalid-did", args{Did: "invalid-did"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestSign, errSign := SignRequest(walletHandle, tt.args.Did, message)
			// hasError := errSign != nil
			if errSign != nil {
				t.Errorf("SignRequest() error = '%v'", errSign)
				return
			}
			if requestSign == "" {
				if tt.wantErr {
					t.Logf("Expected error: %s", indyUtils.GetIndyError(113))
					return
				}
			} else {
				// Check if signature is the same with the expected one
				parsedRequest, _ := gabs.ParseJSON([]byte(requestSign))
				signature := parsedRequest.Path(`signature`).Data()
				if signature == expectedSignature {
					fmt.Println(signature)
				}
			}
		})
	}
}

func TestSignAndSubmitRequest(t *testing.T) {
	poolHandle, errPool := getPoolLedger("test-sign-submit")
	if errPool != nil {
		t.Errorf("getPoolLedger() error = '%v'", errPool)
		return
	}
	defer ClosePoolHandle(poolHandle)

	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil {
		t.Errorf("createWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	did, verKey, errDid := CreateAndStoreDID(walletHandle, "")
	if errDid != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDid)
		return
	}

	whTrustee, errCreate := createWallet(trusteeConfig(), trusteeCredentials())
	if errCreate != nil {
		t.Errorf("createWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(whTrustee, trusteeConfig(), trusteeCredentials())

	trusteeDid, _, errDid := CreateAndStoreDID(whTrustee, seedTrustee1)
	if errDid != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDid)
		return
	}

	nymRequest, errNym := BuildNymRequest(trusteeDid, did, verKey, "", "")
	if errNym != nil {
		t.Errorf("BuildNymRequest() error = '%v'", errNym)
		return
	}

	poolCfgRequest, errPoolCfg := BuildPoolConfigRequest(trusteeDid, true, false)
	if errPoolCfg != nil {
		t.Errorf("BuildPoolConfigRequest() error = '%v'", errPoolCfg)
		return
	}

	type args struct {
		Request string
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"test-sign-and-submit-request-works", args{Request: nymRequest}, false},
		{"test-sign-and-submit-pool-config-request-works", args{Request: poolCfgRequest}, false},
		{"test-sign-and-submit-request-invalid-request", args{Request: "invalid-request"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.Request == "invalid-request" {
				_, errSign := SignAndSubmitRequest(poolHandle, whTrustee, trusteeDid, tt.args.Request)
				hasError := errSign != nil
				if hasError != tt.wantErr {
					t.Errorf("SignAndSubmitRequest() error = '%v'", errSign)
					return
				}
				if tt.wantErr {
					t.Log("Expected error: ", errSign)
					return
				}
			}
			if tt.args.Request == nymRequest {
				response, errSign := SignAndSubmitRequest(poolHandle, whTrustee, trusteeDid, tt.args.Request)
				if errSign != nil {
					t.Errorf("SignAndSubmitRequest() error = '%v'", errSign)
					return
				}

				signAndSubmitResponse, errParse := gabs.ParseJSON([]byte(response))
				if errParse != nil {
					t.Errorf("Gabs Parse error = '%v'", errParse)
					return
				}

				if signAndSubmitResponse.S("op").Data() != "REPLY" {
					t.Errorf("Values are not correct")
				}
			}
			if tt.args.Request == poolCfgRequest {
				response, errSign := SignAndSubmitRequest(poolHandle, whTrustee, trusteeDid, tt.args.Request)
				if errSign != nil {
					t.Errorf("SignAndSubmitRequest() error = '%v'", errSign)
					return
				}
				signAndSubmitResponse, errParse := gabs.ParseJSON([]byte(response))
				if errParse != nil {
					t.Errorf("Gabs Parse error = '%v'", errParse)
					return
				}
				resWrites := signAndSubmitResponse.Path("result.txn.data.writes").Data()
				if resWrites == false {
					t.Errorf("Values are not correct")
				}
			}
		})
	}
	return
}

func TestSubmitRequest(t *testing.T) {
	poolHandle, errPool := getPoolLedger("testpool2")
	if errPool != nil {
		t.Errorf("getPoolLedger() error = '%v'", errPool)
		return
	}
	defer ClosePoolHandle(poolHandle)

	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("createWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	whTrustee, errCreate := createWallet(trusteeConfig(), trusteeCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("createWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(whTrustee, trusteeConfig(), trusteeCredentials())

	did, verKey, errDid := CreateAndStoreDID(walletHandle, "")
	if errDid != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDid)
		return
	}

	trusteeDid, _, errDid := CreateAndStoreDID(whTrustee, seedTrustee1)
	if errDid != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDid)
		return
	}

	getNymRequest, errNym := prepareGetNymReq(poolHandle, whTrustee, trusteeDid, did, verKey, "STEWARD"); if errNym != nil {
		t.Errorf("BuildAndSendNymReq() error = '%v'", errNym)
		return
	}

	getAttribRequest, errAttrib := prepareGetAttribReq(poolHandle, walletHandle, trusteeDid, did, "", `{"test":"name"}`, ""); if errAttrib != nil {
		t.Errorf("prepareGetAttribReq() error = '%v'", errAttrib)
		return
	}

	schemaId, schemaJson, getSchemaRequest, errSchema := prepareGetSchemaReq(poolHandle, walletHandle, did, "gvt", "1.0", `["name", "age", "location"]`)
	if errSchema != nil {
		t.Errorf("BuildAndSendGetSchemaReq() error = '%v'", errSchema)
		return
	}
	schemaJsonParsed, _ := gabs.ParseJSON([]byte(schemaJson))
	schemaJsonParsed.DeleteP("seqNo")

	getSchemaResponse, errSign := SignAndSubmitRequest(poolHandle, walletHandle, did, getSchemaRequest); if errSign != nil {
		t.Errorf("SignAndSubmitRequest() error = '%v'", errSign)
		return
	}

	getSchemaResponseParsed, errParse := gabs.ParseJSON([]byte(getSchemaResponse)); if errParse != nil {
		t.Errorf("Gabs Parse error = '%v'", errParse)
		return
	}
	seqNo := getSchemaResponseParsed.Path("result.seqNo").String()
	seqNoInt, _ := strconv.Atoi(seqNo)
	getTxnReq, errBuild := BuildGetTxnRequest(did, "", seqNoInt); if errBuild != nil {
		t.Errorf("BuildGetTxnRequest() error = '%v'", errBuild)
		return
	}

	type args struct {
		Request string
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"test-submit-get-nym-request-works", args{Request: getNymRequest}, false},
		{"test-submit-get-attrib-request-works", args{Request: getAttribRequest}, false},
		{"test-submit-get-schema-request-works", args{Request: getSchemaRequest}, false},
		{"test-submit-get-txn-request-works", args{Request: getTxnReq}, false},
		{"test-submit-request-invalid-request", args{Request: "invalid-request"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.Request == "invalid-request" {
				_, errSubmit := SubmitRequest(poolHandle, tt.args.Request)
				hasError := errSubmit != nil
				if hasError != tt.wantErr {
					t.Errorf("SubmitRequest() error = '%v'", errSubmit)
					return
				}
				if tt.wantErr {
					t.Log("Expected error: ", errSubmit)
					return
				}
			}
			if tt.args.Request == getNymRequest {
				getNymResponse, errSubmit := SubmitRequest(poolHandle, tt.args.Request); if errSubmit != nil {
					t.Errorf("SubmitRequest() error = '%v'", errSubmit)
					return
				}

				parsedGetNym, errParseGetNym := ParseGetNymResponse(getNymResponse); if errParseGetNym != nil {
					t.Errorf("ParseGetNymResponse() error = '%v'", errParseGetNym)
					return
				}
				returnedNymData, _ := gabs.ParseJSON([]byte(parsedGetNym))
				if did != returnedNymData.Path("did").Data() || verKey != returnedNymData.Path("verkey").Data() {
					t.Errorf("Values are not correct")
					return
				}
			}
			if tt.args.Request == getAttribRequest {
				getAttribResponse, errSubmit := SubmitRequest(poolHandle, tt.args.Request)
				if errSubmit != nil {
					t.Errorf("SubmitRequest() error = '%v'", errSubmit)
					return
				}

				parsedGetAttribResponse, _ := gabs.ParseJSON([]byte(getAttribResponse))
				returnedRaw := parsedGetAttribResponse.Path("result.raw").String()

				expectedRaw := gabs.New()
				expectedRaw.Set(`{"test":"name"}`)
				if returnedRaw != expectedRaw.String() {
					t.Errorf("Values are not correct")
				}
			}
			if tt.args.Request == getSchemaRequest {
				getSchemaResponse, errSubmit := SubmitRequest(poolHandle, tt.args.Request)
				if errSubmit != nil {
					t.Errorf("SubmitRequest() error = '%v'", errSubmit)
					return
				}
				retSchemaId, retSchemaJson, errParse := ParseGetSchemaResponse(getSchemaResponse); if errParse != nil {
					t.Errorf("ParseGetSchemaResponse() error = '%v'", errParse)
					return
				}

				retSchemaJsonParsed, errParse2 := gabs.ParseJSON([]byte(retSchemaJson)); if errParse2 != nil {
					t.Errorf("Gabs Parse error = '%v'", errParse2)
					return
				}

				if schemaId != retSchemaId || !isIncluded(schemaJsonParsed, retSchemaJsonParsed) {
					t.Errorf("Values are not correct")
				}
			}
			if tt.args.Request == getTxnReq {
				getTxnResponse, errSubmit := SubmitRequest(poolHandle, tt.args.Request)
				if errSubmit != nil {
					t.Errorf("SubmitRequest() error = '%v'", errSubmit)
					return
				}
				getTxnResponseParsed, _ := gabs.ParseJSON([]byte(getTxnResponse))
				resName := getTxnResponseParsed.Path("result.data.txn.data.data.name").Data()
				resVersion := getTxnResponseParsed.Path("result.data.txn.data.data.version").Data()

				if resName != schemaJsonParsed.Path("name").Data().(string) || resVersion != schemaJsonParsed.Path("version").Data().(string) {
					t.Errorf("Values are not correct")
				}
			}
		})
	}
	return
}

func TestMultiSignRequest(t *testing.T) {
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("createWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	trusteeDid, _, errDid := CreateAndStoreDID(walletHandle, seedTrustee1)
	if errDid != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDid)
		return
	}

	did, _, errDid2 := CreateAndStoreDID(walletHandle, seedMy1)
	if errDid2 != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDid2)
		return
	}

	message := fmt.Sprintf(`{
			"reqId": 1496822211362017764, 
			"identifier": "%s", 
			"operation": { 
				"type": "1",
           		"dest": "VsKV7grR1BUE29mG2Fm2kX",
            	"verkey": "GjZWsBLgZCR18aL468JAT7w9CZRiBnpxUPPgyQxh4voa"
				}
			}`, trusteeDid)
	expectedSign1 := "3YnLxoUd4utFLzeXUkeGefAqAdHUD7rBprpSx2CJeH7gRYnyjkgJi7tCnFgUiMo62k6M2AyUDtJrkUSgHfcq3vua"
	expectedSign2 := "4EyvSFPoeQCJLziGVqjuMxrbuoWjAWUGPd6LdxeZuG9w3Bcbt7cSvhjrv8SX5e8mGf8jrf3K6xd9kEhXsQLqUg45"

	type args struct {
		WalletHandle int
		Did string
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"multi-sign-request-works", args{WalletHandle: walletHandle, Did: trusteeDid}, false},
		{"multi-sign-request-invalid-did", args{WalletHandle: walletHandle, Did: "8wZcEriaNLNKtteJvx7f8i"}, true},
		{"multi-sign-request-invalid-wallet-handle", args{WalletHandle: walletHandle + 100, Did: did}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			multiSignRequest, errMultiSign := MultiSignRequest(tt.args.WalletHandle, tt.args.Did, message)
			hasError := errMultiSign != nil
			if hasError != tt.wantErr {
				t.Errorf("MultiSignRequest() error = '%v'", errMultiSign)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errMultiSign)
				return
			}
			multiSignMsg, errMultiSign2 := MultiSignRequest(tt.args.WalletHandle, did, multiSignRequest)
			if errMultiSign2 != nil {
				t.Errorf("MultiSignRequest() error = '%v'", errMultiSign2)
				return
			} else {
				// Checking if signatures are the same as the expected ones
				parsed, _ := gabs.ParseJSON([]byte(multiSignMsg))
				sign1 := parsed.Path("signatures.V4SGRU86Z58d6TV7PBUe6f").Data()
				sign2 := parsed.Path("signatures.VsKV7grR1BUE29mG2Fm2kX").Data()
				if sign1 == expectedSign1 && sign2 == expectedSign2 {
					fmt.Println("Success")
				} else {
					t.Errorf("Values are not correct")
					return
				}
			}
		})
	}
	return
}

// TODO: Check function test. https://jira.hyperledger.org/browse/INDY-604
func TestBuildGetDdoRequest(t *testing.T) {
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("createWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	did, _, errDid := CreateAndStoreDID(walletHandle, seedTrustee1)
	if errDid != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDid)
		return
	}

	targetDid, _, errDid2 := CreateAndStoreDID(walletHandle, seedMy1)
	if errDid2 != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDid2)
		return
	}

	type args struct {
		TargetDid string
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"build-get-ddo-request-works", args{TargetDid: targetDid}, false},
		{"build-get-ddo-request-invalid-format-did", args{TargetDid: "invalid-id-string"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, errDdo := BuildGetDdoRequest(did, tt.args.TargetDid)
			hasError := errDdo != nil
			if hasError != tt.wantErr {
				t.Errorf("BuildGetDdoRequest() error = '%v'", errDdo)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errDdo)
				return
			}
			fmt.Println(request)
		})
	}
	return
}

func TestBuildNymRequest(t *testing.T) {
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("createWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	did, _, errDid := CreateAndStoreDID(walletHandle, seedTrustee1)
	if errDid != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDid)
		return
	}

	targetDid, targetVerKey, errDid2 := CreateAndStoreDID(walletHandle, seedMy1)
	if errDid2 != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDid2)
		return
	}

	type Operation struct {
		Type string `json:"type"`
		Dest string `json:"dest"`
		VerKey string `json:"verkey"`
		Role string `json:"role"`
	}
	type NymRequest struct {
		Identifier string `json:"identifier"`
		Op Operation	  `json:"operation"`
	}
	var op = Operation{Type: "1", Dest: targetDid, VerKey: targetVerKey, Role: "2"}
	var expectedNymRequest = NymRequest{did, op}

	type args struct {
		Did string
		Role string
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"build-nym-request-works", args{Did: did, Role: "STEWARD"}, false},
		{"build-nym-request-invalid-did", args{Did: "invalid-did", Role: "STEWARD"}, true},
		{"build-nym-role-invalid-role", args{Did: did, Role: "INVALID_ROLE"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestNym, errBuild := BuildNymRequest(tt.args.Did, targetDid, targetVerKey, "", tt.args.Role)
			hasError := errBuild != nil
			if hasError != tt.wantErr {
				t.Errorf("BuildNymRequest() error = '%v'", errBuild)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errBuild)
				return
			}

			var returnedNymRequest NymRequest
			err := json.Unmarshal([]byte(requestNym), &returnedNymRequest)
			if err != nil {
				t.Errorf("Unmarshal Error: %v", err)
				return
			}

			if expectedNymRequest == returnedNymRequest {
				fmt.Println("Success")
			} else {
				t.Errorf("Values are not correct")
			}
		})
	}
	return
}

func TestBuildAttribRequest(t *testing.T) {
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("createWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	did, _, errDid := CreateAndStoreDID(walletHandle, seedTrustee1)
	if errDid != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDid)
		return
	}

	targetDid, _, errDid2 := CreateAndStoreDID(walletHandle, seedMy1)
	if errDid2 != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDid2)
		return	}

	raw := `{"name":"test"}`
	hash := "83d907821df1c87db829e96569a11f6fc2e7880acba5e43d07ab786959e13bd3"
	enc := "aa3f41f619aa7e5e6b6d0de555e05331787f9bf9aa672b94b57ab65b9b66c3ea960b18a98e3834b1fc6cebf49f463b81fd6e3181"

	type Operation struct {
		Type string `json:"type"`
		Dest string `json:"dest"`
		Raw	string `json:"raw"`
		Hash string `json:"hash"`
		Enc string `json:"enc"`
	}
	type AttributeRequest struct {
		Identifier string    `json:"identifier"`
		Op         Operation `json:"operation"`
	}

	type args struct {
		Hash string
		Raw string
		Enc string
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"build-attrib-request-works-with-raw", args{Hash: "", Raw: raw, Enc: ""}, false},
		{"build-attrib-request-works-with-hash", args{Hash: hash, Raw: "", Enc: ""}, false},
		{"build-attrib-request-works-with-enc", args{Hash: "", Raw: "", Enc: enc}, false},
		{"build-attrib-request-missed-attribute", args{Hash: "", Raw: "", Enc: ""}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestAttrib, errBuild := BuildAttribRequest(did, targetDid, tt.args.Hash, tt.args.Raw, tt.args.Enc)
			hasError := errBuild != nil
			if hasError != tt.wantErr {
				t.Errorf("BuildAttribRequest() error = '%v'", errBuild)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errBuild)
				return
			}

			var op = Operation{Type: "100", Dest: targetDid, Raw: tt.args.Raw, Hash: tt.args.Hash, Enc: tt.args.Enc}
			var expectedAttributeRequest = AttributeRequest{did, op}

			var returnedAttributeRequest AttributeRequest
			err := json.Unmarshal([]byte(requestAttrib), &returnedAttributeRequest)
			if err != nil {
				t.Errorf("Unmarshal Error: %v", err)
				return
			}

			if expectedAttributeRequest == returnedAttributeRequest {
				fmt.Println("Success")
			} else {
				t.Errorf("Values are not correct")
			}
		})
	}
	return
}

func TestBuildGetAttribRequest(t *testing.T) {
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("createWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	did, _, errDid := CreateAndStoreDID(walletHandle, seedTrustee1)
	if errDid != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDid)
		return
	}

	targetDid, _, errDid2 := CreateAndStoreDID(walletHandle, seedMy1)
	if errDid2 != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDid2)
		return	}

	raw := `{"name":"test"}`
	hash := "83d907821df1c87db829e96569a11f6fc2e7880acba5e43d07ab786959e13bd3"
	enc := "aa3f41f619aa7e5e6b6d0de555e05331787f9bf9aa672b94b57ab65b9b66c3ea960b18a98e3834b1fc6cebf49f463b81fd6e3181"

	type Operation struct {
		Type string `json:"type"`
		Dest string `json:"dest"`
		Raw	string `json:"raw"`
		Hash string `json:"hash"`
		Enc string `json:"enc"`
	}
	type AttributeRequest struct {
		Identifier string    `json:"identifier"`
		Op         Operation `json:"operation"`
	}

	type args struct {
		Hash string
		Raw string
		Enc string
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"build-get-attrib-request-works-with-raw", args{Hash: "", Raw: raw, Enc: ""}, false},
		{"build-get-attrib-request-works-with-hash", args{Hash: hash, Raw: "", Enc: ""}, false},
		{"build-get-attrib-request-works-with-enc", args{Hash: "", Raw: "", Enc: enc}, false},
		{"build-get-attrib-request-missed-attribute", args{Hash: "", Raw: "", Enc: ""}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestAttrib, errBuild := BuildGetAttribRequest(did, targetDid, tt.args.Hash, tt.args.Raw, tt.args.Enc)
			hasError := errBuild != nil
			if hasError != tt.wantErr {
				t.Errorf("BuildGetAttribRequest() error = '%v'", errBuild)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errBuild)
				return
			}

			var op = Operation{Type: "104", Dest: targetDid, Raw: tt.args.Raw, Hash: tt.args.Hash, Enc: tt.args.Enc}
			var expectedAttributeRequest = AttributeRequest{did, op}

			var returnedAttributeRequest AttributeRequest
			err := json.Unmarshal([]byte(requestAttrib), &returnedAttributeRequest)
			if err != nil {
				t.Errorf("Unmarshal Error: %v", err)
				return
			}

			if expectedAttributeRequest == returnedAttributeRequest {
				fmt.Println("Success")
			} else {
				t.Errorf("Values are not correct")
			}
		})
	}
	return
}

func TestBuildGetNymRequest(t *testing.T) {
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil && errCreate.Error() != indyUtils.GetIndyError(203) {
		t.Errorf("createWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	did, _, errDid := CreateAndStoreDID(walletHandle, seedTrustee1)
	if errDid != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDid)
		return
	}

	targetDid, _, errDid2 := CreateAndStoreDID(walletHandle, seedMy1)
	if errDid2 != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDid2)
		return
	}

	type Operation struct {
		Type string `json:"type"`
		Dest string `json:"dest"`
	}

	type NymRequest struct {
		Identifier string `json:"identifier"`
		Op Operation `json:"operation"`
	}
	var op = Operation{Type: "105", Dest: targetDid}
	var expectedNymRequest = NymRequest{did, op}

	type args struct {
		Did string
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"build-get-nym-request-works", args{Did: did}, false},
		{"build-get-nym-request-invalid-did", args{Did: "invalid-did"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestNym, errBuild := BuildGetNymRequest(tt.args.Did, targetDid)
			hasError := errBuild != nil
			if hasError != tt.wantErr {
				t.Errorf("BuildGetNymRequest() error = '%v'", errBuild)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errBuild)
				return
			}

			var returnedNymRequest NymRequest
			err := json.Unmarshal([]byte(requestNym), &returnedNymRequest)
			if err != nil {
				t.Errorf("Unmarshal Error: %v", err)
				return
			}

			if expectedNymRequest == returnedNymRequest {
				fmt.Println("Success")
			} else {
				t.Errorf("Values are not correct")
			}
		})
	}
	return
}

func TestBuildSchemaRequest(t *testing.T) {
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil {
		t.Errorf("createWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	did, _, errDid := CreateAndStoreDID(walletHandle, seedTrustee1)
	if errDid != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDid)
		return
	}

	data := `{
		"id": "1",
		"name": "name",
		"version": "1.0",
		"attrNames": ["male"],
		"ver": "1.0"
	}`
	expectedRequest := `{
		"operation": {
			"type": "101", 
			"data": {
				"name": "name", 
				"version": "1.0", 
				"attr_names": ["male"]
			}
		}
	}`
	expectedSchemaRequest, err := gabs.ParseJSON([]byte(expectedRequest)); if err != nil {
		t.Errorf("Parse JSON error = '%v'", err)
		return
	}

	type args struct {
		SchemaJson string
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"build-schema-request-works", args{SchemaJson: data}, false},
		{"build-schema-request-wrong-schema-json", args{SchemaJson: "invalid-schema-json"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestSchema, errBuild := BuildSchemaRequest(did, tt.args.SchemaJson)
			hasError := errBuild != nil
			if hasError != tt.wantErr {
				t.Errorf("BuildSchemaRequest() error = '%v'", errBuild)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errBuild)
				return
			}

			returnedSchemaRequest, errParse := gabs.ParseJSON([]byte(requestSchema)); if errParse != nil {
				t.Errorf("Parse JSON error = '%v'", errParse)
				return
			}

			ok := isIncluded(expectedSchemaRequest, returnedSchemaRequest)
			if ok == true {
				fmt.Println("Success")
			} else {
				t.Errorf("Values are not correct")
			}
		})
	}
	return
}

func TestBuildGetSchemaRequest(t *testing.T) {
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil {
		t.Errorf("createWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	did, _, errDid := CreateAndStoreDID(walletHandle, seedTrustee1)
	if errDid != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDid)
		return
	}

	type Data struct {
		Name string `json:"name"`
		Version string `json:"version"`
	}
	type Operation struct {
		Type string `json:"type"`
		Dest string `json:"dest"`
		Data Data `json:"data"`
	}
	type SchemaRequest struct {
		Identifier string `json:"identifier"`
		Op	Operation `json:"operation"`
	}
	var data = Data{Name: "name", Version: "1.0"}
	var op = Operation{Type: "107", Dest: did, Data: data}
	var expectedSchemaRequest = SchemaRequest{did, op}

	type args struct {
		SchemaId string
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"build-get-schema-request-works", args{SchemaId: "V4SGRU86Z58d6TV7PBUe6f:2:name:1.0"}, false},
		{"build-get-schema-request-wrong-schema-id", args{SchemaId: "invalid-schema-id"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestSchema, errBuild := BuildGetSchemaRequest(did, tt.args.SchemaId)
			hasError := errBuild != nil
			if hasError != tt.wantErr {
				t.Errorf("BuildGetSchemaRequest() error = '%v'", errBuild)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errBuild)
				return
			}

			var returnedSchemaRequest SchemaRequest
			err := json.Unmarshal([]byte(requestSchema), &returnedSchemaRequest)
			if err != nil {
				t.Errorf("Unmarshal Error: %v", err)
				return
			}

			if expectedSchemaRequest == returnedSchemaRequest {
				fmt.Println("Success")
			} else {
				t.Errorf("Values are not correct")
			}
		})
	}
	return
}

func TestBuildCredentialDefinitionRequest(t *testing.T) {
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil {
		t.Errorf("createWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	did, _, errDid := CreateAndStoreDID(walletHandle, seedTrustee1)
	if errDid != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDid)
		return
	}

	data := `{
		"ver": "1.0",
		"id": "NcYxiDXkpYi6ov5FcYDi1e:3:CL:1",
		"schemaId": "1",
		"type": "CL",
		"tag": "TAG_1",
		"value": {
			"primary": {
				"n": "1",
				"s": "2",
				"r": {
					"name": "1",
					"master_secret": "3"
				},
			"rctxt": "1",
			"z": "1"
			}
		}
	}`
	expectedJson := fmt.Sprintf(`{
		"identifier":"%s",
		"operation": {
			"ref": 1, 
			"data": {
				"primary": {
					"n": "1",
					"s": "2",
					"r": {
						"name": "1",
						"master_secret": "3"
					}, 
				"rctxt": "1",
				"z": "1"
				}
			},
		"type": "102",
		"signature_type": "CL",
		"tag": "TAG_1"
		}
	}`, did)

	expectedCredDefRequest, err := gabs.ParseJSON([]byte(expectedJson)); if err != nil {
		t.Errorf("Parse JSON error = '%v'", err)
		return
	}

	type args struct {
		CredDefJson string
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"build-cred-definition-request-works", args{CredDefJson: data}, false},
		{"build-cred-definition-request-invalid-cred-def", args{CredDefJson: "invalid-cred-def-json"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestCredDef, errBuild := BuildCredentialDefinitionRequest(did, tt.args.CredDefJson)
			hasError := errBuild != nil
			if hasError != tt.wantErr {
				t.Errorf("BuildCredentialDefinitionRequest() error = '%v'", errBuild)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errBuild)
				return
			}

			returnedCredDefRequest, errParse := gabs.ParseJSON([]byte(requestCredDef)); if errParse != nil {
				t.Errorf("Parse JSON error = '%v'", errParse)
				return
			}

			ok := isIncluded(expectedCredDefRequest, returnedCredDefRequest)
			if ok == true {
				fmt.Println("Success")
			} else {
				t.Errorf("Values are not correct")
			}
		})
	}
	return
}

func TestBuildGetCredentialDefinitionRequest(t *testing.T) {
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil {
		t.Errorf("createWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	did, _, errDid := CreateAndStoreDID(walletHandle, seedTrustee1)
	if errDid != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDid)
		return
	}

	id := fmt.Sprintf("%s:3:CL:1:TAG_1", did)

	type Operation struct {
		Type string `json:"type"`
		Ref int `json:"ref"`
		SignatureType string `json:"signature_type"`
		Origin string `json:"origin"`
		Tag string `json:"tag"`
	}
	type CredDefRequest struct {
		Identifier string `json:"identifier"`
		Op Operation `json:"operation"`
	}
	var op = Operation{Type: "108", Ref: 1, SignatureType: "CL", Origin: did, Tag: "TAG_1"}
	var expectedCredDefRequest = CredDefRequest{Identifier: did, Op: op}

	type args struct {
		Did string
		CredDefId string
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"build-get-cred-definition-request-works", args{Did: did, CredDefId: id}, false},
		{"build-get-cred-definition-request-invalid-id", args{Did: did, CredDefId: "invalid-cred-def-id"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestCredDef, errBuild := BuildGetCredentialDefinitionRequest(tt.args.Did, tt.args.CredDefId)
			hasError := errBuild != nil
			if hasError != tt.wantErr {
				t.Errorf("BuildGetCredentialDefinition() error = '%v'", errBuild)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errBuild)
				return
			}

			var returnedCredDefRequest CredDefRequest
			err := json.Unmarshal([]byte(requestCredDef), &returnedCredDefRequest)
			if err != nil {
				t.Errorf("Unmarshal Error: %v", err)
				return
			}

			if expectedCredDefRequest == returnedCredDefRequest {
				fmt.Println("Success")
			} else {
				t.Errorf("Values are not correct")
			}
		})
	}
	return
}

func TestBuildNodeRequest(t *testing.T) {
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil {
		t.Errorf("createWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	did, _, errDid := CreateAndStoreDID(walletHandle, seedTrustee1)
	if errDid != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDid)
		return
	}

	targetDid := "VsKV7grR1BUE29mG2Fm2kX"
	data := `{
		"node_ip": "ip",
		"node_port": 1,
		"client_ip": "ip",
		"client_port": 1,
		"alias": "some",
		"services": ["VALIDATOR"],
		"blskey": "CnEDk9HrMnmiHXEV1WFgbVCRteYnPqsJwrTdcZaNhFVW"
	}`
	expectedRequest := fmt.Sprintf(`{
		"identifier": "%s", 
		"operation": {
			"type": "0", 
			"dest": "%s", 
			"data": %s
		}
	}`, did, targetDid, data)
	expectedNodeRequest, err := gabs.ParseJSON([]byte(expectedRequest)); if err != nil {
		t.Errorf("Parse JSON error = '%v'", err)
		return
	}

	type args struct {
		Data string
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"build-node-request-works", args{Data: data}, false},
		{"build-node-request-empty-data", args{Data: "{}"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestNode, errBuild := BuildNodeRequest(did, targetDid, tt.args.Data)
			hasError := errBuild != nil
			if hasError != tt.wantErr {
				t.Errorf("BuildNodeRequest() error = '%v'", errBuild)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errBuild)
				return
			}

			returnedNodeRequest, errParse := gabs.ParseJSON([]byte(requestNode)); if errParse != nil {
				t.Errorf("Parse JSON error = '%v'", errParse)
				return
			}

			ok := isIncluded(expectedNodeRequest, returnedNodeRequest)
			if ok == true {
				fmt.Println("Success")
			} else {
				t.Errorf("Values are not correct")
			}
		})
	}

	return
}

func TestBuildGetValidatorInfoRequest(t *testing.T) {
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil {
		t.Errorf("createWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	did, _, errDid := CreateAndStoreDID(walletHandle, seedTrustee1)
	if errDid != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDid)
		return
	}

	expectedRequest := fmt.Sprintf(`{
		"identifier": "%s",
		"operation": {
			"type": "119"
		}
	}`, did)

	type Operation struct {
		Type string `json:"type"`
	}
	type ValidatorInfoRequest struct {
		Identifier string `json:"identifier"`
		Op Operation `json:"operation"`
	}

	var expectedValidatorInfoRequest ValidatorInfoRequest
	errUnmarshal := json.Unmarshal([]byte(expectedRequest), &expectedValidatorInfoRequest)
	if errUnmarshal != nil {
		t.Errorf("Unmarshal error = '%v'", errUnmarshal)
		return
	}

	type args struct {
		Did string
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"build-get-validator-info-request-works", args{Did: did}, false},
		{"build-get-validator-info-invalid-did", args{Did: "invalid-did"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestValidatorInfo, errBuild := BuildGetValidatorInfoRequest(tt.args.Did)
			hasError := errBuild != nil
			if hasError != tt.wantErr {
				t.Errorf("BuildGetValidatorInfoRequest() error = '%v'", errBuild)
				return
			}
			if tt.wantErr {
				t.Logf("Expected error: '%v'", errBuild)
				return
			}

			var returnedValidatorInfoRequest ValidatorInfoRequest
			errUnmarshal2 := json.Unmarshal([]byte(requestValidatorInfo), &returnedValidatorInfoRequest)
			if errUnmarshal2 != nil {
				t.Errorf("Unmarshal error = '%v'", errUnmarshal)
				return
			}

			if expectedValidatorInfoRequest == returnedValidatorInfoRequest {
				fmt.Println("Success")
			} else {
				t.Errorf("Values are not correct")
			}
		})
	}
	return
}

func TestBuildGetTxnRequest(t *testing.T) {
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil {
		t.Errorf("createWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	did, _, errDid := CreateAndStoreDID(walletHandle, seedTrustee1)
	if errDid != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDid)
		return
	}

	type Operation struct {
		Type string `json:"type"`
		Data int `json:"data"`
		LedgerId int `json:"ledgerId"`
	}
	type TxnRequest struct {
		Identifier string `json:"identifier"`
		Op Operation `json:"operation"`
	}
	seqNo := 1

	// Two expected request, one for a given DID and the other one for no DID provided
	expectedRequest := fmt.Sprintf(`{
		"identifier": "%s",
		"operation": {
			"type": "3",
			"data": 1,
			"ledgerId": 0
		}
	}`, did)
	expectedRequest2 := `{
		"identifier": "LibindyDid111111111111",
		"operation": {
			"type": "3",
			"data": 1,
			"ledgerId": 1
		}
	}`

	var expectedTxnRequest TxnRequest
	var expectedTxnRequest2 TxnRequest

	errUnmarshal := json.Unmarshal([]byte(expectedRequest), &expectedTxnRequest)
	if errUnmarshal != nil {
		t.Errorf("Unmarshal error = '%v'", errUnmarshal)
		return
	}
	errUnmarshal2 := json.Unmarshal([]byte(expectedRequest2), &expectedTxnRequest2)
	if errUnmarshal2 != nil {
		t.Errorf("Unmarshal error = '%v'", errUnmarshal2)
		return
	}

	type args struct {
		Did string
		LedgerType string
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"build-get-txn-request-works", args{Did: did, LedgerType: "POOL"}, false},
		{"build-get-txn-request-works-optional-params", args{Did: "", LedgerType: ""}, false},
		{"build-get-txn-request-invalid-did", args{Did: "invalid-did", LedgerType: ""}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestTxn, errBuild := BuildGetTxnRequest(tt.args.Did, tt.args.LedgerType, seqNo)
			hasError := errBuild != nil
			if hasError != tt.wantErr {
				t.Errorf("BuildGetTxnRequest() error = '%v'", errBuild)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errBuild)
				return
			}

			var returnedTxnRequest TxnRequest
			errUnmarshal3 := json.Unmarshal([]byte(requestTxn), &returnedTxnRequest)
			if errUnmarshal3 != nil {
				t.Errorf("Unmarshal error = '%v'", errUnmarshal3)
				return
			}

			if len(tt.args.Did) > 0 {
				if expectedTxnRequest == returnedTxnRequest {
					fmt.Println("Success")
				} else {
					t.Errorf("Values are not correct")
				}
			} else {
				if expectedTxnRequest2 == returnedTxnRequest {
					fmt.Println("Success")
				} else {
					t.Errorf("Values are not correct")
				}
			}

		})
	}
	return
}

func TestBuildPoolConfigRequest(t *testing.T) {
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil {
		t.Errorf("createWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	did, _, errDid := CreateAndStoreDID(walletHandle, seedTrustee1)
	if errDid != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDid)
		return
	}

	type Operation struct {
		Type string `json:"type"`
		Writes bool `json:"writes"`
		Force bool `json:"force"`
	}
	type PoolConfigRequest struct {
		Identifier string `json:"identifier"`
		Op Operation `json:"operation"`
	}

	expectedRequest := fmt.Sprintf(`{
		"identifier": "%s", 
		"operation": {
			"type": "111",
			"writes": true,
			"force": false
		}
	}`, did)
	var expectedPoolConfigRequest PoolConfigRequest
	errUnmarshal := json.Unmarshal([]byte(expectedRequest), &expectedPoolConfigRequest)
	if errUnmarshal != nil {
		t.Errorf("Unmarshal error = '%v'", errUnmarshal)
		return
	}

	type args struct {
		Did string
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"build-pool-config-request-works", args{Did: did}, false},
		{"build-pool-config-request-invalid-did", args{Did: "invalid-did"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestPoolConfig, errBuild := BuildPoolConfigRequest(tt.args.Did, true, false)
			hasError := errBuild != nil
			if hasError != tt.wantErr {
				t.Errorf("BuildPoolConfigRequest() error = '%v'", errBuild)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errBuild)
				return
			}

			var returnedPoolConfigRequest PoolConfigRequest
			errUnmarshal2 := json.Unmarshal([]byte(requestPoolConfig), &returnedPoolConfigRequest)
			if errUnmarshal2 != nil {
				t.Errorf("Unmarshal error = '%v'", errUnmarshal2)
				return
			}

			if expectedPoolConfigRequest == returnedPoolConfigRequest {
				fmt.Println("Success")
			} else {
				t.Errorf("Values are not correct")
			}
		})
	}
	return
}

func TestBuildPoolRestartRequest(t *testing.T) {
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil {
		t.Errorf("createWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	did, _, errDid := CreateAndStoreDID(walletHandle, seedTrustee1)
	if errDid != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDid)
		return
	}

	type Operation struct {
		Type string `json:"type"`
		Action string `json:"action"`
		DateTime string `json:"dateTime"`
	}
	type PoolRestartRequest struct {
		Identifier string `json:"identifier"`
		Op Operation `json:"operation"`
	}

	type args struct {
		Did string
		Action string
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"build-pool-restart-request-works-for-start-action", args{Did: did, Action: "start"}, false},
		{"build-pool-restart-request-works-for-cancel-action", args{Did: did, Action: "cancel"}, false},
		{"build-pool-restart-request-invalid-did", args{Did: "invalid-did"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestPoolRestart, errBuild := BuildPoolRestartRequest(tt.args.Did, tt.args.Action, "0")
			hasError := errBuild != nil
			if hasError != tt.wantErr {
				t.Errorf("BuildPoolRestartRequest() error = '%v'", errBuild)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errBuild)
				return
			}

			expectedRequest := fmt.Sprintf(`{
				"identifier": "%s", 
				"operation": {
					"type": "118", 
					"action": "%s", 
					"datetime": "0"
				}
			}`, did, tt.args.Action)
			var expectedPoolRestartRequest PoolRestartRequest
			errUnmarshal := json.Unmarshal([]byte(expectedRequest), &expectedPoolRestartRequest)
			if errUnmarshal != nil {
				t.Errorf("Unmarshal error = '%v'", errUnmarshal)
				return
			}

			var returnedPoolRestartRequest PoolRestartRequest
			errUnmarshal2 := json.Unmarshal([]byte(requestPoolRestart), &returnedPoolRestartRequest)
			if errUnmarshal2 != nil {
				t.Errorf("Unmarshal error = '%v'", errUnmarshal2)
				return
			}

			if expectedPoolRestartRequest == returnedPoolRestartRequest {
				fmt.Println("Success")
			} else {
				t.Errorf("Values are not correct")
			}
		})
	}
	return
}

func TestBuildPoolUpgradeRequest(t *testing.T) {
	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil {
		t.Errorf("createWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	did, _, errDid := CreateAndStoreDID(walletHandle, seedTrustee1)
	if errDid != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDid)
		return
	}

	type args struct {
		Did string
		Action string
		Package string
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"build-pool-upgrade-request-works-for-start-action", args{Did: did, Action: "start", Package: ""}, false},
		{"build-pool-upgrade-request-works-for-cancel-action", args{Did: did, Action: "cancel", Package: ""}, false},
		{"build-pool-upgrade-request-works-optional-param", args{Did: did, Action: "start", Package: "some_package"}, false},
		{"build-pool-upgrade-request-invalid-did", args{Did: "invalid-did"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestPoolUpgrade, errBuild := BuildPoolUpgradeRequest(tt.args.Did, "upgrade-go", "2.5.0", tt.args.Action, "abc12345", 1,
				"{}","{}", false, false, tt.args.Package)
			hasError := errBuild != nil
			if hasError != tt.wantErr {
				t.Errorf("BuildPoolUpgradeRequest() error = '%v'", errBuild)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errBuild)
				return
			}

			expectedRequest := fmt.Sprintf(`{
				"identifier": "%s", 
				"operation": {
					"type": "109", 
					"name": "upgrade-go", 
					"version": "2.5.0", 
					"action": "%s", 
					"sha256": "abc12345", 
					"timeout": 1, 
					"schedule": [], 
					"justification": "{}", 
					"reinstall": false, 
					"force": false, 
					"package": "%s"
				}
			}`,
			tt.args.Did, tt.args.Action, tt.args.Package)
			expectedPoolUpgradeReq, errParse := gabs.ParseJSON([]byte(expectedRequest))
			if errParse != nil {
				t.Errorf("Gabs Parse error = '%v'", errParse)
				return
			}

			returnedPoolUpgradeReq, errParse2 := gabs.ParseJSON([]byte(requestPoolUpgrade))
			if errParse2 != nil {
				t.Errorf("Gabs Parse error = '%v'", errParse2)
				return
			}

			ok := isIncluded(expectedPoolUpgradeReq, returnedPoolUpgradeReq)
			if ok == true {
				fmt.Println("Success")
			} else {
				t.Errorf("Values are not correct")
			}
		})
	}
	return
}

func TestBuildRevocRegDefRequest(t *testing.T) {
	did := `Th7MpTaRZVRYnPiabds81Y`

	data := `{
		"ver": "1.0", 
		"id": "NcYxiDXkpYi6ov5FcYDi1e:4:NcYxiDXkpYi6ov5FcYDi1e:3:CL:1:CL_ACCUM:TAG_1", 
		"revocDefType": "CL_ACCUM", 
		"tag": "TAG1", 
		"credDefId": "NcYxiDXkpYi6ov5FcYDi1e:3:CL:1", 
		"value": {
			"issuanceType": "ISSUANCE_ON_DEMAND", 
			"maxCredNum": 5, 
			"tailsHash": "s",
			"tailsLocation": "http://tails.location.com", 
			"publicKeys": {
				"accumKey": {
					"z": "1 0000000000000000000000000000000000000000000000000000000000001111 1 0000000000000000000000000000000000000000000000000000000000000000 1 0000000000000000000000000000000000000000000000000000000000000000 1 0000000000000000000000000000000000000000000000000000000000000000 1 0000000000000000000000000000000000000000000000000000000000000000 1 0000000000000000000000000000000000000000000000000000000000000000 1 0000000000000000000000000000000000000000000000000000000000000000 1 0000000000000000000000000000000000000000000000000000000000000000 1 0000000000000000000000000000000000000000000000000000000000000000 1 0000000000000000000000000000000000000000000000000000000000000000 1 0000000000000000000000000000000000000000000000000000000000000000 1 0000000000000000000000000000000000000000000000000000000000000000"
				}
			}
		}
	}`

	expectedRequest := `{
		"operation": {
			"credDefId": "NcYxiDXkpYi6ov5FcYDi1e:3:CL:1", 
			"id": "NcYxiDXkpYi6ov5FcYDi1e:4:NcYxiDXkpYi6ov5FcYDi1e:3:CL:1:CL_ACCUM:TAG_1", 
			"revocDefType": "CL_ACCUM", 
			"tag": "TAG1", 
			"type": "113", 
			"value": {
				"issuanceType": "ISSUANCE_ON_DEMAND", 
				"maxCredNum": 5, 
				"tailsHash": "s",
				"tailsLocation": "http://tails.location.com", 
				"publicKeys": {
					"accumKey": {
						"z": "1 0000000000000000000000000000000000000000000000000000000000001111 1 0000000000000000000000000000000000000000000000000000000000000000 1 0000000000000000000000000000000000000000000000000000000000000000 1 0000000000000000000000000000000000000000000000000000000000000000 1 0000000000000000000000000000000000000000000000000000000000000000 1 0000000000000000000000000000000000000000000000000000000000000000 1 0000000000000000000000000000000000000000000000000000000000000000 1 0000000000000000000000000000000000000000000000000000000000000000 1 0000000000000000000000000000000000000000000000000000000000000000 1 0000000000000000000000000000000000000000000000000000000000000000 1 0000000000000000000000000000000000000000000000000000000000000000 1 0000000000000000000000000000000000000000000000000000000000000000"
					}
				}
			}
		}
	}`
	expectedRevocRegDefReq, errGabs := gabs.ParseJSON([]byte(expectedRequest))
	if errGabs != nil {
		t.Errorf("Gabs Parse error = '%v'", errGabs)
		return
	}

	type args struct {
		Did string
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"test-build-revoc-reg-def-req-works", args{Did: did}, false},
		{"test-build-revoc-reg-def-req-invalid-did", args{Did: "invalid-did"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestRevocRegDef, errBuild := BuildRevocRegDefRequest(tt.args.Did, data)
			hasError := errBuild != nil
			if hasError != tt.wantErr {
				t.Errorf("BuildRevocRegDefRequest() error = '%v'", errBuild)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errBuild)
				return
			}
			returnedRevocRegDefReq, errGabs2 := gabs.ParseJSON([]byte(requestRevocRegDef))
			if errGabs2 != nil {
				t.Errorf("Gabs Parse error = '%v'", errGabs2)
				return
			}
			ok := isIncluded(expectedRevocRegDefReq, returnedRevocRegDefReq)
			if ok == true {
				fmt.Println("Success")
			} else {
				t.Errorf("Values are not correct")
			}
		})
	}
	return
}

func TestBuildGetRevRegDefRequest(t *testing.T) {
	did := `Th7MpTaRZVRYnPiabds81Y`

	revRegDefId := "NcYxiDXkpYi6ov5FcYDi1e:4:NcYxiDXkpYi6ov5FcYDi1e:3:CL:1:CL_ACCUM:TAG_1"
	expectedRequest := `{
		"operation": {
			"type": "115", 
			"id": "NcYxiDXkpYi6ov5FcYDi1e:4:NcYxiDXkpYi6ov5FcYDi1e:3:CL:1:CL_ACCUM:TAG_1"
		}
	}`
	expectedGetRevRegDefReq, errGabs := gabs.ParseJSON([]byte(expectedRequest))
	if errGabs != nil {
		t.Errorf("Gabs Parse error = '%v'", errGabs)
		return
	}

	type args struct {
		Did string
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"test-build-get-rev-reg-def-req-works", args{Did: did}, false},
		{"test-build-get-rev-reg-def-req-empty-did", args{Did: ""}, false},
		{"test-build-get-rev-reg-def-req-invalid-did", args{Did: "invalid-did"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestGetRevRegDef, errBuild := BuildGetRevRegDefRequest(tt.args.Did, revRegDefId)
			hasError := errBuild != nil
			if hasError != tt.wantErr {
				t.Errorf("BuildGetRevRegDefRequest() error = '%v'", errBuild)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errBuild)
				return
			}
			returnedGetRevRegDefReq, errGabs2 := gabs.ParseJSON([]byte(requestGetRevRegDef))
			if errGabs2 != nil {
				t.Errorf("Gabs Parse error = '%v'", errGabs2)
				return
			}
			ok := isIncluded(expectedGetRevRegDefReq, returnedGetRevRegDefReq)
			if ok == true {
				fmt.Println("Success")
			} else {
				t.Errorf("Values are not correct")
			}
		})
	}
	return
}

func TestBuildRevocRegEntryRequest(t *testing.T) {
	did := `Th7MpTaRZVRYnPiabds81Y`

	revRegEntryValue := `{
		"ver": "1.0", 
		"value": {
			"accum": "1 0000000000000000000000000000000000000000000000000000000000000000 1 0000000000000000000000000000000000000000000000000000000000000000 1 0000000000000000000000000000000000000000000000000000000000000000 1 0000000000000000000000000000000000000000000000000000000000000000 1 0000000000000000000000000000000000000000000000000000000000000000 1 0000000000000000000000000000000000000000000000000000000000000000"
		}
	}`
	revRegDefId := "NcYxiDXkpYi6ov5FcYDi1e:4:NcYxiDXkpYi6ov5FcYDi1e:3:CL:1:CL_ACCUM:TAG_1"
	revRegType := "CL_ACCUM"
	expectedRequest := `{
		"operation": {
			"type": "114", 
			"revocRegDefId": "NcYxiDXkpYi6ov5FcYDi1e:4:NcYxiDXkpYi6ov5FcYDi1e:3:CL:1:CL_ACCUM:TAG_1", 
			"revocDefType": "CL_ACCUM", 
			"value": {
				"accum": "1 0000000000000000000000000000000000000000000000000000000000000000 1 0000000000000000000000000000000000000000000000000000000000000000 1 0000000000000000000000000000000000000000000000000000000000000000 1 0000000000000000000000000000000000000000000000000000000000000000 1 0000000000000000000000000000000000000000000000000000000000000000 1 0000000000000000000000000000000000000000000000000000000000000000"
			}
		}
	}`
	expectedRevRegEntryReq, errGabs := gabs.ParseJSON([]byte(expectedRequest))
	if errGabs != nil {
		t.Errorf("Gabs Parse error = '%v'", errGabs)
		return
	}

	type args struct {
		Did string
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"test-build-rev-reg-entry-req-works", args{Did: did}, false},
		{"test-build-rev-reg-entry-req-invalid-did", args{Did: "invalid-did"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestRevRegEntry, errBuild := BuildRevocRegEntryRequest(tt.args.Did, revRegDefId, revRegType, revRegEntryValue)
			hasError := errBuild != nil
			if hasError != tt.wantErr {
				t.Errorf("BuildRevocRegEntryRequest() error = '%v'", errBuild)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errBuild)
				return
			}
			returnedRevRegEntryReq, errGabs2 := gabs.ParseJSON([]byte(requestRevRegEntry))
			if errGabs2 != nil {
				t.Errorf("Gabs Parse error = '%v'", errGabs2)
				return
			}
			ok := isIncluded(expectedRevRegEntryReq, returnedRevRegEntryReq)
			if ok == true {
				fmt.Println("Success")
			} else {
				t.Errorf("Values are not correct")
			}
		})
	}
	return
}

func TestBuildGetRevocRegRequest(t *testing.T) {
	did := "Th7MpTaRZVRYnPiabds81Y"

	revRegDefId := "NcYxiDXkpYi6ov5FcYDi1e:4:NcYxiDXkpYi6ov5FcYDi1e:3:CL:1:CL_ACCUM:TAG_1"
	timestamp := 100
	expectedRequest := `{
		"operation": {
			"type": "116", 
			"revocRegDefId": "NcYxiDXkpYi6ov5FcYDi1e:4:NcYxiDXkpYi6ov5FcYDi1e:3:CL:1:CL_ACCUM:TAG_1", 
			"timestamp": 100
		}
	}`

	type Operation struct {
		Type string `json:"type"`
		RevocRegDefId string `json:"revocRegDefId"`
		Timestamp int64 `json:"timestamp"`
	}
	type GetRevocRegRequest struct {
		Op Operation `json:"operation"`
	}
	var expectedGetRevocRegRequest GetRevocRegRequest
	errUnmarshal := json.Unmarshal([]byte(expectedRequest), &expectedGetRevocRegRequest)
	if errUnmarshal != nil {
		t.Errorf("Json Unmarshal error = '%v'", errUnmarshal)
		return
	}

	type args struct {
		Did string
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"test-build-get-revoc-reg-request-works", args{Did: did}, false},
		{"test-build-get-revoc-reg-request-empty-did", args{Did: ""}, false},
		{"test-build-get-revoc-reg-invalid-did", args{Did: "invalid-did"}, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestGetRevocReg, errBuild := BuildGetRevocRegRequest(tt.args.Did, revRegDefId, int64(timestamp))
			hasError := errBuild != nil
			if hasError != tt.wantErr {
				t.Errorf("BuildGetRevocRegRequest() error = '%v'", errBuild)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errBuild)
				return
			}
			var returnedGetRevocRegRequest GetRevocRegRequest
			errUnmarshal2 := json.Unmarshal([]byte(requestGetRevocReg), &returnedGetRevocRegRequest)
			if errUnmarshal2 != nil {
				t.Errorf("Json Unmarshal error = '%v'", errUnmarshal)
				return
			}

			if expectedGetRevocRegRequest == returnedGetRevocRegRequest {
				fmt.Println("Success")
			} else {
				t.Errorf("Values are not correct")
			}
		})
	}
	return
}

func TestBuildGetRevocRegDeltaRequest(t *testing.T) {
	did := "Th7MpTaRZVRYnPiabds81Y"

	revRegDefId := "NcYxiDXkpYi6ov5FcYDi1e:4:NcYxiDXkpYi6ov5FcYDi1e:3:CL:1:CL_ACCUM:TAG_1"
	to := 100
	expectedRequest := `{	
		"operation": {
			"type": "117", 
			"revocRegDefId": "NcYxiDXkpYi6ov5FcYDi1e:4:NcYxiDXkpYi6ov5FcYDi1e:3:CL:1:CL_ACCUM:TAG_1", 
			"to": 100
		}
	}`

	type Operation struct {
		Type string `json:"type"`
		RevocRegDefId string `json:"revocRegDefId"`
		To int64 `json:"to"`
	}
	type GetRevocRegDeltaReq struct {
		Op Operation `json:"operation"`
	}
	var expectedGetRevocRegDeltaReq GetRevocRegDeltaReq
	errUnmarshal := json.Unmarshal([]byte(expectedRequest), &expectedGetRevocRegDeltaReq)
	if errUnmarshal != nil {
		t.Errorf("Json Unmarshal error = '%v'", errUnmarshal)
		return
	}

	type args struct {
		Did string
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"test-build-get-revoc-reg-delta-request-works", args{Did: did}, false},
		{"test-build-get-revoc-reg-delta-empty-did", args{Did: ""}, false},
		{"test-build-get-revoc-reg-delta-invalid-did", args{Did: "invalid-did"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestGetRevocRegDelta, errBuild := BuildGetRevocRegDeltaRequest(tt.args.Did, revRegDefId, 0, int64(to))
			hasError := errBuild != nil
			if hasError != tt.wantErr {
				t.Errorf("BuildGetRevocRegDeltaRequest() error = '%v'", errBuild)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errBuild)
				return
			}
			var returnedGetRevocRegDeltaReq GetRevocRegDeltaReq
			errUnmarshal2 := json.Unmarshal([]byte(requestGetRevocRegDelta), &returnedGetRevocRegDeltaReq)
			if errUnmarshal2 != nil {
				t.Errorf("Json Unmarshal error = '%v'", errUnmarshal)
				return
			}

			if expectedGetRevocRegDeltaReq == returnedGetRevocRegDeltaReq {
				fmt.Println("Success")
			} else {
				t.Errorf("Values are not correct")
			}
		})
	}
	return
}

func TestGetResponseMetadata(t *testing.T) {
	poolHandle, errPool := getPoolLedger("indypool")
	if errPool != nil {
		t.Errorf("getPoolLedger() error = '%v'", errPool)
		return
	}
	defer ClosePoolHandle(poolHandle)

	walletHandle, errCreate := createWallet(testConfig(), testCredentials())
	if errCreate != nil {
		t.Errorf("createWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(walletHandle, testConfig(), testCredentials())

	did, verKey, errDid := CreateAndStoreDID(walletHandle, "")
	if errDid != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDid)
		return
	}

	whTrustee, errCreate := createWallet(trusteeConfig(), trusteeCredentials())
	if errCreate != nil {
		t.Errorf("createWallet() error = '%v'", errCreate)
		return
	}
	defer walletCleanup(whTrustee, trusteeConfig(), trusteeCredentials())

	trusteeDid, _, errDid := CreateAndStoreDID(whTrustee, seedTrustee1)
	if errDid != nil {
		t.Errorf("CreateAndStoreDid() error = '%v'", errDid)
		return
	}

	nymRequest, errNym := BuildNymRequest(trusteeDid, did, verKey, "", "")
	if errNym != nil {
		t.Errorf("BuildNymRequest() error = '%v'", errNym)
		return
	}

	nymResponse, errSign := SignAndSubmitRequest(poolHandle, whTrustee, trusteeDid, nymRequest)
	if errSign != nil {
		t.Errorf("SignAndSubmitRequest() error = '%v'", errSign)
		return
	}

	getNymRequest, errGetNym := BuildGetNymRequest(trusteeDid, did)
	if errGetNym != nil {
		t.Errorf("BuildGetNymRequest() error = '%v'", errGetNym)
		return
	}

	getNymResponse, errSign2 := SignAndSubmitRequest(poolHandle, whTrustee, trusteeDid, getNymRequest)
	if errSign2 != nil {
		t.Errorf("SignAndSubmitRequest() error = '%v'", errSign2)
		return
	}

	type args struct {
		Response string
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"test-get-response-metadata-works-with-nym-req", args{Response: nymResponse}, false},
		{"test-get-response-metadata-works-with-get-nym-req", args{Response: getNymResponse}, false},
		{"test-get-response-metadata-invalid-response", args{Response: "invalid-response"}, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			responseMetadata, errMeta := GetResponseMetadata(tt.args.Response)
			hasError := errMeta != nil
			if hasError != tt.wantErr {
				t.Errorf("GetResponseMetadata() error = '%v'", errMeta)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errMeta)
				return
			}

			assert.Containsf(t, responseMetadata, "seqNo", "doesn't contain %s", "seqNo")
			assert.Containsf(t, responseMetadata, "txnTime", "doesn't contain %s", "txnTime")
			if tt.args.Response == getNymResponse {
				assert.Containsf(t, responseMetadata, "lastTxnTime", "doesn't contain %s", "lastTxnTime")
			}
			assert.NotContainsf(t, responseMetadata, "lastSeqNo", "contains %s", "lastSeqNo")
		})
	}
	
	return
}

func TestBuildAuthRuleRequest(t *testing.T) {
	did := "Th7MpTaRZVRYnPiabds81Y"
	constraint := `{
		"sig_count": 1, 
		"metadata": {},
		"role": "0",
		"constraint_id": "ROLE",
		"need_to_be_owner": false
	}`
	var expectedRequest string
	type args struct {
		Did string
		Action string
		OldValue string
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"test-build-auth-rule-request-works-for-add-action", args{Did: did, Action: "ADD", OldValue: ""}, false},
		{"test-build-auth-rule-request-works-for-edit-action", args{Did: did, Action: "EDIT", OldValue: "0"}, false},
		{"test-build-auth-rule-request-invalid-did", args{Did: "invalid-did"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestBuildAuthRule, errBuild := BuildAuthRuleRequest(tt.args.Did, "NYM", tt.args.Action, "role", tt.args.OldValue,
				"101", constraint)
			hasError := errBuild != nil
			if hasError != tt.wantErr {
				t.Errorf("BuildAuthRuleRequest() error = '%v'", errBuild)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errBuild)
				return
			}

			if len(tt.args.OldValue) != 0 {
				expectedRequest = fmt.Sprintf(`{
				"identifier": "%s", 
				"operation": {
					"type": "120",
					"auth_type": "1",
					"auth_action": "%s",
					"field": "role", 
					"old_value": "%s",
					"new_value": "101",
					"constraint": %s
					}
				}`, tt.args.Did, tt.args.Action, tt.args.OldValue, constraint)
			} else {
				expectedRequest = fmt.Sprintf(`{
				"identifier": "%s", 
				"operation": {
					"type": "120",
					"auth_type": "1",
					"auth_action": "%s",
					"field": "role",
					"new_value": "101",
					"constraint": %s
					}
				}`, tt.args.Did, tt.args.Action, constraint)
			}
			expectedBuiltAuthRuleReq, errParse := gabs.ParseJSON([]byte(expectedRequest)); if errParse != nil {
				t.Errorf("Gabs Parse Json error = '%v'", errParse)
				return
			}

			returnedBuildAuthRuleReq, errParse2 := gabs.ParseJSON([]byte(requestBuildAuthRule)); if errParse2 != nil {
				t.Errorf("Gabs Parse Json error = '%v'", errParse2)
				return
			}

			if !isIncluded(expectedBuiltAuthRuleReq, returnedBuildAuthRuleReq) {
				t.Errorf("Values are not correct")
			}
		})
	}

	return
}

func TestBuildAuthRulesRequest(t *testing.T) {
	did := "Th7MpTaRZVRYnPiabds81Y"
	constraint := `{
		"sig_count": 1,
		"metadata": {},
		"role": "0",
		"constraint_id": "ROLE",
		"need_to_be_owner": false
	}`
	data := fmt.Sprintf(`[{
		"auth_type": "1",
		"auth_action": "ADD",
		"field": "role",
		"new_value": "101",
		"constraint": %s
	},
	{
		"auth_type": "1",
		"auth_action": "EDIT",
		"field": "role",
		"old_value": "0",
		"new_value": "101",
		"constraint": %s
	}]`, constraint, constraint)
	expectedRequest := fmt.Sprintf(`{
		"identifier": "%s", 
		"operation": {
			"type": "122", 
			"rules": %s
		}
	}`, did, data)
	expectedBuildAuthRulesReq, errParse := gabs.ParseJSON([]byte(expectedRequest))
	if errParse != nil {
		t.Errorf("Gabs Parse error = '%v'", errParse)
		return
	}

	type args struct {
		Did string
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"test-build-auth-rules-request-works", args{Did: did}, false},
		{"test-build-auth-rules-request-invalid-did", args{Did: "invalid-did"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestBuildAuthRules, errBuild := BuildAuthRulesRequest(tt.args.Did, data)
			hasError := errBuild != nil
			if hasError != tt.wantErr {
				t.Errorf("BuildAuthRulesRequest() error = '%v'", errBuild)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errBuild)
				return
			}

			returnedBuildAuthRulesReq, errParse2 := gabs.ParseJSON([]byte(requestBuildAuthRules))
			if errParse2 != nil {
				t.Errorf("Gabs Parse error = '%v'", errParse2)
				return
			}

			if !isIncluded(expectedBuildAuthRulesReq, returnedBuildAuthRulesReq) {
				t.Errorf("Values are not correct")
			}
		})
	}
	return
}

func TestBuildGetAuthRuleRequest(t *testing.T) {
	did := "Th7MpTaRZVRYnPiabds81Y"
	expectedRequest := `{
		"identifier": "Th7MpTaRZVRYnPiabds81Y",
		"operation": {
			"type": "121",
			"auth_type": "1",
			"auth_action": "ADD",
			"field": "role",
			"new_value": "101"
		}
	}`
	expectedBuildGetAuthRuleReq, errParse := gabs.ParseJSON([]byte(expectedRequest))
	if errParse != nil {
		t.Errorf("Gabs Parse error = '%v'", errParse)
		return
	}

	type args struct {
		Did string
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"test-build-get-auth-rule-request-works", args{Did: did}, false},
		{"test-build-get-auth-rule-request-invalid-did", args{Did: "invalid-did"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestBuildGetAuthRule, errBuild := BuildGetAuthRuleRequest(tt.args.Did, "NYM", "ADD", "role", "", "101")
			hasError := errBuild != nil
			if hasError != tt.wantErr {
				t.Errorf("BuildGetAuthRuleRequest() error = '%v'", errBuild)
				return
			}
			if tt.wantErr {
				t.Logf("Expected error: '%v'", errBuild)
				return
			}

			returnedBuildGetAuthRuleReq, errParse2 := gabs.ParseJSON([]byte(requestBuildGetAuthRule))
			if errParse2 != nil {
				t.Errorf("Gabs Parse error = '%v'", errParse2)
				return
			}

			if !isIncluded(expectedBuildGetAuthRuleReq, returnedBuildGetAuthRuleReq) {
				t.Errorf("Values are not correct")
			}
		})
	}

	return
}

func TestBuildTxnAuthorAgreementRequest(t *testing.T) {
	did := "Th7MpTaRZVRYnPiabds81Y"
	expectedRequest := `{
		"identifier": "Th7MpTaRZVRYnPiabds81Y",
		"operation": {
			"type": "4",
			"text": "indy agreement",
			"version": "1.0.0"
		}
	}`
	expectedBuildTxnAuthorAgrmtReq, errParse := gabs.ParseJSON([]byte(expectedRequest))
	if errParse != nil {
		t.Errorf("Gabs Parse error = '%v'", errParse)
		return
	}

	expectedRequest2 := `{
		"identifier": "Th7MpTaRZVRYnPiabds81Y",
		"operation": {
			"type": "4",
			"text": "indy agreement",
			"version": "1.0.0",
			"ratification_ts": 12345,
            "retirement_ts": 54321
		}
	}`
	expectedBuildTxnAuthorAgrmtReq2, errParse2 := gabs.ParseJSON([]byte(expectedRequest2))
	if errParse2 != nil {
		t.Errorf("Gabs Parse error = '%v'", errParse2)
		return
	}

	type args struct {
		Did string
		RatificationTs int64
		RetirementTs int64
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
			{"test-build-txn-author-agreement-req-works", args{Did: did, RatificationTs: 0, RetirementTs: 0}, false},
			{"test-build-txn-author-agreement-req-works-with-ratification-retirement", args{Did: did, RatificationTs: 12345, RetirementTs: 54321},
				false},
			{"test-build-txn-author-agreement-req-invalid-did", args{Did: "invalid-did"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestBuildTxnAuthorAgrmt, errBuild := BuildTxnAuthorAgreementRequest(tt.args.Did, "indy agreement", "1.0.0",
				tt.args.RatificationTs, tt.args.RetirementTs)
			hasError := errBuild != nil
			if hasError != tt.wantErr {
				t.Errorf("BuildTxnAuthorAgreementRequest() error = '%v'", errBuild)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errBuild)
				return
			}
			returnedBuildTxnAuthorAgrmtReq, errParse3 := gabs.ParseJSON([]byte(requestBuildTxnAuthorAgrmt))
			if errParse3 != nil {
				t.Errorf("BuildAuthorAgreementRequest() error = '%v'", errParse3)
				return
			}

			if tt.args.RatificationTs > 0 && tt.args.RetirementTs > 0 {
				if !isIncluded(expectedBuildTxnAuthorAgrmtReq2, returnedBuildTxnAuthorAgrmtReq) {
					t.Errorf("Values are not correct")
					return
				}
			} else {
				if !isIncluded(expectedBuildTxnAuthorAgrmtReq, returnedBuildTxnAuthorAgrmtReq) {
					t.Errorf("Values are not correct")
					return
				}
			}
		})
	}
	return
}

func TestBuildDisableAllTxnAuthorAgreementsRequest(t *testing.T) {
	did := "Th7MpTaRZVRYnPiabds81Y"
	expectedRequest := `{
		"identifier": "Th7MpTaRZVRYnPiabds81Y",
		"operation": {
			"type": "8"
		}
	}`
	expectedBuildDisableAllTxnAuthorAgrmtReq, errParse := gabs.ParseJSON([]byte(expectedRequest))
	if errParse != nil {
		t.Errorf("Gabs Parse error = '%v'", errParse)
		return
	}

	type args struct {
		Did string
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"test-build-disable-all-txn-author-agrmt-req-works", args{Did: did}, false},
		{"test-build-disable-all-txn-author-agrmt-req-invalid-did", args{Did: "invalid-did"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestBuildDisableAllTxnAuthorAgrmt, errBuild := BuildDisableAllTxnAuthorAgreementsRequest(tt.args.Did)
			hasError := errBuild != nil
			if hasError != tt.wantErr {
				t.Errorf("BuildDisableAllTxnAuthorAgreementRequest() error = '%v'", errBuild)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errBuild)
				return
			}

			returnedBuildDisableAllTxnAuthorAgrmtReq, errParse2 := gabs.ParseJSON([]byte(requestBuildDisableAllTxnAuthorAgrmt))
			if errParse2 != nil {
				t.Errorf("Gabs Parse error = '%v'", errParse2)
				return
			}

			if !isIncluded(expectedBuildDisableAllTxnAuthorAgrmtReq, returnedBuildDisableAllTxnAuthorAgrmtReq) {
				t.Errorf("Values are not correct")
				return
			}
		})
	}
	return
}

func TestBuildGetTxnAuthorAgreementRequest(t *testing.T) {
	did := `Th7MpTaRZVRYnPiabds81Y`
	data := `{
		"digest": "83d907821df1c87db829e96569a11f6fc2e7880acba5e43d07ab786959e13bd3"
	}`

	expectedRequest := `{
		"operation": {
			"type": "6",
			"digest": "83d907821df1c87db829e96569a11f6fc2e7880acba5e43d07ab786959e13bd3"
		}
	}`
	expectedBuildGetTxnAuthorAgrmtReq, errParse := gabs.ParseJSON([]byte(expectedRequest))
	if errParse != nil {
		t.Errorf("Gabs Parse error = '%v'", errParse)
		return
	}

	type args struct {
		Did string
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"test-build-get-txn-author-agreement-req-works", args{Did: did}, false},
		{"test-build-get-txn-author-agreement-req-empty-did", args{Did: ""}, false},
		{"test-build-get-txn-author-agreement-req-invalid-did", args{Did: "invalid-did"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestBuildGetTxnAuthorAgrmt, errBuild := BuildGetTxnAuthorAgreementRequest(tt.args.Did, data)
			hasError := errBuild != nil
			if hasError != tt.wantErr {
				t.Errorf("BuildGetTxnAuthorAgreementRequest() error = '%v'", errBuild)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errBuild)
				return
			}

			returnedBuildGetTxnAuthorAgrmtReq, errParse2 := gabs.ParseJSON([]byte(requestBuildGetTxnAuthorAgrmt))
			if errParse2 != nil {
				t.Errorf("Gabs Parse error = '%v'", errParse2)
				return
			}

			if !isIncluded(expectedBuildGetTxnAuthorAgrmtReq, returnedBuildGetTxnAuthorAgrmtReq) {
				t.Errorf("Values are not correct")
				return
			}
		})
	}
	return
}

func TestBuildAcceptanceMechanismsRequest(t *testing.T) {
	did := "Th7MpTaRZVRYnPiabds81Y"
	aml := `{
		"acceptance mechanism label 1": "some acceptance mechanism description 1"
	}`
	amlContext := "some context"
	expectedRequest := fmt.Sprintf(`{
		"identifier": "Th7MpTaRZVRYnPiabds81Y",
		"operation": {
			"type": "5",
			"aml": %s,
			"version": "1.0.0"
		}
	}`, aml)
	expectedBuildAcceptanceMechanismsReq, errParse := gabs.ParseJSON([]byte(expectedRequest))
	if errParse != nil {
		t.Errorf("Gabs Parse error = '%v'", errParse)
		return
	}

	expectedRequest2 := fmt.Sprintf(`{
		"identifier": "Th7MpTaRZVRYnPiabds81Y",
		"operation": {
			"type": "5",
			"aml": %s,
			"version": "1.0.0",
			"amlContext": "some context"
		}
	}`, aml)
	expectedBuildAcceptanceMechanismsReq2, errParse2 := gabs.ParseJSON([]byte(expectedRequest2))
	if errParse2 != nil {
		t.Errorf("Gabs Parse error = '%v'", errParse2)
		return
	}
	type args struct {
		Did string
		AmlContext string
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"test-build-acceptance-mechanisms-request-works", args{Did: did, AmlContext: ""}, false},
		{"test-build-acceptance-mechanisms-request-works-with-context", args{Did: did, AmlContext: amlContext}, false},
		{"test-build-acceptance-mechanisms-request-invalid-did", args{Did: "invalid-did"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestBuildAcceptanceMechanisms, errBuild := BuildAcceptanceMechanismsRequest(tt.args.Did, aml, "1.0.0", tt.args.AmlContext)
			hasError := errBuild != nil
			if hasError != tt.wantErr {
				t.Errorf("BuildAcceptanceMechanismsRequest() error = '%v'", errBuild)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errBuild)
				return
			}

			returnedBuildAcceptanceMechanismsReq, errParse3 := gabs.ParseJSON([]byte(requestBuildAcceptanceMechanisms))
			if errParse3 != nil {
				t.Errorf("Gabs Parse error = '%v'", errParse3)
				return
			}
			if len(tt.args.AmlContext) == 0 {
				if !isIncluded(expectedBuildAcceptanceMechanismsReq, returnedBuildAcceptanceMechanismsReq) {
					t.Error("Values are not correct")
				}
			} else {
				if !isIncluded(expectedBuildAcceptanceMechanismsReq2, returnedBuildAcceptanceMechanismsReq) {
					t.Error("Values are not correct")
				}
			}
		})
	}
	return
}

func TestBuildGetAcceptanceMechanismsRequest (t *testing.T) {
	did := "Th7MpTaRZVRYnPiabds81Y"
	expectedReqEmpty := `{
		"identifier": "LibindyDid111111111111",
		"operation": {
			"type": "7"
		}
	}`

	expectedReqTimestamp := `{
		"identifier": "Th7MpTaRZVRYnPiabds81Y",
		"operation": {
			"type": "7",
			"timestamp": 123456789
		}
	}`

	expectedReqVersion  := `{
		"identifier": "Th7MpTaRZVRYnPiabds81Y",
		"operation": {
			"type": "7",
			"version": "1.0.0"
		}
	}`

	expectedBuildGetAccMechsReq, errParse := gabs.ParseJSON([]byte(expectedReqEmpty)); if errParse != nil {
		t.Errorf("Gabs Parse error = '%v'", errParse)
		return
	}
	expectedBuildGetAccMechsReq2, errParse2 := gabs.ParseJSON([]byte(expectedReqTimestamp)); if errParse2 != nil {
		t.Errorf("Gabs Parse error = '%v'", errParse2)
		return
	}
	expectedBuildGetAccMechsReq3, errParse3 := gabs.ParseJSON([]byte(expectedReqVersion)); if errParse3 != nil {
		t.Errorf("Gabs Parse error = '%v'", errParse3)
		return
	}

	type args struct {
		Did string
		Timestamp int64
		Version string
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"test-build-get-acceptance-mechanisms-req-works-empty", args{Did: "", Timestamp: -1, Version: ""}, false},
		{"test-build-get-acceptance-mechanisms-req-with-timestamp", args{Did: did, Timestamp: 123456789, Version: ""}, false},
		{"test-build-get-acceptance-mechanisms-req-with-version", args{Did: did, Timestamp: -1, Version: "1.0.0"}, false},
		{"test-build-get-acceptance-mechanisms-req-invalid-did", args{Did: "invalid-did"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestBuildGetAccMechsReq, errBuild := BuildGetAcceptanceMechanismsRequest(tt.args.Did, tt.args.Timestamp, tt.args.Version)
			hasError := errBuild != nil
			if hasError != tt.wantErr {
				t.Errorf("BuildGetAcceptanceMechanismsRequest() error = '%v'", errBuild)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errBuild)
				return
			}

			returnedBuildGetAccMechsReq, errParse4 := gabs.ParseJSON([]byte(requestBuildGetAccMechsReq)); if errParse4 != nil {
				t.Errorf("Gabs Parse error = '%v'", errParse4)
				return
			}

			if len(tt.args.Did) == 0 {
				if !isIncluded(expectedBuildGetAccMechsReq, returnedBuildGetAccMechsReq) {
					t.Errorf("Values are not correct")
					return
				}
			}
			if tt.args.Timestamp > 0 {
				if !isIncluded(expectedBuildGetAccMechsReq2, returnedBuildGetAccMechsReq) {
					t.Errorf("Values are not correct")
					return
				}
			}
			if len(tt.args.Version) > 0 {
				if !isIncluded(expectedBuildGetAccMechsReq3, returnedBuildGetAccMechsReq) {
					t.Errorf("Values are not correct")
					return
				}
			}
		})
	}

	return
}

func TestAppendTxnAuthorAgreementAcceptanceToRequest(t *testing.T) {
	taaDigest := `050e52a57837fff904d3d059c8a123e3a04177042bf467db2b2c27abd8045d5e`
	requestJson := `{
		"reqId": 1496822211362017764,
		"identifier": "GJ1SzoWzavQYfNL9XkaJdrQejfztN4XqdsiV4ct3LXKL",
		"operation": {
			"type": "1",
			"dest": "VsKV7grR1BUE29mG2Fm2kX",
       		"verkey": "GjZWsBLgZCR18aL468JAT7w9CZRiBnpxUPPgyQxh4voa"
		}
	}`

	expectedMeta := `{
		"mechanism": "acceptance type 1",
		"taaDigest": "050e52a57837fff904d3d059c8a123e3a04177042bf467db2b2c27abd8045d5e",
		"time": 123379200
	}`
	expectedRequest := fmt.Sprintf(`{
		"identifier": "GJ1SzoWzavQYfNL9XkaJdrQejfztN4XqdsiV4ct3LXKL",
		"operation": {
			"type": "1",
			"dest": "VsKV7grR1BUE29mG2Fm2kX",
			"verkey": "GjZWsBLgZCR18aL468JAT7w9CZRiBnpxUPPgyQxh4voa"
		},
		"taaAcceptance": %s
	}`, expectedMeta)
	expectedAppendTaaAcceptanceToReq, errParse := gabs.ParseJSON([]byte(expectedRequest)); if errParse != nil {
		t.Errorf("Gabs Parse error = '%v'", errParse)
		return
	}

	type args struct {
		Text string
		Version string
		TaaDigest string
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"test-append-txn-author-agrmt-acceptance-to-req-works-for-text-version", args{Text: "some agreement text", Version: "1.0.0", TaaDigest: ""},
			false},
		{"test-append-txn-author-agrmt-acceptance-to-req-works-for-hash", args{Text: "", Version: "", TaaDigest: taaDigest}, false},
		{"test-append-txn-author-agrmt-acceptance-to-req-missing-params", args{Text: "some agreement text", Version: "", TaaDigest: ""}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestAppendTaaAgreement, errAppend := AppendTxnAuthorAgreementAcceptanceToRequest(requestJson, tt.args.Text, tt.args.Version, tt.args.TaaDigest,
				"acceptance type 1", 123379200)
			hasError := errAppend != nil
			if hasError != tt.wantErr {
				t.Errorf("AppendTxnAuthorAgreementAcceptanceToRequest() error = '%v'", errAppend)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errAppend)
				return
			}
			returnedAppendTaaAcceptanceToReq, errParse2 := gabs.ParseJSON([]byte(requestAppendTaaAgreement)); if errParse2 != nil {
				t.Errorf("Gabs Parse error = '%v'", errAppend)
				return
			}

			if !isIncluded(expectedAppendTaaAcceptanceToReq, returnedAppendTaaAcceptanceToReq) {
				t.Errorf("Values are not correct")
			}
		})
	}
	return
}

func TestAppendRequestEndorser(t *testing.T) {
	did := "V4SGRU86Z58d6TV7PBUe6f"
	requestJson := `{
		"reqId": 1496822211362017764,
		"identifier": "GJ1SzoWzavQYfNL9XkaJdrQejfztN4XqdsiV4ct3LXKL",
		"operation": {
			"type": "1",
			"dest": "VsKV7grR1BUE29mG2Fm2kX",
       		"verkey": "GjZWsBLgZCR18aL468JAT7w9CZRiBnpxUPPgyQxh4voa"
		}
	}`
	expectedRequest := `{
		"endorser": "V4SGRU86Z58d6TV7PBUe6f",
		"identifier": "GJ1SzoWzavQYfNL9XkaJdrQejfztN4XqdsiV4ct3LXKL",
		"operation": {
			"type": "1",
			"dest": "VsKV7grR1BUE29mG2Fm2kX", 
			"verkey": "GjZWsBLgZCR18aL468JAT7w9CZRiBnpxUPPgyQxh4voa"
		}
	}`

	expectedAppendReqEndorser, errParse := gabs.ParseJSON([]byte(expectedRequest)); if errParse != nil {
		t.Errorf("Gabs Parse error = '%v'", errParse)
		return
	}

	type args struct {
		Did string
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"test-append-request-endorser-works", args{Did: did}, false},
		{"test-append-request-endorser-invalid-did", args{Did: "invalid-did"}, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestAppendEndorser, errBuild := AppendRequestEndorser(requestJson, tt.args.Did)
			hasError := errBuild != nil
			if hasError != tt.wantErr {
				t.Errorf("AppendRequestEndorser() error = '%v'", errBuild)
				return
			}
			if tt.wantErr {
				t.Log("Expected error: ", errBuild)
				return
			}

			returnedAppendReqEndorser, errParse2 := gabs.ParseJSON([]byte(requestAppendEndorser)); if errParse2 != nil {
				t.Errorf("Gabs Parse error = '%v'", errParse2)
				return
			}

			if !isIncluded(expectedAppendReqEndorser, returnedAppendReqEndorser) {
				t.Errorf("Values are not correct")
			}
		})
	}
	return
}