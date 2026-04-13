package isa

import "testing"

func TestSignExtend12(t *testing.T) {
	tests := []struct {
		input uint16
		want  int32
	}{
		{0x001, 1},
		{0x7FF, 2047},
		{0x800, -2048},
		{0xFFF, -1},
	}
	for _, tt := range tests {
		got := SignExtend12(tt.input)
		if got != tt.want {
			t.Errorf("SignExtend12(0x%03X) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestSignExtend13(t *testing.T) {
	tests := []struct {
		input uint16
		want  int32
	}{
		{0x0004, 4},
		{0x1000, -4096},
		{0x1FFC, -4},
	}
	for _, tt := range tests {
		got := SignExtend13(tt.input)
		if got != tt.want {
			t.Errorf("SignExtend13(0x%04X) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestDecodeJImm(t *testing.T) {
	// JAL with offset +40 = 0x28
	// imm[20]=0, imm[19:12]=00000000, imm[11]=0, imm[10:1]=0000010100
	// Raw field from inst[31:12]: [imm20|imm10:1|imm11|imm19:12]
	// = 0 | 0000010100 | 0 | 00000000 = 0x02800
	raw := uint32(0x02800)
	got := DecodeJImm(raw)
	if got != 40 {
		t.Errorf("DecodeJImm(0x%05X) = %d, want 40", raw, got)
	}

	// Test negative offset: -4 = 0x1FFFFC in 21 bits
	// imm[20]=1, imm[19:12]=0xFF, imm[11]=1, imm[10:1]=0x3FE
	// Raw [imm20|imm10:1|imm11|imm19:12]:
	//   bit19=imm20=1, bits18:9=imm10:1=1111111110, bit8=imm11=1, bits7:0=imm19:12=0xFF
	//   = 0xFFDFF
	rawNeg := uint32(0xFFDFF)
	gotNeg := DecodeJImm(rawNeg)
	if gotNeg != -4 {
		t.Errorf("DecodeJImm(0x%05X) = %d, want -4", rawNeg, gotNeg)
	}
}
