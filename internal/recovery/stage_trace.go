package recovery

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hash"
	"sort"
	"strings"
)

const maximumNestedStageSamples = 3

func nestedStageSampleLimit(name string) int {
	switch name {
	case "outer.optional":
		return 256
	case "outer.even.iteration", "outer.odd.iteration":
		return 50
	case "outer.tail.iteration", "outer.writer.iteration":
		return 140
	}
	return maximumNestedStageSamples
}

type RecoveryTraceSpan struct {
	Address      uint64 `json:"address"`
	Bytes        int    `json:"bytes"`
	Region       string `json:"region"`
	BeforeSHA256 string `json:"before_sha256"`
	AfterSHA256  string `json:"after_sha256"`
	Changed      bool   `json:"changed"`
	BeforeHex    string `json:"before_hex,omitempty"`
	AfterHex     string `json:"after_hex,omitempty"`
}

type RecoveryTraceRegister struct {
	Register string `json:"register"`
	Before   uint64 `json:"before"`
	After    uint64 `json:"after"`
}

type RecoveryTraceBlock struct {
	Address uint64 `json:"address"`
	Count   int    `json:"count"`
}

type RecoveryStageTrace struct {
	Name              string                  `json:"name"`
	Invocation        int                     `json:"invocation"`
	TotalInvocations  int                     `json:"total_invocations"`
	EntryStateSHA256  string                  `json:"entry_state_sha256"`
	ExitStateSHA256   string                  `json:"exit_state_sha256"`
	PathSHA256        string                  `json:"path_sha256"`
	Blocks            int                     `json:"blocks"`
	UniqueBlocks      int                     `json:"unique_blocks"`
	TopBlocks         []RecoveryTraceBlock    `json:"top_blocks"`
	LiveInSpans       []RecoveryTraceSpan     `json:"live_in_spans"`
	WriteSpans        []RecoveryTraceSpan     `json:"write_spans"`
	RegisterChanges   []RecoveryTraceRegister `json:"register_changes"`
	EntryStackPointer uint64                  `json:"entry_stack_pointer"`
	ExitStackPointer  uint64                  `json:"exit_stack_pointer"`
	Error             string                  `json:"error,omitempty"`
	EntryRegisters    []uint64                `json:"entry_registers,omitempty"`
	ExitRegisters     []uint64                `json:"exit_registers,omitempty"`
}

type recoveryStageCapture struct {
	name       string
	entryState string
	entryX     [31]uint64
	entrySP    uint64
	before     map[uint64][]byte
	read       map[uint64][]byte
	liveIn     map[uint64][]byte
	write      map[uint64][]byte
	written    map[uint64][]byte
	blocks     map[uint64]int
	path       hash.Hash
	blockCount int
	invocation int
}

type recoveryTraceRecorder struct {
	active        []*recoveryStageCapture
	results       []RecoveryStageTrace
	invocations   map[string]int
	includeValues bool
}

func tracePageBits(pages map[uint64][]byte, base uint64) []byte {
	bits := pages[base]
	if bits == nil {
		bits = make([]byte, pageSize)
		pages[base] = bits
	}
	return bits
}

func traceMark(pages map[uint64][]byte, address uint64, size int) {
	for size > 0 {
		base := address &^ (pageSize - 1)
		offset := int(address - base)
		count := min(size, int(pageSize)-offset)
		bits := tracePageBits(pages, base)
		for index := offset; index < offset+count; index++ {
			bits[index] = 1
		}
		address += uint64(count)
		size -= count
	}
}

func traceMarkRead(capture *recoveryStageCapture, address uint64, size int) {
	traceMark(capture.read, address, size)
	for size > 0 {
		base := address &^ (pageSize - 1)
		offset := int(address - base)
		count := min(size, int(pageSize)-offset)
		written := capture.written[base]
		live := tracePageBits(capture.liveIn, base)
		for index := offset; index < offset+count; index++ {
			if written == nil || written[index] == 0 {
				live[index] = 1
			}
		}
		address += uint64(count)
		size -= count
	}
}

func traceMarkWrite(capture *recoveryStageCapture, address uint64, size int) {
	traceMark(capture.write, address, size)
	traceMark(capture.written, address, size)
}

