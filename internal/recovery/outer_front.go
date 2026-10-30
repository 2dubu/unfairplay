package recovery

func applyOuterFront(table *[outerTableWords]uint32, base uint32, state *outerTableRoundState, x1, counterSeed uint32, hooks *outerHookSeeds) bool {
	if table == nil || state == nil || base >= outerTableWords || state.HardwareIndex >= 6 {
		return false
	}
	counter, ok := outerDigestInitialCounter(x1, counterSeed)
	if !ok {
		return false
	}
	digest := outerDigestState{
		X0: state.Words[0], X4: state.Words[3], X10: state.Words[2],
		X27: state.Words[1], X28: state.Words[4], Counter: counter,
		CodecSelector: 0x60319d93,
	}
	hardwareBase := (uint32(state.HardwareIndex) + 2) % 6
	var digestCodec [33]byte
	copy(digestCodec[:], state.Codec[:len(digestCodec)])
	if !applyOuterDigestAtHardwareBase(&digest, &state.Hardware, &digestCodec, hardwareBase) {
		return false
	}
	oldWord1, oldWord4 := state.Words[1], state.Words[4]
	words, word2Core, ok := outerFrontWords(table, base, &digest)
	if !ok {
		return false
	}
	derivedHooks, ok := deriveOuterHookSeedsWithCore(table, base, &digest, word2Core)
	if !ok {
		return false
	}
	if hooks != nil {
		*hooks = derivedHooks
	}

	prefix, ok := outerFrontSourcePrefix(table, base, words)
	if !ok {
		return false
	}
	tail, ok := outerFrontSourceTail(table, base, oldWord1, oldWord4)
	if !ok {
		return false
	}
	state.Words = words
	copy(state.Source[:12], prefix[:])
	copy(state.Source[12:16], tail[:])
	return true
}
