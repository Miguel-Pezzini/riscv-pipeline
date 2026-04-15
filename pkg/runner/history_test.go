package runner

import (
	"testing"

	"riscv-instruction-encoder/pkg/cpu"
	"riscv-instruction-encoder/pkg/executor"
)

func TestPipelineFromExecutionHistoryUsesDynamicPath(t *testing.T) {
	state := cpu.NewState()
	program := []uint32{
		0x00100093, // addi x1, x0, 1
		0x00100113, // addi x2, x0, 1
		0x00208463, // beq  x1, x2, +8
		0x06300193, // addi x3, x0, 99 (skipped)
		0x00700213, // addi x4, x0, 7
		0x00100073, // ebreak
	}

	for i, word := range program {
		addr := cpu.TextBase + uint32(i*4)
		if err := state.StoreWord(addr, int32(word)); err != nil {
			t.Fatalf("StoreWord(0x%08X): %v", addr, err)
		}
	}

	exec := executor.New(state, executor.Config{MaxSteps: 16})
	if err := exec.Run(); err != nil {
		t.Fatalf("Run() unexpected error: %v", err)
	}

	pipeline := NewPipelineFromExecutionHistory(exec.History, true, true, true, "")
	wantPCs := []int{
		int(cpu.TextBase),
		int(cpu.TextBase + 4),
		int(cpu.TextBase + 8),
		int(cpu.TextBase + 16),
	}

	if len(pipeline.Instructions) != len(wantPCs) {
		t.Fatalf("len(Instructions) = %d, want %d", len(pipeline.Instructions), len(wantPCs))
	}

	for i, want := range wantPCs {
		if got := pipeline.Instructions[i].OriginalPC; got != want {
			t.Fatalf("Instructions[%d].OriginalPC = 0x%08X, want 0x%08X", i, got, want)
		}
	}

	if got := state.ReadReg(3); got != 0 {
		t.Fatalf("x3 = %d, want 0 because skipped instruction must not execute", got)
	}
	if got := state.ReadReg(4); got != 7 {
		t.Fatalf("x4 = %d, want 7", got)
	}

	pipeline.Run()

	for _, instruction := range pipeline.Instructions {
		if instruction.OriginalPC == int(cpu.TextBase+12) {
			t.Fatalf("found skipped PC 0x%08X in pipeline trace", instruction.OriginalPC)
		}
	}

	nops := 0
	for _, instruction := range pipeline.Instructions {
		if instruction.Instruction.GetMeta().Name == "NOP" {
			nops++
		}
	}
	if nops == 0 {
		t.Fatal("expected control hazard to insert at least one NOP")
	}
}
