package aws

const filetemplate = `
//go:build !ignore_autogenerated
// +build !ignore_autogenerated

package main

import cco "github.com/openshift/cloud-credential-operator/pkg/apis/cloudcredential/v1"

type IAMPolicy struct {
	Version   string
	Statement []cco.StatementEntry
}

func GetIAMPolicy() IAMPolicy {
	return IAMPolicy{}
}
`
