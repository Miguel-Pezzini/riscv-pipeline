package cpu

import (
	"fmt"
	"riscv-instruction-encoder/pkg/memory"
	"strings"
)

const (
	TextBase = uint32(0x00400000)
	DataBase = uint32(0x10000000)
	StackTop = uint32(0x7FFFFFFC)
)

type State struct {
	Regs [32]int32
	PC   uint32
	Mem  *memory.Memory
}

func NewState() *State {
	s := &State{
		PC:  TextBase,
		Mem: memory.New(),
	}
	s.Regs[2] = int32(StackTop) // sp
	return s
}

func (s *State) ReadReg(n int) int32 {
	if n == 0 {
		return 0
	}
	return s.Regs[n]
}

func (s *State) WriteReg(n int, v int32) {
	if n == 0 {
		return
	}
	s.Regs[n] = v
}

func (s *State) GetPC() uint32 {
	return s.PC
}

func (s *State) SetPC(pc uint32) {
	s.PC = pc
}

func (s *State) LoadWord(addr uint32) (int32, error) {
	return s.Mem.LoadWord(addr)
}

func (s *State) LoadHalf(addr uint32) (int16, error) {
	return s.Mem.LoadHalf(addr)
}

func (s *State) LoadByte(addr uint32) (int8, error) {
	return s.Mem.LoadByte(addr)
}

func (s *State) LoadHalfU(addr uint32) (uint16, error) {
	return s.Mem.LoadHalfU(addr)
}

func (s *State) LoadByteU(addr uint32) (uint8, error) {
	return s.Mem.LoadByteU(addr)
}

func (s *State) StoreWord(addr uint32, v int32) error {
	return s.Mem.StoreWord(addr, v)
}

func (s *State) StoreHalf(addr uint32, v int16) error {
	return s.Mem.StoreHalf(addr, v)
}

func (s *State) StoreByte(addr uint32, v int8) error {
	return s.Mem.StoreByte(addr, v)
}

func (s *State) DumpRegs() string {
	var sb strings.Builder
	for i := 0; i < 32; i++ {
		fmt.Fprintf(&sb, "x%-2d = 0x%08X (%d)", i, uint32(s.Regs[i]), s.Regs[i])
		if (i+1)%4 == 0 {
			sb.WriteString("\n")
		} else {
			sb.WriteString("  ")
		}
	}
	fmt.Fprintf(&sb, "PC  = 0x%08X\n", s.PC)
	return sb.String()
}