func currentTracePages(s *state) map[uint64][]byte {
	result := make(map[uint64][]byte, len(s.model.pages)+len(s.mapped))
	for address, page := range s.model.pages {
		if writable := s.writable[address]; writable != nil {
			page = writable
		}
		result[address] = append([]byte(nil), page...)
	}
	for address, page := range s.mapped {
		if s.model.pages[address] != nil {
			continue
		}
		if writable := s.writable[address]; writable != nil {
			page = writable
		}
		result[address] = append([]byte(nil), page...)
	}
	return result
}

func recoveryTraceRegion(address uint64) string {
	switch {
	case address >= 0xf000000000000000:
		return "external"
	case address >= 0x710000000000 && address < 0x720000000000:
		return "heap"
	case address >= 0x700000000000 && address < 0x710000000000:
		return "stack"
	case address >= 0x100000000 && address < 0x200000000:
		return "image"
	default:
		return "other"
	}
}

func traceSpans(pages map[uint64][]byte, liveIn bool) [][2]uint64 {
	var addresses []uint64
	for base, bits := range pages {
		for offset, marked := range bits {
			if marked == 0 || liveIn && marked != 1 {
				continue
			}
			addresses = append(addresses, base+uint64(offset))
		}
	}
	sort.Slice(addresses, func(left, right int) bool { return addresses[left] < addresses[right] })
	var spans [][2]uint64
	for _, address := range addresses {
		if len(spans) != 0 && spans[len(spans)-1][1] == address &&
			recoveryTraceRegion(spans[len(spans)-1][0]) == recoveryTraceRegion(address) {
			spans[len(spans)-1][1]++
			continue
		}
		spans = append(spans, [2]uint64{address, address + 1})
	}
	return spans
}

