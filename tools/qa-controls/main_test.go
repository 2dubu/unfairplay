package main

import "testing"

func TestValidateManifest(t *testing.T) {
	hash := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	iterations, candidates := 123, 1
	manifest := qaManifest{Schema: qaSchema, Controls: []qaControl{{
		Name: "control", Capture: "/private/capture", ExpectedModelSHA256: hash,
		ExpectedStateSHA256: hash, ExpectedRecoveryIterations: &iterations, ExpectedRegistryCandidates: &candidates,
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
	control := qaControl{
		Name: "control", ExpectedModelSHA256: hash, ExpectedIPASHA256: hash,
		ExpectedReferenceSource: "typed-go",
		ExpectedImages:          2, ExpectedPages: 3, ExpectedPlaintextVerified: &verified,
		ExpectedRecoveryIterations: &iterations, ExpectedStateSHA256: hash, ExpectedRegistryCandidates: &candidates,
	}
	result := recoveryResult{
		Success: true, DecryptedIPASHA256: hash,
		ReferenceCodeSource: "typed-go", EncryptedImages: 2, DecryptedPages: 3, PlaintextVerified: true,
		RecoveryIterations: 123, RecoveryStateSHA256: hash, DecodedRegistryCandidates: 1,
	}
	if err := checkControl(control, result, hash); err != nil {
		t.Fatal(err)
	}
	result.DecryptedPages++
	if err := checkControl(control, result, hash); err == nil {
		t.Fatal("accepted changed QA result")
	}
}
