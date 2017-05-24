package config

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestLocalBasicTargetsParse(t *testing.T) {
	spec, err := ParseTargetSpecFile("testdata/targets/localbasic.hcl")
	if err != nil {
		t.Error(err)
	}
	if spec == nil {
		t.Error("Non-nil spec expected")
	}

	if spec.Name != "" {
		t.Error("Spec name should be empty")
	}

	m1 := spec.Machine["local-1"]
	if m1.Kind != KindLocal || m1.Destination != "" {
		t.Error("Incorrect data, got: ", spew.Sdump(spec.Machine))
	}

	m2 := spec.Machine["local-2"]
	// kind = 'local' is not specified, so it should be populated by normalizeMachineSpec()
	if m2.Kind != KindLocal || m2.Destination != "" {
		t.Error("Incorrect data, got: ", spew.Sdump(spec.Machine))
	}
}

func TestSSHBasicTargetsParse(t *testing.T) {
	spec, err := ParseTargetSpecFile("testdata/targets/sshbasic.hcl")
	if err != nil {
		t.Error(err)
	}
	if spec == nil {
		t.Error("Non-nil spec expected")
	}

	if spec.Name != "Frontend servers" {
		t.Error("Spec name is incorrect")
	}

	m1 := spec.Machine["frontend-1"]

	if m1.Kind != KindSSH || m1.Destination != "10.5.32.1" {
		t.Error("Incorrect data, got: ", spew.Sdump(spec.Machine))
	}

	// This also tests our validation, which sets the kind to AuthKindPassword
	if len(m1.Auth) != 1 || m1.Auth[0].Kind != AuthKindPassword || m1.Auth[0].Password != "1234" {
		t.Error("Incorrect auth data, got: ", spew.Sdump(spec.Machine))
	}
}

func TestBasicTargetsParseErrorCases(t *testing.T) {
	spec, err := ParseTargetSpecFile("testdata/targets/doesnt_exist.hcl")
	if err == nil {
		t.Error("Error expected")
	}
	if spec != nil {
		t.Error("nil spec expected")
	}

	spec, err = ParseTargetSpecFile("testdata/targets/invalid_hcl.hcl")
	if err == nil {
		t.Error("Error expected")
	}
	if spec != nil {
		t.Error("nil spec expected")
	}

	spec, err = ParseTargetSpecFile("testdata/targets/invalid_machinespec_kind.hcl")
	if err == nil {
		t.Error("Error expected")
	}
	if spec != nil {
		t.Error("nil spec expected")
	}
}
