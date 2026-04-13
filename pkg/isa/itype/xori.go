package itype

import isa "riscv-instruction-encoder/pkg/isa"

type XORI struct {
	Type
}

func newXORI(t Type) *XORI {
	inst := &XORI{Type: t}

	inst.InstructionMeta = isa.InstructionMeta{
		Name:           "XORI",
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

func (x *XORI) Execute(state isa.CPUState) error {
	rs1 := state.ReadReg(int(x.Rs1))
	imm := isa.SignExtend12(x.Imm)
	state.WriteReg(int(x.Rd), rs1^imm)
	return nil
}
