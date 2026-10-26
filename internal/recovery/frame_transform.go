package recovery

import (
	"encoding/hex"
	"fmt"
)

// frameRoundKeys are build-global constants, not account material. They are
// the ten independent AddRoundKey values recovered from the 26A428 frame
// mixer. This is an AES-shaped SPN, not standard AES: every round includes
// MixColumns and the keys do not use the AES schedule.
var frameRoundKeys = [10][16]byte{
	{0x45, 0xda, 0x26, 0xeb, 0x35, 0xe7, 0x30, 0x12, 0x1e, 0x00, 0x0c, 0xb3, 0x8f, 0xcb, 0xc0, 0x30},
	{0xf9, 0x6f, 0x5c, 0x1a, 0x51, 0xad, 0x34, 0xa5, 0x18, 0x91, 0x15, 0x60, 0xca, 0xc6, 0x70, 0x9c},
	{0x4f, 0x3e, 0x82, 0x6e, 0x1e, 0x93, 0xb6, 0xcb, 0x06, 0x02, 0xa3, 0xab, 0xcc, 0xc4, 0xd3, 0x37},
	{0x57, 0x58, 0x18, 0x25, 0x49, 0xcb, 0xae, 0xee, 0x4f, 0xc9, 0x0d, 0x45, 0x83, 0x0d, 0xde, 0x72},
	{0x88, 0x45, 0x58, 0xc9, 0xc1, 0x8e, 0xf6, 0x27, 0x8e, 0x47, 0xfb, 0x62, 0x0d, 0x4a, 0x25, 0x10},
	{0x4e, 0x7a, 0x92, 0x1e, 0x8f, 0xf4, 0x64, 0x39, 0x01, 0xb3, 0x9f, 0x5b, 0x0c, 0xf9, 0xba, 0x4b},
	{0xf7, 0x8e, 0x21, 0xe0, 0x78, 0x7a, 0x45, 0xd9, 0x79, 0xc9, 0xda, 0x82, 0x75, 0x30, 0x60, 0xc9},
	{0xb3, 0x5e, 0xfc, 0x7d, 0xcb, 0x24, 0xb9, 0xa4, 0xb2, 0xed, 0x63, 0x26, 0xc7, 0xdd, 0x03, 0xef},
	{0xf2, 0x25, 0x23, 0xbb, 0x39, 0x01, 0x9a, 0x1f, 0x8b, 0xec, 0xf9, 0x39, 0x4c, 0x31, 0xfa, 0xd6},
	{0x2e, 0x08, 0xd5, 0x92, 0x17, 0x09, 0x4f, 0x8d, 0x9c, 0xe5, 0xb6, 0xb4, 0xd0, 0xd4, 0x4c, 0x62},
}

var frameBoundaryConstants = [4]uint32{0x195f4092, 0x26446852, 0xdcadcad2, 0x45bf86c4}

var frameBoundaryDeltas = [32]uint32{
	0xafee0942, 0x3e68be6a, 0x5f453292, 0x2474a137, 0xdc28f1ed, 0x2d568a06, 0xc55ebddd, 0xc36ef968,
	0xcd8dd239, 0x28e5abe6, 0x02a084ee, 0x9fd4fe38, 0x20fc1ab8, 0xe9aa26af, 0xf4f4c46c, 0x2f9ad64f,
	0x01f1d3d8, 0xf2fd788b, 0x8542fcff, 0xa3030226, 0xe0f818f3, 0xed253e4d, 0xc101fa3f, 0xf6e82c3b,
	0xb87bcf16, 0x916f709d, 0x61329162, 0x7ba44844, 0xf85b45b7, 0xf109cdfa, 0xe8d868c5, 0x0643affe,
}

var frameBoundaryRegisters = [4]uint8{9, 10, 8, 11}

var frameBoundaryState = [31]uint64{
	0x13aa1aa0, 0x85f8e4db, 0xd196065f, 0x4a, 0x100df7e20, 0x100facd50, 0xe01f5560, 0x62af52c0,
	0, 0, 0, 0, 0x100dffd60, 0x100d43b60, 0xe19bdcb8, 0xe1eaa86f,
	0x60de3e0b, 0x7ff05555, 0, 0x52bcef82, 0x29, 0xe2c9261f, 0xe2, 0x100d61250,
	0x8d358793, 0x13aa1aa0, 0, 0xbb9f5d0e, 0x69e0c67a, 0x7000000ffdd0, 0x867db14e,
}

