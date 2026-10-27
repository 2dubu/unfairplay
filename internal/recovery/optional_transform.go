package recovery

// optionalTransform26A428 evaluates the semantic form of the optional outer
// hook for the runtime domain observed in the 26A428 enrollment profile.
//
// The lifted kernel accepts arbitrary register values, but the natural path
// reaches it with a fixed selector and a fixed upper input word. Keep those
// constraints explicit until another independent enrollment proves a wider
// contract.
func optionalTransform26A428(input, selector uint32) (uint32, bool) {
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
