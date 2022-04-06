/*
// ******************************************************************
// Purpose: Wrapper to call libindy, imports functions from indy_anoncreds.h
// Author:  alexandru.leonte@siemens.com
// Notes:
// Copyright (c): Siemens SRL
// This work is licensed under the terms of the Apache License Version 2.0.  See
// the LICENSE.txt file in the top-level directory.
// ******************************************************************
*/

package anoncreds

import "time"

type CredentialConfig struct {
	SupportsRevocation bool `json:"support_revocation"`
}

type RevocRegConfig struct {
	MaxCredNumber int    `json:"max_cred_num"`
	IssuanceType  string `json:"issuance_type"`
}

// CredentialDefinitionInfo - helper structure to hold data about credential definition
type CredentialDefinitionInfo struct {
	CredentialDefinitionId   string
	CredentialDefinitionJson string
	RevocationRegistryId     string
	SchemaJson               string
}

// CredentialInfo - helper structure to hold data about credential
type CredentialInfo struct {
	Id                     int64
	IssuerDid              string
	SubjectDid             string
	CredentialDefinitionId string
	MasterSecretId         string
	CredentialId           string
	RevocRegId             string
	CredRevocId            string
	Valid                  bool
	CreationDate           time.Time
	RevocationDate         *time.Time
}
