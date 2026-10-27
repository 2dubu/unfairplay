package recovery

import "fmt"

const (
	bootstrapRecordFrameStart  = uint64(0x10037fa0c)
	bootstrapRecordFrameFinish = uint64(0x10037fc00)
	bootstrapRecordFrameReturn = uint64(0x10037fc4c)
)

func runBootstrapRecordFrame(s *state) (uint64, error) {
	s.traceBlock(bootstrapRecordFrameStart)
	s.sp -= 0x60
	for _, pair := range [...]struct {
		first, second int
		offset        uint64
	}{{22, 21, 0x30}, {20, 19, 0x40}, {29, 30, 0x50}} {
		if err := s.storeValue(s.sp+pair.offset, s.x[pair.first], 8); err != nil {
			return 0, fmt.Errorf("bootstrap record-frame register %d: %w", pair.first, err)
		}
		if err := s.storeValue(s.sp+pair.offset+8, s.x[pair.second], 8); err != nil {
			return 0, fmt.Errorf("bootstrap record-frame register %d: %w", pair.second, err)
		}
	}
	s.x[29] = s.sp + 0x50
	s.x[19] = s.x[0]
	s.x[8] = 0x10106d000
	value, err := s.load(s.x[8]+0x408, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap record-frame guard pointer: %w", err)
	}
	s.x[8] = value
	value, err = s.load(value, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap record-frame guard: %w", err)
	}
	s.x[8] = value
	if err := s.storeValue(s.sp+0x28, value, 8); err != nil {
		return 0, fmt.Errorf("bootstrap record-frame saved guard: %w", err)
	}
	s.x[8] = 0x101074000
	word, err := s.load(s.x[8]+0x1b0, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap record-frame first state word: %w", err)
	}
	s.x[10] = uint64(uint32(word))
	s.x[9] = bootstrapStateBase
	word, err = s.load(s.x[9]+0xb40, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap record-frame second state word: %w", err)
	}
	s.x[14] = uint64(uint32(word))
	s.x[11] = 0xa6
	s.x[10] = uint64(uint32(s.x[14]) - uint32(s.x[10]))
	s.x[12] = uint64(uint32(s.x[10]) ^ uint32(s.x[11]))
	s.x[10] = 0x0a5cd4cf
	s.x[12] = uint64(uint32(s.x[12]) * uint32(s.x[10]))
	s.x[13] = 0x100fe9390 + uint64(byte(s.x[12]))
	value, err = s.load(s.x[13], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap record-frame first byte: %w", err)
	}
	s.x[13] = uint64(byte(value)) ^ 0x84
	s.x[15] = 0x100de2af0
	value, err = s.load(s.x[15]+uint64(uint32(s.x[13])), 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap record-frame second byte: %w", err)
	}
	s.x[13] = uint64(byte(value))
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[13]))
	s.x[12] &= 0xff
	s.x[13] = uint64(uint32(s.x[12]) + 0xc41)
	s.x[12] = bootstrapTargetTable
	value, err = s.load(s.x[12]+uint64(uint32(s.x[13]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap record-frame first target: %w", err)
	}
	s.x[13] = value
	s.x[15] = 0xffffffffdbd5e9bd
	s.x[13] += s.x[15]
	value, err = s.load(s.x[13], 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap record-frame first protected state: %w", err)
	}
	s.x[15] = value
	s.x[14] = uint64(uint32(s.x[14]) ^ uint32(s.x[11]))
	s.x[14] = uint64(uint32(s.x[15]) ^ uint32(s.x[14]))
	s.x[14] = uint64(uint32(s.x[14]) * uint32(s.x[10]))
	value, err = s.load(s.x[12]+0x55c0, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap record-frame second table: %w", err)
	}
	s.x[16] = value
	s.x[17] = 0xffffffffc10da2f9
	s.x[16] += s.x[17]
	s.x[17] = s.x[16] + uint64(byte(s.x[14]))
	value, err = s.load(s.x[17], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap record-frame third byte: %w", err)
	}
	s.x[17] = uint64(byte(value)) ^ 4
	value, err = s.load(s.x[12]+0x2df0, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap record-frame third table: %w", err)
	}
	s.x[0] = value
	s.x[1] = 0xffffffff9f79fd81
	s.x[0] += s.x[1]
	value, err = s.load(s.x[0]+uint64(uint32(s.x[17])), 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap record-frame fourth byte: %w", err)
	}
	s.x[17] = uint64(byte(value))
	s.x[14] = uint64(uint32(s.x[17]) ^ uint32(s.x[14]))
	s.x[17] = 0x1331
	s.x[14] = uint64(uint32(s.x[17]) + uint32(byte(s.x[14])))
	value, err = s.load(s.x[12]+uint64(uint32(s.x[14]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap record-frame second target: %w", err)
	}
	s.x[14] = value
	s.x[17] = 0xffffffff9608a825
	s.x[14] += s.x[17]
	value, err = s.load(s.x[14], 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap record-frame second protected state: %w", err)
	}
	s.x[17] = value
	s.x[15] ^= value
	s.x[17] = s.sp + 8
	s.x[1] = 0x673e3a80075f02c6
	s.x[2] = 0x73d046208eaed1a6
	s.x[15] += s.x[17]
	s.x[17] = s.x[15] ^ s.x[2]
	s.x[15] = s.x[15]*s.x[10] + s.x[1]
	s.x[17] *= s.x[10]
	if err := s.storeValue(s.x[13], s.x[15], 8); err != nil {
		return 0, fmt.Errorf("bootstrap record-frame first protected output: %w", err)
	}
	if err := s.storeValue(s.x[14], s.x[17], 8); err != nil {
		return 0, fmt.Errorf("bootstrap record-frame second protected output: %w", err)
	}
	s.x[15] = uint64(uint32(s.x[19]) ^ 0xdc40790b)
	s.x[17] = 0x5973a279
	s.x[15] = uint64(uint32(s.x[15]) * uint32(s.x[17]))
	word, err = s.load(s.x[19], 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap record-frame descriptor word: %w", err)
	}
	s.x[17] = uint64(uint32(word))
	s.x[20] = uint64(uint32(s.x[17]) ^ uint32(s.x[15]))
	value, err = s.load(s.x[19]+8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap record-frame descriptor pointer: %w", err)
	}
	s.x[21] = value
	s.x[15] = s.sp + 0x10
	s.x[17] = 0x0d5e801f
	s.x[17] = uint64(uint32(s.x[15]) | uint32(s.x[17]))
	s.x[1] = 0xf2a17fe0
	s.x[15] = uint64(uint32(s.x[15]) | uint32(s.x[1]))
	s.x[15] = uint64(uint32(s.x[17]) - uint32(s.x[15]))
	s.x[15] = uint64(uint32(s.x[15]) + uint32(s.x[1]))
	s.x[17] = 0xea6c527d
	s.x[15] = uint64(uint32(s.x[15]) ^ uint32(s.x[17]))
	s.x[17] = 0xba209ef7
	if err := s.storeValue(s.sp+0x18, s.x[21], 8); err != nil {
		return 0, fmt.Errorf("bootstrap record-frame saved descriptor pointer: %w", err)
	}
	s.x[15] = uint64(uint32(s.x[15])*uint32(s.x[17]) + uint32(s.x[20]))
	s.x[17] = 0x2f6a
	s.x[15] = uint64(uint32(s.x[15]) + uint32(s.x[17]))
	if err := s.storeValue(s.sp+0x10, uint64(uint32(s.x[15])), 4); err != nil {
		return 0, fmt.Errorf("bootstrap record-frame saved selector: %w", err)
	}
	word, err = s.load(s.x[13], 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap record-frame first updated word: %w", err)
	}
	s.x[13] = uint64(uint32(word))
	word, err = s.load(s.x[14], 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap record-frame second updated word: %w", err)
	}
	s.x[14] = uint64(uint32(word))
	s.x[13] = uint64(uint32(s.x[14]) - uint32(s.x[13]))
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[11]))
	s.x[13] = uint64(uint32(s.x[13]) * uint32(s.x[10]))
	s.x[14] = s.x[16] + uint64(byte(s.x[13]))
	value, err = s.load(s.x[14], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap record-frame fifth byte: %w", err)
	}
	s.x[14] = uint64(byte(value)) ^ 4
	value, err = s.load(s.x[0]+uint64(uint32(s.x[14])), 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap record-frame sixth byte: %w", err)
	}
	s.x[14] = uint64(byte(value))
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[14]))
	s.x[14] = 0xd07e
	s.x[13] = uint64(uint32(s.x[14]) + uint32(byte(s.x[13])))
	word, err = s.load(s.x[8]+0x1b0, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap record-frame refreshed first word: %w", err)
	}
	s.x[8] = uint64(uint32(word))
	word, err = s.load(s.x[9]+0xb40, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap record-frame refreshed second word: %w", err)
	}
	s.x[9] = uint64(uint32(word))
	s.x[8] = uint64(uint32(s.x[9]) + uint32(s.x[8]))
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[11]))
	s.x[9] = 0x100f49630

	s.traceBlock(bootstrapRecordFrameFinish)
	s.x[8] = uint64(uint32(s.x[8]) * uint32(s.x[10]))
	s.x[9] += uint64(byte(s.x[8]))
	value, err = s.load(s.x[9], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap record-frame seventh byte: %w", err)
	}
	s.x[9] = uint64(byte(value)) ^ 0x17
	s.x[11] = 0x100fe9490
	value, err = s.load(s.x[11]+uint64(uint32(s.x[9])), 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap record-frame eighth byte: %w", err)
	}
	s.x[9] = uint64(byte(value))
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[9]))
	s.x[8] &= 0xff
	s.x[8] = uint64(uint32(s.x[8]) + 0x202)
	value, err = s.load(s.x[12]+uint64(uint32(s.x[8]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap record-frame final table: %w", err)
	}
	s.x[8] = value
	s.x[9] = 0xffffffff98245355
	s.x[22] = s.x[8] + s.x[9]
	next, err := s.load(s.x[22]+uint64(uint32(s.x[13]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap record-frame call target: %w", err)
	}
	s.x[8] = next
	s.x[0] = s.sp + 0x10
	s.x[17] = 0xda16
	s.x[30] = bootstrapRecordFrameReturn
	return next, nil
}
