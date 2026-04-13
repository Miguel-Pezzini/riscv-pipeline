package isa

// SignExtend12 sign-extends a 12-bit immediate to int32.
func SignExtend12(imm uint16) int32 {
	if imm&0x800 != 0 {
		return int32(imm) | ^int32(0xFFF)
	}
	return int32(imm)
}

// SignExtend13 sign-extends a 13-bit immediate (B-type offset) to int32.
func SignExtend13(imm uint16) int32 {
	if imm&0x1000 != 0 {
		return int32(imm) | ^int32(0x1FFF)
	}
	return int32(imm)
}

// SignExtend21 sign-extends a 21-bit immediate (J-type offset) to int32.
func SignExtend21(imm uint32) int32 {
	if imm&0x100000 != 0 {
		return int32(imm) | ^int32(0x1FFFFF)
	}
	return int32(imm)
}

// DecodeJImm reconstructs the 21-bit signed offset from the raw 20-bit J-type immediate field.
// The raw field is inst[31:12] stored as-is. The actual offset layout is:
// imm[20|10:1|11|19:12]
func DecodeJImm(raw uint32) int32 {
	bit20 := (raw >> 19) & 0x1
	bits10_1 := (raw >> 9) & 0x3FF
	bit11 := (raw >> 8) & 0x1
	bits19_12 := raw & 0xFF
	imm := (bit20 << 20) | (bits19_12 << 12) | (bit11 << 11) | (bits10_1 << 1)
	return SignExtend21(imm)
}
