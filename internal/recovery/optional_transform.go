package recovery

// optionalTransform evaluates the semantic form of the optional outer stage.
//
// The recovered relation accepts arbitrary words, but the natural path reaches
// it with a fixed selector and a fixed upper input word. Keep those
// constraints explicit until another independent enrollment proves a wider
// contract.
func optionalTransform(input, selector uint32) (uint32, bool) {
	const (
		inputBase     = uint32(0x976dec00)
		selectorValue = uint32(0x5d5c)
	)
	if selector != selectorValue || input&0xffffff00 != inputBase {
		return 0, false
	}

	value := input & 0xff
	var delta uint32
	switch {
	case value&2 == 0:
		delta = 0x182
	case value&0x0c != 0:
		delta = 0x17a
	case value&0x70 == 0x70:
		delta = 0x09a
	default:
		delta = 0x19a
	}
	return inputBase + ((value + delta) & 0x1ff), true
}

// outerOptionalActive is the semantic dispatch predicate. Values in the upper
// byte half select the optional stage;
// 0x0d is its single low-half control value.
func outerOptionalActive(input uint32) bool {
	return input&0xffffff00 == 0x976dec00 && (byte(input) >= 0x80 || byte(input) == 0x0d)
}
