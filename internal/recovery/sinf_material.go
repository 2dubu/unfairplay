package recovery

// The final native material writer applies this build-global byte mask to the
// selected 16-byte precursor. It contains no account material.
var currentSINFMaterialMask = [16]byte{
	0xa3, 0x14, 0x69, 0x2b, 0x4e, 0x24, 0xba, 0x48,
	0xfc, 0x07, 0x9a, 0x34, 0x7e, 0xf2, 0x2a, 0xe5,
}

func currentSINFMaterialFromPrecursor(precursor [16]byte) (material [16]byte) {
	for index := range material {
		material[index] = precursor[index] ^ currentSINFMaterialMask[index]
	}
	return material
}

// The enrollment pack applies generatedCurrentSINFKeyMask to the material.
// Both masks differ by 0x1b at every position, so the final model key has this
// simpler direct projection from the selected precursor.
func currentSINFModelKeyFromPrecursor(precursor [16]byte) (key [16]byte) {
	for index := range key {
		key[index] = precursor[index] ^ 0x1b
	}
	return key
}
