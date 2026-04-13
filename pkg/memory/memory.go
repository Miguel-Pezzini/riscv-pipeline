package memory

import (
	"fmt"
	"strings"
)

type Memory struct {
	data map[uint32]byte
}

func New() *Memory {
	return &Memory{data: make(map[uint32]byte)}
}

func (m *Memory) LoadByte(addr uint32) (int8, error) {
	return int8(m.data[addr]), nil
}

func (m *Memory) LoadHalf(addr uint32) (int16, error) {
	if addr%2 != 0 {
		return 0, fmt.Errorf("misaligned half-word access at 0x%08X", addr)
	}
	lo := uint16(m.data[addr])
	hi := uint16(m.data[addr+1])
	return int16(hi<<8 | lo), nil
}

func (m *Memory) LoadWord(addr uint32) (int32, error) {
	if addr%4 != 0 {
		return 0, fmt.Errorf("misaligned word access at 0x%08X", addr)
	}
	b0 := uint32(m.data[addr])
	b1 := uint32(m.data[addr+1])
	b2 := uint32(m.data[addr+2])
	b3 := uint32(m.data[addr+3])
	return int32(b3<<24 | b2<<16 | b1<<8 | b0), nil
}

func (m *Memory) LoadByteU(addr uint32) (uint8, error) {
	return m.data[addr], nil
}

func (m *Memory) LoadHalfU(addr uint32) (uint16, error) {
	if addr%2 != 0 {
		return 0, fmt.Errorf("misaligned half-word access at 0x%08X", addr)
	}
	lo := uint16(m.data[addr])
	hi := uint16(m.data[addr+1])
	return hi<<8 | lo, nil
}

func (m *Memory) StoreByte(addr uint32, v int8) error {
	m.data[addr] = byte(v)
	return nil
}

func (m *Memory) StoreHalf(addr uint32, v int16) error {
	if addr%2 != 0 {
		return fmt.Errorf("misaligned half-word access at 0x%08X", addr)
	}
	m.data[addr] = byte(v)
	m.data[addr+1] = byte(v >> 8)
	return nil
}

func (m *Memory) StoreWord(addr uint32, v int32) error {
	if addr%4 != 0 {
		return fmt.Errorf("misaligned word access at 0x%08X", addr)
	}
	m.data[addr] = byte(v)
	m.data[addr+1] = byte(v >> 8)
	m.data[addr+2] = byte(v >> 16)
	m.data[addr+3] = byte(v >> 24)
	return nil
}

func (m *Memory) Dump(from, to uint32) string {
	var sb strings.Builder
	for addr := from; addr < to; addr += 4 {
		w, _ := m.LoadWord(addr)
		fmt.Fprintf(&sb, "0x%08X: 0x%08X\n", addr, uint32(w))
	}
	return sb.String()
}
