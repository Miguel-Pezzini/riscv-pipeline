package itype

import isa "riscv-instruction-encoder/pkg/isa"

type SLLI struct {
	Type
}

func newSLLI(t Type) *SLLI {
	inst := &SLLI{Type: t}

	inst.InstructionMeta = isa.InstructionMeta{
		Name:           "SLLI",
		OpCode:         uint32(t.OpCode),
		IsLoad:         false,
		IsStore:        false,
		IsBranch:       false,
		IsJump:         false,
		WritesRegister: true,
		ReadsRegister:  true,
		Rs:             []int{int(t.Rs1)},
		Rd:             isa.IntPtr(int(t.Rd)),
		ProduceStage:   isa.EX,
		ConsumeStage:   isa.ID,
	}

	return inst
}

func (s *SLLI) Execute(state isa.CPUState) error {
	rs1 := state.ReadReg(int(s.Rs1))
	shamt := s.Imm & 0x1F
	state.WriteReg(int(s.Rd), rs1<<shamt)
	return nil
}