func traceSpanHash(pages map[uint64][]byte, start, end uint64) string {
	hasher := sha256.New()
	var zero [pageSize]byte
	for address := start; address < end; {
		base := address &^ (pageSize - 1)
		offset := int(address - base)
		count := min(int(end-address), int(pageSize)-offset)
		page := pages[base]
		if page == nil {
			page = zero[:]
		}
		_, _ = hasher.Write(page[offset : offset+count])
		address += uint64(count)
	}
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

func traceSpanBytes(pages map[uint64][]byte, start, end uint64) []byte {
	result := make([]byte, 0, end-start)
	var zero [pageSize]byte
	for address := start; address < end; {
		base := address &^ (pageSize - 1)
		offset := int(address - base)
		count := min(int(end-address), int(pageSize)-offset)
		page := pages[base]
		if page == nil {
			page = zero[:]
		}
		result = append(result, page[offset:offset+count]...)
		address += uint64(count)
	}
	return result
}

func buildTraceSpans(bits map[uint64][]byte, liveIn bool, before, after map[uint64][]byte, includeValues bool) []RecoveryTraceSpan {
	raw := traceSpans(bits, liveIn)
	result := make([]RecoveryTraceSpan, 0, len(raw))
	for _, span := range raw {
		beforeHash := traceSpanHash(before, span[0], span[1])
		afterHash := traceSpanHash(after, span[0], span[1])
		item := RecoveryTraceSpan{
			Address: span[0], Bytes: int(span[1] - span[0]), Region: recoveryTraceRegion(span[0]),
			BeforeSHA256: beforeHash, AfterSHA256: afterHash, Changed: beforeHash != afterHash,
		}
		if includeValues {
			beforeBytes := traceSpanBytes(before, span[0], span[1])
			afterBytes := traceSpanBytes(after, span[0], span[1])
			item.BeforeHex = hex.EncodeToString(beforeBytes)
			item.AfterHex = hex.EncodeToString(afterBytes)
			clear(beforeBytes)
			clear(afterBytes)
		}
		result = append(result, item)
	}
	return result
}

func (r *recoveryTraceRecorder) begin(s *state, name string) *recoveryStageCapture {
	if r.invocations == nil {
		r.invocations = make(map[string]int)
	}
	r.invocations[name]++
	invocation := r.invocations[name]
	if strings.Contains(name, ".") && invocation > nestedStageSampleLimit(name) {
		return nil
	}
	capture := &recoveryStageCapture{
		name: name, invocation: invocation, entryState: s.digest(), entryX: s.x, entrySP: s.sp,
		before: currentTracePages(s), read: make(map[uint64][]byte), liveIn: make(map[uint64][]byte), write: make(map[uint64][]byte),
		written: make(map[uint64][]byte), blocks: make(map[uint64]int), path: sha256.New(),
	}
	r.active = append(r.active, capture)
	return capture
}

func (r *recoveryTraceRecorder) finish(s *state, capture *recoveryStageCapture, stageErr error) {
	if len(r.active) == 0 || r.active[len(r.active)-1] != capture {
		panic("recovery trace stage stack is corrupt")
	}
	r.active = r.active[:len(r.active)-1]
	after := currentTracePages(s)
	trace := RecoveryStageTrace{
		Name: capture.name, Invocation: capture.invocation,
		EntryStateSHA256: capture.entryState, ExitStateSHA256: s.digest(),
		PathSHA256: fmt.Sprintf("%x", capture.path.Sum(nil)), Blocks: capture.blockCount,
		UniqueBlocks: len(capture.blocks), LiveInSpans: buildTraceSpans(capture.liveIn, false, capture.before, after, r.includeValues),
		WriteSpans:        buildTraceSpans(capture.write, false, capture.before, after, r.includeValues),
		EntryStackPointer: capture.entrySP, ExitStackPointer: s.sp,
	}
	if stageErr != nil {
		trace.Error = stageErr.Error()
	}
	if r.includeValues {
		trace.EntryRegisters = append(trace.EntryRegisters, capture.entryX[:]...)
		trace.ExitRegisters = append(trace.ExitRegisters, s.x[:]...)
	}
	for index, before := range capture.entryX {
		if afterValue := s.x[index]; afterValue != before {
			trace.RegisterChanges = append(trace.RegisterChanges, RecoveryTraceRegister{
				Register: fmt.Sprintf("x%d", index), Before: before, After: afterValue,
			})
		}
	}
	blocks := make([]RecoveryTraceBlock, 0, len(capture.blocks))
	for address, count := range capture.blocks {
		blocks = append(blocks, RecoveryTraceBlock{Address: address, Count: count})
	}
	sort.Slice(blocks, func(left, right int) bool {
		if blocks[left].Count != blocks[right].Count {
			return blocks[left].Count > blocks[right].Count
		}
		return blocks[left].Address < blocks[right].Address
	})
	trace.TopBlocks = blocks[:min(len(blocks), 32)]
	r.results = append(r.results, trace)
	for _, page := range capture.before {
		clear(page)
	}
}

func (r *recoveryTraceRecorder) completedResults() []RecoveryStageTrace {
	result := append([]RecoveryStageTrace(nil), r.results...)
	for index := range result {
		result[index].TotalInvocations = r.invocations[result[index].Name]
	}
	return result
}

func (s *state) traceAccess(address uint64, size int, write bool) {
	if s.tracer == nil || size <= 0 {
		return
	}
	for _, capture := range s.tracer.active {
		if write {
			traceMarkWrite(capture, address, size)
		} else {
			traceMarkRead(capture, address, size)
		}
	}
}

func (s *state) traceBlock(address uint64) {
	if s.tracer == nil {
		return
	}
	var encoded [8]byte
	binary.LittleEndian.PutUint64(encoded[:], address)
	for _, capture := range s.tracer.active {
		capture.blockCount++
		capture.blocks[address]++
		_, _ = capture.path.Write(encoded[:])
	}
}

func withRecoveryStage(s *state, name string, run func() error) (err error) {
	if s.tracer == nil {
		return run()
	}
	capture := s.tracer.begin(s, name)
	if capture == nil {
		return run()
	}
	defer func() { s.tracer.finish(s, capture, err) }()
	return run()
}

func withRecoveryValueTrace(s *state, name string, run func() error) error {
	if s.tracer == nil || !s.tracer.includeValues {
		return run()
	}
	return withRecoveryStage(s, name, run)
}
