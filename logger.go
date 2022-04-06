/*
// ******************************************************************
// Purpose: exported public functions that handles logger functions
// from libindy
// Author:  alexandru.leonte@siemens.com
// Notes:
// Copyright (c): Siemens SRL
// This work is licensed under the terms of the Apache License Version 2.0.  See
// the LICENSE.txt file in the top-level directory.
// ******************************************************************
*/

package indySDK

import "github.com/joyride9999/IndySdkGoBindings/logger"

func IndySetLogger() {
	logger.IndySetLogger()
}
