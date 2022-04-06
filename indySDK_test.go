/*
// ******************************************************************
// Purpose: unit testing for helper functions
// Author:  alexandru.leonte@siemens.com
// Notes:
// Copyright (c): Siemens SRL
// This work is licensed under the terms of the Apache License Version 2.0.  See
// the LICENSE.txt file in the top-level directory.
// ******************************************************************
*/

package indySDK

import (
	"fmt"
	"testing"
)

func Test_EncodeValue(t *testing.T) {
	// Prepare and run test cases
	type args struct {
		Raw     interface{}
		Encoded string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"bool true", args{true, "1"}, false},
		{"bool false", args{false, "0"}, false},
		{"none/nil", args{nil, "99769404535520360775991420569103450442789945655240760487761322098828903685777"}, false},
		{"str none/nil", args{"None", "99769404535520360775991420569103450442789945655240760487761322098828903685777"}, false},
		{"empty", args{"", "102987336249554097029535212322581322789799900648198034993379397001115665086549"}, false},
		{"maxint32", args{2147483647, "2147483647"}, false},
		{"maxint32 as int64", args{int64(2147483647), "2147483647"}, false},
		{"maxint32 as uint64", args{uint64(2147483647), "2147483647"}, false},
		{"maxint32+1", args{2147483648, "26221484005389514539852548961319751347124425277437769688639924217837557266135"}, false},
		{"minint32", args{-2147483648, "-2147483648"}, false},
		{"minint32-1", args{-2147483649, "-68956915425095939579909400566452872085353864667122112803508671228696852865689"}, false},
		{"float", args{0.0, "62838607218564353630028473473939957328943626306458686867332534889076311281879"}, false},
		{"str float", args{"0.0", "62838607218564353630028473473939957328943626306458686867332534889076311281879"}, false},
		{"addr2", args{"101 Wilson Lane", "68086943237164982734333428280784300550565381723532936263016368251445461241953"}, false},
		{"zip", args{"87121", "87121"}, false},
		{"city", args{"SLC", "101327353979588246869873249766058188995681113722618593621043638294296500696424"}, false},
		{"addr1", args{"101 Tela Lane", "63690509275174663089934667471948380740244018358024875547775652380902762701972"}, false},
		{"state", args{"UT", "93856629670657830351991220989031130499313559332549427637940645777813964461231"}, false},
	}
	fmt.Println("Test")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded := EncodeValue(tt.args.Raw)
			if encoded != tt.args.Encoded {
				if tt.wantErr {
					t.Errorf("EncodeValue() error = encoed = '%s', want = '%s'", encoded, tt.args.Encoded)
				}

				return
			}
		})
	}
}