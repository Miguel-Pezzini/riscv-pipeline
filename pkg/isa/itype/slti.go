package itype

import isa "riscv-instruction-encoder/pkg/isa"

type SLTI struct {
	Type
}

func newSLTI(t Type) *SLTI {
	inst := &SLTI{Type: t}

	inst.InstructionMeta = isa.InstructionMeta{
		Name:           "SLTI",
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

func (s *SLTI) Execute(state isa.CPUState) error {
	rs1 := state.ReadReg(int(s.Rs1))
	imm := isa.SignExtend12(s.Imm)
	if rs1 < imm {
		state.WriteReg(int(s.Rd), 1)
	} else {
		state.WriteReg(int(s.Rd), 0)
	}
	return nil
}
