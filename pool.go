/*
// ******************************************************************
// Purpose: exported public functions that handles pool functions
// from libindy
// Author:  alexandru.leonte@siemens.com
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
	"github.com/joyride9999/IndySdkGoBindings/pool"
	"encoding/json"
	"errors"
	"unsafe"
)

func SetPoolProtocolVersion(pb uint64) error {
	channel := pool.IndySetProtocolVersion(pb)
	result := <-channel
	return result.Error
}

func CreatePoolLedgerConfig(config pool.Pool) error {

	poolName := config.Name
	upPoolName := unsafe.Pointer(C.CString(poolName))
	defer C.free(upPoolName)

	poolConfig, err := json.Marshal(config)
	if err != nil {
		return errors.New("cant read json")
	}
	upPoolCfg := unsafe.Pointer(C.CString(string(poolConfig)))
	defer C.free(upPoolCfg)

	channel := pool.IndyCreatePoolLedgerConfig(upPoolName, upPoolCfg)
	result := <-channel
	return result.Error
}

func OpenPoolLedgerConfig(config pool.Pool) (int, error) {

	poolName := config.Name
	upPoolName := unsafe.Pointer(C.CString(poolName))
	defer C.free(upPoolName)

	type cfg struct {
		Timeout int `json:"timeout"`
	}

	t := cfg{
		Timeout: 10,
	}

	jsonConfig, err := json.Marshal(t)
	if err != nil {
		return 0, errors.New("cant read json")
	}
	upPoolCfg := unsafe.Pointer(C.CString(string(jsonConfig)))
	defer C.free(upPoolCfg)

	channel := pool.IndyOpenPoolLedger(upPoolName, upPoolCfg)
	result := <-channel
	if result.Error != nil {
		return 0, result.Error
	}
	return result.Results[0].(int), result.Error
}
func ClosePoolHandle(ph int) error {
	channel := pool.IndyClosePoolHandle(ph)
	result := <-channel
	return result.Error
}
