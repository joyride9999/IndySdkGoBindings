/*
// ******************************************************************
// Purpose: unit testing
// Author: angel.draghici@siemens.com
// Notes:
// Copyright (c): Siemens SRL
// This work is licensed under the terms of the Apache License Version 2.0.  See
// the LICENSE.txt file in the top-level directory.
// ******************************************************************
*/


package indySDK

import (
	"github.com/Jeffail/gabs/v2"
	"testing"
)

func TestIsIncluded(t *testing.T) {
	type args struct {
		ExpectedJson string
		Json string
	}
	tests := []struct {
		name string
		args args
		expectedOk bool
	}{
		{"test-empty-json-objects", args{ExpectedJson: `{"arg": {}}`, Json: `{"arg": {}}`}, true},
		{"test-one-empty-json-object", args{ExpectedJson: `{"arg": {}}`, Json: `{"arg": {"a": "1"}}`}, false},
		{"test-json-object", args{ExpectedJson: `{"arg1": 1, "arg2": 2}`, Json: `{"arg1": 1, "arg2": 2}`}, true},
		{"test-diff-object", args{ExpectedJson: `{"arg1": 1, "arg2": 2}`, Json: `{"arg1": 1, "arg2": 0}`}, false},
		{"test-nested-objects", args{ExpectedJson: `{"arg1": {"op": {"t": 1, "t2": 1, "t3": {"a": 1, "b": 2}}}}`,
			Json: `{"arg1": {"op": {"t": 1, "t2": 1, "t3": {"a": 1, "b": 2}}}}`}, true},
		{"test-diff-nested-value", args{ExpectedJson: `{"arg1": {"op": {"t": 1, "t2": 1, "t3": {"a": 1, "b": 2}}}}`,
			Json: `{"arg1": {"op": {"t": 1, "t2": 1, "t3": {"a": 1, "b": 0}}}}`}, false},
		{"test-diff-nested-arg", args{ExpectedJson: `{"arg1": {"op": {"t": 1, "t2": 1, "t3": {"a": 1, "b": 2}}}}`,
			Json: `{"arg1": {"op": {"t": 1, "t2": 1}}, "t3": {"a": {"b": 2}}}`}, false},
		{"test-object-with-array-bigger", args{ExpectedJson: `{"arg1": [1, 2, 3]}`, Json: `{"arg1": [1, 2, 3, 4]}`}, false},
		{"test-object-with-array-smaller", args{ExpectedJson: `{"arg1": [1, 2, 3]}`, Json: `{"arg1": [1, 2]}`}, false},
		{"test-object-with-array-equal", args{ExpectedJson: `{"arg1": [1, 2]}`, Json: `{"arg1": [1, 2]}`}, true},
		{"test-object-with-first-array-item-diff", args{ExpectedJson: `{"arg1": [0, 2]}`, Json: `{"arg1": [1, 2]}`}, false},
		{"test-object-with-second-array-item-diff", args{ExpectedJson: `{"arg1": [1, 0]}`, Json: `{"arg1": [1, 2]}`}, false},
		{"test-object-with-array-index-diff", args{ExpectedJson: `{"arg1": [1, 2]}`, Json: `{"arg1": [2, 1]}`}, true},
		{"test-object-with-array-same-obj", args{ExpectedJson: `{"arg1": [{"a": 1, "b": 2}]}`, Json: `{"arg1": [{"b": 2, "a": 1}]}`}, true},
		{"test-object-with-array-with-diff-obj", args{ExpectedJson: `{"arg1": [{"a": 1}]}`, Json: `{"arg1": [{"a": 1, "b": 2}]}`}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsedRequest, _ := gabs.ParseJSON([]byte(tt.args.Json))
			parsedExpected, _ := gabs.ParseJSON([]byte(tt.args.ExpectedJson))
			if parsedExpected == nil || parsedRequest == nil {
				t.Errorf("Invalid JSON passed for parsing.")
			} else {
				ok := isIncluded(parsedExpected, parsedRequest)

				if ok != tt.expectedOk {
					t.Errorf("Test failed")
				}
			}
		})
	}
	return
}

