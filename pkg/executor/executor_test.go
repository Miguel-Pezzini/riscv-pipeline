package executor

import (
	"strings"
	"testing"

	"riscv-instruction-encoder/pkg/cpu"
)

func TestECALLExitServiceHalts(t *testing.T) {
	state := cpu.NewState()
	if err := state.StoreWord(cpu.TextBase, 0x00000073); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	state.WriteReg(17, 10)

	exec := New(state, Config{MaxSteps: 1})
	if err := exec.Run(); err != nil {
		t.Fatalf("Run() unexpected error: %v", err)
	}
	if !exec.Halted {
		t.Fatal("expected executor to halt on ECALL exit")
	}
}

func TestECALLUnsupportedServiceReturnsError(t *testing.T) {
	state := cpu.NewState()
	if err := state.StoreWord(cpu.TextBase, 0x00000073); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	state.WriteReg(17, 1)

	exec := New(state, Config{MaxSteps: 1})
	err := exec.Run()
	if err == nil {
		t.Fatal("expected error for unsupported ECALL service")
	}
	if !strings.Contains(err.Error(), "unsupported ECALL service 1") {
		t.Fatalf("unexpected error: %v", err)
	}
}
