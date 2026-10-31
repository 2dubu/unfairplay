package main

import (
	"testing"

	"github.com/DrWhax/unfairplay/internal/recovery"
)

func testRecoveryLayout() recovery.LayoutEvidence {
	return recovery.LayoutEvidence{
		ContainerBytes: 991, CodecBytes: 27, PayloadBytes: 944, SourceBytes: 224,
		SelectorModulus: 12, SelectedSlot: 5, RootCursor: 2, RootLength: 942,
		RootChildCount: 2, TreeNodes: 17, RegistryRecords: 3,
		RootVersion: 0x10001, AccountVersion: 0x10002,
		TypeCursor: 100, TypeLength: 24, SignCursor: 124, SignLength: 148,
		ValueCursor: 272, ValueLength: 36, SelectedRecordIndex: 2,
		SelectedPayloadOffset: 292, SelectedWindowOffset: 0, SINFSelector: "12345678:00000002",
	}
}

func TestValidateManifest(t *testing.T) {
	hash := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	iterations, candidates := 123, 1
	layout := testRecoveryLayout()
	manifest := qaManifest{Schema: qaSchema, Controls: []qaControl{{
		Name: "control", Capture: "/private/capture", ExpectedModelSHA256: hash,
		ExpectedStateSHA256: hash, ExpectedRecoveryIterations: &iterations, ExpectedRegistryCandidates: &candidates,
		ExpectedRecoveryLayout: &layout,
	}}}
	if err := validateManifest(manifest); err != nil {
		t.Fatal(err)
	}
	manifest.Controls = []qaControl{{Name: "verify", IPA: "/private/App.ipa", Model: "/private/model.bin"}}
	if err := validateManifest(manifest); err != nil {
		t.Fatal(err)
	}
	manifest.Controls[0].Capture = "/private/capture"
	if err := validateManifest(manifest); err == nil {
		t.Fatal("accepted mixed recover and verify control")
	}
	manifest.Controls[0].Capture = ""
	manifest.Controls[0].ExpectedRecoveryIterations = &iterations
	if err := validateManifest(manifest); err == nil {
		t.Fatal("accepted recovery geometry on verify control")
	}
	manifest.Controls[0].ExpectedRecoveryIterations = nil
	manifest.Controls = append(manifest.Controls, manifest.Controls[0])
	if err := validateManifest(manifest); err == nil {
		t.Fatal("accepted duplicate QA control")
	}
}

func TestCheckControl(t *testing.T) {
	hash := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	verified, iterations, candidates := true, 123, 1
	layout := testRecoveryLayout()
	control := qaControl{
		Name: "control", ExpectedModelSHA256: hash, ExpectedIPASHA256: hash,
		ExpectedReferenceSource: "typed-go",
		ExpectedImages:          2, ExpectedPages: 3, ExpectedPlaintextVerified: &verified,
		ExpectedRecoveryIterations: &iterations, ExpectedStateSHA256: hash, ExpectedRegistryCandidates: &candidates,
		ExpectedRecoveryLayout: &layout,
	}
	result := recoveryResult{
		Success: true, DecryptedIPASHA256: hash,
		ReferenceCodeSource: "typed-go", EncryptedImages: 2, DecryptedPages: 3, PlaintextVerified: true,
		RecoveryIterations: 123, RecoveryStateSHA256: hash, DecodedRegistryCandidates: 1,
		RecoveryLayout: &layout,
	}
	if err := checkControl(control, result, hash); err != nil {
		t.Fatal(err)
	}
	result.DecryptedPages++
	if err := checkControl(control, result, hash); err == nil {
		t.Fatal("accepted changed QA result")
	}
}
