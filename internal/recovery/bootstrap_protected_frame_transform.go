package recovery

import "fmt"

const (
	bootstrapDescriptorSetup      = uint64(0x1007b0ce0)
	bootstrapProtectedFrameStart  = uint64(0x1007b2e8c)
	bootstrapProtectedFrameFinish = uint64(0x1007b3000)
	bootstrapProtectedFrameReturn = uint64(0x1007b30e4)
	bootstrapStateBase            = uint64(0x101075000)
	bootstrapTargetTable          = uint64(0x10105b6f0)
)

func runBootstrapDescriptorSetup(s *state) (uint64, error) {
	s.traceBlock(bootstrapDescriptorSetup)
	s.sp -= 0x60
	for _, pair := range [...]struct {
		first, second int
		offset        uint64
	}{{28, 27, 0}, {26, 25, 0x10}, {24, 23, 0x20}, {22, 21, 0x30}, {20, 19, 0x40}, {29, 30, 0x50}} {
		if err := s.storeValue(s.sp+pair.offset, s.x[pair.first], 8); err != nil {
			return 0, fmt.Errorf("bootstrap descriptor-setup register %d: %w", pair.first, err)
		}
		if err := s.storeValue(s.sp+pair.offset+8, s.x[pair.second], 8); err != nil {
			return 0, fmt.Errorf("bootstrap descriptor-setup register %d: %w", pair.second, err)
		}
	}
	s.x[29] = s.sp + 0x50
	s.sp -= 0x500
	s.x[8] = 0x10106d000
	value, err := s.load(s.x[8]+0x408, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-setup guard pointer: %w", err)
	}
	s.x[8] = value
	value, err = s.load(value, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-setup guard: %w", err)
	}
	s.x[8] = value
	if err := s.storeValue(addOffset(s.x[29], -0x60), value, 8); err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-setup saved guard: %w", err)
	}
	s.x[9] = bootstrapStateBase
	word, err := s.load(s.x[9]+0x428, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-setup first state word: %w", err)
	}
	s.x[8] = uint64(uint32(word))
	s.x[10] = bootstrapStateBase
	word, err = s.load(s.x[10]+0xb40, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-setup second state word: %w", err)
	}
	s.x[14] = uint64(uint32(word))
	s.x[12] = 0xa6
	s.x[8] = uint64(uint32(s.x[14]) + uint32(s.x[8]))
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[12]))
	s.x[11] = 0x0a5cd4cf
	s.x[8] = uint64(uint32(s.x[8]) * uint32(s.x[11]))
	s.x[13] = 0x100fe9390 + uint64(byte(s.x[8]))
	value, err = s.load(s.x[13], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-setup first byte: %w", err)
	}
	s.x[13] = uint64(byte(value))
	s.x[15] = 0x28
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[15]))
	s.x[15] = 0x100de2af0
	value, err = s.load(s.x[15]+uint64(uint32(s.x[13])), 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-setup second byte: %w", err)
	}
	s.x[13] = uint64(byte(value))
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[13]))
	s.x[8] &= 0xff
	s.x[8] = uint64(uint32(s.x[8]) - 0x72)
	s.x[13] = bootstrapTargetTable
	value, err = loadBootstrapPointer(s, s.x[13], uint32(s.x[8]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-setup first target: %w", err)
	}
	s.x[8] = value
	s.x[15] = 0xffffffff82c2ff2d
	s.x[3] = s.x[8] + s.x[15]
	value, err = s.load(s.x[3], 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-setup first protected state: %w", err)
	}
	s.x[8] = value
	s.x[14] = uint64(uint32(s.x[14]) + uint32(s.x[8]))
	s.x[14] = uint64(uint32(s.x[14]) ^ uint32(s.x[12]))
	s.x[14] = uint64(uint32(s.x[14]) * uint32(s.x[11]))
	value, err = s.load(s.x[13]+0x5e20, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-setup second table: %w", err)
	}
	s.x[15] = value
	s.x[15] -= 0x9fd000
	s.x[15] -= 0xcfb
	if err := s.storeValue(s.sp+0x90, s.x[15], 8); err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-setup saved second table: %w", err)
	}
	s.x[15] += uint64(byte(s.x[14]))
	value, err = s.load(s.x[15], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-setup third byte: %w", err)
	}
	s.x[15] = uint64(byte(value)) ^ 0xf0
	s.x[16] = s.x[13] + 0x9898
	value, err = s.load(s.x[16], 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-setup third table: %w", err)
	}
	s.x[16] = value
	s.x[17] = 0xfffffffff3e6c916
	s.x[16] += s.x[17]
	if err := s.storeValue(s.sp+0x88, s.x[16], 8); err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-setup saved third table: %w", err)
	}
	value, err = s.load(s.x[16]+uint64(uint32(s.x[15])), 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-setup fourth byte: %w", err)
	}
	s.x[15] = uint64(byte(value))
	s.x[14] = uint64(uint32(s.x[15]) ^ uint32(s.x[14]))
	s.x[15] = 0x1314
	s.x[14] = uint64(uint32(s.x[15]) + uint32(byte(s.x[14])))
	value, err = s.load(s.x[13]+uint64(uint32(s.x[14]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-setup second target: %w", err)
	}
	s.x[14] = value
	s.x[15] = 0xffffffff9608a825
	s.x[4] = s.x[14] + s.x[15]
	value, err = s.load(s.x[4], 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-setup second protected state: %w", err)
	}
	s.x[14] = value
	s.x[15] = s.sp + 0xc0
	s.x[8] ^= s.x[15]
	s.x[8] ^= s.x[14]
	s.x[14] = s.x[8] * s.x[11]
	s.x[15] = 0x73d046208eaed1a6
	s.x[14] ^= s.x[15]
	s.x[8] ^= s.x[15]
	s.x[8] *= s.x[11]
	if err := s.storeValue(s.x[3], s.x[14], 8); err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-setup first protected output: %w", err)
	}
	if err := s.storeValue(s.x[4], s.x[8], 8); err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-setup second protected output: %w", err)
	}
	value, err = s.load(s.x[1], 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-setup descriptor length: %w", err)
	}
	s.x[8] = value + 0x40
	if err := s.storeValue(s.x[1], s.x[8], 8); err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-setup descriptor length output: %w", err)
	}
	s.x[8] = uint64(uint32(s.x[0]) & 0x7cbed57f)
	s.x[8] = uint64(uint32(s.x[8]) << 1)
	s.x[14] = 0x03412a80
	s.x[8] = uint64(uint32(s.x[8]) - uint32(s.x[0]))
	s.x[8] = uint64(uint32(s.x[8]) + uint32(s.x[14]))
	s.x[14] = 0x79d36d76
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[14]))
	s.x[14] = 0x31b6b333
	word, err = s.load(s.x[0]+8, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-setup descriptor word: %w", err)
	}
	s.x[15] = uint64(uint32(word))
	s.x[17] = uint64(uint32(s.x[8])*uint32(s.x[14]) + uint32(s.x[15]))
	value, err = s.load(s.x[0]+0x20, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-setup descriptor pointer: %w", err)
	}
	s.x[8] = value
	if err := s.storeValue(s.sp+0x98, value, 8); err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-setup saved descriptor pointer: %w", err)
	}
	word, err = s.load(value+4, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-setup descriptor status: %w", err)
	}
	s.x[8] = uint64(uint32(word))
	s.x[14] = 0x842a8d93
	s.x[15] = s.addWithCarry(uint64(uint32(s.x[8])), uint64(uint32(s.x[14])), false, 32)
	s.x[14] = 0x7bd5726d
	s.x[8] = uint64(uint32(s.x[14]) - uint32(s.x[8]))
	if s.condition("pl") {
		s.x[8] = uint64(uint32(s.x[15]))
	}
	s.x[15] = 0x82
	if err := s.storeValue(s.sp+0xb0, uint64(uint32(s.x[15])), 4); err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-setup saved status: %w", err)
	}
	s.x[16] = s.sp + 0x2d0
	if err := s.storeValue(s.sp+0xb8, s.x[16], 8); err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-setup saved buffer: %w", err)
	}
	if err := s.storeValue(s.sp+0xa0, uint64(uint32(s.x[15])), 4); err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-setup saved pair first: %w", err)
	}
	if err := s.storeValue(s.sp+0xa4, uint64(uint32(s.x[14])), 4); err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-setup saved pair second: %w", err)
	}
	s.x[15] = s.sp + 0xc8
	if err := s.storeValue(s.sp+0xa8, s.x[15], 8); err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-setup saved output: %w", err)
	}
	value, err = s.load(s.x[0]+0x10, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-setup first descriptor value: %w", err)
	}
	s.x[5] = value
	value, err = s.load(s.x[0]+0x18, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-setup second descriptor value: %w", err)
	}
	s.x[15] = value
	if err := s.storeValue(s.sp+0x30, value, 8); err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-setup saved descriptor value: %w", err)
	}
	word, err = s.load(value+4, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-setup comparison word: %w", err)
	}
	s.x[16] = uint64(uint32(word))
	s.compare32(uint32(s.x[16]), uint32(s.x[14]))
	gate := uint32(boolValue(s.condition("eq")))
	s.x[1] = uint64(gate)
	s.x[15] = 0xffff99d7
	s.x[2] = 0x57f3
	s.x[15] = uint64(uint32(s.x[17]) + uint32(s.x[15]))
	s.x[2] = uint64(uint32(s.x[15]) ^ uint32(s.x[2]))
	s.x[1] = uint64(gate * uint32(s.x[2]))
	s.x[17] = uint64(uint32(s.x[1]) ^ uint32(s.x[17]))
	word, err = s.load(s.x[9]+0x428, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-setup refreshed first word: %w", err)
	}
	s.x[9] = uint64(uint32(word))
	word, err = s.load(s.x[10]+0xb40, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-setup refreshed second word: %w", err)
	}
	s.x[10] = uint64(uint32(word))
	s.x[9] = uint64(uint32(s.x[10]) + uint32(s.x[9]))
	s.x[9] = uint64(uint32(s.x[9]) ^ uint32(s.x[12]))
	s.x[10] = 0x100fe9190
	s.x[9] = uint64(uint32(s.x[9]) * uint32(s.x[11]))
	s.x[10] += uint64(byte(s.x[9]))
	value, err = s.load(s.x[10], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-setup fifth byte: %w", err)
	}
	s.x[10] = uint64(byte(value))
	s.x[11] = 0xfa
	s.x[12] = 0x100de28f0
	s.x[10] = uint64(uint32(s.x[10]) ^ uint32(s.x[11]))
	value, err = s.load(s.x[12]+uint64(uint32(s.x[10])), 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-setup sixth byte: %w", err)
	}
	s.x[10] = uint64(byte(value))
	s.x[9] = uint64(uint32(s.x[9]) ^ uint32(s.x[10]))
	s.x[9] &= 0xff
	s.x[9] = uint64(uint32(s.x[9]) + 0x231)
	value, err = s.load(s.x[13]+uint64(uint32(s.x[9]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-setup final table: %w", err)
	}
	s.x[9] = value
	s.x[10] = 0xffffffff98245355
	s.x[25] = s.x[9] + s.x[10]
	next, err := loadBootstrapPointer(s, s.x[25], uint32(s.x[17]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-setup target: %w", err)
	}
	s.x[9] = next
	s.x[17] = 0x70d4
	return next, nil
}

func runBootstrapProtectedFrame(s *state) (uint64, error) {
	s.traceBlock(bootstrapProtectedFrameStart)
	s.sp -= 0x30
	for _, pair := range [...]struct {
		first, second int
		offset        uint64
	}{{22, 21, 0}, {20, 19, 0x10}, {29, 30, 0x20}} {
		if err := s.storeValue(s.sp+pair.offset, s.x[pair.first], 8); err != nil {
			return 0, fmt.Errorf("bootstrap protected-frame register %d: %w", pair.first, err)
		}
		if err := s.storeValue(s.sp+pair.offset+8, s.x[pair.second], 8); err != nil {
			return 0, fmt.Errorf("bootstrap protected-frame register %d: %w", pair.second, err)
		}
	}
	s.x[29] = s.sp + 0x20
	s.sp -= 0x10
	s.x[19] = s.x[1]
	s.x[8] = bootstrapStateBase
	firstWord, err := s.load(s.x[8]+0x2d0, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap protected-frame first word: %w", err)
	}
	s.x[11] = uint64(uint32(firstWord))
	s.x[9] = bootstrapStateBase
	secondWord, err := s.load(s.x[9]+0xb40, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap protected-frame second word: %w", err)
	}
	s.x[13] = uint64(uint32(secondWord))
	s.x[10] = 0xa6
	s.x[11] = uint64(uint32(s.x[13]) - uint32(s.x[11]))
	s.x[12] = uint64(uint32(s.x[11]) ^ uint32(s.x[10]))
	s.x[11] = 0x0a5cd4cf
	s.x[12] = uint64(uint32(s.x[12]) * uint32(s.x[11]))
	s.x[14] = s.x[12] & 0xff
	s.x[15] = 0x100fe9190
	value, err := s.load(s.x[15]+s.x[14], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap protected-frame first byte: %w", err)
	}
	s.x[14] = uint64(byte(value))
	s.x[15] = 0x29
	s.x[14] = uint64(uint32(s.x[14]) ^ uint32(s.x[15]))
	s.x[15] = 0x100de28f0
	value, err = s.load(s.x[15]+uint64(uint32(s.x[14])), 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap protected-frame second byte: %w", err)
	}
	s.x[14] = uint64(byte(value))
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[14]))
	s.x[12] &= 0xff
	s.x[14] = uint64(uint32(s.x[12]) + 0xe31)
	s.x[12] = bootstrapTargetTable
	value, err = s.load(s.x[12]+uint64(uint32(s.x[14]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap protected-frame first target: %w", err)
	}
	s.x[14] = value
	s.x[15] = 0xffffffff862b81bd
	value, err = s.load(s.x[14]+s.x[15], 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap protected-frame first state: %w", err)
	}
	s.x[16] = value
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[10]))
	s.x[13] = uint64(uint32(s.x[16]) ^ uint32(s.x[13]))
	s.x[17] = uint64(uint32(s.x[13]) * uint32(s.x[11]))
	value, err = s.load(s.x[12]+0x1af8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap protected-frame second table: %w", err)
	}
	s.x[13] = value + uint64(byte(s.x[17]))
	s.x[1] = 0xffffffffcd114862
	value, err = s.load(s.x[13]+s.x[1], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap protected-frame third byte: %w", err)
	}
	s.x[1] = uint64(byte(value))
	s.x[13] = 0xcd
	s.x[1] = uint64(uint32(s.x[1]) ^ uint32(s.x[13]))
	value, err = s.load(s.x[12]+0x2a58, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap protected-frame third table: %w", err)
	}
	s.x[2] = value
	s.x[1] += value
	s.x[2] = 0xffffffffe5b0614e
	value, err = s.load(s.x[1]+s.x[2], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap protected-frame fourth byte: %w", err)
	}
	s.x[1] = uint64(byte(value))
	s.x[17] = uint64(uint32(s.x[1]) ^ uint32(s.x[17]))
	s.x[1] = 0x12bf
	s.x[17] = uint64(uint32(s.x[1]) + uint32(byte(s.x[17])))
	value, err = s.load(s.x[12]+uint64(uint32(s.x[17]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap protected-frame second target: %w", err)
	}
	s.x[17] = value
	s.x[1] = 0xffffffff9608a825
	value, err = s.load(s.x[17]+s.x[1], 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap protected-frame second state: %w", err)
	}
	s.x[2] = value
	s.x[16] ^= value
	s.x[2] = s.x[29] - 0x28
	s.x[3] = 0x73d046208eaed1a6
	s.x[16] += s.x[2]
	s.x[2] = s.x[16] ^ s.x[3]
	s.x[16] = s.x[16]*s.x[11] + s.x[3]
	s.x[2] *= s.x[11]
	if err := s.storeValue(s.x[14]+s.x[15], s.x[16], 8); err != nil {
		return 0, fmt.Errorf("bootstrap protected-frame first state output: %w", err)
	}
	if err := s.storeValue(s.x[17]+s.x[1], s.x[2], 8); err != nil {
		return 0, fmt.Errorf("bootstrap protected-frame second state output: %w", err)
	}
	value, err = s.load(s.x[19], 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap protected-frame input: %w", err)
	}
	s.x[21] = value
	s.x[16] = 0x2ed5
	s.x[2] = 1
	s.x[16] = s.x[21]*s.x[16] + s.x[2]
	word, err := s.load(s.x[14]+s.x[15], 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap protected-frame first updated word: %w", err)
	}
	s.x[2] = uint64(uint32(word))
	word, err = s.load(s.x[9]+0xb40, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap protected-frame current word: %w", err)
	}
	s.x[3] = uint64(uint32(word))
	s.x[3] = uint64(uint32(s.x[3]) ^ uint32(s.x[10]))
	s.x[2] = uint64(uint32(s.x[2]) ^ uint32(s.x[3]))
	s.x[2] = uint64(uint32(s.x[2]) * uint32(s.x[11]))
	value, err = s.load(s.x[12]+0x3a08, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap protected-frame fourth table: %w", err)
	}
	s.x[3] = value + uint64(byte(s.x[2]))
	s.x[4] = 0xfffffffffdc52ec2
	value, err = s.load(s.x[3]+s.x[4], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap protected-frame fifth byte: %w", err)
	}
	s.x[3] = uint64(byte(value)) ^ 0xb5
	value, err = s.load(s.x[12]+0x828, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap protected-frame fifth table: %w", err)
	}
	s.x[4] = value
	s.x[3] += value
	s.x[4] = 0xffffffffba636d86
	value, err = s.load(s.x[3]+s.x[4], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap protected-frame sixth byte: %w", err)
	}
	s.x[3] = uint64(byte(value))
	s.x[2] = uint64(uint32(s.x[3]) ^ uint32(s.x[2]))
	s.x[3] = 0x13f3
	s.x[2] = uint64(uint32(s.x[3]) + uint32(byte(s.x[2])))

	s.traceBlock(bootstrapProtectedFrameFinish)
	value, err = s.load(s.x[12]+uint64(uint32(s.x[2]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap protected-frame sixth table: %w", err)
	}
	s.x[2] = value
	s.x[3] = 0xffffffffe2130ede
	if err := s.storeValue(s.x[2]+s.x[3], s.x[16], 8); err != nil {
		return 0, fmt.Errorf("bootstrap protected-frame third state output: %w", err)
	}
	s.x[16] >>= 60
	s.x[16] = uint64(uint32(s.x[16]) << 4)
	s.x[16] = uint64(uint32(s.x[16]) + 0x1f)
	s.x[16] &= 0x1f0
	s.x[2] = s.sp
	s.x[20] = s.x[2] - s.x[16]
	s.sp = s.x[20]
	s.x[16] = s.x[21] - 0x40
	if err := s.storeValue(s.x[20], s.x[16], 8); err != nil {
		return 0, fmt.Errorf("bootstrap protected-frame dynamic frame: %w", err)
	}
	word, err = s.load(s.x[14]+s.x[15], 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap protected-frame first output word: %w", err)
	}
	s.x[14] = uint64(uint32(word))
	word, err = s.load(s.x[17]+s.x[1], 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap protected-frame second output word: %w", err)
	}
	s.x[15] = uint64(uint32(word))
	s.x[14] = uint64(uint32(s.x[15]) + uint32(s.x[14]))
	s.x[14] = uint64(uint32(s.x[14]) ^ uint32(s.x[10]))
	s.x[14] = uint64(uint32(s.x[14]) * uint32(s.x[11]))
	value, err = s.load(s.x[12]+0x5e20, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap protected-frame sixth state: %w", err)
	}
	s.x[15] = value + uint64(byte(s.x[14]))
	s.x[16] = 0xffffffffff602305
	value, err = s.load(s.x[15]+s.x[16], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap protected-frame seventh byte: %w", err)
	}
	s.x[15] = uint64(byte(value)) ^ 0xf0
	s.x[16] = 0x9898
	value, err = s.load(s.x[12]+s.x[16], 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap protected-frame seventh state: %w", err)
	}
	s.x[16] = value
	s.x[15] += value
	s.x[16] = 0xfffffffff3e6c916
	value, err = s.load(s.x[15]+s.x[16], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap protected-frame eighth byte: %w", err)
	}
	s.x[15] = uint64(byte(value))
	s.x[14] = uint64(uint32(s.x[15]) ^ uint32(s.x[14]))
	s.x[15] = 0xcfec
	s.x[14] = uint64(uint32(s.x[15]) + uint32(byte(s.x[14])))
	word, err = s.load(s.x[8]+0x2d0, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap protected-frame refreshed first word: %w", err)
	}
	s.x[8] = uint64(uint32(word))
	word, err = s.load(s.x[9]+0xb40, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap protected-frame refreshed second word: %w", err)
	}
	s.x[9] = uint64(uint32(word))
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[9]))
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[10]))
	s.x[8] = uint64(uint32(s.x[8]) * uint32(s.x[11]))
	s.x[9] = s.x[8] & 0xff
	s.x[10] = 0x100f49630
	value, err = s.load(s.x[10]+s.x[9], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap protected-frame ninth byte: %w", err)
	}
	s.x[9] = uint64(byte(value))
	s.x[10] = 0x100fe9490
	s.x[9] = uint64(uint32(s.x[9]) ^ uint32(s.x[13]))
	value, err = s.load(s.x[10]+uint64(uint32(s.x[9])), 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap protected-frame tenth byte: %w", err)
	}
	s.x[9] = uint64(byte(value))
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[9]))
	s.x[8] &= 0xff
	s.x[8] = uint64(uint32(s.x[8]) + 0x1ef)
	value, err = s.load(s.x[12]+uint64(uint32(s.x[8]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap protected-frame final table: %w", err)
	}
	s.x[8] = value + uint64(uint32(s.x[14]))<<3
	s.x[9] = 0xffffffff98245355
	value, err = s.load(s.x[8]+s.x[9], 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap protected-frame call target: %w", err)
	}
	s.x[8] = value
	s.x[1] = s.x[20]
	s.x[17] = 0x9c1c
	s.x[30] = bootstrapProtectedFrameReturn
	return s.x[8], nil
}
