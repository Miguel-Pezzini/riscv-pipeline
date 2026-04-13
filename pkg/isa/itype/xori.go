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
