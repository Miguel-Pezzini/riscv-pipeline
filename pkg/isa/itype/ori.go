package itype

import isa "riscv-instruction-encoder/pkg/isa"

type ORI struct {
	Type
}

func newORI(t Type) *ORI {
	inst := &ORI{Type: t}
	inst.InstructionMeta = isa.InstructionMeta{
		Name:           "ORI",
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

func (o *ORI) Execute(state isa.CPUState) error {
	rs1 := state.ReadReg(int(o.Rs1))
	imm := isa.SignExtend12(o.Imm)
	state.WriteReg(int(o.Rd), rs1|imm)
	return nil
}