const frameBoundaryStackHex = "000000000000000000000000a4b70000663150500000000066000000000000009091fe000100000061000000000000009077d30001000000c66d0d7c0000000060fddf00010000009073f40001000000000000009387358dbebcd4b3505ba7679a2f2d05000000009a000000000000009096fe00010000003098f400010000003096f40001000000405b07010100000000000000795f00009094fe000100000024b70f0000700000f01eff00010000000000000000000000c8fc0f0000700000d086ef0001000000f029de00010000009092fe0001000000bc000000010000004eb17d8601000000"

var frameBoundaryStack = func() []byte {
	value, err := hex.DecodeString(frameBoundaryStackHex)
	if err != nil {
		panic(err)
	}
	return value
}()

var frameSBox = func() (result [256]byte) {
	for input := range 256 {
		inverse := frameGFInverse(byte(input))
		rotated, value := inverse, inverse
		for range 4 {
			rotated = rotated<<1 | rotated>>7
			value ^= rotated
		}
		result[input] = value ^ 0x63
	}
	return result
}()

func frameGFMul(left, right byte) byte {
	var result byte
	for right != 0 {
		if right&1 != 0 {
			result ^= left
		}
		high := left & 0x80
		left <<= 1
		if high != 0 {
			left ^= 0x1b
		}
		right >>= 1
	}
	return result
}

func frameGFInverse(value byte) byte {
	if value == 0 {
		return 0
	}
	result, base, exponent := byte(1), value, 254
	for exponent != 0 {
		if exponent&1 != 0 {
			result = frameGFMul(result, base)
		}
		base = frameGFMul(base, base)
		exponent >>= 1
	}
	return result
}

func frameCore(input [16]byte) [16]byte {
	state := input
	shift := [16]byte{0, 5, 10, 15, 4, 9, 14, 3, 8, 13, 2, 7, 12, 1, 6, 11}
	for _, key := range frameRoundKeys {
		var substituted, mixed [16]byte
		for target, source := range shift {
			substituted[target] = frameSBox[state[source]^key[source]]
		}
		for column := range 4 {
			offset := column * 4
			a, b, c, d := substituted[offset], substituted[offset+1], substituted[offset+2], substituted[offset+3]
			mixed[offset] = frameGFMul(2, a) ^ frameGFMul(3, b) ^ c ^ d
			mixed[offset+1] = a ^ frameGFMul(2, b) ^ frameGFMul(3, c) ^ d
			mixed[offset+2] = a ^ b ^ frameGFMul(2, c) ^ frameGFMul(3, d)
			mixed[offset+3] = frameGFMul(3, a) ^ b ^ c ^ frameGFMul(2, d)
		}
		state = mixed
	}
	return state
}

func frameBoundaryWord(column int, input []byte) uint32 {
	result := frameBoundaryConstants[column]
	for bit, delta := range frameBoundaryDeltas {
		if input[bit/8]&(byte(1)<<uint(bit%8)) != 0 {
			result ^= delta
		}
	}
	return result
}

func runFrameSuffix(s *state) error {
	for blockAddress := uint64(0x1002fd800); ; {
		block, ok := compact26A428SelectionBlocks[blockAddress]
		if !ok {
			return fmt.Errorf("frame suffix lacks block %#x", blockAddress)
		}
		start, end := 0, len(block.instructions)
		if blockAddress == 0x1002fd800 {
			start = int((uint64(0x1002fdb64) - blockAddress) / 4)
		}
		if blockAddress == 0x1002ff800 {
			end = int((uint64(0x1002ff8d8) - blockAddress) / 4)
		}
		if err := runCompactKernel(s, compactKernel{start: blockAddress + uint64(start)*4, instructions: block.instructions[start:end]}); err != nil {
			return err
		}
		if blockAddress == 0x1002ff800 {
			return nil
		}
		if block.terminator != compactBlockDirect {
			return fmt.Errorf("frame suffix block %#x is not direct", blockAddress)
		}
		blockAddress = block.target
	}
}

func frameTransform(value *model, input [16]byte) ([16]byte, error) {
	worker := newState(value)
	worker.reset()
	worker.sp = 0x7000000eea20
	worker.n, worker.z, worker.c, worker.v = false, true, false, false
	worker.x = frameBoundaryState
	if err := worker.store(0x7000000eea20, frameBoundaryStack); err != nil {
		return [16]byte{}, err
	}
	core := frameCore(input)
	for column, register := range frameBoundaryRegisters {
		worker.x[register] = uint64(frameBoundaryWord(column, core[column*4:column*4+4]))
	}
	if err := runFrameSuffix(worker); err != nil {
		return [16]byte{}, err
	}
	var result [16]byte
	bytes, err := worker.outputBytes(0x7000000fb724, len(result))
	copy(result[:], bytes)
	return result, err
}
