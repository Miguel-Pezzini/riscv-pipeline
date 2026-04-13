package utype

import (
	"fmt"
	isa "riscv-instruction-encoder/pkg/isa"
)

type Type struct {
	isa.BaseInstruction
	Opcode uint8  // 7 bits
	Rd     uint8  // 5 bits
	Imm    uint32 // 20 bits
}

func (u *Type) Decode(inst uint32) isa.Instruction {
	u.Opcode = uint8(inst & 0x7F)
	u.Rd = uint8((inst >> 7) & 0x1F)
	u.Imm = uint32(inst>>12) & 0xFFFFF

	switch u.Opcode {
	case 0x37: // LUI
		u.InstructionMeta = isa.InstructionMeta{
			Name:           "LUI",
			OpCode:         uint32(u.Opcode),
			WritesRegister: true,
			ReadsRegister:  false,
			Rs:             []int{},
			Rd:             isa.IntPtr(int(u.Rd)),
			ProduceStage:   isa.EX,
			ConsumeStage:   isa.ID,
		}
	case 0x17: // AUIPC
		u.InstructionMeta = isa.InstructionMeta{
			Name:           "AUIPC",
			OpCode:         uint32(u.Opcode),
			WritesRegister: true,
			ReadsRegister:  false,
			Rs:             []int{},
			Rd:             isa.IntPtr(int(u.Rd)),
			ProduceStage:   isa.EX,
			ConsumeStage:   isa.ID,
		}
	default:
		u.InstructionMeta = isa.InstructionMeta{}
	}

	return u
}

func (u *Type) Execute(state isa.CPUState) error {
	switch u.Opcode {
	case 0x37: // LUI
		state.WriteReg(int(u.Rd), int32(u.Imm<<12))
	case 0x17: // AUIPC
		pc := state.GetPC()
		state.WriteReg(int(u.Rd), int32(pc)+int32(u.Imm<<12))
	}
	return nil
}

func (u *Type) String() string {
	name := u.InstructionMeta.Name
	if name == "" {
		name = fmt.Sprintf("U?(%02X)", u.Opcode)
	}
	return fmt.Sprintf("%s rd=%d, imm=0x%X", name, u.Rd, u.Imm)
}
