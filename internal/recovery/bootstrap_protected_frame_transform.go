package recovery

import "fmt"

const (
	bootstrapProtectedFrameStart  = uint64(0x1007b2e8c)
	bootstrapProtectedFrameFinish = uint64(0x1007b3000)
	bootstrapProtectedFrameReturn = uint64(0x1007b30e4)
	bootstrapStateBase            = uint64(0x101075000)
	bootstrapTargetTable          = uint64(0x10105b6f0)
)

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
