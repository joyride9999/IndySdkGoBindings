/*
// ******************************************************************
// Purpose: pool unit testing
// Author: angel.draghici@siemens.com, adrian.toader@siemens.com
// Notes:
// Copyright (c): Siemens SRL
// This work is licensed under the terms of the Apache License Version 2.0.  See
// the LICENSE.txt file in the top-level directory.
// ******************************************************************
*/


package indySDK

import (
	"fmt"
	"indySDK/indyUtils"
	"indySDK/pool"
	"testing"
)

func TestSetPoolProtocolVersion(t *testing.T) {

	// Prepare and run test cases
	type args struct {
		protocolVersion uint64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"valid-protocol", args{2}, false},
		{"invalid-protocol", args{99}, true},
	}
	fmt.Println("Test")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SetPoolProtocolVersion(tt.args.protocolVersion)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetPoolProtocolVersion() error = '%v', wantErr = '%v'", err, tt.wantErr)
				return
			}
		})
	}
}

func TestOpenPoolLedgerConfig(t *testing.T) {
	var poolLedger pool.Pool
	poolLedger.Name = "Siemens4"
	poolLedger.GenesisTxn = "pool.txn"

	errSP := SetPoolProtocolVersion(2)
	if errSP != nil {
		t.Errorf("SetPoolProtocolVersion() error ")
	}

	errC := CreatePoolLedgerConfig(poolLedger)
	if errC != nil && errC.Error() != indyUtils.GetIndyError(306) {
		t.Errorf("CreatePoolLedgerConfig() error ")
	}

	hPool, errOp := OpenPoolLedgerConfig(poolLedger)
	if errOp != nil {
		t.Errorf("OpenPoolLedgerConfig() error ")
	}

	ClosePoolHandle(hPool)
	return
}