package recovery

import "testing"

func TestRecoveryStageTraceSeparatesLiveInputsAndWrites(t *testing.T) {
	const base = uint64(0x7000000f0000)
	value := &model{pages: map[uint64][]byte{base: make([]byte, pageSize)}, traceStages: true}
	value.pages[base][0x10] = 0x42
	worker := newState(value)
	worker.reset()
	worker.x[1] = 7
	worker.sp = base + 0x800
	err := withRecoveryStage(worker, "test-stage", func() error {
		worker.traceBlock(0x1000)
		worker.traceBlock(0x1000)
		worker.traceBlock(0x2000)
		if _, err := worker.load(base+0x10, 1); err != nil {
			return err
		}
		if err := worker.storeValue(base+0x20, 0x99, 1); err != nil {
			return err
		}
		if _, err := worker.load(base+0x20, 1); err != nil {
			return err
		}
		worker.x[1] = 9
		worker.sp += 0x20
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(worker.tracer.results) != 1 {
		t.Fatalf("trace results = %d, want 1", len(worker.tracer.results))
	}
	trace := worker.tracer.results[0]
	if trace.Name != "test-stage" || trace.Blocks != 3 || trace.UniqueBlocks != 2 || len(trace.TopBlocks) != 2 ||
		trace.TopBlocks[0].Address != 0x1000 || trace.TopBlocks[0].Count != 2 {
		t.Fatalf("unexpected block trace: %+v", trace)
	}
	if len(trace.LiveInSpans) != 1 || trace.LiveInSpans[0].Address != base+0x10 || trace.LiveInSpans[0].Bytes != 1 {
		t.Fatalf("unexpected live inputs: %+v", trace.LiveInSpans)
	}
	if len(trace.WriteSpans) != 1 || trace.WriteSpans[0].Address != base+0x20 || trace.WriteSpans[0].Bytes != 1 ||
		!trace.WriteSpans[0].Changed {
		t.Fatalf("unexpected writes: %+v", trace.WriteSpans)
	}
	if len(trace.RegisterChanges) != 1 || trace.RegisterChanges[0].Register != "x1" ||
		trace.RegisterChanges[0].Before != 7 || trace.RegisterChanges[0].After != 9 ||
		trace.ExitStackPointer != trace.EntryStackPointer+0x20 {
		t.Fatalf("unexpected register trace: %+v", trace)
	}
}

func TestRecoveryStageTraceNestsHookContracts(t *testing.T) {
	const base = uint64(0x710000000000)
	value := &model{pages: map[uint64][]byte{base: make([]byte, pageSize)}, traceStages: true}
	worker := newState(value)
	worker.reset()
	if err := withRecoveryStage(worker, "outer", func() error {
		if _, err := worker.load(base+1, 1); err != nil {
			return err
		}
		return withRecoveryStage(worker, "outer.writer", func() error {
			worker.traceBlock(0x3000)
			return worker.storeValue(base+2, 0xaa, 1)
		})
	}); err != nil {
		t.Fatal(err)
	}
	if len(worker.tracer.results) != 2 || worker.tracer.results[0].Name != "outer.writer" ||
		worker.tracer.results[1].Name != "outer" {
		t.Fatalf("nested stage order: %+v", worker.tracer.results)
	}
	if len(worker.tracer.results[1].WriteSpans) != 1 || worker.tracer.results[1].Blocks != 1 {
		t.Fatalf("outer trace did not include nested work: %+v", worker.tracer.results[1])
	}
}
