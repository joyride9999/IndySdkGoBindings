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

import "indySDK/pool"

func SetPoolProtocolVersion(pb uint64) error {
	channel := pool.IndySetProtocolVersion(pb)
	result := <-channel
	return result.Error
}

func CreatePoolLedgerConfig(config pool.Pool) error {
	channel := pool.IndyCreatePoolLedgerConfig(config)
	result := <-channel
	return result.Error
}

func OpenPoolLedgerConfig(config pool.Pool) (int, error) {
	channel := pool.IndyOpenPoolLedger(config)
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
