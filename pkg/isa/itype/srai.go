package itype

import isa "riscv-instruction-encoder/pkg/isa"

type SRAI struct {
	Type
}

func newSRAI(t Type) *SRAI {
	inst := &SRAI{Type: t}

	inst.InstructionMeta = isa.InstructionMeta{
		Name:           "SRAI",
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
