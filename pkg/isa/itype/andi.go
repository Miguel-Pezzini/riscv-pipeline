package itype

import isa "riscv-instruction-encoder/pkg/isa"

type ANDI struct {
	Type
}

func newANDI(t Type) *ANDI {
	inst := &ANDI{Type: t}
	inst.InstructionMeta = isa.InstructionMeta{
		Name:           "ANDI",
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

func (a *ANDI) Execute(state isa.CPUState) error {
	rs1 := state.ReadReg(int(a.Rs1))
	imm := isa.SignExtend12(a.Imm)
	state.WriteReg(int(a.Rd), rs1&imm)
	return nil
}
