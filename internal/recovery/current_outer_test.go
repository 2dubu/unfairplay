package recovery

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"testing"
)

func syntheticCurrentOuterInputs() ([]byte, []byte) {
	container := make([]byte, selectionSIDBContainerBytes)
	binary.BigEndian.PutUint32(container[:4], currentSIDBContainerVersion)
	for index := 4; index < len(container); index++ {
		container[index] = byte(index*37 + 11)
	}
	hardware := []byte{0x02, 0x13, 0x24, 0x35, 0x46, 0x57}
	return container, hardware
}

func syntheticExtendedCurrentOuterInputs() ([]byte, []byte) {
	container := make([]byte, selectionSIDBExtendedBytes)
	binary.BigEndian.PutUint32(container[:4], currentSIDBContainerVersion)
	for index := 4; index < len(container); index++ {
		container[index] = byte(index*29 + 7)
	}
	hardware := []byte{0x03, 0x14, 0x25, 0x36, 0x47, 0x58}
	return container, hardware
}

func syntheticStructuralCurrentOuterInputs(codecBytes, payloadBytes int) ([]byte, []byte) {
	container := make([]byte, 4+codecBytes+payloadBytes+selectionSIDBTrailerBytes)
	binary.BigEndian.PutUint32(container[:4], currentSIDBContainerVersion)
	for index := 4; index < len(container); index++ {
		container[index] = byte(index*17 + 3)
	}
	hardware := []byte{0x04, 0x15, 0x26, 0x37, 0x48, 0x59}
	return container, hardware
}

func TestCurrentOuterInputsAreTypedAndStructural(t *testing.T) {
	container, hardware := syntheticCurrentOuterInputs()
	table, state, control, err := currentOuterInputs(container, hardware)
	if err != nil {
		t.Fatal(err)
	}
	if state.Words != currentOuterInitialWords || state.OptionalWord != currentOuterInitialOptionalWord ||
		control != (outerIterationControl{X1: currentOuterInitialX1, CounterSeed: currentOuterInitialCounterSeed}) {
		t.Fatal("current outer constants differ")
	}
	if !bytes.Equal(state.Hardware[:], hardware) || !bytes.Equal(state.Codec[4:33], container[4:selectionSIDBPayloadOffset]) || !bytes.Equal(state.Codec[33:], make([]byte, 10)) {
		t.Fatal("current outer runtime inputs differ")
	}
	var encoded [outerTableWords * 4]byte
	for index, value := range table {
		binary.LittleEndian.PutUint32(encoded[index*4:], value)
	}
	if digest := fmt.Sprintf("%x", sha256.Sum256(encoded[:])); digest != "ca6a4b1c9df5c6d1f8e05d7d567a6e001e6b7c3f842ed4f7d7a5275f627ef9f5" {
		t.Fatalf("current outer table digest = %s", digest)
	}
}

func TestDeriveCurrentSelectionSourceDependencies(t *testing.T) {
	container, hardware := syntheticCurrentOuterInputs()
	baseline, iterations, err := deriveCurrentSelectionSource(container, hardware)
	if err != nil || iterations != outerNaturalIterationCount {
		t.Fatalf("derive current source = %d/%v", iterations, err)
	}
	if digest := fmt.Sprintf("%x", sha256.Sum256(baseline[:])); digest != "795392fd1df6f47b42fc755e0fc52348d5cf3516728b4c67e6ab9886847941b8" {
		t.Fatalf("synthetic current source digest = %s", digest)
	}

	payloadMutation := append([]byte(nil), container...)
	payloadMutation[selectionSIDBPayloadOffset+17] ^= 0x80
	payloadSource, _, err := deriveCurrentSelectionSource(payloadMutation, hardware)
	if err != nil || payloadSource != baseline {
		t.Fatal("SIDB encrypted payload changed the outer source")
	}

	codecMutation := append([]byte(nil), container...)
	codecMutation[4] ^= 1
	codecSource, _, err := deriveCurrentSelectionSource(codecMutation, hardware)
	if err != nil || codecSource == baseline {
		t.Fatal("SIDB codec header did not change the outer source")
	}

	hardwareMutation := append([]byte(nil), hardware...)
	hardwareMutation[0] ^= 1
	hardwareSource, _, err := deriveCurrentSelectionSource(container, hardwareMutation)
	if err != nil || hardwareSource == baseline {
		t.Fatal("hardware identity did not change the outer source")
	}
}

func TestDeriveExtendedCurrentSelectionSourceDependencies(t *testing.T) {
	container, hardware := syntheticExtendedCurrentOuterInputs()
	_, state, _, err := currentOuterInputs(container, hardware)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(state.Codec[4:], container[4:selectionSIDBExtendedPayloadOffset]) {
		t.Fatal("extended current outer codec differs")
	}
	baseline, iterations, err := deriveCurrentSelectionSource(container, hardware)
	if err != nil || iterations != selectionSIDBExtendedIterations {
		t.Fatalf("derive extended current source = %d/%v", iterations, err)
	}

	payloadMutation := append([]byte(nil), container...)
	payloadMutation[selectionSIDBExtendedPayloadOffset+17] ^= 0x80
	payloadSource, _, err := deriveCurrentSelectionSource(payloadMutation, hardware)
	if err != nil || payloadSource != baseline {
		t.Fatal("extended SIDB encrypted payload changed the outer source")
	}

	trailerMutation := append([]byte(nil), container...)
	trailerMutation[len(trailerMutation)-1] ^= 0x80
	trailerSource, _, err := deriveCurrentSelectionSource(trailerMutation, hardware)
	if err != nil || trailerSource != baseline {
		t.Fatal("extended SIDB trailer changed the outer source")
	}

	codecMutation := append([]byte(nil), container...)
	codecMutation[selectionSIDBPayloadOffset+2] ^= 1
	codecSource, _, err := deriveCurrentSelectionSource(codecMutation, hardware)
	if err != nil || codecSource == baseline {
		t.Fatal("extended SIDB codec tail did not change the outer source")
	}
}

func TestDeriveStructuralCurrentSelectionSource(t *testing.T) {
	container, hardware, layout := syntheticNativeSIDB(t, selectionSIDBPayloadBytes)
	baseline, _, iterations, err := deriveNativeSelectionSource(container, hardware)
	if err != nil || iterations != layout.OuterIterations {
		t.Fatalf("derive structural source = %d/%v", iterations, err)
	}
	payloadMutation := append([]byte(nil), container...)
	payloadMutation[layout.PayloadOffset+17] ^= 1
	unchanged, _, _, err := deriveNativeSelectionSource(payloadMutation, hardware)
	if err != nil || !bytes.Equal(unchanged, baseline) {
		t.Fatal("encrypted payload changed the outer source")
	}
}

func TestCurrentOuterInputsRejectUnknownGeometry(t *testing.T) {
	container, hardware := syntheticCurrentOuterInputs()
	unsupported := make([]byte, 4+12+64+selectionSIDBTrailerBytes)
	binary.BigEndian.PutUint32(unsupported[:4], currentSIDBContainerVersion)
	if _, _, _, err := currentOuterInputs(unsupported, hardware); err == nil {
		t.Fatal("SIDB with an undersized codec was accepted")
	}
	badVersion := append([]byte(nil), container...)
	badVersion[3] ^= 1
	if _, _, _, err := currentOuterInputs(badVersion, hardware); err == nil {
		t.Fatal("unknown SIDB version was accepted")
	}
	if _, _, _, err := currentOuterInputs(container, hardware[:5]); err == nil {
		t.Fatal("short hardware identity was accepted")
	}
}
