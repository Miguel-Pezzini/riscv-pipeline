package cpu

import (
	"testing"
)

func TestNewState(t *testing.T) {
	s := NewState()
	if s.PC != TextBase {
		t.Errorf("PC = 0x%08X, want 0x%08X", s.PC, TextBase)
	}
	if s.Regs[2] != int32(StackTop) {
		t.Errorf("sp = 0x%08X, want 0x%08X", uint32(s.Regs[2]), StackTop)
	}
}

func TestWriteRegZeroIsNoop(t *testing.T) {
	s := NewState()
	s.WriteReg(0, 42)
	if s.ReadReg(0) != 0 {
		t.Errorf("x0 = %d, want 0", s.ReadReg(0))
	}
}

func TestReadWriteReg(t *testing.T) {
	s := NewState()
	s.WriteReg(5, 7)
	if s.ReadReg(5) != 7 {
		t.Errorf("x5 = %d, want 7", s.ReadReg(5))
	}
}

func TestGetSetPC(t *testing.T) {
	s := NewState()
	s.SetPC(0x1000)
	if s.GetPC() != 0x1000 {
		t.Errorf("PC = 0x%08X, want 0x1000", s.GetPC())
	}
}

func TestMemoryAccess(t *testing.T) {
	s := NewState()
	if err := s.StoreWord(DataBase, int32(-559038737)); err != nil { // 0xDEADBEEF
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := s.LoadWord(DataBase)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != int32(-559038737) {
		t.Errorf("LoadWord = 0x%08X, want 0xDEADBEEF", uint32(got))
	}
}
