package memory

import (
	"testing"
)

func TestStoreLoadWord(t *testing.T) {
	m := New()
	if err := m.StoreWord(0x10000000, 0x12345678); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := m.LoadWord(0x10000000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 0x12345678 {
		t.Errorf("LoadWord = 0x%08X, want 0x12345678", uint32(got))
	}
}

func TestStoreLoadNegativeWord(t *testing.T) {
	m := New()
	if err := m.StoreWord(0x100, int32(-1)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := m.LoadWord(0x100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != -1 {
		t.Errorf("LoadWord = %d, want -1", got)
	}
}

func TestStoreLoadHalf(t *testing.T) {
	m := New()
	if err := m.StoreHalf(0x100, 0x1234); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := m.LoadHalf(0x100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 0x1234 {
		t.Errorf("LoadHalf = 0x%04X, want 0x1234", uint16(got))
	}
}

func TestStoreLoadByte(t *testing.T) {
	m := New()
	if err := m.StoreByte(0x100, -5); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := m.LoadByte(0x100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != -5 {
		t.Errorf("LoadByte = %d, want -5", got)
	}
}

func TestLoadByteU(t *testing.T) {
	m := New()
	if err := m.StoreByte(0x100, -1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := m.LoadByteU(0x100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 0xFF {
		t.Errorf("LoadByteU = 0x%02X, want 0xFF", got)
	}
}

func TestLoadHalfU(t *testing.T) {
	m := New()
	if err := m.StoreHalf(0x100, -1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := m.LoadHalfU(0x100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 0xFFFF {
		t.Errorf("LoadHalfU = 0x%04X, want 0xFFFF", got)
	}
}

func TestUnwrittenMemoryReturnsZero(t *testing.T) {
	m := New()
	got, err := m.LoadWord(0x1000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 0 {
		t.Errorf("LoadWord unwritten = %d, want 0", got)
	}
}

func TestMisalignedWordAccess(t *testing.T) {
	m := New()
	_, err := m.LoadWord(0x101)
	if err == nil {
		t.Error("expected misaligned error for word access at 0x101")
	}
}

func TestMisalignedHalfAccess(t *testing.T) {
	m := New()
	_, err := m.LoadHalf(0x101)
	if err == nil {
		t.Error("expected misaligned error for half access at 0x101")
	}
}

func TestLittleEndian(t *testing.T) {
	m := New()
	if err := m.StoreWord(0x100, 0x04030201); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b0, _ := m.LoadByteU(0x100)
	b1, _ := m.LoadByteU(0x101)
	b2, _ := m.LoadByteU(0x102)
	b3, _ := m.LoadByteU(0x103)
	if b0 != 0x01 || b1 != 0x02 || b2 != 0x03 || b3 != 0x04 {
		t.Errorf("little-endian failed: got [%02X %02X %02X %02X], want [01 02 03 04]", b0, b1, b2, b3)
	}
}

func TestMisalignedStoreWord(t *testing.T) {
	m := New()
	if err := m.StoreWord(0x101, 0x12345678); err == nil {
		t.Error("expected misaligned error for word store at 0x101")
	}
}

func TestMisalignedStoreHalf(t *testing.T) {
	m := New()
	if err := m.StoreHalf(0x101, 0x1234); err == nil {
		t.Error("expected misaligned error for half store at 0x101")
	}
}
