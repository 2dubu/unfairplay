package recovery

// semanticBuildProfile contains only typed enrollment data. It has no native
// addresses, registers, stack layout, or control-flow metadata.
type semanticBuildProfile struct {
	currentSINFWhite [16]byte
	currentSINFInput [4096]byte
}

var currentSINFWhitening = [16]byte{
	0x92, 0xd4, 0x61, 0xb7, 0x84, 0xf8, 0xb5, 0x95,
	0x06, 0x7e, 0x68, 0xc6, 0x9d, 0x6b, 0xbf, 0x8a,
}

func portableCurrentSemanticProfile() *semanticBuildProfile {
	result := &semanticBuildProfile{currentSINFWhite: currentSINFWhitening}
	for position, mask := range result.currentSINFWhite {
		for value := range 256 {
			result.currentSINFInput[position*256+value] = byte(value) ^ mask
		}
	}
	return result
}
