package rtype

import isa "riscv-instruction-encoder/pkg/isa"

type SRA struct {
	Type
}

func newSRA(t Type) *SRA {
	inst := &SRA{Type: t}

	inst.InstructionMeta = isa.InstructionMeta{
		Name:           "SRA",
		OpCode:         uint32(t.Opcode),
		IsLoad:         false,
		IsStore:        false,
		IsBranch:       false,
		IsJump:         false,
		WritesRegister: true,
		ReadsRegister:  true,
		Rs:             []int{int(t.Rs1), int(t.Rs2)},
		Rd:             isa.IntPtr(int(t.Rd)),
		ProduceStage:   isa.EX,
		ConsumeStage:   isa.ID,
	}

	return inst
}

func (s *SRA) Execute(state isa.CPUState) error {
	rs1 := state.ReadReg(int(s.Rs1))
	shamt := state.ReadReg(int(s.Rs2)) & 0x1F
	state.WriteReg(int(s.Rd), rs1>>shamt) // Go arithmetic right shift on signed
	return nil
}
