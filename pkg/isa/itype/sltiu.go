package itype

import isa "riscv-instruction-encoder/pkg/isa"

type SLTIU struct {
	Type
}

func newSLTIU(t Type) *SLTIU {
	inst := &SLTIU{Type: t}

	inst.InstructionMeta = isa.InstructionMeta{
		Name:           "SLTIU",
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

func (s *SLTIU) Execute(state isa.CPUState) error {
	rs1 := uint32(state.ReadReg(int(s.Rs1)))
	imm := uint32(isa.SignExtend12(s.Imm))
	if rs1 < imm {
		state.WriteReg(int(s.Rd), 1)
	} else {
		state.WriteReg(int(s.Rd), 0)
	}
	return nil
}
