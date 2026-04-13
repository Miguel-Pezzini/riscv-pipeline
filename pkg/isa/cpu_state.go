package isa

// CPUState é a interface que o executor fornece às instruções.
// Evita import cycle entre isa e cpu.
type CPUState interface {
	ReadReg(n int) int32
	WriteReg(n int, v int32)
	GetPC() uint32
	SetPC(pc uint32)
	LoadWord(addr uint32) (int32, error)
	LoadHalf(addr uint32) (int16, error)
	LoadByte(addr uint32) (int8, error)
	LoadHalfU(addr uint32) (uint16, error)
	LoadByteU(addr uint32) (uint8, error)
	StoreWord(addr uint32, v int32) error
	StoreHalf(addr uint32, v int16) error
	StoreByte(addr uint32, v int8) error
}
